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

      <el-divider />

      <div>
        <h3>双因素认证 (2FA)</h3>
        <div v-if="!twoFactorEnabled">
          <p>通过手机 App (例如 Google Authenticator, Authy) 启用两步验证，保护您的账户安全。</p>
          <el-button type="primary" @click="setup2FA">启用 2FA</el-button>
        </div>
        <div v-else>
          <p>双因素认证已启用。</p>
          <el-button type="danger" @click="disable2FA">禁用 2FA</el-button>
        </div>
      </div>
    </el-card>

    <el-dialog v-model="dialogVisible" title="设置两步验证">
      <div style="text-align: center;">
        <p>请使用您的认证应用扫描下方的二维码：</p>
        <div v-if="qrCode">
          <img :src="qrCode" alt="QR Code" />
        </div>
        <p>或者手动输入密钥：<code>{{ secret }}</code></p>
        <el-form @submit.prevent style="max-width: 300px; margin: 0 auto;">
          <el-form-item label-width="0">
            <el-input v-model="otp" placeholder="请输入6位数字验证码" />
          </el-form-item>
          <el-button type="primary" @click="verify2FA">验证并启用</el-button>
        </el-form>
      </div>
    </el-dialog>
  </Shell>
</template>
<script setup lang="ts">
import { reactive, ref, onMounted } from 'vue'
import http from '../api/http'
import Shell from '../components/Shell.vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { useAuthStore } from '../store/auth'

const auth = useAuthStore()

// Password change form
const form = reactive({ old_password: '', new_password: '' })
const submit = async () => {
  if (!form.old_password || !form.new_password) return ElMessage.warning('请输入完整信息')
  await http.post('/change-password', form)
  ElMessage.success('密码修改成功')
  form.old_password = ''; form.new_password = ''
}

// 2FA state
const twoFactorEnabled = ref(false)
const dialogVisible = ref(false)
const qrCode = ref('')
const secret = ref('')
const otp = ref('')

const loadUser = async () => {
  try {
    const { data } = await http.get('/me')
    twoFactorEnabled.value = data.two_factor_enabled
  } catch (e) {
    console.error("Failed to fetch user 2fa status", e)
  }
}

onMounted(loadUser)

const setup2FA = async () => {
  try {
    const { data } = await http.post('/2fa/setup')
    qrCode.value = `data:image/png;base64,${data.qr_code}`
    secret.value = data.secret
    dialogVisible.value = true
  } catch (error) {
    ElMessage.error('无法启动2FA设置')
  }
}

const verify2FA = async () => {
  if (!otp.value) return ElMessage.warning('请输入验证码')
  try {
    await http.post('/2fa/verify', { token: otp.value })
    ElMessage.success('2FA 已成功启用！')
    dialogVisible.value = false
    twoFactorEnabled.value = true // Update UI
    otp.value = ''
  } catch (error) {
    ElMessage.error('验证码错误或无效')
  }
}

const disable2FA = async () => {
    ElMessageBox.confirm(
    '您确定要禁用双因素认证吗？这会降低您账户的安全性。',
    '警告',
    {
      confirmButtonText: '确定禁用',
      cancelButtonText: '取消',
      type: 'warning',
    }
  ).then(async () => {
      try {
        await http.post('/2fa/disable')
        ElMessage.success('2FA 已禁用')
        twoFactorEnabled.value = false
      } catch (error) {
        ElMessage.error('禁用失败')
      }
  })
}
</script>
