<template>
  <div class="detail-container">
    <el-container>
      <el-header class="header">
        <div class="header-content">
          <el-button @click="goBack">← 返回</el-button>
          <h2>工单详情</h2>
        </div>
      </el-header>
      <el-main class="main-content">
        <el-card v-if="order" class="detail-card">
          <div class="detail-section">
            <h3>基本信息</h3>
            <el-row :gutter="20">
              <el-col :span="12">
                <div class="info-item">
                  <span class="label">工单ID：</span>
                  <span class="value">{{ order.id }}</span>
                </div>
              </el-col>
              <el-col :span="12">
                <div class="info-item">
                  <span class="label">报修用户：</span>
                  <span class="value">{{ order.username }}</span>
                </div>
              </el-col>
              <el-col :span="12">
                <div class="info-item">
                  <span class="label">报修类型：</span>
                  <span class="value">{{ order.repair_type }}</span>
                </div>
              </el-col>
              <el-col :span="12">
                <div class="info-item">
                  <span class="label">状态：</span>
                  <el-tag :type="getStatusType(order.status)">{{ order.status }}</el-tag>
                </div>
              </el-col>
              <el-col :span="12">
                <div class="info-item">
                  <span class="label">提交时间：</span>
                  <span class="value">{{ formatDate(order.created_at) }}</span>
                </div>
              </el-col>
              <el-col :span="12">
                <div class="info-item">
                  <span class="label">更新时间：</span>
                  <span class="value">{{ formatDate(order.updated_at) }}</span>
                </div>
              </el-col>
            </el-row>
          </div>
          
          <div class="detail-section">
            <h3>报修描述</h3>
            <p class="description">{{ order.description }}</p>
          </div>
          
          <div v-if="order.image_url" class="detail-section">
            <h3>报修图片</h3>
            <img :src="order.image_url" class="repair-image" />
          </div>
          
          <div v-if="order.worker_name" class="detail-section">
            <h3>派单信息</h3>
            <div class="info-item">
              <span class="label">派工：</span>
              <span class="value">{{ order.worker_name }}</span>
            </div>
          </div>
          
          <div v-if="order.repair_result" class="detail-section">
            <h3>维修结果</h3>
            <p class="description">{{ order.repair_result }}</p>
            <div v-if="order.repair_imgs" class="repair-imgs">
              <img v-for="(img, index) in repairImgsList" :key="index" :src="img" class="result-image" />
            </div>
          </div>
          
          <div v-if="order.status === '处理中'" class="submit-section">
            <h3>提交维修结果</h3>
            <el-form :model="resultForm" ref="resultFormRef">
              <el-form-item label="维修结果">
                <el-textarea v-model="resultForm.result" :rows="4" placeholder="请描述维修过程和结果" />
              </el-form-item>
              <el-form-item label="维修后图片（可选）">
                <el-upload
                  action="#"
                  :auto-upload="false"
                  :show-file-list="true"
                  :file-list="uploadFiles"
                  accept="image/*"
                >
                  <el-button size="small" type="primary">点击上传</el-button>
                </el-upload>
              </el-form-item>
              <el-form-item>
                <el-button type="success" @click="submitResult">提交维修结果</el-button>
              </el-form-item>
            </el-form>
          </div>
        </el-card>
      </el-main>
    </el-container>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { ElMessage } from 'element-plus'
import { repairAPI } from '@/api'

const router = useRouter()
const route = useRoute()
const order = ref(null)
const loading = ref(false)
const resultForm = ref({
  result: ''
})
const uploadFiles = ref([])
const resultFormRef = ref(null)

const repairImgsList = computed(() => {
  if (!order.value?.repair_imgs) return []
  return order.value.repair_imgs.split(',')
})

const loadOrder = async () => {
  const orderId = route.params.id
  if (!orderId) return
  
  loading.value = true
  try {
    const response = await repairAPI.getWorkerRepairs(parseInt(localStorage.getItem('workerId')))
    if (response.success) {
      const data = response.data || []
      order.value = data.find(o => o.id === parseInt(orderId))
    }
  } catch (error) {
    ElMessage.error('加载工单详情失败')
  } finally {
    loading.value = false
  }
}

const submitResult = async () => {
  if (!resultForm.value.result.trim()) {
    ElMessage.error('请填写维修结果')
    return
  }
  
  try {
    const response = await repairAPI.submitRepairResult(
      parseInt(route.params.id),
      parseInt(localStorage.getItem('workerId')),
      resultForm.value.result,
      ''
    )
    if (response.success) {
      ElMessage.success('提交成功')
      router.push('/worker/home')
    } else {
      ElMessage.error(response.error || '提交失败')
    }
  } catch (error) {
    ElMessage.error(error.response?.data?.error || '提交失败')
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

const goBack = () => {
  router.push('/worker/home')
}

onMounted(() => {
  loadOrder()
})
</script>

<style scoped>
.detail-container {
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

.detail-card {
  max-width: 800px;
  margin: 0 auto;
}

.detail-section {
  margin-bottom: 24px;
}

.detail-section h3 {
  margin: 0 0 16px 0;
  padding-bottom: 8px;
  border-bottom: 1px solid #eee;
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

.description {
  color: #666;
  line-height: 1.8;
}

.repair-image, .result-image {
  max-width: 300px;
  max-height: 200px;
  object-fit: cover;
  border-radius: 8px;
  margin-right: 12px;
  margin-bottom: 12px;
}

.repair-imgs {
  margin-top: 12px;
}

.submit-section {
  margin-top: 24px;
  padding-top: 24px;
  border-top: 1px solid #eee;
}
</style>