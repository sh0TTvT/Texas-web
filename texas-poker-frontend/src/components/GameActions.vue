<!--
游戏操作按钮组件
作用：提供玩家在游戏中的各种操作选项，如下注、跟注、加注、弃牌、全押等
-->

<template>
  <div class="game-actions" v-if="isMyTurn && availableActions.length > 0">
    <!-- 快速操作按钮 -->
    <div class="quick-actions">
      <q-btn
        v-for="action in quickActions"
        :key="action.type"
        :label="action.label"
        :color="action.color"
        :disable="!action.enabled"
        @click="handleQuickAction(action)"
        class="action-btn"
        unelevated
      >
        <q-icon :name="action.icon" class="q-mr-xs" />
      </q-btn>
    </div>

    <!-- 下注控制区域 -->
    <div v-if="canRaise" class="bet-control">
      <!-- 下注金额滑块 -->
      <div class="bet-slider">
        <q-slider
          v-model="betAmount"
          :min="minRaise"
          :max="maxRaise"
          :step="blindSize"
          :label-value="`${betAmount} 筹码`"
          label-always
          color="primary"
          class="bet-range"
        />
      </div>

      <!-- 预设金额按钮 -->
      <div class="preset-amounts">
        <q-btn
          v-for="preset in presetAmounts"
          :key="preset.value"
          :label="preset.label"
          :disable="preset.value > playerChips"
          @click="setBetAmount(preset.value)"
          size="sm"
          outline
          color="primary"
          class="preset-btn"
        />
      </div>

      <!-- 下注/加注按钮 -->
      <div class="bet-actions">
        <q-btn
          :label="raiseButtonText"
          :disable="betAmount < minRaise || betAmount > maxRaise"
          @click="handleRaise"
          color="primary"
          size="lg"
          class="raise-btn"
          unelevated
        >
          <q-icon name="trending_up" class="q-mr-xs" />
        </q-btn>

        <!-- 全押按钮 -->
        <q-btn
          label="全押"
          @click="handleAllIn"
          color="red"
          size="lg"
          class="allin-btn"
          unelevated
        >
          <q-icon name="local_fire_department" class="q-mr-xs" />
        </q-btn>
      </div>
    </div>

    <!-- 确认对话框 -->
    <q-dialog v-model="showConfirm" persistent>
      <q-card style="min-width: 300px">
        <q-card-section>
          <div class="text-h6">确认操作</div>
        </q-card-section>

        <q-card-section class="q-pt-none">
          <p>{{ confirmMessage }}</p>
          <div v-if="pendingAction.amount" class="confirm-amount">
            <q-chip color="primary" text-color="white" size="lg">
              {{ pendingAction.amount }} 筹码
            </q-chip>
          </div>
        </q-card-section>

        <q-card-actions align="right">
          <q-btn flat label="取消" color="grey" @click="cancelAction" />
          <q-btn 
            :label="confirmButtonText" 
            :color="pendingAction.color || 'primary'" 
            @click="confirmAction"
            unelevated
          />
        </q-card-actions>
      </q-card>
    </q-dialog>

    <!-- 操作提示 -->
    <div v-if="showHints" class="action-hints">
      <q-card flat class="hint-card">
        <q-card-section class="hint-content">
          <q-icon name="lightbulb" color="orange" />
          <span class="hint-text">{{ currentHint }}</span>
        </q-card-section>
      </q-card>
    </div>
  </div>

  <!-- 等待提示 -->
  <div v-else-if="gameInProgress" class="waiting-indicator">
    <q-chip color="orange" text-color="white" icon="hourglass_empty">
      {{ waitingMessage }}
    </q-chip>
  </div>
</template>

<script>
import { defineComponent, ref, computed, watch } from 'vue'
import { useQuasar } from 'quasar'

