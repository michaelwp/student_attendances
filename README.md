# Student Attendance System

A comprehensive full-stack student attendance management system with a REST API backend built with Go/Fiber and a modern React/TypeScript frontend. Features include teacher management, student enrollment, class organization, attendance tracking, and a comprehensive admin dashboard with real-time statistics.

## Features

### Backend API Features
- **JWT Authentication**: Secure multi-user authentication (Admin, Teacher, Student) with Redis caching
- **Role-Based Access Control**: Different permission levels for admins, teachers, and students
- **Teacher Management**: Full CRUD operations with photo upload, password reset, and status management
- **Class Management**: Create and manage classes with homeroom teacher assignments
- **Student Management**: Complete student lifecycle management with class assignments and status tracking
- **Attendance Tracking**: Comprehensive attendance recording with multiple status types and filtering
- **Absence Requests**: Student absence request workflow with teacher/admin approval
- **Photo Management**: Profile photo upload/retrieval for teachers and students via AWS S3
- **Admin Dashboard**: Real-time statistics and comprehensive user management
- **Password Security**: Automated password reset and secure update functionality
- **Status Management**: Active/inactive status control for all user types
- **RESTful API**: Clean, well-documented REST endpoints with consistent patterns
- **Database Migrations**: Automated database schema management
- **Comprehensive Documentation**: Full Swagger/OpenAPI specification with updated routes

### Frontend Web Application Features
- **Modern React/TypeScript**: Built with React 19 and TypeScript for type safety
- **Multi-language Support**: Full internationalization (i18n) with English and Indonesian
- **Dark Mode**: System preference detection with manual toggle
- **Responsive Design**: Mobile-first responsive design with Tailwind CSS
- **Admin Dashboard**: Real-time statistics with active/inactive breakdowns for all entities
- **Student Homepage**: Public attendance marking page for students with simple ID/password authentication
- **Student Dashboard**: Comprehensive student portal with profile management, attendance statistics, and absent request functionality
- **Complete CRUD Management**: Full create, read, update, delete operations for all business entities
- **Photo Upload**: Drag-and-drop photo upload with preview and validation
- **Password Management**: Update/reset password functionality for teachers, students, and admins
- **Status Management**: Quick toggle for activating/deactivating users
- **Advanced Data Tables**: Sortable, paginated tables with custom actions
- **State Management**: Zustand for efficient state management with persistence
- **Authentication Flow**: Secure login/logout with JWT token management
- **Professional UI**: Clean, modern interface with accessibility features
- **Error Handling**: Comprehensive error handling with user-friendly messages

## Architecture

The application follows clean architecture principles with clear separation of concerns:

### Backend Structure
```
   cmd/student_attendance/     # Application entry point
   internal/
      api/
         handlers/          # HTTP handlers with interfaces
         middleware/        # JWT authentication, CORS, logging
         router.go          # Route definitions and middleware setup
      config/                # Configuration management
      models/                # Data models and business entities
      repository/            # Data access layer with interfaces
   db/                        # Database migrations and utilities
   docs/                      # API documentation (Swagger/OpenAPI)
   pkg/                       # Shared utilities (JWT, password hashing)
```

### Frontend Structure
```
   web/
      src/
         components/        # Reusable React components (Layout, Forms)
         pages/            # Page components (Dashboard, Login)
         stores/           # Zustand state management stores
         services/         # API service layer with type-safe calls
         types/            # TypeScript type definitions
         hooks/            # Custom React hooks (theme, auth)
         i18n/            # Internationalization setup and translations
         styles/           # Global styles and Tailwind configuration
```

## Prerequisites

- Go 1.21 or higher
- PostgreSQL 12 or higher
- Redis (for JWT token caching and session management)
- Node.js 18 or higher (for web application)
- npm or yarn (for frontend dependencies)
- AWS S3 bucket (for photo storage)
- Make (optional, for using Makefile commands)

## Installation & Setup

1. **Clone the repository**
   ```bash
   git clone <repository-url>
   cd student_attendance
   ```

2. **Install dependencies**
   ```bash
   go mod download
   ```

3. **Set up environment variables**
   ```bash
   cp .env.example .env
   # Edit .env with your database credentials
   ```

4. **Environment Configuration**
   Create or update your `.env` file:
   ```env
   ENVIRONMENT=development
   PORT=8080
   
   # Database configuration
   DB_HOST=localhost
   DB_PORT=5432
   DB_NAME=student_attendances
   DB_USER=postgres
   DB_PASSWORD=postgres
   DB_SSL_MODE=disable
   
   # Connection pool settings
   MAX_CONNECTIONS=10
   MAX_IDLE_CONNECTIONS=5
   MAX_LIFETIME_CONNECTIONS=30
   
   # AWS S3 configuration (for photo uploads)
   AWS_S3_REGION=us-east-1
   AWS_ACCESS_KEY_ID=your-access-key-id
   AWS_SECRET_ACCESS_KEY=your-secret-access-key
   AWS_S3_BUCKET=your-bucket-name
   
   # Redis configuration
   REDIS_HOST=localhost
   REDIS_PORT=6379
   REDIS_PASSWORD=
   REDIS_DB=0
   
   # JWT and encryption (for future authentication)
   JWT_SECRET=your-secret-key
   SALT=your-salt
   ROUND=12
   
   # Logging
   LOG_LEVEL=debug
   ```

5. **Set up PostgreSQL database**
   ```bash
   # Create database
   createdb student_attendances
   
   # Run migrations
   make migration-up
   # OR
   go run db/migration.go up
   ```

6. **Build and run the application**
   ```bash
   # Build
   go build -o bin/student_attendance cmd/student_attendance/main.go
   
   # Run
   ./bin/student_attendance
   # OR
   go run cmd/student_attendance/main.go
   ```

7. **Access the API Documentation**
   Once the server is running, you can access:
   - **API Health Check**: http://localhost:8080/health
   - **Swagger UI**: http://localhost:8080/swagger/index.html
   - **Interactive API Testing**: Use the Swagger UI to test all endpoints

## API Endpoints

### Base URL
- Development: `http://localhost:8080`
- All API endpoints are prefixed with `/api/v1`

### 🚨 **IMPORTANT: Authentication Required for All Endpoints** 

