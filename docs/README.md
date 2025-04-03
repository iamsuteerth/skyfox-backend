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

## Field Validation Rules

- **Name**: Must be 3-70 characters, max 4 words, letters only, no consecutive spaces
- **Username**: Must be 3-30 characters, lowercase, no spaces, cannot start with a number, no consecutive special characters
- **Password**: Must be at least 8 characters with at least one uppercase letter and one special character
- **Phone Number**: Must be exactly 10 digits
- **Email**: Must be in valid email format
- **Security Answer**: Must be at least 3 characters long
