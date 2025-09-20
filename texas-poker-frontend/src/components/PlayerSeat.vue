<!--
玩家座位组件
作用：显示单个玩家的座位信息，包括头像、筹码、手牌、状态等
-->

<template>
  <div 
    class="player-seat"
    :class="seatClasses"
  >
    <!-- 空座位 -->
    <div v-if="!player" class="empty-seat">
      <q-btn
        round
        color="primary"
        icon="add"
        size="lg"
        @click="$emit('join-seat', seatIndex)"
        class="join-btn"
      >
        <q-tooltip>加入游戏</q-tooltip>
      </q-btn>
      <div class="seat-label">空位</div>
    </div>

    <!-- 有玩家的座位 -->
    <div v-else class="occupied-seat">
      <!-- 玩家头像和信息 -->
      <div class="player-info">
        <q-avatar 
          :size="avatarSize"
          :color="player.is_current_user ? 'primary' : 'grey-7'"
          text-color="white"
          class="player-avatar"
        >
          <q-icon name="person" :size="avatarIconSize" />
          
          <!-- 当前玩家指示器 -->
          <q-badge
            v-if="isCurrentTurn"
            color="warning"
            floating
            rounded
            class="turn-indicator"
          >
            <q-icon name="play_arrow" size="xs" />
          </q-badge>
          
          <!-- 庄家标识 -->
          <q-badge
            v-if="isDealer"
            color="orange"
            floating
            rounded
            class="dealer-badge"
          >
            D
          </q-badge>
        </q-avatar>
        
        <div class="player-details">
          <div class="player-name">
            {{ player.username }}
            <q-chip
              v-if="player.is_current_user"
              dense
              color="primary"
              text-color="white"
              label="你"
              size="sm"
            />
          </div>
          
          <div class="player-chips">
            <q-icon name="account_balance" size="xs" />
            {{ formatChips(player.chips) }}
          </div>
        </div>
      </div>

      <!-- 玩家手牌 -->
      <div v-if="showCards" class="player-cards">
        <poker-card
          v-for="(card, index) in playerCards"
          :key="`card-${index}`"
          :suit="card.suit"
          :rank="card.rank"
          :is-hidden="!player.is_current_user && !showAllCards"
          size="small"
          class="hand-card"
        />
      </div>

      <!-- 玩家当前下注 -->
      <div v-if="currentBet > 0" class="player-bet">
        <q-chip
          color="orange"
          text-color="white"
          :label="formatChips(currentBet)"
          icon="casino"
          dense
        />
      </div>

      <!-- 玩家状态指示 -->
      <div v-if="playerStatus" class="player-status">
        <q-chip
          :color="statusColor"
          :label="statusText"
          dense
          size="sm"
        />
      </div>

      <!-- 玩家操作历史 -->
      <div v-if="lastAction" class="last-action">
        <q-chip
          color="grey-6"
          text-color="white"
          :label="lastAction"
          dense
          size="sm"
          class="action-chip"
        />
      </div>

      <!-- 连接状态 -->
      <div v-if="!player.is_connected" class="connection-status">
        <q-icon name="wifi_off" color="red" size="sm">
          <q-tooltip>玩家已断线</q-tooltip>
        </q-icon>
      </div>
    </div>
  </div>
</template>

<script>
import { defineComponent, computed } from 'vue'
import PokerCard from './PokerCard.vue'

