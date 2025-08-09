import { defineStore } from 'pinia'
import http from '../api/http'
import router from '../router'

export const useAuthStore = defineStore('auth', {
  state: () => ({ token: localStorage.getItem('token') || '', role: localStorage.getItem('role') || '' }),
  actions: {
    async login(username: string, password: string) {
      const { data } = await http.post('/login', { username, password })
      this.token = data.token
      this.role = data.role
      localStorage.setItem('token', this.token)
      localStorage.setItem('role', this.role)
      await router.push(this.role === 'admin' ? '/admin/projects' : '/timesheet')
    },
    logout() {
      this.token = ''
      this.role = ''
      localStorage.removeItem('token')
      localStorage.removeItem('role')
      router.push('/login')
    }
  }
})
