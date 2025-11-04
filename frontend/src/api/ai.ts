import http from './http';

export interface GenerateReportPayload {
  date: string;
  extra_content: string;
  worklogs: string[];
}

export interface GenerateReportResponse {
  report: string;
}

export const aiApi = {
  generateReport: (payload: GenerateReportPayload) => {
    return http.post<GenerateReportResponse>('/ai/generate-report', payload);
  },
};
