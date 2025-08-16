import axios from 'axios'
import { useAuthStore } from '../store/auth'

const base = import.meta.env.VITE_API_BASE || '/api'
const http = axios.create({ baseURL: base })

http.interceptors.request.use(config => {
  const auth = useAuthStore()
  if (auth.token) config.headers.Authorization = `Bearer ${auth.token}`
  return config
})

export default http
