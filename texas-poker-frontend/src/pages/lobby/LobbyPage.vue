<!--
游戏大厅页面
作用：显示房间列表，用户信息，提供创建房间和加入房间功能
-->

<template>
  <q-layout view="hHh lpR fFf">
    <!-- 顶部导航栏 -->
    <q-header elevated class="bg-primary text-white">
      <q-toolbar>
        <!-- 品牌Logo -->
        <q-toolbar-title class="brand-title">
          <q-icon name="casino" size="sm" class="q-mr-sm" />
          德州扑克
        </q-toolbar-title>

        <!-- 用户信息 -->
        <div class="user-info">
          <q-chip color="white" text-color="primary" icon="account_balance">
            {{ authStore.userChips }} 筹码
          </q-chip>
          
          <q-btn flat round icon="person" class="q-ml-sm">
            <q-menu>
              <q-list style="min-width: 200px">
                <q-item>
                  <q-item-section>
                    <q-item-label class="text-weight-bold">
                      {{ authStore.username }}
                    </q-item-label>
                    <q-item-label caption>
                      {{ authStore.user?.email }}
                    </q-item-label>
                  </q-item-section>
                </q-item>
                
                <q-separator />
                
                <q-item clickable @click="showProfile = true">
                  <q-item-section avatar>
                    <q-icon name="settings" />
                  </q-item-section>
                  <q-item-section>个人设置</q-item-section>
                </q-item>
                
                <q-item clickable @click="handleLogout">
                  <q-item-section avatar>
                    <q-icon name="logout" />
                  </q-item-section>
                  <q-item-section>退出登录</q-item-section>
                </q-item>
              </q-list>
            </q-menu>
          </q-btn>
        </div>
      </q-toolbar>
    </q-header>

    <!-- 主内容区 -->
    <q-page-container>
      <q-page class="lobby-page">
        <div class="page-container">
          <!-- 页面标题和操作按钮 -->
          <div class="page-header">
            <div class="header-content">
              <h4 class="page-title">游戏大厅</h4>
              <p class="page-subtitle">选择房间开始游戏</p>
            </div>
            
            <div class="header-actions">
              <q-btn 
                color="primary" 
                icon="add" 
                label="创建房间"
                @click="showCreateRoom = true"
                unelevated
              />
              
              <q-btn 
                color="secondary" 
                icon="refresh" 
                label="刷新"
                @click="refreshRooms"
                :loading="gameStore.isLoadingRooms"
                unelevated
                class="q-ml-sm"
              />
            </div>
          </div>

          <!-- 房间列表 -->
          <div class="rooms-container">
            <!-- 房间筛选标签 -->
            <q-tabs
              v-model="selectedTab"
              dense
              class="room-tabs"
              active-color="primary"
              indicator-color="primary"
              align="left"
            >
              <q-tab name="all" label="全部房间" />
              <q-tab name="low" label="低级场" />
              <q-tab name="medium" label="中级场" />
              <q-tab name="high" label="高级场" />
            </q-tabs>

            <!-- 房间卡片列表 -->
            <div class="rooms-grid">
              <div v-if="gameStore.isLoadingRooms" class="loading-container">
                <q-spinner size="48px" color="primary" />
                <p class="loading-text">加载房间列表...</p>
              </div>
              
              <div v-else-if="filteredRooms.length === 0" class="empty-container">
                <q-icon name="casino" size="64px" color="grey-5" />
                <p class="empty-text">暂无可用房间</p>
                <q-btn 
                  color="primary" 
                  label="创建房间" 
                  @click="showCreateRoom = true"
                  unelevated
                />
              </div>
              
              <room-card
                v-else
                v-for="room in filteredRooms"
                :key="room.id"
                :room="room"
                @join="handleJoinRoom"
                class="room-card"
              />
            </div>
          </div>
        </div>
      </q-page>
    </q-page-container>

    <!-- 创建房间对话框 -->
    <q-dialog v-model="showCreateRoom" persistent>
      <create-room-dialog 
        @close="showCreateRoom = false"
        @created="handleRoomCreated"
      />
    </q-dialog>

    <!-- 个人设置对话框 -->
    <q-dialog v-model="showProfile">
      <profile-dialog @close="showProfile = false" />
    </q-dialog>
  </q-layout>
</template>

<script>
import { defineComponent, ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useQuasar } from 'quasar'
import { useAuthStore } from '../../stores/auth'
import { useGameStore } from '../../stores/game'
import RoomCard from '../../components/RoomCard.vue'
import CreateRoomDialog from '../../components/CreateRoomDialog.vue'
import ProfileDialog from '../../components/ProfileDialog.vue'

