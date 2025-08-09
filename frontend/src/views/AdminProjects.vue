<template>
  <Shell>
    <el-card>
      <template #header><b>项目管理</b></template>
      <div class="flex gap-3 mb-3">
        <el-input v-model="name" placeholder="项目名称" style="max-width: 260px" />
        <el-input v-model="desc" placeholder="项目描述" style="max-width: 360px" />
        <el-button type="primary" @click="create">创建</el-button>
      </div>
      <el-table :data="projects" size="small" stripe>
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column label="名称">
          <template #default="{ row }"><el-input v-model="row.name" /></template>
        </el-table-column>
        <el-table-column label="描述">
          <template #default="{ row }"><el-input v-model="row.desc" /></template>
        </el-table-column>
        <el-table-column prop="is_active" label="启用" width="100">
          <template #default="{ row }"><el-switch v-model="row.is_active" /></template>
        </el-table-column>
        <el-table-column label="操作" width="160">
          <template #default="{ row }">
            <el-button type="primary" link @click="save(row)">保存</el-button>
            <el-popconfirm title="确认删除该项目？" @confirm="del(row.id)">
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

type Project = { id:number; name:string; desc:string; is_active:boolean }
const projects = ref<Project[]>([])
const name = ref('')
const desc = ref('')

const load = async () => { const { data } = await http.get('/projects'); projects.value = data }

const create = async () => {
  try { await http.post('/admin/projects', { name: name.value, desc: desc.value }); name.value=''; desc.value=''; await load(); ElMessage.success('创建成功') }
  catch (err:any) { ElMessage.error(err.response?.data?.error || '创建项目失败') }
}

const save = async (p: Project) => { await http.put(`/admin/projects/${p.id}`, { name: p.name, desc: p.desc, is_active: p.is_active }); ElMessage.success('已保存') }

const del = async (id: number) => { await http.delete(`/admin/projects/${id}`); ElMessage.success('已删除'); await load() }

onMounted(load)
</script>

<style scoped>
.flex { display:flex; }
.gap-3 { gap: .75rem; }
.mb-3 { margin-bottom: .75rem; }
</style>