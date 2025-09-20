<!--
游戏房间页面
作用：德州扑克游戏的主要游戏界面，包含扑克桌、玩家座位、公共牌、操作区域等
-->

<template>
  <q-layout view="hHh lpR fFf">
    <!-- 顶部工具栏 -->
    <q-header elevated class="bg-dark text-white">
      <q-toolbar>
        <!-- 返回按钮 -->
        <q-btn 
          flat 
          round 
          icon="arrow_back" 
          @click="handleLeaveRoom"
        />
        
        <!-- 房间信息 -->
        <q-toolbar-title class="room-title">
          <div class="room-name">{{ currentRoom?.name || '游戏房间' }}</div>
          <div class="room-info">
            {{ currentRoom?.chip_level }} · {{ currentRoom?.small_blind }}/{{ currentRoom?.big_blind }}
          </div>
        </q-toolbar-title>

        <!-- 房间状态和玩家数 -->
        <q-chip 
          :color="roomStatusColor" 
          :label="roomStatusText"
          text-color="white"
        />
        
        <q-chip 
          color="white" 
          text-color="dark"
          :label="`${playerCount}/${currentRoom?.max_players || 6}`"
          icon="people"
          class="q-ml-sm"
        />

        <!-- 设置菜单 -->
        <q-btn flat round icon="more_vert" class="q-ml-sm">
          <q-menu>
            <q-list style="min-width: 150px">
              <q-item clickable @click="showRoomInfo = true">
                <q-item-section avatar>
                  <q-icon name="info" />
                </q-item-section>
                <q-item-section>房间信息</q-item-section>
              </q-item>
              
              <q-item clickable @click="showSettings = true">
                <q-item-section avatar>
                  <q-icon name="settings" />
                </q-item-section>
                <q-item-section>游戏设置</q-item-section>
              </q-item>
              
              <q-separator />
              
              <q-item clickable @click="handleLeaveRoom">
                <q-item-section avatar>
                  <q-icon name="exit_to_app" />
                </q-item-section>
                <q-item-section>离开房间</q-item-section>
              </q-item>
            </q-list>
          </q-menu>
        </q-btn>
      </q-toolbar>
    </q-header>

    <!-- 游戏主区域 -->
    <q-page-container>
      <q-page class="game-page">
        <div class="poker-table-container">
          <!-- 扑克桌 -->
          <div class="poker-table">
            <!-- 背景装饰 -->
            <div class="table-background"></div>
            
            <!-- 玩家座位 -->
            <div class="player-seats">
              <!-- 顶部座位 -->
              <div class="seats-top">
                <player-seat
                  v-for="seat in topSeats"
                  :key="seat.index"
                  :seat-index="seat.index"
                  :player="seat.player"
                  :is-current-turn="seat.player?.id === currentPlayer"
                  :is-dealer="seat.player?.id === dealerPlayer"
                  :current-bet="seat.player?.current_bet || 0"
                  :player-status="seat.player?.status"
                  :last-action="seat.player?.last_action"
                  :show-cards="gameState === 'showdown' || gameState === 'river'"
                  :show-all-cards="gameState === 'showdown'"
                  position="top"
                  @join-seat="handleJoinSeat"
                />
              </div>

              <!-- 左侧座位 -->
              <div class="seats-left">
                <player-seat
                  v-for="seat in leftSeats"
                  :key="seat.index"
                  :seat-index="seat.index"
                  :player="seat.player"
                  :is-current-turn="seat.player?.id === currentPlayer"
                  :is-dealer="seat.player?.id === dealerPlayer"
                  :current-bet="seat.player?.current_bet || 0"
                  :player-status="seat.player?.status"
                  :last-action="seat.player?.last_action"
                  :show-cards="gameState === 'showdown' || gameState === 'river'"
                  :show-all-cards="gameState === 'showdown'"
                  position="left"
                  @join-seat="handleJoinSeat"
                />
              </div>

              <!-- 右侧座位 -->
              <div class="seats-right">
                <player-seat
                  v-for="seat in rightSeats"
                  :key="seat.index"
                  :seat-index="seat.index"
                  :player="seat.player"
                  :is-current-turn="seat.player?.id === currentPlayer"
                  :is-dealer="seat.player?.id === dealerPlayer"
                  :current-bet="seat.player?.current_bet || 0"
                  :player-status="seat.player?.status"
                  :last-action="seat.player?.last_action"
                  :show-cards="gameState === 'showdown' || gameState === 'river'"
                  :show-all-cards="gameState === 'showdown'"
                  position="right"
                  @join-seat="handleJoinSeat"
                />
              </div>

              <!-- 底部座位（当前用户） -->
              <div class="seats-bottom">
                <player-seat
                  v-for="seat in bottomSeats"
                  :key="seat.index"
                  :seat-index="seat.index"
                  :player="seat.player"
                  :is-current-turn="seat.player?.id === currentPlayer"
                  :is-dealer="seat.player?.id === dealerPlayer"
                  :current-bet="seat.player?.current_bet || 0"
                  :player-status="seat.player?.status"
                  :last-action="seat.player?.last_action"
                  :show-cards="true"
                  :show-all-cards="gameState === 'showdown'"
                  position="bottom"
                  @join-seat="handleJoinSeat"
                />
              </div>
            </div>

            <!-- 桌面中央区域 -->
            <div class="table-center">
              <!-- 公共牌区域 -->
              <div class="community-cards">
                <div class="cards-title">公共牌</div>
                <div class="cards-area">
                  <poker-card
                    v-for="(card, index) in communityCards"
                    :key="`community-${index}`"
                    :suit="card.suit"
                    :rank="card.rank"
                    size="normal"
                    class="community-card"
                  />
                  
                  <!-- 占位卡片 -->
                  <poker-card
                    v-for="index in (5 - communityCards.length)"
                    :key="`placeholder-${index}`"
                    :is-hidden="true"
                    size="normal"
                    class="community-card placeholder"
                  />
                </div>
              </div>

              <!-- 底池信息 -->
              <div class="pot-info">
                <q-chip 
                  color="orange" 
                  text-color="white" 
                  size="lg"
                  icon="account_balance"
                >
                  底池: {{ formatChips(pot) }}
                </q-chip>
              </div>

              <!-- 游戏状态信息 -->
              <div class="game-status">
                <q-chip 
                  :color="gameStateColor" 
                  :label="gameStateText"
                  size="md"
                />
              </div>
            </div>
          </div>
        </div>

        <!-- 操作区域 -->
        <div class="action-area">
          <game-actions
            :is-my-turn="isMyTurn"
            :available-actions="availableActions"
            :current-bet="currentBet"
            :player-chips="authStore.userChips"
            :pot-size="pot"
            :big-blind="currentRoom?.big_blind || 10"
            :small-blind="currentRoom?.small_blind || 5"
            :game-in-progress="gameInProgress"
            @action="handlePlayerAction"
          />
        </div>
      </q-page>
    </q-page-container>

    <!-- 房间信息对话框 -->
    <q-dialog v-model="showRoomInfo">
      <q-card style="min-width: 400px">
        <q-card-section>
          <div class="text-h6">房间信息</div>
        </q-card-section>

        <q-card-section class="q-pt-none">
          <q-list dense>
            <q-item>
              <q-item-section avatar>
                <q-icon name="title" />
              </q-item-section>
              <q-item-section>
                <q-item-label>房间名称</q-item-label>
                <q-item-label caption>{{ currentRoom?.name }}</q-item-label>
              </q-item-section>
            </q-item>

            <q-item>
              <q-item-section avatar>
                <q-icon name="star" />
              </q-item-section>
              <q-item-section>
                <q-item-label>筹码级别</q-item-label>
                <q-item-label caption>{{ currentRoom?.chip_level }}</q-item-label>
              </q-item-section>
            </q-item>

            <q-item>
              <q-item-section avatar>
                <q-icon name="visibility" />
              </q-item-section>
              <q-item-section>
                <q-item-label>盲注</q-item-label>
                <q-item-label caption>
                  {{ currentRoom?.small_blind }}/{{ currentRoom?.big_blind }}
                </q-item-label>
              </q-item-section>
            </q-item>

            <q-item>
              <q-item-section avatar>
                <q-icon name="people" />
              </q-item-section>
              <q-item-section>
                <q-item-label>玩家数量</q-item-label>
                <q-item-label caption>
                  {{ playerCount }}/{{ currentRoom?.max_players }}
                </q-item-label>
              </q-item-section>
            </q-item>
          </q-list>
        </q-card-section>

        <q-card-actions align="right">
          <q-btn flat label="关闭" color="primary" v-close-popup />
        </q-card-actions>
      </q-card>
    </q-dialog>

    <!-- 游戏设置对话框 -->
    <q-dialog v-model="showSettings">
      <q-card style="min-width: 400px">
        <q-card-section>
          <div class="text-h6">游戏设置</div>
        </q-card-section>

        <q-card-section class="q-pt-none">
          <q-list>
            <q-item>
              <q-item-section>
                <q-item-label>音效</q-item-label>
              </q-item-section>
              <q-item-section side>
                <q-toggle v-model="gameSettings.sound" color="primary" />
              </q-item-section>
            </q-item>

            <q-item>
              <q-item-section>
                <q-item-label>自动操作提示</q-item-label>
              </q-item-section>
              <q-item-section side>
                <q-toggle v-model="gameSettings.autoHints" color="primary" />
              </q-item-section>
            </q-item>

            <q-item>
              <q-item-section>
                <q-item-label>快速操作</q-item-label>
              </q-item-section>
              <q-item-section side>
                <q-toggle v-model="gameSettings.quickActions" color="primary" />
              </q-item-section>
            </q-item>
          </q-list>
        </q-card-section>

        <q-card-actions align="right">
          <q-btn flat label="关闭" color="primary" v-close-popup />
        </q-card-actions>
      </q-card>
    </q-dialog>
  </q-layout>
