<!--
用户个人资料对话框组件
作用：显示和编辑用户个人信息，查看游戏统计
-->

<template>
  <q-card style="min-width: 400px; max-width: 600px">
    <q-card-section>
      <div class="text-h6">个人资料</div>
      <div class="text-subtitle2 text-grey-7">查看和编辑您的信息</div>
    </q-card-section>

    <q-card-section class="q-pt-none">
      <q-tabs
        v-model="activeTab"
        dense
        class="profile-tabs"
        active-color="primary"
        indicator-color="primary"
        align="left"
      >
        <q-tab name="info" label="基本信息" />
        <q-tab name="stats" label="游戏统计" />
        <q-tab name="settings" label="设置" />
      </q-tabs>

      <q-tab-panels v-model="activeTab" animated>
        <!-- 基本信息 -->
        <q-tab-panel name="info" class="q-pa-none">
          <div class="info-panel">
            <!-- 用户头像和基本信息 -->
            <div class="user-header">
              <q-avatar size="80px" color="primary" text-color="white">
                <q-icon name="person" size="40px" />
              </q-avatar>
              
              <div class="user-basic">
                <h6 class="username">{{ authStore.username }}</h6>
                <p class="email">{{ authStore.user?.email }}</p>
                <q-chip color="primary" text-color="white" icon="account_balance">
                  {{ authStore.userChips }} 筹码
                </q-chip>
              </div>
            </div>

            <!-- 编辑表单 -->
            <q-form @submit="handleUpdateProfile" class="profile-form">
              <q-input
                v-model="profileForm.username"
                label="用户名"
                outlined
                :rules="usernameRules"
                class="form-field"
              >
                <template v-slot:prepend>
                  <q-icon name="person" />
                </template>
              </q-input>

              <q-input
                v-model="profileForm.email"
                label="邮箱地址"
                type="email"
                outlined
                :rules="emailRules"
                class="form-field"
              >
                <template v-slot:prepend>
                  <q-icon name="email" />
                </template>
              </q-input>

              <q-input
                v-model="profileForm.currentPassword"
                label="当前密码"
                type="password"
                outlined
                class="form-field"
                hint="修改信息时需要验证当前密码"
              >
                <template v-slot:prepend>
                  <q-icon name="lock" />
                </template>
              </q-input>

              <q-input
                v-model="profileForm.newPassword"
                label="新密码 (可选)"
                type="password"
                outlined
                :rules="newPasswordRules"
                class="form-field"
                hint="留空表示不修改密码"
              >
                <template v-slot:prepend>
                  <q-icon name="lock_open" />
                </template>
              </q-input>

              <q-input
                v-if="profileForm.newPassword"
                v-model="profileForm.confirmPassword"
                label="确认新密码"
                type="password"
                outlined
                :rules="confirmPasswordRules"
                class="form-field"
              >
                <template v-slot:prepend>
                  <q-icon name="lock" />
                </template>
              </q-input>

              <div class="form-actions">
                <q-btn
                  label="保存修改"
                  color="primary"
                  type="submit"
                  :loading="isUpdating"
                  unelevated
                >
                  <template v-slot:loading>
                    <q-spinner-hourglass class="on-left" />
                    保存中...
                  </template>
                </q-btn>
              </div>
            </q-form>
          </div>
        </q-tab-panel>

        <!-- 游戏统计 -->
        <q-tab-panel name="stats" class="q-pa-none">
          <div class="stats-panel">
            <div class="stats-grid">
              <q-card flat bordered class="stat-card">
                <q-card-section class="stat-content">
                  <q-icon name="casino" size="32px" color="primary" />
                  <div class="stat-info">
                    <div class="stat-value">{{ authStore.user?.total_games || 0 }}</div>
                    <div class="stat-label">总游戏局数</div>
                  </div>
                </q-card-section>
              </q-card>

              <q-card flat bordered class="stat-card">
                <q-card-section class="stat-content">
                  <q-icon name="emoji_events" size="32px" color="orange" />
                  <div class="stat-info">
                    <div class="stat-value">{{ authStore.user?.total_wins || 0 }}</div>
                    <div class="stat-label">获胜局数</div>
                  </div>
                </q-card-section>
              </q-card>

              <q-card flat bordered class="stat-card">
                <q-card-section class="stat-content">
                  <q-icon name="trending_up" size="32px" color="green" />
                  <div class="stat-info">
                    <div class="stat-value">{{ winRate }}%</div>
                    <div class="stat-label">胜率</div>
                  </div>
                </q-card-section>
              </q-card>

              <q-card flat bordered class="stat-card">
                <q-card-section class="stat-content">
                  <q-icon name="account_balance" size="32px" color="blue" />
                  <div class="stat-info">
                    <div class="stat-value">{{ authStore.userChips }}</div>
                    <div class="stat-label">当前筹码</div>
                  </div>
                </q-card-section>
              </q-card>
            </div>

            <!-- 成就徽章 -->
            <div class="achievements-section">
              <h6 class="section-title">成就徽章</h6>
              <div class="achievements-grid">
                <q-card flat bordered class="achievement-card">
                  <q-card-section class="achievement-content">
                    <q-icon name="star" size="24px" color="gold" />
                    <div class="achievement-name">新手上路</div>
                  </q-card-section>
                </q-card>

                <q-card flat bordered class="achievement-card">
                  <q-card-section class="achievement-content">
                    <q-icon name="timeline" size="24px" color="silver" />
                    <div class="achievement-name">坚持不懈</div>
                  </q-card-section>
                </q-card>

                <q-card flat bordered class="achievement-card disabled">
                  <q-card-section class="achievement-content">
                    <q-icon name="emoji_events" size="24px" color="grey-5" />
                    <div class="achievement-name">常胜将军</div>
                  </q-card-section>
                </q-card>
              </div>
            </div>
          </div>
        </q-tab-panel>

        <!-- 设置 -->
        <q-tab-panel name="settings" class="q-pa-none">
          <div class="settings-panel">
            <q-list>
              <q-item>
                <q-item-section avatar>
                  <q-icon name="volume_up" />
                </q-item-section>
                <q-item-section>
                  <q-item-label>音效设置</q-item-label>
                  <q-item-label caption>游戏音效和背景音乐</q-item-label>
                </q-item-section>
                <q-item-section side>
                  <q-toggle v-model="settings.sound" color="primary" />
                </q-item-section>
              </q-item>

              <q-item>
                <q-item-section avatar>
                  <q-icon name="notifications" />
                </q-item-section>
                <q-item-section>
                  <q-item-label>通知设置</q-item-label>
                  <q-item-label caption>游戏通知和消息提醒</q-item-label>
                </q-item-section>
                <q-item-section side>
                  <q-toggle v-model="settings.notifications" color="primary" />
                </q-item-section>
              </q-item>

              <q-item>
                <q-item-section avatar>
                  <q-icon name="help" />
                </q-item-section>
                <q-item-section>
                  <q-item-label>新手提示</q-item-label>
                  <q-item-label caption>显示操作提示和手牌强度</q-item-label>
                </q-item-section>
                <q-item-section side>
                  <q-toggle v-model="settings.beginnerTips" color="primary" />
                </q-item-section>
              </q-item>

              <q-separator />

              <q-item>
                <q-item-section avatar>
                  <q-icon name="palette" />
                </q-item-section>
                <q-item-section>
                  <q-item-label>主题设置</q-item-label>
                  <q-item-label caption>选择您喜欢的主题颜色</q-item-label>
                </q-item-section>
                <q-item-section side>
                  <q-select
                    v-model="settings.theme"
                    :options="themeOptions"
                    dense
                    outlined
                    emit-value
                    map-options
                    style="min-width: 120px"
                  />
                </q-item-section>
              </q-item>
            </q-list>
          </div>
        </q-tab-panel>
      </q-tab-panels>
    </q-card-section>

    <q-card-actions align="right" class="q-pt-none">
      <q-btn
        flat
        label="关闭"
        color="grey"
        @click="$emit('close')"
      />
    </q-card-actions>
  </q-card>
