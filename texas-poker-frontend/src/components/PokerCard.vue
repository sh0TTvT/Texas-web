<!--
扑克牌组件
作用：显示单张扑克牌，支持花色、点数、隐藏状态
-->

<template>
  <div 
    class="poker-card"
    :class="cardClasses"
    @click="handleClick"
  >
    <div v-if="!isHidden" class="card-content">
      <!-- 左上角 -->
      <div class="card-corner top-left">
        <div class="rank">{{ rankDisplay }}</div>
        <div class="suit" :style="{ color: suitColor }">{{ suitSymbol }}</div>
      </div>
      
      <!-- 中央花色 -->
      <div class="card-center">
        <div class="center-suit" :style="{ color: suitColor }">
          {{ suitSymbol }}
        </div>
      </div>
      
      <!-- 右下角（旋转180度） -->
      <div class="card-corner bottom-right">
        <div class="rank">{{ rankDisplay }}</div>
        <div class="suit" :style="{ color: suitColor }">{{ suitSymbol }}</div>
      </div>
    </div>
    
    <!-- 背面 -->
    <div v-else class="card-back">
      <div class="card-pattern"></div>
    </div>
  </div>
</template>

<script>
import { defineComponent, computed } from 'vue'

export default defineComponent({
  name: 'PokerCard',
  
  props: {
    // 花色：hearts, diamonds, clubs, spades
    suit: {
      type: String,
      default: null
    },
    // 点数：A, 2-10, J, Q, K
    rank: {
      type: String, 
      default: null
    },
    // 是否隐藏（显示背面）
    isHidden: {
      type: Boolean,
      default: false
    },
    // 卡片大小
    size: {
      type: String,
      default: 'normal', // small, normal, large
      validator: val => ['small', 'normal', 'large'].includes(val)
    },
    // 是否可点击
    clickable: {
      type: Boolean,
      default: false
    },
    // 是否选中状态
    selected: {
      type: Boolean,
      default: false
    },
    // 是否禁用状态
    disabled: {
      type: Boolean,
      default: false
    }
  },
  
  emits: ['click'],
  
  setup(props, { emit }) {
    // 计算属性
    const suitSymbol = computed(() => {
      const symbols = {
        hearts: '♥',
        diamonds: '♦',
        clubs: '♣',
        spades: '♠'
      }
      return symbols[props.suit] || '?'
    })
    
    const suitColor = computed(() => {
      const redSuits = ['hearts', 'diamonds']
      return redSuits.includes(props.suit) ? '#e53e3e' : '#2d3748'
    })
    
    const rankDisplay = computed(() => {
      if (!props.rank) return '?'
      return props.rank
    })
    
    const cardClasses = computed(() => {
      return {
        [`card-${props.size}`]: true,
        'card-clickable': props.clickable,
        'card-selected': props.selected,
        'card-disabled': props.disabled,
        'card-hidden': props.isHidden,
        'card-red': ['hearts', 'diamonds'].includes(props.suit),
        'card-black': ['clubs', 'spades'].includes(props.suit)
      }
    })
    
    // 方法
    const handleClick = () => {
      if (props.clickable && !props.disabled) {
        emit('click', { suit: props.suit, rank: props.rank })
      }
    }
    
    return {
      suitSymbol,
      suitColor,
      rankDisplay,
      cardClasses,
      handleClick
    }
  }
})
</script>

<style scoped lang="scss">
.poker-card {
  position: relative;
  border-radius: 8px;
  background: white;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.15);
  transition: all 0.3s ease;
  user-select: none;
  overflow: hidden;
  
  // 默认尺寸（normal）
  width: 60px;
  height: 84px;
  
  &.card-small {
    width: 45px;
    height: 63px;
    
    .card-corner {
      font-size: 0.7rem;
      
      .rank { font-size: 0.65rem; }
      .suit { font-size: 0.6rem; }
    }
    
    .center-suit {
      font-size: 1.2rem;
    }
  }
  
  &.card-large {
    width: 80px;
    height: 112px;
    
    .card-corner {
      font-size: 1rem;
      
      .rank { font-size: 0.9rem; }
      .suit { font-size: 0.8rem; }
    }
    
    .center-suit {
      font-size: 2rem;
    }
  }
  
  &.card-clickable {
    cursor: pointer;
    
    &:hover:not(.card-disabled) {
      transform: translateY(-2px);
      box-shadow: 0 4px 16px rgba(0, 0, 0, 0.25);
    }
  }
  
  &.card-selected {
    transform: translateY(-4px);
    box-shadow: 0 6px 20px rgba(25, 118, 210, 0.4);
    border: 2px solid #1976d2;
  }
  
  &.card-disabled {
    opacity: 0.5;
    cursor: not-allowed;
  }
  
  &.card-hidden {
    background: linear-gradient(135deg, #1976d2, #1565c0);
  }
}

.card-content {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  padding: 4px;
}

.card-corner {
  position: absolute;
  display: flex;
  flex-direction: column;
  align-items: center;
  line-height: 1;
  font-weight: 600;
  font-size: 0.8rem;
  
  &.top-left {
    top: 4px;
    left: 4px;
  }
  
  &.bottom-right {
    bottom: 4px;
    right: 4px;
    transform: rotate(180deg);
  }
  
  .rank {
    font-size: 0.75rem;
    font-weight: 700;
  }
  
  .suit {
    font-size: 0.7rem;
    margin-top: 1px;
  }
}

.card-center {
  position: absolute;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
  
  .center-suit {
    font-size: 1.5rem;
    font-weight: bold;
    opacity: 0.8;
  }
}

.card-back {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  display: flex;
  align-items: center;
  justify-content: center;
  
  .card-pattern {
    width: 80%;
    height: 80%;
    background-image: 
      repeating-linear-gradient(
        45deg,
        rgba(255, 255, 255, 0.1),
        rgba(255, 255, 255, 0.1) 2px,
        transparent 2px,
        transparent 8px
      ),
      repeating-linear-gradient(
        -45deg,
        rgba(255, 255, 255, 0.1),
        rgba(255, 255, 255, 0.1) 2px,
        transparent 2px,
        transparent 8px
      );
    border-radius: 4px;
  }
}

// 动画效果
.poker-card {
  &.card-flip-enter-active,
  &.card-flip-leave-active {
    transition: transform 0.6s ease-in-out;
  }
  
  &.card-flip-enter-from {
    transform: rotateY(-90deg);
  }
  
  &.card-flip-leave-to {
    transform: rotateY(90deg);
  }
}

// 特殊状态样式
.card-red {
  border-left: 3px solid #e53e3e;
}

.card-black {
  border-left: 3px solid #2d3748;
}
</style> 