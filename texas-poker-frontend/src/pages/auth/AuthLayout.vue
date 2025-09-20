<!--
认证页面布局组件
作用：为登录和注册页面提供统一的布局和样式
-->

<template>
  <q-layout view="hHh lpR fFf">
    <q-page-container>
      <q-page class="auth-page">
        <!-- 背景装饰 -->
        <div class="auth-background">
          <div class="bg-pattern"></div>
        </div>
        
        <!-- 主要内容区域 -->
        <div class="auth-container">
          <!-- 品牌Logo和标题 -->
          <div class="auth-header">
            <div class="brand-logo">
              <q-icon name="casino" size="64px" color="primary" />
            </div>
            <h1 class="brand-title">德州扑克</h1>
            <p class="brand-subtitle">Texas Hold'em Poker</p>
          </div>

          <!-- 认证表单容器 -->
          <div class="auth-form-container">
            <router-view />
          </div>

          <!-- 页面底部 -->
          <div class="auth-footer">
            <p class="footer-text">
              © 2024 德州扑克在线. 保留所有权利.
            </p>
          </div>
        </div>
      </q-page>
    </q-page-container>
  </q-layout>
</template>

<script>
import { defineComponent } from 'vue'

export default defineComponent({
  name: 'AuthLayout'
})
</script>

<style scoped lang="scss">
.auth-page {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  position: relative;
  background: linear-gradient(135deg, #1e3c72 0%, #2a5298 100%);
}

.auth-background {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  overflow: hidden;
  
  .bg-pattern {
    position: absolute;
    top: -50%;
    left: -50%;
    width: 200%;
    height: 200%;
    background-image: 
      radial-gradient(circle at 25% 25%, rgba(255, 255, 255, 0.1) 2px, transparent 2px),
      radial-gradient(circle at 75% 75%, rgba(255, 255, 255, 0.1) 2px, transparent 2px);
    background-size: 60px 60px;
    animation: float 20s ease-in-out infinite;
  }
}

@keyframes float {
  0%, 100% { transform: translateY(0px) rotate(0deg); }
  50% { transform: translateY(-20px) rotate(1deg); }
}

.auth-container {
  width: 100%;
  max-width: 400px;
  padding: 2rem;
  position: relative;
  z-index: 1;
}

.auth-header {
  text-align: center;
  margin-bottom: 2rem;
  
  .brand-logo {
    margin-bottom: 1rem;
  }
  
  .brand-title {
    font-size: 2.5rem;
    font-weight: 700;
    color: white;
    margin: 0 0 0.5rem 0;
    text-shadow: 2px 2px 4px rgba(0, 0, 0, 0.3);
  }
  
  .brand-subtitle {
    font-size: 1.1rem;
    color: rgba(255, 255, 255, 0.8);
    margin: 0;
    font-weight: 300;
  }
}

.auth-form-container {
  background: rgba(255, 255, 255, 0.95);
  backdrop-filter: blur(10px);
  border-radius: 16px;
  padding: 2rem;
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.3);
  border: 1px solid rgba(255, 255, 255, 0.2);
}

.auth-footer {
  text-align: center;
  margin-top: 2rem;
  
  .footer-text {
    color: rgba(255, 255, 255, 0.7);
    font-size: 0.875rem;
    margin: 0;
  }
}

// 响应式设计
@media (max-width: 480px) {
  .auth-container {
    padding: 1rem;
    max-width: 90%;
  }
  
  .auth-header .brand-title {
    font-size: 2rem;
  }
  
  .auth-form-container {
    padding: 1.5rem;
  }
}

@media (max-height: 700px) {
  .auth-header {
    margin-bottom: 1rem;
  }
  
  .auth-header .brand-title {
    font-size: 2rem;
  }
}
</style> 