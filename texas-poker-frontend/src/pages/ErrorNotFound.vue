<!--
404错误页面
作用：当用户访问不存在的路由时显示的错误页面
-->

<template>
  <div class="error-page">
    <div class="error-container">
      <!-- 错误图标和代码 -->
      <div class="error-header">
        <q-icon name="casino" size="120px" color="primary" class="error-icon" />
        <h1 class="error-code">404</h1>
        <h2 class="error-title">页面未找到</h2>
        <p class="error-message">
          抱歉，您访问的页面不存在。<br>
          可能是链接错误或页面已被移动。
        </p>
      </div>

      <!-- 操作按钮 -->
      <div class="error-actions">
        <q-btn
          label="返回首页"
          color="primary"
          size="lg"
          @click="goHome"
          unelevated
          class="action-btn"
        >
          <q-icon name="home" class="q-mr-sm" />
        </q-btn>
        
        <q-btn
          label="返回上页"
          color="secondary"
          size="lg"
          @click="goBack"
          outline
          class="action-btn"
        >
          <q-icon name="arrow_back" class="q-mr-sm" />
        </q-btn>
      </div>

      <!-- 装饰元素 -->
      <div class="decorations">
        <div class="floating-card card-1">
          <q-icon name="style" size="24px" color="red" />
        </div>
        <div class="floating-card card-2">
          <q-icon name="favorite" size="24px" color="red" />
        </div>
        <div class="floating-card card-3">
          <q-icon name="dashboard" size="24px" color="black" />
        </div>
        <div class="floating-card card-4">
          <q-icon name="star" size="24px" color="black" />
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import { defineComponent } from 'vue'
import { useRouter } from 'vue-router'

export default defineComponent({
  name: 'ErrorNotFound',
  
  setup() {
    const router = useRouter()
    
    const goHome = () => {
      router.push('/lobby')
    }
    
    const goBack = () => {
      if (window.history.length > 1) {
        router.go(-1)
      } else {
        router.push('/lobby')
      }
    }
    
    return {
      goHome,
      goBack
    }
  }
})
</script>

<style scoped lang="scss">
.error-page {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  position: relative;
  overflow: hidden;
}

.error-container {
  text-align: center;
  max-width: 600px;
  padding: 2rem;
  position: relative;
  z-index: 1;
}

.error-header {
  margin-bottom: 3rem;
  
  .error-icon {
    margin-bottom: 1rem;
    animation: bounce 2s infinite;
  }
  
  .error-code {
    font-size: 6rem;
    font-weight: 900;
    color: white;
    margin: 0;
    text-shadow: 4px 4px 8px rgba(0, 0, 0, 0.3);
    line-height: 1;
  }
  
  .error-title {
    font-size: 2rem;
    font-weight: 600;
    color: white;
    margin: 1rem 0;
    text-shadow: 2px 2px 4px rgba(0, 0, 0, 0.3);
  }
  
  .error-message {
    font-size: 1.1rem;
    color: rgba(255, 255, 255, 0.9);
    line-height: 1.6;
    margin: 0;
  }
}

.error-actions {
  display: flex;
  gap: 1rem;
  justify-content: center;
  flex-wrap: wrap;
  
  .action-btn {
    min-width: 150px;
    height: 48px;
    font-weight: 600;
  }
}

.decorations {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  pointer-events: none;
  
  .floating-card {
    position: absolute;
    width: 40px;
    height: 40px;
    background: white;
    border-radius: 8px;
    display: flex;
    align-items: center;
    justify-content: center;
    box-shadow: 0 4px 12px rgba(0, 0, 0, 0.2);
    
    &.card-1 {
      top: 10%;
      left: 10%;
      animation: float 3s ease-in-out infinite;
    }
    
    &.card-2 {
      top: 20%;
      right: 15%;
      animation: float 3s ease-in-out infinite 0.5s;
    }
    
    &.card-3 {
      bottom: 30%;
      left: 8%;
      animation: float 3s ease-in-out infinite 1s;
    }
    
    &.card-4 {
      bottom: 15%;
      right: 10%;
      animation: float 3s ease-in-out infinite 1.5s;
    }
  }
}

// 动画
@keyframes bounce {
  0%, 20%, 50%, 80%, 100% {
    transform: translateY(0);
  }
  40% {
    transform: translateY(-20px);
  }
  60% {
    transform: translateY(-10px);
  }
}

@keyframes float {
  0%, 100% {
    transform: translateY(0px) rotate(0deg);
  }
  50% {
    transform: translateY(-20px) rotate(5deg);
  }
}

// 响应式设计
@media (max-width: 768px) {
  .error-header {
    .error-code {
      font-size: 4rem;
    }
    
    .error-title {
      font-size: 1.5rem;
    }
    
    .error-message {
      font-size: 1rem;
    }
  }
  
  .error-actions {
    flex-direction: column;
    align-items: center;
    
    .action-btn {
      width: 100%;
      max-width: 300px;
    }
  }
  
  .decorations .floating-card {
    width: 30px;
    height: 30px;
  }
}

@media (max-width: 480px) {
  .error-container {
    padding: 1rem;
  }
  
  .error-header {
    margin-bottom: 2rem;
    
    .error-icon {
      font-size: 80px !important;
    }
    
    .error-code {
      font-size: 3rem;
    }
    
    .error-title {
      font-size: 1.25rem;
    }
  }
}
</style> 