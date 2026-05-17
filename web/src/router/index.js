import { createRouter, createWebHistory } from 'vue-router'

const routes = [
  {
    path: '/',
    redirect: '/user/login'
  },
  {
    path: '/user',
    redirect: '/user/login'
  },
  {
    path: '/user/login',
    name: 'UserLogin',
    component: () => import('../views/user/Login.vue')
  },
  {
    path: '/user/register',
    name: 'UserRegister',
    component: () => import('../views/user/Register.vue')
  },
  {
    path: '/user/home',
    name: 'UserHome',
    component: () => import('../views/user/Home.vue'),
    meta: { requiresAuth: true }
  },
  {
    path: '/user/repair',
    name: 'UserRepair',
    component: () => import('../views/user/Repair.vue'),
    meta: { requiresAuth: true }
  },
  {
    path: '/admin',
    name: 'Admin',
    redirect: '/admin/login'
  },
  {
    path: '/admin/login',
    name: 'AdminLogin',
    component: () => import('../views/admin/Login.vue')
  },
  {
    path: '/admin/home',
    name: 'AdminHome',
    component: () => import('../views/admin/Home.vue'),
    meta: { requiresAuth: true }
  },
  {
    path: '/admin/workers',
    name: 'AdminWorkers',
    component: () => import('../views/admin/Workers.vue'),
    meta: { requiresAuth: true }
  },
  {
    path: '/admin/stats',
    name: 'AdminStats',
    component: () => import('../views/admin/Stats.vue'),
    meta: { requiresAuth: true }
  },
  {
    path: '/worker',
    name: 'Worker',
    redirect: '/worker/login'
  },
  {
    path: '/worker/login',
    name: 'WorkerLogin',
    component: () => import('../views/worker/Login.vue')
  },
  {
    path: '/worker/home',
    name: 'WorkerHome',
    component: () => import('../views/worker/Home.vue'),
    meta: { requiresAuth: true }
  },
  {
    path: '/worker/order/:id',
    name: 'WorkerOrderDetail',
    component: () => import('../views/worker/OrderDetail.vue'),
    meta: { requiresAuth: true }
  },
  {
    path: '/worker/profile',
    name: 'WorkerProfile',
    component: () => import('../views/worker/Profile.vue'),
    meta: { requiresAuth: true }
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

router.beforeEach((to, from, next) => {
  const isAdmin = localStorage.getItem('isAdmin') === 'true'
  const isWorker = localStorage.getItem('isWorker') === 'true'
  const userId = localStorage.getItem('userId')
  const workerId = localStorage.getItem('workerId')

  if (to.meta.requiresAuth) {
    if (to.path.startsWith('/admin')) {
      if (!userId || !isAdmin) {
        next('/admin/login')
        return
      }
    } else if (to.path.startsWith('/worker')) {
      if (!workerId || !isWorker) {
        next('/worker/login')
        return
      }
    } else {
      if (!userId || isAdmin || isWorker) {
        next('/user/login')
        return
      }
    }
    next()
  } else {
    next()
  }
})

export default router
