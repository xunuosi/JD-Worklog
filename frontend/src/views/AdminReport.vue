<template>
  <Shell>
    <el-card>
      <template #header><b>工时报表（项目总工时）</b></template>
      <div class="flex gap-3 mb-3 items-center">
        <el-date-picker v-model="range" type="daterange" range-separator="至" start-placeholder="起始"
          end-placeholder="结束" />
        <el-select v-model="userId" placeholder="全体成员" clearable style="width: 220px">
          <el-option :value="0" label="全体成员" />
          <el-option v-for="u in users" :key="u.id" :label="u.nickname ? `${u.nickname}（${u.username}）` : u.username"
            :value="u.id" />
        </el-select>
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
import { ref, onMounted } from 'vue'
import http from '../api/http'
import Shell from '../components/Shell.vue'
import { ElMessage } from 'element-plus'

type Row = { project_id: number; project_name: string; total_hours: number }
const range = ref<[Date, Date] | null>(null)
const rows = ref<Row[]>([])
const users = ref<Array<{ id: number; username: string; nickname?: string }>>([])
const userId = ref<number>(0) // 0 = 全体成员（默认）

const fmt = (d: Date) => `${d.getFullYear()}-${String(d.getMonth() + 1).padStart(2, '0')}-${String(d.getDate()).padStart(2, '0')}`

const loadUsers = async () => {
  const { data } = await http.get('/admin/users')
  users.value = data
}

const run = async () => {
  if (!range.value) return ElMessage.warning('请选择日期范围')
  const [s, e] = range.value
  const payload: any = { from: fmt(s), to: fmt(e) }
  if (userId.value && userId.value !== 0) payload.user_id = userId.value
  const { data } = await http.post('/admin/reports/project-totals', payload)
  rows.value = data
}

const exportCsv = async () => {
  if (!range.value) return
  const [s, e] = range.value
  const params: any = { from: fmt(s), to: fmt(e) }
  if (userId.value && userId.value !== 0) params.user_id = userId.value
  const { data } = await http.get('/admin/reports/project-export-xlsx', { params, responseType: 'blob' })
  const url = URL.createObjectURL(new Blob([data]))
  const a = document.createElement('a')
  a.href = url
  a.download = `worklog_${params.from}_${params.to}.xlsx`
  a.click()
  URL.revokeObjectURL(url)
}

onMounted(loadUsers)
</script>
<style scoped>
.flex {
  display: flex;
}

.gap-3 {
  gap: .75rem;
}

.mb-3 {
  margin-bottom: .75rem;
}

.items-center {
  align-items: center;
}
</style>