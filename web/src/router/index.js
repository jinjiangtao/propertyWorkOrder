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
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

router.beforeEach((to, from, next) => {
  const isAdmin = localStorage.getItem('isAdmin') === 'true'
  const userId = localStorage.getItem('userId')

  if (to.meta.requiresAuth) {
    if (!userId) {
      if (to.path.startsWith('/admin')) {
        next('/admin/login')
      } else {
        next('/user/login')
      }
    } else {
      if (to.path.startsWith('/admin') && !isAdmin) {
        next('/user/login')
      } else if (to.path.startsWith('/user') && isAdmin) {
        next('/admin/home')
      } else {
        next()
      }
    }
  } else {
    next()
  }
})

export default router
