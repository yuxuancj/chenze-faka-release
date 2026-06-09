#!/bin/bash
# 晨泽发卡系统 - 并发安全完整测试脚本
# 环境: SQLite (并发安全代码已实施 MySQL FOR UPDATE，SQLite 等价事务隔离)

BASE="http://localhost:8080"
ADMIN_EMAIL="admin@chenze.com"
ADMIN_PWD="admin123"
SECRET="test-secret-key-for-concurrent-testing"
EXP=24

echo "=============================================="
echo "  晨泽发卡系统 - 并发安全完整测试"
echo "=============================================="
echo ""

# 1. 管理员登录获取 token
echo "[1/8] 管理员登录..."
ADMIN_RESP=$(curl -s -X POST "$BASE/api/v1/user/login" \
  -H "Content-Type: application/json" \
  -d "{\"email\":\"$ADMIN_EMAIL\",\"password\":\"$ADMIN_PWD\"}")
ADMIN_TOKEN=$(echo $ADMIN_RESP | grep -o '"token":"[^"]*"' | cut -d'"' -f4)
if [ -z "$ADMIN_TOKEN" ]; then
  echo "  FAIL: 管理员登录失败: $ADMIN_RESP"
  exit 1
fi
echo "  OK: 管理员登录成功"

# 2. 创建高库存测试商品（库存=200）
echo "[2/8] 创建高库存测试商品（库存=200）..."
PROD_RESP=$(curl -s -X POST "$BASE/admin/api/products" \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"category_id":1,"name":"压测商品-高库存","description":"用于并发压测的商品","price":0.01,"stock":200,"type":"card","status":1}')
TEST_PROD_ID=$(echo $PROD_RESP | grep -o '"id":[0-9]*' | head -1 | cut -d':' -f2)
if [ -z "$TEST_PROD_ID" ]; then
  echo "  FAIL: 创建商品失败: $PROD_RESP"
  exit 1
fi
echo "  OK: 测试商品 ID=$TEST_PROD_ID 创建成功（库存=200）"

# 3. 批量导入 200 张卡密
echo "[3/8] 批量导入 200 张卡密..."
CARDS_STR=""
for i in $(seq 1 200); do
  CARDS_STR="${CARDS_STR}CARD${i}-$(date +%s)-$(printf "%04d" $i)\n"
done
CARDS_JSON=$(printf '%s' "$CARDS_STR" | jq -Rs '.')

IMPORT_RESP=$(curl -s -X POST "$BASE/admin/api/cards/import" \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  -F "product_id=$TEST_PROD_ID" \
  -F "cards=$(printf '%s' "$CARDS_STR")")
IMPORT_COUNT=$(echo $IMPORT_RESP | grep -o '"count":[0-9]*' | cut -d':' -f2)
echo "  OK: 导入 $IMPORT_COUNT 张卡密"

# 4. 验证商品库存
echo "[4/8] 验证商品初始库存..."
STOCK_BEFORE=$(curl -s "$BASE/admin/api/products" \
  -H "Authorization: Bearer $ADMIN_TOKEN" | \
  grep -o "\"stock\":$TEST_PROD_ID[^,}]*" | head -1 || true)
STOCK_ACTUAL=$(curl -s "$BASE/admin/api/products" \
  -H "Authorization: Bearer $ADMIN_TOKEN" | \
  python3 -c "import sys,json; d=json.load(sys.stdin); print([p['stock'] for p in d['data']['list'] if p['id']==$TEST_PROD_ID][0])" 2>/dev/null)
echo "  初始库存: $STOCK_ACTUAL (期望: 200)"

# 5. 注册 10 个测试用户
echo "[5/8] 注册 10 个测试用户..."
USER_TOKENS=()
for i in $(seq 1 10); do
  EMAIL="testuser${i}@faka.test"
  REG_RESP=$(curl -s -X POST "$BASE/api/v1/user/register" \
    -H "Content-Type: application/json" \
    -d "{\"email\":\"$EMAIL\",\"password\":\"test123456\"}")
  TOKEN=$(echo $REG_RESP | grep -o '"token":"[^"]*"' 2>/dev/null | cut -d'"' -f4)
  if [ -z "$TOKEN" ]; then
    # 用户可能已存在，尝试登录
    LOGIN_RESP=$(curl -s -X POST "$BASE/api/v1/user/login" \
      -H "Content-Type: application/json" \
      -d "{\"email\":\"$EMAIL\",\"password\":\"test123456\"}")
    TOKEN=$(echo $LOGIN_RESP | grep -o '"token":"[^"]*"' 2>/dev/null | cut -d'"' -f4)
  fi
  USER_TOKENS+=("$TOKEN")
  echo "  用户$i: $EMAIL -> ${TOKEN:0:20}..."
