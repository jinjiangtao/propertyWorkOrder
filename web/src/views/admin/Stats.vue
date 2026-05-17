<template>
  <div class="stats-container">
    <el-container>
      <el-header class="header">
        <div class="header-content">
          <h2>物业报修管理系统 - 统计分析</h2>
          <div class="user-info">
            <span>管理员：{{ username }}</span>
            <el-button type="danger" size="small" @click="handleLogout">退出登录</el-button>
          </div>
        </div>
      </el-header>
      <el-container>
        <el-aside width="200px" class="sidebar">
          <el-menu :default-active="activeMenu" @select="handleMenuSelect">
            <el-menu-item index="home">
              <span>工单管理</span>
            </el-menu-item>
            <el-menu-item index="workers">
              <span>工人管理</span>
            </el-menu-item>
            <el-menu-item index="stats">
              <span>统计分析</span>
            </el-menu-item>
          </el-menu>
        </el-aside>
        <el-main class="main-content">
          <div class="stats-header">
            <h3>工人工作量统计</h3>
            <el-button type="primary" @click="loadStats">刷新数据</el-button>
          </div>
          
          <div class="stats-summary">
            <div class="summary-card">
              <div class="summary-value">{{ totalOrders }}</div>
              <div class="summary-label">总工单数</div>
            </div>
            <div class="summary-card">
              <div class="summary-value completed">{{ completedOrders }}</div>
              <div class="summary-label">已完成</div>
            </div>
            <div class="summary-card">
              <div class="summary-value processing">{{ processingOrders }}</div>
              <div class="summary-label">处理中</div>
            </div>
            <div class="summary-card">
              <div class="summary-value pending">{{ pendingOrders }}</div>
              <div class="summary-label">待处理</div>
            </div>
          </div>
          
          <div class="stats-table">
            <el-table :data="workerStats" style="width: 100%" v-loading="loading">
              <el-table-column prop="name" label="工人姓名" width="120" />
              <el-table-column prop="work_no" label="工号" width="120" />
              <el-table-column prop="skill_type" label="工种" width="120" />
              <el-table-column prop="total_count" label="派单总数" width="120" />
              <el-table-column prop="completed_count" label="已完成数" width="120" />
              <el-table-column label="完成率" width="150">
                <template #default="scope">
                  <el-progress 
                    :percentage="getCompletionRate(scope.row)" 
                    :color="getProgressColor(scope.row)"
                  />
                </template>
              </el-table-column>
            </el-table>
            <el-empty v-if="!loading && workerStats.length === 0" description="暂无统计数据" />
          </div>
        </el-main>
      </el-container>
    </el-container>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { repairAPI } from '@/api'

const router = useRouter()
const username = ref(localStorage.getItem('username') || '')
const workerStats = ref([])
const loading = ref(false)
const activeMenu = ref('stats')

const totalOrders = computed(() => {
  return workerStats.value.reduce((sum, stat) => sum + stat.total_count, 0)
})

const completedOrders = computed(() => {
  return workerStats.value.reduce((sum, stat) => sum + stat.completed_count, 0)
})

const processingOrders = ref(0)
const pendingOrders = ref(0)

const loadStats = async () => {
  loading.value = true
  try {
    const response = await repairAPI.getWorkerStats()
    if (response.success) {
      workerStats.value = response.data || []
    }
  } catch (error) {
    ElMessage.error('加载统计数据失败')
  } finally {
    loading.value = false
  }
}

const getCompletionRate = (row) => {
  if (row.total_count === 0) return 0
  return Math.round((row.completed_count / row.total_count) * 100)
}

const getProgressColor = (row) => {
  const rate = getCompletionRate(row)
  if (rate >= 80) return '#67c23a'
  if (rate >= 50) return '#e6a23c'
  return '#f56c6c'
}

const handleMenuSelect = (index) => {
  activeMenu.value = index
  if (index === 'home') {
    router.push('/admin/home')
  } else if (index === 'workers') {
    router.push('/admin/workers')
  } else if (index === 'stats') {
    router.push('/admin/stats')
  }
}

const handleLogout = () => {
  localStorage.clear()
  ElMessage.success('已退出登录')
  router.push('/admin/login')
}

onMounted(() => {
  if (!localStorage.getItem('userId')) {
    router.push('/admin/login')
    return
  }
  loadStats()
})
</script>

<style scoped>
.stats-container {
  min-height: 100vh;
}

.header {
  background: #f5576c;
  color: white;
  display: flex;
  align-items: center;
  padding: 0 20px;
}

.header-content {
  width: 100%;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.header h2 {
  margin: 0;
  font-size: 20px;
}

.user-info {
  display: flex;
  align-items: center;
  gap: 15px;
}

.sidebar {
  background: #f5f7fa;
  min-height: calc(100vh - 60px);
}

.main-content {
  padding: 20px;
  background: white;
}

.stats-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}

h3 {
  margin: 0;
  color: #333;
}

.stats-summary {
  display: flex;
  gap: 20px;
  margin-bottom: 20px;
}

.summary-card {
  flex: 1;
  background: #f8f9fa;
  padding: 20px;
  border-radius: 12px;
  text-align: center;
}

.summary-value {
  font-size: 32px;
  font-weight: bold;
  color: #333;
}

.summary-value.completed {
  color: #67c23a;
}

.summary-value.processing {
  color: #409eff;
}

.summary-value.pending {
  color: #e6a23c;
}

.summary-label {
  font-size: 14px;
  color: #999;
  margin-top: 8px;
}

.stats-table {
  background: white;
  border-radius: 12px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.06);
}
</style>