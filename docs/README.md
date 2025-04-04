# SkyFox Backend API Documentation

## Authentication

### Login
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

## Security Questions

### Get All Security Questions
- **URL**: `/security-questions`
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

### Get Security Question By Email
- **URL**: `/security-question/by-email`
- **Method**: `POST`
- **Authentication**: None
- **Request Body**:
  ```json
  {
    "email": "user@example.com"
  }
  ```
- **Success Response (200 OK)**:
  ```json
  {
    "message": "Security question retrieved successfully",
    "request_id": "unique-request-id",
    "status": "SUCCESS",
    "data": {
      "question_id": 1,
      "question": "What was the name of your first pet?",
      "email": "user@example.com"
    }
  }
  ```
- **Error Response (404 Not Found)**:
  ```json
  {
    "status": "ERROR",
    "code": "USER_NOT_FOUND",
    "message": "No user found with the provided email",
    "request_id": "unique-request-id"
  }
  ```
### Verify Security Answer
- **URL**: `/verify-security-answer`
- **Method**: `POST`
- **Authentication**: None
- **Request Body**:
  ```json
  {
    "email": "user@example.com",
    "security_answer": "string"
  }
  ```
- **Success Response (200 OK)**:
  ```json
  {
    "message": "Security answer verified successfully",
    "request_id": "unique-request-id",
    "status": "SUCCESS",
    "data": {
      "reset_token": "f200d83b-ab92-41a7-ba59-929d1472e692",
      "expires_in_seconds": 300
    }
  }
  ```
- **Error Response (400 Bad Request)**:
  ```json
  {
    "status": "ERROR",
    "code": "INVALID_ANSWER",
    "message": "The security answer provided is incorrect",
    "request_id": "unique-request-id"
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
        "field": "SecurityAnswer",
        "message": "Security answer must be at least 3 characters long"
      }
    ]
  }
  ```

## User Management

