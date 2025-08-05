# Database Migrations

This directory contains the database migration system for the Student Attendance API.

## Migration Tool

The migration tool (`migration.go`) provides a simple interface to manage database schema changes using the [golang-migrate](https://github.com/golang-migrate/migrate) library.

### Usage

```bash
# Apply all pending migrations (up)
go run db/migration.go up

# Rollback all migrations (down)
go run db/migration.go down
```

### How it Works

1. The tool reads database configuration from environment variables (same as main application)
2. Connects to PostgreSQL using the connection string format
3. Looks for migration files in the `db/migrations/` directory
4. Applies migrations in sequential order based on version numbers

### Migration Files

Migration files should be placed in the `db/migrations/` directory with the following naming convention:

```
000001_initial_schema.up.sql
000001_initial_schema.down.sql
000002_add_users_table.up.sql
000002_add_users_table.down.sql
```

- **Version number**: 6-digit sequential number (000001, 000002, etc.)
- **Description**: Brief description of the migration
- **Direction**: `.up.sql` for forward migrations, `.down.sql` for rollbacks

### Environment Variables

The migration tool uses the same database configuration as the main application:

- `DB_HOST`: Database host
- `DB_PORT`: Database port
- `DB_NAME`: Database name
- `DB_USER`: Database username
- `DB_PASSWORD`: Database password
- `DB_SSL_MODE`: SSL mode (disable, require, etc.)

### Example Migration Files

**000001_create_students_table.up.sql:**
```sql
CREATE TABLE students (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

**000001_create_students_table.down.sql:**
```sql
DROP TABLE IF EXISTS students;
```

### Error Handling

- If no migrations are pending, the tool will complete successfully
- Migration errors will be logged and the process will exit
- The tool validates the migration direction parameter (up/down)

### Best Practices

1. Always create both up and down migration files
2. Test migrations on a copy of production data before applying
3. Use descriptive names for migration files
4. Keep migrations small and focused on single changes
5. Never modify existing migration files once they've been applied to production