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
    "status": "ERROR",
    "code": "INVALID_CREDENTIALS",
    "message": "Invalid username or password",
    "request_id": "unique-request-id"
  }
  ```

#### Get Security Questions
- **URL**: `/api/security-questions`
- **Method**: `GET`
- **Authentication**: None
- **Success Response (200 OK)**:
  ```json
  {
    "message": "Security questions retrieved successfully",
    "request_id": "unique-request-id",
    "status": "SUCCESS",
    "data": [
      {
        "id": 1,
        "question": "What was the name of your first pet?"
      },
      {
        "id": 2,
        "question": "What was your childhood nickname?"
      }
      // Additional security questions...
    ]
  }
  ```

#### Customer Signup
- **URL**: `/customer/signup`
- **Method**: `POST`
- **Authentication**: None
- **Request Body**:
  ```json
  {
    "name": "string",
    "username": "string",
    "password": "string",
    "number": "string",
    "email": "string",
    "profile_img": null,
    "security_question_id": 1,
    "security_answer": "string"
  }
  ```
- **Success Response (201 Created)**:
  ```json
  {
    "message": "User registered successfully",
    "request_id": "unique-request-id",
    "status": "SUCCESS",
    "data": {
      "username": "string",
      "name": "string"
    }
  }
  ```
- **Validation Error Response (400 Bad Request)**:
  ```json
  {
    "status": "ERROR",
    "code": "VALIDATION_ERROR",
    "message": "Validation failed",
    "request_id": "unique-request-id",
    "errors": [
      {
        "field": "Name",
        "message": "Name must be 3-70 characters, max 4 words, letters only, no consecutive spaces"
      },
      {
        "field": "SecurityAnswer",
        "message": "Security answer must be at least 3 characters long"
      }
      // Other validation errors...
    ]
  }
  ```
- **Security Question Error Response (400 Bad Request)**:
  ```json
  {
    "status": "ERROR",
    "code": "INVALID_SECURITY_QUESTION",
    "message": "The selected security question does not exist",
    "request_id": "unique-request-id"
  }
  ```
- The validation system enforces strict rules for different field types:
  - **Name**: Must be 3-70 characters, max 4 words, letters only, no consecutive spaces
  - **Username**: Must be 3-30 characters, lowercase, no spaces, cannot start with a number, no consecutive special characters
  - **Password**: Must be at least 8 characters with at least one uppercase letter and one special character
  - **Phone Number**: Must be exactly 10 digits
  - **Email**: Must be in valid email format
  - **Security Answer**: Must be at least 3 characters long

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
