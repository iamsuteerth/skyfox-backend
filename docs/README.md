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