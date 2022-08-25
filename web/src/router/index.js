import { createRouter, createWebHistory } from 'vue-router'
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
        component: () => import("views/home/Hello.vue")
      },
      {
        path: '/dashboard',
        name: 'Dashboard',
        component: () => import("views/Dashboard.vue")
      },
      {
        path: '/about',
        name: 'About',
        component: () => import("views/others/About.vue")
      },
      {
        path: '/users',
        name: 'Users',
        component: () => import("views/user/User.vue")
      },
      {
        path: '/users/:id',
        name: 'UserDetail',
        component: () => import("views/user/UserDetail.vue")
      },
      {
        path: '/user_groups',
        name: 'UserGroups',
        component: () => import("views/user/UserGroup.vue")
      },
      {
        path: '/groups',
        name: 'Groups',
        component: () => import("views/user/SysGroup.vue")
      },
      {
        path: '/groups/:id',
        name: 'GroupDetail',
        component: () => import("views/user/GroupDetail.vue")
      },
      {
        path: '/rbac',
        name: 'RBAC',
        component: () => import("views/auth/RBAC.vue")
      },
      {
        path: '/apps',
        name: 'Applications',
        component: () => import("views/container/Application.vue")
      },
      {
        path: '/apps/:id/exec',
        name: 'ContainerTerminal',
        component: () => import("views/container/ContainerTerminal.vue")
      },
      {
        path: '/apps/:id/proxy:proxyPath(.*)',
        name: 'Proxy',
        component: () => import("views/container/Proxy.vue"),
      },
      {
        path: '/namespaces',
        name: 'Namespace',
        component: () => import("views/kube/Namespace.vue")
      },
      {
        path: '/workloads',
        name: 'Workload',
        component: () => import("views/kube/Workload.vue")
      },
      {
        path: '/pods',
        name: 'Pod',
        component: () => import("views/kube/Pod.vue")
      },
      {
        path: '/services',
        name: 'Service',
        component: () => import("views/kube/Service.vue")
      },
      {
        path: '/ingresses',
        name: 'Ingress',
        component: () => import("views/kube/Ingress.vue")
      },
      {
        path: '/namespaces/:namespace/pods/:pod/exec',
        name: 'PodTerminal',
        component: () => import("views/kube/PodTerminal.vue")
      },
      {
        path: '/webcode',
        name: 'WebCode',
        component: () => import('views/others/WebCode.vue')
      },
      {
        path: '/markdown',
        name: 'MarkDown',
        component: () => import('views/others/MarkDownEditor.vue')
      },
      {
        path: '/404',
        name: '404',
        component: () => import('views/others/404.vue')
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
    component: () => import("views/auth/Login.vue")
  },
  {
    path: '/oauth',
    name: 'OAuth',
    component: () => import("views/auth/OAuth.vue")
  },
  {
    path: '/docs',
    name: 'Document',
    component: () => import("views/doc/List.vue"),
    redirect: '/docs/introduce',
    children: [
      {
        path: ':page',
        name: 'DocDetail',
        component: () => import('views/doc/Detail.vue')
      }
    ]
  },
  {
    path: '/posts',
    name: 'Post',
    component: () => import('views/post/Home.vue'),
    redirect: '/posts/list',
    children: [
      {
        path: 'editor',
        name: 'Editor',
        component: () => import('views/post/Editor.vue')
      },
      {
        path: 'editor/:id',
        name: 'Draft',
        component: () => import('views/post/Editor.vue')
      },
      {
        path: 'list',
        name: 'PostList',
        component: () => import('views/post/List.vue')
      },
      {
        path: ':post',
        name: 'PostDetail',
        component: () => import('views/post/Detail.vue')
      }
    ]
  },
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

router.beforeEach((to, from, next) => {
  let isAuthenticated = false;
  let user = getUser();
  if (user && user.name) {
    isAuthenticated = true;
  }

  if (!isAuthenticated && to.name !== 'Login' && to.name !== 'OAuth') next({ name: 'Login' })
  else if(isAuthenticated && (to.name == 'Login' || to.name == 'OAuth' )) next({ name: 'Index'})
  else next()
})

export default router