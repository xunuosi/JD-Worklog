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
  ID: number;
  CreatedAt: string;
  UpdatedAt: string;
  DeletedAt: string | null;
  UserID: number;
  ProjectID: number;
  Project: Project;
  Date: string;
  Hours: number;
  Content: string;
}

export interface Project {
  id: number;
  name: string;
  desc: string;
  contract_num: string;
  is_active: boolean;
}