export default defineComponent({
  name: 'PlayerSeat',
  
  components: {
    PokerCard
  },
  
  props: {
    // 座位索引
    seatIndex: {
      type: Number,
      required: true
    },
    // 玩家信息
    player: {
      type: Object,
      default: null
    },
    // 是否是当前轮到的玩家
    isCurrentTurn: {
      type: Boolean,
      default: false
    },
    // 是否是庄家
    isDealer: {
      type: Boolean,
      default: false
    },
    // 当前下注额
    currentBet: {
      type: Number,
      default: 0
    },
    // 玩家状态
    playerStatus: {
      type: String,
      default: null // active, folded, all_in, waiting
    },
    // 最后操作
    lastAction: {
      type: String,
      default: null // call, raise, fold, check, all_in
    },
    // 是否显示手牌
    showCards: {
      type: Boolean,
      default: false
    },
    // 是否显示所有玩家的手牌（摊牌阶段）
    showAllCards: {
      type: Boolean,
      default: false
    },
    // 座位位置（用于样式）
    position: {
      type: String,
      default: 'bottom', // top, bottom, left, right
      validator: val => ['top', 'bottom', 'left', 'right'].includes(val)
    }
  },
  
  emits: ['join-seat'],
  
  setup(props) {
    // 计算属性
    const seatClasses = computed(() => {
      return {
        [`seat-${props.position}`]: true,
        'seat-empty': !props.player,
        'seat-occupied': !!props.player,
        'seat-current-turn': props.isCurrentTurn,
        'seat-current-user': props.player?.is_current_user,
        'seat-disconnected': props.player && !props.player.is_connected
      }
    })
    
    const avatarSize = computed(() => {
      return props.position === 'bottom' ? '64px' : '48px'
    })
    
    const avatarIconSize = computed(() => {
      return props.position === 'bottom' ? 'md' : 'sm'
    })
    
    const playerCards = computed(() => {
      if (!props.player?.cards) return []
      
      // 解析卡牌字符串为对象
      return props.player.cards.map(cardStr => {
        if (typeof cardStr === 'string') {
          // 假设格式为 "AS" (Ace of Spades)
          const rank = cardStr.slice(0, -1)
          const suitChar = cardStr.slice(-1)
          
          const suitMap = {
            'S': 'spades',
            'H': 'hearts', 
            'D': 'diamonds',
            'C': 'clubs'
          }
          
          return {
            rank,
            suit: suitMap[suitChar] || 'spades'
          }
        }
        return cardStr
      })
    })
    
    const statusColor = computed(() => {
      const colors = {
        active: 'green',
        folded: 'grey',
        all_in: 'red',
        waiting: 'orange'
      }
      return colors[props.playerStatus] || 'grey'
    })
    
    const statusText = computed(() => {
      const texts = {
        active: '游戏中',
        folded: '已弃牌',
        all_in: '全押',
        waiting: '等待中'
      }
      return texts[props.playerStatus] || ''
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
    
    return {
      seatClasses,
      avatarSize,
      avatarIconSize,
      playerCards,
      statusColor,
      statusText,
      formatChips
    }
  }
})
</script>

<style scoped lang="scss">
.player-seat {
  position: relative;
  display: flex;
  flex-direction: column;
  align-items: center;
  min-width: 120px;
  padding: 0.5rem;
  transition: all 0.3s ease;
  
  &.seat-current-turn {
    .player-avatar {
      animation: pulse 2s infinite;
    }
  }
  
  &.seat-current-user {
    .occupied-seat {
      border: 2px solid #1976d2;
      border-radius: 12px;
      background: rgba(25, 118, 210, 0.1);
    }
  }
  
  &.seat-disconnected {
    opacity: 0.6;
  }
}

.empty-seat {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 0.5rem;
  
  .join-btn {
    opacity: 0.7;
    transition: opacity 0.3s ease;
    
    &:hover {
      opacity: 1;
    }
  }
  
  .seat-label {
    font-size: 0.8rem;
    color: #666;
  }
}

.occupied-seat {
  position: relative;
  padding: 0.75rem;
  border-radius: 12px;
  background: rgba(255, 255, 255, 0.9);
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
  min-width: 140px;
}

.player-info {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 0.5rem;
  
  .player-avatar {
    position: relative;
    
    .turn-indicator {
      top: -5px;
      right: -5px;
    }
    
    .dealer-badge {
      top: -8px;
      right: -8px;
      font-size: 0.7rem;
      font-weight: bold;
    }
  }
  
  .player-details {
    text-align: center;
    
    .player-name {
      font-size: 0.9rem;
      font-weight: 600;
      color: #333;
      margin-bottom: 0.25rem;
      display: flex;
      align-items: center;
      gap: 0.25rem;
      justify-content: center;
    }
    
    .player-chips {
      font-size: 0.8rem;
      color: #666;
      display: flex;
      align-items: center;
      gap: 0.25rem;
      justify-content: center;
    }
  }
}

.player-cards {
  display: flex;
  gap: 0.25rem;
  margin-top: 0.5rem;
  justify-content: center;
  
  .hand-card {
    transform: rotate(-2deg);
    
    &:nth-child(2) {
      transform: rotate(2deg);
      margin-left: -8px;
    }
  }
}

.player-bet {
  position: absolute;
  top: -10px;
  right: -10px;
}

.player-status {
  margin-top: 0.5rem;
}

.last-action {
  position: absolute;
  bottom: -8px;
  left: 50%;
  transform: translateX(-50%);
  
  .action-chip {
    font-size: 0.7rem;
    opacity: 0.8;
  }
}

.connection-status {
  position: absolute;
  top: 5px;
  left: 5px;
}

// 位置特定样式
.seat-top {
  .player-info {
    flex-direction: column-reverse;
  }
  
  .player-cards {
    margin-top: 0;
    margin-bottom: 0.5rem;
  }
}

.seat-left,
.seat-right {
  .player-info {
    flex-direction: row;
    gap: 0.75rem;
  }
  
  .player-details {
    text-align: left;
  }
  
  .player-cards {
    position: absolute;
    top: 50%;
    transform: translateY(-50%);
  }
}

.seat-left {
  .player-cards {
    right: -30px;
  }
}

.seat-right {
  .player-cards {
    left: -30px;
  }
}

// 动画
@keyframes pulse {
  0%, 100% {
    box-shadow: 0 0 0 0 rgba(25, 118, 210, 0.7);
  }
  50% {
    box-shadow: 0 0 0 8px rgba(25, 118, 210, 0);
  }
}

// 响应式设计
@media (max-width: 768px) {
  .player-seat {
    min-width: 100px;
  }
  
  .occupied-seat {
    min-width: 120px;
    padding: 0.5rem;
  }
  
  .player-info {
    gap: 0.25rem;
  }
  
  .player-details {
    .player-name {
      font-size: 0.8rem;
    }
    
    .player-chips {
      font-size: 0.7rem;
    }
  }
}
</style> 