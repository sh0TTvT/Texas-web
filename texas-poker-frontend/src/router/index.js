// Vue Router路由配置
// 作用：定义应用的路由规则，管理页面导航和权限控制

import { createRouter, createWebHashHistory } from 'vue-router'
import { useAuthStore } from '../stores/auth'

const routes = [
  {
    path: '/',
    redirect: '/auth/login'
  },
  {
    path: '/auth',
    component: () => import('../pages/auth/AuthLayout.vue'),
    children: [
      {
        path: 'login',
        name: 'Login',
        component: () => import('../pages/auth/LoginPage.vue')
      },
      {
        path: 'register', 
        name: 'Register',
        component: () => import('../pages/auth/RegisterPage.vue')
      }
    ]
  },
  {
    path: '/lobby',
    name: 'Lobby',
    component: () => import('../pages/lobby/LobbyPage.vue'),
    meta: { requiresAuth: true }
  },
  {
    path: '/game/:roomId',
    name: 'GameRoom',
    component: () => import('../pages/game/GameRoomPage.vue'),
    meta: { requiresAuth: true },
    props: true
  },
  {
    path: '/admin',
    component: () => import('../pages/admin/AdminLayout.vue'),
    meta: { requiresAuth: true, requiresAdmin: true },
    children: [
      {
        path: '',
        redirect: '/admin/dashboard'
      },
      {
        path: 'dashboard',
        name: 'AdminDashboard',
        component: () => import('../pages/admin/DashboardPage.vue')
      },
      {
        path: 'users',
        name: 'AdminUsers',
        component: () => import('../pages/admin/UsersPage.vue')
      },
      {
        path: 'rooms',
        name: 'AdminRooms', 
        component: () => import('../pages/admin/RoomsPage.vue')
      }
    ]
  },
  // 始终保留这个作为最后一个，否则会覆盖其他路由
  {
    path: '/:catchAll(.*)*',
    component: () => import('../pages/ErrorNotFound.vue')
  }
]

const router = createRouter({
  scrollBehavior: () => ({ left: 0, top: 0 }),
  routes,
  // Leave this as is and make changes in quasar.conf.js instead!
  // quasar.conf.js -> build -> vueRouterMode
  // quasar.conf.js -> build -> publicPath
  history: createWebHashHistory(process.env.VUE_ROUTER_BASE)
})

// 路由守卫
router.beforeEach((to, from, next) => {
  const authStore = useAuthStore()
  
  // 检查是否需要认证
  if (to.meta.requiresAuth && !authStore.isAuthenticated) {
    next('/auth/login')
    return
  }
  
  // 检查是否需要管理员权限
  if (to.meta.requiresAdmin && !authStore.isAdmin) {
    next('/lobby')
    return
  }
  
  // 如果已登录且访问认证页面，重定向到大厅
  if (authStore.isAuthenticated && to.path.startsWith('/auth')) {
    next('/lobby')
    return
  }
  
  next()
})

export default router 