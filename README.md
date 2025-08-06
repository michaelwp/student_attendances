# Student Attendance API

A comprehensive REST API for managing student attendance, classes, teachers, and absence requests built with Go, Fiber web framework, and PostgreSQL.

## Features

- **Teacher Management**: CRUD operations for teacher accounts with authentication
- **Class Management**: Create and manage classes with homeroom teacher assignments
- **Student Management**: Student enrollment with class assignments
- **Attendance Tracking**: Record and track student attendance with multiple status types
- **Absence Requests**: Students can request absences with approval workflow
- **Photo Upload**: Upload and manage profile photos for teachers and students with AWS S3 integration
- **RESTful API**: Clean, well-documented REST endpoints
- **Database Migrations**: Automated database schema management
- **Comprehensive Documentation**: Full Swagger/OpenAPI specification

## Architecture

The application follows clean architecture principles with clear separation of concerns:

```
   cmd/student_attendance/     # Application entry point
   internal/
      api/                   # HTTP layer
         handlers/          # HTTP handlers with interfaces
         router.go          # Route definitions
      config/                # Configuration management
      models/                # Data models
      repository/            # Data access layer with interfaces
   db/                        # Database migrations and utilities
   docs/                      # API documentation
   web/                       # Static files (if needed)
```

## Prerequisites

- Go 1.24.5 or higher
- PostgreSQL 12 or higher
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

### Health Check
- `GET /health` - Check API health status

### Teachers
- `POST /api/v1/teachers` - Create a new teacher
- `GET /api/v1/teachers` - Get all teachers (paginated)
- `GET /api/v1/teachers/{id}` - Get teacher by database ID
- `GET /api/v1/teachers/teacher-id/{teacherId}` - Get teacher by teacher ID
- `PUT /api/v1/teachers/{id}` - Update teacher
- `DELETE /api/v1/teachers/{id}` - Delete teacher
- `PUT /api/v1/teachers/{id}/photo` - Upload teacher profile photo
- `GET /api/v1/teachers/{id}/photo` - Get teacher profile photo (signed URL)

### Classes
- `POST /api/v1/classes` - Create a new class
- `GET /api/v1/classes` - Get all classes (paginated)
- `GET /api/v1/classes/{id}` - Get class by ID
- `GET /api/v1/classes/teacher-id/{teacherId}` - Get classes by teacher
- `PUT /api/v1/classes/{id}` - Update class
- `DELETE /api/v1/classes/{id}` - Delete class

### Students
- `POST /api/v1/students` - Create a new student
- `GET /api/v1/students` - Get all students (paginated)
- `GET /api/v1/students/{id}` - Get student by database ID
- `GET /api/v1/students/student-id/{studentId}` - Get student by student ID
- `GET /api/v1/students/class-id/{classId}` - Get students by class
- `PUT /api/v1/students/{id}` - Update student
- `DELETE /api/v1/students/{id}` - Delete student
- `PUT /api/v1/students/{id}/photo` - Upload student profile photo
- `GET /api/v1/students/{id}/photo` - Get student profile photo (signed URL)

### Attendances
- `POST /api/v1/attendances` - Create attendance record
- `GET /api/v1/attendances/{id}` - Get attendance by ID
- `GET /api/v1/attendances/student-id/{studentId}` - Get attendance by student
- `GET /api/v1/attendances/class-id/{classId}` - Get attendance by class
- `GET /api/v1/attendances/date-range?start_date=YYYY-MM-DD&end_date=YYYY-MM-DD` - Get attendance by date range
- `PUT /api/v1/attendances/{id}` - Update attendance record
- `DELETE /api/v1/attendances/{id}` - Delete attendance record

### Absent Requests
- `POST /api/v1/absent-requests` - Create absence request
- `GET /api/v1/absent-requests/{id}` - Get absent request by ID
- `GET /api/v1/absent-requests/student-id/{studentId}` - Get requests by student
- `GET /api/v1/absent-requests/class-id/{classId}` - Get requests by class
- `GET /api/v1/absent-requests/pending` - Get all pending requests
- `PATCH /api/v1/absent-requests/{id}/status` - Update request status
- `DELETE /api/v1/absent-requests/{id}` - Delete absent request

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
  "created_at": "2024-01-01T00:00:00Z",
  "updated_at": "2024-01-01T00:00:00Z"
}
```

**Attendance Status Options:**
- `present`: Student was present
- `absent`: Student was absent
- `late`: Student was late

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

#### Create a Teacher
```bash
curl -X POST http://localhost:8080/api/v1/teachers \
  -H "Content-Type: application/json" \
  -d '{
    "teacher_id": "TCH001",
    "first_name": "John",
    "last_name": "Doe",
    "email": "john.doe@school.com",
    "phone": "+1234567890",
    "password": "securepassword123"
  }'
```

#### Create a Class
```bash
curl -X POST http://localhost:8080/api/v1/classes \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Grade 10A",
    "homeroom_teacher": "TCH001",
    "description": "Advanced mathematics class"
  }'
```

#### Create a Student
```bash
curl -X POST http://localhost:8080/api/v1/students \
  -H "Content-Type: application/json" \
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

#### Record Attendance
```bash
curl -X POST http://localhost:8080/api/v1/attendances \
  -H "Content-Type: application/json" \
  -d '{
    "student_id": "STU001",
    "class_id": 1,
    "date": "2024-01-15T09:00:00Z",
    "status": "present",
    "description": "Student was on time"
  }'
```

#### Create Absence Request
```bash
curl -X POST http://localhost:8080/api/v1/absent-requests \
  -H "Content-Type: application/json" \
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

## Database Schema

The application uses PostgreSQL with the following main tables:

- **teachers**: Teacher information and credentials
- **classes**: Class information with homeroom teacher references
- **students**: Student information and class assignments
- **attendances**: Daily attendance records
- **absent_requests**: Student absence requests with approval workflow

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
  "translate.key": "success.photo_uploaded",
  "message": "Photo uploaded successfully",
  "path": "https://bucket-name.s3.region.amazonaws.com/photos/teachers/1/teacher_1_1704067200.jpg"
}
```

### Get Photo Response Format
```json
{
  "translate.key": "success.photo_url_retrieved",
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

## Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Support

For support and questions:
- Create an issue in the GitHub repository
- Contact: support@studentattendance.com

## Roadmap

- [ ] Authentication and authorization (JWT)
- [ ] Role-based access control
- [ ] Email notifications for absence requests
- [ ] Attendance reports and analytics
- [ ] Mobile app support
- [ ] Integration with external calendar systems
- [ ] Bulk operations for attendance
- [ ] Parent/guardian access portal