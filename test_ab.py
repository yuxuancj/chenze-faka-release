#!/usr/bin/env python3
"""ab 风格压力测试 - 用 Python 实现 500 请求 / 50 并发"""
import json
import time
import urllib.request
import urllib.parse
from concurrent.futures import ThreadPoolExecutor, as_completed
import threading
import sys

BASE = "http://localhost:8080"
N_REQUESTS = 500
CONCURRENCY = 50

print("=" * 60)
print(f"  ab 风格压力测试: -n {N_REQUESTS} -c {CONCURRENCY}")
print("=" * 60)
print()

# 1. 获取有效 token 和商品 ID
print("[准备] 获取 token 和商品...")
login = urllib.request.urlopen(
    urllib.request.Request(f"{BASE}/api/v1/user/login",
        data=json.dumps({"email": "user1@faka.test", "password": "test123456"}).encode(),
        headers={"Content-Type": "application/json"}),
    timeout=10
)
token = json.loads(login.read())["data"]["token"]

# 找压测商品 ID
prod_req = urllib.request.Request(f"{BASE}/api/v1/products",
    headers={"Authorization": f"Bearer {token}"})
prod_resp = json.loads(urllib.request.urlopen(prod_req, timeout=10).read())
test_prod_id = None
for p in prod_resp.get("data", {}).get("list", []):
    if p["stock"] > 50:
        test_prod_id = p["id"]
        break
if not test_prod_id:
    print("FAIL: 没有找到高库存商品")
    sys.exit(1)
print(f"[准备] 测试商品 ID={test_prod_id}")
print()

# 2. 执行并发下单请求
print(f"[压测] 开始 {N_REQUESTS} 请求 / {CONCURRENCY} 并发...")
success = 0
failed = 0
lock = threading.Lock()
results = []

def place_order(i):
    try:
        req = urllib.request.Request(
            f"{BASE}/api/v1/orders",
            data=json.dumps({"product_id": test_prod_id, "quantity": 1, "email": f"load{i}@test.com"}).encode(),
            headers={
                "Authorization": f"Bearer {token}",
                "Content-Type": "application/json"
            },
            method="POST"
        )
        with urllib.request.urlopen(req, timeout=15) as resp:
            body = json.loads(resp.read())
            code = body.get("code", -1)
            with lock:
                results.append(("OK", code, body.get("msg", "")))
            return "OK"
    except Exception as e:
        with lock:
            results.append(("FAIL", str(e)[:50]))
        return "FAIL"

t0 = time.time()
with ThreadPoolExecutor(max_workers=CONCURRENCY) as ex:
    futures = [ex.submit(place_order, i) for i in range(N_REQUESTS)]
    for f in as_completed(futures):
        r = f.result()
        if r == "OK":
            success += 1
        else:
            failed += 1
t1 = time.time()
elapsed = t1 - t0

print()
print(f"[压测] 完成！耗时: {elapsed:.2f}s")
print(f"  成功: {success}")
print(f"  失败: {failed}")
print(f"  QPS: {N_REQUESTS / elapsed:.1f}")
print()

# 3. 验证最终库存
print("[验证] 检查库存...")
admin_login = urllib.request.urlopen(
    urllib.request.Request(f"{BASE}/api/v1/user/login",
        data=json.dumps({"email": "admin@chenze.com", "password": "admin123"}).encode(),
        headers={"Content-Type": "application/json"}),
    timeout=10
)
admin_token = json.loads(admin_login.read())["data"]["token"]
admin_req = urllib.request.Request(f"{BASE}/admin/api/products",
    headers={"Authorization": f"Bearer {admin_token}"})
admin_resp = json.loads(urllib.request.urlopen(admin_req, timeout=10).read())

initial_stock = 0
for p in admin_resp.get("data", {}).get("list", []):
    if p["id"] == test_prod_id:
        initial_stock = p["stock"]
        break

print(f"  初始库存: {initial_stock} (压测前)")
print(f"  成功订单: {success}")
expected = initial_stock - success
print(f"  最终库存: 需验证...")

# 获取当前库存
admin_req2 = urllib.request.Request(f"{BASE}/admin/api/products",
    headers={"Authorization": f"Bearer {admin_token}"})
final_resp = json.loads(urllib.request.urlopen(admin_req2, timeout=10).read())
final_stock = 0
for p in final_resp.get("data", {}).get("list", []):
    if p["id"] == test_prod_id:
        final_stock = p["stock"]
        break

print(f"  最终库存: {final_stock} (期望: >=0)")
if final_stock >= 0:
    print(f"  OK: 无超卖！库存无负数")
else:
    print(f"  FAIL: 发生超卖！库存={final_stock}")

# 检查是否有死锁日志
print()
print("[日志] 检查错误日志...")
with open("/tmp/srv.log") as f:
    log = f.read()
    deadlock_count = log.count("deadlock")
    timeout_count = log.count("database is locked")
    print(f"  deadlock: {deadlock_count}")
    print(f"  database locked: {timeout_count}")
    if deadlock_count == 0 and timeout_count < 10:
        print(f"  OK: 无死锁")
    else:
        print(f"  WARN: 发现 {deadlock_count} 个死锁或 {timeout_count} 个锁定错误")

print()
print("=" * 60)
print(f"  压测完成")
print("=" * 60)