### Customer Signup
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
      // Other Validation Errors
    ]
  }
  ```

### Reset Password with Token
- **URL**: `/forgot-password`
- **Method**: `POST`
- **Authentication**: None
- **Request Body**:
  ```json
  {
    "email": "user@example.com",
    "reset_token": "f200d83b-ab92-41a7-ba59-929d1472e692",
    "new_password": "SecurePass@123"
  }
  ```
- **Success Response (200 OK)**:
  ```json
  {
    "message": "Password has been reset successfully",
    "request_id": "unique-request-id",
    "status": "SUCCESS"
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
        "field": "Email",
        "message": "Invalid email format"
      },
      {
        "field": "ResetToken",
        "message": "This field is required"
      },
      {
        "field": "NewPassword",
        "message": "Password must be at least 8 characters with at least one uppercase letter and one special character"
      }
    ]
  }
  ```
- **Invalid Token Error Response (400 Bad Request)**:
  ```json
  {
    "status": "ERROR",
    "code": "INVALID_RESET_TOKEN",
    "message": "The reset token is invalid, expired, or has already been used",
    "request_id": "unique-request-id"
  }
  ```
- **Password Reuse Error Response (400 Bad Request)**:
  ```json
  {
    "status": "ERROR",
    "code": "PASSWORD_REUSE",
    "message": "New password cannot match any of your previous passwords",
    "request_id": "unique-request-id"
  }
  ```

## Show Management

### Get Shows
- **URL**: `/shows`
- **Method**: `GET`
- **Authentication**: Required
- **Query Parameters**:
  - `date`: Date in YYYY-MM-DD format (optional, defaults to current date)
- **Notes**: 
  - Admin and staff can view shows for any date
  - Customers can only view shows from today to the next 6 days
- **Success Response (200 OK)**:
  ```json
  {
    "message": "Shows retrieved successfully",
    "request_id": "unique-request-id",
    "status": "SUCCESS",
    "data": [
      {
        "movie": {
          "movieId": "tt1375666",
          "name": "Inception",
          "duration": "2h 28min",
          "plot": "A thief who steals corporate secrets through the use of dream-sharing technology...",
          "imdbRating": "8.8",
          "moviePoster": "https://example.com/inception_poster.jpg",
          "genre": "Action, Adventure, Sci-Fi"
        },
        "slot": {
          "id": 1,
          "name": "Morning",
          "startTime": "09:00:00.000000",
          "endTime": "12:00:00.000000"
        },
        "id": 1,
        "date": "2025-04-30T00:00:00Z",
        "cost": 250.50,
        "availableseats": 100
      }
    ]
  }
  ```
- **Date Range Error Response (400 Bad Request)**:
  ```json
  {
    "status": "ERROR",
    "code": "DATE_OUT_OF_RANGE",
    "message": "Customers can only view shows from today to the next 6 days",
    "request_id": "unique-request-id"
  }
  ```

### Get Movies
- **URL**: `/shows/movies`
- **Method**: `GET`
- **Authentication**: Required (Admin only)
- **Success Response (200 OK)**:
  ```json
  {
    "message": "Movies retrieved successfully",
    "request_id": "unique-request-id",
    "status": "SUCCESS",
    "data": [
      {
        "imdbid": "tt6644200",
        "title": "A Quiet Place",
        "runtime": "1h30m0s",
        "plot": "In a post-apocalyptic world, a family is forced to live in silence while hiding from monsters with ultra-sensitive hearing.",
        "imdbRating": "7.5",
        "poster": "https://example.com/quiet_place.jpg",
        "genre": "Drama, Horror, Sci-Fi"
      },
      {
        "imdbid": "tt1375666",
        "title": "Inception",
        "runtime": "2h 28min",
        "plot": "A thief who steals corporate secrets through the use of dream-sharing technology...",
        "imdbRating": "8.8",
        "poster": "https://example.com/inception_poster.jpg",
        "genre": "Action, Adventure, Sci-Fi"
      }
      // Additional movies...
    ]
  }
  ```
- **Error Response (401 Unauthorized)**:
  ```json
  {
    "status": "ERROR",
    "code": "UNAUTHORIZED",
    "message": "Missing Authorization header",
    "request_id": "unique-request-id"
  }
  ```
- **Error Response (403 Forbidden)**:
  ```json
  {
    "status": "ERROR",
    "code": "FORBIDDEN",
    "message": "Access denied. Admin role required",
    "request_id": "unique-request-id"
  }
  ```

### Create Show
- **URL**: `/shows`
- **Method**: `POST`
- **Authentication**: Required (Admin only)
- **Request Body**:
  ```json
  {
    "movieId": "tt1375666",
    "date": "2025-05-01",
    "slotId": 2,
    "cost": 250.50
  }
  ```
- **Notes**:
  - Date must be in YYYY-MM-DD format and not in the past
  - Cost must be greater than 0 and less than or equal to 3000
  - SlotId must refer to an available slot for the selected date
  - MovieId must refer to a valid movie in the movie service
- **Success Response (201 Created)**:
  ```json
  {
    "message": "Show created successfully",
    "request_id": "unique-request-id",
    "status": "SUCCESS",
    "data": {
      "id": 2,
      "movie": "tt1375666",
      "slot": {
        "id": 2,
        "name": "Afternoon",
        "startTime": "13:00:00.000000",
        "endTime": "16:00:00.000000"
      },
      "date": "2025-05-01",
      "cost": 250.50
    }
  }
  ```
- **Error Responses (400 Bad Request)**:
  ```json
  {
    "status": "ERROR",
    "code": "INVALID_COST",
    "message": "The cost must be less than or equal to 3000",
    "request_id": "unique-request-id"
  }
  ```
  ```json
  {
    "status": "ERROR",
    "code": "PAST_DATETIME",
    "message": "The show cannot be scheduled for a time in the past",
    "request_id": "unique-request-id"
  }
  ```
  ```json
  {
    "status": "ERROR",
    "code": "SLOT_NOT_AVAILABLE",
    "message": "The selected slot is not available on 2025-05-01",
    "request_id": "unique-request-id"
  }
  ```
  ```json
  {
    "status": "ERROR",
    "code": "INVALID_MOVIE",
    "message": "The selected movie does not exist",
    "request_id": "unique-request-id"
  }
  ```
- **Error Response (401 Unauthorized)**:
  ```json
  {
    "status": "ERROR",
    "code": "UNAUTHORIZED",
    "message": "Missing Authorization header",
    "request_id": "unique-request-id"
  }
  ```
- **Error Response (403 Forbidden)**:
  ```json
  {
    "status": "ERROR",
    "code": "FORBIDDEN",
    "message": "Access denied. Admin role required",
    "request_id": "unique-request-id"
  }
  ```

## Slot Management

### Get Available Slots
- **URL**: `/slot`
- **Method**: `GET`
- **Authentication**: Required (Admin only)
- **Query Parameters**:
  - `date`: Date in YYYY-MM-DD format (optional, defaults to current date)
- **Success Response (200 OK)**:
  ```json
  {
    "message": "Available slots retrieved successfully",
    "request_id": "unique-request-id",
    "status": "SUCCESS",
    "data": [
      {
        "id": 1,
        "name": "Morning",
        "startTime": "09:00:00.000000",
        "endTime": "12:00:00.000000"
      },
      {
        "id": 2,
        "name": "Afternoon",
        "startTime": "13:00:00.000000",
        "endTime": "16:00:00.000000"
      }
    ]
  }
  ```
- **Error Response (401 Unauthorized)**:
  ```json
  {
    "status": "ERROR",
    "code": "UNAUTHORIZED",
    "message": "Missing Authorization header",
    "request_id": "unique-request-id"
  }
  ```
- **Error Response (403 Forbidden)**:
  ```json
  {
    "status": "ERROR",
    "code": "FORBIDDEN",
    "message": "Access denied. Admin role required",
    "request_id": "unique-request-id"
  }
  ```

## Field Validation Rules

- **Name**: Must be 3-70 characters, max 4 words, letters only, no consecutive spaces
- **Username**: Must be 3-30 characters, lowercase, no spaces, cannot start with a number, no consecutive special characters
- **Password**: Must be at least 8 characters with at least one uppercase letter and one special character
- **Phone Number**: Must be exactly 10 digits
- **Email**: Must be in valid email format
- **Security Answer**: Must be at least 3 characters long
