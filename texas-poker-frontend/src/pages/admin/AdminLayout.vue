<!--
后台管理布局组件
作用：为后台管理页面提供统一的布局和导航
-->

<template>
  <q-layout view="hHh lpR fFf">
    <!-- 侧边栏 -->
    <q-drawer
      v-model="leftDrawerOpen"
      side="left"
      behavior="desktop"
      elevated
      :width="280"
      :breakpoint="1024"
    >
      <q-scroll-area class="fit">
        <!-- 品牌标识 -->
        <div class="admin-header">
          <q-icon name="casino" size="32px" color="primary" />
          <div class="brand-text">
            <div class="brand-name">德州扑克</div>
            <div class="brand-subtitle">管理后台</div>
          </div>
        </div>

        <!-- 导航菜单 -->
        <q-list padding>
          <q-item
            v-for="item in menuItems"
            :key="item.name"
            :to="item.to"
            clickable
            v-ripple
            exact-active-class="menu-item-active"
            class="menu-item"
          >
            <q-item-section avatar>
              <q-icon :name="item.icon" />
            </q-item-section>
            <q-item-section>
              <q-item-label>{{ item.label }}</q-item-label>
              <q-item-label caption v-if="item.description">
                {{ item.description }}
              </q-item-label>
            </q-item-section>
          </q-item>
        </q-list>
      </q-scroll-area>
    </q-drawer>

    <!-- 顶部导航栏 -->
    <q-header elevated class="bg-white text-dark">
      <q-toolbar>
        <!-- 菜单按钮 -->
        <q-btn
          flat
          dense
          round
          icon="menu"
          aria-label="Menu"
          @click="leftDrawerOpen = !leftDrawerOpen"
        />

        <!-- 页面标题 -->
        <q-toolbar-title class="page-title">
          <q-icon :name="currentPageIcon" class="q-mr-sm" />
          {{ currentPageTitle }}
        </q-toolbar-title>

        <q-space />

        <!-- 通知 -->
        <q-btn flat round icon="notifications" class="q-mr-sm">
          <q-badge color="red" floating>2</q-badge>
          <q-menu>
            <q-list style="min-width: 300px">
              <q-item>
                <q-item-section>
                  <q-item-label>系统通知</q-item-label>
                  <q-item-label caption>暂无新通知</q-item-label>
                </q-item-section>
              </q-item>
            </q-list>
          </q-menu>
        </q-btn>

        <!-- 用户菜单 -->
        <q-btn flat round icon="account_circle">
          <q-menu>
            <q-list style="min-width: 200px">
              <q-item>
                <q-item-section>
                  <q-item-label class="text-weight-bold">
                    {{ authStore.username }}
                  </q-item-label>
                  <q-item-label caption>
                    管理员
                  </q-item-label>
                </q-item-section>
              </q-item>
              
              <q-separator />
              
              <q-item clickable @click="showProfile = true">
                <q-item-section avatar>
                  <q-icon name="person" />
                </q-item-section>
                <q-item-section>个人资料</q-item-section>
              </q-item>
              
              <q-item clickable @click="showSettings = true">
                <q-item-section avatar>
                  <q-icon name="settings" />
                </q-item-section>
                <q-item-section>系统设置</q-item-section>
              </q-item>
              
              <q-separator />
              
              <q-item clickable @click="handleLogout">
                <q-item-section avatar>
                  <q-icon name="logout" />
                </q-item-section>
                <q-item-section>退出登录</q-item-section>
              </q-item>
            </q-list>
          </q-menu>
        </q-btn>
      </q-toolbar>
    </q-header>

    <!-- 主内容区 -->
    <q-page-container>
      <div class="admin-content">
        <router-view />
      </div>
    </q-page-container>

    <!-- 个人资料对话框 -->
    <q-dialog v-model="showProfile">
      <profile-dialog @close="showProfile = false" />
    </q-dialog>

    <!-- 系统设置对话框 -->
    <q-dialog v-model="showSettings">
      <q-card style="min-width: 400px">
        <q-card-section>
          <div class="text-h6">系统设置</div>
        </q-card-section>

        <q-card-section class="q-pt-none">
          <q-list>
            <q-item>
              <q-item-section>
                <q-item-label>自动刷新</q-item-label>
                <q-item-label caption>自动刷新数据</q-item-label>
              </q-item-section>
              <q-item-section side>
                <q-toggle v-model="settings.autoRefresh" color="primary" />
              </q-item-section>
            </q-item>

            <q-item>
              <q-item-section>
                <q-item-label>桌面通知</q-item-label>
                <q-item-label caption>接收桌面通知</q-item-label>
              </q-item-section>
              <q-item-section side>
                <q-toggle v-model="settings.notifications" color="primary" />
              </q-item-section>
            </q-item>

            <q-item>
              <q-item-section>
                <q-item-label>数据刷新间隔</q-item-label>
                <q-item-label caption>{{ settings.refreshInterval }}秒</q-item-label>
              </q-item-section>
              <q-item-section side>
                <q-slider
                  v-model="settings.refreshInterval"
                  :min="10"
                  :max="60"
                  :step="5"
                  style="width: 100px"
                />
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
import { defineComponent, ref, computed } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useQuasar } from 'quasar'
import { useAuthStore } from '../../stores/auth'
import ProfileDialog from '../../components/ProfileDialog.vue'

