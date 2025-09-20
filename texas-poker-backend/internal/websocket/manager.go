// WebSocket连接管理器
// 作用：管理WebSocket连接，实现连接池和消息广播，处理游戏实时通信

package websocket

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

// Message WebSocket消息结构
type Message struct {
	Type string      `json:"type"`
	Data interface{} `json:"data"`
}

// Client WebSocket客户端结构
type Client struct {
	ID     string          // 客户端唯一标识
	UserID int64           // 用户ID
	Conn   *websocket.Conn // WebSocket连接
	Send   chan Message    // 发送消息通道
	Hub    *Hub            // 所属的Hub
}

// Hub WebSocket连接中心
type Hub struct {
	// 注册的客户端
	clients map[*Client]bool

	// 客户端注册通道
	register chan *Client

	// 客户端注销通道
	unregister chan *Client

	// 广播消息通道
	broadcast chan Message

	// 房间消息通道
	roomMessages chan RoomMessage

	// 用户ID到客户端的映射
	userClients map[int64]*Client

	// 互斥锁
	mu sync.RWMutex
}

// RoomMessage 房间消息结构
type RoomMessage struct {
	RoomID  int64   `json:"room_id"`
	Message Message `json:"message"`
}

// Manager WebSocket管理器
type Manager struct {
	Hub      *Hub
	upgrader websocket.Upgrader
}

// NewManager 创建新的WebSocket管理器
func NewManager() *Manager {
	hub := &Hub{
		clients:      make(map[*Client]bool),
		register:     make(chan *Client),
		unregister:   make(chan *Client),
		broadcast:    make(chan Message),
		roomMessages: make(chan RoomMessage),
		userClients:  make(map[int64]*Client),
	}

	// 启动Hub
	go hub.run()

	return &Manager{
		Hub: hub,
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				// 允许所有源的连接（生产环境需要限制）
				return true
			},
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
		},
	}
}

// run 运行Hub主循环
func (h *Hub) run() {
	for {
		select {
		case client := <-h.register:
			// 注册新客户端
			h.mu.Lock()
			h.clients[client] = true
			h.userClients[client.UserID] = client
			h.mu.Unlock()
			
			log.Printf("Client %s (UserID: %d) connected", client.ID, client.UserID)
			
			// 发送连接成功消息
			client.Send <- Message{
				Type: "connected",
				Data: map[string]interface{}{
					"client_id": client.ID,
					"message":   "WebSocket连接成功",
				},
			}

		case client := <-h.unregister:
			// 注销客户端
			h.mu.Lock()
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				delete(h.userClients, client.UserID)
				close(client.Send)
			}
			h.mu.Unlock()
			
			log.Printf("Client %s (UserID: %d) disconnected", client.ID, client.UserID)

		case message := <-h.broadcast:
			// 广播消息给所有客户端
			h.mu.RLock()
			for client := range h.clients {
				select {
				case client.Send <- message:
				default:
					// 客户端发送通道已满，断开连接
					delete(h.clients, client)
					delete(h.userClients, client.UserID)
					close(client.Send)
				}
			}
			h.mu.RUnlock()

		case roomMsg := <-h.roomMessages:
			// 发送房间消息（后续实现房间功能时使用）
			log.Printf("Room message for room %d: %+v", roomMsg.RoomID, roomMsg.Message)
		}
	}
}

// HandleWebSocket 处理WebSocket连接
func (m *Manager) HandleWebSocket(c *gin.Context) {
	// 从上下文获取用户ID（需要认证中间件）
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// 升级HTTP连接为WebSocket
	conn, err := m.upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("WebSocket upgrade failed: %v", err)
		return
	}

	// 创建客户端
	client := &Client{
		ID:     generateClientID(),
		UserID: userID.(int64),
		Conn:   conn,
		Send:   make(chan Message, 256),
		Hub:    m.Hub,
	}

	// 注册客户端
	m.Hub.register <- client

	// 启动goroutine处理客户端
	go client.writePump()
	go client.readPump()
}

// readPump 读取客户端消息
func (c *Client) readPump() {
	defer func() {
		c.Hub.unregister <- c
		c.Conn.Close()
	}()

	// 设置读取超时
	c.Conn.SetReadDeadline(time.Now().Add(60 * time.Second))
	c.Conn.SetPongHandler(func(string) error {
		c.Conn.SetReadDeadline(time.Now().Add(60 * time.Second))
		return nil
	})

	for {
		var msg Message
		err := c.Conn.ReadJSON(&msg)
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("WebSocket error: %v", err)
			}
			break
		}

		// 处理接收到的消息
		c.handleMessage(msg)
	}
}

// writePump 发送消息给客户端
func (c *Client) writePump() {
	ticker := time.NewTicker(54 * time.Second)
	defer func() {
		ticker.Stop()
		c.Conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.Send:
			c.Conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if !ok {
				c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			if err := c.Conn.WriteJSON(message); err != nil {
				log.Printf("WebSocket write error: %v", err)
				return
			}

		case <-ticker.C:
			c.Conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

// handleMessage 处理客户端消息
func (c *Client) handleMessage(msg Message) {
	log.Printf("Received message from client %s: type=%s", c.ID, msg.Type)

	switch msg.Type {
	case "ping":
		// 心跳响应
		c.Send <- Message{
			Type: "pong",
			Data: time.Now().Unix(),
		}

	case "game_action":
		// 游戏操作消息（后续实现游戏逻辑时处理）
		log.Printf("Game action from user %d: %+v", c.UserID, msg.Data)

	default:
		log.Printf("Unknown message type: %s", msg.Type)
	}
}

// BroadcastToAll 广播消息给所有客户端
func (m *Manager) BroadcastToAll(msg Message) {
	m.Hub.broadcast <- msg
}

// SendToUser 发送消息给特定用户
func (m *Manager) SendToUser(userID int64, msg Message) {
	m.Hub.mu.RLock()
	client, exists := m.Hub.userClients[userID]
	m.Hub.mu.RUnlock()

	if exists {
		select {
		case client.Send <- msg:
		default:
			log.Printf("Failed to send message to user %d: channel full", userID)
		}
	}
}

// GetConnectedUsers 获取当前连接的用户列表
func (m *Manager) GetConnectedUsers() []int64 {
	m.Hub.mu.RLock()
	defer m.Hub.mu.RUnlock()

	users := make([]int64, 0, len(m.Hub.userClients))
	for userID := range m.Hub.userClients {
		users = append(users, userID)
	}
	return users
}

// generateClientID 生成客户端ID
func generateClientID() string {
	return fmt.Sprintf("client_%d", time.Now().UnixNano())
} 