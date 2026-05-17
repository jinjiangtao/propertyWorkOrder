<template>
  <div class="home-container">
    <el-container>
      <el-header class="header">
        <div class="header-content">
          <h2>物业报修管理系统 - 用户中心</h2>
          <div class="user-info">
            <span>欢迎，{{ username }}</span>
            <el-button type="danger" size="small" @click="handleLogout">退出登录</el-button>
          </div>
        </div>
      </el-header>
      <el-container>
        <el-aside width="200px" class="sidebar">
          <el-menu :default-active="activeMenu" router>
            <el-menu-item index="/user/home">
              <span>我的报修</span>
            </el-menu-item>
            <el-menu-item index="/user/repair">
              <span>提交报修</span>
            </el-menu-item>
          </el-menu>
        </el-aside>
        <el-main class="main-content">
          <h3>我的报修记录</h3>
          <el-table :data="repairs" style="width: 100%" v-loading="loading">
            <el-table-column prop="id" label="ID" width="80" />
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
          </el-table>
          <el-empty v-if="!loading && repairs.length === 0" description="暂无报修记录" />
        </el-main>
      </el-container>
    </el-container>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { repairAPI } from '@/api'

const router = useRouter()
const username = ref(localStorage.getItem('username') || '')
const repairs = ref([])
const loading = ref(false)
const activeMenu = ref('/user/home')

const loadRepairs = async () => {
  const userId = localStorage.getItem('userId')
  if (!userId) return

  loading.value = true
  try {
    const response = await repairAPI.getUserRepairs(userId)
    if (response.success) {
      repairs.value = response.data || []
    }
  } catch (error) {
    ElMessage.error('加载报修记录失败')
  } finally {
    loading.value = false
  }
}

const getStatusType = (status) => {
  const statusMap = {
    '未处理': 'warning',
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

const handleLogout = () => {
  localStorage.clear()
  ElMessage.success('已退出登录')
  router.push('/user/login')
}

onMounted(() => {
  loadRepairs()
})
</script>

<style scoped>
.home-container {
  min-height: 100vh;
}

.header {
  background: #409eff;
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

h3 {
  margin-bottom: 20px;
  color: #333;
}
</style>