**ALL API endpoints require valid JWT authentication except for:**
- `GET /health` (Health check)
- `POST /api/v1/auth/login` (Login)

### 🔒 **Authentication Required**
**IMPORTANT**: All API endpoints require authentication (JWT token) except for:
- `GET /health` - Health check endpoint
- `POST /api/v1/auth/login` - User login endpoint

**Authentication Methods**: Include JWT token via:
- **Authorization Header**: `Authorization: Bearer <token>`
- **HTTP-Only Cookie**: Automatically set after login

### 🛡️ **Role-Based Access Control**
The API implements role-based access control with three user types:

| User Type | Access Level | Can Access |
|-----------|-------------|------------|
| **Admin** | Full Access | All endpoints including admin management |
| **Teacher** | Limited | Teachers, Classes, Students, Attendances, Absent Requests |
| **Student** | Restricted | Limited access to Students, Attendances, Absent Requests |

**Special Restrictions:**
- **Admin endpoints** (`/admins/*`) - Admin authentication required
- **Absent Request endpoints** (`/absent-requests/*`) - Student or Teacher authentication required
- **All other endpoints** - Any authenticated user can access

### Public Endpoints (No Authentication Required)
- `GET /health` - Check API health status
- `POST /api/v1/attendance/mark` - Student self-attendance marking (student ID + password)

### Authentication Endpoints
- `POST /api/v1/auth/login` - User login (admin, teacher, or student) - Returns JWT token
- `POST /api/v1/auth/logout` - User logout (requires authentication) - Invalidates JWT token

### Teachers (🔒 Authentication Required)
- `POST /api/v1/teachers` - Create a new teacher
- `GET /api/v1/teachers` - Get all teachers (paginated)
- `GET /api/v1/teachers/{id}` - Get teacher by database ID
- `GET /api/v1/teachers/teacher-id/{teacherId}` - Get teacher by teacher ID
- `PUT /api/v1/teachers/{id}` - Update teacher
- `DELETE /api/v1/teachers/{id}` - Delete teacher
- `PUT /api/v1/teachers/{id}/photo` - Upload teacher profile photo
- `GET /api/v1/teachers/{id}/photo` - Get teacher profile photo (signed URL)
- `PUT /api/v1/teachers/teacher-id/{teacherId}/reset-password` - Reset teacher password (generates new password)
- `PUT /api/v1/teachers/teacher-id/{teacherId}/password` - Update teacher password (user provides old and new password)

### Classes (🔒 Authentication Required)
- `POST /api/v1/classes` - Create a new class
- `GET /api/v1/classes` - Get all classes (paginated)
- `GET /api/v1/classes/{id}` - Get class by ID
- `GET /api/v1/classes/teacher-id/{teacherId}` - Get classes by teacher
- `PUT /api/v1/classes/{id}` - Update class
- `DELETE /api/v1/classes/{id}` - Delete class

### Students (🔒 Authentication Required)
- `POST /api/v1/students` - Create a new student
- `GET /api/v1/students` - Get all students (paginated)
- `GET /api/v1/students/{id}` - Get student by database ID
- `GET /api/v1/students/student-id/{studentId}` - Get student by student ID
- `GET /api/v1/students/class-id/{classId}` - Get students by class
- `PUT /api/v1/students/{id}` - Update student
- `DELETE /api/v1/students/{id}` - Delete student
- `PUT /api/v1/students/{id}/photo` - Upload student profile photo
- `GET /api/v1/students/{id}/photo` - Get student profile photo (signed URL)
- `PUT /api/v1/students/student-id/{studentId}/reset-password` - Reset student password (generates new password)
- `PUT /api/v1/students/student-id/{studentId}/password` - Update student password (user provides old and new password)

### Student Dashboard (🔒 Student Authentication Required)
- `GET /api/v1/student/profile` - Get authenticated student's profile with attendance statistics
- `PUT /api/v1/student/profile` - Update authenticated student's profile (first name, last name, email, phone)
- `PUT /api/v1/student/password` - Update authenticated student's password (with old password verification)
- `GET /api/v1/student/absent-requests` - Get authenticated student's absent requests (paginated)
- `POST /api/v1/student/absent-requests` - Create new absent request for authenticated student
- `DELETE /api/v1/student/absent-requests/{id}` - Delete student's own absent request (pending requests only)

### Teacher Dashboard (🔒 Teacher Authentication Required)
- `GET /api/v1/teacher/profile` - Get authenticated teacher's profile with assigned classes and statistics
- `PUT /api/v1/teacher/password` - Update authenticated teacher's password (with old password verification)
- `GET /api/v1/absent-requests/current-teacher` - Get absent requests from students in teacher's classes (paginated)
- `PUT /api/v1/absent-requests/absent-request-id/{id}/approve` - Approve a student's absent request
- `PUT /api/v1/absent-requests/absent-request-id/{id}/reject` - Reject a student's absent request

### Attendances (🔒 Authentication Required)
- `POST /api/v1/attendances` - Create attendance record
- `GET /api/v1/attendances/all` - Get all attendance records (paginated)
- `GET /api/v1/attendances/attendances-id/{id}` - Get attendance by database ID
- `GET /api/v1/attendances/student-id/{studentId}` - Get attendance by student
- `GET /api/v1/attendances/class-id/{classId}` - Get attendance by class
- `GET /api/v1/attendances/date-range?start_date=YYYY-MM-DD&end_date=YYYY-MM-DD` - Get attendance by date range
- `PUT /api/v1/attendances/attendances-id/{id}` - Update attendance record
- `DELETE /api/v1/attendances/attendances-id/{id}` - Delete attendance record (soft delete)

### Absent Requests (🔒 Authentication Required - Student/Teacher Only)
- `POST /api/v1/absent-requests` - Create absence request
- `GET /api/v1/absent-requests/{id}` - Get absent request by ID
- `GET /api/v1/absent-requests/student-id/{studentId}` - Get requests by student
- `GET /api/v1/absent-requests/class-id/{classId}` - Get requests by class
- `GET /api/v1/absent-requests/pending` - Get all pending requests
- `PATCH /api/v1/absent-requests/{id}/status` - Update request status
- `DELETE /api/v1/absent-requests/{id}` - Delete absent request

