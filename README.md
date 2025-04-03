# SkyFox Backend

A Go-based backend service for the SkyFox movie booking system with Supabase PostgreSQL integration.

## Overview

SkyFox Backend is a modern, well-structured API service that provides authentication, customer management, and security question features for the SkyFox movie booking application. Built with Go and the Gin framework, it implements clean architecture principles with proper separation of controllers, services, and repositories.

## Features

- JWT-based authentication with role-based authorization
- Customer signup with comprehensive validation
- Security question system for account recovery
- Standardized error responses
- PostgreSQL database integration via Supabase

## Project Structure

The project follows a clean architecture approach:
- Controllers: Handle HTTP requests and responses
- Services: Implement business logic
- Repositories: Manage data access
- Models: Define data structures
- DTOs: Handle data transfer objects
- Middleware: Process requests (validation, authentication)

## API Documentation

For detailed API documentation, please see the [API Documentation](./docs/README.md).

## Prerequisites

- Go 1.20+
- Supabase account with a project
- Git

## Environment Setup

1. Clone the repository:
```bash
git clone https://github.com/iamsuteerth/skyfox-backend.git
cd skyfox-backend
```

2. Create a `.env` file in the root directory with the following variables:
```
# Database Configuration
DATABASE_URL=postgresql://postgres:[YOUR-PASSWORD]@db.[YOUR-PROJECT-REF].supabase.co:5432/postgres

# Auth Configuration
JWT_SECRET_KEY=your_secure_jwt_secret

# Application Configuration
PORT=8080
APP_ENV=development  # Options: development, production
LOG_LEVEL=info       # Options: debug, info, warn, error
```

3. Install dependencies:
```bash
go mod download
```

## Database Migration

The migration files are stored in the `supabase/migration` directory:
- `000001_initial_schema.up.sql` - Creates the initial schema
- `000001_initial_schema.down.sql` - Drops the initial schema
- `000002_seat_types.up.sql` - Adds seat types
- `000002_seat_types.down.sql` - Removes seat types
- `000003_add_security_questions.up.sql` - Adds security questions system
- `000003_add_security_questions.down.sql` - Removes security questions system

To apply these migrations to your Supabase project, use the Supabase SQL Editor to execute the SQL scripts.

## Running the Application

To run the application:

```bash
cd server
go run main.go
```

The server will start on the port specified in your `.env` file (default 8080).

## Logging

The application uses structured logging with zerolog. Log levels can be configured in the `.env` file:
- `debug` - Detailed information for debugging
- `info` - General information about application flow
- `warn` - Warning events that might need attention
- `error` - Error events that might still allow the application to continue running

## Database Seeding

The application automatically seeds the database with initial data on startup:

- Admin users:
  - Username: `seed-user-1`, Password: `foobar`, Role: `admin`
  - Username: `seed-user-2`, Password: `foobar`, Role: `admin`
- Staff user:
  - Username: `staff-1`, Password: `foobar`, Role: `staff`

The seeding process creates records in both the user table and staff table, ensuring proper relationships.

## Authentication

The application uses JWT (JSON Web Token) for authentication. Tokens are valid for 24 hours and include the user's role for authorization purposes.

## Error Handling

The application uses standardized error responses:

```json
{
  "status": "ERROR",
  "code": "ERROR_CODE",
  "message": "Error message",
  "request_id": "unique-request-id",
  "errors": [
    {
      "field": "field_name",
      "message": "Error message for this field"
    }
  ]
}
```

- `status`: Always "ERROR" for error responses
- `code`: Machine-readable error code (e.g., "VALIDATION_ERROR", "INVALID_CREDENTIALS")
- `message`: Human-readable error message
- `request_id`: Unique identifier for the request
- `errors`: Array of field-specific validation errors (only present for validation failures)

## Database Schema

The database includes the following tables:

1. **usertable** - User authentication and roles
   - Contains username, password (hashed), and role

2. **password_history** - Password management
   - Tracks previous passwords for security measures

3. **stafftable** - Staff information
   - Links staff members to user accounts

4. **customertable** - Customer information
   - Contains security_question_id and security_answer_hash for account security

5. **security_questions** - Predefined security questions
   - Used for account recovery and additional security

6. **admin_booked_customer** - Customers booked by admins

7. **seat** - Theater seats (Standard and Deluxe types)

8. **slot** - Movie time slots

9. **show** - Movie screenings

10. **booking** - Ticket reservations

11. **booking_seat_mapping** - Mapping between bookings and seats

## Security Features

### Password Management
- Passwords are securely hashed before storage
- Password history is maintained to prevent reuse

### Security Questions
- Used as an additional layer of security for account recovery
- Security answers are hashed before storage
- Required during signup and used for password reset

## License

See the [LICENSE](LICENSE) file for details.