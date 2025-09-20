<!--
用户注册页面
作用：提供用户注册表单，处理注册逻辑，表单验证
-->

<template>
  <div class="register-page">
    <!-- 页面标题 -->
    <div class="page-header">
      <h2 class="page-title">创建账户</h2>
      <p class="page-subtitle">加入德州扑克，开始您的游戏之旅</p>
    </div>

    <!-- 注册表单 -->
    <q-form @submit="handleRegister" class="register-form">
      <!-- 用户名输入框 -->
      <q-input
        v-model="form.username"
        label="用户名"
        type="text"
        outlined
        :rules="usernameRules"
        :loading="isLoading"
        class="form-input"
        hint="3-20个字符，只能包含字母、数字和下划线"
      >
        <template v-slot:prepend>
          <q-icon name="person" />
        </template>
      </q-input>

      <!-- 邮箱输入框 -->
      <q-input
        v-model="form.email"
        label="邮箱地址"
        type="email"
        outlined
        :rules="emailRules"
        :loading="isLoading"
        class="form-input"
        hint="用于账户验证和密码找回"
      >
        <template v-slot:prepend>
          <q-icon name="email" />
        </template>
      </q-input>

      <!-- 密码输入框 -->
      <q-input
        v-model="form.password"
        label="密码"
        :type="showPassword ? 'text' : 'password'"
        outlined
        :rules="passwordRules"
        :loading="isLoading"
        class="form-input"
        hint="至少8个字符，包含字母和数字"
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

      <!-- 确认密码输入框 -->
      <q-input
        v-model="form.confirmPassword"
        label="确认密码"
        :type="showConfirmPassword ? 'text' : 'password'"
        outlined
        :rules="confirmPasswordRules"
        :loading="isLoading"
        class="form-input"
      >
        <template v-slot:prepend>
          <q-icon name="lock" />
        </template>
        <template v-slot:append>
          <q-icon
            :name="showConfirmPassword ? 'visibility_off' : 'visibility'"
            class="cursor-pointer"
            @click="showConfirmPassword = !showConfirmPassword"
          />
        </template>
      </q-input>

      <!-- 服务条款同意 -->
      <q-checkbox
        v-model="form.agreeTerms"
        color="primary"
        class="terms-checkbox"
      >
        <span class="terms-text">
          我已阅读并同意
          <a href="#" class="terms-link" @click.prevent="showTerms = true">
            服务条款
          </a>
          和
          <a href="#" class="terms-link" @click.prevent="showPrivacy = true">
            隐私政策
          </a>
        </span>
      </q-checkbox>

      <!-- 注册按钮 -->
      <q-btn
        type="submit"
        label="创建账户"
        color="primary"
        size="lg"
        :loading="isLoading"
        :disable="!canSubmit"
        class="register-btn"
        unelevated
      >
        <template v-slot:loading>
          <q-spinner-hourglass class="on-left" />
          注册中...
        </template>
      </q-btn>

      <!-- 分隔线 -->
      <q-separator class="form-separator" />

      <!-- 登录链接 -->
      <div class="auth-links">
        <p class="link-text">
          已经有账户？
          <router-link to="/auth/login" class="auth-link">
            立即登录
          </router-link>
        </p>
      </div>
    </q-form>

    <!-- 服务条款对话框 -->
    <q-dialog v-model="showTerms">
      <q-card style="min-width: 400px">
        <q-card-section>
          <div class="text-h6">服务条款</div>
        </q-card-section>
        <q-card-section class="q-pt-none">
          <div class="text-body2">
            <p>欢迎使用德州扑克在线游戏平台。在使用我们的服务之前，请仔细阅读以下条款：</p>
            <ul>
              <li>本平台仅供娱乐，不涉及真实货币交易</li>
              <li>禁止使用外挂、作弊软件等违法行为</li>
              <li>用户应维护良好的游戏环境，禁止恶意言论</li>
              <li>平台有权对违规账户进行处罚</li>
              <li>用户数据将得到妥善保护</li>
            </ul>
          </div>
        </q-card-section>
        <q-card-actions align="right">
          <q-btn flat label="关闭" color="primary" v-close-popup />
        </q-card-actions>
      </q-card>
    </q-dialog>

    <!-- 隐私政策对话框 -->
    <q-dialog v-model="showPrivacy">
      <q-card style="min-width: 400px">
        <q-card-section>
          <div class="text-h6">隐私政策</div>
        </q-card-section>
        <q-card-section class="q-pt-none">
          <div class="text-body2">
            <p>我们重视您的隐私，承诺保护您的个人信息：</p>
            <ul>
              <li>我们收集必要的账户信息用于游戏服务</li>
              <li>不会向第三方泄露您的个人信息</li>
              <li>使用加密技术保护数据传输安全</li>
              <li>您可以随时查看和修改个人信息</li>
              <li>如需删除账户，请联系客服</li>
            </ul>
          </div>
        </q-card-section>
        <q-card-actions align="right">
          <q-btn flat label="关闭" color="primary" v-close-popup />
        </q-card-actions>
      </q-card>
    </q-dialog>
  </div>
