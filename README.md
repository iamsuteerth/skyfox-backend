# SkyFox Backend

A Go-based backend service for the SkyFox movie booking system with Supabase PostgreSQL integration.

## Overview

SkyFox Backend is a modern, well-structured API service that provides authentication, customer management, and a security question-based password recovery system for the SkyFox movie booking application. Built with Go and the Gin framework, it implements clean architecture principles with proper separation of controllers, services, and repositories.

![Architecture Diagram](./skyfox-backend-arch.png)

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
- Secure profile image management with S3 and presigned URLs

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
- AWS S3 bucket for profile image storage

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

# AWS S3 Configuration
AWS_ACCESS_KEY_ID=your-access-key
AWS_SECRET_ACCESS_KEY=your-secret-key
AWS_REGION=your-region
S3_BUCKET=your-project-profile-images
```

3. Install dependencies:
```bash
go mod download
```

## AWS Elastic Beanstalk Deployment

The SkyFox Backend is deployed on AWS Elastic Beanstalk, providing a scalable and managed environment for the application.

### Deployment Architecture

Our deployment uses a containerized approach with Docker, running on ARM64 architecture for optimal cost-efficiency and performance:

- **Platform**: Docker running on 64bit Amazon Linux 2023
- **Instance Type**: ARM64-based t4g instances (t4g.micro, t4g.small)
- **Environment Type**: Single instance (suitable for development and initial production)

### Deployment Process

1. **Multi-Stage Docker Build**:
   - First stage: Builds the Go application with proper optimization
   - Second stage: Creates a lightweight image with Nginx and the compiled binary

2. **Image Storage**:
   - The Docker image is stored in Amazon ECR (Elastic Container Registry)
   - Tagged and versioned for consistent deployments

3. **Environment Configuration**:
   - Environment variables are securely managed through Elastic Beanstalk
   - IAM roles provide secure access to AWS services like S3

4. **Monitoring and Logging**:
   - Basic CloudWatch monitoring is enabled
   - Application logs are retained for 7 days
   - Supervisor manages process logs within the container

## Nginx Reverse Proxy with API Key Security

The SkyFox Backend implements a robust security layer using Nginx as a reverse proxy with API key validation.

### Key Features

- **API Gateway Pattern**: Nginx functions as an API Gateway, providing a security layer in front of the Go application
- **API Key Validation**: All requests (except health checks) require a valid API key in the `X-Api-Key` header
- **Rate Limiting**: Protects against abuse with configurable request rate limiting
- **Process Management**: Supervisor ensures both Nginx and the Go application run reliably
- **Security Headers**: Implements best practices for HTTP security headers
- **Performance Optimization**: Includes Gzip compression, connection optimizations, and efficient file serving

### How It Works

1. **Request Flow**:
   - Incoming requests arrive at Nginx on port 80
   - Nginx validates the API key in the `X-Api-Key` header
   - Valid requests are proxied to the Go application on port 8080
   - The API key is stripped before forwarding to the backend for security

2. **Configuration Management**:
   - Nginx configuration uses environment variable substitution
   - The API Gateway key is securely stored as an environment variable
   - Configuration is generated at container startup

3. **Health Checks**:
   - A dedicated `/health` endpoint bypasses API key requirements
   - Returns a simple JSON response: `{"status":"healthy"}`
   - Enables monitoring systems to verify application health

### Security Benefits

- **Attack Surface Reduction**: The Go application is never directly exposed to the internet
- **Secret Management**: API keys are never hard-coded in the application
- **Request Filtering**: Unauthorized requests are blocked before reaching the application code
- **Header Sanitization**: Security-sensitive headers are managed by Nginx
- **Connection Control**: Timeouts and buffer sizes are optimized to prevent abuse
- **Standardized Errors**: Unauthorized requests receive consistent error responses

### Testing with Docker

To test this setup locally:

1. **Build the Docker image**:
   ```bash
   docker build -t skyfox-backend:local .
   ```

2. **Run the container**:
   ```bash
   docker run -p 80:80 \
     -e PORT=8080 \
     -e API_GATEWAY_KEY=your_api_key_here \
     -e OTHER_ENV_VARIABLES=values \
     skyfox-backend:local
   ```

3. **Test the API key validation**:
   ```bash
   # This should return 403
   curl http://localhost/

   # This should succeed
   curl -H "X-Api-Key: your_api_key_here" http://localhost/

   # Health check should always work
   curl http://localhost/health
   ```

This deployment architecture provides a secure, scalable, and maintainable infrastructure for the SkyFox Backend service, leveraging AWS managed services while implementing industry best practices for API security.

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

### AWS S3 Integration for Profile Images
The application implements a sophisticated profile image management system using **AWS S3**:

1. **Secure Upload Process**:
   - Frontend scales the image to 64x64 pixels and converts it to base64
   - A SHA-256 hash is calculated from the original image bytes
   - Both the base64 string and hash are sent to the backend

2. **Verification and Processing**:
   - Backend decodes the base64 string back to a byte array
   - Recalculates the SHA-256 hash to verify integrity
   - Rescales the image to ensure consistent 64x64 dimensions
   - Uploads the processed image to AWS S3 with private ACL

3. **Efficient Storage**:
   - Only the S3 URL is stored in the database, not the image data
   - This keeps the database lightweight and optimized

4. **Secure Retrieval**:
   - When an image is requested, a presigned URL is generated with a 24-hour expiration
   - This approach prevents direct access to S3 resources

5. **Performance Benefits**:
   - Users directly interact with AWS S3 when fetching images
   - This offloads bandwidth and processing from the application server
   - Results in faster image loading and better scalability

6. **Cleanup Handling**:
   - When profile images are updated, old S3 objects are automatically deleted
   - Failed operations include proper rollback to maintain data consistency

This implementation provides an optimal balance of security, performance, and user experience by leveraging AWS S3's capabilities while maintaining proper access controls.

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

![Supabase Database Schema](./database_schema.png)

The database includes the following tables:

1. **usertable** - User authentication and roles
   - Contains username, password (hashed), and role

2. **password_history** - Password management
   - Tracks previous passwords for security measures

3. **stafftable** - Staff information
   - Links staff members to user accounts

4. **customertable** - Customer information
   - Contains `security_question_id` and `security_answer_hash` for account security
   - Stores S3 URLs for profile images

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