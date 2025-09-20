// WebSocket处理器
// 作用：处理WebSocket连接，游戏实时通信，房间消息广播

package handlers

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	
	"texas-poker-backend/internal/utils"
)

// WebSocketHandler 处理WebSocket连接升级
func (h *Handler) WebSocketHandler(c *gin.Context) {
	// 从查询参数或Header获取token进行认证
	token := c.Query("token")
	if token == "" {
		token = c.GetHeader("Authorization")
		if token != "" && len(token) > 7 && token[:7] == "Bearer " {
			token = token[7:]
		}
	}

	if token == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "缺少认证token",
		})
		return
	}

	// 验证token
	claims, err := utils.ValidateToken(token, h.config.JWTSecret)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "无效的token",
		})
		return
	}

	// 将用户信息设置到上下文
	c.Set("user_id", claims.UserID)
	c.Set("username", claims.Username)

	// 升级WebSocket连接
	h.wsManager.HandleWebSocket(c)
}

// BroadcastToRoom 向房间广播消息
func (h *Handler) BroadcastToRoom(roomID int64, messageType string, data interface{}) {
	// 这里可以根据房间ID获取房间内的所有用户，然后发送消息
	// 暂时简单实现，后续可以扩展
	log.Printf("Broadcasting to room %d: type=%s, data=%+v", roomID, messageType, data)
}

// BroadcastToUser 向特定用户发送消息
func (h *Handler) BroadcastToUser(userID int64, messageType string, data interface{}) {
	message := struct {
		Type string      `json:"type"`
		Data interface{} `json:"data"`
	}{
		Type: messageType,
		Data: data,
	}

	h.wsManager.SendToUser(userID, message)
}

// BroadcastToAll 向所有连接的用户广播消息
func (h *Handler) BroadcastToAll(messageType string, data interface{}) {
	message := struct {
		Type string      `json:"type"`
		Data interface{} `json:"data"`
	}{
		Type: messageType,
		Data: data,
	}

	h.wsManager.BroadcastToAll(message)
} 