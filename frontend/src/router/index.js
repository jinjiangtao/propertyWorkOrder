import { createRouter, createWebHistory } from 'vue-router'

const routes = [
  {
    path: '/',
    name: 'UserLogin',
    component: () => import('../views/UserLogin.vue')
  },
  {
    path: '/user/register',
    name: 'UserRegister',
    component: () => import('../views/UserRegister.vue')
  },
  {
    path: '/user/home',
    name: 'UserHome',
    component: () => import('../views/UserHome.vue')
  },
  {
    path: '/admin/login',
    name: 'AdminLogin',
    component: () => import('../views/AdminLogin.vue')
  },
  {
    path: '/admin/home',
    name: 'AdminHome',
    component: () => import('../views/AdminHome.vue')
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

export default router