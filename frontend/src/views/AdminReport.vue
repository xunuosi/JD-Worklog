<template>
  <Shell>
    <el-card>
      <template #header><b>工时报表（项目总工时）</b></template>
      <div class="flex gap-3 mb-3 items-center">
        <el-date-picker v-model="range" type="daterange" range-separator="至" start-placeholder="起始" end-placeholder="结束" />
        <el-button type="primary" @click="run">查询</el-button>
        <el-button @click="exportCsv" :disabled="!range">导出 CSV</el-button>
      </div>
      <el-table :data="rows" size="small" stripe>
        <el-table-column prop="project_id" label="项目ID" width="120" />
        <el-table-column prop="project_name" label="项目名称" />
        <el-table-column prop="total_hours" label="总工时" width="120" />
      </el-table>
    </el-card>
  </Shell>
</template>
<script setup lang="ts">
import { ref } from 'vue'
import http from '../api/http'
import Shell from '../components/Shell.vue'
import { ElMessage } from 'element-plus'

type Row = { project_id:number; project_name:string; total_hours:number }
const range = ref<[Date,Date] | null>(null)
const rows = ref<Row[]>([])

const fmt = (d: Date) => `${d.getFullYear()}-${String(d.getMonth()+1).padStart(2,'0')}-${String(d.getDate()).padStart(2,'0')}`

const run = async () => {
  if (!range.value) return ElMessage.warning('请选择日期范围')
  const [s, e] = range.value
  const { data } = await http.post('/admin/reports/project-totals', { from: fmt(s), to: fmt(e) })
  rows.value = data
}

const exportCsv = async () => {
  if (!range.value) return
  const [s, e] = range.value
  const { data } = await http.get('/admin/reports/project-totals.csv', { params: { from: fmt(s), to: fmt(e) }, responseType: 'blob' })
  const url = URL.createObjectURL(new Blob([data]))
  const a = document.createElement('a'); a.href = url; a.download = 'project_totals.csv'; a.click(); URL.revokeObjectURL(url)
}
</script>
<style scoped>
.flex { display:flex; }
.gap-3 { gap:.75rem; }
.mb-3 { margin-bottom:.75rem; }
.items-center { align-items:center; }
</style>