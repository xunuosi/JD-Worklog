<template>
  <Shell>
    <el-card>
      <template #header><b>用户管理</b></template>
      <el-form :inline="true" @submit.prevent>
        <el-form-item><el-input v-model="username" placeholder="新用户名" /></el-form-item>
        <el-form-item><el-input v-model="nickname" placeholder="昵称（可选）" /></el-form-item>
        <el-form-item><el-input v-model="password" type="password" placeholder="初始密码" /></el-form-item>
        <el-form-item><el-button type="primary" @click="create">创建普通用户</el-button></el-form-item>
      </el-form>

      <el-table :data="users" size="small" stripe class="mt-4">
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="username" label="用户名" />
        <el-table-column prop="nickname" label="昵称" />
        <el-table-column prop="role" label="角色" width="120" />
        <el-table-column label="操作" width="220">
          <template #default="{ row }">
            <!-- 其他按钮... -->
            <el-popconfirm title="确认将该用户密码重置为 root？" @confirm="resetPwd(row)">
              <template #reference>
                <el-button type="primary" link :disabled="row.role === 'admin'">重置为 root</el-button>
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

const resetPwd = async (u: { id: number; username: string }) => {
  await http.post('/admin/users/reset-password', { user_id: u.id }) // 不传 new_password，后端默认 root
  ElMessage.success('已重置为 root')
}

type User = { id: number; username: string; role: string; nickname?: string }
const users = ref<User[]>([])
const username = ref('')
const password = ref('')
const nickname = ref('')

const load = async () => { const { data } = await http.get('/admin/users'); users.value = data }

const create = async () => {
  try {
    await http.post('/admin/users', { username: username.value, password: password.value, nickname: nickname.value })
    username.value = ''; password.value = ''; nickname.value = '';
    await load(); ElMessage.success('创建成功')
  } catch (err: any) { ElMessage.error(err.response?.data?.error || '创建失败') }
}

const remove = async (id: number) => { await http.delete(`/admin/users/${id}`); ElMessage.success('已删除'); await load() }

onMounted(load)
</script>
<style scoped>
.mt-4 {
  margin-top: 1rem;
}
</style>