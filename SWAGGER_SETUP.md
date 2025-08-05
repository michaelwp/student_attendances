# Swagger Integration Guide

This document explains how Swagger is integrated with the Student Attendance API.

## What's Included

### Dependencies Added
- `github.com/swaggo/swag` - Swagger documentation generator
- `github.com/swaggo/fiber-swagger` - Fiber middleware for Swagger UI
- `github.com/swaggo/files` - Static file serving for Swagger UI

### Files Structure
```
docs/
├── docs.go          # Generated Go documentation
├── swagger.json     # Generated JSON specification
└── swagger.yaml     # Generated YAML specification (original manually created)

cmd/student_attendance/main.go  # Main file with API info annotations
internal/api/router.go          # Router with Swagger UI endpoint
internal/api/handlers/          # Handlers with Swagger annotations
.swaggo.yml                     # Swagger configuration file
```

## How It Works

### 1. Swagger Annotations
The API documentation is generated from special comments in the code:

```go
// CreateTeacher godoc
// @Summary Create a new teacher
// @Description Create a new teacher in the system
// @Tags Teachers
// @Accept json
// @Produce json
// @Param teacher body models.Teacher true "Teacher data"
// @Success 201 {object} map[string]interface{} "Teacher created successfully"
// @Failure 400 {object} map[string]interface{} "Invalid request body"
// @Router /teachers [post]
func (h *teacherHandler) Create(c *fiber.Ctx) error {
    // Handler implementation
}
```

### 2. Documentation Generation
Run this command to generate documentation:
```bash
swag init -g cmd/student_attendance/main.go --output docs --parseDependency --parseInternal
```

### 3. Swagger UI Integration
The Swagger UI is served at `/swagger/*` endpoint:
```go
app.Get("/swagger/*", fiberSwagger.WrapHandler)
```

## Adding New Endpoints

When you add new API endpoints:

1. **Add Swagger annotations** to your handler functions
2. **Regenerate documentation**:
   ```bash
   make swagger
   # OR
   swag init -g cmd/student_attendance/main.go --output docs --parseDependency --parseInternal
   ```
3. **Test the endpoint** in Swagger UI at http://localhost:8080/swagger/

## Swagger Annotation Reference

### Common Annotations
- `@Summary` - Brief description
- `@Description` - Detailed description  
- `@Tags` - Group endpoints by functionality
- `@Accept` - Request content type (json, xml, etc.)
- `@Produce` - Response content type
- `@Param` - Parameter definition
- `@Success` - Success response
- `@Failure` - Error response
- `@Router` - Route path and HTTP method

### Parameter Types
- `path` - URL path parameter: `@Param id path int true "User ID"`
- `query` - Query parameter: `@Param limit query int false "Limit"`
- `body` - Request body: `@Param user body models.User true "User data"`
- `header` - Header parameter: `@Param Authorization header string true "Bearer token"`

### Response Examples
```go
// @Success 200 {object} models.User "Success response"
// @Success 200 {array} models.User "List of users"
// @Success 200 {object} map[string]interface{} "Generic response"
// @Failure 400 {object} map[string]interface{} "Bad request"
```

## Configuration

### Main API Info (main.go)
```go
// Student Attendance API
//
// A comprehensive API for managing student attendance
//
// Version: 1.0.0
// Host: localhost:8080
// BasePath: /api/v1
// Schemes: http, https
//
// swagger:meta
```

### Swaggo Config (.swaggo.yml)
```yaml
dir: ./
generalInfo:
  title: Student Attendance API
  version: "1.0.0"
  host: localhost:8080
  basePath: /api/v1
output: docs
mainAPIFile: cmd/student_attendance/main.go
```

## Accessing Documentation

### During Development
1. Start the server: `go run cmd/student_attendance/main.go`
2. Open browser: http://localhost:8080/swagger/index.html
3. Test API endpoints directly from the UI

### Static Files
- **Swagger YAML**: Copy `docs/swagger.yaml` to Swagger Editor
- **Postman Collection**: Import `docs/swagger.json` into Postman
- **Client Generation**: Use swagger-codegen with the spec files

## Best Practices

1. **Always add annotations** when creating new endpoints
2. **Use consistent tags** to group related endpoints
3. **Provide clear descriptions** for parameters and responses
4. **Test annotations** by checking the generated documentation
5. **Regenerate docs** after making changes
6. **Use proper HTTP status codes** in success/failure responses

## Troubleshooting

### Common Issues

1. **Documentation not updating**:
   - Regenerate: `make swagger`
   - Restart the server

2. **Swagger UI not loading**:
   - Check the import: `_ "github.com/michaelwp/student_attendance/docs"`
   - Verify the route: `app.Get("/swagger/*", fiberSwagger.WrapHandler)`

3. **Model not showing**:
   - Use `--parseDependency --parseInternal` flags
   - Check model struct tags

4. **Build errors**:
   - Run `go mod tidy`
   - Check imports in docs/docs.go

### Validation
- Visit `/swagger/index.html` to see if UI loads
- Check console for JavaScript errors
- Verify JSON validity at `/swagger/doc.json`