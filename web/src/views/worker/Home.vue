<template>
  <div class="home-container">
    <el-container>
      <el-header class="header">
        <div class="header-content">
          <h2>工人工作台</h2>
          <div class="user-info">
            <span>{{ workerName }} ({{ skillType }})</span>
            <el-button type="danger" size="small" @click="handleLogout">退出登录</el-button>
          </div>
        </div>
      </el-header>
      <el-container>
        <el-aside width="200px" class="sidebar">
          <el-menu :default-active="activeMenu" @select="handleMenuSelect">
            <el-menu-item index="all">
              <span>全部工单</span>
            </el-menu-item>
            <el-menu-item index="assigned">
              <span>已派单</span>
            </el-menu-item>
            <el-menu-item index="processing">
              <span>处理中</span>
            </el-menu-item>
            <el-menu-item index="completed">
              <span>已完成</span>
            </el-menu-item>
            <el-menu-item index="profile">
              <span>个人中心</span>
            </el-menu-item>
          </el-menu>
        </el-aside>
        <el-main class="main-content">
          <div class="toolbar">
            <h3>{{ currentTitle }}</h3>
            <el-button type="primary" @click="loadOrders">刷新</el-button>
          </div>
          <el-table :data="filteredOrders" style="width: 100%" v-loading="loading" @row-click="handleRowClick">
            <el-table-column prop="id" label="工单ID" width="100" />
            <el-table-column prop="username" label="报修用户" width="120" />
            <el-table-column prop="repair_type" label="报修类型" width="150" />
            <el-table-column prop="description" label="报修描述" />
            <el-table-column prop="status" label="状态" width="120">
              <template #default="scope">
                <el-tag :type="getStatusType(scope.row.status)">
                  {{ scope.row.status }}
                </el-tag>
              </template>
            </el-table-column>
            <el-table-column prop="created_at" label="提交时间" width="180">
              <template #default="scope">
                {{ formatDate(scope.row.created_at) }}
              </template>
            </el-table-column>
            <el-table-column label="操作" width="200" fixed="right">
              <template #default="scope">
                <el-button v-if="scope.row.status === '已派单'" size="small" type="primary" @click.stop="handleAccept(scope.row)">接单</el-button>
                <el-button v-if="scope.row.status === '已派单'" size="small" type="danger" @click.stop="handleReject(scope.row)">拒单</el-button>
                <el-button v-if="scope.row.status === '处理中'" size="small" type="success" @click.stop="goToDetail(scope.row.id)">提交结果</el-button>
                <el-button v-if="scope.row.status === '已完成'" size="small" type="info" @click.stop="goToDetail(scope.row.id)">查看详情</el-button>
              </template>
            </el-table-column>
          </el-table>
          <el-empty v-if="!loading && filteredOrders.length === 0" description="暂无工单" />
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
const workerId = ref(localStorage.getItem('workerId') || '')
const workerName = ref(localStorage.getItem('workerName') || '')
const skillType = ref(localStorage.getItem('workerSkillType') || '')
const orders = ref([])
const loading = ref(false)
const activeMenu = ref('all')

const currentTitle = computed(() => {
  const titleMap = {
    'all': '全部工单',
    'assigned': '已派单',
    'processing': '处理中',
    'completed': '已完成',
    'profile': '个人中心'
  }
  return titleMap[activeMenu.value] || '工单列表'
})

const filteredOrders = computed(() => {
  if (activeMenu.value === 'all') {
    return orders.value
  }
  const statusMap = {
    'assigned': '已派单',
    'processing': '处理中',
    'completed': '已完成'
  }
  return orders.value.filter(order => order.status === statusMap[activeMenu.value])
})

const loadOrders = async () => {
  if (!workerId.value) return

  loading.value = true
  try {
    const response = await repairAPI.getWorkerRepairs(workerId.value)
    if (response.success) {
      orders.value = response.data || []
    }
  } catch (error) {
    ElMessage.error('加载工单失败')
  } finally {
    loading.value = false
  }
}

const handleAccept = async (order) => {
  try {
    const response = await repairAPI.acceptOrder(order.id, workerId.value)
    if (response.success) {
      ElMessage.success('接单成功')
      loadOrders()
    } else {
      ElMessage.error(response.error || '接单失败')
    }
  } catch (error) {
    ElMessage.error(error.response?.data?.error || '接单失败')
  }
}

const handleReject = async (order) => {
  try {
    const response = await repairAPI.rejectOrder(order.id, workerId.value)
    if (response.success) {
      ElMessage.success(response.message)
      loadOrders()
    } else {
      ElMessage.error(response.error || '拒单失败')
    }
  } catch (error) {
    ElMessage.error(error.response?.data?.error || '拒单失败')
  }
}

const handleRowClick = (row) => {
  if (row.status === '已完成') {
    goToDetail(row.id)
  }
}

const goToDetail = (id) => {
  router.push(`/worker/order/${id}`)
}

const getStatusType = (status) => {
  const statusMap = {
    '未处理': 'warning',
    '已派单': 'info',
    '处理中': 'primary',
    '已完成': 'success'
  }
  return statusMap[status] || 'info'
}

const formatDate = (dateStr) => {
  if (!dateStr) return ''
  const date = new Date(dateStr)
  return date.toLocaleString('zh-CN')
}

const handleMenuSelect = (index) => {
  if (index === 'profile') {
    router.push('/worker/profile')
    return
  }
  activeMenu.value = index
}

const handleLogout = () => {
  localStorage.clear()
  ElMessage.success('已退出登录')
  router.push('/worker/login')
}

onMounted(() => {
  if (!workerId.value) {
    router.push('/worker/login')
    return
  }
  loadOrders()
})
</script>

<style scoped>
.home-container {
  min-height: 100vh;
}

.header {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
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

.toolbar {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}

h3 {
  margin: 0;
  color: #333;
}
</style>