# 晨泽发卡系统 - 完整架构设计文档

## 一、项目概述

晨泽发卡是一款企业级虚拟商品自动发卡平台，采用 Go 语言开发，支持高并发、分布式部署。系统涵盖用户管理、商品管理、订单管理、支付对接、分销推广、商户入驻、工单售后、营销促销、数据报表、插件扩展等完整功能。

### 技术栈
| 层级 | 技术选型 | 说明 |
|------|----------|------|
| 后端语言 | Go 1.21+ | 高性能、静态编译 |
| Web 框架 | Gin | 轻量、高性能 HTTP 框架 |
| ORM | GORM | 功能完善的 ORM 库 |
| 数据库 | MySQL 5.7+ | 关系型数据库 |
| 缓存 | Redis（可选） | 缓存、队列、分布式锁 |
| 前端-后台 | Layui | 经典企业级 UI 框架 |
| 前端-前台 | 响应式 HTML/CSS | 适配 PC 和移动端 |
| 部署 | 单二进制 + systemd | 无依赖、一键安装 |

## 二、系统架构图
┌─────────────────────────────────────────────────────────────────┐
│ 前端层 │
├──────────────┬──────────────┬──────────────┬───────────────────┤
│ PC 前台 │ H5 移动端 │ 管理后台 │ 微信小程序 │
│ (响应式) │ (自适应) │ (Layui) │ (独立项目) │
└──────────────┴──────────────┴──────────────┴───────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────────┐
│ API 网关层 │
├─────────────────────────────────────────────────────────────────┤
│ JWT 认证 │ 限流中间件 │ 日志中间件 │ 跨域处理 │ 恢复中间件 │
└─────────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────────┐
│ 业务逻辑层 (Service) │
├─────────┬─────────┬─────────┬─────────┬─────────┬─────────────┤
│UserSvc │ProductSvc│OrderSvc │PaymentSvc│DistSvc │TicketSvc │
│MerchantSvc│AgentSvc│SupplierSvc│CouponSvc│SeckillSvc│PointsSvc │
│PluginSvc │WebhookSvc│ReportSvc│BackupSvc│CronSvc │AISvc │
└─────────┴─────────┴─────────┴─────────┴─────────┴─────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────────┐
│ 数据访问层 (DAO/Model) │
├─────────────────────────────────────────────────────────────────┤
│ GORM │ 数据库连接池 │ 事务管理 │ 软删除 │ 钩子 │ 预加载 │
└─────────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────────┐
│ 基础设施层 │
├──────────────┬──────────────┬──────────────┬───────────────────┤
│ MySQL │ Redis │ 文件存储 │ 队列(可选) │
└──────────────┴──────────────┴──────────────┴───────────────────┘

