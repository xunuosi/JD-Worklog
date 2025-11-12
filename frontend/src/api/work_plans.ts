import http from './http';
import type { WorkPlan } from '../types';

export function getWorkPlans(start_date?: string, end_date?: string, user_id?: number) {
  return http.post<WorkPlan[]>('/work-plans/list', { start_date, end_date, user_id });
}

export function getMyWorkPlans(start_date: string, end_date: string) {
  return http.post<WorkPlan[]>('/work-plans/mine', { start_date, end_date });
}

export function getWorkPlansByProject(project_id: number, start_date?: string, end_date?: string, user_id?: number) {
  return http.post<WorkPlan[]>('/work-plans/by-project', { project_id, start_date, end_date, user_id });
}

export function createWorkPlan(data: Partial<WorkPlan>) {
  return http.post<WorkPlan>('/work-plans/create', data);
}

export function updateWorkPlan(data: Partial<WorkPlan>) {
  return http.post<WorkPlan>('/work-plans/update', data);
}

export function deleteWorkPlan(id: number) {
  return http.post('/work-plans/delete', { id });
}