export default defineComponent({
  name: 'LobbyPage',
  
  components: {
    RoomCard,
    CreateRoomDialog,
    ProfileDialog
  },
  
  setup() {
    const router = useRouter()
    const $q = useQuasar()
    const authStore = useAuthStore()
    const gameStore = useGameStore()
    
    // 响应式数据
    const selectedTab = ref('all')
    const showCreateRoom = ref(false)
    const showProfile = ref(false)
    
    // 计算属性
    const filteredRooms = computed(() => {
      if (selectedTab.value === 'all') {
        return gameStore.rooms
      }
      
      return gameStore.rooms.filter(room => 
        room.chip_level === selectedTab.value
      )
    })
    
    // 方法
    const refreshRooms = async () => {
      const result = await gameStore.fetchRooms()
      if (result.success) {
        $q.notify({
          type: 'positive',
          message: '房间列表已刷新',
          position: 'top'
        })
      } else {
        $q.notify({
          type: 'negative',
          message: result.message,
          position: 'top'
        })
      }
    }
    
    const handleJoinRoom = async (roomId) => {
      // 检查筹码是否足够
      const room = gameStore.rooms.find(r => r.id === roomId)
      if (room && authStore.userChips < room.min_chips) {
        $q.notify({
          type: 'negative',
          message: `筹码不足，需要至少 ${room.min_chips} 筹码`,
          position: 'top'
        })
        return
      }
      
      const result = await gameStore.joinRoom(roomId)
      if (result.success) {
        $q.notify({
          type: 'positive',
          message: '成功加入房间',
          position: 'top'
        })
        
        // 跳转到游戏房间
        router.push(`/game/${roomId}`)
      } else {
        $q.notify({
          type: 'negative',
          message: result.message,
          position: 'top'
        })
      }
    }
    
    const handleRoomCreated = (room) => {
      showCreateRoom.value = false
      $q.notify({
        type: 'positive',
        message: '房间创建成功',
        position: 'top'
      })
      
      // 跳转到新创建的房间
      router.push(`/game/${room.id}`)
    }
    
    const handleLogout = async () => {
      $q.dialog({
        title: '确认',
        message: '确定要退出登录吗？',
        cancel: true,
        persistent: true
      }).onOk(async () => {
        await authStore.logout()
        $q.notify({
          type: 'positive',
          message: '已退出登录',
          position: 'top'
        })
        router.push('/auth/login')
      })
    }
    
    // 生命周期
    onMounted(() => {
      // 初始化用户认证状态
      authStore.initAuth()
      
      // 获取房间列表
      gameStore.fetchRooms()
    })
    
    return {
      authStore,
      gameStore,
      selectedTab,
      showCreateRoom,
      showProfile,
      filteredRooms,
      refreshRooms,
      handleJoinRoom,
      handleRoomCreated,
      handleLogout
    }
  }
})
</script>

<style scoped lang="scss">
.lobby-page {
  background: linear-gradient(135deg, #f5f7fa 0%, #c3cfe2 100%);
  min-height: calc(100vh - 64px);
}

.brand-title {
  font-weight: 600;
  font-size: 1.2rem;
}

.user-info {
  display: flex;
  align-items: center;
}

.page-container {
  max-width: 1200px;
  margin: 0 auto;
  padding: 2rem;
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 2rem;
  
  .header-content {
    .page-title {
      margin: 0 0 0.5rem 0;
      color: #1976d2;
      font-weight: 600;
    }
    
    .page-subtitle {
      margin: 0;
      color: #666;
      font-size: 1rem;
    }
  }
  
  .header-actions {
    display: flex;
    gap: 0.5rem;
  }
}

.rooms-container {
  background: white;
  border-radius: 12px;
  padding: 1.5rem;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.1);
  
  .room-tabs {
    margin-bottom: 1.5rem;
  }
}

.rooms-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
  gap: 1.5rem;
  min-height: 200px;
}

.loading-container,
.empty-container {
  grid-column: 1 / -1;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 3rem;
  
  .loading-text,
  .empty-text {
    margin: 1rem 0;
    color: #666;
    font-size: 1rem;
  }
}

// 响应式设计
@media (max-width: 768px) {
  .page-container {
    padding: 1rem;
  }
  
  .page-header {
    flex-direction: column;
    align-items: flex-start;
    gap: 1rem;
    
    .header-actions {
      width: 100%;
      justify-content: center;
    }
  }
  
  .rooms-grid {
    grid-template-columns: 1fr;
  }
  
  .user-info {
    .q-chip {
      font-size: 0.8rem;
    }
  }
}

@media (max-width: 480px) {
  .brand-title {
    font-size: 1rem;
  }
  
  .header-actions {
    flex-direction: column;
    width: 100%;
    
    .q-btn {
      width: 100%;
    }
  }
}
</style> 