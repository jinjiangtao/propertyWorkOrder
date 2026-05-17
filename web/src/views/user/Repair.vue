<template>
  <div class="repair-container">
    <div class="header">
      <div class="header-content">
        <el-button text @click="goBack" class="back-btn">← 返回</el-button>
        <h2>提交报修</h2>
        <div style="width: 50px;"></div>
      </div>
    </div>
    
    <div class="content">
      <el-form :model="repairForm" :rules="rules" ref="repairFormRef" label-position="top">
        <el-form-item label="报修类型" prop="repairType">
          <el-select 
            v-model="repairForm.repairType" 
            placeholder="请选择报修类型" 
            size="large"
            class="type-select"
          >
            <el-option label="🏠 水电维修" value="水电维修" />
            <el-option label="🚪 门窗维修" value="门窗维修" />
            <el-option label="🔌 家电维修" value="家电维修" />
            <el-option label="🚿 管道疏通" value="管道疏通" />
            <el-option label="📦 其他" value="其他" />
          </el-select>
        </el-form-item>
        
        <el-form-item label="报修描述" prop="description">
          <el-input
            v-model="repairForm.description"
            type="textarea"
            :rows="6"
            placeholder="请详细描述您的问题，例如：具体位置、故障情况等"
            size="large"
            class="description-input"
          />
          <div class="word-count">{{ repairForm.description.length }}/200</div>
        </el-form-item>
        
        <el-form-item label="图片（可选）">
          <el-input 
            v-model="repairForm.imageUrl" 
            placeholder="请输入图片URL（可选）"
            size="large"
          />
        </el-form-item>
        
        <el-form-item>
          <el-button 
            type="primary" 
            @click="handleSubmit" 
            :loading="loading" 
            size="large"
            class="submit-button"
          >
            提交报修
          </el-button>
          <el-button 
            @click="resetForm" 
            size="large"
            class="reset-button"
          >
            重置
          </el-button>
        </el-form-item>
      </el-form>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { repairAPI } from '@/api'

const router = useRouter()
const username = ref(localStorage.getItem('username') || '')
const activeMenu = ref('/user/repair')
const repairFormRef = ref(null)
const loading = ref(false)

const repairForm = reactive({
  repairType: '',
  description: '',
  imageUrl: ''
})

const rules = {
  repairType: [
    { required: true, message: '请选择报修类型', trigger: 'change' }
  ],
  description: [
    { required: true, message: '请输入报修描述', trigger: 'blur' },
    { min: 5, message: '描述至少5个字符', trigger: 'blur' },
    { max: 200, message: '描述不能超过200个字符', trigger: 'blur' }
  ]
}

const handleSubmit = async () => {
  if (!repairFormRef.value) return

  await repairFormRef.value.validate(async (valid) => {
    if (valid) {
      const userId = localStorage.getItem('userId')
      if (!userId) {
        ElMessage.error('用户未登录')
        return
      }

      loading.value = true
      try {
        const response = await repairAPI.createRepair({
          user_id: parseInt(userId),
          username: username.value,
          repair_type: repairForm.repairType,
          description: repairForm.description,
          image_url: repairForm.imageUrl
        })

        if (response.success) {
          ElMessage.success('报修提交成功')
          setTimeout(() => {
            router.push('/user/home')
          }, 1500)
        } else {
          ElMessage.error(response.error || '提交失败')
        }
      } catch (error) {
        ElMessage.error('提交失败，请重试')
      } finally {
        loading.value = false
      }
    }
  })
}

const resetForm = () => {
  if (repairFormRef.value) {
    repairFormRef.value.resetFields()
  }
  repairForm.imageUrl = ''
}

const goBack = () => {
  router.push('/user/home')
}

onMounted(() => {
  if (!localStorage.getItem('userId')) {
    router.push('/user/login')
  }
})
</script>

<style scoped>
.repair-container {
  width: 100%;
  min-height: 100vh;
  min-height: -webkit-fill-available;
  background: #f5f5f5;
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

.back-btn {
  color: white;
  font-size: 16px;
}

.header h2 {
  margin: 0;
  font-size: 20px;
}

.content {
  padding: 20px 16px;
}

.type-select {
  width: 100%;
}

.type-select :deep(.el-input__wrapper) {
  border-radius: 8px;
  padding: 12px 16px;
}

.description-input :deep(.el-textarea__inner) {
  border-radius: 8px;
  padding: 12px;
  font-size: 15px;
  line-height: 1.6;
}

.word-count {
  text-align: right;
  font-size: 12px;
  color: #999;
  margin-top: 4px;
}

.submit-button,
.reset-button {
  width: 100%;
  height: 48px;
  font-size: 16px;
  border-radius: 8px;
  margin-bottom: 12px;
}

.submit-button {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  border: none;
}

.submit-button:hover {
  background: linear-gradient(135deg, #5a6fd6 0%, #6a4190 100%);
}

:deep(.el-form-item__label) {
  font-size: 15px;
  font-weight: 500;
  color: #333;
  padding-bottom: 8px !important;
}

:deep(.el-input__wrapper) {
  border-radius: 8px;
  padding: 12px 16px;
  box-shadow: 0 0 0 1px #dcdfe6;
}

:deep(.el-input__wrapper:hover) {
  box-shadow: 0 0 0 1px #c0c4cc;
}

:deep(.el-input__wrapper.is-focus) {
  box-shadow: 0 0 0 1px #667eea;
}

@media screen and (min-width: 768px) {
  .content {
    max-width: 600px;
    margin: 0 auto;
    padding: 30px 20px;
  }
  
  .submit-button,
  .reset-button {
    height: 50px;
    font-size: 17px;
  }
}
</style>