</template>

<script>
import { defineComponent, ref, computed } from 'vue'
import { useRouter } from 'vue-router'
import { useQuasar } from 'quasar'
import { useAuthStore } from '../../stores/auth'

export default defineComponent({
  name: 'RegisterPage',
  
  setup() {
    const router = useRouter()
    const $q = useQuasar()
    const authStore = useAuthStore()
    
    // 响应式数据
    const form = ref({
      username: '',
      email: '',
      password: '',
      confirmPassword: '',
      agreeTerms: false
    })
    
    const showPassword = ref(false)
    const showConfirmPassword = ref(false)
    const isLoading = ref(false)
    const showTerms = ref(false)
    const showPrivacy = ref(false)
    
    // 计算属性
    const canSubmit = computed(() => 
      form.value.username.length >= 3 && 
      form.value.email.length > 0 && 
      form.value.password.length >= 8 && 
      form.value.password === form.value.confirmPassword &&
      form.value.agreeTerms &&
      !isLoading.value
    )
    
    // 验证规则
    const usernameRules = [
      val => !!val || '请输入用户名',
      val => val.length >= 3 || '用户名至少3个字符',
      val => val.length <= 20 || '用户名不能超过20个字符',
      val => /^[a-zA-Z0-9_]+$/.test(val) || '用户名只能包含字母、数字和下划线'
    ]
    
    const emailRules = [
      val => !!val || '请输入邮箱地址',
      val => /.+@.+\..+/.test(val) || '请输入有效的邮箱地址'
    ]
    
    const passwordRules = [
      val => !!val || '请输入密码',
      val => val.length >= 8 || '密码至少8个字符',
      val => /(?=.*[a-zA-Z])(?=.*\d)/.test(val) || '密码必须包含字母和数字'
    ]
    
    const confirmPasswordRules = [
      val => !!val || '请确认密码',
      val => val === form.value.password || '两次输入的密码不一致'
    ]
    
    // 方法
    const handleRegister = async () => {
      isLoading.value = true
      
      try {
        const result = await authStore.register({
          username: form.value.username,
          email: form.value.email,
          password: form.value.password
        })
        
        if (result.success) {
          $q.notify({
            type: 'positive',
            message: `注册成功！欢迎加入，${result.user.username}！`,
            position: 'top'
          })
          
          // 注册成功后跳转到游戏大厅
          router.push('/lobby')
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
          message: '注册失败，请稍后重试',
          position: 'top'
        })
      } finally {
        isLoading.value = false
      }
    }
    
    return {
      form,
      showPassword,
      showConfirmPassword,
      isLoading,
      showTerms,
      showPrivacy,
      canSubmit,
      usernameRules,
      emailRules,
      passwordRules,
      confirmPasswordRules,
      handleRegister
    }
  }
})
</script>

<style scoped lang="scss">
.register-page {
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

.register-form {
  .form-input {
    margin-bottom: 1.5rem;
  }
  
  .terms-checkbox {
    margin-bottom: 2rem;
    align-items: flex-start;
    
    .terms-text {
      font-size: 0.9rem;
      line-height: 1.4;
      color: #666;
    }
    
    .terms-link {
      color: #1976d2;
      text-decoration: none;
      
      &:hover {
        text-decoration: underline;
      }
    }
  }
  
  .register-btn {
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

// 响应式设计
@media (max-width: 480px) {
  .page-header .page-title {
    font-size: 1.5rem;
  }
  
  .register-form .form-input {
    margin-bottom: 1.25rem;
  }
  
  .terms-checkbox {
    margin-bottom: 1.5rem;
  }
}
</style> 