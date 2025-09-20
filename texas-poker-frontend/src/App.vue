<!--
德州扑克应用根组件
作用：应用的根级容器，管理全局布局和路由视图
-->

<template>
  <div id="q-app">
    <router-view />
  </div>
</template>

<script>
import { defineComponent } from 'vue'

export default defineComponent({
  name: 'App'
})
</script>

<style>
#q-app {
  height: 100vh;
  width: 100vw;
}

body {
  margin: 0;
  padding: 0;
  font-family: 'Roboto', sans-serif;
}
</style> 