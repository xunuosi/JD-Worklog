<template>
  <Navbar />
  <div style="max-width:820px;margin:20px auto;">
    <h2>项目管理</h2>
    <form @submit.prevent="create">
      <input v-model="name" placeholder="项目名称" required />
      <input v-model="desc" placeholder="项目描述" />
      <button type="submit">创建</button>
    </form>

    <table border="1" cellspacing="0" cellpadding="6" style="margin-top:12px;width:100%;">
      <thead><tr><th>ID</th><th>名称</th><th>描述</th><th>启用</th><th>操作</th></tr></thead>
      <tbody>
        <tr v-for="p in projects" :key="p.id">
          <td>{{p.id}}</td>
          <td><input v-model="p.name" /></td>
          <td><input v-model="p.desc" /></td>
          <td><input type="checkbox" v-model="p.is_active" /></td>
          <td>
            <button @click="save(p)">保存</button>
            <button @click="del(p.id)">删除</button>
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

type Project = { id:number; name:string; desc:string; is_active:boolean }
const projects = ref<Project[]>([])
const name = ref('')
const desc = ref('')

const load = async () => {
  const { data } = await http.get('/projects')
  projects.value = data
}

const create = async () => {
  try {
    await http.post('/admin/projects', { name: name.value, desc: desc.value })
    name.value = ''
    desc.value = ''
    await load()
  } catch (err: any) {
    if (err.response?.data?.error) {
      alert(err.response.data.error)
    } else {
      alert('创建项目未知错误')
    }
  }
}

const save = async (p: Project) => {
  await http.put(`/admin/projects/${p.id}`, { name: p.name, desc: p.desc, is_active: p.is_active })
}

const del = async (id: number) => {
  await http.delete(`/admin/projects/${id}`)
  await load()
}

onMounted(load)
</script>
