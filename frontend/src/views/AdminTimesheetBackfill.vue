<template>
  <Shell>
    <el-card>
      <template #header><b>工时补入</b></template>
      <el-form :model="form" label-width="120px" style="max-width: 600px">
        <el-form-item label="项目">
          <el-select v-model="form.project_id" placeholder="请选择项目">
            <el-option v-for="p in projects" :key="p.id" :label="p.name" :value="p.id"></el-option>
          </el-select>
        </el-form-item>
        <el-form-item label="人员">
          <el-select v-model="form.user_id" placeholder="请选择人员">
            <el-option v-for="u in users" :key="u.id" :label="u.nickname" :value="u.id"></el-option>
          </el-select>
        </el-form-item>
        <el-form-item label="总工时(人天)">
          <el-input-number v-model="form.total_days" :min="1"></el-input-number>
        </el-form-item>
        <el-form-item label="时间范围">
          <el-date-picker
            v-model="dateRange"
            type="daterange"
            range-separator="至"
            start-placeholder="开始日期"
            end-placeholder="结束日期"
            value-format="YYYY-MM-DD"
          ></el-date-picker>
        </el-form-item>
        <el-form-item label="工作内容">
          <el-input v-model="form.content" type="textarea"></el-input>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="onSubmit">提交</el-button>
        </el-form-item>
      </el-form>
    </el-card>
  </Shell>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import Shell from '../components/Shell.vue'
import { getAllProjects, getAllUsers, backfillTimesheets, type BackfillRequest } from '../api/admin'

const projects = ref<any[]>([])
const users = ref<any[]>([])
const dateRange = ref<[string, string] | null>(null)

const initialFormState = {
  project_id: null,
  user_id: null,
  total_days: 1,
  start_date: '',
  end_date: '',
  content: '工时补入'
}

const form = reactive<BackfillRequest>({ ...initialFormState })

const resetForm = () => {
  Object.assign(form, initialFormState)
  dateRange.value = null
}

onMounted(async () => {
  const [projRes, userRes] = await Promise.all([getAllProjects(), getAllUsers()])
  projects.value = projRes.data
  users.value = userRes.data
})

const onSubmit = async () => {
  if (!dateRange.value) {
    ElMessage.error('请选择时间范围')
    return
  }
  form.start_date = dateRange.value[0]
  form.end_date = dateRange.value[1]

  ElMessageBox.confirm('确认补入工时？', '提示', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning'
  }).then(async () => {
    try {
      await backfillTimesheets(form)
      ElMessage.success('工时补入成功')
      resetForm()
    } catch (error: any) {
      ElMessage.error(error.response?.data?.error || '工时补入失败')
    }
  })
}
</script>