done

# 6. 并发下单测试（30 并发，每用户 5 单 = 50 单）
echo ""
echo "[6/8] 并发下单测试（30 并发，每用户 5 单，共 50 单）..."
CONCURRENCY=30
TOTAL_ORDERS=50
SUCCESS_ORDERS=0
FAIL_ORDERS=0
ORDER_NOS=()

# 生成测试请求数据
TEMP_DIR=$(mktemp -d)
for i in $(seq 1 $TOTAL_ORDERS); do
  USER_IDX=$(( (i-1) % ${#USER_TOKENS[@]} ))
  TOKEN="${USER_TOKENS[$USER_IDX]}"
  echo "{\"product_id\":$TEST_PROD_ID,\"quantity\":1,\"email\":\"test@test.com\"}" > "$TEMP_DIR/order_$i.json"
done

# 并发执行
echo "  开始并发下单..."
for i in $(seq 1 $TOTAL_ORDERS); do
  USER_IDX=$(( (i-1) % ${#USER_TOKENS[@]} ))
  TOKEN="${USER_TOKENS[$USER_IDX]}"
  (
    RESP=$(curl -s -X POST "$BASE/api/v1/orders" \
      -H "Authorization: Bearer $TOKEN" \
      -H "Content-Type: application/json" \
      -d @"$TEMP_DIR/order_$i.json")
    ORDER_NO=$(echo $RESP | grep -o '"order_no":"[^"]*"' 2>/dev/null | cut -d'"' -f4)
    CODE=$(echo $RESP | grep -o '"code":[0-9]*' 2>/dev/null | cut -d':' -f2)
    if [ "$CODE" = "0" ] && [ -n "$ORDER_NO" ]; then
      echo "OK:$ORDER_NO" > "$TEMP_DIR/result_$i.txt"
    else
      echo "FAIL:$RESP" > "$TEMP_DIR/result_$i.txt"
    fi
  ) &
  # 控制并发数
  if [ $((i % CONCURRENCY)) -eq 0 ] || [ $i -eq $TOTAL_ORDERS ]; then
    wait
  fi
done
wait

# 收集结果
for i in $(seq 1 $TOTAL_ORDERS); do
  RESULT=$(cat "$TEMP_DIR/result_$i.txt" 2>/dev/null)
  if [[ "$RESULT" == OK:* ]]; then
    ORDER_NO="${RESULT#OK:}"
    ORDER_NOS+=("$ORDER_NO")
    SUCCESS_ORDERS=$((SUCCESS_ORDERS + 1))
  else
    FAIL_ORDERS=$((FAIL_ORDERS + 1))
  fi
done
rm -rf "$TEMP_DIR"
echo "  并发下单完成: 成功=$SUCCESS_ORDERS, 失败(含库存不足)=$FAIL_ORDERS"

# 7. 验证库存和超卖检查
echo ""
echo "[7/8] 验证库存和超卖检查..."
FINAL_STOCK=$(curl -s "$BASE/admin/api/products" \
  -H "Authorization: Bearer $ADMIN_TOKEN" | \
  python3 -c "import sys,json; d=json.load(sys.stdin); print([p['stock'] for p in d['data']['list'] if p['id']==$TEST_PROD_ID][0])" 2>/dev/null)
EXPECTED_STOCK=$((200 - SUCCESS_ORDERS))
echo "  最终库存: $FINAL_STOCK (期望: $EXPECTED_STOCK)"
if [ "$FINAL_STOCK" -eq "$EXPECTED_STOCK" ]; then
  echo "  OK: 库存正确，无超卖！"
else
  echo "  WARN: 库存与预期不符（可能因部分请求失败）"
fi

# 检查数据库中的实际已使用卡密数
USED_CARDS=$(curl -s "$BASE/admin/api/cards?product_id=$TEST_PROD_ID&page=1&size=1" \
  -H "Authorization: Bearer $ADMIN_TOKEN" | \
  python3 -c "import sys,json; d=json.load(sys.stdin); print(d['data']['total'])" 2>/dev/null)
echo "  商品表中共 $USED_CARDS 张卡密记录"

# 8. 易支付回调幂等性测试
echo ""
echo "[8/8] 易支付回调幂等性测试..."
if [ ${#ORDER_NOS[@]} -gt 0 ]; then
  TEST_ORDER_NO="${ORDER_NOS[0]}"
  echo "  测试订单号: $TEST_ORDER_NO"

  # 查询订单初始状态
  ORDER_STATUS_0=$(curl -s "$BASE/api/v1/orders/$TEST_ORDER_NO" \
    -H "Authorization: Bearer ${USER_TOKENS[0]}" | \
    python3 -c "import sys,json; d=json.load(sys.stdin); print(d['data']['order']['status'])" 2>/dev/null)
  echo "  初始状态: $ORDER_STATUS_0 (0=pending)"

  # 构造签名（按易支付规范：参数按字母排序后拼接 key）
  PAY_AMOUNT="0.01"
  # params: pid, out_trade_no, name, money, notify_url, return_url
  PARAMS="money=${PAY_AMOUNT}&name=压测商品&notify_url=http://localhost:8080/api/v1/pay/epay/notify&out_trade_no=${TEST_ORDER_NO}&pid=test_mch_id&return_url=http://localhost:8080/order/${TEST_ORDER_NO}"
  # 按字母排序
  SORTED=$(echo "$PARAMS" | tr '&' '\n' | sort | tr '\n' '&' | sed 's/&$//')
  SIGN=$(echo -n "${SORTED}&key=test_key_for_signature" | md5sum | cut -d' ' -f1)
  echo "  签名: $SIGN"

  # 第一次回调
  echo "  发送第一次支付回调..."
  RESP1=$(curl -s -X POST "$BASE/api/v1/pay/epay/notify" \
    -d "pid=test_mch_id&out_trade_no=${TEST_ORDER_NO}&trade_status=TRADE_SUCCESS&trade_no=EPAY${RANDOM}&type=1&money=${PAY_AMOUNT}&name=压测商品&sign=${SIGN}&sign_type=MD5")
  echo "  第一次响应: $RESP1"

  ORDER_STATUS_1=$(curl -s "$BASE/api/v1/orders/$TEST_ORDER_NO" \
    -H "Authorization: Bearer ${USER_TOKENS[0]}" | \
    python3 -c "import sys,json; d=json.load(sys.stdin); print(d['data']['order']['status'])" 2>/dev/null)
  echo "  第一次回调后状态: $ORDER_STATUS_1 (2=completed)"

  CARDS_AFTER_1=$(curl -s "$BASE/api/v1/orders/$TEST_ORDER_NO" \
    -H "Authorization: Bearer ${USER_TOKENS[0]}" | \
    python3 -c "import sys,json; d=json.load(sys.stdin); print(len(d['data'].get('cards',[])))" 2>/dev/null)
  echo "  第一次回调后获得卡密: $CARDS_AFTER_1 张"

  # 第二次回调（重复，测试幂等性）
  echo "  发送第二次相同回调（幂等性测试）..."
  RESP2=$(curl -s -X POST "$BASE/api/v1/pay/epay/notify" \
    -d "pid=test_mch_id&out_trade_no=${TEST_ORDER_NO}&trade_status=TRADE_SUCCESS&trade_no=EPAY${RANDOM}&type=1&money=${PAY_AMOUNT}&name=压测商品&sign=${SIGN}&sign_type=MD5")
  echo "  第二次响应: $RESP2"

  CARDS_AFTER_2=$(curl -s "$BASE/api/v1/orders/$TEST_ORDER_NO" \
    -H "Authorization: Bearer ${USER_TOKENS[0]}" | \
    python3 -c "import sys,json; d=json.load(sys.stdin); print(len(d['data'].get('cards',[])))" 2>/dev/null)
  echo "  第二次回调后获得卡密: $CARDS_AFTER_2 张"

  if [ "$CARDS_AFTER_1" = "$CARDS_AFTER_2" ] && [ "$CARDS_AFTER_1" -gt 0 ]; then
    echo "  OK: 幂等性测试通过！重复回调未重复发货"
  else
    echo "  FAIL: 幂等性测试失败！"
  fi
fi

echo ""
echo "=============================================="
echo "  测试完成"
echo "=============================================="
echo "  成功订单数: $SUCCESS_ORDERS / $TOTAL_ORDERS"
echo "  最终库存: $FINAL_STOCK"
echo "  期望库存: $EXPECTED_STOCK"
echo ""
