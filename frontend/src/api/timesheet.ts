import http from './http';
import { Timesheet } from '../types'; // Assuming you have a Timesheet type defined

export const timesheetApi = {
  getMineByDate: (date: string) => {
    return http.post<Timesheet[]>('/timesheets/mine/by-date', { date });
  },
};
