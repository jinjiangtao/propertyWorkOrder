<template>
  <div class="workers-container">
    <el-container>
      <el-header class="header">
        <div class="header-content">
          <h2>物业报修管理系统 - 工人管理</h2>
          <div class="user-info">
            <span>管理员：{{ username }}</span>
            <el-button type="danger" size="small" @click="handleLogout">退出登录</el-button>
          </div>
        </div>
      </el-header>
      <el-container>
        <el-aside width="200px" class="sidebar">
          <el-menu :default-active="activeMenu" @select="handleMenuSelect">
            <el-menu-item index="home">
              <span>工单管理</span>
            </el-menu-item>
            <el-menu-item index="workers">
              <span>工人管理</span>
            </el-menu-item>
            <el-menu-item index="stats">
              <span>统计分析</span>
            </el-menu-item>
          </el-menu>
        </el-aside>
        <el-main class="main-content">
          <div class="toolbar">
            <h3>工人列表</h3>
            <el-button type="primary" @click="showAddDialog = true">新增工人</el-button>
          </div>
          <el-table :data="workers" style="width: 100%" v-loading="loading">
            <el-table-column prop="id" label="ID" width="80" />
            <el-table-column prop="work_no" label="工号" width="120" />
            <el-table-column prop="name" label="姓名" width="100" />
            <el-table-column prop="phone" label="手机号" width="130" />
            <el-table-column prop="skill_type" label="擅长工种" width="120" />
            <el-table-column prop="status" label="状态" width="100">
              <template #default="scope">
                <el-tag :type="scope.row.status === 1 ? 'success' : 'danger'">
                  {{ scope.row.status === 1 ? '在职' : '离职' }}
                </el-tag>
              </template>
            </el-table-column>
            <el-table-column prop="created_at" label="创建时间" width="180">
              <template #default="scope">
                {{ formatDate(scope.row.created_at) }}
              </template>
            </el-table-column>
            <el-table-column label="操作" width="200" fixed="right">
              <template #default="scope">
                <el-button size="small" @click="handleEdit(scope.row)">编辑</el-button>
                <el-button 
                  size="small" 
                  :type="scope.row.status === 1 ? 'danger' : 'success'" 
                  @click="handleToggleStatus(scope.row)"
                >
                  {{ scope.row.status === 1 ? '禁用' : '启用' }}
                </el-button>
              </template>
            </el-table-column>
          </el-table>
          <el-empty v-if="!loading && workers.length === 0" description="暂无工人" />
        </el-main>
      </el-container>
    </el-container>

    <el-dialog title="新增工人" v-model="showAddDialog" width="450px">
      <el-form :model="workerForm" :rules="rules" ref="workerFormRef">
        <el-form-item label="工号" prop="workNo">
          <el-input v-model="workerForm.workNo" placeholder="请输入工号" />
        </el-form-item>
        <el-form-item label="姓名" prop="name">
          <el-input v-model="workerForm.name" placeholder="请输入姓名" />
        </el-form-item>
        <el-form-item label="手机号" prop="phone">
          <el-input v-model="workerForm.phone" placeholder="请输入手机号（11位，以1开头）" />
        </el-form-item>
        <el-form-item label="密码" prop="password">
          <el-input v-model="workerForm.password" type="password" placeholder="请输入密码" />
        </el-form-item>
        <el-form-item label="擅长工种" prop="skillType">
          <el-select v-model="workerForm.skillType" placeholder="请选择工种">
            <el-option label="水电" value="水电" />
            <el-option label="木工" value="木工" />
            <el-option label="保洁" value="保洁" />
            <el-option label="综合维修" value="综合维修" />
          </el-select>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="showAddDialog = false">取消</el-button>
        <el-button type="primary" @click="handleAdd">确认添加</el-button>
      </template>
    </el-dialog>

    <el-dialog title="编辑工人" v-model="showEditDialog" width="450px">
      <el-form :model="editForm" :rules="editRules" ref="editFormRef">
        <el-form-item label="姓名" prop="name">
          <el-input v-model="editForm.name" placeholder="请输入姓名" />
        </el-form-item>
        <el-form-item label="手机号" prop="phone">
          <el-input v-model="editForm.phone" placeholder="请输入手机号（11位，以1开头）" />
        </el-form-item>
        <el-form-item label="擅长工种" prop="skillType">
          <el-select v-model="editForm.skillType" placeholder="请选择工种">
            <el-option label="水电" value="水电" />
            <el-option label="木工" value="木工" />
            <el-option label="保洁" value="保洁" />
            <el-option label="综合维修" value="综合维修" />
          </el-select>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="showEditDialog = false">取消</el-button>
        <el-button type="primary" @click="handleUpdate">确认更新</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { workerAPI } from '@/api'

