<template>
  <Shell>
    <el-card>
      <template #header><b>AI 生成工作日报</b></template>
      <el-form :inline="true" @submit.prevent>
        <el-form-item label="选择日期">
          <el-date-picker v-model="selectedDate" type="date" placeholder="选择日期" value-format="YYYY-MM-DD" @change="fetchWorklogs" />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="generateReport" :loading="loading">生成日报</el-button>
        </el-form-item>
      </el-form>

      <div v-if="worklogs.length === 0 && selectedDate">
        <el-alert title="当前日期没有工作内容" type="info" show-icon :closable="false" />
      </div>

      <el-form v-else-if="worklogs.length > 0" label-position="top">
        <el-form-item label="今日工作内容">
          <ul>
            <li v-for="log in worklogs" :key="log">{{ log }}</li>
          </ul>
        </el-form-item>
        <el-form-item label="附加内容">
          <el-input v-model="extraContent" type="textarea" :rows="3" placeholder="请输入附加内容，用于AI上下文" />
        </el-form-item>
      </el-form>

      <el-form v-if="generatedReport" label-position="top">
        <el-form-item label="AI 生成的日报">
          <el-input v-model="generatedReport" type="textarea" :rows="10" readonly />
        </el-form-item>
        <el-form-item>
          <el-button type="success" @click="copyReport">复制内容</el-button>
        </el-form-item>
      </el-form>
    </el-card>
  </Shell>
</template>

<script setup lang="ts">
import { ref } from 'vue';
import { ElMessage } from 'element-plus';
import { aiApi } from '../api/ai';
import { timesheetApi } from '../api/timesheet';
import { getMyWorkPlans } from '../api/work_plans';
import Shell from '../components/Shell.vue';
import type { WorkPlan, Timesheet } from '../types';

const selectedDate = ref('');
const extraContent = ref('');
const worklogs = ref<string[]>([]);
const generatedReport = ref('');
const loading = ref(false);

// Helper function for formatting worklog entries
function formatWorklogEntry(worklog: Timesheet): string {
  const projectName = worklog.project ? worklog.project.name : '未知项目';
  return `${projectName}: ${worklog.content} (${worklog.hours}小时)`;
}

async function fetchNextDayWorkPlans() {
  if (!selectedDate.value) {
    extraContent.value = '';
    return;
  }

  try {
    const currentDate = new Date(selectedDate.value);
    currentDate.setDate(currentDate.getDate() + 1);
    const nextDay = currentDate.toISOString().split('T')[0];

    const response = await getMyWorkPlans(nextDay, nextDay);
    const plans: WorkPlan[] = response.data;

    if (plans && plans.length > 0) {
      let content = '明日工作计划\n';
      content += plans.map(plan => `${plan.project.name}:${plan.content}`).join('\n');
      extraContent.value = content;
    } else {
      extraContent.value = ''; // Clear if no plans
    }
  } catch (error) {
    console.error('Failed to fetch next day work plans', error);
    extraContent.value = '';
  }
}

async function fetchWorklogs() {
  if (!selectedDate.value) return;

  worklogs.value = [];
  generatedReport.value = '';

  try {
    const response = await timesheetApi.getMineByDate(selectedDate.value);
    worklogs.value = response.data.map(formatWorklogEntry);
  } catch (error) {
    console.error(error);
    ElMessage.error('获取工作内容失败');
  }

  await fetchNextDayWorkPlans();
}

async function generateReport() {
  if (worklogs.value.length === 0) {
    ElMessage.warning('当前日期没有工作内容可用于生成日报');
    return;
  }

  loading.value = true;
  try {
    const response = await aiApi.generateReport({
      date: selectedDate.value,
      extra_content: extraContent.value,
      worklogs: worklogs.value,
    });
    generatedReport.value = response.data.report;
  } catch (error) {
    ElMessage.error('生成日报失败');
  } finally {
    loading.value = false;
  }
}

function copyReport() {
  navigator.clipboard.writeText(generatedReport.value).then(() => {
    ElMessage.success('内容已复制到剪贴板');
  });
}
</script>

<style scoped>
ul {
  list-style-type: none;
  padding-left: 0;
}
li {
  margin-bottom: 5px;
}
</style>
