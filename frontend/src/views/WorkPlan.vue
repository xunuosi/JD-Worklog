<template>
  <Shell>
    <el-card>
      <template #header>
        <div class="flex justify-between items-center">
          <b>工作计划</b>
          <el-button type="primary" @click="openCreateDialog"> 新增计划 </el-button>
        </div>
      </template>
      <div v-if="authStore.role === 'admin'" class="mb-4">
        <el-date-picker
          v-model="fetchDateRange"
          type="daterange"
          range-separator="至"
          start-placeholder="开始日期"
          end-placeholder="结束日期"
          @change="fetchWorkPlans"
          style="margin-right: 20px"
        />
        <el-select
          v-model="selectedProject"
          placeholder="选择项目"
          @change="fetchWorkPlans"
          clearable
          filterable
          style="width: 240px"
        >
          <el-option
            v-for="project in projects"
            :key="project.id"
            :label="project.name"
            :value="project.id"
          />
        </el-select>
      </div>
      <FullCalendar :options="calendarOptions" />
      <div class="legend-container">
        <div v-for="item in legendItems" :key="item.name" class="legend-item">
          <div :style="{ backgroundColor: item.color }" class="legend-color-box"></div>
          <span>{{ item.name }}</span>
        </div>
      </div>
    </el-card>

    <el-dialog v-if="showCreateDialog" v-model="showCreateDialog" title="新增计划">
      <el-form :model="createForm" label-width="100px">
        <el-form-item label="项目">
          <el-select v-model="createForm.project_id" placeholder="选择项目" clearable filterable>
            <el-option
              v-for="project in projects"
              :key="project.id"
              :label="project.name"
              :value="project.id"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="起止日期">
          <el-date-picker
            v-model="createDateRange"
            type="daterange"
            range-separator="至"
            start-placeholder="开始日期"
            end-placeholder="结束日期"
            :disabled-date="disabledDate"
          />
        </el-form-item>
        <el-form-item label="计划内容">
          <el-input v-model="createForm.content" type="textarea" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="showCreateDialog = false">取消</el-button>
        <el-button type="primary" @click="handleCreate">确定</el-button>
      </template>
    </el-dialog>

    <el-dialog v-if="showEditDialog" v-model="showEditDialog" title="编辑计划" @close="resetForm">
      <el-form :model="form" label-width="100px">
        <el-form-item label="项目">
          <el-select v-model="form.project_id" placeholder="选择项目" clearable filterable>
            <el-option
              v-for="project in projects"
              :key="project.id"
              :label="project.name"
              :value="project.id"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="起止日期">
          <el-date-picker
            v-model="dateRange"
            type="daterange"
            range-separator="至"
            start-placeholder="开始日期"
            end-placeholder="结束日期"
          />
        </el-form-item>
        <el-form-item label="计划内容">
          <el-input v-model="form.content" type="textarea" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="showEditDialog = false">取消</el-button>
        <el-button type="primary" @click="handleUpdate">确定</el-button>
        <el-button type="danger" @click="handleDelete">删除</el-button>
      </template>
    </el-dialog>
  </Shell>
</template>
<script setup lang="ts">
import { ref, onMounted, computed } from 'vue';
import FullCalendar from '@fullcalendar/vue3';
import dayGridPlugin from '@fullcalendar/daygrid';
import interactionPlugin from '@fullcalendar/interaction';
import listPlugin from '@fullcalendar/list';
import type { EventClickArg } from '@fullcalendar/core';
import zhCnLocale from '@fullcalendar/core/locales/zh-cn';
import { ElMessage, ElMessageBox } from 'element-plus';
import { useAuthStore } from '../store/auth';
import { getAllProjects } from '../api/admin';
import { getProjects } from '../api/projects';
import {
  getWorkPlans,
  getWorkPlansByProject,
  createWorkPlan,
  updateWorkPlan,
  deleteWorkPlan,
} from '../api/work_plans';
import type { WorkPlan, Project } from '../types';
import Shell from '../components/Shell.vue';

const authStore = useAuthStore();
const projects = ref<Project[]>([]);
const selectedProject = ref<number | undefined>();
const workPlans = ref<WorkPlan[]>([]);

// State for Edit Dialog
const showEditDialog = ref(false);
const form = ref<Partial<WorkPlan>>({});
const dateRange = ref<[Date, Date] | null>(null);

// State for Create Dialog
const showCreateDialog = ref(false);
const createForm = ref<Partial<WorkPlan>>({});
const createDateRange = ref<[Date, Date] | null>(null);

const fetchDateRange = ref<[Date, Date] | null>(null);

function formatDate(date: Date): string {
  const year = date.getFullYear();
  const month = (date.getMonth() + 1).toString().padStart(2, '0');
  const day = date.getDate().toString().padStart(2, '0');
  return `${year}-${month}-${day}`;
}

const colorPalette = [
  '#1abc9c', '#2ecc71', '#3498db', '#9b59b6', '#f1c40f',
  '#e67e22', '#e74c3c', '#16a085', '#27ae60', '#2980b9',
  '#8e44ad', '#f39c12', '#d35400', '#c0392b'
];

function getColorForEntity(id: number): string {
  if (id === null || id === undefined) {
    return '#7f8c8d'; // A neutral default color
  }
  return colorPalette[id % colorPalette.length];
}

const legendItems = computed(() => {
  if (!workPlans.value) return [];
  const isUser = authStore.role !== 'admin';
  const items = new Map<string, string>();

  workPlans.value.forEach((plan) => {
    const name = isUser ? plan.project.name : plan.user.nickname;
    if (!items.has(name)) {
      const id = isUser ? plan.project.id : plan.user.id;
      items.set(name, getColorForEntity(id));
    }
  });

  return Array.from(items, ([name, color]) => ({ name, color }));
});

