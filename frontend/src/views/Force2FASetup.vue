<template>
  <div class="flex items-center justify-center min-h-screen bg-gray-100">
    <el-card header="必须启用两步验证 (2FA)">
      <div class="p-6">
        <p class="mb-4">您的账户已被管理员要求启用两步验证以增强安全性。</p>
        <p class="mb-4">请使用您的 Authenticator 应用扫描下面的二维码，然后输入6位验证码以完成设置。</p>

        <div v-if="qrCode" class="flex flex-col items-center">
          <img :src="`data:image/png;base64,${qrCode}`" alt="2FA QR Code" class="mb-4" />
          <p class="mb-2">或者，手动输入密钥:</p>
          <el-tag type="info">{{ secret }}</el-tag>
        </div>

        <el-form @submit.prevent="verify" class="mt-6">
          <el-form-item label="验证码">
            <OtpInput ref="otpInputRef" v-model="token" />
          </el-form-item>
          <el-form-item>
            <el-button type="primary" @click="verify" native-type="submit">验证并启用</el-button>
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

const qrCode = ref('')
const secret = ref('')
const token = ref('')
const auth = useAuthStore()
const otpInputRef = ref<InstanceType<typeof OtpInput> | null>(null)

onMounted(async () => {
  try {
    const { data } = await http.post('/2fa/setup')
    qrCode.value = data.qr_code
    secret.value = data.secret
    otpInputRef.value?.focus()
  } catch (error) {
    ElMessage.error('无法获取2FA设置，请重试')
  }
})

const verify = async () => {
  if (!token.value || token.value.length !== 6) {
    ElMessage.error('请输入有效的6位验证码')
    return
  }
  try {
    const { data } = await http.post('/2fa/verify', { token: token.value })
    ElMessage.success('两步验证已成功启用！')

    // Update store with the new full token
    auth.token = data.token
    auth.role = data.role
    localStorage.setItem('token', data.token)
    localStorage.setItem('role', data.role)
    auth.is2FASetupRequired = false // Reset flag

    await router.push(data.role === 'admin' ? '/admin/projects' : '/timesheet')
  } catch (error: any) {
    ElMessage.error(error.response?.data?.error || '验证失败')
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
.flex-col {
  flex-direction: column;
}
.p-6 {
  padding: 1.5rem;
}
.mb-4 {
  margin-bottom: 1rem;
}
.mt-6 {
  margin-top: 1.5rem;
}
</style>
