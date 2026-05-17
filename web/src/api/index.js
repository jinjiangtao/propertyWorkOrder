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
    api.put('/repair/status', { repair_id: repairId, status })
}

export default api