</template>

<script>
import { defineComponent, ref, computed, onMounted, onUnmounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useQuasar } from 'quasar'
import { useAuthStore } from '../../stores/auth'
import { useGameStore } from '../../stores/game'
import PlayerSeat from '../../components/PlayerSeat.vue'
import PokerCard from '../../components/PokerCard.vue'
import GameActions from '../../components/GameActions.vue'
import wsClient from '../../utils/websocket'

export default defineComponent({
  name: 'GameRoomPage',
  
  components: {
    PlayerSeat,
    PokerCard,
    GameActions
  },
  
  setup() {
    const router = useRouter()
    const route = useRoute()
    const $q = useQuasar()
    const authStore = useAuthStore()
    const gameStore = useGameStore()
    
    // 响应式数据
    const showRoomInfo = ref(false)
    const showSettings = ref(false)
    const gameSettings = ref({
      sound: true,
      autoHints: true,
      quickActions: true
    })
    
    // 计算属性
    const roomId = computed(() => route.params.roomId)
    const currentRoom = computed(() => gameStore.currentRoom)
    const players = computed(() => gameStore.players || [])
    const playerCount = computed(() => players.value.length)
    const gameState = computed(() => gameStore.gameState)
    const communityCards = computed(() => gameStore.communityCards || [])
    const pot = computed(() => gameStore.pot || 0)
    const currentPlayer = computed(() => gameStore.currentPlayer)
    const dealerPlayer = computed(() => {
      return players.value.find(p => p.is_dealer)?.id
    })
    const availableActions = computed(() => gameStore.availableActions || [])
    const currentBet = computed(() => gameStore.currentBet || 0)
    const isMyTurn = computed(() => gameStore.isMyTurn)
    const gameInProgress = computed(() => ['preflop', 'flop', 'turn', 'river', 'showdown'].includes(gameState.value))
    
    // 座位布局
    const allSeats = computed(() => {
      const seats = Array.from({ length: 6 }, (_, index) => ({
        index,
        player: players.value.find(p => p.seat_index === index)
      }))
      return seats
    })
    
    const topSeats = computed(() => allSeats.value.slice(1, 3))      // 座位 1, 2
    const leftSeats = computed(() => [allSeats.value[0]])           // 座位 0
    const rightSeats = computed(() => [allSeats.value[3]])          // 座位 3
    const bottomSeats = computed(() => allSeats.value.slice(4, 6))   // 座位 4, 5
    
    const roomStatusColor = computed(() => {
      const colors = {
        waiting: 'green',
        playing: 'orange',
        paused: 'blue'
      }
      return colors[currentRoom.value?.status] || 'grey'
    })
    
    const roomStatusText = computed(() => {
      const texts = {
        waiting: '等待中',
        playing: '游戏中',
        paused: '暂停'
      }
      return texts[currentRoom.value?.status] || '未知'
    })
    
    const gameStateColor = computed(() => {
      const colors = {
        waiting: 'grey',
        preflop: 'blue',
        flop: 'green',
        turn: 'orange',
        river: 'red',
        showdown: 'purple'
      }
      return colors[gameState.value] || 'grey'
    })
    
    const gameStateText = computed(() => {
      const texts = {
        waiting: '等待开始',
        preflop: '翻牌前',
        flop: '翻牌',
        turn: '转牌',
        river: '河牌',
        showdown: '摊牌'
      }
      return texts[gameState.value] || '未知'
    })
    
    // 方法
    const formatChips = (chips) => {
      if (chips >= 1000000) {
        return (chips / 1000000).toFixed(1) + 'M'
      } else if (chips >= 1000) {
        return (chips / 1000).toFixed(1) + 'K'
      }
      return chips.toString()
    }
    
    const handleJoinSeat = (seatIndex) => {
      if (wsClient.isConnected) {
        wsClient.send({
          type: 'join_seat',
          seat_index: seatIndex
        })
      } else {
        $q.notify({
          type: 'negative',
          message: 'WebSocket连接未建立',
          position: 'top'
        })
      }
    }
    
    const handlePlayerAction = (action) => {
      if (!wsClient.isConnected) {
        $q.notify({
          type: 'negative',
          message: 'WebSocket连接未建立',
          position: 'top'
        })
        return
      }
      
      // 发送玩家操作到服务器
      const success = wsClient.playerAction(action.type, action.amount)
      
      if (!success) {
        $q.notify({
          type: 'negative',
          message: '操作发送失败，请重试',
          position: 'top'
        })
      }
    }
    
    const handleLeaveRoom = () => {
      $q.dialog({
        title: '确认',
        message: '确定要离开房间吗？',
        cancel: true,
        persistent: true
      }).onOk(async () => {
        const result = await gameStore.leaveRoom()
        if (result.success) {
          $q.notify({
            type: 'positive',
            message: '已离开房间',
            position: 'top'
          })
          router.push('/lobby')
        } else {
          $q.notify({
            type: 'negative',
            message: '离开房间失败',
            position: 'top'
          })
        }
      })
    }
    
    // WebSocket事件监听器
    const setupWebSocketListeners = () => {
      // 监听游戏状态更新
      wsClient.on('game_state_update', (data) => {
        console.log('Game state updated:', data)
      })
      
      // 监听房间更新
      wsClient.on('room_update', (data) => {
        console.log('Room updated:', data)
      })
      
      // 监听玩家操作
      wsClient.on('player_action', (data) => {
        console.log('Player action:', data)
      })
      
      // 监听错误消息
      wsClient.on('error_message', (data) => {
        console.error('WebSocket error:', data)
      })
      
      // 监听连接状态变化
      wsClient.on('connected', () => {
        console.log('WebSocket connected')
        // 连接成功后加入房间
        wsClient.joinRoom(roomId.value)
      })
      
      wsClient.on('closed', () => {
        console.log('WebSocket disconnected')
      })
    }
    
    // 清理WebSocket监听器
    const cleanupWebSocketListeners = () => {
      wsClient.off('game_state_update')
      wsClient.off('room_update')
      wsClient.off('player_action')
      wsClient.off('error_message')
      wsClient.off('connected')
      wsClient.off('closed')
    }
    
    // 生命周期
    onMounted(async () => {
      // 检查用户是否已登录
      if (!authStore.isAuthenticated) {
        router.push('/auth/login')
        return
      }
      
      // 设置WebSocket监听器
      setupWebSocketListeners()
      
      // 建立WebSocket连接
      try {
        await wsClient.connect(authStore.token)
        
        // 如果不在房间中，尝试加入房间
        if (!gameStore.isInRoom) {
          const result = await gameStore.joinRoom(roomId.value)
          if (!result.success) {
            $q.notify({
              type: 'negative',
              message: result.message,
              position: 'top'
            })
            router.push('/lobby')
            return
          }
        }
      } catch (error) {
        console.error('WebSocket connection failed:', error)
        $q.notify({
          type: 'negative',
          message: 'WebSocket连接失败',
          position: 'top'
        })
      }
    })
    
    onUnmounted(() => {
      // 清理WebSocket监听器和连接
      cleanupWebSocketListeners()
      
      // 如果还在房间中，先离开房间
      if (gameStore.isInRoom) {
        wsClient.leaveRoom()
      }
      
      // 断开WebSocket连接
      wsClient.disconnect()
    })
    
    return {
      authStore,
      gameStore,
      showRoomInfo,
      showSettings,
      gameSettings,
      currentRoom,
      players,
      playerCount,
      gameState,
      communityCards,
      pot,
      currentPlayer,
      dealerPlayer,
      availableActions,
      currentBet,
      isMyTurn,
      gameInProgress,
      topSeats,
      leftSeats,
      rightSeats,
      bottomSeats,
      roomStatusColor,
      roomStatusText,
      gameStateColor,
      gameStateText,
      formatChips,
      handleJoinSeat,
      handlePlayerAction,
      handleLeaveRoom
    }
  }
})
</script>