export default defineComponent({
  name: 'GameActions',
  
  props: {
    // 是否轮到当前玩家
    isMyTurn: {
      type: Boolean,
      default: false
    },
    // 可用操作列表
    availableActions: {
      type: Array,
      default: () => []
    },
    // 当前下注额
    currentBet: {
      type: Number,
      default: 0
    },
    // 玩家筹码
    playerChips: {
      type: Number,
      default: 0
    },
    // 底池大小
    potSize: {
      type: Number,
      default: 0
    },
    // 大盲注大小
    bigBlind: {
      type: Number,
      default: 10
    },
    // 小盲注大小
    smallBlind: {
      type: Number,
      default: 5
    },
    // 游戏是否进行中
    gameInProgress: {
      type: Boolean,
      default: false
    },
    // 是否显示提示
    showHints: {
      type: Boolean,
      default: true
    }
  },
  
  emits: ['action'],
  
  setup(props, { emit }) {
    const $q = useQuasar()
    
    // 响应式数据
    const betAmount = ref(0)
    const showConfirm = ref(false)
    const pendingAction = ref({})
    
    // 计算属性
    const blindSize = computed(() => props.smallBlind || 5)
    
    const minRaise = computed(() => {
      const minBet = Math.max(props.currentBet * 2, props.bigBlind)
      return Math.min(minBet, props.playerChips)
    })
    
    const maxRaise = computed(() => props.playerChips)
    
    const canCall = computed(() => 
      props.availableActions.includes('call') && props.currentBet <= props.playerChips
    )
    
    const canCheck = computed(() => 
      props.availableActions.includes('check')
    )
    
    const canRaise = computed(() => 
      props.availableActions.includes('raise') && props.playerChips > props.currentBet
    )
    
    const canFold = computed(() => 
      props.availableActions.includes('fold')
    )
    
    const canAllIn = computed(() => 
      props.availableActions.includes('all_in')
    )
    
    const callAmount = computed(() => 
      Math.min(props.currentBet, props.playerChips)
    )
    
    const quickActions = computed(() => {
      const actions = []
      
      // 弃牌
      if (canFold.value) {
        actions.push({
          type: 'fold',
          label: '弃牌',
          icon: 'close',
          color: 'grey',
          enabled: true
        })
      }
      
      // 看牌/跟注
      if (canCheck.value) {
        actions.push({
          type: 'check',
          label: '看牌',
          icon: 'visibility',
          color: 'green',
          enabled: true
        })
      } else if (canCall.value) {
        actions.push({
          type: 'call',
          label: `跟注 ${callAmount.value}`,
          icon: 'check',
          color: 'green',
          enabled: true,
          amount: callAmount.value
        })
      }
      
      return actions
    })
    
    const presetAmounts = computed(() => {
      const pot = props.potSize || 0
      const presets = [
        { label: '1/2底池', value: Math.round(pot / 2) },
        { label: '底池', value: pot },
        { label: '2倍底池', value: pot * 2 },
        { label: '3倍底池', value: pot * 3 }
      ]
      
      return presets.filter(preset => preset.value >= minRaise.value && preset.value <= maxRaise.value)
    })
    
    const raiseButtonText = computed(() => {
      if (props.currentBet === 0) {
        return `下注 ${betAmount.value}`
      } else {
        return `加注到 ${betAmount.value}`
      }
    })
    
    const confirmMessage = computed(() => {
      const action = pendingAction.value
      switch (action.type) {
        case 'fold':
          return '确定要弃牌吗？您将失去本局游戏的机会。'
        case 'call':
          return `确定要跟注 ${action.amount} 筹码吗？`
        case 'raise':
          return `确定要加注到 ${action.amount} 筹码吗？`
        case 'all_in':
          return `确定要全押所有筹码（${props.playerChips}）吗？`
        case 'check':
          return '确定要看牌吗？'
        default:
          return '确定要执行此操作吗？'
      }
    })
    
    const confirmButtonText = computed(() => {
      const action = pendingAction.value
      switch (action.type) {
        case 'fold': return '弃牌'
        case 'call': return '跟注'
        case 'raise': return '加注'
        case 'all_in': return '全押'
        case 'check': return '看牌'
        default: return '确认'
      }
    })
    
    const currentHint = computed(() => {
      if (!props.isMyTurn) return ''
      
      const hints = [
        '根据手牌强度和位置决定操作',
        '观察其他玩家的下注模式',
        '考虑底池赔率和隐含赔率',
        '注意控制筹码管理'
      ]
      
      return hints[Math.floor(Math.random() * hints.length)]
    })
    
    const waitingMessage = computed(() => {
      if (!props.gameInProgress) return '等待游戏开始'
      return '等待其他玩家操作'
    })
    
    // 监听器
    watch(() => props.isMyTurn, (newVal) => {
      if (newVal) {
        // 重置下注金额为最小加注
        betAmount.value = minRaise.value
      }
    })
    
    // 方法
    const handleQuickAction = (action) => {
      pendingAction.value = {
        type: action.type,
        amount: action.amount,
        color: action.color
      }
      showConfirm.value = true
    }
    
    const handleRaise = () => {
      pendingAction.value = {
        type: 'raise',
        amount: betAmount.value,
        color: 'primary'
      }
      showConfirm.value = true
    }
    
    const handleAllIn = () => {
      pendingAction.value = {
        type: 'all_in',
        amount: props.playerChips,
        color: 'red'
      }
      showConfirm.value = true
    }
    
    const setBetAmount = (amount) => {
      betAmount.value = Math.min(amount, maxRaise.value)
    }
    
    const confirmAction = () => {
      const action = pendingAction.value
      
      emit('action', {
        type: action.type,
        amount: action.amount
      })
      
      showConfirm.value = false
      pendingAction.value = {}
      
      // 显示操作反馈
      $q.notify({
        type: 'positive',
        message: `${confirmButtonText.value}成功`,
        position: 'top'
      })
    }
    
    const cancelAction = () => {
      showConfirm.value = false
      pendingAction.value = {}
    }
    
    return {
      betAmount,
      showConfirm,
      pendingAction,
      minRaise,
      maxRaise,
      canRaise,
      quickActions,
      presetAmounts,
      raiseButtonText,
      confirmMessage,
      confirmButtonText,
      currentHint,
      waitingMessage,
      blindSize,
      handleQuickAction,
      handleRaise,
      handleAllIn,
      setBetAmount,
      confirmAction,
      cancelAction
    }
  }
})
</script>

