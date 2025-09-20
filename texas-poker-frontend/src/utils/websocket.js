// WebSocket客户端管理工具
// 作用：管理与后端的WebSocket连接，处理消息收发、重连机制、事件监听

import { useGameStore } from '../stores/game'
import { useAuthStore } from '../stores/auth'
import { Notify } from 'quasar'

class WebSocketClient {
  constructor() {
    this.ws = null
    this.url = ''
    this.token = ''
    this.reconnectAttempts = 0
    this.maxReconnectAttempts = 5
    this.reconnectDelay = 1000
    this.heartbeatInterval = null
    this.isConnecting = false
    this.isManualClose = false
    this.messageQueue = []
    
    // 事件监听器
    this.eventListeners = new Map()
    
    // 绑定方法
    this.connect = this.connect.bind(this)
    this.disconnect = this.disconnect.bind(this)
    this.send = this.send.bind(this)
    this.onMessage = this.onMessage.bind(this)
    this.onError = this.onError.bind(this)
    this.onClose = this.onClose.bind(this)
    this.onOpen = this.onOpen.bind(this)
  }

  // 连接WebSocket
  connect(token) {
    if (this.isConnecting || (this.ws && this.ws.readyState === WebSocket.OPEN)) {
      return Promise.resolve()
    }

    this.token = token
    this.isConnecting = true
    this.isManualClose = false
    
    // 构建WebSocket URL
    const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:'
    const host = import.meta.env.VITE_WS_HOST || window.location.host.replace(':3000', ':8080')
    this.url = `${protocol}//${host}/ws?token=${encodeURIComponent(token)}`

    return new Promise((resolve, reject) => {
      try {
        this.ws = new WebSocket(this.url)
        
        this.ws.onopen = (event) => {
          this.onOpen(event)
          resolve()
        }
        
        this.ws.onmessage = this.onMessage
        this.ws.onerror = (event) => {
          this.onError(event)
          reject(event)
        }
        this.ws.onclose = this.onClose
        
        // 连接超时处理
        setTimeout(() => {
          if (this.isConnecting) {
            this.isConnecting = false
            reject(new Error('WebSocket connection timeout'))
          }
        }, 10000)
        
      } catch (error) {
        this.isConnecting = false
        reject(error)
      }
    })
  }

  // 断开连接
  disconnect() {
    this.isManualClose = true
    this.clearHeartbeat()
    
    if (this.ws) {
      this.ws.close(1000, 'Manual disconnect')
      this.ws = null
    }
    
    this.reconnectAttempts = 0
    this.isConnecting = false
    
    // 更新游戏store状态
    const gameStore = useGameStore()
    gameStore.setWSConnected(false)
  }

  // 发送消息
  send(message) {
    if (!message) return false

    const messageStr = typeof message === 'string' ? message : JSON.stringify(message)
    
    if (this.ws && this.ws.readyState === WebSocket.OPEN) {
      try {
        this.ws.send(messageStr)
        return true
      } catch (error) {
        console.error('Failed to send WebSocket message:', error)
        this.messageQueue.push(messageStr)
        return false
      }
    } else {
      // 连接未就绪，将消息加入队列
      this.messageQueue.push(messageStr)
      return false
    }
  }

  // 处理连接打开
  onOpen(event) {
    console.log('WebSocket connected:', event)
    this.isConnecting = false
    this.reconnectAttempts = 0
    
    // 发送队列中的消息
    this.flushMessageQueue()
    
    // 启动心跳
    this.startHeartbeat()
    
    // 更新游戏store状态
    const gameStore = useGameStore()
    gameStore.setWSConnected(true)
    
    // 触发连接成功事件
    this.emit('connected', event)
    
    Notify.create({
      type: 'positive',
      message: 'WebSocket连接成功',
      timeout: 2000,
      position: 'top'
    })
  }

  // 处理消息接收
  onMessage(event) {
    try {
      const data = JSON.parse(event.data)
      console.log('WebSocket message received:', data)
      
      // 根据消息类型处理
      switch (data.type) {
        case 'game_state_update':
          this.handleGameStateUpdate(data.payload)
          break
          
        case 'room_update':
          this.handleRoomUpdate(data.payload)
          break
          
        case 'player_action':
          this.handlePlayerAction(data.payload)
          break
          
        case 'error':
          this.handleError(data.payload)
          break
          
        case 'notification':
          this.handleNotification(data.payload)
          break
          
        case 'heartbeat':
          // 心跳响应，不需要特殊处理
          break
          
        default:
          console.warn('Unknown message type:', data.type)
      }
      
      // 触发通用消息事件
      this.emit('message', data)
      
    } catch (error) {
      console.error('Failed to parse WebSocket message:', error, event.data)
    }
  }

  // 处理连接错误
  onError(event) {
    console.error('WebSocket error:', event)
    this.isConnecting = false
    
    // 触发错误事件
    this.emit('error', event)
  }

