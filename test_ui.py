#!/usr/bin/env python3
"""前端UI全面测试脚本（基于 HTTP 请求的黑盒测试）"""
import urllib.request
import urllib.parse
import json
import sys

BASE = "http://127.0.0.1:8080"
results = []

def log(name, ok, detail=""):
    status = "✅" if ok else "❌"
    print(f"{status} {name} {detail}")
    results.append((name, ok, detail))

def get(path):
    req = urllib.request.Request(BASE + path)
    try:
        with urllib.request.urlopen(req, timeout=5) as r:
            body = r.read().decode('utf-8', errors='ignore')
            return r.status, body
    except Exception as e:
        return 0, str(e)

def post(path, data_json, token=None):
    body = json.dumps(data_json).encode('utf-8')
    headers = {'Content-Type': 'application/json'}
    if token:
        headers['Authorization'] = 'Bearer ' + token
    req = urllib.request.Request(BASE + path, data=body, headers=headers, method='POST')
    try:
        with urllib.request.urlopen(req, timeout=5) as r:
            data = json.loads(r.read().decode('utf-8', errors='ignore'))
            return r.status, data
    except Exception as e:
        if hasattr(e, 'read'):
            try:
                return getattr(e, 'code', 500), json.loads(e.read().decode('utf-8', errors='ignore'))
            except:
                pass
        return 0, {}

print("=" * 60)
print("第一部分：页面加载与路由测试")
print("=" * 60)

pages = [
    ("/", "首页"),
    ("/products", "商品列表"),
    ("/product/1", "商品详情(ID=1)"),
    ("/cart", "购物车"),
    ("/checkout", "结算页"),
    ("/user/login", "登录页"),
    ("/user/register", "注册页"),
    ("/admin/login", "后台登录页"),
    ("/admin/", "后台仪表盘"),
    ("/nothing", "不存在的路径 (404)"),
]

for path, name in pages:
    code, body = get(path)
    if code == 200:
        is_vue = '<div id="app"' in body or '<div id="app">' in body or 'vue' in body.lower() or 'tailwind' in body.lower() or len(body) > 2000
        log(f"{name} ({path})", code == 200, f"HTTP {code}, 内容长度 {len(body)}")
    else:
        log(f"{name} ({path})", False, f"HTTP {code}")

print()
print("=" * 60)
print("第二部分：API 连通性测试")
print("=" * 60)

api_tests = [
    ("GET /api/v1/products", "/api/v1/products?page=1&size=10"),
    ("GET /api/v1/products/1", "/api/v1/products/1"),
]
for name, path in api_tests:
    code, body = get(path)
    if isinstance(body, str):
        try:
            body = json.loads(body)
        except:
            body = {"_raw": body[:100]}
    ok = code == 200 and (isinstance(body, dict)) and body.get('code') == 0
    log(name, ok, f"HTTP {code}, code={body.get('code', 'n/a') if isinstance(body, dict) else 'n/a'}")

print()
print("=" * 60)
print("第三部分：用户注册/登录/下单测试")
print("=" * 60)

# 注册
code, data = post("/api/v1/user/register", {"email": "uitest@example.com", "password": "test123", "nickname": "UI测试用户"})
log(f"注册 (POST /api/v1/user/register)", code == 200 and isinstance(data, dict), f"HTTP {code}, data={data}")

# 登录
code, data = post("/api/v1/user/login", {"email": "uitest@example.com", "password": "test123"})
ok_login = code == 200 and isinstance(data, dict) and data.get('code') == 0
log("登录", ok_login, f"HTTP {code}, resp={json.dumps(data)[:120]}")
token = ""
is_admin = False
if ok_login:
    payload = data.get('data') or data
    token = payload.get('token', '') or (payload.get('data') or {}).get('token', '') if isinstance(payload, dict) else ''
    if not token and isinstance(data.get('data'), dict):
        token = data['data'].get('token', '')
    user = data.get('data', {}).get('user', {}) if isinstance(data.get('data'), dict) else data.get('user', {})
    is_admin = user.get('is_admin', False)

# 错误密码登录
code2, data2 = post("/api/v1/user/login", {"email": "uitest@example.com", "password": "wrongpass"})
log("错误密码登录应拒绝", code2 == 200 and isinstance(data2, dict) and data2.get('code', 0) != 0, f"HTTP {code2}, code={data2.get('code', 'n/a')}")

# 创建订单
if token:
    code, data = post("/api/v1/orders", {"product_id": 1, "quantity": 1, "pay_type": "epay", "email": "uitest@example.com", "remark": "UI测试订单"}, token=token)
    log("创建订单", code == 200 and isinstance(data, dict) and data.get('code') == 0, f"HTTP {code}, resp={json.dumps(data)[:120]}")

# 管理员登录
code, data = post("/api/v1/user/login", {"email": "admin@example.com", "password": "admin123"})
admin_login_ok = code == 200 and isinstance(data, dict) and data.get('code') == 0
log("管理员登录", admin_login_ok, f"HTTP {code}, resp={json.dumps(data)[:120]}")
admin_token = ""
if admin_login_ok:
    if isinstance(data.get('data'), dict):
        admin_token = data['data'].get('token', '')
    else:
        admin_token = data.get('token', '')

# 后台仪表盘
if admin_token:
    code, body = get("/admin/api/dashboard")
    req = urllib.request.Request(BASE + "/admin/api/dashboard", headers={"Authorization": "Bearer " + admin_token})
    try:
        with urllib.request.urlopen(req, timeout=5) as r:
            d = json.loads(r.read().decode('utf-8'))
            ok = d.get('code') == 0 and isinstance(d.get('data'), dict)
            log("后台仪表盘", ok, f"code={d.get('code')}, data_keys={list(d.get('data', {}).keys())[:4]}")
    except Exception as e:
        log("后台仪表盘", False, str(e))

print()
print("=" * 60)
print("测试总结")
print("=" * 60)
total = len(results)
passed = sum(1 for _, ok, _ in results if ok)
failed = total - passed
print(f"总计: {total} 项, 通过: {passed}, 失败: {failed}")
if failed > 0:
    print("失败项:")
    for name, ok, detail in results:
        if not ok:
            print(f"  - {name}: {detail}")
    sys.exit(1)
else:
    print("所有测试通过")
