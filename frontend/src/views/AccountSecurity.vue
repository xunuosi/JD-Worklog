<template>
  <Shell>
    <el-card>
      <template #header><b>账号安全</b></template>
      <el-form :model="form" label-width="100px" @submit.prevent>
        <el-form-item label="旧密码"><el-input v-model="form.old_password" type="password" /></el-form-item>
        <el-form-item label="新密码"><el-input v-model="form.new_password" type="password" /></el-form-item>
        <el-form-item>
          <el-button type="primary" @click="submit">修改密码</el-button>
        </el-form-item>
      </el-form>
    </el-card>
  </Shell>
</template>
<script setup lang="ts">
import { reactive } from 'vue'
import http from '../api/http'
import Shell from '../components/Shell.vue'
import { ElMessage } from 'element-plus'

const form = reactive({ old_password: '', new_password: '' })
const submit = async () => {
  if (!form.old_password || !form.new_password) return ElMessage.warning('请输入完整信息')
  await http.post('/change-password', form)
  ElMessage.success('密码修改成功')
  form.old_password = ''; form.new_password = ''
}
</script>
