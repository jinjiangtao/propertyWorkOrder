<template>
  <div class="home-container">
    <div class="header">
      <div class="header-content">
        <div class="header-left">
          <h2>报修记录</h2>
        </div>
        <div class="header-right">
          <span class="username">{{ username }}</span>
          <el-button type="danger" size="small" @click="handleLogout">退出</el-button>
        </div>
      </div>
    </div>
    
    <div class="content-wrapper">
      <div class="tab-bar">
        <div class="tab-item active">我的报修</div>
        <div class="tab-item" @click="goToRepair">提交报修</div>
      </div>
      
      <div class="repair-list">
        <div v-if="loading" class="loading">
          <el-icon class="is-loading"><Loading /></el-icon>
          <p>加载中...</p>
        </div>
        
        <div v-else-if="repairs.length === 0" class="empty-state">
          <div class="empty-icon">📋</div>
          <p>暂无报修记录</p>
          <el-button type="primary" @click="goToRepair">立即报修</el-button>
        </div>
        
        <div v-else class="repair-cards">
          <div 
            v-for="repair in repairs" 
            :key="repair.id" 
            class="repair-card"
          >
            <div class="card-header">
              <span class="repair-type">{{ repair.repair_type }}</span>
              <el-tag :type="getStatusType(repair.status)" size="small">
                {{ repair.status }}
              </el-tag>
            </div>
            <div class="card-body">
              <p class="description">{{ repair.description }}</p>
              <div class="meta">
                <span class="time">{{ formatDate(repair.created_at) }}</span>
              </div>
            </div>
          </div>
        </div>
      </div>
      
      <div class="bottom-nav">
        <div class="nav-item active" @click="refresh">
          <div class="nav-icon">📋</div>
          <span>报修记录</span>
        </div>
        <div class="nav-item" @click="goToRepair">
          <div class="nav-icon">➕</div>
          <span>提交报修</span>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { repairAPI } from '@/api'
import { Loading } from '@element-plus/icons-vue'

const router = useRouter()
const username = ref(localStorage.getItem('username') || '')
const repairs = ref([])
const loading = ref(false)

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

const refresh = () => {
  loadRepairs()
}

const goToRepair = () => {
  router.push('/user/repair')
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
  width: 100%;
  min-height: 100vh;
  min-height: -webkit-fill-available;
  background: #f5f5f5;
  display: flex;
  flex-direction: column;
}

.header {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  padding: 20px 16px;
  position: sticky;
  top: 0;
  z-index: 100;
}

.header-content {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.header h2 {
  margin: 0;
  font-size: 20px;
}

.header-right {
  display: flex;
  align-items: center;
  gap: 12px;
}

.username {
  font-size: 14px;
  opacity: 0.9;
}

.content-wrapper {
  flex: 1;
  padding-bottom: 70px;
}

.tab-bar {
  display: flex;
  background: white;
  border-bottom: 1px solid #eee;
}

.tab-item {
  flex: 1;
  text-align: center;
  padding: 14px 0;
  font-size: 15px;
  color: #666;
  position: relative;
}

.tab-item.active {
  color: #667eea;
  font-weight: bold;
}

.tab-item.active::after {
  content: '';
  position: absolute;
  bottom: 0;
  left: 50%;
  transform: translateX(-50%);
  width: 40px;
  height: 3px;
  background: #667eea;
  border-radius: 2px;
}

.repair-list {
  padding: 16px;
}

.loading {
  text-align: center;
  padding: 60px 20px;
  color: #999;
}

.loading .el-icon {
  font-size: 32px;
  margin-bottom: 12px;
}

.empty-state {
  text-align: center;
  padding: 80px 20px;
}

.empty-icon {
  font-size: 64px;
  margin-bottom: 16px;
}

.empty-state p {
  color: #999;
  margin-bottom: 24px;
  font-size: 15px;
}

.repair-cards {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.repair-card {
  background: white;
  border-radius: 12px;
  padding: 16px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.06);
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 12px;
}

.repair-type {
  font-size: 16px;
  font-weight: bold;
  color: #333;
}

.card-body .description {
  font-size: 14px;
  color: #666;
  line-height: 1.6;
  margin-bottom: 12px;
}

.meta {
  display: flex;
  justify-content: space-between;
  font-size: 12px;
  color: #999;
}

.bottom-nav {
  position: fixed;
  bottom: 0;
  left: 0;
  right: 0;
  background: white;
  display: flex;
  box-shadow: 0 -2px 10px rgba(0, 0, 0, 0.08);
  padding: 8px 0;
  padding-bottom: calc(8px + env(safe-area-inset-bottom));
}

.nav-item {
  flex: 1;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 4px;
  padding: 8px 0;
  color: #999;
  cursor: pointer;
}

.nav-item.active {
  color: #667eea;
}

.nav-icon {
  font-size: 24px;
}

.nav-item span {
  font-size: 12px;
}

@media screen and (min-width: 768px) {
  .repair-cards {
    max-width: 600px;
    margin: 0 auto;
  }
  
  .repair-card {
    border-radius: 16px;
    padding: 20px;
  }
}
</style>
