<template>
  <div class="home-container">
    <header class="header">
      <h1>物业报修管理系统</h1>
      <div class="user-info">
        <span>欢迎, {{ user?.username }}</span>
        <el-button @click="handleLogout">退出登录</el-button>
      </div>
    </header>

    <div class="content">
      <el-tabs v-model="activeTab" type="card">
        <el-tab-pane label="新增报修" name="create">
          <div class="form-container">
            <el-form :model="workOrderForm" label-width="100px">
              <el-form-item label="报修类型">
                <el-select v-model="workOrderForm.type" placeholder="请选择报修类型">
                  <el-option label="水电维修" value="水电维修"></el-option>
                  <el-option label="家电维修" value="家电维修"></el-option>
                  <el-option label="房屋维修" value="房屋维修"></el-option>
                  <el-option label="其他" value="其他"></el-option>
                </el-select>
              </el-form-item>
              <el-form-item label="报修描述">
                <el-textarea v-model="workOrderForm.description" placeholder="请详细描述报修问题"></el-textarea>
              </el-form-item>
              <el-form-item label="报修图片">
                <el-input v-model="workOrderForm.images" placeholder="请输入图片URL（多个用逗号分隔）"></el-input>
              </el-form-item>
              <el-form-item>
                <el-button type="primary" @click="handleCreate">提交报修</el-button>
              </el-form-item>
            </el-form>
          </div>
        </el-tab-pane>

        <el-tab-pane label="报修记录" name="list">
          <el-table :data="workOrders" border>
            <el-table-column prop="id" label="报修编号"></el-table-column>
            <el-table-column prop="type" label="报修类型"></el-table-column>
            <el-table-column prop="description" label="报修描述"></el-table-column>
            <el-table-column prop="images" label="报修图片">
              <template #default="scope">
                <img v-if="scope.row.images" :src="scope.row.images.split(',')[0]" width="50" height="50" />
                <span v-else>无</span>
              </template>
            </el-table-column>
            <el-table-column prop="status" label="状态">
              <template #default="scope">
                <el-tag :type="getStatusType(scope.row.status)">{{ getStatusText(scope.row.status) }}</el-tag>
              </template>
            </el-table-column>
            <el-table-column prop="created_at" label="创建时间"></el-table-column>
          </el-table>
        </el-tab-pane>
      </el-tabs>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { getUser, removeUser } from '../utils/storage'
import { createWorkOrder, getWorkOrders } from '../api'

const router = useRouter()
const user = ref(getUser())
const activeTab = ref('create')
const workOrders = ref([])

const workOrderForm = reactive({
  type: '',
  description: '',
  images: ''
})

onMounted(() => {
  if (!user.value) {
    router.push('/')
    return
  }
  loadWorkOrders()
})

const loadWorkOrders = async () => {
  try {
    const res = await getWorkOrders(user.value.id)
    workOrders.value = res.data.data
  } catch (err) {
    ElMessage.error('获取报修记录失败')
  }
}

const handleCreate = async () => {
  if (!workOrderForm.type || !workOrderForm.description) {
    ElMessage.warning('请填写完整信息')
    return
  }

  try {
    await createWorkOrder({
      user_id: user.value.id,
      type: workOrderForm.type,
      description: workOrderForm.description,
      images: workOrderForm.images
    })
    ElMessage.success('报修提交成功')
    workOrderForm.type = ''
    workOrderForm.description = ''
    workOrderForm.images = ''
    loadWorkOrders()
    activeTab.value = 'list'
  } catch (err) {
    ElMessage.error('提交失败')
  }
}

const handleLogout = () => {
  removeUser()
  router.push('/')
}

const getStatusText = (status) => {
  const statusMap = {
    1: '未处理',
    2: '处理中',
    3: '已完成'
  }
  return statusMap[status] || '未知'
}

const getStatusType = (status) => {
  const typeMap = {
    1: 'warning',
    2: 'info',
    3: 'success'
  }
  return typeMap[status] || 'default'
}
</script>

<style scoped>
.home-container {
  min-height: 100vh;
  background: #f5f5f5;
}

.header {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  padding: 20px 30px;
  display: flex;
  justify-content: space-between;
  align-items: center;
  color: white;
}

.header h1 {
  margin: 0;
}

.user-info {
  display: flex;
  align-items: center;
  gap: 20px;
}

.content {
  padding: 20px;
}

.form-container {
  background: white;
  padding: 20px;
  border-radius: 10px;
  box-shadow: 0 0 10px rgba(0, 0, 0, 0.05);
}
</style>