import http from './http'
import type { Project } from '../types'

export interface BackfillRequest {
  project_id: number
  user_id: number
  total_days: number
  start_date: string
  end_date: string
  content: string
}

export const backfillTimesheets = (data: BackfillRequest) => http.post('/admin/timesheets/backfill', data)

export const getBackfillHistory = (params: { page: number, page_size: number }) => http.get('/admin/timesheets/backfill/history', { params })

export const deleteBackfill = (id: number) => http.delete(`/admin/timesheets/backfill/${id}`)

export const getAllProjects = () => http.get<Project[]>('/admin/allprojects')

export const getAllUsers = () => http.get('/admin/users')

export const require2FA = (userId: number, require: boolean) => http.post(`/admin/users/${userId}/require-2fa`, { require })
