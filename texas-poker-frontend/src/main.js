// Vue应用主入口文件
// 作用：初始化Vue应用，配置Quasar、Pinia、路由等

import { createApp } from 'vue'
import { Quasar } from 'quasar'
import { createPinia } from 'pinia'
import router from './router'

// Import icon libraries
import '@quasar/extras/material-icons/material-icons.css'

// Import Quasar css
import 'quasar/src/css/index.sass'

// Assumes your root component is App.vue
// and placed in same folder as main.js
import App from './App.vue'

const myApp = createApp(App)

// 安装Pinia状态管理
const pinia = createPinia()
myApp.use(pinia)

// 安装Quasar
myApp.use(Quasar, {
  plugins: {
    // import Quasar plugins and add here
  }
})

// 安装Vue Router
myApp.use(router)

// Assumes you have a <div id="app"></div> in your index.html
myApp.mount('#app') 