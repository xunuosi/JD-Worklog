<template>
  <div style="max-width:360px;margin:60px auto;">
    <h2>登录</h2>
    <form @submit.prevent="submit">
      <div>
        <label>用户名</label>
        <input v-model="username" required />
      </div>
      <div>
        <label>密码</label>
        <input v-model="password" type="password" required />
      </div>
      <button type="submit">登录</button>
    </form>
    <p style="margin-top:12px;color:#666;">示例账号：admin/admin123，或 alice/alice123</p>
  </div>
</template>
<script setup lang="ts">
import { ref } from 'vue'
import { useAuthStore } from '../store/auth'
const auth = useAuthStore()
const username = ref('admin')
const password = ref('admin123')
const submit = async () => { await auth.login(username.value, password.value) }
</script>
