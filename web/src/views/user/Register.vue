<template>
  <div class="register-container">
    <div class="register-wrapper">
      <div class="register-header">
        <h1>物业报修</h1>
        <p>创建您的账号</p>
      </div>
      <el-card class="register-card" shadow="never">
        <template #header>
          <div class="card-header">
            <span>用户注册</span>
          </div>
        </template>
        <el-form :model="registerForm" :rules="rules" ref="registerFormRef" label-position="top">
          <el-form-item label="用户名" prop="username">
            <el-input 
              v-model="registerForm.username" 
              placeholder="请输入用户名（3-20个字符）"
              size="large"
              clearable
            />
          </el-form-item>
          <el-form-item label="密码" prop="password">
            <el-input 
              v-model="registerForm.password" 
              type="password" 
              placeholder="请输入密码（至少6个字符）" 
              show-password
              size="large"
            />
          </el-form-item>
          <el-form-item label="确认密码" prop="confirmPassword">
            <el-input 
              v-model="registerForm.confirmPassword" 
              type="password" 
              placeholder="请再次输入密码" 
              show-password
              size="large"
            />
          </el-form-item>
          <el-form-item>
            <el-button 
              type="primary" 
              @click="handleRegister" 
              :loading="loading" 
              size="large"
              class="register-button"
            >
              注册
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
        <div class="login-link">
          <el-link type="primary" @click="goToLogin">已有账号？去登录</el-link>
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
const registerFormRef = ref(null)
const loading = ref(false)

const registerForm = reactive({
  username: '',
  password: '',
  confirmPassword: ''
})

const validateConfirmPassword = (rule, value, callback) => {
  if (value !== registerForm.password) {
    callback(new Error('两次输入的密码不一致'))
  } else {
    callback()
  }
}

const rules = {
  username: [
    { required: true, message: '请输入用户名', trigger: 'blur' },
    { min: 3, max: 20, message: '用户名长度在 3 到 20 个字符', trigger: 'blur' }
  ],
  password: [
    { required: true, message: '请输入密码', trigger: 'blur' },
    { min: 6, message: '密码长度至少 6 个字符', trigger: 'blur' }
  ],
  confirmPassword: [
    { required: true, message: '请再次输入密码', trigger: 'blur' },
    { validator: validateConfirmPassword, trigger: 'blur' }
  ]
}

const handleRegister = async () => {
  if (!registerFormRef.value) return

  await registerFormRef.value.validate(async (valid) => {
    if (valid) {
      loading.value = true
      try {
        const response = await authAPI.register(registerForm.username, registerForm.password)
        if (response.success) {
          ElMessage.success('注册成功，请登录')
          router.push('/user/login')
        } else {
          ElMessage.error(response.error || '注册失败')
        }
      } catch (error) {
        ElMessage.error(error.response?.data?.error || '注册失败')
      } finally {
        loading.value = false
      }
    }
  })
}

const resetForm = () => {
  if (registerFormRef.value) {
    registerFormRef.value.resetFields()
  }
}

const goToLogin = () => {
  router.push('/user/login')
}
</script>

<style scoped>
.register-container {
  width: 100%;
  min-height: 100vh;
  min-height: -webkit-fill-available;
  display: flex;
  align-items: center;
  justify-content: center;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  padding: 20px;
}

.register-wrapper {
  width: 100%;
  max-width: 400px;
}

.register-header {
  text-align: center;
  margin-bottom: 30px;
  color: white;
}

.register-header h1 {
  font-size: 28px;
  margin-bottom: 8px;
}

.register-header p {
  font-size: 14px;
  opacity: 0.9;
}

.register-card {
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

.register-button,
.reset-button {
  width: 100%;
  height: 48px;
  font-size: 16px;
  border-radius: 8px;
  margin-bottom: 10px;
}

.login-link {
  text-align: center;
  margin-top: 16px;
  padding-bottom: 10px;
}

@media screen and (max-width: 480px) {
  .register-container {
    padding: 15px;
  }

  .register-header h1 {
    font-size: 24px;
  }

  .register-card :deep(.el-card__header) {
    padding: 16px 20px;
  }

  .register-card :deep(.el-form-item__label) {
    font-size: 14px;
  }

  .register-button,
  .reset-button {
    height: 44px;
    font-size: 15px;
  }
}

@media screen and (min-width: 768px) {
  .register-card {
    border-radius: 16px;
    box-shadow: 0 8px 24px rgba(0, 0, 0, 0.2);
  }
}
</style>
