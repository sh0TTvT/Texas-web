<!--
管理员仪表盘页面
作用：显示系统概览、统计数据和关键指标
-->

<template>
  <div class="dashboard-page">
    <!-- 页面标题 -->
    <div class="page-header">
      <h4 class="page-title">系统仪表盘</h4>
      <p class="page-subtitle">实时监控系统运行状态和关键数据</p>
    </div>

    <!-- 统计卡片 -->
    <div class="stats-grid">
      <!-- 在线用户 -->
      <q-card class="stat-card">
        <q-card-section>
          <div class="stat-content">
            <div class="stat-icon">
              <q-icon name="people" size="32px" color="primary" />
            </div>
            <div class="stat-details">
              <div class="stat-value">{{ stats.onlineUsers }}</div>
              <div class="stat-label">在线用户</div>
              <div class="stat-change positive">
                <q-icon name="trending_up" size="sm" />
                +12%
              </div>
            </div>
          </div>
        </q-card-section>
      </q-card>

      <!-- 活跃房间 -->
      <q-card class="stat-card">
        <q-card-section>
          <div class="stat-content">
            <div class="stat-icon">
              <q-icon name="meeting_room" size="32px" color="orange" />
            </div>
            <div class="stat-details">
              <div class="stat-value">{{ stats.activeRooms }}</div>
              <div class="stat-label">活跃房间</div>
              <div class="stat-change positive">
                <q-icon name="trending_up" size="sm" />
                +8%
              </div>
            </div>
          </div>
        </q-card-section>
      </q-card>

      <!-- 游戏局数 -->
      <q-card class="stat-card">
        <q-card-section>
          <div class="stat-content">
            <div class="stat-icon">
              <q-icon name="casino" size="32px" color="green" />
            </div>
            <div class="stat-details">
              <div class="stat-value">{{ stats.totalGames }}</div>
              <div class="stat-label">今日游戏</div>
              <div class="stat-change positive">
                <q-icon name="trending_up" size="sm" />
                +25%
              </div>
            </div>
          </div>
        </q-card-section>
      </q-card>

      <!-- 筹码流通 -->
      <q-card class="stat-card">
        <q-card-section>
          <div class="stat-content">
            <div class="stat-icon">
              <q-icon name="account_balance" size="32px" color="purple" />
            </div>
            <div class="stat-details">
              <div class="stat-value">{{ formatNumber(stats.chipsCirculation) }}</div>
              <div class="stat-label">筹码流通</div>
              <div class="stat-change positive">
                <q-icon name="trending_up" size="sm" />
                +15%
              </div>
            </div>
          </div>
        </q-card-section>
      </q-card>
    </div>

    <!-- 图表和数据 -->
    <div class="charts-grid">
      <!-- 用户活跃度图表 -->
      <q-card class="chart-card">
        <q-card-section>
          <div class="chart-header">
            <h6 class="chart-title">用户活跃度</h6>
            <q-btn flat icon="refresh" @click="refreshCharts" />
          </div>
          <div class="chart-placeholder">
            <q-icon name="show_chart" size="48px" color="grey-5" />
            <p class="placeholder-text">图表数据加载中...</p>
          </div>
        </q-card-section>
      </q-card>

      <!-- 游戏分布图表 -->
      <q-card class="chart-card">
        <q-card-section>
          <div class="chart-header">
            <h6 class="chart-title">游戏分布</h6>
            <q-btn flat icon="more_vert">
              <q-menu>
                <q-list style="min-width: 100px">
                  <q-item clickable>
                    <q-item-section>导出数据</q-item-section>
                  </q-item>
                  <q-item clickable>
                    <q-item-section>打印图表</q-item-section>
                  </q-item>
                </q-list>
              </q-menu>
            </q-btn>
          </div>
          <div class="chart-placeholder">
            <q-icon name="pie_chart" size="48px" color="grey-5" />
            <p class="placeholder-text">饼图数据加载中...</p>
          </div>
        </q-card-section>
      </q-card>
    </div>

    <!-- 最近活动 -->
    <div class="activity-section">
      <q-card>
        <q-card-section>
          <div class="section-header">
            <h6 class="section-title">最近活动</h6>
            <q-btn flat label="查看全部" color="primary" />
          </div>
          
          <q-list>
            <q-item
              v-for="(activity, index) in recentActivities"
              :key="index"
              class="activity-item"
            >
              <q-item-section avatar>
                <q-avatar :color="activity.color" text-color="white">
                  <q-icon :name="activity.icon" />
                </q-avatar>
              </q-item-section>
              
              <q-item-section>
                <q-item-label>{{ activity.title }}</q-item-label>
                <q-item-label caption>{{ activity.description }}</q-item-label>
              </q-item-section>
              
              <q-item-section side>
                <q-item-label caption>{{ activity.time }}</q-item-label>
              </q-item-section>
            </q-item>
          </q-list>
        </q-card-section>
      </q-card>
    </div>

    <!-- 系统状态 -->
    <div class="system-status">
      <q-card>
        <q-card-section>
          <h6 class="section-title">系统状态</h6>
          
          <div class="status-grid">
            <div class="status-item">
              <div class="status-label">数据库连接</div>
              <q-chip color="green" text-color="white" icon="check_circle">
                正常
              </q-chip>
            </div>
            
            <div class="status-item">
              <div class="status-label">Redis缓存</div>
              <q-chip color="green" text-color="white" icon="check_circle">
                正常
              </q-chip>
            </div>
            
            <div class="status-item">
              <div class="status-label">WebSocket服务</div>
              <q-chip color="green" text-color="white" icon="check_circle">
                正常
              </q-chip>
            </div>
            
            <div class="status-item">
              <div class="status-label">服务器负载</div>
              <q-chip color="orange" text-color="white" icon="warning">
                中等
              </q-chip>
            </div>
          </div>
        </q-card-section>
      </q-card>
    </div>
  </div>
