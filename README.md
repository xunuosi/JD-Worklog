# Worklog System

- Frontend: Vue3 + Vite + Pinia + Vue Router + Axios
- Backend: Go (Gin + GORM) + JWT
- DB: MySQL

## Run locally
1. Start MySQL and create DB `worklog`
2. Backend
```bash
cd backend
export DB_HOST=127.0.0.1 DB_PORT=3306 DB_USER=worklog DB_PASS=worklog DB_NAME=worklog JWT_SECRET=devsecretchangeit CORS_ORIGINS=http://localhost:5173
go mod tidy
go run ./cmd/server
```
3. Frontend
```bash
cd frontend
npm i
npm run dev
```
Accounts:
- admin/admin123
- alice/alice123

## CSV Export
Admin -> Report:
- Click **导出 CSV** after selecting date range.
- Backend endpoint: `GET /api/admin/reports/project-totals.csv?from=YYYY-MM-DD&to=YYYY-MM-DD`
# JD-Worklog
