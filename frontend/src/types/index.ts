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
  ID: number;
  CreatedAt: string;
  UpdatedAt: string;
  DeletedAt: string | null;
  Name: string;
  Description: string;
}
