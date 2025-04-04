# SkyFox Backend

A Go-based backend service for the SkyFox movie booking system with Supabase PostgreSQL integration.

## Overview

SkyFox Backend is a modern, well-structured API service that provides authentication, customer management, and a security question-based password recovery system for the SkyFox movie booking application. Built with Go and the Gin framework, it implements clean architecture principles with proper separation of controllers, services, and repositories.

## Features

- JWT-based authentication with role-based authorization
- Customer signup with comprehensive validation
- Security question system for account recovery and password reset
- Token-based password reset functionality with expiration and uniqueness
- Standardized error responses
- PostgreSQL database integration via Supabase
- Movie data integration with an external movie service
- Show scheduling and management system
- Role-based content filtering (different views for customers vs. admins)
- Available slot management for preventing double-booking

## Project Structure

The project follows a clean architecture approach:
- **Controllers**: Handle HTTP requests and responses
- **Services**: Implement business logic
- **Repositories**: Manage data access
- **Models**: Define data structures
- **DTOs**: Manage data transfer objects for requests and responses
- **Middleware**: Process requests (e.g., CORS, validation, authentication)

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

# Movie Service Configuration
MOVIE_SERVICE_URL=http://localhost:4567
MOVIE_SERVICE_API_KEY=your_movie_service_api_key
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
- `000004_security_question_temp_tkn.up.sql` - Adds the password reset tokens table
- `000004_security_question_temp_tkn.down.sql` - Drops the password reset tokens table
- `000005_add_unique_constraint_email_token.up.sql` - Adds a unique constraint on `(email, token)` for password reset tokens
- `000005_add_unique_constraint_email_token.down.sql` - Removes the unique constraint on `(email, token)`

To apply these migrations to your Supabase project, use the Supabase SQL Editor or a migration tool.

## Running the Application

To run the application:

```bash
cd server
go run main.go
```

The server will start on the port specified in your `.env` file (default: 8080).

## Logging

The application uses structured logging with `zerolog`. Log levels can be configured in the `.env` file:
- `debug` - Detailed information for debugging
- `info` - General information about application flow
- `warn` - Warning events that might need attention
- `error` - Error events that might still allow the application to continue running

## Database Seeding

The application automatically seeds the database with initial data on startup:

- **Admin users**:
  - Username: `seed-user-1`, Password: `foobar`, Role: `admin`
  - Username: `seed-user-2`, Password: `foobar`, Role: `admin`
- **Staff user**:
  - Username: `staff-1`, Password: `foobar`, Role: `staff`

The seeding process creates records in both the `user` table and `staff` table, ensuring proper relationships.

## External Service Integration

### Movie Service
The application integrates with an external movie service to retrieve movie data:
- Requires `MOVIE_SERVICE_URL` and `MOVIE_SERVICE_API_KEY` environment variables
- Fetches movie details like title, runtime, plot, and poster images
- Caches movie data to minimize external API calls

## Authentication

The application uses JWT (JSON Web Token) for authentication. Tokens are valid for 24 hours and include the user's role for authorization purposes.

## Role-Based Access

The application implements role-based access control:

1. **Customer Role**:
   - Can view shows only for the current date plus 6 days
   - Has access to personal profile and booking history

2. **Staff Role**:
   - Has access to check-in functionality
   - Can download booking data as CSV

3. **Admin Role**:
   - Can view shows for any date
   - Can create and schedule new shows
   - Can view revenue data

## Security Question and Password Reset System

### Security Features:
- Security questions are used as an additional layer of security for account recovery.
- Security answers are hashed before being stored in the database for security purposes.
- Password reset is token-based with the following rules:
  - Tokens expire 5 minutes after creation.
  - Only one valid token is allowed per email at any time.
  - All previous tokens are deleted when generating a new token.

### Password Reset Tokens:
- Tokens are managed in the `password_reset_tokens` table.
- A unique constraint ensures no duplicate `(email, token)` pairs.

### Password Management:
- Passwords are securely hashed before storage, using `bcrypt`.
- Password history is maintained to prevent reuse of old passwords.

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
- `code`: Machine-readable error code (e.g., `VALIDATION_ERROR`, `INVALID_CREDENTIALS`)
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
   - Contains `security_question_id` and `security_answer_hash` for account security

5. **security_questions** - Predefined security questions
   - Used for account recovery and additional security

6. **password_reset_tokens** - Temporary tokens for password reset
   - Contains a unique token, expiration timestamp, and `used` flag for invalidation
   - Enforces a unique constraint on `(email, token)`

7. **admin_booked_customer** - Customers booked by admins

8. **seat** - Theater seats (Standard and Deluxe types)

9. **slot** - Movie time slots

10. **show** - Movie screenings

11. **booking** - Ticket reservations

12. **booking_seat_mapping** - Mapping between bookings and seats

## License

See the [LICENSE](LICENSE) file for details.