<style scoped lang="scss">
.game-actions {
  background: rgba(255, 255, 255, 0.95);
  border-radius: 12px;
  padding: 1rem;
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.15);
  backdrop-filter: blur(10px);
  border: 1px solid rgba(255, 255, 255, 0.2);
}

.quick-actions {
  display: flex;
  gap: 0.75rem;
  margin-bottom: 1rem;
  
  .action-btn {
    flex: 1;
    height: 48px;
    font-weight: 600;
  }
}

.bet-control {
  .bet-slider {
    margin-bottom: 1rem;
    
    .bet-range {
      padding: 0 0.5rem;
    }
  }
  
  .preset-amounts {
    display: flex;
    gap: 0.5rem;
    margin-bottom: 1rem;
    flex-wrap: wrap;
    
    .preset-btn {
      font-size: 0.8rem;
      min-width: 60px;
    }
  }
  
  .bet-actions {
    display: flex;
    gap: 0.75rem;
    
    .raise-btn,
    .allin-btn {
      flex: 1;
      height: 48px;
      font-weight: 600;
    }
  }
}

.confirm-amount {
  text-align: center;
  margin-top: 1rem;
}

.action-hints {
  margin-top: 1rem;
  
  .hint-card {
    background: rgba(255, 193, 7, 0.1);
    border: 1px solid rgba(255, 193, 7, 0.3);
    
    .hint-content {
      display: flex;
      align-items: center;
      gap: 0.5rem;
      padding: 0.75rem;
      
      .hint-text {
        font-size: 0.875rem;
        color: #666;
      }
    }
  }
}

.waiting-indicator {
  text-align: center;
  padding: 1rem;
}

// 响应式设计
@media (max-width: 768px) {
  .game-actions {
    padding: 0.75rem;
  }
  
  .quick-actions {
    flex-direction: column;
    gap: 0.5rem;
    
    .action-btn {
      height: 44px;
    }
  }
  
  .preset-amounts {
    .preset-btn {
      flex: 1;
      min-width: 50px;
    }
  }
  
  .bet-actions {
    flex-direction: column;
    gap: 0.5rem;
    
    .raise-btn,
    .allin-btn {
      height: 44px;
    }
  }
}

@media (max-width: 480px) {
  .preset-amounts {
    display: grid;
    grid-template-columns: 1fr 1fr;
    gap: 0.5rem;
  }
}

// 动画效果
.game-actions {
  animation: slideUp 0.3s ease-out;
}

@keyframes slideUp {
  from {
    transform: translateY(20px);
    opacity: 0;
  }
  to {
    transform: translateY(0);
    opacity: 1;
  }
}

// 按钮状态样式
.action-btn {
  transition: all 0.2s ease;
  
  &:hover:not(.disabled) {
    transform: translateY(-1px);
    box-shadow: 0 4px 12px rgba(0, 0, 0, 0.2);
  }
}

.raise-btn,
.allin-btn {
  transition: all 0.2s ease;
  
  &:hover:not(.disabled) {
    transform: translateY(-1px);
    box-shadow: 0 4px 12px rgba(0, 0, 0, 0.2);
  }
}
</style> 