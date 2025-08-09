<template>
  <nav style="display:flex;gap:12px;padding:12px;border-bottom:1px solid #eee;align-items:center;">
    <strong>工时系统</strong>
    <span style="flex:1"></span>
    <router-link v-if="role==='admin'" to="/admin/projects">项目管理</router-link>
    <router-link v-if="role==='admin'" to="/admin/report">工时报表</router-link>
    <router-link v-if="role==='admin'" to="/admin/users">用户管理</router-link>
    <router-link to="/timesheet">录入工时</router-link>
    <button @click="logout">退出</button>
  </nav>
</template>
<script setup lang="ts">
import { storeToRefs } from 'pinia'
import { useAuthStore } from '../store/auth'
const auth = useAuthStore()
const { role } = storeToRefs(auth)
const logout = () => auth.logout()
</script>
