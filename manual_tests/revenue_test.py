import requests
import time
import json
from datetime import datetime

# Configuration
BASE_URL = "http://localhost:8080"
JWT_TOKEN = ".."
API_GATEWAY_KEY = input("Enter your API Gateway Key: ")
REVENUE_ENDPOINT = "/revenue"
REQUEST_DELAY = 3  # seconds

# Headers
headers = {
    "Authorization": f"Bearer {JWT_TOKEN}",
    "X-Api-Key": API_GATEWAY_KEY
}

# Test cases with thorough coverage
test_cases = [
    # Base cases
    {"name": "Default, all data", "url": "/revenue"},
    {"name": "Explicit all timeframe", "url": "/revenue?timeframe=all"},
    
    # Timeframe variants
    {"name": "Daily timeframe", "url": "/revenue?timeframe=daily"},
    {"name": "Weekly timeframe", "url": "/revenue?timeframe=weekly"},
    {"name": "Monthly timeframe", "url": "/revenue?timeframe=monthly"},
    {"name": "Yearly timeframe", "url": "/revenue?timeframe=yearly"},
    
    # Single filters (movie, slot, genre)
    {"name": "Filter by movie ID", "url": "/revenue?movie_id=tt6823368"},
    {"name": "Filter by slot ID (Morning)", "url": "/revenue?slot_id=1"},
    {"name": "Filter by slot ID (Afternoon)", "url": "/revenue?slot_id=2"},
    {"name": "Filter by slot ID (Evening)", "url": "/revenue?slot_id=3"},
    {"name": "Filter by Crime genre", "url": "/revenue?genre=Crime"},
    {"name": "Filter by Drama genre", "url": "/revenue?genre=Drama"},
    {"name": "Filter by Action genre", "url": "/revenue?genre=Action"},
    
    # Month and Year filters
    {"name": "Filter by April", "url": "/revenue?month=4"},
    {"name": "Filter by May", "url": "/revenue?month=5"},
    {"name": "Filter by year 2025", "url": "/revenue?year=2025"},
    {"name": "Filter by April 2025", "url": "/revenue?month=4&year=2025"},
    {"name": "Filter by May 2025", "url": "/revenue?month=5&year=2025"},
    
    # Movie + Timeframe combinations
    {"name": "Movie with daily timeframe", "url": "/revenue?movie_id=tt6823368&timeframe=daily"},
    {"name": "Movie with weekly timeframe", "url": "/revenue?movie_id=tt6823368&timeframe=weekly"},
    {"name": "Movie with monthly timeframe", "url": "/revenue?movie_id=tt6823368&timeframe=monthly"},
    {"name": "Movie with yearly timeframe", "url": "/revenue?movie_id=tt6823368&timeframe=yearly"},
    
    # Slot + Timeframe combinations
    {"name": "Slot with daily timeframe", "url": "/revenue?slot_id=3&timeframe=daily"},
    {"name": "Slot with weekly timeframe", "url": "/revenue?slot_id=3&timeframe=weekly"},
    {"name": "Slot with monthly timeframe", "url": "/revenue?slot_id=3&timeframe=monthly"},
    {"name": "Slot with yearly timeframe", "url": "/revenue?slot_id=3&timeframe=yearly"},
    
    # Genre + Timeframe combinations
    {"name": "Genre with daily timeframe", "url": "/revenue?genre=Crime&timeframe=daily"},
    {"name": "Genre with weekly timeframe", "url": "/revenue?genre=Crime&timeframe=weekly"},
    {"name": "Genre with monthly timeframe", "url": "/revenue?genre=Crime&timeframe=monthly"},
    {"name": "Genre with yearly timeframe", "url": "/revenue?genre=Crime&timeframe=yearly"},
    
    # Double filter combinations
    {"name": "Movie + slot", "url": "/revenue?movie_id=tt6823368&slot_id=3"},
    {"name": "Movie + genre", "url": "/revenue?movie_id=tt6823368&genre=Crime"},
    {"name": "Slot + genre", "url": "/revenue?slot_id=3&genre=Crime"},
    {"name": "Movie + month", "url": "/revenue?movie_id=tt6823368&month=4"},
    {"name": "Movie + year", "url": "/revenue?movie_id=tt6823368&year=2025"},
    {"name": "Slot + month", "url": "/revenue?slot_id=3&month=4"},
    {"name": "Slot + year", "url": "/revenue?slot_id=3&year=2025"},
    {"name": "Genre + month", "url": "/revenue?genre=Crime&month=4"},
    {"name": "Genre + year", "url": "/revenue?genre=Crime&year=2025"},
    
    # Triple filter combinations
    {"name": "Movie + slot + timeframe", "url": "/revenue?movie_id=tt6823368&slot_id=3&timeframe=monthly"},
    {"name": "Movie + genre + timeframe", "url": "/revenue?movie_id=tt6823368&genre=Crime&timeframe=monthly"},
    {"name": "Slot + genre + timeframe", "url": "/revenue?slot_id=3&genre=Crime&timeframe=monthly"},
    {"name": "Movie + slot + month", "url": "/revenue?movie_id=tt6823368&slot_id=3&month=4"},
    {"name": "Movie + genre + month", "url": "/revenue?movie_id=tt6823368&genre=Crime&month=4"},
    {"name": "Slot + genre + month", "url": "/revenue?slot_id=3&genre=Crime&month=4"},
    {"name": "Movie + slot + year", "url": "/revenue?movie_id=tt6823368&slot_id=3&year=2025"},
    {"name": "Movie + genre + year", "url": "/revenue?movie_id=tt6823368&genre=Crime&year=2025"},
    {"name": "Slot + genre + year", "url": "/revenue?slot_id=3&genre=Crime&year=2025"},
    
    # Quad filter combinations
    {"name": "Movie + slot + month + year", "url": "/revenue?movie_id=tt6823368&slot_id=3&month=4&year=2025"},
    {"name": "Movie + genre + month + year", "url": "/revenue?movie_id=tt6823368&genre=Crime&month=4&year=2025"},
    {"name": "Slot + genre + month + year", "url": "/revenue?slot_id=3&genre=Crime&month=4&year=2025"},
    
    # Parameter order variations
    {"name": "Movie then genre", "url": "/revenue?movie_id=tt6823368&genre=Crime"},
    {"name": "Genre then movie", "url": "/revenue?genre=Crime&movie_id=tt6823368"},
    {"name": "Timeframe then slot", "url": "/revenue?timeframe=monthly&slot_id=3"},
    {"name": "Slot then timeframe", "url": "/revenue?slot_id=3&timeframe=monthly"},
    {"name": "Complex param order", "url": "/revenue?slot_id=3&movie_id=tt6823368&timeframe=daily"},
    {"name": "Month then year", "url": "/revenue?month=4&year=2025"},
    {"name": "Year then month", "url": "/revenue?year=2025&month=4"},
    
    # Invalid parameters (400 errors)
    {"name": "Invalid timeframe", "url": "/revenue?timeframe=invalid"},
    {"name": "Invalid month (too high)", "url": "/revenue?month=13"},
    {"name": "Invalid month (zero)", "url": "/revenue?month=0"},
    {"name": "Non-numeric month", "url": "/revenue?month=abc"},
    
    # Mutually exclusive parameters (400 errors)
    {"name": "Timeframe + month", "url": "/revenue?timeframe=daily&month=4"},
    {"name": "Timeframe + year", "url": "/revenue?timeframe=weekly&year=2025"},
    {"name": "Timeframe + month + year", "url": "/revenue?timeframe=monthly&month=4&year=2025"},
    {"name": "Timeframe + month + movie", "url": "/revenue?timeframe=yearly&month=4&movie_id=tt6823368"},
    {"name": "Timeframe + year + slot", "url": "/revenue?timeframe=daily&year=2025&slot_id=3"},
    {"name": "Timeframe + month + genre", "url": "/revenue?timeframe=weekly&month=4&genre=Crime"},
    
    # Edge cases - non-existent data
    {"name": "Non-existent movie ID", "url": "/revenue?movie_id=nonexistent"},
    {"name": "Non-existent slot ID", "url": "/revenue?slot_id=999"},
    {"name": "Non-existent genre", "url": "/revenue?genre=NonexistentGenre"},
    {"name": "Time period with no data", "url": "/revenue?month=1&year=2020"},
    {"name": "Multiple filters with no data", "url": "/revenue?month=1&year=2020&movie_id=tt6823368"},
    
    # Comprehensive corner cases
    {"name": "Multiple dimensions with valid data", "url": "/revenue?movie_id=tt6823368&slot_id=1&genre=Drama&month=4&year=2025"},
    {"name": "All filters with non-existent value", "url": "/revenue?movie_id=nonexistent&slot_id=999&genre=NonexistentGenre"},
    {"name": "Mixed valid/invalid filters", "url": "/revenue?movie_id=tt6823368&slot_id=999&genre=Drama"}
]

