import { createRouter, createWebHistory } from 'vue-router'
import AboutView from '@/views/AboutView.vue'
import DashboardView from '@/views/DashboardView.vue'
import InstanceListView from '@/views/instance/ListView.vue'
import InstanceCreateView from '@/views/instance/CreateView.vue'
import FrpListView from '@/views/frp/ListView.vue'
import ProfileView from '@/views/setting/ProfileView.vue'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      redirect: '/dashboard/base',
    },
    {
      path: '/dashboard',
      redirect: '/dashboard/base',
    },
    {
      path: '/dashboard/base',
      name: 'dashboard-base',
      component: DashboardView,
    },
    {
      path: '/instance',
      redirect: '/instance/list',
    },
    {
      path: '/instance/list',
      name: 'instance-list',
      component: InstanceListView,
    },
    {
      path: '/instance/create',
      name: 'instance-create',
      component: InstanceCreateView,
    },
    {
      path: '/frp',
      redirect: '/frp/list',
    },

    {
      path: '/frp/list',
      name: 'frp-list',
      component: FrpListView,
    },
    {
      path: '/setting',
      redirect: '/setting/profile',
    },
    {
      path: '/setting/profile',
      name: 'profile',
      component: ProfileView,
    },
    {
      path: '/about',
      name: 'about',
      component: AboutView,
    },
    {
      // 未匹配的路由统一重定向到首页
      path: '/:pathMatch(.*)*',
      redirect: '/',
    },
  ],
})

export default router
