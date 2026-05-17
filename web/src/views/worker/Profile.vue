<template>
  <div class="profile-container">
    <el-container>
      <el-header class="header">
        <div class="header-content">
          <el-button @click="goBack">← 返回</el-button>
          <h2>个人中心</h2>
        </div>
      </el-header>
      <el-main class="main-content">
        <el-card class="profile-card">
          <div class="profile-header">
            <div class="avatar">
              <User />
            </div>
            <div class="profile-info">
              <h3>{{ workerName }}</h3>
              <p>{{ skillType }}</p>
            </div>
          </div>
          
          <div class="profile-section">
            <h4>基本信息</h4>
            <el-row :gutter="20">
              <el-col :span="12">
                <div class="info-item">
                  <span class="label">工号：</span>
                  <span class="value">{{ workNo }}</span>
                </div>
              </el-col>
              <el-col :span="12">
                <div class="info-item">
                  <span class="label">姓名：</span>
                  <span class="value">{{ workerName }}</span>
                </div>
              </el-col>
              <el-col :span="12">
                <div class="info-item">
                  <span class="label">工种：</span>
                  <span class="value">{{ skillType }}</span>
                </div>
              </el-col>
              <el-col :span="12">
                <div class="info-item">
                  <span class="label">状态：</span>
                  <el-tag type="success">在职</el-tag>
                </div>
              </el-col>
            </el-row>
          </div>
          
          <div class="profile-section">
            <h4>统计信息</h4>
            <div class="stats-grid">
              <div class="stat-item">
                <div class="stat-value">{{ stats.total }}</div>
                <div class="stat-label">总工单</div>
              </div>
              <div class="stat-item">
                <div class="stat-value processing">{{ stats.processing }}</div>
                <div class="stat-label">处理中</div>
              </div>
              <div class="stat-item">
                <div class="stat-value completed">{{ stats.completed }}</div>
                <div class="stat-label">已完成</div>
              </div>
            </div>
          </div>
          
          <div class="profile-section">
            <el-button type="danger" @click="handleLogout">退出登录</el-button>
          </div>
        </el-card>
      </el-main>
    </el-container>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { User } from '@element-plus/icons-vue'
import { repairAPI } from '@/api'

const router = useRouter()
const workerId = ref(localStorage.getItem('workerId') || '')
const workerName = ref(localStorage.getItem('workerName') || '')
const skillType = ref(localStorage.getItem('workerSkillType') || '')
const workNo = ref(localStorage.getItem('workerWorkNo') || '')

const stats = ref({
  total: 0,
  processing: 0,
  completed: 0
})

const loadStats = async () => {
  try {
    const response = await repairAPI.getWorkerRepairs(workerId.value)
    if (response.success) {
      const orders = response.data || []
      stats.value = {
        total: orders.length,
        processing: orders.filter(o => o.status === '处理中').length,
        completed: orders.filter(o => o.status === '已完成').length
      }
    }
  } catch (error) {
    console.error('加载统计失败')
  }
}

const goBack = () => {
  router.push('/worker/home')
}

const handleLogout = () => {
  localStorage.clear()
  ElMessage.success('已退出登录')
  router.push('/worker/login')
}

onMounted(() => {
  loadStats()
})
</script>

<style scoped>
.profile-container {
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
  align-items: center;
  gap: 20px;
}

.header h2 {
  margin: 0;
  font-size: 20px;
}

.main-content {
  padding: 20px;
}

.profile-card {
  max-width: 600px;
  margin: 0 auto;
}

.profile-header {
  display: flex;
  align-items: center;
  gap: 20px;
  padding-bottom: 20px;
  border-bottom: 1px solid #eee;
  margin-bottom: 20px;
}

.avatar {
  width: 80px;
  height: 80px;
  border-radius: 50%;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
  font-size: 36px;
}

.profile-info h3 {
  margin: 0 0 8px 0;
  font-size: 20px;
}

.profile-info p {
  margin: 0;
  color: #999;
}

.profile-section {
  margin-bottom: 24px;
}

.profile-section h4 {
  margin: 0 0 16px 0;
  color: #333;
}

.info-item {
  margin-bottom: 12px;
}

.label {
  color: #999;
  margin-right: 8px;
}

.value {
  color: #333;
}

.stats-grid {
  display: flex;
  gap: 20px;
}

.stat-item {
  flex: 1;
  text-align: center;
  padding: 20px;
  background: #f8f9fa;
  border-radius: 12px;
}

.stat-value {
  font-size: 32px;
  font-weight: bold;
  color: #333;
}

.stat-value.processing {
  color: #409eff;
}

.stat-value.completed {
  color: #67c23a;
}

.stat-label {
  font-size: 14px;
  color: #999;
  margin-top: 8px;
}
</style>