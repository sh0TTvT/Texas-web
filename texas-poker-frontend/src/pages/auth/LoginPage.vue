<!--
用户登录页面
作用：提供用户登录表单，处理登录逻辑，支持用户名或邮箱登录
-->

<template>
  <div class="login-page">
    <!-- 页面标题 -->
    <div class="page-header">
      <h2 class="page-title">欢迎回来</h2>
      <p class="page-subtitle">登录您的账户开始游戏</p>
    </div>

    <!-- 登录表单 -->
    <q-form @submit="handleLogin" class="login-form">
      <!-- 用户名/邮箱输入框 -->
      <q-input
        v-model="form.username"
        label="用户名或邮箱"
        type="text"
        outlined
        :rules="usernameRules"
        :loading="isLoading"
        class="form-input"
      >
        <template v-slot:prepend>
          <q-icon name="person" />
        </template>
      </q-input>

      <!-- 密码输入框 -->
      <q-input
        v-model="form.password"
        :label="'密码'"
        :type="showPassword ? 'text' : 'password'"
        outlined
        :rules="passwordRules"
        :loading="isLoading"
        class="form-input"
      >
        <template v-slot:prepend>
          <q-icon name="lock" />
        </template>
        <template v-slot:append>
          <q-icon
            :name="showPassword ? 'visibility_off' : 'visibility'"
            class="cursor-pointer"
            @click="showPassword = !showPassword"
          />
        </template>
      </q-input>

      <!-- 记住我选项 -->
      <q-checkbox
        v-model="form.rememberMe"
        label="记住我"
        color="primary"
        class="remember-checkbox"
      />

      <!-- 登录按钮 -->
      <q-btn
        type="submit"
        label="登录"
        color="primary"
        size="lg"
        :loading="isLoading"
        :disable="!canSubmit"
        class="login-btn"
        unelevated
      >
        <template v-slot:loading>
          <q-spinner-hourglass class="on-left" />
          登录中...
        </template>
      </q-btn>

      <!-- 分隔线 -->
      <q-separator class="form-separator" />

      <!-- 注册链接 -->
      <div class="auth-links">
        <p class="link-text">
          还没有账户？
          <router-link to="/auth/register" class="auth-link">
            立即注册
          </router-link>
        </p>
      </div>
    </q-form>

    <!-- 快速登录提示（开发阶段） -->
    <div class="dev-hints" v-if="isDev">
      <q-expansion-item
        icon="info"
        label="开发测试账号"
        dense
      >
        <q-card>
          <q-card-section class="text-caption">
            <p><strong>测试用户：</strong> testuser / test@example.com</p>
            <p><strong>密码：</strong> test1234</p>
            <p><strong>管理员：</strong> admin / admin@texaspoker.com</p>
            <p><strong>密码：</strong> admin1234</p>
          </q-card-section>
        </q-card>
      </q-expansion-item>
    </div>
  </div>
</template>

<script>
import { defineComponent, ref, computed } from 'vue'
import { useRouter } from 'vue-router'
import { useQuasar } from 'quasar'
import { useAuthStore } from '../../stores/auth'

export default defineComponent({
  name: 'LoginPage',
  
  setup() {
    const router = useRouter()
    const $q = useQuasar()
    const authStore = useAuthStore()
    
    // 响应式数据
    const form = ref({
      username: '',
      password: '',
      rememberMe: false
    })
    
    const showPassword = ref(false)
    const isLoading = ref(false)
    
    // 计算属性
    const isDev = computed(() => process.env.NODE_ENV === 'development')
    
    const canSubmit = computed(() => 
      form.value.username.length > 0 && 
      form.value.password.length > 0 && 
      !isLoading.value
    )
    
    // 验证规则
    const usernameRules = [
      val => !!val || '请输入用户名或邮箱',
      val => val.length >= 3 || '用户名至少3个字符'
    ]
    
    const passwordRules = [
      val => !!val || '请输入密码',
      val => val.length >= 6 || '密码至少6个字符'
    ]
    
    // 方法
    const handleLogin = async () => {
      isLoading.value = true
      
      try {
        const result = await authStore.login({
          username: form.value.username,
          password: form.value.password
        })
        
        if (result.success) {
          $q.notify({
            type: 'positive',
            message: `欢迎回来，${result.user.username}！`,
            position: 'top'
          })
          
          // 根据用户角色跳转
          if (result.user.role === 'admin') {
            router.push('/admin')
          } else {
            router.push('/lobby')
          }
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
          message: '登录失败，请稍后重试',
          position: 'top'
        })
      } finally {
        isLoading.value = false
      }
    }
    
    return {
      form,
      showPassword,
      isLoading,
      isDev,
      canSubmit,
      usernameRules,
      passwordRules,
      handleLogin
    }
  }
})
</script>

<style scoped lang="scss">
.login-page {
  width: 100%;
}

.page-header {
  text-align: center;
  margin-bottom: 2rem;
  
  .page-title {
    font-size: 1.75rem;
    font-weight: 600;
    color: #1976d2;
    margin: 0 0 0.5rem 0;
  }
  
  .page-subtitle {
    color: #666;
    margin: 0;
    font-size: 0.95rem;
  }
}

.login-form {
  .form-input {
    margin-bottom: 1.5rem;
  }
  
  .remember-checkbox {
    margin-bottom: 2rem;
  }
  
  .login-btn {
    width: 100%;
    height: 48px;
    font-weight: 600;
    border-radius: 8px;
  }
  
  .form-separator {
    margin: 2rem 0;
  }
}

.auth-links {
  text-align: center;
  
  .link-text {
    margin: 0;
    color: #666;
    font-size: 0.95rem;
  }
  
  .auth-link {
    color: #1976d2;
    text-decoration: none;
    font-weight: 500;
    
    &:hover {
      text-decoration: underline;
    }
  }
}

.dev-hints {
  margin-top: 2rem;
  padding: 1rem;
  background: rgba(25, 118, 210, 0.1);
  border-radius: 8px;
  border: 1px solid rgba(25, 118, 210, 0.2);
}

// 响应式设计
@media (max-width: 480px) {
  .page-header .page-title {
    font-size: 1.5rem;
  }
  
  .login-form .form-input {
    margin-bottom: 1.25rem;
  }
}
</style> 