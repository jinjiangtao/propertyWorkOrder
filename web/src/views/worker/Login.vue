<template>
  <div class="login-container">
    <el-card class="login-card">
      <template #header>
        <div class="card-header">
          <span>工人登录</span>
        </div>
      </template>
      <el-form :model="loginForm" :rules="rules" ref="loginFormRef" label-width="80px">
        <el-form-item label="工号" prop="workNo">
          <el-input v-model="loginForm.workNo" placeholder="请输入工号" />
        </el-form-item>
        <el-form-item label="密码" prop="password">
          <el-input v-model="loginForm.password" type="password" placeholder="请输入密码" show-password />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleLogin" :loading="loading" style="width: 100%">登录</el-button>
        </el-form-item>
      </el-form>
      <div class="back-link">
        <el-link type="primary" @click="goToUserLogin">返回用户登录</el-link>
      </div>
    </el-card>
  </div>
</template>

<script setup>
import { ref, reactive } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { workerAPI } from '@/api'

const router = useRouter()
const loginFormRef = ref(null)
const loading = ref(false)

const loginForm = reactive({
  workNo: '',
  password: ''
})

const rules = {
  workNo: [
    { required: true, message: '请输入工号', trigger: 'blur' }
  ],
  password: [
    { required: true, message: '请输入密码', trigger: 'blur' }
  ]
}

const handleLogin = async () => {
  if (!loginFormRef.value) return

  await loginFormRef.value.validate(async (valid) => {
    if (valid) {
      loading.value = true
      try {
        const response = await workerAPI.login(loginForm.workNo, loginForm.password)
        if (response.success) {
          localStorage.setItem('workerId', response.worker_id)
          localStorage.setItem('workerName', response.name)
          localStorage.setItem('workerSkillType', response.skill_type)
          localStorage.setItem('workerWorkNo', response.work_no)
          localStorage.setItem('isWorker', 'true')
          ElMessage.success('登录成功')
          router.push('/worker/home')
        } else {
          ElMessage.error(response.error || '登录失败')
        }
      } catch (error) {
        ElMessage.error(error.response?.data?.error || '登录失败')
      } finally {
        loading.value = false
      }
    }
  })
}

const goToUserLogin = () => {
  router.push('/user/login')
}
</script>

<style scoped>
.login-container {
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: 100vh;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
}

.login-card {
  width: 450px;
  box-shadow: 0 8px 16px rgba(0, 0, 0, 0.2);
}

.card-header {
  text-align: center;
  font-size: 24px;
  font-weight: bold;
  color: #333;
}

.back-link {
  text-align: center;
  margin-top: 20px;
}
</style>