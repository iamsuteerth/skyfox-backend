# SkyFox Backend

A Go-based backend service for the SkyFox movie booking system with Supabase PostgreSQL integration.

## Project Setup

This document outlines the setup process, migration, and development workflow for the SkyFox backend project.

## Prerequisites

- Go 1.23+
- PostgreSQL client tools (`psql` command)
- Supabase account with a project

## Environment Setup

1. Clone the repository:
```bash
git clone https://github.com/yourusername/skyfox-backend.git
cd skyfox-backend
```

2. Set up environment variables:
```bash
export DB_HOST=your-supabase-uri
export DB_PORT=5432
export DB_NAME=postgres
export POSTGRES_USERNAME=postgres
export DB_PASSWORD=your-supabase-db-password
```

3. Install dependencies:
```bash
go mod download
```

## Database Migration

The project uses a custom Go-based migration tool to manage database schema changes with Supabase.

### Migration Files

Migration files are stored in the `supabase/migration` directory using the format:
- `000001_initial_schema.up.sql` - Creates the initial schema
- `000001_initial_schema.down.sql` - Drops the initial schema
- `000002_seat_types.up.sql` - Adds seat types
- `000002_seat_types.down.sql` - Removes seat types

### Running Migrations

To apply migrations:
```bash
make run-migrations
```

To roll back migrations:
```bash
make rollback-migrations
```

If you encounter a "dirty database" error, you can force a specific migration version:
```bash
make force-migration
```

## Seeding Data

The project includes a seed script to populate the database with test data.

To seed the database:
```bash
make seedData
```

This will:
- Clear existing data and reset sequences
- Add 4 time slots (Morning, Afternoon, Evening, Night)
- Add 100 seats (A1-J10), with rows A-E as Standard type and F-J as Deluxe type
- Add movie shows for the next 21 days across all time slots

## Development

### Available Make Commands

- `make run-migrations` - Apply database migrations
- `make rollback-migrations` - Revert the most recent migration
- `make force-migration` - Force a dirty migration
- `make deploy` - Build and run the Go server locally
- `make seedData` - Populate the database with test data

## Database Schema

The database includes the following tables:

1. **usertable** - User authentication and roles
2. **password_history** - Password management
3. **stafftable** - Staff member information
4. **customertable** - Customer information
5. **admin_booked_customer** - Customers booked by admins
6. **seat** - Theater seats (Standard and Deluxe types)
7. **slot** - Movie time slots
8. **show** - Movie screenings
9. **booking** - Ticket reservations
10. **booking_seat_mapping** - Mapping between bookings and seats

## External Services

The application integrates with:
1. **Movie Service** - External API for movie information
2. **Payment Service** - External API for payment processing

Both services require API keys stored in environment variables.

## Local Development

For local development:
- The Go server runs directly on the development machine
- Movie and payment services run in containers
- Database is hosted on Supabase

## Deployment

The project is configured for local development only. Future updates will include production deployment configurations.

## License

See the [LICENSE](LICENSE) file for details.
