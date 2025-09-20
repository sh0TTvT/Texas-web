// API接口配置
// 作用：配置axios实例，处理请求拦截器，提供统一的API调用接口

import axios from 'axios'
import { Notify } from 'quasar'

// 创建axios实例
export const api = axios.create({
  baseURL: process.env.API_BASE_URL || 'http://localhost:8080/api',
  timeout: 10000,
  headers: {
    'Content-Type': 'application/json'
  }
})

// 请求拦截器
api.interceptors.request.use(
  (config) => {
    // 从localStorage获取token
    const token = localStorage.getItem('token')
    if (token) {
      config.headers.Authorization = `Bearer ${token}`
    }
    return config
  },
  (error) => {
    return Promise.reject(error)
  }
)

// 响应拦截器
api.interceptors.response.use(
  (response) => {
    return response
  },
  (error) => {
    // 处理通用错误
    if (error.response) {
      const { status, data } = error.response
      
      switch (status) {
        case 401:
          // 未授权，清除token并跳转到登录页
          localStorage.removeItem('token')
          if (window.location.hash !== '#/auth/login') {
            window.location.hash = '#/auth/login'
          }
          Notify.create({
            type: 'negative',
            message: '请重新登录',
            position: 'top'
          })
          break
          
        case 403:
          Notify.create({
            type: 'negative',
            message: '权限不足',
            position: 'top'
          })
          break
          
        case 404:
          Notify.create({
            type: 'negative',
            message: '请求的资源不存在',
            position: 'top'
          })
          break
          
        case 500:
          Notify.create({
            type: 'negative',
            message: '服务器内部错误',
            position: 'top'
          })
          break
          
        default:
          Notify.create({
            type: 'negative',
            message: data?.error || '请求失败',
            position: 'top'
          })
      }
    } else if (error.request) {
      // 网络错误
      Notify.create({
        type: 'negative',
        message: '网络连接失败，请检查网络设置',
        position: 'top'
      })
    } else {
      // 其他错误
      Notify.create({
        type: 'negative',
        message: '请求失败，请稍后重试',
        position: 'top'
      })
    }
    
    return Promise.reject(error)
  }
)

// API接口定义
export const authAPI = {
  // 用户认证
  login: (credentials) => api.post('/auth/login', credentials),
  register: (userData) => api.post('/auth/register', userData),
  logout: () => api.post('/auth/logout'),
  getProfile: () => api.get('/auth/profile'),
  updateProfile: (data) => api.put('/auth/profile', data)
}

export const roomAPI = {
  // 房间管理
  getRooms: () => api.get('/rooms'),
  createRoom: (roomData) => api.post('/rooms', roomData),
  getRoom: (roomId) => api.get(`/rooms/${roomId}`),
  joinRoom: (roomId) => api.post(`/rooms/${roomId}/join`),
  leaveRoom: (roomId) => api.post(`/rooms/${roomId}/leave`)
}

export const adminAPI = {
  // 管理员接口
  login: (credentials) => api.post('/admin/auth/login', credentials),
  getUsers: (params) => api.get('/admin/users', { params }),
  updateUser: (userId, data) => api.put(`/admin/users/${userId}`, data),
  getRooms: () => api.get('/admin/rooms'),
  getStats: () => api.get('/admin/stats')
}

export default api 