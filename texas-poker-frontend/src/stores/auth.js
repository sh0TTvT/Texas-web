// 认证状态管理
// 作用：管理用户登录状态、用户信息、token等认证相关数据

import { defineStore } from 'pinia'
import { api } from '../api'

export const useAuthStore = defineStore('auth', {
  state: () => ({
    user: null,
    token: localStorage.getItem('token'),
    isLoading: false
  }),

  getters: {
    isAuthenticated: (state) => !!state.token && !!state.user,
    isAdmin: (state) => state.user?.role === 'admin',
    userChips: (state) => state.user?.chips || 0,
    username: (state) => state.user?.username || ''
  },

  actions: {
    // 初始化认证状态
    async initAuth() {
      if (this.token) {
        try {
          await this.fetchUserProfile()
        } catch (error) {
          console.error('初始化认证失败:', error)
          this.clearAuth()
        }
      }
    },

    // 用户登录
    async login(credentials) {
      this.isLoading = true
      try {
        const response = await api.post('/auth/login', credentials)
        const { user, token } = response.data
        
        this.user = user
        this.token = token
        localStorage.setItem('token', token)
        
        return { success: true, user }
      } catch (error) {
        const message = error.response?.data?.error || '登录失败'
        return { success: false, message }
      } finally {
        this.isLoading = false
      }
    },

    // 用户注册
    async register(userData) {
      this.isLoading = true
      try {
        const response = await api.post('/auth/register', userData)
        const { user, token } = response.data
        
        this.user = user
        this.token = token
        localStorage.setItem('token', token)
        
        return { success: true, user }
      } catch (error) {
        const message = error.response?.data?.error || '注册失败'
        return { success: false, message }
      } finally {
        this.isLoading = false
      }
    },

    // 用户登出
    async logout() {
      try {
        await api.post('/auth/logout')
      } catch (error) {
        console.error('登出请求失败:', error)
      } finally {
        this.clearAuth()
      }
    },

    // 获取用户信息
    async fetchUserProfile() {
      try {
        const response = await api.get('/auth/profile')
        this.user = response.data.user
        return this.user
      } catch (error) {
        console.error('获取用户信息失败:', error)
        throw error
      }
    },

    // 更新用户信息
    async updateProfile(profileData) {
      this.isLoading = true
      try {
        const response = await api.put('/auth/profile', profileData)
        this.user = response.data.user
        return { success: true, user: this.user }
      } catch (error) {
        const message = error.response?.data?.error || '更新失败'
        return { success: false, message }
      } finally {
        this.isLoading = false
      }
    },

    // 更新用户筹码（游戏结束后）
    updateUserChips(newChips) {
      if (this.user) {
        this.user.chips = newChips
      }
    },

    // 清除认证状态
    clearAuth() {
      this.user = null
      this.token = null
      localStorage.removeItem('token')
    }
  }
}) 