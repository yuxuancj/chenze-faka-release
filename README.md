# 晨泽发卡系统 (chenze-faka)

一个基于 Go + Vue 3 + Arco Design 的轻量虚拟商品发卡系统。支持商品/卡密管理、订单创建、余额支付、易支付对接、安装向导、后台管理。

## 特性

- 单二进制分发（前端静态资源已嵌入）
- 首次访问自动跳转到 `/install` 进行安装
- 商品分类、商品 CRUD、卡密批量导入
- 购物车 → 订单 → 支付 → 自动发卡
- 余额支付 + 易支付（第三方支付跳转）
- 后台管理：商品、卡密、订单、用户
- 深蓝安全主题（SCDN 风格）

## 目录结构

```
chenze-faka/
├── cmd/main.go          # 程序入口（嵌入前端 dist）
├── internal/
│   ├── controller/      # Gin 控制器
│   ├── service/         # 业务逻辑层
│   ├── model/           # GORM 模型
│   ├── middleware/      # 中间件（JWT、安装拦截）
│   └── pkg/             # 工具（DB、响应、JWT）
├── frontend/            # Vue 3 前端
├── config/              # 配置模块
├── config.yaml          # 安装后自动生成
└── install.lock         # 安装后自动生成
```

## 快速开始

### 1. 运行

```bash
# 从成品包直接运行
./chenze_faka
# 默认监听 0.0.0.0:8080
```

然后浏览器打开 `http://127.0.0.1:8080/`，系统会自动跳转到安装向导。

### 2. 安装

1. 填写数据库连接信息（MySQL 5.7+，提前创建数据库）
2. 填写第一个管理员账号
3. 点击「开始安装」，系统将自动建表并写入 `config.yaml` 与 `install.lock`

安装完成后使用管理员账号登录，进入 `/admin` 即可开始管理商品与卡密。

### 3. 配置易支付

安装完成后，编辑 `config.yaml`：

```yaml
pay:
  epay:
    api_url: https://你的易支付域名/mapi.php
    pid: "商户号"
    key: "密钥"
```

修改后重启服务。异步回调地址为 `/api/pay/epay/notify`，跳转地址 `/api/pay/epay/return`。

## 开发构建

```bash
# 后端
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o chenze_faka ./cmd

# 前端
cd frontend
npm install
npm run build
```

前端构建产物输出到 `cmd/dist`，由后端 `go:embed` 嵌入。

## 常见接口

- `POST /api/register` `{email, password, nickname}` 注册
- `POST /api/login` `{email, password}` 登录（返回 token）
- `GET /api/products` 商品列表
- `POST /api/orders` 创建订单 `{items:[{product_id,quantity}]}`
- `POST /api/pay` 发起支付 `{order_no, method:"balance"|"epay"}`
- `POST /api/pay/epay/notify` 易支付异步回调

## 许可证

MIT
