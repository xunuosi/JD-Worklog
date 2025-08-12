<template>
  <Shell>
    <el-card>
      <template #header><b>录入工时</b></template>

      <!-- 新增工时 -->
      <el-form :model="form" label-width="80px" @submit.prevent>
        <el-form-item label="项目">
          <el-select v-model="form.projectId" filterable clearable placeholder="搜索/选择项目" style="width: 260px">
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
          <el-input v-model="form.content" type="textarea" rows="3" />
        </el-form-item>

        <el-form-item>
          <el-button type="primary" @click="submit">提交</el-button>
        </el-form-item>
      </el-form>
    </el-card>

    <el-card class="mt-4">
      <template #header><b>我的工时记录</b></template>

      <el-table :data="list" stripe>
        <el-table-column prop="date" label="日期" width="140">
          <template #default="{ row }">{{ (row.date || '').slice(0, 10) }}</template>
        </el-table-column>
        <el-table-column prop="project.name" label="项目" width="200">
          <template #default="{ row }">{{ row.project?.name }}</template>
        </el-table-column>
        <el-table-column prop="hours" label="小时" width="100" />
        <el-table-column prop="content" label="内容" />
        <el-table-column label="操作" width="160">
          <template #default="{ row }">
            <el-button link type="primary" @click="openEdit(row)">编辑</el-button>
            <el-divider direction="vertical" />
            <el-popconfirm width="190" title="确认删除该条记录？" @confirm="del(row.id)">
              <template #reference>
                <el-button link type="danger">删除</el-button>
              </template>
            </el-popconfirm>
          </template>
        </el-table-column>
      </el-table>

      <el-pagination class="mt-4" background layout="total, prev, pager, next, sizes" :total="total"
        v-model:current-page="page" v-model:page-size="pageSize" :page-sizes="[5, 10, 20, 50]" @current-change="load"
        @size-change="() => { page = 1; load() }" />
    </el-card>

    <!-- 编辑弹窗 -->
    <el-dialog v-model="editVisible" title="编辑工时" width="520px">
      <el-form :model="editForm" label-width="80px">
        <el-form-item label="项目">
          <el-select v-model="editForm.project_id" style="width: 100%">
            <el-option v-for="p in projects" :key="p.id" :label="p.name" :value="p.id" />
          </el-select>
        </el-form-item>
        <el-form-item label="日期">
          <el-date-picker v-model="editForm.date" type="date" value-format="YYYY-MM-DD" style="width: 100%" />
        </el-form-item>
        <el-form-item label="小时数">
          <el-input-number v-model="editForm.hours" :min="0" :step="0.5" />
        </el-form-item>
        <el-form-item label="内容">
          <el-input v-model="editForm.content" type="textarea" rows="3" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="editVisible = false">取消</el-button>
        <el-button type="primary" @click="saveEdit">保存</el-button>
      </template>
    </el-dialog>
  </Shell>
</template>

<script setup lang="ts">
import { onMounted, ref } from 'vue'
import Shell from '../components/Shell.vue'
import http from '../api/http'
import { ElMessage } from 'element-plus'

type Project = { id: number; name: string }
const projects = ref<Project[]>([])

const form = ref({
  projectId: undefined as number | undefined,
  date: '',
  hours: 8,
  content: ''
})

const list = ref<any[]>([])
const total = ref(0)
const page = ref(1)
const pageSize = ref(10)

const load = async () => {
  // 项目
  const { data: ps } = await http.get('/projects')
  projects.value = ps
  if (ps.length && !form.value.projectId) form.value.projectId = ps[0].id

  // 我的工时（分页）
  const { data } = await http.get('/timesheets/mine', {
    params: { page: page.value, page_size: pageSize.value }
  })
  list.value = data.items
  total.value = data.total
}

const submit = async () => {
  if (!form.value.projectId || !form.value.date) {
    ElMessage.warning('请选择项目与日期')
    return
  }
  await http.post('/timesheets', {
    project_id: form.value.projectId,
    date: form.value.date,
    hours: form.value.hours,
    content: form.value.content
  })
  form.value.content = ''
  ElMessage.success('提交成功')
  await load()
}

// ====== 编辑相关 ======
const editVisible = ref(false)
const editForm = ref<{ id: number; project_id: number; date: string; hours: number; content: string }>({
  id: 0,
  project_id: 0,
  date: '',
  hours: 0,
  content: ''
})

const openEdit = (row: any) => {
  editForm.value = {
    id: row.id,
    project_id: row.project?.id ?? row.project_id,
    date: (row.date || '').slice(0, 10),
    hours: row.hours,
    content: row.content
  }
  editVisible.value = true
}

const saveEdit = async () => {
  await http.put(`/timesheets/${editForm.value.id}`, {
    project_id: editForm.value.project_id,
    date: editForm.value.date,
    hours: editForm.value.hours,
    content: editForm.value.content
  })
  ElMessage.success('已保存')
  editVisible.value = false
  await load()
}

const del = async (id: number) => {
  await http.delete(`/timesheets/${id}`)
  ElMessage.success('已删除')
  await load()
}

onMounted(load)
</script>

<style scoped>
.mt-4 {
  margin-top: 1rem;
}
</style>
