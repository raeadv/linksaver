import { createRouter, createWebHistory } from 'vue-router'


import HomePage from '@/pages/Home.vue'
import Register from '@/pages/auth/Register.vue'
import Login from '@/pages/auth/Login.vue'
import { createDiscreteApi } from 'naive-ui'
import SettingsPage from '@/pages/SettingsPage.vue'


export const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    { path: '/', component: HomePage },
    { path: '/settings', component: SettingsPage },
    { path: '/login', component: Login },
    { path: '/register', component: Register },
  ],
})


const { loadingBar } = createDiscreteApi(['loadingBar'])

let loading = false

router.beforeEach((to, from) => {
  if (to.fullPath !== from.fullPath && !loading) {
    loadingBar.start()
    loading = true
  }
  return true
})

router.afterEach(() => {
  if (loading) {
    setTimeout(() => {
      loadingBar.finish()
    }, 500);
    loading = false
  }
})

router.onError(() => {
  if (loading) {
    loadingBar.error()
    loading = false
  }
})