</template>

<script>
import { defineComponent, ref, onMounted } from 'vue'
import { adminAPI } from '../../api'

export default defineComponent({
  name: 'DashboardPage',
  
  setup() {
    // 响应式数据
    const stats = ref({
      onlineUsers: 0,
      activeRooms: 0,
      totalGames: 0,
      chipsCirculation: 0
    })
    
    const recentActivities = ref([
      {
        icon: 'person_add',
        color: 'green',
        title: '新用户注册',
        description: '用户 "player123" 完成注册',
        time: '2分钟前'
      },
      {
        icon: 'casino',
        color: 'blue',
        title: '游戏开始',
        description: '房间 "高级场1" 开始新局游戏',
        time: '5分钟前'
      },
      {
        icon: 'warning',
        color: 'orange',
        title: '系统警告',
        description: '服务器负载达到80%',
        time: '10分钟前'
      },
      {
        icon: 'account_balance',
        color: 'purple',
        title: '大额交易',
        description: '玩家 "highroller" 获得50,000筹码',
        time: '15分钟前'
      }
    ])
    
    // 方法
    const formatNumber = (num) => {
      if (num >= 1000000) {
        return (num / 1000000).toFixed(1) + 'M'
      } else if (num >= 1000) {
        return (num / 1000).toFixed(1) + 'K'
      }
      return num.toString()
    }
    
    const loadStats = async () => {
      try {
        const response = await adminAPI.getStats()
        stats.value = response.data.stats
      } catch (error) {
        console.error('Failed to load stats:', error)
        // 使用模拟数据
        stats.value = {
          onlineUsers: 127,
          activeRooms: 8,
          totalGames: 245,
          chipsCirculation: 1234567
        }
      }
    }
    
    const refreshCharts = () => {
      // TODO: 实现图表刷新逻辑
      console.log('Refreshing charts...')
    }
    
    // 生命周期
    onMounted(() => {
      loadStats()
      
      // 设置定时刷新
      const interval = setInterval(loadStats, 30000) // 30秒刷新一次
      
      // 清理定时器
      return () => clearInterval(interval)
    })
    
    return {
      stats,
      recentActivities,
      formatNumber,
      refreshCharts
    }
  }
})
</script>

