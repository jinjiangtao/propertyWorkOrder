import requests
import json

base_url = "http://localhost:8080/api"

try:
    print("=== Test 1: Register user ===")
    resp = requests.post(f"{base_url}/register", json={"username": "testuser2024", "password": "123456"})
    print(f"Status: {resp.status_code}")
    print(f"Response: {resp.text}")

    print("\n=== Test 2: Login user ===")
    resp = requests.post(f"{base_url}/login", json={"username": "testuser2024", "password": "123456"})
    print(f"Status: {resp.status_code}")
    print(f"Response: {resp.text}")
    if resp.status_code == 200:
        data = resp.json()
        user_id = data.get("user_id")

        print("\n=== Test 3: Create repair ===")
        resp = requests.post(f"{base_url}/repair/create", json={
            "user_id": user_id,
            "username": "testuser2024",
            "repair_type": "水电维修",
            "description": "水龙头漏水"
        })
        print(f"Status: {resp.status_code}")
        print(f"Response: {resp.text}")

        print("\n=== Test 4: Admin get all repairs ===")
        resp = requests.get(f"{base_url}/repair/all")
        print(f"Status: {resp.status_code}")
        print(f"Response: {resp.text}")

        print("\n=== Test 5: Create worker ===")
        resp = requests.post(f"{base_url}/worker/create", json={
            "work_no": "W002",
            "name": "Worker2",
            "phone": "13800138002",
            "password": "123456",
            "skill_type": "水电"
        })
        print(f"Status: {resp.status_code}")
        print(f"Response: {resp.text}")

        print("\n=== Test 6: Assign worker to repair ===")
        resp = requests.post(f"{base_url}/repair/assign", json={
            "repair_id": 1,
            "worker_id": 1
        })
        print(f"Status: {resp.status_code}")
        print(f"Response: {resp.text}")

        print("\n=== Test 7: Worker get repairs ===")
        resp = requests.get(f"{base_url}/repair/worker?worker_id=1")
        print(f"Status: {resp.status_code}")
        print(f"Response: {resp.text}")

        print("\n=== All tests completed ===")

except Exception as e:
    print(f"Error: {e}")