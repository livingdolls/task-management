# description task-management

A Task Management application is an application used to organize, monitor, and complete tasks efficiently.

## Run manual backend and frontend task management

- rename config.yaml.example to config.yaml on folder backend/config
- rename .env.local.example to .env.local

- setup backend, on backend folder, run go mod tidy
- run backend, on backend folder, run go run cmd/api/runner.go

- setup frontend, on folder frontend, run npm install
- run frontend, on folder frontend, run npm run dev

## Or run it easier using docker composer

- run docker compose compose up --build

- frontend url http://localhost:5173
- backend url http://localhost:3000
- docs swagger url http://localhost:3000/swagger/index.html

## Technology

- React Vite typescript
- MySql
- Golang Gin

## Database Struktur

- user {id, name, username, password, created_at}
- task {id, user_id, title, description, status, deadline, created_by, created_at}

## Screenshoot

![Preview](https://res.cloudinary.com/dwg1vtwlc/image/upload/v1760672189/dashboard_task_qglqwj.png)
