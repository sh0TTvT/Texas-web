<!--
房间卡片组件
作用：显示单个房间的信息，提供加入房间功能
-->

<template>
  <q-card class="room-card" bordered>
    <!-- 房间头部 -->
    <q-card-section class="room-header">
      <div class="room-title">
        <h6 class="room-name">{{ room.name }}</h6>
        <q-badge 
          :color="chipLevelColor" 
          :label="chipLevelText"
          class="level-badge"
        />
      </div>
      
      <div class="room-status">
        <q-chip 
          :color="statusColor" 
          :label="statusText"
          size="sm"
          dense
        />
      </div>
    </q-card-section>

    <!-- 房间信息 -->
    <q-card-section class="room-info">
      <div class="info-grid">
        <!-- 玩家数量 -->
        <div class="info-item">
          <q-icon name="people" color="primary" />
          <span class="info-label">玩家</span>
          <span class="info-value">
            {{ room.current_players || 0 }}/{{ room.max_players }}
          </span>
        </div>

        <!-- 最低筹码 -->
        <div class="info-item">
          <q-icon name="account_balance" color="orange" />
          <span class="info-label">最低筹码</span>
          <span class="info-value">{{ room.min_chips }}</span>
        </div>

        <!-- 盲注信息 -->
        <div class="info-item">
          <q-icon name="visibility" color="green" />
          <span class="info-label">盲注</span>
          <span class="info-value">
            {{ room.small_blind }}/{{ room.big_blind }}
          </span>
        </div>

        <!-- 房间类型 -->
        <div class="info-item">
          <q-icon :name="room.is_private ? 'lock' : 'public'" color="grey" />
          <span class="info-label">类型</span>
          <span class="info-value">
            {{ room.is_private ? '私人' : '公开' }}
          </span>
        </div>
      </div>
    </q-card-section>

    <!-- 操作按钮 -->
    <q-card-actions class="room-actions">
      <q-btn
        :color="canJoin ? 'primary' : 'grey'"
        :label="joinButtonText"
        :disable="!canJoin"
        @click="handleJoin"
        unelevated
        class="join-btn"
      >
        <q-icon :name="joinButtonIcon" class="q-mr-xs" />
      </q-btn>
      
      <q-btn
        flat
        color="primary"
        icon="info"
        @click="showRoomInfo = true"
      >
        <q-tooltip>查看详情</q-tooltip>
      </q-btn>
    </q-card-actions>

    <!-- 房间详情对话框 -->
    <q-dialog v-model="showRoomInfo">
      <q-card style="min-width: 400px">
        <q-card-section>
          <div class="text-h6">房间详情</div>
        </q-card-section>

        <q-card-section class="q-pt-none">
          <q-list dense>
            <q-item>
              <q-item-section avatar>
                <q-icon name="title" />
              </q-item-section>
              <q-item-section>
                <q-item-label>房间名称</q-item-label>
                <q-item-label caption>{{ room.name }}</q-item-label>
              </q-item-section>
            </q-item>

            <q-item>
              <q-item-section avatar>
                <q-icon name="star" />
              </q-item-section>
              <q-item-section>
                <q-item-label>级别</q-item-label>
                <q-item-label caption>{{ chipLevelText }}</q-item-label>
              </q-item-section>
            </q-item>

            <q-item>
              <q-item-section avatar>
                <q-icon name="people" />
              </q-item-section>
              <q-item-section>
                <q-item-label>当前玩家</q-item-label>
                <q-item-label caption>
                  {{ room.current_players || 0 }} / {{ room.max_players }}
                </q-item-label>
              </q-item-section>
            </q-item>

            <q-item>
              <q-item-section avatar>
                <q-icon name="account_balance" />
              </q-item-section>
              <q-item-section>
                <q-item-label>最低筹码要求</q-item-label>
                <q-item-label caption>{{ room.min_chips }} 筹码</q-item-label>
              </q-item-section>
            </q-item>

            <q-item>
              <q-item-section avatar>
                <q-icon name="visibility" />
              </q-item-section>
              <q-item-section>
                <q-item-label>盲注设置</q-item-label>
                <q-item-label caption>
                  小盲注 {{ room.small_blind }} / 大盲注 {{ room.big_blind }}
                </q-item-label>
              </q-item-section>
            </q-item>

            <q-item v-if="room.created_at">
              <q-item-section avatar>
                <q-icon name="schedule" />
              </q-item-section>
              <q-item-section>
                <q-item-label>创建时间</q-item-label>
                <q-item-label caption>
                  {{ formatDate(room.created_at) }}
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
  </q-card>