## 三、目录结构
chenze-faka/
├── cmd/
│   └── main.go                                    # 主程序入口
├── internal/
│   ├── controller/                                # API 控制器层
│   │   ├── user.go
│   │   ├── product.go
│   │   ├── order.go
│   │   ├── payment.go
│   │   ├── distribution.go
│   │   ├── merchant.go
│   │   ├── agent.go
│   │   ├── supplier.go
│   │   ├── ticket.go
│   │   ├── coupon.go
│   │   ├── seckill.go
│   │   ├── points.go
│   │   ├── plugin.go
│   │   ├── webhook.go
│   │   ├── report.go
│   │   └── admin.go                               # 后台管理控制器
│   ├── service/                                   # 业务逻辑层
│   │   ├── user_service.go
│   │   ├── product_service.go
│   │   ├── order_service.go
│   │   ├── payment_service.go
│   │   ├── distribution_service.go
│   │   ├── points_service.go
│   │   └── ...
│   ├── model/                                     # 数据模型（GORM）
│   │   ├── user.go
│   │   ├── product.go
│   │   ├── product_sku.go
│   │   ├── order.go
│   │   ├── order_card.go
│   │   ├── card.go
│   │   ├── payment.go
│   │   ├── commission.go
│   │   ├── merchant.go
│   │   ├── agent.go
│   │   ├── supplier.go
│   │   ├── ticket.go
│   │   ├── coupon.go
│   │   ├── seckill.go
│   │   ├── points.go
│   │   ├── points_exchange.go
│   │   ├── plugin.go
│   │   ├── webhook.go
│   │   ├── backup.go
│   │   └── audit_log.go
│   ├── middleware/                                # 中间件
│   │   ├── auth.go           # JWT 认证
│   │   ├── cors.go           # 跨域处理
│   │   ├── ratelimit.go      # 限流
│   │   ├── logger.go         # 日志记录
│   │   ├── recovery.go       # 异常恢复
│   │   └── csrf.go           # CSRF 防护
│   ├── router/                                    # 路由定义
│   │   ├── router.go         # 路由注册
│   │   ├── api.go            # API 路由组
│   │   └── admin.go          # 后台路由组
│   └── pkg/                                       # 内部公共包
│       ├── config/          # 配置管理
│       ├── db/              # 数据库连接
│       ├── cache/           # 缓存封装
│       ├── queue/           # 队列处理
│       ├── crypto/          # 加密工具
│       ├── jwt/             # JWT 工具
│       ├── logger/          # 日志工具
│       ├── response/        # 统一响应
│       ├── validator/       # 参数验证
│       └── cron/            # 定时任务
├── pkg/                                          # 可复用的公共库
│   ├── alipay/          # 支付宝 SDK 封装
│   ├── epay/            # 易支付 SDK
│   ├── payjs/           # PAYJS SDK
│   ├── usdt/            # USDT 支付
│   ├── sms/             # 短信服务
│   ├── email/           # 邮件服务
│   └── storage/         # 存储适配（本地/OSS）
├── web/                                          # 静态资源
│   ├── admin/           # 后台管理界面（Layui）
│   ├── front/           # 前台界面（响应式）
│   └── install/         # 安装向导页面
├── config/
│   ├── config.yaml.example # 配置文件示例
│   └── config.yaml         # 实际配置（安装时生成）
├── storage/                                      # 运行时数据
│   ├── logs/            # 日志文件
│   ├── backups/         # 备份文件
│   └── uploads/         # 上传文件
├── scripts/
│   ├── install.sh      # 一键安装脚本
│   └── chenze_faka.service # systemd 服务文件
├── docs/
│   ├── api/             # API 文档（Swagger）
│   └── guide/           # 用户手册
├── go.mod
├── go.sum
├── README.md
└── CHANGELOG.md

## 四、核心模块边界定义

### 4.1 用户模块 (User)
| 职责 | 禁止 |
|------|------|
| 注册、登录、JWT 签发 | 不直接操作订单表 |
| 个人资料管理 | 不直接操作商品表 |
| 等级计算与升级 | 不直接处理支付回调 |
| 积分增减 |  |
| 签到功能 |  |

### 4.2 商品模块 (Product)
| 职责 | 禁止 |
|------|------|
| 商品 CRUD、分类管理 | 不直接修改订单状态 |
| SKU 管理 | 不直接调用支付接口 |
| 库存管理（扣减、预警） |  |
| 卡密导入导出 |  |
| 货源 API 同步 |  |

### 4.3 订单模块 (Order)
| 职责 | 禁止 |
|------|------|
| 订单创建、状态流转 | 不直接操作商品库存（应调用 ProductService） |
| 订单查询、取消、退款 | 不直接调用支付接口（应调用 PaymentService） |
| 发货处理（调用对应发货服务） |  |
| 超时订单自动关闭 |  |

### 4.4 支付模块 (Payment)
| 职责 | 禁止 |
|------|------|
| 支付接口封装（易支付/支付宝/PAYJS/USDT） | 不直接修改订单状态（应由 OrderService 处理） |
| 支付回调验签、幂等处理 | 不操作商品库存 |
| 退款处理 |  |

### 4.5 分销模块 (Distribution)
| 职责 | 禁止 |
|------|------|
| 邀请码/链接管理 | 不直接修改订单金额 |
| 上下级关系绑定 | 不直接操作支付 |
| 佣金计算与结算 |  |
| 提现申请与审核 |  |

### 4.6 工单模块 (Ticket)
| 职责 | 禁止 |
|------|------|
| 工单创建、回复 | 不直接修改订单状态（可通过工单触发补发） |
| 工单分配、解决 |  |
| 附件管理 |  |

