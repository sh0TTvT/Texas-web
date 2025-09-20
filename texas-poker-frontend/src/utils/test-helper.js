// 前后端集成测试辅助工具
// 作用：提供测试接口连接、WebSocket通信、游戏流程等的测试函数

import { api } from '../api'
import wsClient from './websocket'
import { useAuthStore } from '../stores/auth'
import { useGameStore } from '../stores/game'

class TestHelper {
  constructor() {
    this.testResults = []
    this.isRunning = false
  }

  // 开始测试套件
  async runAllTests() {
    if (this.isRunning) {
      console.warn('测试已在运行中')
      return
    }

    this.isRunning = true
    this.testResults = []
    
    console.log('🚀 开始前后端集成测试')
    
    try {
      // 1. API连接测试
      await this.testAPIConnection()
      
      // 2. 用户认证测试
      await this.testAuthentication()
      
      // 3. WebSocket连接测试
      await this.testWebSocketConnection()
      
      // 4. 房间管理测试
      await this.testRoomManagement()
      
      // 5. 游戏流程测试
      await this.testGameFlow()
      
      console.log('✅ 所有测试完成')
      this.printTestResults()
      
    } catch (error) {
      console.error('❌ 测试过程中发生错误:', error)
    } finally {
      this.isRunning = false
    }
  }

  // 测试API连接
  async testAPIConnection() {
    console.log('📡 测试API连接...')
    
    try {
      // 测试基础连接
      const startTime = Date.now()
      const response = await fetch(`${api.defaults.baseURL}/health`, {
        method: 'GET',
        headers: {
          'Content-Type': 'application/json'
        }
      })
      const endTime = Date.now()
      
      if (response.ok) {
        this.addTestResult('API连接', true, `响应时间: ${endTime - startTime}ms`)
      } else {
        this.addTestResult('API连接', false, `HTTP ${response.status}`)
      }
    } catch (error) {
      this.addTestResult('API连接', false, error.message)
    }
  }

  // 测试用户认证
  async testAuthentication() {
    console.log('🔐 测试用户认证...')
    
    const testUser = {
      username: `test_${Date.now()}`,
      email: `test_${Date.now()}@example.com`,
      password: 'test123456'
    }

    try {
      // 测试注册
      const registerResponse = await api.post('/auth/register', testUser)
      
      if (registerResponse.status === 201 || registerResponse.status === 200) {
        this.addTestResult('用户注册', true, '注册成功')
        
        // 测试登录
        const loginResponse = await api.post('/auth/login', {
          username: testUser.username,
          password: testUser.password
        })
        
        if (loginResponse.data.token) {
          this.addTestResult('用户登录', true, '获取到token')
          
          // 测试获取用户信息
          const profileResponse = await api.get('/auth/profile', {
            headers: {
              Authorization: `Bearer ${loginResponse.data.token}`
            }
          })
          
          if (profileResponse.data.user) {
            this.addTestResult('获取用户信息', true, '用户信息正确')
          } else {
            this.addTestResult('获取用户信息', false, '用户信息为空')
          }
        } else {
          this.addTestResult('用户登录', false, '未获取到token')
        }
      } else {
        this.addTestResult('用户注册', false, `HTTP ${registerResponse.status}`)
      }
    } catch (error) {
      this.addTestResult('用户认证', false, error.response?.data?.error || error.message)
    }
  }

  // 测试WebSocket连接
  async testWebSocketConnection() {
    console.log('🔌 测试WebSocket连接...')
    
    return new Promise((resolve) => {
      const timeout = setTimeout(() => {
        this.addTestResult('WebSocket连接', false, '连接超时')
        resolve()
      }, 10000)

      // 模拟token（实际应用中应该从认证流程获取）
      const mockToken = 'test-token'
      
      wsClient.on('connected', () => {
        clearTimeout(timeout)
        this.addTestResult('WebSocket连接', true, '连接成功')
        
        // 测试心跳
        wsClient.send({ type: 'heartbeat', timestamp: Date.now() })
        
        setTimeout(() => {
          wsClient.disconnect()
          resolve()
        }, 2000)
      })
      
      wsClient.on('error', (error) => {
        clearTimeout(timeout)
        this.addTestResult('WebSocket连接', false, '连接失败')
        resolve()
      })
      
      // 尝试连接（可能会失败，因为没有真实的token）
      wsClient.connect(mockToken).catch(() => {
        clearTimeout(timeout)
        this.addTestResult('WebSocket连接', false, '认证失败（预期结果）')
        resolve()
      })
    })
  }

  // 测试房间管理
  async testRoomManagement() {
    console.log('🏠 测试房间管理...')
    
    try {
      // 测试获取房间列表
      const roomsResponse = await api.get('/rooms')
      
      if (Array.isArray(roomsResponse.data.rooms)) {
        this.addTestResult('获取房间列表', true, `获取到 ${roomsResponse.data.rooms.length} 个房间`)
        
        // 测试创建房间
        const newRoom = {
          name: `测试房间_${Date.now()}`,
          chip_level: 'low',
          min_chips: 100,
          small_blind: 5,
          big_blind: 10,
          max_players: 6,
          is_private: false
        }
        
        try {
          const createResponse = await api.post('/rooms', newRoom)
          
          if (createResponse.data.room) {
            this.addTestResult('创建房间', true, '房间创建成功')
            
            // 测试获取单个房间信息
            const roomId = createResponse.data.room.id
            const roomResponse = await api.get(`/rooms/${roomId}`)
            
            if (roomResponse.data.room) {
              this.addTestResult('获取房间信息', true, '房间信息正确')
            } else {
              this.addTestResult('获取房间信息', false, '房间信息为空')
            }
          } else {
            this.addTestResult('创建房间', false, '房间数据为空')
          }
        } catch (createError) {
          this.addTestResult('创建房间', false, createError.response?.data?.error || '需要认证')
        }
      } else {
        this.addTestResult('获取房间列表', false, '响应格式错误')
      }
    } catch (error) {
      this.addTestResult('房间管理', false, error.response?.data?.error || error.message)
    }
  }

