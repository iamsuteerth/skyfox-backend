# SkyFox Backend

A Go-based backend service for the SkyFox movie booking system with Supabase PostgreSQL integration.

## Project Setup

This document outlines the setup process, migration, and development workflow for the SkyFox backend project.

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

## Project Structure

The project follows a clean architecture approach with the following structure:

- `pkg/` - Core application code
  - `config/` - Configuration management
  - `controllers/` - HTTP request handlers
  - `database/seed/` - Database seeding functionality
  - `dto/` - Data transfer objects (for future API implementations)
  - `models/` - Domain models
  - `repositories/` - Database access layer
  - `services/` - Business logic
  - `utils/` - Utility functions and error handling
- `server/` - Application entry point
- `supabase/migration/` - Database migration scripts
- `scripts/` - Utility scripts for development

## Database Migration

The migration files are stored in the `supabase/migration` directory:
- `000001_initial_schema.up.sql` - Creates the initial schema
- `000001_initial_schema.down.sql` - Drops the initial schema
- `000002_seat_types.up.sql` - Adds seat types
- `000002_seat_types.down.sql` - Removes seat types

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

## API Endpoints

### Authentication

#### Login
- **URL**: `/login`
- **Method**: `POST`
- **Authentication**: None
- **Request Body**:
  ```json
  {
    "username": "string",
    "password": "string"
  }
  ```
- **Success Response (200 OK)**:
  ```json
  {
    "message": "Login successful",
    "request_id": "unique-request-id",
    "status": "SUCCESS",
    "data": {
      "user": {
        "username": "string",
        "role": "string"
      },
      "token": "jwt-token"
    }
  }
  ```
- **Error Response (401 Unauthorized)**:
  ```json
  {
    "error": "Invalid username or password",
    "request_id": "unique-request-id"
  }
  ```

## Error Handling

The application uses standardized error responses:

For general errors:
```json
{
  "error": "Error message",
  "request_id": "unique-request-id"
}
```

For validation errors:
```json
{
  "errors": [
    {
      "field": "field_name",
      "message": "Error message for this field"
    }
  ],
  "request_id": "unique-request-id",
  "status": "REJECT"
}
```

## Database Schema

The database includes the following tables:

1. **usertable** - User authentication and roles
   - Contains username, password (hashed), and role

2. **password_history** - Password management
   - Tracks previous passwords for security measures

3. **stafftable** - Staff information
   - Links staff members to user accounts

4. **customertable** - Customer information

5. **admin_booked_customer** - Customers booked by admins

6. **seat** - Theater seats (Standard and Deluxe types)

7. **slot** - Movie time slots

8. **show** - Movie screenings

9. **booking** - Ticket reservations

10. **booking_seat_mapping** - Mapping between bookings and seats

## License

See the [LICENSE](LICENSE) file for details.