### 4.7 插件模块 (Plugin)
| 职责 | 禁止 |
|------|------|
| 插件上传、安装、卸载 | 插件代码不能直接访问核心数据库 |
| 钩子注册与触发 |  |
| 插件市场接口 |  |

## 五、数据库设计核心表

### 5.1 用户相关表

```sql
-- 用户表
users (id, email, password, nickname, avatar, balance, points, level, parent_id, invite_code, status, created_at, updated_at)

-- 用户等级表
user_levels (id, name, min_amount, discount, created_at)

-- 成长值记录
growth_records (id, user_id, amount, type, description, created_at)

-- 积分记录
points_records (id, user_id, amount, type, balance_after, description, created_at)

-- 签到记录
signin_records (id, user_id, continuous_days, reward_points, created_at)
```

### 5.2 商品相关表

```sql
-- 商品分类表
categories (id, parent_id, name, sort, status, created_at)

-- 商品表
products (id, category_id, name, description, price, stock, sales, type, status, is_hidden, created_at, updated_at)

-- 商品 SKU 表
product_skus (id, product_id, sku_code, spec_names, price, stock, sales, image, weight, status, created_at)

-- 卡密表
cards (id, product_id, card_data, status, order_id, is_premium, premium_price, expire_days, created_at)

-- 卡密授权记录表（用于核销策略）
card_authorizations (id, card_id, user_id, device_id, ip_address, used_count, expire_at, created_at)

-- API 授权商品配置
product_apis (id, product_id, method, url, headers, body, success_field, auth_code_field, expire_days, created_at)
```

### 5.3 订单相关表

```sql
-- 订单表
orders (id, order_no, user_id, product_id, sku_id, sku_snapshot, amount, pay_type, status, paid_at, completed_at, created_at)

-- 订单卡密发货记录
order_cards (id, order_id, card_id, card_data, created_at)

-- 支付记录
payments (id, order_id, payment_method, transaction_id, amount, status, callback_data, created_at)

-- 退款记录
refunds (id, order_id, amount, reason, status, created_at)
```

### 5.4 分销相关表

```sql
-- 分销关系树
distribution_tree (id, user_id, parent_id, level, path, created_at)

-- 佣金记录
commissions (id, user_id, from_user_id, order_id, amount, level, status, settled_at, created_at)

-- 提现申请
withdraws (id, user_id, amount, account_type, account, real_name, status, remark, processed_at, created_at)
```

### 5.5 营销相关表

```sql
-- 优惠券表
coupons (id, name, type, value, min_amount, total, used, expire_start, expire_end, status, created_at)

-- 用户优惠券
user_coupons (id, user_id, coupon_id, order_id, status, used_at, created_at)

-- 秒杀活动表
seckills (id, product_id, seckill_price, stock, sold, preheat_start, preheat_end, start_time, end_time, limit_per_user, status, created_at)

-- 积分兑换商品表
points_exchange (id, name, type, points, stock, exchange_value, status, created_at)
```

### 5.6 多角色相关表

```sql
-- 商户表
merchants (id, user_id, name, slug, logo, commission_rate, balance, frozen_balance, status, created_at)

-- 代理商表
agents (id, user_id, name, slug, commission_rate, balance, status, created_at)

-- 供应商表
suppliers (id, name, contact, phone, email, commission_rate, balance, status, created_at)

-- 供应商商品表
supplier_products (id, supplier_id, product_id, supply_price, stock, status, created_at)
```

### 5.7 工单相关表

```sql
-- 工单表
tickets (id, user_id, merchant_id, order_id, type, priority, title, content, status, assigned_to, closed_by, closed_at, created_at)

-- 工单回复表
ticket_replies (id, ticket_id, user_id, content, attachments, created_at)

-- 工单附件表
ticket_attachments (id, ticket_id, file_name, file_path, file_size, created_at)
```

### 5.8 扩展功能表

```sql
-- 插件表
plugins (id, name, identifier, version, author, description, price, status, installed_at, created_at)

-- Webhook 配置表
webhooks (id, name, event, url, secret, status, retry_count, created_at)

-- Webhook 发送日志
webhook_logs (id, webhook_id, event, request, response, status_code, duration, created_at)

-- 备份记录表
backups (id, name, type, file_path, file_size, status, created_at)

-- 操作审计日志
audit_logs (id, user_id, username, action, target, ip_address, user_agent, old_data, new_data, created_at)
```

