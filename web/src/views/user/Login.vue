<template>
  <div class="login-container">
    <div class="login-wrapper">
      <div class="login-header">
        <h1>物业报修</h1>
        <p>便捷报修服务</p>
      </div>
      <el-card class="login-card" shadow="never">
        <template #header>
          <div class="card-header">
            <span>用户登录</span>
          </div>
        </template>
        <el-form :model="loginForm" :rules="rules" ref="loginFormRef" label-position="top">
          <el-form-item label="用户名" prop="username">
            <el-input 
              v-model="loginForm.username" 
              placeholder="请输入用户名"
              size="large"
              clearable
            />
          </el-form-item>
          <el-form-item label="密码" prop="password">
            <el-input 
              v-model="loginForm.password" 
              type="password" 
              placeholder="请输入密码" 
              show-password
              size="large"
            />
          </el-form-item>
          <el-form-item>
            <el-button 
              type="primary" 
              @click="handleLogin" 
              :loading="loading" 
              size="large"
              class="login-button"
            >
              登录
            </el-button>
          </el-form-item>
        </el-form>
        <div class="register-link">
          <el-link type="primary" @click="goToRegister">还没有账号？去注册</el-link>
        </div>
      </el-card>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { authAPI } from '@/api'

const router = useRouter()
const loginFormRef = ref(null)
const loading = ref(false)

const loginForm = reactive({
  username: '',
  password: ''
})

const rules = {
  username: [
    { required: true, message: '请输入用户名', trigger: 'blur' }
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
        const response = await authAPI.login(loginForm.username, loginForm.password)
        if (response.success) {
          localStorage.setItem('userId', response.user_id)
          localStorage.setItem('username', loginForm.username)
          localStorage.setItem('isAdmin', response.is_admin)
          ElMessage.success('登录成功')
          router.push('/user/home')
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

const goToRegister = () => {
  router.push('/user/register')
}
</script>

<style scoped>
.login-container {
  width: 100%;
  min-height: 100vh;
  min-height: -webkit-fill-available;
  display: flex;
  align-items: center;
  justify-content: center;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  padding: 20px;
}

.login-wrapper {
  width: 100%;
  max-width: 400px;
}

.login-header {
  text-align: center;
  margin-bottom: 30px;
  color: white;
}

.login-header h1 {
  font-size: 28px;
  margin-bottom: 8px;
}

.login-header p {
  font-size: 14px;
  opacity: 0.9;
}

.login-card {
  border-radius: 12px;
  overflow: hidden;
}

.card-header {
  text-align: center;
  font-size: 20px;
  font-weight: bold;
  color: #333;
  padding: 10px 0;
}

.login-button {
  width: 100%;
  height: 48px;
  font-size: 16px;
  border-radius: 8px;
}

.register-link {
  text-align: center;
  margin-top: 16px;
  padding-bottom: 10px;
}

@media screen and (max-width: 480px) {
  .login-container {
    padding: 15px;
  }

  .login-header h1 {
    font-size: 24px;
  }

  .login-card :deep(.el-card__header) {
    padding: 16px 20px;
  }

  .login-card :deep(.el-form-item__label) {
    font-size: 14px;
  }

  .login-button {
    height: 44px;
    font-size: 15px;
  }
}

@media screen and (min-width: 768px) {
  .login-card {
    border-radius: 16px;
    box-shadow: 0 8px 24px rgba(0, 0, 0, 0.2);
  }
}
</style>
