import base64
import hashlib
import requests
import json
import time

# Base URL for the API
BASE_URL = "http://localhost:8080"

# JWT tokens
CUSTOMER_JWT = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NDQ4NzEzMTcsInJvbGUiOiJjdXN0b21lciIsInVzZXJuYW1lIjoic3V0ZWVydGgifQ.udEDPpw1qJSrnvn1iLol7S2b-OAUAkrBEC1QoSIS-H4"
STAFF_JWT = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NDQ4NzEyNDYsInJvbGUiOiJzdGFmZiIsInVzZXJuYW1lIjoic3RhZmYtMSJ9.OFa0OLMRz3BtpvNQWI4KrR2CGRA0tZ_ug_yawL_SdvY"

# Helper function to print response
def print_response(test_name, response):
    print(f"\n=== {test_name} ===")
    print(f"Status Code: {response.status_code}")
    print("Response Body:")
    try:
        print(json.dumps(response.json(), indent=4))
    except:
        print(response.text)
    print("=" * (len(test_name) + 8))

# Test 1: RBAC Check - Using Staff JWT (should fail)
def test_rbac():
    headers = {
        "Authorization": f"Bearer {STAFF_JWT}",
        "Content-Type": "application/json"
    }
    
    payload = {
        "security_answer": "test_answer",
        "profile_img": "test_base64",
        "profile_img_sha": "test_sha"
    }
    
    response = requests.post(
        f"{BASE_URL}/customer/update-profile-image",
        headers=headers,
        json=payload
    )
    
    print_response("RBAC Test (Staff JWT)", response)
    return response.status_code == 403  # Should be forbidden

# Test 2: Validation Errors - Missing Fields
def test_validation_missing_fields():
    headers = {
        "Authorization": f"Bearer {CUSTOMER_JWT}",
        "Content-Type": "application/json"
    }
    
    payload = {
        # Missing all required fields
    }
    
    response = requests.post(
        f"{BASE_URL}/customer/update-profile-image",
        headers=headers,
        json=payload
    )
    
    print_response("Validation Test (Missing Fields)", response)
    return response.status_code == 400  # Should be bad request

# Test 3: Invalid Base64 Image
def test_invalid_base64():
    headers = {
        "Authorization": f"Bearer {CUSTOMER_JWT}",
        "Content-Type": "application/json"
    }
    
    payload = {
        "security_answer": "your_security_answer",  # Replace with actual answer
        "profile_img": "not_valid_base64!@#",
        "profile_img_sha": "invalid_sha"
    }
    
    response = requests.post(
        f"{BASE_URL}/customer/update-profile-image",
        headers=headers,
        json=payload
    )
    
    print_response("Invalid Base64 Test", response)
    return response.status_code == 400  # Should be bad request

# Test 4: Invalid Security Answer
def test_invalid_security_answer():
    # Read the image file
    with open('test_image.jpg', 'rb') as image_file:
        image_bytes = image_file.read()

    # Convert image to base64
    base64_encoded = base64.b64encode(image_bytes).decode('utf-8')

    # Calculate SHA-256 hash
    sha256_hash = hashlib.sha256(image_bytes).hexdigest()
    
    headers = {
        "Authorization": f"Bearer {CUSTOMER_JWT}",
        "Content-Type": "application/json"
    }
    
    payload = {
        "security_answer": "wrong_answer",  # Deliberately wrong answer
        "profile_img": base64_encoded,
        "profile_img_sha": sha256_hash
    }
    
    response = requests.post(
        f"{BASE_URL}/customer/update-profile-image",
        headers=headers,
        json=payload
    )
    
    print_response("Invalid Security Answer Test", response)
    return response.status_code == 400  # Should be bad request

# Test 5: Successful Update
def test_successful_update():
    # Read the image file
    with open('test_image.jpg', 'rb') as image_file:
        image_bytes = image_file.read()

    # Convert image to base64
    base64_encoded = base64.b64encode(image_bytes).decode('utf-8')

    # Calculate SHA-256 hash
    sha256_hash = hashlib.sha256(image_bytes).hexdigest()
    
    headers = {
        "Authorization": f"Bearer {CUSTOMER_JWT}",
        "Content-Type": "application/json"
    }
    
    payload = {
        "security_answer": "blr",  # Replace with actual answer
        "profile_img": base64_encoded,
        "profile_img_sha": sha256_hash
    }
    
    response = requests.post(
        f"{BASE_URL}/customer/update-profile-image",
        headers=headers,
        json=payload
    )
    
    print_response("Successful Update Test", response)
    return response.status_code == 200  # Should be success

# Run all tests
def run_tests():
    print("Starting API tests for update-profile-image endpoint...")
    
    tests = [
        ("RBAC Check", test_rbac),
        ("Validation - Missing Fields", test_validation_missing_fields),
        ("Invalid Base64", test_invalid_base64),
        ("Invalid Security Answer", test_invalid_security_answer),
        ("Successful Update", test_successful_update)
    ]
    
    results = {}
    for name, test_func in tests:
        print(f"\nRunning test: {name}")
        try:
            result = test_func()
            results[name] = "✅ PASSED" if result else "❌ FAILED"
        except Exception as e:
            print(f"Error during test: {str(e)}")
            results[name] = f"❌ ERROR: {str(e)}"
    
    # Print summary
    print("\n=== TEST SUMMARY ===")
    for name, result in results.items():
        print(f"{name}: {result}")

# For successful test, you'll need to prepare:
# 1. A valid customer JWT token
# 2. The actual security answer for the user
# 3. A test image file named 'test_image.jpg'

if __name__ == "__main__":
    run_tests()