## 六、API 设计规范

### 6.1 统一响应格式

```json
{
  "code": 0,        // 0=成功，非0=错误码
  "msg": "success", // 提示信息
  "data": {}        // 响应数据
}
```

### 6.2 HTTP 状态码

| 状态码 | 说明 |
|--------|------|
| 200 | 成功 |
| 400 | 参数错误 |
| 401 | 未登录/Token 失效 |
| 403 | 无权限 |
| 404 | 资源不存在 |
| 500 | 服务器内部错误 |

### 6.3 路由分组规范

```text
/api/v1/                 # 公开 API（无需登录）
/api/v1/user/            # 用户 API（需登录）
/api/v1/merchant/        # 商户 API（需商户权限）
/api/v1/agent/           # 代理商 API
/admin/api/v1/           # 管理后台 API（需管理员权限）
```

## 七、安全规范

### 7.1 认证与授权

- 所有 API（除公开接口）必须携带 JWT Token（Header: `Authorization: Bearer <token>`）
- Token 有效期：Access Token 2 小时，Refresh Token 7 天
- 后台敏感操作（修改余额、导出订单）需二次验证（输入登录密码）

### 7.2 数据安全

- 用户密码使用 bcrypt 加密（cost >= 10）
- 敏感配置（支付密钥、数据库密码）存储在 `config.yaml`，禁止硬编码
- 日志中不得记录明文密码、支付密钥、用户身份证号

### 7.3 防攻击

- SQL 注入：所有查询使用 GORM 参数化
- XSS：前端输出使用 `html.EscapeString` 或模板引擎自动转义
- CSRF：后台表单提交需携带 CSRF Token
- 限流：下单接口每 IP 每分钟最多 10 次
- 暴力破解：登录失败 5 次锁定 IP 15 分钟

## 八、性能优化指南

### 8.1 数据库优化

- 高频查询字段建立索引：`order_no`, `user_id`, `product_id`, `status`
- 订单表按月分区（`PARTITION BY RANGE (YEAR(created_at))`）
- 分页查询使用 `WHERE id > last_id` 优化大偏移量

### 8.2 缓存策略

- 商品详情：Redis 缓存 5 分钟
- 商品分类：Redis 缓存 10 分钟
- 系统设置：Redis 缓存 30 分钟
- 用户会话：Redis 存储（支持多实例）

### 8.3 异步处理

- 邮件发送：放入队列异步处理
- Webhook 推送：放入队列异步处理
- 卡密批量导入：后台任务异步处理

## 九、部署架构

```text
┌─────────────────────────────────────────────────────────────┐
│                        用户请求                             │
└─────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────┐
│                      Nginx（反向代理）                       │
│                    SSL 终止 / 静态资源 / 负载均衡            │
└─────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────┐
│                    chenze_faka（Go 二进制）                  │
│                    多实例部署（端口 8080-8082）              │
└─────────────────────────────────────────────────────────────┘
                              │
              ┌───────────────┼───────────────┐
              ▼               ▼               ▼
        ┌──────────┐   ┌──────────┐   ┌──────────┐
        │  MySQL   │   │  Redis   │   │   OSS    │
        │ 主从/读写 │   │ 缓存/队列 │   │ 文件存储 │
        └──────────┘   └──────────┘   └──────────┘
```

## 十、开发与贡献指南

### 10.1 新增功能流程

1. 在 `internal/model/` 定义数据模型（如需新表）
2. 在 `internal/service/` 实现业务逻辑
3. 在 `internal/controller/` 实现 API 控制器
4. 在 `internal/router/` 注册路由
5. 编写单元测试
6. 更新 API 文档（Swagger 注释）

### 10.2 代码审查要点

- 是否有单元测试覆盖核心逻辑？
- 数据库查询是否有索引？
- 是否有并发安全问题？
- 错误处理是否完整？
- 日志是否记录关键操作？
- 是否遵循编码规范？

---

文档版本: v1.0
最后更新: 2026-06-09
维护者: 晨泽科技
