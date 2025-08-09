<template>
  <Navbar />
  <div style="max-width:820px;margin:20px auto;">
    <h2>录入工时</h2>
    <form @submit.prevent="submit">
      <label>项目</label>
      <select v-model.number="projectId" required>
        <option v-for="p in projects" :key="p.id" :value="p.id">{{p.name}}</option>
      </select>
      <label>日期</label>
      <input type="date" v-model="date" required />
      <label>小时数</label>
      <input type="number" step="0.5" min="0" v-model.number="hours" required />
      <label>内容</label>
      <textarea v-model="content" rows="3" placeholder="本日工作内容"></textarea>
      <button type="submit">提交</button>
    </form>

    <h3 style="margin-top:24px;">我的工时记录</h3>
    <table border="1" cellspacing="0" cellpadding="6" style="margin-top:8px;width:100%;">
      <thead><tr><th>日期</th><th>项目</th><th>小时</th><th>内容</th><th>操作</th></tr></thead>
      <tbody>
        <tr v-for="t in list" :key="t.id">
          <td>{{t.date?.slice(0,10)}}</td>
          <td>{{t.project?.name}}</td>
          <td>{{t.hours}}</td>
          <td style="white-space:pre-wrap;">{{t.content}}</td>
          <td><button @click="del(t.id)">删除</button></td>
        </tr>
      </tbody>
    </table>
  </div>
</template>
<script setup lang="ts">
import { onMounted, ref } from 'vue'
import http from '../api/http'
import Navbar from '../components/Navbar.vue'

type Project = { id:number; name:string }
const projects = ref<Project[]>([])
const projectId = ref<number>()
const date = ref<string>('')
const hours = ref<number>(8)
const content = ref<string>('')

const list = ref<any[]>([])

const load = async () => {
  const { data: ps } = await http.get('/projects')
  projects.value = ps
  if (ps.length && !projectId.value) projectId.value = ps[0].id
  const { data: mine } = await http.get('/timesheets/mine')
  list.value = mine
}

const submit = async () => {
  await http.post('/timesheets', { project_id: projectId.value, date: date.value, hours: hours.value, content: content.value })
  content.value = ''
  await load()
}

const del = async (id: number) => {
  await http.delete(`/timesheets/${id}`)
  await load()
}

onMounted(load)
</script>