</template>

<script>
import { defineComponent, ref, computed, reactive } from 'vue'
import { useQuasar } from 'quasar'
import { useAuthStore } from '../stores/auth'

export default defineComponent({
  name: 'ProfileDialog',
  
  emits: ['close'],
  
  setup() {
    const $q = useQuasar()
    const authStore = useAuthStore()
    
    const activeTab = ref('info')
    const isUpdating = ref(false)
    
    // 个人资料表单
    const profileForm = reactive({
      username: authStore.username,
      email: authStore.user?.email || '',
      currentPassword: '',
      newPassword: '',
      confirmPassword: ''
    })
    
    // 设置选项
    const settings = reactive({
      sound: true,
      notifications: true,
      beginnerTips: true,
      theme: 'blue'
    })
    
    const themeOptions = [
      { label: '经典蓝', value: 'blue' },
      { label: '深邃绿', value: 'green' },
      { label: '典雅紫', value: 'purple' },
      { label: '暖橙色', value: 'orange' }
    ]
    
    // 计算属性
    const winRate = computed(() => {
      const totalGames = authStore.user?.total_games || 0
      const totalWins = authStore.user?.total_wins || 0
      
      if (totalGames === 0) return 0
      return Math.round((totalWins / totalGames) * 100)
    })
    
    // 验证规则
    const usernameRules = [
      val => !!val || '请输入用户名',
      val => val.length >= 3 || '用户名至少3个字符',
      val => val.length <= 20 || '用户名不能超过20个字符'
    ]
    
    const emailRules = [
      val => !!val || '请输入邮箱地址',
      val => /.+@.+\..+/.test(val) || '请输入有效的邮箱地址'
    ]
    
    const newPasswordRules = [
      val => !val || val.length >= 8 || '新密码至少8个字符',
      val => !val || /(?=.*[a-zA-Z])(?=.*\d)/.test(val) || '新密码必须包含字母和数字'
    ]
    
    const confirmPasswordRules = [
      val => !profileForm.newPassword || val === profileForm.newPassword || '两次输入的密码不一致'
    ]
    
    // 方法
    const handleUpdateProfile = async () => {
      if (!profileForm.currentPassword) {
        $q.notify({
          type: 'negative',
          message: '请输入当前密码进行验证',
          position: 'top'
        })
        return
      }
      
      isUpdating.value = true
      
      try {
        const updateData = {
          username: profileForm.username,
          email: profileForm.email,
          current_password: profileForm.currentPassword
        }
        
        // 如果有新密码，添加到更新数据中
        if (profileForm.newPassword) {
          updateData.new_password = profileForm.newPassword
        }
        
        const result = await authStore.updateProfile(updateData)
        
        if (result.success) {
          $q.notify({
            type: 'positive',
            message: '个人资料更新成功',
            position: 'top'
          })
          
          // 清空密码字段
          profileForm.currentPassword = ''
          profileForm.newPassword = ''
          profileForm.confirmPassword = ''
        } else {
          $q.notify({
            type: 'negative',
            message: result.message,
            position: 'top'
          })
        }
      } catch (error) {
        $q.notify({
          type: 'negative',
          message: '更新失败，请稍后重试',
          position: 'top'
        })
      } finally {
        isUpdating.value = false
      }
    }
    
    return {
      authStore,
      activeTab,
      isUpdating,
      profileForm,
      settings,
      themeOptions,
      winRate,
      usernameRules,
      emailRules,
      newPasswordRules,
      confirmPasswordRules,
      handleUpdateProfile
    }
  }
})
</script>

