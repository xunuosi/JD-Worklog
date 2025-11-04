// vite.config.ts
import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue' // 如果是 React 替换为 react()
import { resolve } from 'path'

export default defineConfig(({ mode }) => ({
  plugins: [vue()],
  base: mode === 'development' ? '/' : '/worklog/',
  resolve: {
    alias: {
      '@': resolve(__dirname, 'src'),
    },
  },
  server: {
    host: '0.0.0.0',
    port: 5173,
    proxy: {
      '/api': { target: 'http://localhost:10081', changeOrigin: true },
    },
  },
}))
