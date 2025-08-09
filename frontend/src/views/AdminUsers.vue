<template>
  <Navbar />
  <div style="max-width:640px;margin:20px auto;">
    <h2>用户管理</h2>
    <form @submit.prevent="create">
      <input v-model="username" placeholder="新用户名" required />
      <input v-model="password" type="password" placeholder="初始密码" required />
      <button type="submit">创建普通用户</button>
    </form>

    <table border="1" cellspacing="0" cellpadding="6" style="margin-top:12px;width:100%;">
      <thead><tr><th>ID</th><th>用户名</th><th>角色</th><th>操作</th></tr></thead>
      <tbody>
        <tr v-for="u in users" :key="u.id">
          <td>{{u.id}}</td>
          <td>{{u.username}}</td>
          <td>{{u.role}}</td>
          <td>
            <button @click="remove(u.id)" :disabled="u.role==='admin'">删除</button>
          </td>
        </tr>
      </tbody>
    </table>
  </div>
</template>

<script setup lang="ts">
import { onMounted, ref } from 'vue'
import http from '../api/http'
import Navbar from '../components/Navbar.vue'

type User = { id:number; username:string; role:string }
const users = ref<User[]>([])
const username = ref('')
const password = ref('')

const load = async () => {
  const { data } = await http.get('/admin/users')
  users.value = data
}

const create = async () => {
  try {
    await http.post('/admin/users', { username: username.value, password: password.value })
    username.value = ''
    password.value = ''
    await load()
  } catch (err: any) {
    alert(err.response?.data?.error || '创建失败')
  }
}

const remove = async (id: number) => {
  if (!confirm('确认删除该用户？此操作为软删除，可后续从数据库恢复。')) return
  try {
    await http.delete(`/admin/users/${id}`)
    await load()
  } catch (err: any) {
    alert(err.response?.data?.error || '删除失败')
  }
}

onMounted(load)
</script>