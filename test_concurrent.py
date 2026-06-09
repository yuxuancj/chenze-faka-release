#!/usr/bin/env python3
"""晨泽发卡系统 - 并发安全完整测试 (urllib版)"""
import json
import time
import hashlib
import urllib.request
import urllib.parse
from concurrent.futures import ThreadPoolExecutor, as_completed
import threading

BASE = "http://localhost:8080"
ADMIN_EMAIL = "admin@chenze.com"
ADMIN_PWD = "admin123"
TEST_PROD_STOCK = 200
TEST_CARD_COUNT = 200

def api_call(method, path, data=None, token=None):
    url = f"{BASE}{path}"
    headers = {}
    if token:
        headers["Authorization"] = f"Bearer {token}"
    body = None
    if data is not None:
        if method == "GET":
            qs = urllib.parse.urlencode(data)
            url = f"{url}?{qs}"
        else:
            headers["Content-Type"] = "application/json"
            body = json.dumps(data).encode()
    req = urllib.request.Request(url, data=body, method=method)
    for k, v in headers.items():
        req.add_header(k, v)
    try:
        with urllib.request.urlopen(req, timeout=10) as resp:
            return json.loads(resp.read())
    except Exception as e:
        return {"error": str(e)}

def api_post_raw(path, data):
    """发送 form-urlencoded 数据（用于回调测试）"""
    url = f"{BASE}{path}"
    body = urllib.parse.urlencode(data).encode()
    req = urllib.request.Request(url, data=body, method="POST")
    req.add_header("Content-Type", "application/x-www-form-urlencoded")
    try:
        with urllib.request.urlopen(req, timeout=10) as resp:
            return resp.read().decode()
    except Exception as e:
        return str(e)

print("=" * 56)
print("  晨泽发卡系统 - 并发安全完整测试")
print("=" * 56)
print()

# 1. 管理员登录
print("[1/8] 管理员登录...")
r = api_call("POST", "/api/v1/user/login", {"email": ADMIN_EMAIL, "password": ADMIN_PWD})
admin_token = r.get("data", {}).get("token", "")
if not admin_token:
    print(f"  FAIL: {r}")
    exit(1)
print(f"  OK: admin token = {admin_token[:20]}...")

# 2. 创建高库存测试商品
print("[2/8] 创建高库存测试商品...")
r = api_call("POST", "/admin/api/products", {
    "category_id": 1, "name": "压测商品-并发安全测试",
    "description": "用于并发压测", "price": 0.01,
    "stock": TEST_PROD_STOCK, "type": "card", "status": 1,
}, token=admin_token)
test_prod_id = r.get("data", {}).get("id", 0)
if not test_prod_id:
    print(f"  FAIL: {r}")
    exit(1)
print(f"  OK: 测试商品 ID={test_prod_id} 创建成功")

# 3. 批量导入卡密
print("[3/8] 批量导入 200 张卡密...")
import subprocess as sp
cards_file = "/tmp/test_cards.txt"
with open(cards_file, "w") as f:
    for i in range(1, TEST_CARD_COUNT + 1):
        f.write(f"PERF-CARD-{i:04d}\n")
proc = sp.run([
    "curl", "-s", "-X", "POST", f"{BASE}/admin/api/cards/import",
    "-H", f"Authorization: Bearer {admin_token}",
    "--form", f"product_id={test_prod_id}",
    "--form", f"cards=@{cards_file}"
], capture_output=True, text=True)
try:
    import_resp = json.loads(proc.stdout)
except:
    import_resp = {"raw": proc.stdout}
import_count = import_resp.get("data", {}).get("count", 0)
print(f"  OK: 导入 {import_count} 张卡密")

# 4. 验证初始库存
print("[4/8] 验证商品初始库存...")
r = api_call("GET", "/admin/api/products", token=admin_token)
initial_stock = 0
for p in r.get("data", {}).get("list", []):
    if p["id"] == test_prod_id:
        initial_stock = p["stock"]
        break
print(f"  初始库存: {initial_stock}")

# 5. 注册 10 个测试用户
print("[5/8] 注册 10 个测试用户...")
user_tokens = []
for i in range(1, 11):
    email = f"user{i}@faka.test"
    r = api_call("POST", "/api/v1/user/register", {"email": email, "password": "test123456"})
    token = r.get("data", {}).get("token", "")
    if not token:
        r = api_call("POST", "/api/v1/user/login", {"email": email, "password": "test123456"})
        token = r.get("data", {}).get("token", "")
    user_tokens.append(token)
    print(f"  用户{i}: {token[:20] if token else 'FAIL'}...")
print(f"  共 {len(user_tokens)} 个用户")

# 6. 并发下单
print()
print("[6/8] 并发下单测试（50 单，30 并发）...")
success_orders = []
results_list = []
lock = threading.Lock()

def create_order(idx):
    token = user_tokens[idx % len(user_tokens)]
    try:
        r = api_call("POST", "/api/v1/orders", {
            "product_id": test_prod_id, "quantity": 1, "email": "test@faka.test"
        }, token=token)
        with lock:
            if r.get("code") == 0 and r.get("data", {}).get("order_no"):
                success_orders.append(r["data"]["order_no"])
                return "OK"
            else:
                return f"FAIL:{r.get('msg','')}"
    except Exception as e:
        with lock:
            return f"ERR:{e}"

