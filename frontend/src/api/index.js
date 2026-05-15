import axios from 'axios'

const instance = axios.create({
  baseURL: '/api',
  timeout: 10000
})

export const register = (username, password) => {
  return instance.post('/register', { username, password })
}

export const login = (username, password) => {
  return instance.post('/login', { username, password })
}

export const adminLogin = (username, password) => {
  return instance.post('/admin/login', { username, password })
}

export const createWorkOrder = (data) => {
  return instance.post('/workorder', data)
}

export const getWorkOrders = (userId) => {
  const params = userId ? { user_id: userId } : {}
  return instance.get('/workorders', { params })
}

export const updateWorkOrderStatus = (id, status) => {
  return instance.put(`/workorder/${id}`, { status })
}