export default defineComponent({
  name: 'AdminLayout',
  
  components: {
    ProfileDialog
  },
  
  setup() {
    const router = useRouter()
    const route = useRoute()
    const $q = useQuasar()
    const authStore = useAuthStore()
    
    // 响应式数据
    const leftDrawerOpen = ref(true)
    const showProfile = ref(false)
    const showSettings = ref(false)
    
    const settings = ref({
      autoRefresh: true,
      notifications: true,
      refreshInterval: 30
    })
    
    // 菜单项配置
    const menuItems = [
      {
        name: 'dashboard',
        label: '仪表盘',
        description: '系统总览',
        icon: 'dashboard',
        to: '/admin/dashboard'
      },
      {
        name: 'users',
        label: '用户管理',
        description: '用户账号管理',
        icon: 'people',
        to: '/admin/users'
      },
      {
        name: 'rooms',
        label: '房间管理',
        description: '游戏房间管理',
        icon: 'meeting_room',
        to: '/admin/rooms'
      },
      {
        name: 'games',
        label: '游戏记录',
        description: '游戏数据分析',
        icon: 'casino',
        to: '/admin/games'
      },
      {
        name: 'finance',
        label: '财务统计',
        description: '筹码流水分析',
        icon: 'account_balance',
        to: '/admin/finance'
      },
      {
        name: 'logs',
        label: '系统日志',
        description: '操作和错误日志',
        icon: 'list_alt',
        to: '/admin/logs'
      }
    ]
    
    // 计算属性
    const currentPageTitle = computed(() => {
      const currentItem = menuItems.find(item => 
        route.path.startsWith(item.to)
      )
      return currentItem?.label || '管理后台'
    })
    
    const currentPageIcon = computed(() => {
      const currentItem = menuItems.find(item => 
        route.path.startsWith(item.to)
      )
      return currentItem?.icon || 'dashboard'
    })
    
    // 方法
    const handleLogout = () => {
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
    
    return {
      authStore,
      leftDrawerOpen,
      showProfile,
      showSettings,
      settings,
      menuItems,
      currentPageTitle,
      currentPageIcon,
      handleLogout
    }
  }
})
</script>

<style scoped lang="scss">
.admin-header {
  display: flex;
  align-items: center;
  padding: 1.5rem 1rem;
  border-bottom: 1px solid rgba(0, 0, 0, 0.12);
  background: #f8f9fa;
  
  .brand-text {
    margin-left: 0.75rem;
    
    .brand-name {
      font-size: 1.1rem;
      font-weight: 600;
      color: #1976d2;
      line-height: 1.2;
    }
    
    .brand-subtitle {
      font-size: 0.85rem;
      color: #666;
      line-height: 1.2;
    }
  }
}

.menu-item {
  margin: 0.25rem 0.5rem;
  border-radius: 8px;
  transition: all 0.2s ease;
  
  &:hover {
    background: rgba(25, 118, 210, 0.08);
  }
  
  &.menu-item-active {
    background: rgba(25, 118, 210, 0.12);
    color: #1976d2;
    
    .q-icon {
      color: #1976d2;
    }
  }
}

.page-title {
  font-size: 1.25rem;
  font-weight: 600;
  display: flex;
  align-items: center;
}

.admin-content {
  padding: 1.5rem;
  min-height: calc(100vh - 64px);
  background: #f5f5f5;
}

// 响应式设计
@media (max-width: 1024px) {
  .admin-content {
    padding: 1rem;
  }
}

@media (max-width: 768px) {
  .admin-content {
    padding: 0.75rem;
  }
  
  .page-title {
    font-size: 1.1rem;
  }
}

// 抽屉样式
:deep(.q-drawer) {
  .q-scrollarea__content {
    display: flex;
    flex-direction: column;
    min-height: 100%;
  }
}

// 工具栏样式
:deep(.q-toolbar) {
  min-height: 64px;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.12);
}

// 菜单项动画
.menu-item {
  &.q-router-link--exact-active {
    .q-item__section--avatar .q-icon {
      transform: scale(1.1);
    }
  }
}

// 品牌区域动画
.admin-header {
  .q-icon {
    transition: transform 0.3s ease;
    
    &:hover {
      transform: rotate(5deg) scale(1.05);
    }
  }
}
</style> 