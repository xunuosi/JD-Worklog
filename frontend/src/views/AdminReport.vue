<template>
  <Navbar />
  <div style="max-width:820px;margin:20px auto;">
    <h2>工时报表（项目总工时）</h2>
    <form @submit.prevent="run">
      <label>起始</label>
      <input type="date" v-model="from" required />
      <label>结束</label>
      <input type="date" v-model="to" required />
      <button type="submit">查询</button>
      <button type="button" @click="exportCsv" :disabled="!from || !to">导出 CSV</button>
    </form>

    <table v-if="rows.length" border="1" cellspacing="0" cellpadding="6" style="margin-top:12px;width:100%;">
      <thead><tr><th>项目ID</th><th>项目名称</th><th>总工时</th></tr></thead>
      <tbody>
        <tr v-for="r in rows" :key="r.project_id">
          <td>{{r.project_id}}</td>
          <td>{{r.project_name}}</td>
          <td>{{r.total_hours}}</td>
        </tr>
      </tbody>
    </table>
  </div>
</template>
<script setup lang="ts">
import { ref } from 'vue'
import http from '../api/http'
import Navbar from '../components/Navbar.vue'

type Row = { project_id:number; project_name:string; total_hours:number }
const from = ref<string>('')
const to = ref<string>('')
const rows = ref<Row[]>([])

const run = async () => {
  const { data } = await http.post('/admin/reports/project-totals', { from: from.value, to: to.value })
  rows.value = data
}

const exportCsv = async () => {
  const { data } = await http.get('/admin/reports/project-totals.csv', {
    params: { from: from.value, to: to.value },
    responseType: 'blob'
  })
  const url = URL.createObjectURL(new Blob([data]))
  const a = document.createElement('a')
  a.href = url
  a.download = 'project_totals.csv'
  document.body.appendChild(a)
  a.click()
  a.remove()
  URL.revokeObjectURL(url)
}
</script>