  // 处理连接关闭
  onClose(event) {
    console.log('WebSocket closed:', event.code, event.reason)
    this.isConnecting = false
    this.clearHeartbeat()
    
    // 更新游戏store状态
    const gameStore = useGameStore()
    gameStore.setWSConnected(false)
    
    // 触发关闭事件
    this.emit('closed', event)
    
    // 如果不是手动关闭，尝试重连
    if (!this.isManualClose && this.reconnectAttempts < this.maxReconnectAttempts) {
      this.attemptReconnect()
    } else if (this.reconnectAttempts >= this.maxReconnectAttempts) {
      Notify.create({
        type: 'negative',
        message: 'WebSocket连接失败，请刷新页面重试',
        timeout: 5000,
        position: 'top'
      })
    }
  }

  // 尝试重连
  attemptReconnect() {
    this.reconnectAttempts++
    const delay = this.reconnectDelay * Math.pow(2, this.reconnectAttempts - 1) // 指数退避
    
    console.log(`Attempting to reconnect (${this.reconnectAttempts}/${this.maxReconnectAttempts}) in ${delay}ms`)
    
    setTimeout(() => {
      if (!this.isManualClose && this.token) {
        this.connect(this.token).catch(error => {
          console.error('Reconnection failed:', error)
        })
      }
    }, delay)
  }

  // 发送队列中的消息
  flushMessageQueue() {
    while (this.messageQueue.length > 0) {
      const message = this.messageQueue.shift()
      if (this.ws && this.ws.readyState === WebSocket.OPEN) {
        this.ws.send(message)
      } else {
        // 连接又断开了，重新加入队列
        this.messageQueue.unshift(message)
        break
      }
    }
  }

  // 启动心跳
  startHeartbeat() {
    this.clearHeartbeat()
    this.heartbeatInterval = setInterval(() => {
      if (this.ws && this.ws.readyState === WebSocket.OPEN) {
        this.send({ type: 'heartbeat', timestamp: Date.now() })
      }
    }, 30000) // 30秒心跳
  }

  // 清除心跳
  clearHeartbeat() {
    if (this.heartbeatInterval) {
      clearInterval(this.heartbeatInterval)
      this.heartbeatInterval = null
    }
  }

  // 处理游戏状态更新
  handleGameStateUpdate(payload) {
    const gameStore = useGameStore()
    gameStore.updateGameState(payload)
    
    // 触发特定事件
    this.emit('game_state_update', payload)
  }

  // 处理房间更新
  handleRoomUpdate(payload) {
    const gameStore = useGameStore()
    gameStore.updateRoomInfo(payload)
    
    // 触发特定事件
    this.emit('room_update', payload)
  }

  // 处理玩家操作
  handlePlayerAction(payload) {
    // 触发特定事件
    this.emit('player_action', payload)
    
    // 显示操作通知
    if (payload.player_name && payload.action) {
      const actionTexts = {
        call: '跟注',
        raise: '加注',
        fold: '弃牌',
        check: '看牌',
        all_in: '全押'
      }
      
      const actionText = actionTexts[payload.action] || payload.action
      
      Notify.create({
        message: `${payload.player_name} ${actionText}`,
        color: 'info',
        timeout: 2000,
        position: 'top-right'
      })
    }
  }

  // 处理错误消息
  handleError(payload) {
    console.error('WebSocket error message:', payload)
    
    Notify.create({
      type: 'negative',
      message: payload.message || '发生错误',
      timeout: 5000,
      position: 'top'
    })
    
    // 触发特定事件
    this.emit('error_message', payload)
  }

  // 处理通知消息
  handleNotification(payload) {
    const notifyConfig = {
      message: payload.message,
      timeout: payload.timeout || 3000,
      position: payload.position || 'top'
    }
    
    // 根据通知类型设置样式
    if (payload.type) {
      notifyConfig.type = payload.type
    }
    
    Notify.create(notifyConfig)
    
    // 触发特定事件
    this.emit('notification', payload)
  }

  // 游戏操作方法
  joinRoom(roomId) {
    return this.send({
      type: 'join_room',
      room_id: roomId
    })
  }

  leaveRoom() {
    return this.send({
      type: 'leave_room'
    })
  }

  playerAction(action, amount = null) {
    const message = {
      type: 'player_action',
      action: action
    }
    
    if (amount !== null) {
      message.amount = amount
    }
    
    return this.send(message)
  }

  // 事件监听器管理
  on(event, callback) {
    if (!this.eventListeners.has(event)) {
      this.eventListeners.set(event, [])
    }
    this.eventListeners.get(event).push(callback)
  }

  off(event, callback) {
    if (this.eventListeners.has(event)) {
      const listeners = this.eventListeners.get(event)
      const index = listeners.indexOf(callback)
      if (index > -1) {
        listeners.splice(index, 1)
      }
    }
  }

  emit(event, data) {
    if (this.eventListeners.has(event)) {
      this.eventListeners.get(event).forEach(callback => {
        try {
          callback(data)
        } catch (error) {
          console.error('Error in event listener:', error)
        }
      })
    }
  }

  // 获取连接状态
  get readyState() {
    return this.ws ? this.ws.readyState : WebSocket.CLOSED
  }

  get isConnected() {
    return this.ws && this.ws.readyState === WebSocket.OPEN
  }
}

// 创建单例实例
const wsClient = new WebSocketClient()

export default wsClient 