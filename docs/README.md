# SkyFox Backend API Documentation

## Field Validation Rules

- **Name**: Must be 3-70 characters, max 4 words, letters only, no consecutive spaces
- **Username**: Must be 3-30 characters, lowercase, no spaces, cannot start with a number, no consecutive special characters
- **Password**: Must be at least 8 characters with at least one uppercase letter and one special character. Cannot match the previous 3 passwords stored in the database.
- **Phone Number**: Must be exactly 10 digits
- **Email**: Must be in valid email format
- **Security Answer**: Must be at least 3 characters long

## Authentication

### Login
- **URL**: `/login`
- **Method**: `POST`
- **Authentication**: None
- **Description**: Login and generate a JWT token with valid credentials.
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
- **Error Response (400 Unauthorized)**:
  ```json
  {
    "status": "ERROR",
    "code": "INVALID_REQUEST",
    "message": "Invalid request body",
    "request_id": "unique-request-id"
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
- **Description**: Get all security questions.
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
- **URL**: `/security-question/by-email?email=`
- **Method**: `POST`
- **Authentication**: None
- **Query Parameters**:
  - `email`: Validation performed
- **Description**: Get the security question of a valid user's correct email.
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
- **Description**: Verify the answer to the security question of a user and return a reset-token which expires in 300 seconds back to the user which is used for forgot password request validation.
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
- **Description**: Create a user as customer.
- **Notes**: 
  - This endpoint cannot be easily tested through Postman. Please refer to the Python script [here](../manual_tests/signup_test.py)
- **Request Body**:
  ```json
  {
    "name": "string",
    "username": "string",
    "password": "string",
    "number": "string",
    "email": "string",
    "profile_img": "base64 encoded image",
    "profile_img_sha":"sha256 hash of image",
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
- **Description**: Change the password of a user. Checks password history behind the scenes.
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

### Change Password
- **URL**: `/change-password`
- **Method**: `POST`
- **Authentication**: Required
- **Description**: Changes the user's password with current password verification and checks against password history.
- **Request Body**:
  ```json
  {
    "current_password": "string",
    "new_password": "string"
  }
  ```
- **Success Response (200 OK)**:
  ```json
  {
    "message": "Password updated successfully",
    "request_id": "4b16d5ca-6342-4c0f-ab18-3998bfb250c5",
    "status": "SUCCESS"
  }
  ```
- **Validation Error Response (400 Bad Request)**:
  ```json
  {
    "status": "ERROR",
    "code": "VALIDATION_ERROR",
    "message": "Validation failed",
    "request_id": "0a55356a-906f-48b7-b8c6-4a43ddf959ef",
    "errors": [
        {
            "field": "CurrentPassword",
            "message": "This field is required"
        },
        {
            "field": "NewPassword",
            "message": "This field is required"
        }
    ]
  }
  ```
- **Password Validation Error Response (400 Bad Request)**:
  ```json
  {
    "status": "ERROR",
    "code": "VALIDATION_ERROR",
    "message": "Validation failed",
    "request_id": "6fe2ef94-36e4-4f1d-987e-60a0fa7e7338",
    "errors": [
        {
            "field": "NewPassword",
            "message": "Password must be at least 8 characters with at least one uppercase letter and one special character"
        }
    ]
  }
  ```
- **Current Password Error Response (400 Bad Request)**:
  ```json
  {
    "status": "ERROR",
    "code": "INCORRECT PASSWORD",
    "message": "Current password doesn't match user's password",
    "request_id": "4350aaef-d290-4fdd-b7b4-c61baf220dfb"
  }
  ```
- **Password Reuse Error Response (400 Bad Request)**:
  ```json
  {
    "status": "ERROR",
    "code": "PASSWORD_REUSE",
    "message": "New password cannot match any of your previous passwords",
    "request_id": "503da576-f006-4bce-a1bb-55fa6f6d444a"
  }
  ```
- **Unauthorized Error Response (401 Unauthorized)**:
  ```json
  {
    "status": "ERROR",
    "code": "INVALID_TOKEN",
    "message": "Unauthorized",
    "request_id": "92bce3e3-2719-42ca-8388-1ec46fb1f677"
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
- **Description**: Fetch all shows for a given day.
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

### Get Show By ID

- **URL**: `/show`
- **Method**: `GET`
- **Authentication**: Required
- **Query Parameters**:
  - `id`: Show ID (must be a valid integer)
- **Description**: Retrieves detailed information for a specific show, including movie details, slot information, show cost, and available seats.
- **Success Response (200 OK)**
  ```json
  {
    "message": "Show fetched successfully",
    "request_id": "1319fff8-1d27-46c3-a062-5fd3ae61da18",
    "status": "SUCCESS",
    "data": {
      "movie": {
        "movieId": "tt1833116",
        "name": "The Informer",
        "duration": "1h53m0s",
        "plot": "An ex-convict working undercover intentionally gets himself incarcerated again in order to infiltrate the mob at a maximum security prison.",
        "imdbRating": "6.5",
        "moviePoster": "https://m.media-amazon.com/images/M/MV5BN2YyYTgxYmYtNjg3My00YzI4LWJlZWItYmZhZGEyYTYxNWY3XkEyXkFqcGdeQXVyMjAwNTYzNDg@._V1_SX300.jpg",
        "genre": "Crime, Drama, Thriller"
      },
      "slot": {
        "id": 3,
        "name": "Evening",
        "startTime": "17:00:00.000000",
        "endTime": "20:00:00.000000"
      },
      "id": 27,
      "date": "2025-04-27T00:00:00Z",
      "cost": 245.21,
      "availableseats": 99
    }
  }
  ```
- **Error Response (400 Bad Request)**
  ```json
  {
    "status": "ERROR",
    "code": "INVALID_SHOW_ID",
    "message": "Show id must be a valid integer",
    "request_id": "b327950f-ba86-432d-802f-c21256111cda"
  }
  ```
- **Error Response (401 Unauthorized)**
  ```json
  {
    "status": "ERROR",
    "code": "INVALID_TOKEN",
    "message": "Unauthorized",
    "request_id": "4dc04dae-fb48-4eaf-9499-9d75e1fdbd35"
  }
  ```
- **Error Response (404 Not Found)**
  ```json
  {
    "status": "ERROR",
    "code": "SHOW_NOT_FOUND",
    "message": "Show not found for id: 100",
    "request_id": "1e9fa455-17a0-4541-af8d-8d71ab33229b"
  }
  ```

### Get Movies
- **URL**: `/shows/movies`
- **Method**: `GET`
- **Authentication**: Required (Admin only)
- **Description**: Get all movies available via the movie service.
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
- **Description**: Schedule a show.
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
- **Description**: Retrieves all available slots on the requested date.
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

## Profile Management

### Get Customer Profile
- **URL**: `/customer/profile`
- **Method**: `GET`
- **Authentication**: Required (Customer Only)
- **Description**: Retrieves all customer profile information except the profile image.
- **Success Response (200 OK)**:
  ```json
  {
    "message": "Customer profile retrieved successfully",
    "request_id": "unique-request-id",
    "status": "SUCCESS",
    "data": {
        "username": "string",
        "name": "string",
        "email": "string",
        "phone_number": "string",
        "security_question_exists": boolean,
        "created_at": "ISO-8601 timestamp"
    }
  }
  ```
- **Error Response (401 Unauthorized)**:
  ```json
  {
    "status": "ERROR",
    "code": "INVALID_TOKEN",
    "message": "Unauthorized",
    "request_id": "unique-request-id"
  }
  ```
- **Error Response (403 Forbidden)**:
  ```json
  {
    "status": "ERROR",
    "code": "FORBIDDEN",
    "message": "Access denied. Customer role required",
    "request_id": "2326c83e-5a3d-4cd9-bcea-61deaa5bbedf"
  }
  ```

### Get Profile Image Presigned URL
- **URL**: `/customer/profile-image`
- **Method**: `GET`
- **Authentication**: Required (Customer Only)
- **Notes**: 
  - Generates a presigned URL for accessing the user's profile image
  - URL expires after 24 hours (1440 minutes)
- **Description**: Retrieve the image url of a valid customer.
- **Success Response (200 OK)**:
  ```json
  {
    "message": "Presigned URL generated successfully",
    "request_id": "unique-request-id",
    "status": "SUCCESS",
    "data": {
      "presigned_url": "https://bucket-name.s3.region.amazonaws.com/profile-images/username_timestamp.jpg?X-Amz-Algorithm=AWS4-HMAC-SHA256&X-Amz-Credential=...",
      "expires_at": "2025-04-08T00:17:00Z"
    }
  }
  ```
- **Error Response (401 Unauthorized)**:
  ```json
  {
    "status": "ERROR",
    "code": "UNAUTHORIZED",
    "message": "Unable to verify credentials",
    "request_id": "unique-request-id"
  }
  ```
- **Error Response (401 Unauthorized)**:
  ```json
  {
    "status": "ERROR",
    "code": "INVALID_TOKEN",
    "message": "Invalid token claims",
    "request_id": "unique-request-id"
  }
  ```
- **Error Response (404 Not Found)**:
  ```json
  {
    "status": "ERROR",
    "code": "PROFILE_IMAGE_NOT_FOUND",
    "message": "No profile image found for this user",
    "request_id": "unique-request-id"
  }
  ```

### Update Customer Profile
- **URL**: `/customer/update-profile`
- **Method**: `POST`
- **Authentication**: Required (Customer only)
- **Description**: Updates customer profile information with security question verification
- **Notes**:
  - The security answer must match the answer provided during signup
  - Email and phone number must be unique across all users
- **Request Body**:
  ```json
  {
    "name": "string",
    "email": "string",
    "phone_number": "string",
    "security_answer": "string"
  }
  ```
- **Success Response (200 OK)**:
  ```json
  {
    "message": "Profile updated successfully",
    "request_id": "99a4762c-4d27-424f-83be-40a592cfbc28",
    "status": "SUCCESS",
    "data": {
        "username": "suteerth",
        "name": "Suteerth S",
        "email": "suteerth1@gmail.com",
        "phone_number": "1234567892"
    }
  }
  ```
- **Error Response (403 Forbidden)**:
  ```json
  {
    "status": "ERROR",
    "code": "FORBIDDEN",
    "message": "Access denied. Customer role required",
    "request_id": "2326c83e-5a3d-4cd9-bcea-61deaa5bbedf"
  }
  ```
- **Invalid Request Error Response (400 Bad Request)**:
  ```json
  {
    "status": "ERROR",
    "code": "INVALID_REQUEST",
    "message": "Invalid request data",
    "request_id": "1ca237fd-a6d6-419e-b1a2-05ccde846447"
  }
  ```
- **Missing Fields Validation Error Response (400 Bad Request)**:
  ```json
  {
    "status": "ERROR",
    "code": "VALIDATION_ERROR",
    "message": "Validation failed",
    "request_id": "70d23976-ab8e-432e-9dda-62109eed19b5",
    "errors": [
        {
            "field": "Name",
            "message": "This field is required"
        },
        {
            "field": "Email",
            "message": "This field is required"
        },
        {
            "field": "PhoneNumber",
            "message": "This field is required"
        },
        {
            "field": "SecurityAnswer",
            "message": "This field is required"
        }
    ]
  }
  ```
- **Format Validation Error Response (400 Bad Request)**:
  ```json
  {
    "status": "ERROR",
    "code": "VALIDATION_ERROR",
    "message": "Validation failed",
    "request_id": "48cf1f2e-fe7f-4887-a75f-ba359801bd5f",
    "errors": [
        {
            "field": "Email",
            "message": "Invalid email format"
        },
        {
            "field": "PhoneNumber",
            "message": "Phone number must be exactly 10 digits"
        }
    ]
  }
  ```
- **Security Answer Error Response (400 Bad Request)**:
  ```json
  {
    "status": "ERROR",
    "code": "INVALID_SECURITY_ANSWER",
    "message": "The security answer provided is incorrect",
    "request_id": "3fabb098-188a-4c96-87f0-c5730b04cded"
  }
  ```
- **Duplicate Email Error Response (400 Bad Request)**:
  ```json
  {
    "status": "ERROR",
    "code": "email_EXISTS",
    "message": "Email already exists",
    "request_id": "77b0fc67-0b01-4494-bdd8-1e187c7f756c"
  }
  ```
- **Duplicate Phone Number Error Response (400 Bad Request)**:
  ```json
  {
    "status": "ERROR",
    "code": "mobilenumber_EXISTS",
    "message": "Mobile number already exists",
    "request_id": "418e8e21-3ba6-4010-8449-15ce24a12da6"
  }
  ```

### Update Profile Image
- **URL**: `/customer/update-profile-image`
- **Method**: `POST`
- **Authentication**: Required (Customer only)
- **Description**: Updates the customer's profile image with security question verification
- **Notes**: 
  - This endpoint cannot be easily tested through Postman. Please refer to the Python script [here](../manual_tests/update_prof_image_test.py)
  - The old image in S3 is automatically deleted when a new one is uploaded
- **Request Body**:
  ```json
  {
    "security_answer": "string",
    "profile_img": "base64 encoded image",
    "profile_img_sha": "sha256 hash of image"
  }
  ```
- **Success Response (200 OK)**:
  ```json
  {
    "message": "Profile image updated successfully",
    "request_id": "37137e7d-549f-4644-9abd-f00516d76a4a",
    "status": "SUCCESS",
  }
  ```
- **Error Response (403 Forbidden)**:
  ```json
  {
    "status": "ERROR",
    "code": "FORBIDDEN",
    "message": "Access denied. Customer role required",
    "request_id": "c7e82ff5-b9f6-4c61-b343-81af58d585fb"
  }
  ```
- **Validation Error Response (400 Bad Request)**:
  ```json
  {
    "status": "ERROR",
    "code": "VALIDATION_ERROR",
    "message": "Validation failed",
    "request_id": "60ff21e3-63a7-43b4-a1f7-80f1364b35f9",
    "errors": [
        {
            "field": "SecurityAnswer",
            "message": "This field is required"
        },
        {
            "field": "ProfileImg",
            "message": "This field is required"
        },
        {
            "field": "ProfileImgSHA",
            "message": "This field is required"
        }
    ]
  }
  ```
- **Security Answer Error Response (400 Bad Request)**:
  ```json
  {
    "status": "ERROR",
    "code": "INVALID_SECURITY_ANSWER",
    "message": "The security answer provided is incorrect",
    "request_id": "e3c50ab1-6dd6-41a6-97ea-f49428ea3192"
  }
  ```
- **Invalid Image Error Response (400 Bad Request)**:
  ```json
  {
    "status": "ERROR",
    "code": "INVALID_IMAGE",
    "message": "Invalid base64 image data",
    "request_id": "unique-request-id"
  }
  ```
- **Invalid Image Hash Error Response (400 Bad Request)**:
  ```json
  {
    "status": "ERROR",
    "code": "INVALID_IMAGE_HASH",
    "message": "The image hash does not match the provided image",
    "request_id": "unique-request-id"
  }
  ```
- **S3 Error Response (500 Internal Server Error)**:
  ```json
  {
    "status": "ERROR",
    "code": "S3_DELETE_FAILED",
    "message": "Failed to delete profile image",
    "request_id": "unique-request-id"
  }
  ```

### Get Admin Profile
- **URL**:`/admin/profile`
- **Method**: `GET`
- **Authentication:** Required (Admin role only)
- **Description:** Retrieves the profile information for an authenticated admin user.
- **Success Response (200 OK)**:
  ```json
  {
      "message": "Admin profile retrieved successfully",
      "request_id": "64bea111-78e1-4e12-8c2b-8452463a86f4",
      "status": "SUCCESS",
      "data": {
          "username": "seed-user-1",
          "name": "Admin One",
          "counter_no": 101,
          "created_at": "2025-04-15T11:30:06+05:30"
      }
  }
  ```
- **Error Response (401 Unauthorized)**:
  ```json
  {
      "status": "ERROR",
      "code": "INVALID_TOKEN",
      "message": "Unauthorized",
      "request_id": "aed33148-d8a0-4657-9816-64b0e64c0467"
  }
  ```
- **Error Response (403 Forbidden)**:
  ```json
  {
      "status": "ERROR",
      "code": "FORBIDDEN",
      "message": "Access denied. Admin or staff role required!",
      "request_id": "938b985d-d729-4b67-9509-077d3a30ab74"
  }
  ```

### Get Staff Profile
- **URL**:`/staff/profile`
- **Method**: `GET`
- **Authentication:** Required (Staff role only)
- **Description:** Retrieves the profile information for an authenticated staff user.
- **Success Response (200 OK)**:
  ```json
  {
      "message": "Staff profile retrieved successfully",
      "request_id": "dad89523-8ba7-476d-bafe-1db1a0c0d05d",
      "status": "SUCCESS",
      "data": {
          "username": "staff-1",
          "name": "Staff One",
          "counter_no": 501,
          "created_at": "2025-04-15T11:30:07+05:30"
      }
  }
  ```
- **Error Response (401 Unauthorized)**:
  ```json
  {
      "status": "ERROR",
      "code": "INVALID_TOKEN",
      "message": "Unauthorized",
      "request_id": "e4a44212-1a2b-4c48-9be8-ec5e935b7ebe"
  }
  ```
- **Error Response (403 Forbidden)**:
  ```json
  {
      "status": "ERROR",
      "code": "FORBIDDEN",
      "message": "Access denied. Admin or staff role required!",
      "request_id": "938b985d-d729-4b67-9509-077d3a30ab74"
  }
  ```

## Booking Management

### Get Seat Map
- **URL**: `/shows/{show_id}/seat-map`
- **Method**: `GET`
- **Authentication**: Required
- **Parameters**:
  - `show_id`: ID of the show (must be a valid integer)
- **Description**: Retrieves a complete seat map for a specific show, including seat availability status, seat type, and pricing information.
- **Notes**: 
  - Standard seats are priced at the show's base cost
  - Deluxe seats have an additional price premium (currently 150.00)
  - Rows A-E are Standard seats, rows F-J are Deluxe seats
  - Each row has 10 seats numbered 1-10
  - Occupied seats cannot be booked

- **Success Response (200 OK)**:
  ```json
  {
    "message": "Seat map retrieved successfully",
    "request_id": "c2b01fde-fbe5-4d85-830e-e8b84b798b34",
    "status": "SUCCESS",
    "data": {
      "seat_map": {
        "A": [
          {
            "column": "1",
            "occupied": false,
            "price": 208.97,
            "seat_number": "A1",
            "type": "Standard"
          },
          {
            "column": "2",
            "occupied": false,
            "price": 208.97,
            "seat_number": "A2",
            "type": "Standard"
          },
          // Additional seats...
        ],
        // Additional rows B through J...
      }
    }
  }
  ```

- **Error Response (400 Bad Request)**:
  ```json
  {
    "status": "ERROR",
    "code": "INVALID_SHOW_ID",
    "message": "Show ID must be a valid integer",
    "request_id": "2870e4f4-3918-42f6-b28a-304193ebc720"
  }
  ```

- **Error Response (401 Unauthorized)**:
  ```json
  {
    "status": "ERROR",
    "code": "INVALID_TOKEN",
    "message": "Unauthorized",
    "request_id": "8c55aba3-fd1a-4ed5-9655-5dd28318d28e"
  }
  ```

- **Error Response (404 Not Found)**:
  ```json
  {
    "status": "ERROR",
    "code": "SHOW_NOT_FOUND",
    "message": "Show not found for id: 1000",
    "request_id": "fb8b880a-4d97-4e33-ad7a-5f137e205312"
  }
  ```

### Create Booking For Customer as Admin
- **URL**: `/admin/create-customer-booking`
- **Method**: `POST`
- **Authentication**: Required (Admin only)
- **Description**: Create a booking directly by an admin (offline booking).
- **Notes**:
  - Admin bookings are immediately confirmed with "Cash" payment type
  - Amount paid can be provided or calculated automatically based on selected seats
  - Booking cannot be created for shows that have already started
  - Each admin booking creates a new admin_booked_customer record
- **Request Body**:
  ```json
  {
    "show_id": 22,
    "customer_name": "John Doe",
    "phone_number": "1234567890",
    "seat_numbers": ["A1", "J1"],
    "amount_paid": 553.30
  }
  ```
- **Success Response (201 Created)**:
  ```json
  {
    "message": "Booking created successfully",
    "request_id": "6c3bcc3d-f146-4afc-b03e-60f218c379e8",
    "status": "SUCCESS",
    "data": {
      "booking_id": 1,
      "show_id": 22,
      "customer_name": "John Doe",
      "phone_number": "1234567890",
      "seat_numbers": [
        "A1",
        "J1"
      ],
      "amount_paid": 553.3,
      "payment_type": "Cash",
      "booking_time": "2025-04-22T16:26:53.835139+05:30",
      "status": "Confirmed"
    }
  }
  ```
- **Error Response (400 Bad Request) - Invalid Input**:
  ```json
  {
    "status": "ERROR",
    "code": "VALIDATION_ERROR",
    "message": "Validation failed",
    "request_id": "33749b29-e38f-4aa3-91f6-336ce9d8fc2a",
    "errors": [
      {
        "field": "CustomerName",
        "message": "Name must be 3-70 characters, max 4 words, letters only, no consecutive spaces"
      },
      {
        "field": "PhoneNumber",
        "message": "Phone number must be exactly 10 digits"
      },
      {
        "field": "AmountPaid",
        "message": "Invalid value"
      }
    ]
  }
  ```
- **Error Response (400 Bad Request) - Too Many Seats In a Booking**:
  ```json
  {
    "status": "ERROR",
    "code": "TOO_MANY_SEATS",
    "message": "Maximum 10 seats can be booked per booking",
    "request_id": "5fce2cc9-eb5e-4c01-8f4e-8c818431dbb7"
  } 
  ```
- **Error Response (400 Bad Request) - Price Mismatch**:
  ```json
  {
    "status": "ERROR",
    "code": "INVALID_AMOUNT",
    "message": "The amount paid does not match the expected price",
    "request_id": "446ed540-c84f-4264-b1f6-f3958792cc90"
  }
  ```
- **Error Response (400 Bad Request) - Show Started**:
  ```json
  {
    "status": "ERROR",
    "code": "SHOW_ALREADY_STARTED",
    "message": "Cannot book tickets for a show that has already started",
    "request_id": "5d3a483e-1489-4011-a5de-7192ba65a560"
  }
  ```
- **Error Response (400 Bad Request) - Seats Unavailable**:
  ```json
  {
    "status": "ERROR",
    "code": "SEATS_UNAVAILABLE",
    "message": "One or more selected seats are not available",
    "request_id": "ea471ac4-072a-456e-b1d6-b757c9abbf9e"
  }
  ```
- **Error Response (400 Bad Request) - Invalid Request**:
  ```json
  {
    "status": "ERROR",
    "code": "INVALID_REQUEST",
    "message": "Invalid request data",
    "request_id": "2de78e1f-95a1-44c9-b6fa-7fa0ff1d9d60"
  }
  ```
- **Error Response (403 Forbidden)**:
  ```json
  {
    "status": "ERROR",
    "code": "FORBIDDEN",
    "message": "Access denied. Admin role required",
    "request_id": "869a6343-051c-4063-8f01-524521e64cfb"
  }
  ```

### Initialize Customer Booking
- **URL**: `/customer/booking/initialize`
- **Method**: `POST`
- **Authentication**: Required (Customer only)
- **Description**: Creates a temporary booking with reserved seats and a 5-minute expiration time.
- **Notes**:
  - Seats are temporarily reserved for 5 minutes, allowing time for payment
  - If payment is not completed before expiration, seats are automatically released
  - Booking status is set to "Pending" until payment is processed
- **Request Body**:
  ```json
  {
    "show_id": 22,
    "seat_numbers": ["A1", "J1"]
  }
  ```
- **Success Response (201 Created)**:
  ```json
  {
    "message": "Booking initialized successfully",
    "request_id": "b4408a16-4fc7-4caf-88ff-b1340968947d",
    "status": "SUCCESS",
    "data": {
      "booking_id": 7,
      "show_id": 22,
      "seat_numbers": [
        "A1",
        "J1"
      ],
      "amount_due": 553.3,
      "expiration_time": "2025-04-23T16:32:07.831347491+05:30",
      "time_remaining_ms": 300000
    }
  }
  ```
- **Error Response (403 Forbidden)**:
  ```json
  {
    "status": "ERROR",
    "code": "FORBIDDEN",
    "message": "Access denied. Customer role required",
    "request_id": "19c3079b-15c9-4671-aa0c-1f0ebdfdb436"
  }
  ```
- **Error Response (400 Bad Request) - Invalid JSON**:
  ```json
  {
    "status": "ERROR",
    "code": "INVALID_REQUEST",
    "message": "Invalid request data",
    "request_id": "bff152e2-3e9a-479f-949a-a184b72adc40"
  }
  ```
- **Error Response (400 Bad Request) - Too Many Seats In a Booking**:
  ```json
  {
    "status": "ERROR",
    "code": "TOO_MANY_SEATS",
    "message": "Maximum 10 seats can be booked per booking",
    "request_id": "22f9a1b5-75c5-4e49-88ff-bc10d033b933"
  }
  ```
- **Error Response (400 Bad Request) - Missing Fields**:
  ```json
  {
    "status": "ERROR",
    "code": "VALIDATION_ERROR",
    "message": "Validation failed",
    "request_id": "d0063421-5ff1-4f34-bacd-87daf76a7c6b",
    "errors": [
      {
        "field": "SeatNumbers",
        "message": "This field is required"
      }
    ]
  }
  ```
- **Error Response (404 Not Found) - Show Not Found**:
  ```json
  {
    "status": "ERROR",
    "code": "SHOW_NOT_FOUND",
    "message": "Show not found for id: 999",
    "request_id": "249e34df-6f52-47c5-943d-0da4e4115197"
  }
  ```
- **Error Response (400 Bad Request) - Show Already Started**:
  ```json
  {
    "status": "ERROR",
    "code": "SHOW_ALREADY_STARTED",
    "message": "Cannot book tickets for a show that has already started",
    "request_id": "3a212b2a-e47a-4cf5-9b86-8a2e8cf714da"
  }
  ```
- **Error Response (400 Bad Request) - Seats Unavailable**:
  ```json
  {
    "status": "ERROR",
    "code": "SEATS_UNAVAILABLE",
    "message": "One or more selected seats are not available",
    "request_id": "03867a26-d12b-47b6-9d20-bdacd13ae466"
  }
  ```
- **Error Response (400 Bad Request) - Invalid Seat Format**:
  ```json
  {
    "status": "ERROR",
    "code": "VALIDATION_ERROR",
    "message": "Validation failed",
    "request_id": "c13a7284-4b3e-44d8-aad4-1bc0eec404ef",
    "errors": [
      {
        "field": "SeatNumbers[1]",
        "message": "Must be at most 3 characters"
      }
    ]
  }
  ```

### Process Payment for Booking
- **URL**: `/customer/booking/:id/payment`
- **Method**: `POST`
- **Authentication**: Required (Customer only)
- **Description**: Processes payment for a pending booking and confirms it.
- **Notes**:
  - Must be completed within the 5-minute expiration window
  - Only the customer who created the booking can process payment
  - Successfully processed bookings are set to "Confirmed" status
  - Payment validation includes Luhn check for card number and expiry date validation
- **Request Body**:
  ```json
  {
    "booking_id": 10,
    "card_number": "4111111111111111",
    "cvv": "123",
    "expiry_month": "12",
    "expiry_year": "25",
    "cardholder_name": "John Doe"
  }
  ```
- **Success Response (200 OK)**:
  ```json
  {
    "message": "Payment processed successfully",
    "request_id": "969a7924-f48e-4ee0-9b1a-f20b046f2adc",
    "status": "SUCCESS",
    "data": {
      "booking_id": 10,
      "show_id": 22,
      "show_date": "2025-04-26",
      "show_time": "13:00:00.000000",
      "movie": {
        "movieId": "tt7349950",
        "name": "It Chapter Two",
        "duration": "2h49m0s",
        "plot": "Twenty-seven years after their first encounter with the terrifying Pennywise, the Losers Club have grown up and moved away, until a devastating phone call brings them back.",
        "imdbRating": "6.6",
        "moviePoster": "https://m.media-amazon.com/images/M/MV5BYTJlNjlkZTktNjEwOS00NzI5LTlkNDAtZmEwZDFmYmM2MjU2XkEyXkFqcGdeQXVyNjg2NjQwMDQ@._V1_SX300.jpg",
        "genre": "Drama, Fantasy, Horror"
      },
      "seat_numbers": [
        "B1",
        "I1"
      ],
      "amount_paid": 553.3,
      "payment_type": "Card",
      "booking_time": "2025-04-23T16:46:18.154511+05:30",
      "status": "Confirmed",
      "transaction_id": "71d2b74c-4f79-49d1-9a6c-7eddaaf454a5"
    }
  }
  ```
- **Error Response (403 Forbidden)**:
  ```json
  {
    "status": "ERROR",
    "code": "FORBIDDEN",
    "message": "Access denied. Customer role required",
    "request_id": "19c3079b-15c9-4671-aa0c-1f0ebdfdb436"
  }
  ```
- **Error Response (400 Bad Request) - Missing Fields**:
  ```json
  {
    "status": "ERROR",
    "code": "VALIDATION_ERROR",
    "message": "Validation failed",
    "request_id": "d67c5c25-a9e6-4285-9467-be545283c18c",
    "errors": [
      {
        "field": "CVV",
        "message": "This field is required"
      },
      {
        "field": "ExpiryMonth",
        "message": "This field is required"
      },
      {
        "field": "ExpiryYear",
        "message": "This field is required"
      },
      {
        "field": "CardholderName",
        "message": "This field is required"
      }
    ]
  }
  ```
- **Error Response (404 Not Found) - Booking Not Found**:
  ```json
  {
    "status": "ERROR",
    "code": "BOOKING_NOT_FOUND",
    "message": "Booking not found",
    "request_id": "40f12206-504b-49d1-97b2-11a89fc8ab37"
  }
  ```
- **Error Response (400 Bad Request) - Booking Expired**:
  ```json
  {
    "status": "ERROR",
    "code": "BOOKING_EXPIRED",
    "message": "This booking has expired. Please make a new booking",
    "request_id": "da3b799e-e464-4d68-b8eb-c417234477ae"
  }
  ```
- **Error Response (400 Bad Request) - Invalid Card Number**:
  ```json
  {
    "status": "ERROR",
    "code": "PAYMENT_VALIDATION_FAILED",
    "message": "Payment validation failed: Card number failed Luhn check (card_number)",
    "request_id": "67b0d312-c13c-41aa-adfe-8f5a5fd69e7a"
  }
  ```
- **Error Response (400 Bad Request) - Invalid Expiry Format**:
  ```json
  {
    "status": "ERROR",
    "code": "PAYMENT_VALIDATION_FAILED",
    "message": "Payment validation failed: Card number failed Luhn check (card_number), Expiry must be in MM/YY format with valid month (01-12) (expiry)",
    "request_id": "71d9406d-bef1-4ffa-999e-658d78b95822"
  }
  ```
- **Error Response (400 Bad Request) - Expired Card**:
  ```json
  {
    "status": "ERROR",
    "code": "PAYMENT_VALIDATION_FAILED",
    "message": "Payment validation failed: Card number failed Luhn check (card_number), Card has expired (expiry)",
    "request_id": "e86e6041-e979-462c-be73-28af37cb4b65"
  }
  ```
- **Error Response (400 Bad Request) - Invalid Booking Status**:
  ```json
  {
    "status": "ERROR",
    "code": "INVALID_BOOKING_STATUS",
    "message": "Payment can only be processed for bookings in pending state",
    "request_id": "1f3a7932-b416-4335-b178-383b6fa8d5d2"
  }
  ```
- **Error Response (403 Forbidden) - Unauthorized Access**:
  ```json
  {
    "status": "ERROR",
    "code": "UNAUTHORIZED_ACCESS",
    "message": "You don't have permission to access this booking",
    "request_id": "896e64d8-cbba-4aed-acc2-ae17f1f8348f"
  }
  ```

### Cancel Pending Booking
- **URL**: `/customer/booking/:id/cancel`
- **Method**: `DELETE`
- **Authentication**: Required (Customer only)
- **Parameters**:
  - `id`: Booking ID (must be a valid integer)
- **Description**: Cancels a pending booking, releases reserved seats, and terminates the associated monitoring goroutine.
- **Notes**:
  - Only the customer who created the booking can cancel it
  - Only bookings with "Pending" status can be cancelled
  - The endpoint is typically triggered when a user refreshes the page, closes the dialog, or explicitly cancels a transaction

- **Success Response (200 OK)**:
  ```json
  {
    "message": "Booking cancelled successfully",
    "request_id": "95950d5d-2571-4ca8-adb2-a66a444415a8",
    "status": "SUCCESS"
  }
  ```

- **Error Response (400 Bad Request)**:
  ```json
  {
    "status": "ERROR",
    "code": "INVALID_BOOKING_ID",
    "message": "Booking ID must be a valid integer",
    "request_id": "66301062-8152-4a17-9093-3bde55d147c3"
  }
  ```

- **Error Response (401 Unauthorized)**:
  ```json
  {
    "status": "ERROR",
    "code": "INVALID_TOKEN",
    "message": "Unauthorized",
    "request_id": "1290b3cd-edf0-474c-adae-3a0855f08795"
  }
  ```

- **Error Response (403 Forbidden) - Role Restriction**:
  ```json
  {
    "status": "ERROR",
    "code": "FORBIDDEN",
    "message": "Access denied. Customer role required",
    "request_id": "6b9308b4-aff0-4aa3-bdfd-828b7ad1d3f8"
  }
  ```

- **Error Response (403 Forbidden) - Unauthorized Access**:
  ```json
  {
    "status": "ERROR",
    "code": "UNAUTHORIZED_ACCESS",
    "message": "You don't have permission to access this booking",
    "request_id": "3e3668ea-3068-4bb3-a454-664783d1d55b"
  }
  ```

- **Error Response (404 Not Found)**:
  ```json
  {
    "status": "ERROR",
    "code": "BOOKING_NOT_FOUND",
    "message": "Booking not found",
    "request_id": "ec213609-f4ae-491c-a9f9-0b083a68fe2f"
  }
  ```

### Get QR Code
- **URL**: `/booking/:id/qr`
- **Method**: `GET`
- **Authentication**: Required
- **Parameters**:
  - `id`: Booking ID (must be a valid integer)
- **Description**: Generates and returns a QR code for a booking.
- **Notes**:
  - Admin and staff can access QR codes for any booking
  - Customers can only access QR codes for their own bookings
  - QR code contains essential booking information for verification
- **Success Response (200 OK)**:
  ```json
  {
    "message": "QR code generated successfully",
    "request_id": "f02b2ef2-e849-4d7a-b4fd-7c212048c56d",
    "status": "SUCCESS",
    "data": {
        "encoding": "base64",
        "mime_type": "image/png",
        "qr_code": "Base64 Content"
    }
  }
  ```
- **Error Response (400 Bad Request)**:
  ```json
  {
    "status": "ERROR",
    "code": "INVALID_BOOKING_ID",
    "message": "Invalid booking ID format",
    "request_id": "413c5464-17b6-435f-b4f3-350525f0fac3"
  }
  ```
- **Error Response (404 Not Found)**:
  ```json
  {
    "status": "ERROR",
    "code": "BOOKING_NOT_FOUND",
    "message": "Booking not found",
    "request_id": "08e35df2-25a4-422b-b66e-2114e72aa288"
  }
  ```
- **Error Response (403 Forbidden)**:
  ```json
  {
    "status": "ERROR",
    "code": "FORBIDDEN",
    "message": "Access denied to this booking",
    "request_id": "3987dfd9-7fca-4a0c-8a9a-847e1dc606ff"
  }
  ```

### Get PDF Ticket
- **URL**: `/booking/:id/pdf`
- **Method**: `GET`
- **Authentication**: Required
- **Parameters**:
  - `id`: Booking ID (must be a valid integer)
- **Description**: Generates and returns a PDF ticket for a booking.
- **Notes**:
  - Admin and staff can access PDF tickets for any booking
  - Customers can only access PDF tickets for their own bookings
  - PDF includes booking details, QR code, and theater information
- **Success Response (200 OK)**:
  ```json
  {
    "message": "PDF ticket generated successfully",
    "request_id": "e935191f-5883-492c-86b7-a6adb19008aa",
    "status": "SUCCESS",
    "data": {
        "encoding": "base64",
        "mime_type": "application/pdf",
        "pdf": "Base64 Content"
    }
  }
  ```
- **Error Response (400 Bad Request)**:
  ```json
  {
    "status": "ERROR",
    "code": "INVALID_BOOKING_ID",
    "message": "Invalid booking ID format",
    "request_id": "413c5464-17b6-435f-b4f3-350525f0fac3"
  }
  ```
- **Error Response (404 Not Found)**:
  ```json
  {
    "status": "ERROR",
    "code": "BOOKING_NOT_FOUND",
    "message": "Booking not found",
    "request_id": "08e35df2-25a4-422b-b66e-2114e72aa288"
  }
  ```
- **Error Response (403 Forbidden)**:
  ```json
  {
    "status": "ERROR",
    "code": "FORBIDDEN",
    "message": "Access denied to this booking",
    "request_id": "3987dfd9-7fca-4a0c-8a9a-847e1dc606ff"
  }
  ```

### Get All Bookings for Customer
- **URL**: `/customer/bookings`
- **Method**: `GET`
- **Authentication**: Required (Customer role only)
- **Description**: Retrieves all bookings made by the authenticated customer, ordered by booking time (latest first).
- **Success Response (200 OK) - Multiple Bookings**
  ```json
  {
      "message": "Bookings fetched successfully",
      "request_id": "a727277a-b0bc-4fb3-bf84-d8b65b8542b6",
      "status": "SUCCESS",
      "data": [
          {
              "booking_id": 54,
              "show_id": 25,
              "show_date": "2025-04-27",
              "show_time": "09:00:00.000000",
              "seat_numbers": ["E6"],
              "amount_paid": 297.34,
              "payment_type": "Card",
              "booking_time": "2025-04-26T23:56:20.062442+05:30",
              "status": "Confirmed"
          },
          {
              "booking_id": 11,
              "show_id": 22,
              "show_date": "2025-04-26",
              "show_time": "13:00:00.000000",
              "seat_numbers": ["A1","J1"],
              "amount_paid": 553.3,
              "payment_type": "Card",
              "booking_time": "2025-04-23T16:57:41.680409+05:30",
              "status": "Confirmed"
          }
          // ...more bookings
      ]
  }
  ```
- **Success Response (200 OK) – No Bookings**
  ```json
  {
      "message": "Bookings fetched successfully",
      "request_id": "54959558-4f68-49a0-998e-d0077b723638",
      "status": "SUCCESS",
      "data": []
  }
  ```
- **Error Response (401 Unauthorized)**
  ```json
  {
      "status": "ERROR",
      "code": "INVALID_TOKEN",
      "message": "Unauthorized",
      "request_id": "74153dac-fb26-4370-b718-bcf01df71b1b"
  }
  ```
- **Error Response (403 Forbidden)**
  ```json
  {
      "status": "ERROR",
      "code": "FORBIDDEN",
      "message": "Access denied. Customer role required",
      "request_id": "e9b5a26f-6bf1-4a5d-a04d-2a64ab837374"
  }
  ```

### Get Latest Booking for Customer
- **URL**: `/customer/bookings/latest`
- **Method**: `GET`
- **Authentication**: Required (Customer role only)
- **Description**: Retrieves the most recent booking made by the authenticated customer.
- **Success Response (200 OK) – With Booking**
  ```json
  {
      "message": "Latest booking fetched successfully",
      "request_id": "a69b2734-0222-4f5d-9688-6fa125208cb5",
      "status": "SUCCESS",
      "data": {
          "booking_id": 48,
          "show_id": 25,
          "show_date": "2025-04-27",
          "show_time": "09:00:00.000000",
          "seat_numbers": ["F4"],
          "amount_paid": 447.34,
          "payment_type": "Card",
          "booking_time": "2025-04-26T23:34:44.643438+05:30",
          "status": "Confirmed"
      }
  }
  ```
- **Success Response (200 OK) – No Booking Found**
  ```json
  {
      "message": "Latest booking fetched successfully",
      "request_id": "77fb3ba2-33e1-4fbb-bb3b-8cfdf112a4b5",
      "status": "SUCCESS",
      "data": null
  }
  ```
- **Error Response (401 Unauthorized)**
  ```json
  {
      "status": "ERROR",
      "code": "INVALID_TOKEN",
      "message": "Unauthorized",
      "request_id": "9f3ec42d-480a-4693-977d-97ba4332969b"
  }
  ```
- **Error Response (403 Forbidden)**
  ```json
  {
      "status": "ERROR",
      "code": "FORBIDDEN",
      "message": "Access denied. Customer role required",
      "request_id": "e9b5a26f-6bf1-4a5d-a04d-2a64ab837374"
  }
  ```

### Get All Bookings Requiring Check-in 
- **URL:** `/check-in/bookings`  
- **Method:** `GET`  
- **Authentication:** Required (Admin/Staff role)  
- **Description:** Returns all bookings in "Confirmed" status awaiting check-in.
- **Success Response (200)**
  ```json
  {
      "message": "Confirmed bookings fetched successfully",
      "request_id": "4ee5868f-299b-4fb8-9b62-cab29b4677fe",
      "status": "SUCCESS",
      "data": [
          {
              "id": 63,
              "date": "2025-04-28T00:00:00Z",
              "show_id": 32,
              "customer_id": null,
              "customer_username": "venkat",
              "no_of_seats": 1,
              "amount_paid": 224.29,
              "status": "Confirmed",
              "booking_time": "2025-04-28T16:52:41.215358+05:30",
              "payment_type": "Card"
          }
          // ... more bookings
      ]
  }
  ```
- **Forbidden (403)**
  ```json
  {
      "status": "ERROR",
      "code": "FORBIDDEN",
      "message": "Access denied. Admin or staff role required!",
      "request_id": "79fc1d07-2235-4ba2-a43e-f63ac7adf779"
  }
  ```

### Bulk Check-In Bookings (Admin/Staff only)
- **URL:** `/check-in/bookings`  
- **Method:** `POST`  
- **Authentication:** Required (Admin/Staff role)  
- **Description:** Attempts to check in the list of booking IDs. Already checked-in, invalid, or expired bookings are skipped, and the response details which bookings were checked in, already done, or invalid.
- **Request Body**
  ```json
  {
      "booking_ids": [65, 64, 66]
  }
  ```
- **Success Response (200)**
  ```json
  {
      "message": "Bulk check-in attempted",
      "request_id": "a98d74ce-b903-4ef8-a1fb-3cc4cc5c26a8",
      "status": "SUCCESS",
      "data": {
          "checked_in": [65, 64],
          "already_done": [],
          "invalid": [66]
      }
  }
  ```
- **Invalid Input (400)**
  ```json
  {
      "status": "ERROR",
      "code": "INVALID_INPUT",
      "message": "Invalid input",
      "request_id": "c376b8fa-ea53-45aa-8067-cc2aa7407980"
  }
  ```
- **Forbidden (403)**
  ```json
  {
      "status": "ERROR",
      "code": "FORBIDDEN",
      "message": "Access denied. Admin or staff role required!",
      "request_id": "79fc1d07-2235-4ba2-a43e-f63ac7adf779"
  }
  ```

### Single Booking Check-In (Admin/Staff only)
- **URL:** `/check-in/booking`  
- **Method:** `POST`  
- **Authentication:** Required (Admin/Staff role)  
- **Description:** Attempts to check in a single booking by ID. Returns which category the booking falls into (`checked_in`, `already_done`, or `invalid`)
- **Request Body:**
  ```json
  {
      "booking_id": 68
  }
  ```
- **Success Response (200) – Valid, Already Done, or Invalid**
  ```json
  {
      "message": "Check-in failed: invalid booking (already expired/invalid status/or show ended)",
      "request_id": "6ba3026d-e175-46bd-bbe7-8121055df759",
      "status": "SUCCESS",
      "data": {
          "checked_in": [],
          "already_done": [],
          "invalid": [68]
      }
  }
  ```
- **Invalid Input (400)**
  ```json
  {
      "status": "ERROR",
      "code": "INVALID_INPUT",
      "message": "Invalid input",
      "request_id": "c376b8fa-ea53-45aa-8067-cc2aa7407980"
  }
  ```
- **Forbidden (403)**
  ```json
  {
      "status": "ERROR",
      "code": "FORBIDDEN",
      "message": "Access denied. Admin or staff role required!",
      "request_id": "79fc1d07-2235-4ba2-a43e-f63ac7adf779"
  }
  ```

## Dashboard - Revenue

The Revenue Dashboard API provides a powerful way to analyze booking revenue data across various dimensions. This API supports dynamic filtering, grouping, and aggregation to help you understand booking patterns and revenue trends.

### How to Use Query Parameters

The Revenue API uses query parameters to filter and group data:

- **Timeframe Parameters**: Group data by time periods (`timeframe=daily|weekly|monthly|yearly`)
  - `daily`: Past 30 days
  - `weekly`: Past 16 weeks
  - `monthly`: Past 12 months
  - `yearly`: No soft limit
- **Period Filters**: Filter by specific time periods (`month=1-12`, `year=YYYY`)
- **Dimension Filters**: Filter by booking properties (`movie_id`, `slot_id`, `genre`)

### Important Rules

1. **Parameter Order Matters**: The order of parameters in your query determines the order of components in the response labels (separated by semicolons)
2. **Mutual Exclusivity**: `timeframe` parameter cannot be combined with `month` or `year` parameters
3. **Authentication**: This API requires admin authentication
4. **Empty Results**: When no data matches your filters, the API returns empty arrays (not null)

### Suggested Usage

For the best frontend experience, implement the API with dropdown filters that dynamically build the query URL:

1. Start with overall revenue statistics (no filters)
2. Add dropdown selectors for:
   - Timeframe (daily/weekly/monthly/yearly)
   - Movie selection
   - Slot selection
   - Genre selection
   - Month/year selection (when timeframe isn't used)
3. Update the URL and fetch data when filters change
4. Display aggregated stats in cards and grouped data in charts

### Revenue - All Data (No Filters)

- **URL**: `/revenue`
- **Method**: `GET`
- **Authentication**: Required (Admin role)
- **Description**: Returns aggregated revenue statistics across all bookings with "Confirmed" or "CheckedIn" status.

- **Success Response (200 OK)**:
  ```json
  {
    "message": "Revenue data fetched successfully",
    "request_id": "7823a569-c8dd-4a54-813a-84eb5ece9ae3",
    "status": "SUCCESS",
    "data": {
      "total_revenue": 8271.09,
      "mean_revenue": 689.2575,
      "median_revenue": 447.34,
      "total_bookings": 12,
      "total_seats_booked": 25,
      "groups": [
        {
          "label": "All",
          "total_revenue": 8271.09,
          "mean_revenue": 689.2575,
          "median_revenue": 447.34,
          "total_bookings": 12,
          "total_seats_booked": 25
        }
      ]
    }
  }
  ```

- **Unauthorized Response (401)**:
  ```json
  {
    "status": "ERROR",
    "code": "INVALID_TOKEN",
    "message": "Unauthorized",
    "request_id": "761234df-ff53-4bbb-882a-039925807c74"
  }
  ```

- **Forbidden Response (403)**:
  ```json
  {
    "status": "ERROR",
    "code": "FORBIDDEN",
    "message": "Access denied. Admin role required",
    "request_id": "4ccefb6e-f3bd-464b-8ce4-2f30ee4055f0"
  }
  ```

### Revenue - Timeframe Grouping

- **URL**: `/revenue?timeframe=daily|weekly|monthly|yearly`
- **Method**: `GET`
- **Authentication**: Required (Admin role)
- **Description**: Groups revenue data by the specified timeframe (daily, weekly, monthly, or yearly).

- **Success Response (200 OK) - Daily Timeframe**:
  ```json
  {
    "message": "Revenue data fetched successfully",
    "request_id": "8dfb6e11-0a2b-47e2-98b5-1e05c6f3bda9",
    "status": "SUCCESS",
    "data": {
      "total_revenue": 8271.09,
      "mean_revenue": 689.2575,
      "median_revenue": 447.34,
      "total_bookings": 12,
      "total_seats_booked": 25,
      "groups": [
        {
          "label": "2025-04-28",
          "total_revenue": 4425.23,
          "mean_revenue": 885.046,
          "median_revenue": 447.34,
          "total_bookings": 5,
          "total_seats_booked": 10
        },
        {
          "label": "2025-04-29",
          "total_revenue": 3845.86,
          "mean_revenue": 549.4085714285714,
          "median_revenue": 489.525,
          "total_bookings": 7,
          "total_seats_booked": 15
        }
      ]
    }
  }
  ```

- **Success Response (200 OK) - Monthly Timeframe**:
  ```json
  {
    "message": "Revenue data fetched successfully",
    "request_id": "5a3f8c9d-73d1-4bcc-9e6d-f89c4f1e3708",
    "status": "SUCCESS",
    "data": {
      "total_revenue": 8271.09,
      "mean_revenue": 689.2575,
      "median_revenue": 447.34,
      "total_bookings": 12,
      "total_seats_booked": 25,
      "groups": [
        {
          "label": "2025-04",
          "total_revenue": 8271.09,
          "mean_revenue": 689.2575,
          "median_revenue": 447.34,
          "total_bookings": 12,
          "total_seats_booked": 25
        }
      ]
    }
  }
  ```

### Revenue - Movie Filtering

- **URL**: `/revenue?movie_id=tt6823368`
- **Method**: `GET`
- **Authentication**: Required (Admin role)
- **Description**: Returns revenue data filtered by the specified movie ID.

- **Success Response (200 OK)**:
  ```json
  {
    "message": "Revenue data fetched successfully",
    "request_id": "3e7b2c1d-e8a9-4f5b-9c6d-2e1d3f4c5b6a",
    "status": "SUCCESS",
    "data": {
      "total_revenue": 8271.09,
      "mean_revenue": 919.01,
      "median_revenue": 447.34,
      "total_bookings": 9,
      "total_seats_booked": 21,
      "groups": [
        {
          "label": "Glass",
          "total_revenue": 8271.09,
          "mean_revenue": 919.01,
          "median_revenue": 447.34,
          "total_bookings": 9,
          "total_seats_booked": 21
        }
      ]
    }
  }
  ```

### Revenue - Slot Filtering

- **URL**: `/revenue?slot_id=3`
- **Method**: `GET`
- **Authentication**: Required (Admin role)
- **Description**: Returns revenue data filtered by the specified slot ID.

- **Success Response (200 OK)**:
  ```json
  {
    "message": "Revenue data fetched successfully",
    "request_id": "1c2d3e4f-5a6b-7c8d-9e0f-1a2b3c4d5e6f",
    "status": "SUCCESS",
    "data": {
      "total_revenue": 6599.68,
      "mean_revenue": 1099.9466666666667,
      "median_revenue": 489.52500000000003,
      "total_bookings": 6,
      "total_seats_booked": 17,
      "groups": [
        {
          "label": "Evening",
          "total_revenue": 6599.68,
          "mean_revenue": 1099.9466666666667,
          "median_revenue": 489.52500000000003,
          "total_bookings": 6,
          "total_seats_booked": 17
        }
      ]
    }
  }
  ```

### Revenue - Genre Filtering

- **URL**: `/revenue?genre=Crime`
- **Method**: `GET`
- **Authentication**: Required (Admin role)
- **Description**: Returns revenue data filtered by the specified genre.

- **Success Response (200 OK)**:
  ```json
  {
    "message": "Revenue data fetched successfully",
    "request_id": "9a8b7c6d-5e4f-3g2h-1i0j-9k8l7m6n5o4p",
    "status": "SUCCESS",
    "data": {
      "total_revenue": 245.21,
      "mean_revenue": 245.21,
      "median_revenue": 245.21,
      "total_bookings": 1,
      "total_seats_booked": 1,
      "groups": [
        {
          "label": "Crime",
          "total_revenue": 245.21,
          "mean_revenue": 245.21,
          "median_revenue": 245.21,
          "total_bookings": 1,
          "total_seats_booked": 1
        }
      ]
    }
  }
  ```

### Revenue - Month and Year Filtering

- **URL**: `/revenue?month=4`
- **Method**: `GET`
- **Authentication**: Required (Admin role)
- **Description**: Returns revenue data filtered by the specified month (1-12).

- **Success Response (200 OK)**:
  ```json
  {
    "message": "Revenue data fetched successfully",
    "request_id": "2d3e4f5g-6h7i-8j9k-0l1m-2n3o4p5q6r7s",
    "status": "SUCCESS",
    "data": {
      "total_revenue": 8271.09,
      "mean_revenue": 689.2575,
      "median_revenue": 447.34,
      "total_bookings": 12,
      "total_seats_booked": 25,
      "groups": [
        {
          "label": "April",
          "total_revenue": 8271.09,
          "mean_revenue": 689.2575,
          "median_revenue": 447.34,
          "total_bookings": 12,
          "total_seats_booked": 25
        }
      ]
    }
  }
  ```

- **URL**: `/revenue?year=2025`
- **Method**: `GET`
- **Authentication**: Required (Admin role)
- **Description**: Returns revenue data filtered by the specified year.

- **Success Response (200 OK)**:
  ```json
  {
    "message": "Revenue data fetched successfully",
    "request_id": "3e4f5g6h-7i8j-9k0l-1m2n-3o4p5q6r7s8t",
    "status": "SUCCESS",
    "data": {
      "total_revenue": 8271.09,
      "mean_revenue": 689.2575,
      "median_revenue": 447.34,
      "total_bookings": 12,
      "total_seats_booked": 25,
      "groups": [
        {
          "label": "2025",
          "total_revenue": 8271.09,
          "mean_revenue": 689.2575,
          "median_revenue": 447.34,
          "total_bookings": 12,
          "total_seats_booked": 25
        }
      ]
    }
  }
  ```

- **URL**: `/revenue?month=4&year=2025`
- **Method**: `GET`
- **Authentication**: Required (Admin role)
- **Description**: Returns revenue data filtered by both month and year.

- **Success Response (200 OK)**:
  ```json
  {
    "message": "Revenue data fetched successfully",
    "request_id": "4f5g6h7i-8j9k-0l1m-2n3o-4p5q6r7s8t9u",
    "status": "SUCCESS",
    "data": {
      "total_revenue": 8271.09,
      "mean_revenue": 689.2575,
      "median_revenue": 447.34,
      "total_bookings": 12,
      "total_seats_booked": 25,
      "groups": [
        {
          "label": "April;2025",
          "total_revenue": 8271.09,
          "mean_revenue": 689.2575,
          "median_revenue": 447.34,
          "total_bookings": 12,
          "total_seats_booked": 25
        }
      ]
    }
  }
  ```

### Revenue - Combined Filters (Movie and Timeframe)

- **URL**: `/revenue?movie_id=tt6823368&timeframe=monthly`
- **Method**: `GET`
- **Authentication**: Required (Admin role)
- **Description**: Returns revenue data for a specific movie, grouped by month.

- **Success Response (200 OK)**:
  ```json
  {
    "message": "Revenue data fetched successfully",
    "request_id": "b43fe9b8-328a-4f57-a4a5-a219b2a6c239",
    "status": "SUCCESS",
    "data": {
      "total_revenue": 8271.09,
      "mean_revenue": 919.01,
      "median_revenue": 447.34,
      "total_bookings": 9,
      "total_seats_booked": 21,
      "groups": [
        {
          "label": "Glass;2025-04",
          "total_revenue": 8271.09,
          "mean_revenue": 919.01,
          "median_revenue": 447.34,
          "total_bookings": 9,
          "total_seats_booked": 21
        }
      ]
    }
  }
  ```

### Revenue - Combined Filters (Slot and Timeframe)

- **URL**: `/revenue?slot_id=3&timeframe=monthly`
- **Method**: `GET`
- **Authentication**: Required (Admin role)
- **Description**: Returns revenue data for a specific slot, grouped by month.

- **Success Response (200 OK)**:
  ```json
  {
    "message": "Revenue data fetched successfully",
    "request_id": "c4d03d87-d388-4b70-9f56-7da3705d2c8e",
    "status": "SUCCESS",
    "data": {
      "total_revenue": 6599.68,
      "mean_revenue": 1099.9466666666667,
      "median_revenue": 489.52500000000003,
      "total_bookings": 6,
      "total_seats_booked": 17,
      "groups": [
        {
          "label": "Evening;2025-04",
          "total_revenue": 6599.68,
          "mean_revenue": 1099.9466666666667,
          "median_revenue": 489.52500000000003,
          "total_bookings": 6,
          "total_seats_booked": 17
        }
      ]
    }
  }
  ```

### Revenue - Parameter Order Variation

- **URL**: `/revenue?timeframe=monthly&slot_id=3`
- **Method**: `GET`
- **Authentication**: Required (Admin role)
- **Description**: Same filters as previous example but in different order, affecting the response label order.

- **Success Response (200 OK)**:
  ```json
  {
    "message": "Revenue data fetched successfully",
    "request_id": "1eabcde0-e42a-4043-b33b-7a8358a3d775",
    "status": "SUCCESS",
    "data": {
      "total_revenue": 6599.68,
      "mean_revenue": 1099.9466666666667,
      "median_revenue": 489.52500000000003,
      "total_bookings": 6,
      "total_seats_booked": 17,
      "groups": [
        {
          "label": "2025-04;Evening",
          "total_revenue": 6599.68,
          "mean_revenue": 1099.9466666666667,
          "median_revenue": 489.52500000000003,
          "total_bookings": 6,
          "total_seats_booked": 17
        }
      ]
    }
  }
  ```

### Revenue - Combined Filters (Genre and Timeframe)

- **URL**: `/revenue?genre=Crime&timeframe=yearly`
- **Method**: `GET`
- **Authentication**: Required (Admin role)
- **Description**: Returns revenue data for a specific genre, grouped by year.

- **Success Response (200 OK)**:
  ```json
  {
    "message": "Revenue data fetched successfully",
    "request_id": "6feb3104-5182-45af-add3-d7127925a94b",
    "status": "SUCCESS",
    "data": {
      "total_revenue": 245.21,
      "mean_revenue": 245.21,
      "median_revenue": 245.21,
      "total_bookings": 1,
      "total_seats_booked": 1,
      "groups": [
        {
          "label": "Crime;2025",
          "total_revenue": 245.21,
          "mean_revenue": 245.21,
          "median_revenue": 245.21,
          "total_bookings": 1,
          "total_seats_booked": 1
        }
      ]
    }
  }
  ```

### Revenue - No Matching Data

- **URL**: `/revenue?movie_id=nonexistent`
- **Method**: `GET`
- **Authentication**: Required (Admin role)
- **Description**: Returns empty results when no data matches the filter criteria.

- **Success Response (200 OK)**:
  ```json
  {
    "message": "Revenue data fetched successfully",
    "request_id": "5f6g7h8i-9j0k-1l2m-3n4o-5p6q7r8s9t0u",
    "status": "SUCCESS",
    "data": {
      "total_revenue": 0,
      "mean_revenue": 0,
      "median_revenue": 0,
      "total_bookings": 0,
      "total_seats_booked": 0,
      "groups": []
    }
  }
  ```

### Revenue - Mutually Exclusive Parameters (Error)

- **URL**: `/revenue?timeframe=daily&month=4`
- **Method**: `GET`
- **Authentication**: Required (Admin role)
- **Description**: Returns an error when mutually exclusive parameters are combined.

- **Error Response (400 Bad Request)**:
  ```json
  {
    "status": "ERROR",
    "code": "INVALID_PARAMS",
    "message": "Timeframe cannot be combined with month or year filters",
    "request_id": "6g7h8i9j-0k1l-2m3n-4o5p-6q7r8s9t0u1v"
  }
  ```

### Revenue - Invalid Parameter Values (Error)

- **URL**: `/revenue?timeframe=invalid`
- **Method**: `GET`
- **Authentication**: Required (Admin role)
- **Description**: Returns an error for invalid parameter values.

- **Error Response (400 Bad Request)**:
  ```json
  {
    "status": "ERROR",
    "code": "INVALID_PARAMS",
    "message": "Invalid query parameters",
    "request_id": "7h8i9j0k-1l2m-3n4o-5p6q-7r8s9t0u1v2w"
  }
  ```

- **URL**: `/revenue?month=13`
- **Method**: `GET`
- **Authentication**: Required (Admin role)
- **Description**: Returns an error for month value outside valid range.

- **Error Response (400 Bad Request)**:
  ```json
  {
    "status": "ERROR",
    "code": "INVALID_PARAMS",
    "message": "Invalid query parameters",
    "request_id": "8i9j0k1l-2m3n-4o5p-6q7r-8s9t0u1v2w3x"
  }
  ```

### Additional Notes

- The API always returns an aggregate view of statistics (`total_revenue`, `mean_revenue`, `median_revenue`, etc.) regardless of filter combinations.
- The `groups` array contains breakdowns according to your filter and grouping choices.
- When using multiple filter dimensions (e.g., movie + slot + timeframe), data is aggregated according to all filters.
- Label components are always separated by semicolons (`;`) and ordered according to the query parameter order.
- All revenue calculations consider only bookings with `Confirmed` or `CheckedIn` status.

## Dashboard - Booking Data as CSV

The Booking CSV API provides a way to export booking data for analysis or record-keeping purposes. This endpoint allows administrators to download booking information as a CSV file, with optional filtering by month and year.

### How to Use Query Parameters

The Booking CSV API supports the following optional query parameters for filtering:

- **Month Filter**: Filter by specific month (`month=1-12`)
- **Year Filter**: Filter by specific year (`year=YYYY`)

When no parameters are provided, all booking data is exported.

### Important Notes

1. **Authentication**: This API requires admin authentication
2. **Content Type**: Unlike other APIs that return JSON, this endpoint returns a CSV file with appropriate headers
3. **Response Format**: The response is a downloadable CSV file, not a JSON object
4. **Data Scope**: Only "Confirmed" and "CheckedIn" bookings are included in the export

### Download Bookings as CSV

- **URL**: `/booking-csv`
- **Method**: `GET`
- **Authentication**: Required (Admin role)
- **Query Parameters**:
  - `month` (1-12, optional) - Filter by month
  - `year` (e.g., 2025, optional) - Filter by year
- **Description**: Download a CSV file containing booking data filtered by the specified month and/or year.

#### Response Headers

```
Content-Type: text/csv
Content-Disposition: attachment; filename="bookings.csv"
```

The filename varies based on filters:
- All bookings: `bookings.csv`
- Month filter: `bookings_April.csv`
- Year filter: `bookings_2025.csv`
- Month and year filter: `bookings_April_2025.csv`

#### CSV Content

The CSV file contains the following columns:
```
Booking ID, Show ID, Show Date, Customer Name, Phone Number, Number of Seats, Amount Paid, Payment Type, Booking Time, Status
```

#### Example CSV Content

```csv
Booking ID,Show ID,Show Date,Customer Name,Phone Number,Number of Seats,Amount Paid,Payment Type,Booking Time,Status
54,25,2025-04-27,John Smith,9876543210,2,297.34,Card,2025-04-26 23:56:20,Confirmed
11,22,2025-04-26,Jane Doe,8765432109,3,553.30,Card,2025-04-23 16:57:41,Confirmed
```

#### Error Responses

- **Unauthorized (401)**:
  ```json
  {
    "status": "ERROR",
    "code": "INVALID_TOKEN",
    "message": "Unauthorized",
    "request_id": "761234df-ff53-4bbb-882a-039925807c74"
  }
  ```

- **Forbidden (403)**:
  ```json
  {
    "status": "ERROR",
    "code": "FORBIDDEN",
    "message": "Access denied. Admin role required",
    "request_id": "4ccefb6e-f3bd-464b-8ce4-2f30ee4055f0"
  }
  ```

- **Invalid Month (400)**:
  ```json
  {
    "status": "ERROR",
    "code": "INVALID_MONTH",
    "message": "Month must be a number between 1 and 12",
    "request_id": "33749b29-e38f-4aa3-91f6-336ce9d8fc2a"
  }
  ```

- **Invalid Year (400)**:
  ```json
  {
    "status": "ERROR",
    "code": "INVALID_YEAR",
    "message": "Year must be a valid number",
    "request_id": "446ed540-c84f-4264-b1f6-f3958792cc90"
  }
  ```

### Example Usage

#### Download All Bookings
```
GET /booking-csv
```

#### Download Bookings for April 2025
```
GET /booking-csv?month=4&year=2025
```

#### Download Bookings for 2025
```
GET /booking-csv?year=2025
```

#### Download Bookings for April (All Years)
```
GET /booking-csv?month=4
```