def run_tests():
    results = []
    
    for i, test_case in enumerate(test_cases):
        print(f"Running test {i+1}/{len(test_cases)}: {test_case['name']}")
        
        full_url = BASE_URL + test_case["url"]
        
        try:
            response = requests.get(full_url, headers=headers)
            
            # Format response for readability
            if response.status_code == 200:
                try:
                    response_data = response.json()
                    formatted_json = json.dumps(response_data, indent=4)
                except:
                    formatted_json = response.text
            else:
                formatted_json = response.text
                
            result = {
                "test_name": test_case["name"],
                "url": test_case["url"],
                "status_code": response.status_code,
                "response": formatted_json
            }
            
            # Verify expected response patterns
            if response.status_code == 200:
                try:
                    data = response.json()
                    if "data" in data and "groups" not in data["data"]:
                        print(f"  ❌ ERROR: Missing 'groups' in response for test {i+1}")
                    elif "data" in data and data["data"]["groups"] is None:
                        print(f"  ❌ ERROR: 'groups' is null instead of empty array in test {i+1}")
                except:
                    print(f"  ❌ ERROR: Could not parse JSON response for test {i+1}")
            
            results.append(result)
            
        except Exception as e:
            results.append({
                "test_name": test_case["name"],
                "url": test_case["url"],
                "status_code": "ERROR",
                "response": str(e)
            })
            print(f"  ❌ ERROR: {str(e)}")
        
        # Wait before next request
        if i < len(test_cases) - 1:
            time.sleep(REQUEST_DELAY)
    
    return results