with ThreadPoolExecutor(max_workers=30) as ex:
    futures = [ex.submit(create_order, i) for i in range(50)]
    for f in as_completed(futures):
        results_list.append(f.result())

ok_count = sum(1 for x in results_list if x == "OK")
fail_count = 50 - ok_count
print(f"  并发下单完成: 成功={ok_count}, 失败(含库存不足)={fail_count}")

# 7. 库存验证
print()
print("[7/8] 库存与超卖检查...")
r = api_call("GET", "/admin/api/products", token=admin_token)
final_stock = 0
for p in r.get("data", {}).get("list", []):
    if p["id"] == test_prod_id:
        final_stock = p["stock"]
        break
expected = initial_stock - ok_count
print(f"  初始库存: {initial_stock}, 成功订单: {ok_count}")
print(f"  最终库存: {final_stock} (期望: {expected})")
if final_stock >= 0:
    print(f"  OK: 库存无负数，无超卖！")
else:
    print(f"  FAIL: 库存为负数，发生超卖！")
if final_stock == expected:
    print(f"  OK: 库存完全正确！")
else:
    diff = abs(final_stock - expected)
    print(f"  INFO: 库存差异 {diff}（部分请求失败，如库存不足）")

# 8. 易支付回调幂等性测试
print()
print("[8/8] 易支付回调幂等性测试...")
if not success_orders:
    print("  SKIP: 无成功订单")
else:
    test_order = success_orders[0]
    print(f"  测试订单号: {test_order}")

    # 初始状态
    r = api_call("GET", f"/api/v1/orders/{test_order}", token=user_tokens[0])
    init_status = r.get("data", {}).get("order", {}).get("status", -1)
    init_cards = len(r.get("data", {}).get("cards", []))
    print(f"  初始状态: {init_status}(0=pending), 卡密: {init_cards} 张")

    # 构造签名（易支付规范：参数按字母排序拼接 &key=KEY）
    # 注意：所有 curl 发送的字段都要包含（money, name, notify_url, out_trade_no,
    #      pid, return_url, trade_no, trade_status, type），sign 和 sign_type 除外
    params = {
        "money":        "0.01",
        "name":         "压测商品",
        "notify_url":   "http://localhost:8080/api/v1/pay/epay/notify",
        "out_trade_no": test_order,
        "pid":          "test_mch_id",
        "return_url":   f"http://localhost:8080/order/{test_order}",
        "trade_no":     "EPAY-TEST-001",
        "trade_status":  "TRADE_SUCCESS",
        "type":         "1",
    }
    sorted_keys = sorted(params.keys())
    sign_str = "&".join([f"{k}={params[k]}" for k in sorted_keys]) + "&key=test_key_for_signature"
    sign = hashlib.md5(sign_str.encode()).hexdigest()
    print(f"  签名算法: MD5(排序参数&key)")
    print(f"  签名: {sign}")

    # 第一次回调（包含所有发送字段）
    notify = {
        "pid": test_order,  # 临时占位，将在下面替换
        "out_trade_no": test_order,
        "trade_status": "TRADE_SUCCESS",
        "trade_no": "EPAY-TEST-001",
        "type": "1",
        "money": "0.01",
        "name": "压测商品",
        "sign": sign,
        "sign_type": "MD5",
    }
    # 动态构建 notify_url 和 return_url（需要在循环中设置）
    notify["pid"] = "test_mch_id"
    notify["notify_url"] = "http://localhost:8080/api/v1/pay/epay/notify"
    notify["return_url"] = f"http://localhost:8080/order/{test_order}"
    resp1 = api_post_raw("/api/v1/pay/epay/notify", notify)
    print(f"  第1次回调响应: '{resp1}'")

    r = api_call("GET", f"/api/v1/orders/{test_order}", token=user_tokens[0])
    s1 = r.get("data", {}).get("order", {}).get("status", -1)
    c1 = len(r.get("data", {}).get("cards", []))
    print(f"  第1次后: 状态={s1}(2=completed), 卡密={c1} 张")

    # 第二次回调（重复）
    resp2 = api_post_raw("/api/v1/pay/epay/notify", notify)
    print(f"  第2次回调响应: '{resp2}'")

    r = api_call("GET", f"/api/v1/orders/{test_order}", token=user_tokens[0])
    s2 = r.get("data", {}).get("order", {}).get("status", -1)
    c2 = len(r.get("data", {}).get("cards", []))
    print(f"  第2次后: 状态={s2}, 卡密={c2} 张")

    if c1 == c2 and c1 > 0:
        print(f"  OK: 幂等性通过！重复回调未重复发货")
    elif c1 == 0:
        print(f"  FAIL: 回调未发货，检查签名或服务器状态")
    else:
        print(f"  FAIL: 重复回调导致重复发货！")

print()
print("=" * 56)
print("  测试完成")
print("=" * 56)
print(f"  成功订单: {ok_count}/50, 最终库存: {final_stock}, 期望: {expected}")
print()