const calendarOptions = computed(() => ({
  plugins: [dayGridPlugin, interactionPlugin, listPlugin],
  initialView: 'dayGridMonth',
  locale: zhCnLocale,
  headerToolbar: {
    left: 'prev,next today',
    center: 'title',
    right: 'dayGridMonth,listWeek',
  },
  buttonText: {
    today: '今日',
    dayGridMonth: '月视图',
    listWeek: '周视图',
  },
  events: workPlans.value.map((plan) => {
    const isUser = authStore.role !== 'admin';
    const title = `${plan.project.name} - ${plan.content}`;
    const id = isUser ? plan.project.id : plan.user.id;
    const color = getColorForEntity(id);
    return {
      id: plan.id.toString(),
      title: title,
      start: plan.start_date,
      end: new Date(new Date(plan.end_date).setDate(new Date(plan.end_date).getDate() + 1)),
      allDay: true, // Explicitly set to all-day
      backgroundColor: color,
      borderColor: color,
      extendedProps: {
        content: plan.content,
      },
    };
  }),
  eventClick: (info: EventClickArg) => {
    const plan = workPlans.value.find((p) => p.id.toString() === info.event.id);
    if (plan) {
      form.value = { ...plan };
      dateRange.value = [new Date(plan.start_date), new Date(plan.end_date)];
      showEditDialog.value = true;
    }
  },
}));

async function fetchProjects() {
  try {
    if (authStore.role === 'admin') {
      projects.value = (await getAllProjects()).data;
    } else {
      projects.value = (await getProjects()).data;
    }
  } catch (error) {
    ElMessage.error('获取项目列表失败');
  }
}

async function fetchWorkPlans() {
  try {
    let start, end;
    if (fetchDateRange.value) {
      start = formatDate(fetchDateRange.value[0]);
      end = formatDate(fetchDateRange.value[1]);
    }

    if (selectedProject.value) {
      workPlans.value = (await getWorkPlansByProject(selectedProject.value, start, end)).data;
    } else {
      workPlans.value = (await getWorkPlans(start, end)).data;
    }
  } catch (error) {
    ElMessage.error('获取工作计划失败');
  }
}

function resetForm() {
  form.value = {};
  dateRange.value = null;
}

const disabledDate = (time: Date) => {
  const today = new Date();
  today.setHours(0, 0, 0, 0); // Set to the beginning of today
  return time.getTime() < today.getTime(); // Disable dates before the beginning of today
};

async function openCreateDialog() {
  createForm.value = {};
  createDateRange.value = null;
  if (projects.value.length === 0) {
    await fetchProjects();
  }
  showCreateDialog.value = true;
}

async function handleCreate() {
  if (!createForm.value.project_id || !createDateRange.value) {
    ElMessage.error('请填写完整信息');
    return;
  }
  const payload = {
    ...createForm.value,
    start_date: formatDate(createDateRange.value[0]),
    end_date: formatDate(createDateRange.value[1]),
  };

  try {
    await createWorkPlan(payload);
    ElMessage.success('创建成功');
    showCreateDialog.value = false;
    fetchWorkPlans();
  } catch (error) {
    ElMessage.error('创建失败');
  }
}

async function handleUpdate() {
  if (!form.value.id || !form.value.project_id || !dateRange.value) {
    ElMessage.error('请填写完整信息');
    return;
  }
  const payload = {
    id: form.value.id,
    project_id: form.value.project_id,
    content: form.value.content,
    start_date: formatDate(dateRange.value[0]),
    end_date: formatDate(dateRange.value[1]),
  };

  try {
    await updateWorkPlan(payload);
    ElMessage.success('更新成功');
    showEditDialog.value = false;
    fetchWorkPlans();
  } catch (error) {
    ElMessage.error('更新失败');
  }
}

async function handleDelete() {
  if (!form.value.id) return;

  ElMessageBox.confirm('确定删除此计划吗？', '提示', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning',
  })
    .then(async () => {
      try {
        await deleteWorkPlan(form.value.id!);
        ElMessage.success('删除成功');
        showEditDialog.value = false;
        fetchWorkPlans();
      } catch (error) {
        ElMessage.error('删除失败');
      }
    })
    .catch(() => {});
}

onMounted(() => {
  const now = new Date();
  const start = new Date(now.getFullYear(), now.getMonth(), 1);
  const end = new Date(now.getFullYear(), now.getMonth() + 1, 0);
  fetchDateRange.value = [start, end];

  fetchProjects();
  fetchWorkPlans();
});
</script>
<style scoped>
.flex {
  display: flex;
}
.justify-between {
  justify-content: space-between;
}
.items-center {
  align-items: center;
}
.ml-4 {
  margin-left: 1rem;
}
.mr-4 {
  margin-right: 1rem;
}
.mb-4 {
  margin-bottom: 1rem;
}
.gap-4 {
  gap: 1rem;
}
.w-3\.5 {
  width: 0.875rem;
}
.h-3\.5 {
  height: 0.875rem;
}
.mr-1 {
  margin-right: 0.25rem;
}
.rounded-sm {
  border-radius: 0.125rem;
}
.legend-container {
  display: flex;
  flex-wrap: wrap;
  gap: 1rem;
  margin-top: 1rem;
}
.legend-item {
  display: flex;
  align-items: center;
}
.legend-color-box {
  width: 0.875rem;
  height: 0.875rem;
  margin-right: 0.25rem;
  border-radius: 0.125rem;
}
</style>