### Admins (🔒 Authentication Required - Admin Only)
- `POST /api/v1/admins` - Create a new admin
- `GET /api/v1/admins` - Get all admins (paginated)
- `GET /api/v1/admins/{id}` - Get admin by database ID
- `GET /api/v1/admins/email/{email}` - Get admin by email
- `PUT /api/v1/admins/{id}` - Update admin
- `DELETE /api/v1/admins/{id}` - Delete admin
- `PUT /api/v1/admins/{id}/password` - Update admin password (with old password verification)
- `PUT /api/v1/admins/{id}/status` - Set admin active status (activate/deactivate)

## Data Models

### Teacher
```json
{
  "id": 1,
  "teacher_id": "TCH001",
  "first_name": "John",
  "last_name": "Doe",
  "email": "john.doe@school.com",
  "phone": "+1234567890",
  "created_at": "2024-01-01T00:00:00Z",
  "updated_at": "2024-01-01T00:00:00Z"
}
```

### Class
```json
{
  "id": 1,
  "name": "Grade 10A",
  "homeroom_teacher": "TCH001",
  "description": "Advanced mathematics class",
  "created_at": "2024-01-01T00:00:00Z",
  "updated_at": "2024-01-01T00:00:00Z"
}
```

### Student
```json
{
  "id": 1,
  "student_id": "STU001",
  "classes_id": 1,
  "first_name": "Jane",
  "last_name": "Smith",
  "email": "jane.smith@student.com",
  "phone": "+1234567890",
  "created_at": "2024-01-01T00:00:00Z",
  "updated_at": "2024-01-01T00:00:00Z"
}
```

### Attendance
```json
{
  "id": 1,
  "student_id": "STU001",
  "class_id": 1,
  "date": "2024-01-15T00:00:00Z",
  "status": "present",
  "description": "Student was on time",
  "time_in": "2024-01-15T08:00:00Z",
  "time_out": "2024-01-15T15:30:00Z",
  "created_by": 1,
  "updated_by": 1,
  "created_at": "2024-01-01T00:00:00Z",
  "updated_at": "2024-01-01T00:00:00Z",
  "deleted_at": null,
  "deleted_by": null
}
```

**Attendance Status Options:**
- `present`: Student was present and on time
- `absent`: Student was absent without prior notice
- `late`: Student arrived late but attended class
- `excused`: Student was absent with valid excuse/permission

**Field Descriptions:**
- `student_id`: Reference to the student's unique identifier
- `class_id`: Reference to the class database ID
- `date`: The date of attendance (ISO 8601 format)
- `status`: Current attendance status (see options above)
- `description`: Optional notes about the attendance record
- `time_in`: Time when student checked in (auto-recorded for self-marking)
- `time_out`: Time when student checked out (future feature)
- `created_by`: ID of user who created the record (admin/teacher ID for manual entry, student ID for self-marking)
- `updated_by`: ID of user who last updated the record
- `deleted_at`: Timestamp when record was soft-deleted (null if active)
- `deleted_by`: ID of admin who deleted the record

### Absent Request
```json
{
  "id": 1,
  "student_id": "STU001",
  "class_id": 1,
  "request_date": "2024-01-15",
  "reason": "Medical appointment",
  "status": "pending",
  "created_at": "2024-01-01T00:00:00Z",
  "updated_at": "2024-01-01T00:00:00Z"
}
```

**Request Status Options:**
- `pending`: Request is waiting for approval
- `approved`: Request has been approved
- `rejected`: Request has been rejected

### Admin
```json
{
  "id": 1,
  "email": "admin@school.com",
  "last_login": "2024-01-15T10:30:00Z",
  "is_active": true,
  "created_at": "2024-01-01T00:00:00Z",
  "updated_at": "2024-01-01T00:00:00Z"
}
```

**Admin Status Options:**
- `is_active: true`: Admin account is active and can log in
- `is_active: false`: Admin account is deactivated and cannot log in

### Authentication Models

#### Login Request
```json
{
  "user_type": "admin",
  "user_id": "admin@school.com",
  "password": "securepassword123"
}
```

**User Type Options:**
- `admin`: Use email as user_id
- `teacher`: Use teacher_id as user_id
- `student`: Use student_id as user_id

#### Login Response
```json
{
  "translate_key": "success.login_successful",
  "message": "Login successful",
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user_type": "admin",
  "user_id": "admin@school.com",
  "expires_at": 1640995200
}
```

**Additional Response Headers:**
- `Set-Cookie`: Secure HTTP-only cookie containing the JWT token
  - `HttpOnly=true` - Cannot be accessed by JavaScript
  - `Secure=true` - Only transmitted over HTTPS
  - `SameSite=Strict` - CSRF protection
  - `Expires` - Same as token expiration (1 hour)

#### JWT Token Claims
```json
{
  "user_id": "admin@school.com",
  "user_type": "admin",
  "exp": 1640995200,
  "iat": 1640991600
}
```

## Development Commands

The project includes a Makefile with useful development commands:

```bash
# Database migrations
make migration-up     # Run migrations up
make migration-down   # Run migrations down

# Code quality
make lint            # Run golangci-lint
make lint-fix        # Run golangci-lint with auto-fix

# Git hooks
make install-hooks   # Install pre-commit hooks
make uninstall-hooks # Remove pre-commit hooks

# Documentation
make swagger         # Generate Swagger documentation
```

### Manual Commands

```bash
# Build the application
go build -o bin/student_attendance cmd/student_attendance/main.go

# Run tests
go test ./...

# Run with coverage
go test -cover ./...

# Format code
go fmt ./...

# Vet code
go vet ./...

# Run migrations manually
go run db/migration.go up
go run db/migration.go down
```

## API Documentation

### Swagger/OpenAPI Documentation

The API comes with integrated Swagger UI for interactive documentation and testing:

#### **Live Interactive Documentation**
- **Swagger UI**: `http://localhost:8080/swagger/index.html` (when server is running)
- **API Documentation**: Interactive interface to test all endpoints
- **Model Schemas**: Complete data model definitions
- **Try it out**: Execute API calls directly from the browser

#### **Static Documentation Files**
- **OpenAPI 3.0 Spec**: `docs/swagger.yaml`
- **Swagger JSON**: `docs/swagger.json`  
- **Generated Docs**: `docs/docs.go`

