import subprocess
import sys

try:
    import requests
except ImportError:
    subprocess.check_call([sys.executable, "-m", "pip", "install", "requests"])
    import requests

base_url = "http://localhost:8080/api"

def test_register():
    print("=== Test 1: 用户注册 ===")
    resp = requests.post(f"{base_url}/register", json={"username": "testuser_new", "password": "123456"})
    print(f"状态码: {resp.status_code}")
    print(f"响应: {resp.text}")
    return resp.status_code == 200

def test_login():
    print("\n=== Test 2: 用户登录 ===")
    resp = requests.post(f"{base_url}/login", json={"username": "testuser_new", "password": "123456"})
    print(f"状态码: {resp.status_code}")
    print(f"响应: {resp.text}")
    if resp.status_code == 200:
        return resp.json().get("user_id")
    return None

def test_create_repair(user_id):
    print("\n=== Test 3: 提交报修 ===")
    resp = requests.post(f"{base_url}/repair/create", json={
        "user_id": user_id,
        "username": "testuser_new",
        "repair_type": "水电维修",
        "description": "水龙头漏水测试"
    })
    print(f"状态码: {resp.status_code}")
    print(f"响应: {resp.text}")
    return resp.status_code == 200

def test_get_all_repairs():
    print("\n=== Test 4: Admin查看所有报修 ===")
    resp = requests.get(f"{base_url}/repair/all")
    print(f"状态码: {resp.status_code}")
    print(f"响应: {resp.text}")
    return resp.status_code == 200

def test_create_worker():
    print("\n=== Test 5: 创建工人 ===")
    resp = requests.post(f"{base_url}/worker/create", json={
        "work_no": "W003",
        "name": "测试工人",
        "phone": "13800138003",
        "password": "123456",
        "skill_type": "水电"
    })
    print(f"状态码: {resp.status_code}")
    print(f"响应: {resp.text}")
    return resp.status_code == 200

def test_assign_worker():
    print("\n=== Test 6: 派单给工人 ===")
    resp = requests.post(f"{base_url}/repair/assign", json={
        "repair_id": 1,
        "worker_id": 1
    })
    print(f"状态码: {resp.status_code}")
    print(f"响应: {resp.text}")
    return resp.status_code == 200

def test_worker_get_repairs():
    print("\n=== Test 7: 工人查看工单 ===")
    resp = requests.get(f"{base_url}/repair/worker?worker_id=1")
    print(f"状态码: {resp.status_code}")
    print(f"响应: {resp.text}")
    return resp.status_code == 200

def test_admin_login():
    print("\n=== Test 8: Admin登录 ===")
    resp = requests.post(f"{base_url}/login", json={"username": "admin", "password": "123456"})
    print(f"状态码: {resp.status_code}")
    print(f"响应: {resp.text}")
    return resp.status_code == 200

def test_worker_login():
    print("\n=== Test 9: 工人登录 ===")
    resp = requests.post(f"{base_url}/worker/login", json={"work_no": "W003", "password": "123456"})
    print(f"状态码: {resp.status_code}")
    print(f"响应: {resp.text}")
    return resp.status_code == 200

if __name__ == "__main__":
    try:
        results = []
        
        results.append(("用户注册", test_register()))
        user_id = test_login()
        results.append(("用户登录", user_id is not None))
        
        if user_id:
            results.append(("提交报修", test_create_repair(user_id)))
        
        results.append(("Admin查看报修", test_get_all_repairs()))
        results.append(("创建工人", test_create_worker()))
        results.append(("派单", test_assign_worker()))
        results.append(("工人查看工单", test_worker_get_repairs()))
        results.append(("Admin登录", test_admin_login()))
        results.append(("工人登录", test_worker_login()))
        
        print("\n" + "="*50)
        print("测试结果汇总:")
        print("="*50)
        all_pass = True
        for test_name, passed in results:
            status = "✅ 通过" if passed else "❌ 失败"
            print(f"{test_name}: {status}")
            if not passed:
                all_pass = False
        
        print("="*50)
        if all_pass:
            print("🎉 所有测试通过!")
        else:
            print("⚠️ 部分测试失败，请检查")
            
    except Exception as e:
        print(f"测试过程出错: {e}")
        import traceback
        traceback.print_exc()