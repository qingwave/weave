import { createRouter, createWebHashHistory } from 'vue-router'
import { getUser } from '@/utils'

const routes = [
  {
    path: '/',
    name: 'Home',
    component: () => import("views/Home.vue"),
    redirect: '/index',
    children: [
      {
        path: '/index',
        name: 'Index',
        component: () => import("views/Hello.vue")
      },
      {
        path: '/dashboard',
        name: 'Dashboard',
        component: () => import("views/Dashboard.vue")
      },
      {
        path: '/about',
        name: 'About',
        component: () => import("views/About.vue")
      },
      {
        path: '/users',
        name: 'Users',
        component: () => import("views/User.vue")
      },
      {
        path: '/apps',
        name: 'Applications',
        component: () => import("views/Application.vue")
      },
      {
        path: '/apps/:id/exec',
        name: 'Terminal',
        component: () => import("views/Terminal.vue")
      },
      {
        path: '/apps/:id/proxy:proxyPath(.*)',
        name: 'Proxy',
        component: () => import("views/Proxy.vue"),
      },
      {
        path: '/404',
        name: '404',
        component: () => import('views/404.vue')
      },
      {
        path: '/:pathMatch(.*)',
        redirect: '/404'
      }
    ]
  },
  {
    path: '/login',
    name: 'Login',
    component: () => import("views/Login.vue")
  }
]

const router = createRouter({
  history: createWebHashHistory(),
  routes
})

router.beforeEach((to, from, next) => {
  let isAuthenticated = false;
  let user = getUser();
  if (user.length > 0) {
    isAuthenticated = true;
  }

  if (!isAuthenticated && to.name !== 'Login') next({ name: 'Login' })
  else next()
})

export default router