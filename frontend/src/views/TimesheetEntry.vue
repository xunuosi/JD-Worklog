<template>
  <Shell>
    <el-card>
      <template #header><b>录入工时</b></template>
      <el-form :model="form" label-width="80px" @submit.prevent>
        <el-form-item label="项目">
          <el-select v-model="form.projectId" placeholder="选择项目" style="width: 260px">
            <el-option v-for="p in projects" :key="p.id" :label="p.name" :value="p.id" />
          </el-select>
        </el-form-item>
        <el-form-item label="日期">
          <el-date-picker v-model="form.date" type="date" placeholder="选择日期" value-format="YYYY-MM-DD"
            style="width: 260px" />
        </el-form-item>
        <el-form-item label="小时数">
          <el-input-number v-model="form.hours" :min="0" :step="0.5" />
        </el-form-item>
        <el-form-item label="内容">
          <el-input type="textarea" v-model="form.content" :rows="3" placeholder="本日工作内容" />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="submit">提交</el-button>
        </el-form-item>
      </el-form>
    </el-card>

    <el-card class="mt-4">
      <template #header><b>我的工时记录</b></template>
      <el-table :data="list" size="small" stripe>
        <el-table-column prop="date" label="日期" width="140">
          <template #default="{ row }">{{ row.date?.slice(0, 10) }}</template>
        </el-table-column>
        <el-table-column prop="project.name" label="项目" />
        <el-table-column prop="hours" label="小时" width="100" />
        <el-table-column prop="content" label="内容" />
        <el-table-column label="操作" width="100">
          <template #default="{ row }">
            <el-popconfirm title="确认删除这条记录？" @confirm="del(row.id)">
              <template #reference>
                <el-button type="danger" link>删除</el-button>
              </template>
            </el-popconfirm>
          </template>
        </el-table-column>
      </el-table>
    </el-card>
  </Shell>
</template>
<script setup lang="ts">
import { onMounted, ref } from 'vue'
import http from '../api/http'
import Shell from '../components/Shell.vue'
import { ElMessage } from 'element-plus'

type Project = { id: number; name: string }
const projects = ref<Project[]>([])
const list = ref<any[]>([])
const form = ref({ projectId: undefined as number | undefined, date: '', hours: 8, content: '' })

const load = async () => {
  const { data: ps } = await http.get('/projects')
  projects.value = ps
  if (ps.length && !form.value.projectId) form.value.projectId = ps[0].id
  const { data: mine } = await http.get('/timesheets/mine')
  list.value = mine
}

const submit = async () => {
  if (!form.value.projectId || !form.value.date) return ElMessage.warning('请选择项目与日期')
  await http.post('/timesheets', { project_id: form.value.projectId, date: form.value.date, hours: form.value.hours, content: form.value.content })
  form.value.content = ''
  ElMessage.success('提交成功')
  await load()
}

const del = async (id: number) => { await http.delete(`/timesheets/${id}`); ElMessage.success('已删除'); await load() }

onMounted(load)
</script>
<style scoped>
.mt-4 {
  margin-top: 1rem;
}
</style>