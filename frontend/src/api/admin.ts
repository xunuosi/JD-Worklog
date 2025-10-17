import http from './http'

export interface BackfillRequest {
  project_id: number
  user_id: number
  total_days: number
  start_date: string
  end_date: string
  content: string
}

export const backfillTimesheets = (data: BackfillRequest) => http.post('/admin/timesheets/backfill', data)

export const getBackfillHistory = () => http.get('/admin/timesheets/backfill/history')

export const deleteBackfill = (id: number) => http.delete(`/admin/timesheets/backfill/${id}`)

export const getAllProjects = () => http.get('/admin/allprojects')

export const getAllUsers = () => http.get('/admin/users')
