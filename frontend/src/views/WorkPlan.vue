<template>
  <Shell>
    <el-card>
      <template #header>
        <div class="flex justify-between">
          <b>工作计划</b>
          <div>
            <el-date-picker
              v-if="authStore.user?.role === 'admin'"
              v-model="fetchDateRange"
              type="daterange"
              range-separator="至"
              start-placeholder="开始日期"
              end-placeholder="结束日期"
              @change="fetchWorkPlans"
              class="mr-4"
            />
            <el-select
              v-if="authStore.user?.role === 'admin'"
              v-model="selectedProject"
              placeholder="选择项目"
              @change="fetchWorkPlans"
              clearable
              filterable
            >
              <el-option
                v-for="project in projects"
                :key="project.id"
                :label="project.name"
                :value="project.id"
              />
            </el-select>
            <el-button type="primary" @click="openCreateDialog" class="ml-4">
              新增计划
            </el-button>
          </div>
        </div>
      </template>
      <FullCalendar :options="calendarOptions" />
    </el-card>

    <el-dialog v-if="showCreateDialog" v-model="showCreateDialog" title="新增计划" @close="resetForm">
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
const showCreateDialog = ref(false);
const showEditDialog = ref(false);
const form = ref<Partial<WorkPlan>>({});
const dateRange = ref<[Date, Date] | null>(null);
const fetchDateRange = ref<[Date, Date] | null>(null);

const calendarOptions = computed(() => ({
  plugins: [dayGridPlugin, interactionPlugin, listPlugin],
  initialView: 'dayGridMonth',
  headerToolbar: {
    left: 'prev,next today',
    center: 'title',
    right: 'dayGridMonth,listWeek',
  },
  events: workPlans.value.map((plan) => ({
    id: plan.id.toString(),
    title: `${plan.user.nickname} - ${plan.project.name}`,
    start: plan.start_date,
    end: plan.end_date,
    extendedProps: {
      content: plan.content,
    },
  })),
  eventClick: (info) => {
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
    if (authStore.user?.role === 'admin') {
      projects.value = await getAllProjects();
    } else {
      projects.value = await getProjects();
    }
  } catch (error) {
    ElMessage.error('获取项目列表失败');
  }
}

async function fetchWorkPlans() {
  try {
    let start, end;
    if (fetchDateRange.value) {
      start = fetchDateRange.value[0].toISOString().split('T')[0];
      end = fetchDateRange.value[1].toISOString().split('T')[0];
    }

    if (selectedProject.value) {
      workPlans.value = await getWorkPlansByProject(selectedProject.value, start, end);
    } else {
      workPlans.value = await getWorkPlans(start, end);
    }
  } catch (error) {
    ElMessage.error('获取工作计划失败');
  }
}

function resetForm() {
  form.value = {};
  dateRange.value = null;
}

async function openCreateDialog() {
  if (projects.value.length === 0) {
    await fetchProjects();
  }
  showCreateDialog.value = true;
}

async function handleCreate() {
  if (!form.value.project_id || !dateRange.value) {
    ElMessage.error('请填写完整信息');
    return;
  }
  form.value.start_date = dateRange.value[0].toISOString().split('T')[0];
  form.value.end_date = dateRange.value[1].toISOString().split('T')[0];

  try {
    await createWorkPlan(form.value);
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
  form.value.start_date = dateRange.value[0].toISOString().split('T')[0];
  form.value.end_date = dateRange.value[1].toISOString().split('T')[0];

  try {
    await updateWorkPlan(form.value);
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
.ml-4 {
  margin-left: 1rem;
}
.mr-4 {
  margin-right: 1rem;
}
</style>
