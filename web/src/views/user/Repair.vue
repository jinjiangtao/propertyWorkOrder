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
          <h3>提交报修</h3>
          <el-card class="repair-form-card">
            <el-form :model="repairForm" :rules="rules" ref="repairFormRef" label-width="100px">
              <el-form-item label="报修类型" prop="repairType">
                <el-select v-model="repairForm.repairType" placeholder="请选择报修类型" style="width: 100%">
                  <el-option label="水电维修" value="水电维修" />
                  <el-option label="门窗维修" value="门窗维修" />
                  <el-option label="家电维修" value="家电维修" />
                  <el-option label="管道疏通" value="管道疏通" />
                  <el-option label="其他" value="其他" />
                </el-select>
              </el-form-item>
              <el-form-item label="报修描述" prop="description">
                <el-input
                  v-model="repairForm.description"
                  type="textarea"
                  :rows="4"
                  placeholder="请详细描述您的问题"
                />
              </el-form-item>
              <el-form-item label="图片">
                <el-input v-model="repairForm.imageUrl" placeholder="请输入图片URL（可选）" />
              </el-form-item>
              <el-form-item>
                <el-button type="primary" @click="handleSubmit" :loading="loading" style="width: 100%">提交报修</el-button>
                <el-button @click="resetForm" style="width: 100%; margin-top: 10px">重置</el-button>
              </el-form-item>
            </el-form>
          </el-card>
        </el-main>
      </el-container>
    </el-container>
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
    { min: 5, message: '描述至少5个字符', trigger: 'blur' }
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
          resetForm()
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

const handleLogout = () => {
  localStorage.clear()
  ElMessage.success('已退出登录')
  router.push('/user/login')
}

onMounted(() => {
  if (!localStorage.getItem('userId')) {
    router.push('/user/login')
  }
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

.repair-form-card {
  max-width: 600px;
}
</style>
