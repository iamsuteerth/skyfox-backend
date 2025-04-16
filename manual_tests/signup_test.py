import base64
import hashlib
import requests
import json

def test_customer_signup():
    print("Testing Customer Signup API")
    
    try:
        # Read the image file
        with open('test_image.jpg', 'rb') as image_file:
            image_bytes = image_file.read()

        # Convert image to base64
        base64_encoded = base64.b64encode(image_bytes).decode('utf-8')

        # Calculate SHA-256 hash of the raw image bytes
        sha256_hash = hashlib.sha256(image_bytes).hexdigest()

        print("\nImage prepared successfully:")
        print(f"SHA-256 Hash: {sha256_hash[:10]}...{sha256_hash[-10:]}")
        
        # Test Case 1: Valid Signup Request
        print("\n=== Test Case 1: Valid Signup Request ===")
        
        # Prepare the request payload with valid data
        payload = {
            "name": "John Doe",
            "username": f"johndoe{int(time.time())}",  # Use timestamp to avoid conflicts
            "password": "SecurePass@123",
            "number": "9876543210",
            "email": f"john.doe123@example.com",  # Use timestamp to avoid conflicts
            "profile_img": base64_encoded,
            "profile_img_sha": sha256_hash,
            "security_question_id": 1,
            "security_answer": "Answer123"
        }

        # Send the request
        response = requests.post('http://localhost:8080/customer/signup', json=payload)

        # Print the response
        print(f"Status Code: {response.status_code}")
        print("Response Body:")
        print(json.dumps(response.json(), indent=4))
        
        # Test Case 2: Validation Error - Missing Required Fields
        print("\n=== Test Case 2: Validation Error - Missing Required Fields ===")
        
        # Prepare payload with missing fields
        invalid_payload = {
            "name": "John Doe",
            # Missing username
            "password": "SecurePass@123",
            # Missing number
            # Missing email
            # Missing profile_img
            # Missing profile_img_sha
            "security_question_id": 1,
            "security_answer": "Answer123"
        }
        
        # Send the request
        response = requests.post('http://localhost:8080/customer/signup', json=invalid_payload)
        
        # Print the response
        print(f"Status Code: {response.status_code}")
        print("Response Body:")
        print(json.dumps(response.json(), indent=4))
        
        # Test Case 3: Validation Error - Invalid Format
        print("\n=== Test Case 3: Validation Error - Invalid Format ===")
        
        # Prepare payload with invalid format
        invalid_format_payload = {
            "name": "123", # Numbers instead of letters
            "username": "j", # Too short
            "password": "weak", # Too weak
            "number": "123", # Not 10 digits
            "email": "not-an-email", # Invalid email format
            "profile_img": base64_encoded,
            "profile_img_sha": sha256_hash,
            "security_question_id": 1,
            "security_answer": "A" # Too short
        }
        
        # Send the request
        response = requests.post('http://localhost:8080/customer/signup', json=invalid_format_payload)
        
        # Print the response
        print(f"Status Code: {response.status_code}")
        print("Response Body:")
        print(json.dumps(response.json(), indent=4))
        
        # Test Case 4: Duplicate Username/Email
        print("\n=== Test Case 4: Duplicate Username/Email ===")
        
        # Use the same payload as the first test (assuming it was successful)
        # Send the request again to trigger duplicate error
        response = requests.post('http://localhost:8080/customer/signup', json=payload)
        
        # Print the response
        print(f"Status Code: {response.status_code}")
        print("Response Body:")
        print(json.dumps(response.json(), indent=4))
        
    except Exception as e:
        print(f"An error occurred: {str(e)}")

# For successful test, you'll need to prepare:
# 1. A test image file named 'test_image.jpg'

if __name__ == "__main__":
    import time
    test_customer_signup()
