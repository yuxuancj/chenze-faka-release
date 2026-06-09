# 晨泽发卡系统

企业级虚拟商品自动发卡平台，Go 语言实现，支持高并发。

## 功能特性

- 用户注册 / 登录（JWT）
- 商品管理（分类、上下架、卡密）
- 订单管理（创建、状态机、自动发货）
- 易支付对接（回调验签 + 幂等）
- 管理后台（Layui）
- 响应式前台

## 技术栈

- Go 1.21+
- Gin
- GORM
- MySQL 5.7+
- Layui（后台）

## 快速开始

```bash
# 1. 修改 config.yaml 填写数据库连接
# 2. 启动
go build -o chenze-faka ./cmd
./chenze-faka
```

首次启动自动建表并写入默认管理员：

```
邮箱: admin@chenze.com
密码: admin123
```

## 目录结构

```
chenze-faka/
├── cmd/              # 程序入口
├── internal/         # 核心实现
│   ├── controller/
│   ├── service/
│   ├── model/
│   ├── middleware/
│   ├── router/
│   └── pkg/          # 内部公共组件
├── templates/        # HTML 模板
└── storage/          # 运行时数据
```

## 默认路由

| 路径 | 说明 |
|------|------|
| `GET /` | 前台首页 |
| `POST /api/v1/user/login` | 登录 |
| `POST /api/v1/user/register` | 注册 |
| `GET /admin/` | 后台仪表盘 |
| `POST /admin/api/products` | 后台添加商品 |
| `POST /admin/api/cards/import` | 卡密导入 |
| `POST /api/v1/pay/epay/notify` | 易支付回调 |

## 开发规范

详见 `.trae/` 目录（`ARCHITECTURE.md` / `CONVENTIONS.md` / `FUNCTIONALITY_CHECKLIST.md`）。
