import { createRouter, createWebHistory } from 'vue-router'
import Login from '../views/Login.vue'
import AdminProjects from '../views/AdminProjects.vue'
import AdminReport from '../views/AdminReport.vue'
import TimesheetEntry from '../views/TimesheetEntry.vue'
import { useAuthStore } from '../store/auth'
import AdminUsers from '../views/AdminUsers.vue'
import AccountSecurity from '../views/AccountSecurity.vue'
import AdminTimesheetBackfill from '../views/AdminTimesheetBackfill.vue'
import AdminTimesheetBackfillHistory from '../views/AdminTimesheetBackfillHistory.vue'

const router = createRouter({
  history: createWebHistory('/worklog/'),
  routes: [
    { path: '/', redirect: '/login' },
    { path: '/login', component: Login },
    { path: '/admin/projects', component: AdminProjects },
    { path: '/admin/report', component: AdminReport },
    { path: '/timesheet', component: TimesheetEntry },
    { path: '/admin/users', component: AdminUsers },
    { path: '/account/security', component: AccountSecurity },
    { path: '/admin/timesheet-backfill', component: AdminTimesheetBackfill },
    { path: '/admin/timesheet-backfill-history', component: AdminTimesheetBackfillHistory },
  ]
})

router.beforeEach((to, _from, next) => {
  const auth = useAuthStore()
  if (to.path !== '/login' && !auth.token) return next('/login')
  if (to.path.startsWith('/admin') && auth.role !== 'admin') return next('/timesheet')
  next()
})

export default router
