import http from './http';
import type { Project } from '../types';

export function getProjects() {
  return http.get<Project[]>('/projects');
}
