<template>
  <div class="flex items-center justify-center min-h-screen bg-gray-100">
    <el-card header="两步验证">
      <div class="p-6">
        <p class="mb-4">请输入您的 Authenticator 应用生成的6位验证码。</p>
        <el-form @submit.prevent="login">
          <el-form-item label="验证码">
            <OtpInput ref="otpInputRef" v-model="token" />
          </el-form-item>
          <el-form-item>
            <el-button type="primary" @click="login" native-type="submit">登录</el-button>
          </el-form-item>
        </el-form>
      </div>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import http from '../api/http'
import { useAuthStore } from '../store/auth'
import router from '../router'
import OtpInput from '../components/OtpInput.vue'

const token = ref('')
const auth = useAuthStore()
const otpInputRef = ref<InstanceType<typeof OtpInput> | null>(null)

onMounted(() => {
  otpInputRef.value?.focus()
})

const login = async () => {
  if (!token.value || token.value.length !== 6) {
    ElMessage.error('请输入有效的6位验证码')
    return
  }
  try {
    const { data } = await http.post('/login/2fa', { token: token.value })
    // Manually update the store with the final token
    auth.token = data.token
    auth.role = data.role
    localStorage.setItem('token', data.token)
    localStorage.setItem('role', data.role)
    auth.is2FASetupRequired = false // Reset flag
    await router.push(data.role === 'admin' ? '/admin/projects' : '/timesheet')
  } catch (error: any) {
    ElMessage.error(error.response?.data?.error || '登录失败')
  }
}
</script>

<style scoped>
.min-h-screen {
  min-height: 100vh;
}
.bg-gray-100 {
  background-color: #f7fafc;
}
.items-center {
  align-items: center;
}
.justify-center {
  justify-content: center;
}
.flex {
  display: flex;
}
.p-6 {
  padding: 1.5rem;
}
.mb-4 {
  margin-bottom: 1rem;
}
</style>
