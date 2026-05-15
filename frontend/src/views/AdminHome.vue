<template>
  <div class="home-container">
    <header class="header">
      <h1>物业报修管理系统 - 管理员</h1>
      <div class="user-info">
        <span>欢迎, {{ user?.username }}</span>
        <el-button @click="handleLogout">退出登录</el-button>
      </div>
    </header>

    <div class="content">
      <div class="stats-row">
        <div class="stat-card warning">
          <div class="stat-icon">📋</div>
          <div class="stat-info">
            <div class="stat-value">{{ pendingCount }}</div>
            <div class="stat-label">未处理</div>
          </div>
        </div>
        <div class="stat-card info">
          <div class="stat-icon">🔄</div>
          <div class="stat-info">
            <div class="stat-value">{{ processingCount }}</div>
            <div class="stat-label">处理中</div>
          </div>
        </div>
        <div class="stat-card success">
          <div class="stat-icon">✅</div>
          <div class="stat-info">
            <div class="stat-value">{{ completedCount }}</div>
            <div class="stat-label">已完成</div>
          </div>
        </div>
      </div>

      <el-table :data="workOrders" border>
        <el-table-column prop="id" label="报修编号"></el-table-column>
        <el-table-column prop="user_id" label="用户ID"></el-table-column>
        <el-table-column prop="type" label="报修类型"></el-table-column>
        <el-table-column prop="description" label="报修描述"></el-table-column>
        <el-table-column prop="images" label="报修图片">
          <template #default="scope">
            <img v-if="scope.row.images" :src="scope.row.images.split(',')[0]" width="50" height="50" />
            <span v-else>无</span>
          </template>
        </el-table-column>
        <el-table-column prop="status" label="状态">
          <template #default="scope">
            <el-tag :type="getStatusType(scope.row.status)">{{ getStatusText(scope.row.status) }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="created_at" label="创建时间"></el-table-column>
        <el-table-column prop="updated_at" label="更新时间"></el-table-column>
        <el-table-column label="操作">
          <template #default="scope">
            <el-button-group>
              <el-button v-if="scope.row.status !== 1" @click="updateStatus(scope.row.id, 1)" size="small">设为未处理</el-button>
              <el-button v-if="scope.row.status !== 2" @click="updateStatus(scope.row.id, 2)" size="small">设为处理中</el-button>
              <el-button v-if="scope.row.status !== 3" @click="updateStatus(scope.row.id, 3)" size="small">设为已完成</el-button>
            </el-button-group>
          </template>
        </el-table-column>
      </el-table>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { getUser, removeUser } from '../utils/storage'
import { getWorkOrders, updateWorkOrderStatus } from '../api'

const router = useRouter()
const user = ref(getUser())
const workOrders = ref([])

const pendingCount = ref(0)
const processingCount = ref(0)
const completedCount = ref(0)

onMounted(() => {
  if (!user.value || user.value.role !== 2) {
    router.push('/admin/login')
    return
  }
  loadWorkOrders()
})

const loadWorkOrders = async () => {
  try {
    const res = await getWorkOrders()
    workOrders.value = res.data.data
    updateStats()
  } catch (err) {
    ElMessage.error('获取报修记录失败')
  }
}

const updateStats = () => {
  pendingCount.value = workOrders.value.filter(o => o.status === 1).length
  processingCount.value = workOrders.value.filter(o => o.status === 2).length
  completedCount.value = workOrders.value.filter(o => o.status === 3).length
}

const updateStatus = async (id, status) => {
  try {
    await updateWorkOrderStatus(id, status)
    ElMessage.success('状态更新成功')
    loadWorkOrders()
  } catch (err) {
    ElMessage.error('更新失败')
  }
}

const handleLogout = () => {
  removeUser()
  router.push('/admin/login')
}

const getStatusText = (status) => {
  const statusMap = {
    1: '未处理',
    2: '处理中',
    3: '已完成'
  }
  return statusMap[status] || '未知'
}

const getStatusType = (status) => {
  const typeMap = {
    1: 'warning',
    2: 'info',
    3: 'success'
  }
  return typeMap[status] || 'default'
}
</script>

<style scoped>
.home-container {
  min-height: 100vh;
  background: #f5f5f5;
}

.header {
  background: linear-gradient(135deg, #11998e 0%, #38ef7d 100%);
  padding: 20px 30px;
  display: flex;
  justify-content: space-between;
  align-items: center;
  color: white;
}

.header h1 {
  margin: 0;
}

.user-info {
  display: flex;
  align-items: center;
  gap: 20px;
}

.content {
  padding: 20px;
}

.stats-row {
  display: flex;
  gap: 20px;
  margin-bottom: 20px;
}

.stat-card {
  flex: 1;
  background: white;
  padding: 20px;
  border-radius: 10px;
  display: flex;
  align-items: center;
  gap: 15px;
  box-shadow: 0 0 10px rgba(0, 0, 0, 0.05);
}

.stat-icon {
  font-size: 40px;
}

.stat-info {
  flex: 1;
}

.stat-value {
  font-size: 24px;
  font-weight: bold;
}

.stat-label {
  font-size: 14px;
  color: #666;
}

.stat-card.warning .stat-value { color: #e6a23c; }
.stat-card.info .stat-value { color: #67c23a; }
.stat-card.success .stat-value { color: #409eff; }
</style>