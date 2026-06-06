import { createRouter, createWebHashHistory } from 'vue-router'

const routes = [
  {
    path: '/login',
    name: 'Login',
    component: () => import('../views/Login.vue'),
    meta: { title: '登录' }
  },
  {
    path: '/',
    component: () => import('../views/Layout.vue'),
    redirect: '/dashboard',
    children: [
      {
        path: 'dashboard',
        name: 'Dashboard',
        component: () => import('../views/Dashboard.vue'),
        meta: { title: '仪表盘' }
      },
      {
        path: 'clusters',
        name: 'Clusters',
        component: () => import('../views/Clusters.vue'),
        meta: { title: '集群管理' }
      },
      {
        path: 'topics',
        name: 'Topics',
        component: () => import('../views/Topics.vue'),
        meta: { title: 'Topic管理' }
      },
      {
        path: 'consumers',
        name: 'Consumers',
        component: () => import('../views/Consumers.vue'),
        meta: { title: '消费组监控' }
      },
      {
        path: 'messages',
        name: 'Messages',
        component: () => import('../views/Messages.vue'),
        meta: { title: '消息追踪' }
      },
      {
        path: 'alerts',
        name: 'Alerts',
        component: () => import('../views/Alerts.vue'),
        meta: { title: '告警中心' }
      },
      {
        path: 'acl',
        name: 'ACL',
        component: () => import('../views/ACL.vue'),
        meta: { title: '权限管理' }
      },
      {
        path: 'ops',
        name: 'Ops',
        component: () => import('../views/Ops.vue'),
        meta: { title: '运维工具' }
      },
      {
        path: 'guide',
        name: 'Guide',
        component: () => import('../views/Guide.vue'),
        meta: { title: '使用指引' }
      }
    ]
  }
]

const router = createRouter({
  history: createWebHashHistory(),
  routes
})

export default router