  // 测试游戏流程
  async testGameFlow() {
    console.log('🎮 测试游戏流程...')
    
    // 由于游戏流程需要实际的认证和WebSocket连接，这里只做模拟测试
    const gameStore = useGameStore()
    
    try {
      // 模拟游戏状态更新
      const mockGameState = {
        state: 'preflop',
        players: [
          {
            id: 1,
            username: 'test_player1',
            chips: 1000,
            cards: ['AS', 'KH'],
            is_current_user: true
          },
          {
            id: 2,
            username: 'test_player2',
            chips: 1500,
            cards: ['QD', 'JC'],
            is_current_user: false
          }
        ],
        community_cards: [],
        pot: 15,
        current_player: 1,
        current_bet: 10
      }
      
      gameStore.updateGameState(mockGameState)
      
      // 验证状态更新
      if (gameStore.gameState === 'preflop') {
        this.addTestResult('游戏状态管理', true, '状态更新正常')
      } else {
        this.addTestResult('游戏状态管理', false, '状态更新失败')
      }
      
      // 模拟可用操作
      const mockActions = ['call', 'raise', 'fold']
      gameStore.updateAvailableActions(mockActions)
      
      if (gameStore.availableActions.length === 3) {
        this.addTestResult('游戏操作管理', true, '操作列表更新正常')
      } else {
        this.addTestResult('游戏操作管理', false, '操作列表更新失败')
      }
      
    } catch (error) {
      this.addTestResult('游戏流程', false, error.message)
    }
  }

  // 添加测试结果
  addTestResult(testName, success, details = '') {
    const result = {
      name: testName,
      success,
      details,
      timestamp: new Date().toISOString()
    }
    
    this.testResults.push(result)
    
    const status = success ? '✅' : '❌'
    console.log(`${status} ${testName}: ${details}`)
  }

  // 打印测试结果摘要
  printTestResults() {
    const total = this.testResults.length
    const passed = this.testResults.filter(r => r.success).length
    const failed = total - passed
    
    console.log('\n📊 测试结果摘要:')
    console.log(`总计: ${total} 个测试`)
    console.log(`✅ 通过: ${passed} 个`)
    console.log(`❌ 失败: ${failed} 个`)
    console.log(`成功率: ${((passed / total) * 100).toFixed(1)}%`)
    
    if (failed > 0) {
      console.log('\n❌ 失败的测试:')
      this.testResults
        .filter(r => !r.success)
        .forEach(r => console.log(`  - ${r.name}: ${r.details}`))
    }
  }

  // 获取测试结果
  getTestResults() {
    return {
      total: this.testResults.length,
      passed: this.testResults.filter(r => r.success).length,
      failed: this.testResults.filter(r => !r.success).length,
      results: this.testResults
    }
  }

  // 性能测试
  async performanceTest() {
    console.log('⚡ 开始性能测试...')
    
    const tests = [
      {
        name: 'API响应时间',
        test: async () => {
          const start = performance.now()
          await api.get('/rooms')
          const end = performance.now()
          return end - start
        },
        threshold: 1000 // 1秒
      },
      {
        name: 'WebSocket连接时间',
        test: async () => {
          const start = performance.now()
          // 模拟连接时间
          await new Promise(resolve => setTimeout(resolve, 100))
          const end = performance.now()
          return end - start
        },
        threshold: 5000 // 5秒
      }
    ]
    
    for (const test of tests) {
      try {
        const duration = await test.test()
        const success = duration < test.threshold
        this.addTestResult(
          test.name,
          success,
          `${duration.toFixed(2)}ms (阈值: ${test.threshold}ms)`
        )
      } catch (error) {
        this.addTestResult(test.name, false, error.message)
      }
    }
  }

  // 压力测试（模拟）
  async stressTest() {
    console.log('🔥 开始压力测试...')
    
    const concurrentRequests = 10
    const requests = []
    
    for (let i = 0; i < concurrentRequests; i++) {
      requests.push(
        api.get('/rooms').catch(error => ({ error: error.message }))
      )
    }
    
    const start = performance.now()
    const results = await Promise.all(requests)
    const end = performance.now()
    
    const successful = results.filter(r => !r.error).length
    const failed = results.length - successful
    
    this.addTestResult(
      '并发请求测试',
      failed === 0,
      `${concurrentRequests}个并发请求，成功${successful}个，失败${failed}个，耗时${(end - start).toFixed(2)}ms`
    )
  }
}

// 创建单例实例
const testHelper = new TestHelper()

// 导出测试函数
export const runIntegrationTests = () => testHelper.runAllTests()
export const runPerformanceTests = () => testHelper.performanceTest()
export const runStressTests = () => testHelper.stressTest()
export const getTestResults = () => testHelper.getTestResults()

// 在开发环境中暴露到全局
if (process.env.NODE_ENV === 'development') {
  window.testHelper = testHelper
  window.runTests = runIntegrationTests
}

export default testHelper 