const router = useRouter()
const username = ref(localStorage.getItem('username') || '')
const workers = ref([])
const loading = ref(false)
const activeMenu = ref('workers')

const showAddDialog = ref(false)
const showEditDialog = ref(false)
const workerFormRef = ref(null)
const editFormRef = ref(null)

const workerForm = reactive({
  workNo: '',
  name: '',
  phone: '',
  password: '',
  skillType: ''
})

const editForm = reactive({
  id: 0,
  name: '',
  phone: '',
  skillType: ''
})

const validatePhone = (rule, value, callback) => {
  const phoneRegex = /^1[3-9]\d{9}$/
  if (!value) {
    callback(new Error('请输入手机号'))
  } else if (!phoneRegex.test(value)) {
    callback(new Error('请输入正确的手机号（11位，以1开头）'))
  } else {
    callback()
  }
}

const rules = {
  workNo: [{ required: true, message: '请输入工号', trigger: 'blur' }],
  name: [{ required: true, message: '请输入姓名', trigger: 'blur' }],
  phone: [{ required: true, validator: validatePhone, trigger: 'blur' }],
  password: [{ required: true, message: '请输入密码', trigger: 'blur' }],
  skillType: [{ required: true, message: '请选择工种', trigger: 'change' }]
}

const editRules = {
  name: [{ required: true, message: '请输入姓名', trigger: 'blur' }],
  phone: [{ required: true, validator: validatePhone, trigger: 'blur' }],
  skillType: [{ required: true, message: '请选择工种', trigger: 'change' }]
}

const loadWorkers = async () => {
  loading.value = true
  try {
    const response = await workerAPI.getWorkers()
    if (response.success) {
      workers.value = response.data || []
    }
  } catch (error) {
    ElMessage.error('加载工人列表失败')
  } finally {
    loading.value = false
  }
}

const handleAdd = async () => {
  if (!workerFormRef.value) return
  
  await workerFormRef.value.validate(async (valid) => {
    if (valid) {
      try {
        const response = await workerAPI.createWorker({
          work_no: workerForm.workNo,
          name: workerForm.name,
          phone: workerForm.phone,
          password: workerForm.password,
          skill_type: workerForm.skillType
        })
        if (response.success) {
          ElMessage.success('添加成功')
          showAddDialog.value = false
          workerForm.workNo = ''
          workerForm.name = ''
          workerForm.phone = ''
          workerForm.password = ''
          workerForm.skillType = ''
          loadWorkers()
        } else {
          ElMessage.error(response.error || '添加失败')
        }
      } catch (error) {
        ElMessage.error(error.response?.data?.error || '添加失败')
      }
    }
  })
}

const handleEdit = (worker) => {
  editForm.id = worker.id
  editForm.name = worker.name
  editForm.phone = worker.phone
  editForm.skillType = worker.skill_type
  showEditDialog.value = true
}

const handleUpdate = async () => {
  if (!editFormRef.value) return
  
  await editFormRef.value.validate(async (valid) => {
    if (valid) {
      try {
        const response = await workerAPI.updateWorker({
          id: editForm.id,
          name: editForm.name,
          phone: editForm.phone,
          skill_type: editForm.skillType
        })
        if (response.success) {
          ElMessage.success('更新成功')
          showEditDialog.value = false
          loadWorkers()
        } else {
          ElMessage.error(response.error || '更新失败')
        }
      } catch (error) {
        ElMessage.error(error.response?.data?.error || '更新失败')
      }
    }
  })
}

const handleToggleStatus = async (worker) => {
  try {
    const response = await workerAPI.toggleStatus(worker.id)
    if (response.success) {
      ElMessage.success(response.message)
      loadWorkers()
    } else {
      ElMessage.error(response.error || '操作失败')
    }
  } catch (error) {
    ElMessage.error(error.response?.data?.error || '操作失败')
  }
}

const formatDate = (dateStr) => {
  if (!dateStr) return ''
  const date = new Date(dateStr)
  return date.toLocaleString('zh-CN')
}

const handleMenuSelect = (index) => {
  activeMenu.value = index
  if (index === 'home') {
    router.push('/admin/home')
  } else if (index === 'workers') {
    router.push('/admin/workers')
  } else if (index === 'stats') {
    router.push('/admin/stats')
  }
}

const handleLogout = () => {
  localStorage.clear()
  ElMessage.success('已退出登录')
  router.push('/admin/login')
}

onMounted(() => {
  if (!localStorage.getItem('userId')) {
    router.push('/admin/login')
    return
  }
  loadWorkers()
})
</script>

<style scoped>
.workers-container {
  min-height: 100vh;
}

.header {
  background: #f5576c;
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

.toolbar {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}

h3 {
  margin: 0;
  color: #333;
}
</style>