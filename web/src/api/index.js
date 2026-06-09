import axios from 'axios'

const API_BASE_URL = 'http://localhost:8080/api'

const api = axios.create({
  baseURL: API_BASE_URL,
  timeout: 10000,
  headers: {
    'Content-Type': 'application/json'
  }
})

api.interceptors.response.use(
  response => response.data,
  error => {
    console.error('API Error:', error)
    return Promise.reject(error)
  }
)

export const authAPI = {
  register: (username, password) =>
    api.post('/register', { username, password }),

  login: (username, password) =>
    api.post('/login', { username, password })
}

export const repairAPI = {
  createRepair: (data) =>
    api.post('/repair/create', data),

  getUserRepairs: (userId) =>
    api.get('/repair/user', { params: { user_id: userId } }),

  getAllRepairs: () =>
    api.get('/repair/all'),

  updateRepairStatus: (repairId, status) =>
    api.put('/repair/status', { repair_id: repairId, status }),

  assignWorker: (repairId, workerId) =>
    api.post('/repair/assign', { repair_id: repairId, worker_id: workerId }),

  getWorkerRepairs: (workerId) =>
    api.get('/repair/worker', { params: { worker_id: workerId } }),

  acceptOrder: (repairId, workerId) =>
    api.post('/repair/accept', { repair_id: repairId, worker_id: workerId }),

  rejectOrder: (repairId, workerId) =>
    api.post('/repair/reject', { repair_id: repairId, worker_id: workerId }),

  submitRepairResult: (repairId, workerId, result, imgs) =>
    api.post('/repair/result', { repair_id: repairId, worker_id: workerId, repair_result: result, repair_imgs: imgs }),

  getWorkerStats: () =>
    api.get('/repair/stats')
}

export const workerAPI = {
  login: (workNo, password) =>
    api.post('/worker/login', { work_no: workNo, password }),

  getWorkers: () =>
    api.get('/worker/list'),

  createWorker: (data) =>
    api.post('/worker/create', data),

  updateWorker: (data) =>
    api.put('/worker/update', data),

  toggleStatus: (workerId) =>
    api.put('/worker/status', null, { params: { worker_id: workerId } })
}

export default api
