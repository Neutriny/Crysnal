import { createRouter, createWebHistory } from 'vue-router'
import DashboardView from '@/views/DashboardView.vue'
import InstanceListView from '@/views/instance/ListView.vue'
import InstanceCreateView from '@/views/instance/CreateView.vue'
import FrpListView from '@/views/frp/ListView.vue'
import Profile from '@/views/setting/ProfileView.vue'

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
      component: InstanceListView
    },
    {
      path: '/instance/create',
      name: 'instance-create',
      component: InstanceCreateView
    },
    {
      path: '/frp',
      redirect: '/frp/list',
    },
    
    {
      path: '/frp/list',
      name: 'frp-list',
      component: FrpListView
    },
    {
      path: '/setting',
      redirect: '/setting/profile',
    },
    {
      path: '/setting/profile',
      name: 'profile',
      component: Profile
    },
    {
      path: '/about',
      name: 'about',
      // 懒加载：只有访问 /about 时才加载这个页面的代码
      component: () => import('@/views/AboutView.vue'),
    },
  ],
})

export default router