</template>

<script>
import { defineComponent, ref, computed } from 'vue'
import { useAuthStore } from '../stores/auth'

export default defineComponent({
  name: 'RoomCard',
  
  props: {
    room: {
      type: Object,
      required: true
    }
  },
  
  emits: ['join'],
  
  setup(props, { emit }) {
    const authStore = useAuthStore()
    const showRoomInfo = ref(false)
    
    // 计算属性
    const chipLevelColor = computed(() => {
      const colors = {
        low: 'green',
        medium: 'orange', 
        high: 'red'
      }
      return colors[props.room.chip_level] || 'grey'
    })
    
    const chipLevelText = computed(() => {
      const texts = {
        low: '低级场',
        medium: '中级场',
        high: '高级场'
      }
      return texts[props.room.chip_level] || '未知'
    })
    
    const statusColor = computed(() => {
      const colors = {
        waiting: 'green',
        playing: 'orange',
        closed: 'red'
      }
      return colors[props.room.status] || 'grey'
    })
    
    const statusText = computed(() => {
      const texts = {
        waiting: '等待中',
        playing: '游戏中',
        closed: '已关闭'
      }
      return texts[props.room.status] || '未知'
    })
    
    const canJoin = computed(() => {
      const hasEnoughChips = authStore.userChips >= props.room.min_chips
      const roomNotFull = (props.room.current_players || 0) < props.room.max_players
      const roomWaiting = props.room.status === 'waiting'
      
      return hasEnoughChips && roomNotFull && roomWaiting
    })
    
    const joinButtonText = computed(() => {
      if (props.room.status === 'closed') return '房间已关闭'
      if ((props.room.current_players || 0) >= props.room.max_players) return '房间已满'
      if (authStore.userChips < props.room.min_chips) return '筹码不足'
      if (props.room.status === 'playing') return '游戏中'
      return '加入游戏'
    })
    
    const joinButtonIcon = computed(() => {
      if (!canJoin.value) return 'block'
      return 'play_arrow'
    })
    
    // 方法
    const handleJoin = () => {
      if (canJoin.value) {
        emit('join', props.room.id)
      }
    }
    
    const formatDate = (dateString) => {
      if (!dateString) return ''
      const date = new Date(dateString)
      return date.toLocaleString('zh-CN')
    }
    
    return {
      authStore,
      showRoomInfo,
      chipLevelColor,
      chipLevelText,
      statusColor,
      statusText,
      canJoin,
      joinButtonText,
      joinButtonIcon,
      handleJoin,
      formatDate
    }
  }
})
</script>

<style scoped lang="scss">
.room-card {
  transition: all 0.3s ease;
  
  &:hover {
    transform: translateY(-2px);
    box-shadow: 0 4px 20px rgba(0, 0, 0, 0.15);
  }
}

.room-header {
  padding: 1rem 1rem 0.5rem 1rem;
  
  .room-title {
    display: flex;
    justify-content: space-between;
    align-items: flex-start;
    margin-bottom: 0.5rem;
    
    .room-name {
      margin: 0;
      font-size: 1.1rem;
      font-weight: 600;
      color: #1976d2;
      flex: 1;
      margin-right: 0.5rem;
    }
    
    .level-badge {
      font-size: 0.7rem;
      font-weight: 500;
    }
  }
  
  .room-status {
    display: flex;
    justify-content: flex-end;
  }
}

.room-info {
  padding: 0.5rem 1rem;
  
  .info-grid {
    display: grid;
    grid-template-columns: 1fr 1fr;
    gap: 0.75rem;
  }
  
  .info-item {
    display: flex;
    align-items: center;
    gap: 0.5rem;
    
    .info-label {
      font-size: 0.85rem;
      color: #666;
      flex: 1;
    }
    
    .info-value {
      font-size: 0.9rem;
      font-weight: 500;
      color: #333;
    }
  }
}

.room-actions {
  padding: 0.5rem 1rem 1rem 1rem;
  display: flex;
  justify-content: space-between;
  
  .join-btn {
    flex: 1;
    margin-right: 0.5rem;
  }
}

// 响应式设计
@media (max-width: 480px) {
  .room-header {
    .room-title {
      flex-direction: column;
      align-items: flex-start;
      
      .level-badge {
        margin-top: 0.25rem;
      }
    }
  }
  
  .info-grid {
    grid-template-columns: 1fr;
    gap: 0.5rem;
  }
}
</style> 