<style scoped lang="scss">
.profile-tabs {
  margin-bottom: 1.5rem;
}

.info-panel {
  padding: 1rem 0;
  
  .user-header {
    display: flex;
    align-items: center;
    gap: 1rem;
    margin-bottom: 2rem;
    padding: 1rem;
    background: #f8f9fa;
    border-radius: 8px;
    
    .user-basic {
      flex: 1;
      
      .username {
        margin: 0 0 0.25rem 0;
        font-size: 1.25rem;
        font-weight: 600;
        color: #1976d2;
      }
      
      .email {
        margin: 0 0 0.5rem 0;
        color: #666;
        font-size: 0.9rem;
      }
    }
  }
  
  .profile-form {
    .form-field {
      margin-bottom: 1rem;
    }
    
    .form-actions {
      margin-top: 1.5rem;
      text-align: right;
    }
  }
}

.stats-panel {
  padding: 1rem 0;
  
  .stats-grid {
    display: grid;
    grid-template-columns: repeat(auto-fit, minmax(150px, 1fr));
    gap: 1rem;
    margin-bottom: 2rem;
  }
  
  .stat-card {
    .stat-content {
      display: flex;
      align-items: center;
      gap: 1rem;
      
      .stat-info {
        .stat-value {
          font-size: 1.5rem;
          font-weight: 600;
          color: #333;
        }
        
        .stat-label {
          font-size: 0.875rem;
          color: #666;
        }
      }
    }
  }
  
  .achievements-section {
    .section-title {
      margin: 0 0 1rem 0;
      color: #1976d2;
      font-weight: 600;
    }
    
    .achievements-grid {
      display: grid;
      grid-template-columns: repeat(auto-fill, minmax(120px, 1fr));
      gap: 1rem;
    }
    
    .achievement-card {
      text-align: center;
      transition: transform 0.2s ease;
      
      &:not(.disabled):hover {
        transform: translateY(-2px);
      }
      
      &.disabled {
        opacity: 0.5;
      }
      
      .achievement-content {
        padding: 1rem;
        
        .achievement-name {
          margin-top: 0.5rem;
          font-size: 0.8rem;
          font-weight: 500;
        }
      }
    }
  }
}

.settings-panel {
  padding: 1rem 0;
}

// 响应式设计
@media (max-width: 480px) {
  .user-header {
    flex-direction: column;
    text-align: center;
  }
  
  .stats-grid {
    grid-template-columns: repeat(2, 1fr) !important;
  }
  
  .achievements-grid {
    grid-template-columns: repeat(2, 1fr) !important;
  }
}
</style> 