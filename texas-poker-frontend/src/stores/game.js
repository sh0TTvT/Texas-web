// 游戏状态管理
// 作用：管理房间列表、当前游戏状态、玩家信息等游戏相关数据

import { defineStore } from 'pinia'
import { roomAPI } from '../api'

export const useGameStore = defineStore('game', {
  state: () => ({
    // 房间列表
    rooms: [],
    isLoadingRooms: false,
    
    // 当前房间信息
    currentRoom: null,
    isInRoom: false,
    
    // 当前游戏状态
    gameState: null,
    players: [],
    playerHand: [],
    communityCards: [],
    pot: 0,
    currentPlayer: null,
    
    // 游戏操作
    availableActions: [],
    currentBet: 0,
    
    // WebSocket连接状态
    wsConnected: false,
    wsConnecting: false
  }),

  getters: {
    // 按筹码级别分组的房间
    roomsByLevel: (state) => {
      const grouped = {
        low: [],
        medium: [],
        high: []
      }
      
      state.rooms.forEach(room => {
        if (grouped[room.chip_level]) {
          grouped[room.chip_level].push(room)
        }
      })
      
      return grouped
    },
    
    // 可加入的房间（未满且状态为等待）
    availableRooms: (state) => {
      return state.rooms.filter(room => 
        room.status === 'waiting' && 
        room.current_players < room.max_players
      )
    },
    
    // 当前玩家信息
    currentPlayerInfo: (state) => {
      return state.players.find(p => p.is_current_user)
    },
    
    // 是否轮到当前玩家操作
    isMyTurn: (state) => {
      const currentPlayerInfo = state.players.find(p => p.is_current_user)
      return currentPlayerInfo && state.currentPlayer === currentPlayerInfo.id
    }
  },

  actions: {
    // 获取房间列表
    async fetchRooms() {
      this.isLoadingRooms = true
      try {
        const response = await roomAPI.getRooms()
        this.rooms = response.data.rooms || []
        return { success: true }
      } catch (error) {
        console.error('获取房间列表失败:', error)
        return { success: false, message: '获取房间列表失败' }
      } finally {
        this.isLoadingRooms = false
      }
    },

    // 创建房间
    async createRoom(roomData) {
      try {
        const response = await roomAPI.createRoom(roomData)
        const newRoom = response.data.room
        
        // 添加到房间列表
        this.rooms.push(newRoom)
        
        return { success: true, room: newRoom }
      } catch (error) {
        console.error('创建房间失败:', error)
        const message = error.response?.data?.error || '创建房间失败'
        return { success: false, message }
      }
    },

    // 加入房间
    async joinRoom(roomId) {
      try {
        const response = await roomAPI.joinRoom(roomId)
        const roomData = response.data
        
        this.currentRoom = roomData.room
        this.isInRoom = true
        
        return { success: true, room: roomData.room }
      } catch (error) {
        console.error('加入房间失败:', error)
        const message = error.response?.data?.error || '加入房间失败'
        return { success: false, message }
      }
    },

    // 离开房间
    async leaveRoom() {
      if (!this.currentRoom) return
      
      try {
        await roomAPI.leaveRoom(this.currentRoom.id)
        
        this.currentRoom = null
        this.isInRoom = false
        this.resetGameState()
        
        return { success: true }
      } catch (error) {
        console.error('离开房间失败:', error)
        return { success: false, message: '离开房间失败' }
      }
    },

    // 重置游戏状态
    resetGameState() {
      this.gameState = null
      this.players = []
      this.playerHand = []
      this.communityCards = []
      this.pot = 0
      this.currentPlayer = null
      this.availableActions = []
      this.currentBet = 0
    },

    // 更新房间信息（通过WebSocket）
    updateRoomInfo(roomData) {
      // 更新房间列表中的房间信息
      const roomIndex = this.rooms.findIndex(room => room.id === roomData.id)
      if (roomIndex !== -1) {
        this.rooms[roomIndex] = { ...this.rooms[roomIndex], ...roomData }
      }
      
      // 如果是当前房间，也更新当前房间信息
      if (this.currentRoom && this.currentRoom.id === roomData.id) {
        this.currentRoom = { ...this.currentRoom, ...roomData }
      }
    },

    // 更新游戏状态（通过WebSocket）
    updateGameState(gameData) {
      this.gameState = gameData.state
      this.players = gameData.players || []
      this.communityCards = gameData.community_cards || []
      this.pot = gameData.pot || 0
      this.currentPlayer = gameData.current_player
      this.currentBet = gameData.current_bet || 0
      
      // 更新当前玩家的手牌
      const currentPlayerInfo = this.players.find(p => p.is_current_user)
      if (currentPlayerInfo) {
        this.playerHand = currentPlayerInfo.cards || []
      }
    },

    // 更新可用操作
    updateAvailableActions(actions) {
      this.availableActions = actions || []
    },

    // WebSocket连接状态管理
    setWSConnected(connected) {
      this.wsConnected = connected
      this.wsConnecting = false
    },

    setWSConnecting(connecting) {
      this.wsConnecting = connecting
    }
  }
}) 