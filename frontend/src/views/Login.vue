<template>
  <div class="login-wrap">
    <el-card class="box-card">
      <template #header>
        <div class="card-header">登录</div>
      </template>
      <el-form :model="form" @submit.prevent>
        <el-form-item label="用户名"><el-input v-model="form.username" /></el-form-item>
        <el-form-item label="密码"><el-input v-model="form.password" type="password" /></el-form-item>
        <el-form-item>
          <el-button type="primary" class="w-full" @click="submit">登录</el-button>
        </el-form-item>
      </el-form>
      <div class="mt-2 text-muted">示例：admin/admin123 或 alice/alice123</div>
    </el-card>
  </div>
</template>
<script setup lang="ts">
import { reactive } from 'vue'
import { useAuthStore } from '../store/auth'
import { ElMessage } from 'element-plus'
const auth = useAuthStore()
const form = reactive({ username: 'admin', password: 'admin123' })
const submit = async () => {
  try { await auth.login(form.username, form.password) } 
  catch { ElMessage.error('登录失败') }
}
</script>
<style scoped>
.login-wrap { min-height: 100vh; display:flex; align-items:center; justify-content:center; background: linear-gradient(135deg,#f5f7fa,#e6ecf5); }
.box-card { width: 360px; }
.card-header { font-weight: 600; }
.w-full { width: 100%; }
.mt-2 { margin-top: .5rem; }
.text-muted { color: #909399; }
</style>