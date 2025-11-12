export interface User {
  id: number;
  username: string;
  nickname: string;
  role: 'admin' | 'user';
}

export interface WorkPlan {
  id: number;
  user_id: number;
  project_id: number;
  start_date: string;
  end_date: string;
  content: string;
  user: User;
  project: Project;
}

export interface Timesheet {
  id: number;
  createdAt: string;
  updatedAt: string;
  deletedAt: string | null;
  userId: number;
  projectId: number;
  project: Project;
  date: string;
  hours: number;
  content: string;
}

export interface Project {
  id: number;
  name: string;
  desc: string;
  contract_num: string;
  is_active: boolean;
}
