<template>
  <div class="home-container">
    <el-container>
      <el-header class="header">
        <div class="header-content">
          <h2>物业报修管理系统 - 管理后台</h2>
          <div class="user-info">
            <span>管理员：{{ username }}</span>
            <el-button type="danger" size="small" @click="handleLogout">退出登录</el-button>
          </div>
        </div>
      </el-header>
      <el-container>
        <el-aside width="200px" class="sidebar">
          <el-menu :default-active="activeMenu" @select="handleMenuSelect">
            <el-menu-item index="all">
              <span>全部报修</span>
            </el-menu-item>
            <el-menu-item index="pending">
              <span>待处理</span>
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
            <el-menu-item index="workers">
              <span>工人管理</span>
            </el-menu-item>
            <el-menu-item index="stats">
              <span>统计分析</span>
            </el-menu-item>
          </el-menu>
        </el-aside>
        <el-main class="main-content">
          <div class="toolbar">
            <h3>{{ currentTitle }}</h3>
            <el-button type="primary" @click="loadRepairs">刷新</el-button>
          </div>
          <el-table :data="filteredRepairs" style="width: 100%" v-loading="loading">
            <el-table-column prop="id" label="ID" width="80" />
            <el-table-column prop="username" label="用户名" width="120" />
            <el-table-column prop="repair_type" label="报修类型" width="150" />
            <el-table-column prop="description" label="报修描述" />
            <el-table-column prop="worker_name" label="派工" width="120" />
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
            <el-table-column label="操作" width="250" fixed="right">
              <template #default="scope">
                <el-button 
                  v-if="scope.row.status === '未处理'" 
                  size="small" 
                  type="primary" 
                  @click="handleAssign(scope.row)"
                >
                  派单
                </el-button>
                <el-select
                  v-model="scope.row.status"
                  placeholder="更新状态"
                  @change="handleStatusChange(scope.row)"
                  style="width: 130px"
                >
                  <el-option label="未处理" value="未处理" />
                  <el-option label="已派单" value="已派单" />
                  <el-option label="处理中" value="处理中" />
                  <el-option label="已完成" value="已完成" />
                </el-select>
              </template>
            </el-table-column>
          </el-table>
          <el-empty v-if="!loading && filteredRepairs.length === 0" description="暂无报修记录" />
        </el-main>
      </el-container>
    </el-container>

    <el-dialog title="派单给工人" v-model="showAssignDialog" width="450px">
      <el-form :model="assignForm" ref="assignFormRef">
        <el-form-item label="工单信息">
          <p>报修类型：{{ selectedRepair?.repair_type }}</p>
          <p>报修描述：{{ selectedRepair?.description }}</p>
        </el-form-item>
        <el-form-item label="选择工人" prop="workerId">
          <el-select v-model="assignForm.workerId" placeholder="请选择工人">
            <el-option 
              v-for="worker in availableWorkers" 
              :key="worker.id" 
              :label="worker.name + ' (' + worker.skill_type + ')'" 
              :value="worker.id" 
            />
          </el-select>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="showAssignDialog = false">取消</el-button>
        <el-button type="primary" @click="handleConfirmAssign">确认派单</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { repairAPI, workerAPI } from '@/api'

const router = useRouter()
const username = ref(localStorage.getItem('username') || '')
const repairs = ref([])
const workers = ref([])
const loading = ref(false)
const activeMenu = ref('all')
const showAssignDialog = ref(false)
const selectedRepair = ref(null)
const assignForm = ref({
  workerId: ''
})
const assignFormRef = ref(null)

const availableWorkers = computed(() => {
  return workers.value.filter(w => w.status === 1)
})

const currentTitle = computed(() => {
  const titleMap = {
    'all': '全部报修记录',
    'pending': '待处理报修',
    'assigned': '已派单报修',
    'processing': '处理中报修',
    'completed': '已完成报修'
  }
  return titleMap[activeMenu.value] || '报修记录'
})

const filteredRepairs = computed(() => {
  if (activeMenu.value === 'all') {
    return repairs.value
  }
  const statusMap = {
    'pending': '未处理',
    'assigned': '已派单',
    'processing': '处理中',
    'completed': '已完成'
  }
  return repairs.value.filter(repair => repair.status === statusMap[activeMenu.value])
})

const loadRepairs = async () => {
  loading.value = true
  try {
    const response = await repairAPI.getAllRepairs()
    if (response.success) {
      repairs.value = response.data || []
    }
  } catch (error) {
    ElMessage.error('加载报修记录失败')
  } finally {
    loading.value = false
  }
}

const loadWorkers = async () => {
  try {
    const response = await workerAPI.getWorkers()
    if (response.success) {
      workers.value = response.data || []
    }
  } catch (error) {
    console.error('加载工人列表失败')
  }
}

const handleStatusChange = async (repair) => {
  try {
    const response = await repairAPI.updateRepairStatus(repair.id, repair.status)
    if (response.success) {
      ElMessage.success('状态更新成功')
    } else {
      ElMessage.error(response.error || '更新失败')
      loadRepairs()
    }
  } catch (error) {
    ElMessage.error('更新失败，请重试')
    loadRepairs()
  }
}

const handleAssign = (repair) => {
  selectedRepair.value = repair
  assignForm.value.workerId = ''
  showAssignDialog.value = true
}

const handleConfirmAssign = async () => {
  if (!assignForm.value.workerId) {
    ElMessage.error('请选择工人')
    return
  }
  
  try {
    const response = await repairAPI.assignWorker(selectedRepair.value.id, assignForm.value.workerId)
    if (response.success) {
      ElMessage.success('派单成功')
      showAssignDialog.value = false
      loadRepairs()
    } else {
      ElMessage.error(response.error || '派单失败')
    }
  } catch (error) {
    ElMessage.error(error.response?.data?.error || '派单失败')
  }
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
  activeMenu.value = index
  if (index === 'workers') {
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
  loadRepairs()
  loadWorkers()
})
</script>

<style scoped>
.home-container {
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