def write_results_to_file(results):
    timestamp = datetime.now().strftime("%Y%m%d_%H%M%S")
    filename = f"revenue_api_test_results_{timestamp}.txt"
    
    with open(filename, "w") as f:
        f.write("REVENUE API TEST RESULTS\n")
        f.write("========================\n\n")
        f.write(f"Test run at: {datetime.now().strftime('%Y-%m-%d %H:%M:%S')}\n\n")
        
        for i, result in enumerate(results):
            f.write(f"TEST #{i+1}: {result['test_name']}\n")
            f.write(f"URL: {result['url']}\n")
            f.write(f"Status Code: {result['status_code']}\n")
            f.write("Response:\n")
            f.write(result['response'])
            f.write("\n\n" + "="*80 + "\n\n")
    
    summary_filename = f"revenue_api_test_summary_{timestamp}.txt"
    
    with open(summary_filename, "w") as f:
        f.write("REVENUE API TEST SUMMARY\n")
        f.write("========================\n\n")
        
        success_count = len([r for r in results if r['status_code'] == 200])
        error_count = len([r for r in results if r['status_code'] != 200])
        
        f.write(f"Total tests: {len(results)}\n")
        f.write(f"Successful tests (200): {success_count}\n")
        f.write(f"Error responses: {error_count}\n\n")
        
        f.write("Error details:\n")
        for i, result in enumerate(results):
            if result['status_code'] != 200:
                f.write(f"  - Test #{i+1}: {result['test_name']} (Status: {result['status_code']})\n")
    
    print(f"Test results written to {filename}")
    print(f"Test summary written to {summary_filename}")

def run_verification_tests():
    """Run specific tests to verify consistency in results"""
    print("\nRunning verification tests...")
    
    # Test parameter order consistency
    order_tests = [
        {"name": "Order test 1", "url": "/revenue?movie_id=tt6823368&genre=Crime"},
        {"name": "Order test 2", "url": "/revenue?genre=Crime&movie_id=tt6823368"}
    ]
    
    responses = []
    for test in order_tests:
        resp = requests.get(BASE_URL + test["url"], headers=headers)
        responses.append(resp.json())
        time.sleep(1)
    
    if responses[0]["data"]["total_revenue"] == responses[1]["data"]["total_revenue"]:
        print("✅ Parameter order filter consistency: PASSED")
    else:
        print("❌ Parameter order filter consistency: FAILED - different results for same filters in different order")
    
    # Test empty results consistency
    print("\nVerifying empty results consistency...")
    resp = requests.get(BASE_URL + "/revenue?movie_id=nonexistent", headers=headers)
    data = resp.json()
    
    if "data" in data and "groups" in data["data"]:
        if data["data"]["groups"] == []:
            print("✅ Empty results handling: PASSED - returns empty array")
        elif data["data"]["groups"] is None:
            print("❌ Empty results handling: FAILED - returns null instead of empty array")
        else:
            print("❌ Empty results handling: FAILED - unexpected groups format")
    else:
        print("❌ Empty results handling: FAILED - missing expected data structure")

if __name__ == "__main__":
    print("Starting Revenue API tests...")
    results = run_tests()
    write_results_to_file(results)
    
    # Run additional verification tests
    run_verification_tests()
    
    print("Testing complete!")
