import { defineStore } from 'pinia'
import http from '../api/http'
import router from '../router'

export const useAuthStore = defineStore('auth', {
  state: () => ({
    token: localStorage.getItem('token') || '',
    role: localStorage.getItem('role') || '',
    is2FASetupRequired: false,
  }),
  actions: {
    async login(username: string, password: string) {
      const { data } = await http.post('/login', { username, password })

      if (data.force_2fa_setup) {
        this.token = data.token
        localStorage.setItem('token', this.token)
        this.is2FASetupRequired = true
        await router.push('/force-2fa-setup')
        return
      }

      if (data.two_factor_required) {
        this.token = data.token
        localStorage.setItem('token', this.token)
        // Redirect to a page where the user enters their 2FA token.
        // This page needs to be created.
        await router.push('/login-2fa')
        return
      }

      this.token = data.token
      this.role = data.role
      localStorage.setItem('token', this.token)
      localStorage.setItem('role', this.role)
      await router.push(this.role === 'admin' ? '/admin/projects' : '/timesheet')
    },
    logout() {
      this.token = ''
      this.role = ''
      this.is2FASetupRequired = false
      localStorage.removeItem('token')
      localStorage.removeItem('role')
      router.push('/login')
    }
  }
})
