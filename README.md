# SkyFox Backend

A Go-based backend service for the SkyFox movie booking system with Supabase PostgreSQL integration.

## Overview

- SkyFox Backend is a modern, well-structured API service that provides authentication, RBAC, customer management, and a security question-based password recovery system for the SkyFox movie booking application. 
- Built with Go and the Gin framework, it implements clean architecture principles with proper separation of controllers, services, and repositories.

> **Note:** This documentation is specific to the backend. For DevOps, please refer to [skyfox-devops](https://github.com/iamsuteerth/skyfox-devops).

### Production AWS Architecture
![Architecture Diagram - Production](./AWS_Skyfox.png)

### Legacy NGINX Architecture
![Architecture Diagram - POC](./skyfox-backend-arch.png)

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
- **Sophisticated Booking System**: Two-phase booking process with temporary seat reservation, automated expiration, and integrated payment processing.
- **Efficient Concurrency Handling**: Each booking gets its own dedicated monitor, allowing thousands of concurrent reservations with precise timing control.
- **Digital Wallet System**: Integrated customer wallet for funds management with secure transaction tracking and support for partial wallet payments.
- **OLTP Support**: Decimal package implementation for precise financial calculations and transaction processing.

## Go Backend Security Defenses

- **API Key Gatekeeping:** Every business endpoint is protected by a mandatory API key, validated on every request.
- **No File-Based Secrets:** The application never reads secrets from the filesystem or in-repo config files. All sensitive data comes from environment variables only.
- **Read-Only Secrets at Runtime:** Environment variables holding secrets are injected at runtime from AWS SSM and are never logged, echoed, or returned in responses.
- **Scanner-Resistant Routing:** Common attack paths (`/secrets/aws/*`, `/env/*`, `/config/*`, etc.) are always met with 403 Forbidden—never a 200, never a file disclosure.
- **Detailed Audit Logging:** Every denied access attempt (403) is logged with source IP and request details, supporting proactive threat monitoring.
- **No Leaky Endpoints:** Even advanced bots probing for AWS credentials, ENV, or config files cannot extract real secrets.

> **Summary:**  
> The SkyFox Go backend is designed with modern SaaS threat models in mind.  
> Production deployments are immune to standard cloud credential harvesting attacks—even if bots hit you 1,000,000 times, they only get a 403.

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
- `000006_created_at_column.up.sql` - Adds a column created_at in usertable to have timestamps of account creation handled automatically by the database
- `000006_created_at_column.down.sql` - Drops the created_at column in the usertable
- `000007_booking_table_modifications.up.sql` - Adds customer_username to booking table
- `000007_booking_table_modifications.down.sql` - Removes customer_username from booking table
- `000008_payment_mode_enum.up.sql` - Creates payment_mode enum type
- `000008_payment_mode_enum.down.sql` - Removes payment_mode enum type
- `000009_booking_timestamp_update.up.sql` - Updates booking_time to use timezone
- `000009_booking_timestamp_update.down.sql` - Reverts booking_time timezone change
- `000010_relationship_btw_abc_booking.up.sql` - Sets up cascading delete from booking to admin_booked_customer
- `000010_relationship_btw_abc_booking.down.sql` - Reverts the cascading delete relationship
- `000011_payment_txd_table.up.sql` - Creates payment transaction table for online payments
- `000011_payment_txd_table.down.sql` - Removes payment transaction table
- `000012_pending_booking_tracker.up.sql` - Creates table to track pending bookings and their expiration times
- `000012_pending_booking_tracker.down.sql` - Removes pending booking tracker
- `000013_add_foreign_key_indices.up.sql` - Adds performance indices on foreign keys
- `000013_add_foreign_key_indices.down.sql` - Removes foreign key indices
- `000014_circular_dependency_fix.up.sql` - Fixes circular dependency between booking and admin_booked_customer
- `000014_circular_dependency_fix.down.sql` - Reverts the circular dependency fix
- `000015_performance_idices.up.sql` - Add performance indices on revenue related fields
- `000015_performance_idices.dowb.sql` - Removes the indices on revenue related fields
- `000016_wallet_payment_type.up.sql` - Updates payment_mode enum to include Wallet as a payment type
- `000016_wallet_payment_type.down.sql` - Postgres doesn't allow us to remove element from an enum.
- `000017_wallet_db_tables.up.sql` - Creates customer_wallet table for storing wallet balances
- `000017_wallet_db_tables.down.sql` - Drops the customer_wallet table
- `000018_wallet_transaction_type_enum.up.sql` - Creates wallet_transaction_type enum for ADD/DEDUCT operations
- `000018_wallet_transaction_type_enum.down.sql` - Drops the wallet_transaction_type enum
- `000019_wallet_transaction_table.up.sql` - Creates wallet_transaction table for tracking wallet operations
- `000019_wallet_transaction_table.down.sql` - Drops the wallet_transaction table
- `000020_add_cascade_delete_behavior.up.sql` - Adds CASCADE delete behavior to customer foreign key constraints for development flexibility
- `000020_add_cascade_delete_behavior.down.sql` - Removes CASCADE delete behavior, restoring original constraint behavior
- `000021_optimize_database_indices.up.sql` - Optimizes database performance by adding composite and partial indexes for common query patterns while removing redundant indexes
- `000021_optimize_database_indices.down.sql` - Reverts index optimizations and restores original index structure

To apply these migrations to your Supabase project, use the Supabase SQL Editor or a migration tool.

### Migration Analysis Tool

The SkyFox Backend includes a custom migration analysis tool that consolidates all database migration files into comprehensive documentation for review and analysis.

#### **Purpose**

The Migration Analyzer Script addresses the common challenge of understanding database schema evolution across multiple migration files. As projects grow and accumulate dozens of migrations, it becomes difficult to:
- Visualize the complete database schema evolution
- Identify migration dependencies and potential conflicts
- Create comprehensive setup scripts for fresh environments
- Plan effective rollback strategies for production deployments

#### **Features**

- **Consolidated View**: Merges all `.up.sql` and `.down.sql` files into single, readable documents
- **Proper Ordering**: Processes UP migrations in ascending order and DOWN migrations in descending order for correct application sequence
- **Comprehensive Headers**: Each migration is clearly labeled with separators for easy navigation
- **Summary Statistics**: Provides counts of total migrations and missing files
- **Error Detection**: Identifies missing migration files and reports warnings
- **Fresh Setup Support**: The consolidated UP file can serve as a complete database setup script

#### **Usage**

```bash
# Make the script executable
chmod +x analyze_migrations.sh

# Run from your project root (adjust MIGRATIONS_DIR path as needed)
./analyze_migrations.sh
```

#### **Output Files**

The tool generates two primary files in the `./migration_analysis/` directory:

1. **`mega_up_migration.sql`**: Contains all UP migrations in chronological order
   - Use for fresh database setup
   - Review complete schema evolution
   - Identify cumulative changes

2. **`mega_down_migration.sql`**: Contains all DOWN migrations in reverse order
   - Use for complete database teardown
   - Plan rollback strategies
   - Test migration reversibility

#### **Development Benefits**

- **Documentation**: Generate comprehensive migration documentation
- **Debugging**: Quickly identify conflicting or problematic migrations
- **Onboarding**: Help new developers understand the database structure evolution
- **Production Planning**: Create reliable deployment and rollback procedures

This tool is particularly valuable for the SkyFox project given its 20+ migrations spanning user management, booking systems, payment processing, and wallet functionality.

## Running the Application

To run the application in debug mode:

```bash
cd server
go run server/main.go
```

To run the application in release mode:
```bash
make deploy
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

## Wallet System

The SkyFox platform includes a comprehensive digital wallet system that enables customers to:

- Add funds to their wallet using card payments
- Use wallet balance for ticket bookings
- View transaction history for all wallet operations
- Combine wallet and card payments for a seamless experience

### Key Features

1. **Precise Financial Calculations**: Utilizes the decimal package to ensure accurate monetary operations without floating-point errors, critical for OLTP systems.

2. **Transaction Tracking**: All wallet operations (ADD/DEDUCT) are tracked with timestamps and transaction IDs for complete auditability.

3. **Payment Integration**: Seamlessly integrates with the existing payment system to support:
   - Full wallet payments (when wallet has sufficient balance)
   - Partial wallet payments (automatically topping up wallet with card payment for the shortfall)
   - Traditional card payments

4. **Automatic Wallet Creation**: Each customer account is provisioned with a wallet upon registration.

5. **Security**: All financial transactions are protected with proper validations and database constraints.

## OLTP Support

The system implements robust OLTP (Online Transaction Processing) capabilities through:

1. **Decimal Package Integration**: Using the decimal package for all monetary values ensures exact arithmetic without the precision issues of floating-point calculations.

2. **Transaction Consistency**: All financial operations maintain ACID properties with proper rollback mechanisms.

3. **Concurrency Control**: The wallet system handles multiple simultaneous transactions without race conditions.

4. **Audit Trail**: Complete transaction history is maintained for reporting and reconciliation.

## Authentication

The application uses JWT (JSON Web Token) for authentication. Tokens are valid for 24 hours and include the user's role for authorization purposes.

## Role-Based Access

The application implements role-based access control:

1. **Customer Role**:
   - Can view shows only for the current date plus 6 days
   - Has access to personal profile and booking history
   - Can manage their wallet and view transaction history

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

## Booking Management System

The SkyFox platform provides a comprehensive booking management system with separate workflows for administrators and customers:

### Admin Booking Workflow
Administrators can create bookings directly through a simplified one-step process:
- Create bookings on behalf of walk-in customers
- Instantly confirm bookings with cash payments
- Manage seating arrangements through an intuitive interface
- No payment processing step required (all admin bookings use cash payment)

### Customer Booking Workflow
Customers follow a two-step booking process:
1. Initialize booking with seat selection (seats reserved for 5 minutes)
2. Complete payment within the reservation window to confirm booking using:
   - Credit/debit card
   - Wallet funds
   - Combined wallet and card payment
3. Automatic seat release if payment is not completed in time

This dual approach accommodates both in-person and online ticket purchases while maintaining consistent data structures.

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

13. **payment_transaction**: Records payment details for online bookings

14. **pending_booking_tracker**: Manages temporary seat reservations

15. **customer_wallet** - Stores wallet balance for each customer

16. **wallet_transaction** - Records all wallet operations (ADD/DEDUCT)

## License

See the [LICENSE](LICENSE) file for details.