<style scoped lang="scss">
.dashboard-page {
  max-width: 1400px;
  margin: 0 auto;
}

.page-header {
  margin-bottom: 2rem;
  
  .page-title {
    margin: 0 0 0.5rem 0;
    color: #1976d2;
    font-weight: 600;
  }
  
  .page-subtitle {
    margin: 0;
    color: #666;
    font-size: 1rem;
  }
}

.stats-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
  gap: 1.5rem;
  margin-bottom: 2rem;
}

.stat-card {
  .stat-content {
    display: flex;
    align-items: center;
    gap: 1rem;
    
    .stat-icon {
      flex-shrink: 0;
    }
    
    .stat-details {
      flex: 1;
      
      .stat-value {
        font-size: 2rem;
        font-weight: 700;
        color: #333;
        line-height: 1;
        margin-bottom: 0.25rem;
      }
      
      .stat-label {
        font-size: 0.9rem;
        color: #666;
        margin-bottom: 0.5rem;
      }
      
      .stat-change {
        display: flex;
        align-items: center;
        gap: 0.25rem;
        font-size: 0.8rem;
        font-weight: 500;
        
        &.positive {
          color: #4caf50;
        }
        
        &.negative {
          color: #f44336;
        }
      }
    }
  }
}

.charts-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(400px, 1fr));
  gap: 1.5rem;
  margin-bottom: 2rem;
}

.chart-card {
  .chart-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 1rem;
    
    .chart-title {
      margin: 0;
      color: #333;
      font-weight: 600;
    }
  }
  
  .chart-placeholder {
    height: 200px;
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    background: #f5f5f5;
    border-radius: 8px;
    
    .placeholder-text {
      margin: 0.5rem 0 0 0;
      color: #666;
      font-size: 0.9rem;
    }
  }
}

.activity-section {
  margin-bottom: 2rem;
  
  .section-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 1rem;
    
    .section-title {
      margin: 0;
      color: #333;
      font-weight: 600;
    }
  }
  
  .activity-item {
    border-bottom: 1px solid rgba(0, 0, 0, 0.12);
    
    &:last-child {
      border-bottom: none;
    }
  }
}

.system-status {
  .section-title {
    margin: 0 0 1rem 0;
    color: #333;
    font-weight: 600;
  }
  
  .status-grid {
    display: grid;
    grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
    gap: 1rem;
  }
  
  .status-item {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 0.75rem;
    background: #f8f9fa;
    border-radius: 8px;
    
    .status-label {
      font-size: 0.9rem;
      color: #333;
      font-weight: 500;
    }
  }
}

// 响应式设计
@media (max-width: 768px) {
  .stats-grid {
    grid-template-columns: 1fr;
    gap: 1rem;
  }
  
  .charts-grid {
    grid-template-columns: 1fr;
    gap: 1rem;
  }
  
  .chart-card .chart-placeholder {
    height: 150px;
  }
  
  .status-grid {
    grid-template-columns: 1fr;
  }
  
  .section-header {
    flex-direction: column;
    align-items: flex-start;
    gap: 0.5rem;
  }
}

@media (max-width: 480px) {
  .stat-card .stat-content {
    flex-direction: column;
    text-align: center;
    gap: 0.75rem;
  }
  
  .activity-item {
    .q-item__section--avatar {
      margin-right: 0.5rem;
    }
  }
}

// 动画效果
.stat-card {
  transition: transform 0.2s ease, box-shadow 0.2s ease;
  
  &:hover {
    transform: translateY(-2px);
    box-shadow: 0 4px 20px rgba(0, 0, 0, 0.15);
  }
}

.chart-card {
  transition: transform 0.2s ease;
  
  &:hover {
    transform: translateY(-1px);
  }
}
</style> 