<style scoped lang="scss">
.game-page {
  background: linear-gradient(135deg, #0f3460 0%, #1e3c72 50%, #2a5298 100%);
  min-height: 100vh;
  display: flex;
  flex-direction: column;
}

.room-title {
  .room-name {
    font-size: 1.1rem;
    font-weight: 600;
  }
  
  .room-info {
    font-size: 0.8rem;
    opacity: 0.8;
  }
}

.poker-table-container {
  flex: 1;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 1rem;
}

.poker-table {
  position: relative;
  width: 100%;
  max-width: 800px;
  height: 500px;
  border-radius: 200px;
  background: linear-gradient(135deg, #2d5a27 0%, #1a3918 100%);
  border: 8px solid #8b4513;
  box-shadow: 
    inset 0 0 50px rgba(0, 0, 0, 0.3),
    0 10px 30px rgba(0, 0, 0, 0.4);
  overflow: visible;
}

.table-background {
  position: absolute;
  top: 20px;
  left: 20px;
  right: 20px;
  bottom: 20px;
  border-radius: 180px;
  background: 
    radial-gradient(ellipse at center, rgba(255, 255, 255, 0.1) 0%, transparent 70%),
    linear-gradient(45deg, rgba(0, 0, 0, 0.1) 25%, transparent 25%),
    linear-gradient(-45deg, rgba(0, 0, 0, 0.1) 25%, transparent 25%);
  background-size: 100% 100%, 20px 20px, 20px 20px;
}

.player-seats {
  position: relative;
  width: 100%;
  height: 100%;
}

.seats-top {
  position: absolute;
  top: -60px;
  left: 50%;
  transform: translateX(-50%);
  display: flex;
  gap: 2rem;
}

.seats-bottom {
  position: absolute;
  bottom: -60px;
  left: 50%;
  transform: translateX(-50%);
  display: flex;
  gap: 2rem;
}

.seats-left {
  position: absolute;
  left: -80px;
  top: 50%;
  transform: translateY(-50%);
  display: flex;
  flex-direction: column;
  gap: 1rem;
}

.seats-right {
  position: absolute;
  right: -80px;
  top: 50%;
  transform: translateY(-50%);
  display: flex;
  flex-direction: column;
  gap: 1rem;
}

.table-center {
  position: absolute;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 1rem;
}

.community-cards {
  text-align: center;
  
  .cards-title {
    color: white;
    font-size: 0.9rem;
    font-weight: 500;
    margin-bottom: 0.5rem;
    opacity: 0.8;
  }
  
  .cards-area {
    display: flex;
    gap: 0.5rem;
    
    .community-card {
      &.placeholder {
        opacity: 0.3;
      }
    }
  }
}

.pot-info {
  text-align: center;
}

.game-status {
  text-align: center;
}

.action-area {
  padding: 1rem;
  background: rgba(255, 255, 255, 0.1);
  backdrop-filter: blur(10px);
}

// 响应式设计
@media (max-width: 1024px) {
  .poker-table {
    max-width: 700px;
    height: 450px;
  }
  
  .seats-left,
  .seats-right {
    left: -60px;
    right: -60px;
  }
}

@media (max-width: 768px) {
  .poker-table-container {
    padding: 0.5rem;
  }
  
  .poker-table {
    max-width: 600px;
    height: 400px;
  }
  
  .seats-top,
  .seats-bottom {
    gap: 1rem;
  }
  
  .seats-left,
  .seats-right {
    left: -40px;
    right: -40px;
  }
  
  .community-cards .cards-area {
    gap: 0.25rem;
  }
}

@media (max-width: 480px) {
  .poker-table {
    max-width: 350px;
    height: 250px;
    border-radius: 120px;
  }
  
  .table-background {
    border-radius: 100px;
  }
  
  .seats-top,
  .seats-bottom {
    top: -50px;
    bottom: -50px;
    gap: 0.5rem;
  }
  
  .seats-left,
  .seats-right {
    left: -30px;
    right: -30px;
  }
  
  .table-center {
    gap: 0.5rem;
  }
  
  .community-cards .cards-area {
    flex-wrap: wrap;
    justify-content: center;
    max-width: 200px;
  }
}

// 动画效果
.poker-table {
  animation: tableGlow 3s ease-in-out infinite alternate;
}

@keyframes tableGlow {
  0% {
    box-shadow: 
      inset 0 0 50px rgba(0, 0, 0, 0.3),
      0 10px 30px rgba(0, 0, 0, 0.4);
  }
  100% {
    box-shadow: 
      inset 0 0 50px rgba(0, 0, 0, 0.3),
      0 10px 30px rgba(0, 0, 0, 0.4),
      0 0 20px rgba(255, 215, 0, 0.2);
  }
}

.community-card {
  animation: cardDeal 0.5s ease-out;
}

@keyframes cardDeal {
  from {
    transform: translateY(-50px) rotate(180deg);
    opacity: 0;
  }
  to {
    transform: translateY(0) rotate(0deg);
    opacity: 1;
  }
}
</style> 