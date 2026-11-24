import { createRouter, createWebHistory } from 'vue-router'
import HomeView from '@/views/HomeView.vue'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      name: 'home',
      component: HomeView,
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