#### **Usage Options**
1. **Interactive Testing**: Visit `/swagger/` when the server is running
2. **Swagger Editor**: Copy `docs/swagger.yaml` to [Swagger Editor](https://editor.swagger.io/)
3. **Postman Import**: Import `docs/swagger.json` directly into Postman
4. **Client SDK Generation**: Use `swagger-codegen` with the spec files

#### **Regenerating Documentation**
When you modify API handlers or add new endpoints:
```bash
# Regenerate Swagger docs
swag init -g cmd/student_attendance/main.go --output docs --parseDependency --parseInternal

# Or use make command (if added to Makefile)
make swagger
```

### Example API Calls

⚠️ **Note**: All examples below (except login) require authentication. Include your JWT token in requests.

#### Authentication Examples

#### Admin Login
```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "user_type": "admin",
    "user_id": "admin@school.com",
    "password": "securepassword123"
  }'
```

#### Teacher Login
```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "user_type": "teacher",
    "user_id": "TCH001",
    "password": "securepassword123"
  }'
```

#### Student Login
```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "user_type": "student",
    "user_id": "STU001",
    "password": "securepassword123"
  }'
```

#### Logout (Requires Authentication)
```bash
curl -X POST http://localhost:8080/api/v1/auth/logout \
  -H "Authorization: Bearer YOUR_JWT_TOKEN_HERE"
```

#### Using Authentication for Protected Requests
```bash
# Method 1: Using Authorization Header
curl -X GET http://localhost:8080/api/v1/teachers \
  -H "Authorization: Bearer YOUR_JWT_TOKEN_HERE"

# Method 2: Using HTTP-Only Cookie (automatically set after login)
curl -X GET http://localhost:8080/api/v1/teachers \
  --cookie-jar cookies.txt --cookie cookies.txt
```

#### Create a Teacher (🔒 Authentication Required)
```bash
curl -X POST http://localhost:8080/api/v1/teachers \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{
    "teacher_id": "TCH001",
    "first_name": "John",
    "last_name": "Doe",
    "email": "john.doe@school.com",
    "phone": "+1234567890",
    "password": "securepassword123"
  }'
```

#### Create a Class (🔒 Authentication Required)
```bash
curl -X POST http://localhost:8080/api/v1/classes \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{
    "name": "Grade 10A",
    "homeroom_teacher": "TCH001",
    "description": "Advanced mathematics class"
  }'
```

#### Create a Student (🔒 Authentication Required)
```bash
curl -X POST http://localhost:8080/api/v1/students \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{
    "student_id": "STU001",
    "classes_id": 1,
    "first_name": "Jane",
    "last_name": "Smith",
    "email": "jane.smith@student.com",
    "phone": "+1234567890",
    "password": "securepassword123"
  }'
```

#### Attendance Management Examples (🔒 Authentication Required)

##### Create Attendance Record
```bash
curl -X POST http://localhost:8080/api/v1/attendances \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{
    "student_id": "STU001",
    "class_id": 1,
    "date": "2024-01-15T09:00:00Z",
    "status": "present",
    "description": "Student was on time"
  }'
```

##### Get All Attendance Records (Paginated)
```bash
curl -X GET "http://localhost:8080/api/v1/attendances/all?limit=20&offset=0" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

##### Get Attendance by ID
```bash
curl -X GET http://localhost:8080/api/v1/attendances/attendances-id/1 \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

##### Update Attendance Record
```bash
curl -X PUT http://localhost:8080/api/v1/attendances/attendances-id/1 \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{
    "student_id": "STU001",
    "class_id": 1,
    "date": "2024-01-15T09:00:00Z",
    "status": "late",
    "description": "Student arrived 10 minutes late"
  }'
```

##### Delete Attendance Record (Soft Delete)
```bash
curl -X DELETE http://localhost:8080/api/v1/attendances/attendances-id/1 \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

##### Get Attendance by Student
```bash
curl -X GET http://localhost:8080/api/v1/attendances/student-id/STU001 \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

##### Get Attendance by Class
```bash
curl -X GET http://localhost:8080/api/v1/attendances/class-id/1 \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

##### Get Attendance by Date Range
```bash
curl -X GET "http://localhost:8080/api/v1/attendances/date-range?start_date=2024-01-01&end_date=2024-01-31" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

#### Student Self-Attendance Marking (No Authentication Required)
```bash
curl -X POST http://localhost:8080/api/v1/attendance/mark \
  -H "Content-Type: application/json" \
  -d '{
    "student_id": "STU001",
    "password": "studentpassword123"
  }'
```

#### Create Absence Request (🔒 Student/Teacher Authentication Required)
```bash
curl -X POST http://localhost:8080/api/v1/absent-requests \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{
    "student_id": "STU001",
    "class_id": 1,
    "request_date": "2024-01-20",
    "reason": "Medical appointment"
  }'
```

#### Approve Absence Request
```bash
curl -X PATCH http://localhost:8080/api/v1/absent-requests/1/status \
  -H "Content-Type: application/json" \
  -d '{
    "status": "approved"
  }'
```

#### Upload Teacher Photo
```bash
curl -X PUT http://localhost:8080/api/v1/teachers/1/photo \
  -F "photo=@/path/to/teacher-photo.jpg"
```

#### Get Teacher Photo
```bash
curl -X GET http://localhost:8080/api/v1/teachers/1/photo
```

#### Upload Student Photo
```bash
curl -X PUT http://localhost:8080/api/v1/students/1/photo \
  -F "photo=@/path/to/student-photo.jpg"
```

#### Get Student Photo
```bash
curl -X GET http://localhost:8080/api/v1/students/1/photo
```

#### Reset Teacher Password
```bash
curl -X PUT http://localhost:8080/api/v1/teachers/teacher-id/TCH001/reset-password
```

#### Update Teacher Password
```bash
curl -X PUT http://localhost:8080/api/v1/teachers/teacher-id/TCH001/password \
  -H "Content-Type: application/json" \
  -d '{
    "old_password": "current_password",
    "new_password": "new_secure_password"
  }'
```

#### Reset Student Password
```bash
curl -X PUT http://localhost:8080/api/v1/students/student-id/STU001/reset-password
```

#### Update Student Password
```bash
curl -X PUT http://localhost:8080/api/v1/students/student-id/STU001/password \
  -H "Content-Type: application/json" \
  -d '{
    "old_password": "current_password",
    "new_password": "new_secure_password"
  }'
```

#### Student Dashboard Examples (🔒 Student Authentication Required)

##### Get Student Profile with Statistics
```bash
curl -X GET http://localhost:8080/api/v1/student/profile \
  -H "Authorization: Bearer YOUR_STUDENT_JWT_TOKEN"
```

##### Update Student Profile
```bash
curl -X PUT http://localhost:8080/api/v1/student/profile \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_STUDENT_JWT_TOKEN" \
  -d '{
    "first_name": "Jane",
    "last_name": "Smith",
    "email": "jane.smith@student.com",
    "phone": "+1234567890"
  }'
```

##### Update Student Password (Self-Service)
```bash
curl -X PUT http://localhost:8080/api/v1/student/password \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_STUDENT_JWT_TOKEN" \
  -d '{
    "old_password": "current_password",
    "new_password": "new_secure_password"
  }'
```

##### Get Student's Absent Requests
```bash
curl -X GET "http://localhost:8080/api/v1/student/absent-requests?limit=10&offset=0" \
  -H "Authorization: Bearer YOUR_STUDENT_JWT_TOKEN"
```

##### Create Absent Request
```bash
curl -X POST http://localhost:8080/api/v1/student/absent-requests \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_STUDENT_JWT_TOKEN" \
  -d '{
    "request_date": "2024-01-20",
    "reason": "Medical appointment"
  }'
```

##### Delete Student's Own Absent Request
```bash
curl -X DELETE http://localhost:8080/api/v1/student/absent-requests/1 \
  -H "Authorization: Bearer YOUR_STUDENT_JWT_TOKEN"
```

#### Teacher Dashboard Examples (🔒 Teacher Authentication Required)

##### Get Teacher Profile with Classes and Statistics
```bash
curl -X GET http://localhost:8080/api/v1/teacher/profile \
  -H "Authorization: Bearer YOUR_TEACHER_JWT_TOKEN"
```

##### Update Teacher Password (Self-Service)
```bash
curl -X PUT http://localhost:8080/api/v1/teacher/password \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_TEACHER_JWT_TOKEN" \
  -d '{
    "old_password": "current_password",
    "new_password": "new_secure_password"
  }'
```

##### Get Absent Requests from Teacher's Classes
```bash
curl -X GET "http://localhost:8080/api/v1/absent-requests/current-teacher?limit=10&offset=0" \
  -H "Authorization: Bearer YOUR_TEACHER_JWT_TOKEN"
```

##### Approve Student Absent Request
```bash
curl -X PUT http://localhost:8080/api/v1/absent-requests/absent-request-id/1/approve \
  -H "Authorization: Bearer YOUR_TEACHER_JWT_TOKEN"
```

##### Reject Student Absent Request
```bash
curl -X PUT http://localhost:8080/api/v1/absent-requests/absent-request-id/1/reject \
  -H "Authorization: Bearer YOUR_TEACHER_JWT_TOKEN"
```

#### Create an Admin (🔒 Admin Authentication Required)
```bash
curl -X POST http://localhost:8080/api/v1/admins \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{
    "email": "admin@school.com",
    "password": "securepassword123",
    "is_active": true
  }'
```

#### Get Admin by Email
```bash
curl -X GET http://localhost:8080/api/v1/admins/email/admin@school.com
```

#### Update Admin Password
```bash
curl -X PUT http://localhost:8080/api/v1/admins/1/password \
  -H "Content-Type: application/json" \
  -d '{
    "old_password": "current_password",
    "new_password": "new_secure_password"
  }'
```

#### Set Admin Active Status
```bash
curl -X PUT http://localhost:8080/api/v1/admins/1/status \
  -H "Content-Type: application/json" \
  -d '{
    "is_active": false
  }'
```

## Database Schema

The application uses PostgreSQL with the following main tables:

- **teachers**: Teacher information and credentials
- **classes**: Class information with homeroom teacher references
- **students**: Student information and class assignments
- **attendances**: Daily attendance records
- **absent_requests**: Student absence requests with approval workflow
- **admins**: Administrator accounts with authentication and status management

All tables include:
- Auto-incrementing primary keys
- Created/updated timestamps
- Proper foreign key relationships
- Check constraints for status fields

## Photo Management

The API supports profile photo uploads for both teachers and students with AWS S3 integration:

### Features
- **File Upload**: Multipart form-data upload support
- **S3 Storage**: Photos are stored securely in AWS S3
- **Signed URLs**: Secure photo access with time-limited signed URLs
- **Format Support**: Common image formats (JPEG, PNG, GIF, etc.)
- **Automatic Path Updates**: Database records are updated with S3 keys
- **Photo Retrieval**: Get secure photo URLs with expiration times
- **Error Handling**: Comprehensive validation and error responses

### Requirements
- AWS S3 bucket configured with proper permissions
- Environment variables for AWS credentials set up
- Maximum file size limits enforced by the server

### API Endpoints
- `PUT /api/v1/teachers/{id}/photo` - Upload teacher profile photo
- `GET /api/v1/teachers/{id}/photo` - Get teacher profile photo (signed URL)
- `PUT /api/v1/students/{id}/photo` - Upload student profile photo
- `GET /api/v1/students/{id}/photo` - Get student profile photo (signed URL)

### Usage Examples
```bash
# Upload teacher photo
curl -X PUT http://localhost:8080/api/v1/teachers/1/photo \
  -F "photo=@teacher-photo.jpg"

# Get teacher photo (signed URL)
curl -X GET http://localhost:8080/api/v1/teachers/1/photo

# Upload student photo  
curl -X PUT http://localhost:8080/api/v1/students/1/photo \
  -F "photo=@student-photo.jpg"

# Get student photo (signed URL)
curl -X GET http://localhost:8080/api/v1/students/1/photo
```

### Upload Response Format
```json
{
  "translate_key": "success.photo_uploaded",
  "message": "Photo uploaded successfully",
  "path": "https://bucket-name.s3.region.amazonaws.com/photos/teachers/1/teacher_1_1704067200.jpg"
}
```

### Get Photo Response Format
```json
{
  "translate_key": "success.photo_url_retrieved",
  "message": "Photo URL retrieved successfully",
  "url": "https://bucket-name.s3.region.amazonaws.com/photos/teachers/1/teacher_1_1704067200.jpg?X-Amz-Algorithm=..."
}
```

**Note**: The GET photo endpoint returns a signed URL that expires after a specified time (15 minutes for students, 1 hour for teachers) for security purposes.

### Technical Implementation
The photo management system uses several new repository methods:

#### Repository Methods
- `UpdatePhotoPath(ctx context.Context, id uint, photoPath string) error` - Updates the photo path in the database
- `GetPhotoPath(ctx context.Context, id uint) (string, error)` - Retrieves the photo path from the database

#### S3 Configuration Methods
- `UploadFile(client *s3.Client, key string, body []byte) error` - Uploads file to S3
- `GetObjectURL(key string) string` - Gets public S3 object URL
- `GetSignedURL(client *s3.Client, key string, expires time.Duration) (string, error)` - Generates presigned URLs

#### File Storage Structure
Photos are organized in S3 with the following structure:
```
photos/
├── teachers/
│   └── {teacher_id}/
│       └── teacher_{id}_{timestamp}.{extension}
└── students/
    └── {student_id}/
        └── student_{id}_{timestamp}.{extension}
```

## Password Management

The API provides secure password management functionality for both teachers and students:

### Features
- **Password Reset**: Administrators can reset user passwords and generate new secure passwords
- **Password Update**: Users can change their own passwords by providing old and new passwords
- **Secure Hashing**: All passwords are hashed using bcrypt with configurable salt rounds
- **Password Generation**: System generates secure random passwords when resetting
- **Validation**: Comprehensive validation for password requirements and user existence

### API Endpoints
- `PUT /api/v1/teachers/teacher-id/{teacherId}/reset-password` - Reset teacher password
- `PUT /api/v1/teachers/teacher-id/{teacherId}/password` - Update teacher password
- `PUT /api/v1/students/student-id/{studentId}/reset-password` - Reset student password  
- `PUT /api/v1/students/student-id/{studentId}/password` - Update student password

### Usage Examples

#### Reset Teacher Password (Admin Function)
```bash
curl -X PUT http://localhost:8080/api/v1/teachers/teacher-id/TCH001/reset-password
```

#### Update Teacher Password (User Function)
```bash
curl -X PUT http://localhost:8080/api/v1/teachers/teacher-id/TCH001/password \
  -H "Content-Type: application/json" \
  -d '{
    "old_password": "current_password",
    "new_password": "new_secure_password"
  }'
```

#### Reset Student Password (Admin Function)
```bash
curl -X PUT http://localhost:8080/api/v1/students/student-id/STU001/reset-password
```

#### Update Student Password (User Function)
```bash
curl -X PUT http://localhost:8080/api/v1/students/student-id/STU001/password \
  -H "Content-Type: application/json" \
  -d '{
    "old_password": "current_password", 
    "new_password": "new_secure_password"
  }'
```

### Response Formats

#### Reset Password Response
```json
{
  "translate_key": "success.password.reset",
  "message": "Password reset successfully",
  "newPassword": "GeneratedSecurePassword123"
}
```

#### Update Password Response
```json
{
  "translate_key": "success.password.updated",
  "message": "Password updated successfully"
}
```

### Security Considerations
- Reset password generates a new secure random password
- Update password requires knowledge of the current password
- All passwords are hashed with bcrypt before storage
- Password complexity is handled by the generation algorithm
- User existence is verified before any password operations

## Admin Management

The API provides comprehensive administrator account management functionality:

### Features
- **Admin Account Creation**: Create new admin accounts with email and password
- **Account Status Management**: Activate/deactivate admin accounts
- **Password Management**: Secure password updates with old password verification
- **Admin Lookup**: Find admins by ID or email address
- **Account Management**: Full CRUD operations for admin accounts
- **Security**: Password hashing with bcrypt, account status controls

### API Endpoints
- `POST /api/v1/admins` - Create a new admin
- `GET /api/v1/admins` - Get all admins (paginated)
- `GET /api/v1/admins/{id}` - Get admin by database ID
- `GET /api/v1/admins/email/{email}` - Get admin by email
- `PUT /api/v1/admins/{id}` - Update admin information
- `DELETE /api/v1/admins/{id}` - Delete admin account
- `PUT /api/v1/admins/{id}/password` - Update admin password
- `PUT /api/v1/admins/{id}/status` - Set admin active status

### Admin Account States
- **Active (`is_active: true`)**: Admin can log in and access the system
- **Inactive (`is_active: false`)**: Admin account is disabled and cannot log in

### Usage Examples

#### Create Admin Account
```bash
curl -X POST http://localhost:8080/api/v1/admins \
  -H "Content-Type: application/json" \
  -d '{
    "email": "admin@school.com",
    "password": "securepassword123",
    "is_active": true
  }'
```

#### Deactivate Admin Account
```bash
curl -X PUT http://localhost:8080/api/v1/admins/1/status \
  -H "Content-Type: application/json" \
  -d '{
    "is_active": false
  }'
```

### Response Format
```json
{
  "translate_key": "success.admin_created",
  "message": "Admin created successfully"
}
```

### Security Features
- Email-based authentication and identification
- Secure password hashing with bcrypt
- Account activation/deactivation controls
- Password update requires current password verification
- Admin existence validation before operations

## Authentication & Authorization

The API provides comprehensive JWT-based authentication for three user types: admin, teacher, and student.

### Features
- **JWT Token Authentication**: Secure token-based authentication with 1-hour expiration
- **Multi-User Type Support**: Admin, teacher, and student authentication with different credentials
- **Redis Token Caching**: Tokens are cached in Redis for fast validation and logout functionality
- **Automatic Token Expiration**: Tokens expire after 1 hour for enhanced security
- **Role-Based Access**: Different user types have different access levels
- **Secure HTTP-Only Cookies**: Tokens are also stored in secure HTTP-only cookies
- **Dual Authentication Support**: Both Bearer tokens and cookies are supported
- **Secure Logout**: Token invalidation through Redis cache removal and cookie clearing

### Authentication Flow
1. **Login**: User provides user_type, user_id, and password
2. **Validation**: System validates credentials based on user type:
   - Admin: Uses email and password from admins table
   - Teacher: Uses teacher_id and password from teachers table
   - Student: Uses student_id and password from students table
3. **Token Generation**: JWT token is generated with user information and 1-hour expiration
4. **Redis Caching**: Token is cached in Redis with expiration
5. **Cookie Setting**: Secure HTTP-only cookie is set with the token
6. **Role Assignment**: User role is embedded in JWT claims for access control
7. **Response**: Client receives JWT token in both response body and secure cookie

### Using Authentication

#### Login Process
```bash
# Admin login (use email as user_id)
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "user_type": "admin",
    "user_id": "admin@school.com", 
    "password": "password123"
  }'

# Teacher login (use teacher_id as user_id)
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "user_type": "teacher",
    "user_id": "TCH001",
    "password": "password123"
  }'

# Student login (use student_id as user_id)
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "user_type": "student", 
    "user_id": "STU001",
    "password": "password123"
  }'
```

#### Using JWT Tokens
There are two ways to authenticate with the API:

**Method 1: Authorization Header (Bearer Token)**
```bash
curl -X GET http://localhost:8080/api/v1/teachers \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
```

**Method 2: HTTP-Only Cookie (Automatic)**
```bash
# Cookie is automatically included in subsequent requests after login
curl -X GET http://localhost:8080/api/v1/teachers \
  --cookie-jar cookies.txt --cookie cookies.txt
```

#### Logout
Logout works with either authentication method:

**With Bearer Token:**
```bash
curl -X POST http://localhost:8080/api/v1/auth/logout \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

**With Cookie:**
```bash
curl -X POST http://localhost:8080/api/v1/auth/logout \
  --cookie-jar cookies.txt --cookie cookies.txt
```

**Note**: Logout clears both the Redis cache and the HTTP-only cookie for complete session termination.

### Token Structure
JWT tokens contain the following claims:
- `user_id`: The user's identifier (email for admin, teacher_id for teacher, student_id for student)
- `user_type`: The type of user (admin, teacher, student)
- `exp`: Token expiration timestamp
- `iat`: Token issued at timestamp

### Middleware Features
The authentication middleware provides:
- **JWTMiddleware**: Requires valid JWT token for all protected routes
- **OptionalJWTMiddleware**: Validates token if provided, continues without if not
- **RequireUserType**: Restricts access to specific user types (Admin, Teacher, Student)

**Middleware Implementation:**
- All API route groups use `JWTMiddleware` for authentication
- Admin routes additionally use `RequireUserType("admin")` 
- Absent Request routes use `RequireUserType("student", "teacher")`

### Security Considerations
- Tokens are cached in Redis and validated on each request
- Expired tokens are automatically removed from Redis
- Admin accounts can be deactivated, preventing login even with valid passwords
- All passwords are hashed with bcrypt before storage
- Token expiration is enforced both in JWT claims and Redis cache
- HTTP-only cookies prevent XSS attacks by making tokens inaccessible to JavaScript
- Secure flag ensures cookies are only transmitted over HTTPS
- SameSite=Strict prevents CSRF attacks
- Logout immediately invalidates tokens by removing them from Redis and clearing cookies

### Authentication Error Responses

#### Missing Authentication
```json
{
  "translate_key": "error.token_required",
  "error": "Authorization token is required"
}
```

#### Invalid Token Format
```json
{
  "translate_key": "error.invalid_token_format", 
  "error": "Token must be in format: Bearer <token>"
}
```

#### Expired Token
```json
{
  "translate_key": "error.token_expired",
  "error": "Token has expired"
}
```

#### Insufficient Permissions
```json
{
  "translate_key": "error.insufficient_permissions",
  "error": "Insufficient permissions for this operation"
}
```

**Common HTTP Status Codes:**
- `401 Unauthorized`: Missing, invalid, or expired token
- `403 Forbidden`: Token valid but insufficient permissions for the endpoint
- `200 OK`: Successful authentication and authorized access

## Security Features

- Password fields are automatically excluded from JSON responses
- Input validation on all endpoints
- SQL injection prevention using parameterized queries
- CORS support for cross-origin requests
- Request logging middleware
- Secure file upload with AWS S3 integration

## Error Handling

The API returns consistent error responses:

```json
{
  "error": "Description of the error"
}
```

Common HTTP status codes:
- `200`: Success
- `201`: Created
- `400`: Bad Request (invalid input)
- `404`: Not Found
- `500`: Internal Server Error

## Pagination

List endpoints support pagination with query parameters:
- `limit`: Number of items to return (default: 10, max: 100)
- `offset`: Number of items to skip (default: 0)

Example:
```bash
GET /api/v1/students?limit=20&offset=40
```

## Web Application

The system includes a modern React/TypeScript web application providing a comprehensive admin interface:

### Features
- **Secure Authentication**: JWT-based login with automatic token management
- **Admin Dashboard**: Real-time statistics and system overview with comprehensive entity breakdowns
- **Multi-language Support**: English and Indonesian with automatic language detection
- **Dark Mode**: System preference detection with manual toggle
- **Responsive Design**: Mobile-first design that works on all devices
- **Professional UI**: Clean, modern interface using Tailwind CSS
- **State Management**: Persistent state management with Zustand
- **Type Safety**: Full TypeScript support with comprehensive type definitions

### Dashboard Statistics Display
The admin dashboard provides comprehensive real-time statistics:

- **Entity Overview Cards**: Display total, active, and inactive counts for:
  - Administrators with color-coded status indicators
  - Teachers with active/inactive breakdown
  - Students with enrollment status
  - Total classes in the system

- **Today's Attendance Summary**: Real-time attendance data showing:
  - Total attendance records for today
  - Present students count
  - Absent students count  
  - Late arrivals count

- **Quick Actions**: Navigation shortcuts to manage:
  - Teachers (add, edit, remove)
  - Students (enrollment, management)
  - Classes (creation, assignments)
  - Attendance tracking

### User Interface Features
- **Navigation**: Sidebar navigation with route highlighting
- **Theme Toggle**: Dark/light mode with system preference detection
- **Language Toggle**: Switch between English and Indonesian
- **Loading States**: Professional loading indicators and error handling
- **Responsive Grid**: Adaptive layout for different screen sizes
- **Accessibility**: ARIA labels and keyboard navigation support

### Technical Implementation
- **React 19**: Latest React with concurrent features
- **TypeScript**: Full type safety with strict mode
- **Tailwind CSS**: Utility-first CSS framework with custom design system
- **i18next**: Internationalization with namespace support
- **Zustand**: Lightweight state management with persistence
- **React Hook Form**: Form handling with validation
- **Custom Hooks**: Reusable hooks for theme and authentication

### Getting Started with Web Application

1. **Install dependencies**
   ```bash
   cd web
   npm install
   ```

2. **Start development server**
   ```bash
   npm run dev
   ```

3. **Build for production**
   ```bash
   npm run build
   ```

4. **Access the application**
   - Development: http://localhost:5173
   - Login with admin credentials to access the dashboard

### Environment Configuration
The web application automatically connects to the API at `http://localhost:8080/api/v1`. Update the API base URL in `web/src/services/api.ts` for different environments.

### Authentication Flow
1. Users select their account type (Admin, Teacher, Student)
2. Enter appropriate credentials (email for admin, ID for teacher/student)
3. System validates credentials and generates JWT token
4. Token is stored securely and used for subsequent API calls
5. Dashboard loads with personalized statistics and navigation

### Student Attendance Homepage
A dedicated public page for students to mark their daily attendance:

**Features:**
- **Public Access**: No admin authentication required - accessible at the root URL (`/`)
- **Simple Authentication**: Students enter their Student ID and password
- **One-Click Attendance**: Quick attendance marking with immediate confirmation
- **Success Feedback**: Visual confirmation with student name and timestamp
- **Multi-Student Support**: Ability to mark attendance for multiple students in sequence
- **Responsive Design**: Works seamlessly on mobile devices and tablets
- **Multilingual**: Supports both English and Indonesian interfaces
- **Instructions**: Clear usage instructions for students
- **Dashboard Access**: Direct navigation link to student login page for accessing the comprehensive student portal

**Usage:**
1. Navigate to the root URL of the application
2. Enter Student ID (e.g., "STU001") and password
3. Click "Mark Attendance" to record presence
4. Receive confirmation with name and timestamp
5. Option to mark attendance for another student or refresh

**API Integration:**
- Uses `/api/v1/attendance/mark` endpoint for attendance submission
- Validates student credentials against the database
- Automatically records attendance with "present" status
- Returns student name and success confirmation

### Student Dashboard Portal
A comprehensive authenticated dashboard for students to manage their academic profile and track attendance:

**Features:**
- **Secure Authentication**: JWT-based login using Student ID and password
- **Profile Management**: View and update personal information (name, email, phone)
- **Password Management**: Self-service password updates with current password verification
- **Attendance Statistics**: Real-time attendance tracking with percentage calculations
- **Absent Request Management**: Create, view, and manage absence requests with status tracking
- **Auto-Refresh Lists**: Absent request list automatically updates when new requests are submitted
- **Responsive Interface**: Optimized for both desktop and mobile devices
- **Real-Time Updates**: Live statistics and immediate feedback for all actions
- **Multilingual Support**: Available in English and Indonesian

**Dashboard Sections:**
1. **Overview Tab**: 
   - Student ID and contact information display
   - Attendance rate with visual percentage indicator
   - Present and absent days counters
   - Quick action buttons for common tasks

2. **Profile Tab**:
   - Personal information editor with validation
   - Password update form with security requirements
   - Real-time form validation and error handling

3. **Absent Requests Tab**:
   - Paginated list of all absence requests with status badges
   - Create new absence requests with date and reason validation
   - Real-time list updates when new requests are created (auto-refresh functionality)
   - Delete pending requests (approved/rejected cannot be modified)
   - Status tracking (Pending, Approved, Rejected) with color-coded indicators

### Teacher Dashboard Portal
A dedicated authenticated dashboard for teachers to manage their classes and review student requests:
**Features:**
- **Teacher Authentication**: JWT-based login using Teacher ID and password
- **Profile Management**: View personal information, assigned classes, and teaching statistics
- **Password Management**: Self-service password updates with security validation
- **Class Management**: View assigned classes and total student counts
- **Absent Request Processing**: Review, approve, or reject student absence requests
- **Request Management**: Real-time list of pending, approved, and rejected requests
- **Batch Actions**: Process multiple requests efficiently
- **Student Information**: View detailed request information with student identification
- **Status Tracking**: Complete audit trail of request decisions with timestamps
- **Responsive Design**: Optimized interface for various screen sizes
- **Multilingual Support**: Available in English and Indonesian
**Dashboard Sections:**
1. **Profile Tab**: 
   - Teacher ID and contact information display
   - Assigned classes with student count statistics
   - Pending requests counter and quick stats
   - Professional information and status indicators
2. **Password Tab**:
   - Secure password update form with current password verification
   - Password strength requirements and validation
   - Real-time form validation and security feedback
3. **Absent Requests Tab**:
   - List all student absence requests from assigned classes
   - Filter by status (pending, approved, rejected)
   - Approve or reject requests with single-click actions
   - View detailed request information with student context
   - Real-time updates when actions are performed

**Technical Implementation:**
- **React Components**: Modular architecture with reusable components
- **Form Validation**: Real-time validation using React Hook Form
- **State Management**: Efficient state updates with callback patterns
- **Error Handling**: Comprehensive error states with user-friendly messages
- **Toast Notifications**: Success and error feedback for all operations
- **Modal Dialogs**: Confirmation dialogs for destructive actions
- **Pagination**: Efficient data loading for large datasets

**Usage Flow:**
1. Student logs in using their Student ID and password
2. System redirects to the student dashboard upon successful authentication
3. Dashboard displays personalized statistics and navigation tabs
4. Students can update their profile, change passwords, and manage absent requests
5. All changes are immediately reflected with toast notifications
6. System maintains session state until logout or token expiration

**Security Features:**
- JWT token authentication for all API calls
- Password updates require current password verification
- Students can only access their own data
- Secure logout with token invalidation
- Client-side validation combined with server-side security