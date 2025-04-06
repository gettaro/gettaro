# EMS.dev

A modern email marketing and automation platform built with Go, React, and PostgreSQL.

## Prerequisites

- Docker and Docker Compose
- Node.js 18+ (for local development)
- Go 1.21+ (for local development)
- Auth0 account

## Setup

1. Clone the repository:

```bash
git clone https://github.com/yourusername/ems.dev.git
cd ems.dev
```

2. Set up Auth0:

   - Create a new Auth0 application
   - Configure the following settings:
     - Allowed Callback URLs: `http://localhost:3000`
     - Allowed Logout URLs: `http://localhost:3000`
     - Allowed Web Origins: `http://localhost:3000`
   - Create an API in Auth0 and configure the necessary scopes

3. Create environment files:

   - Copy `.env.example` to `.env` in both frontend and backend directories
   - Update the environment variables with your Auth0 credentials

4. Start the application:

```bash
docker-compose up --build
```

The application will be available at:

- Frontend: http://localhost:3000
- Backend API: http://localhost:8080
- PostgreSQL: localhost:5432

## Development

### Frontend Development

```bash
cd frontend
npm install
npm run dev
```

### Backend Development

```bash
cd backend
go mod download
go run main.go
```

## Features

- User authentication with Auth0
- Protected API endpoints
- PostgreSQL database integration
- Modern React frontend with Vite
- Go backend with Gin framework

## License

MIT
