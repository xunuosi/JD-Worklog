// vite.config.ts
import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue' // 如果是 React 替换为 react()

export default defineConfig({
  plugins: [vue()],
  base: '/worklog/',
  server: {
    host: '0.0.0.0',
    port: 5173,
    proxy: {
      // 任何以 /api 开头的请求，都转发到本地后端
      '/api': {
        target: 'http://localhost:10081',
        changeOrigin: true,
        // 如果你的后端是以 /api 开头注册的路由（你的 Gin 就是），不要重写
        // 如果后端没有 /api 前缀，可以用下面这行把 /api 去掉：
        // rewrite: path => path.replace(/^\/api/, '')
      },
    },
  },
  // 如果你的前端是部署在根路径，不需要设置 base
  // base: '/'
})
