<template>
  <Shell>
    <template #title>补入历史</template>
    <el-table :data="history" style="width: 100%">
      <el-table-column prop="operator.nickname" label="操作人"></el-table-column>
      <el-table-column prop="user.nickname" label="补入人员"></el-table-column>
      <el-table-column prop="project.name" label="项目"></el-table-column>
      <el-table-column prop="total_days" label="总人天"></el-table-column>
      <el-table-column label="时间范围">
        <template #default="{ row }">
          {{ new Date(row.start_date).toLocaleDateString() }} - {{ new Date(row.end_date).toLocaleDateString() }}
        </template>
      </el-table-column>
      <el-table-column label="操作时间">
        <template #default="{ row }">
          {{ new Date(row.CreatedAt).toLocaleString() }}
        </template>
      </el-table-column>
      <el-table-column label="操作">
        <template #default="{ row }">
          <el-button type="danger" size="small" @click="onDelete(row.id)">删除</el-button>
        </template>
      </el-table-column>
    </el-table>
  </Shell>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import Shell from '../components/Shell.vue'
import { getBackfillHistory, deleteBackfill } from '../api/admin'

const history = ref<any[]>([])

const fetchHistory = async () => {
  const res = await getBackfillHistory()
  history.value = res.data
}

onMounted(fetchHistory)

const onDelete = (id: number) => {
  ElMessageBox.confirm('确认删除此条补入记录及其关联的工时？此操作不可恢复。', '警告', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning'
  }).then(async () => {
    try {
      await deleteBackfill(id)
      ElMessage.success('删除成功')
      fetchHistory()
    } catch (error: any) {
      ElMessage.error(error.response?.data?.error || '删除失败')
    }
  })
}
</script>
