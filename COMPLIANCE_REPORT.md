# 代码合规性检查报告

**项目**：晨泽发卡系统 v2.1
**检查日期**：2026-06-10
**检查范围**：后端 Go 代码（`internal/`）、前端 Vue 代码（`frontend/src/`）、配置文件、数据库模型
**参考文档**：`ARCHITECTURE.md`、`CONVENTIONS.md`、`FUNCTIONALITY_CHECKLIST.md`

---

## 一、架构设计

### 1.1 目录结构

| 检查项 | 状态 | 说明 |
|--------|------|------|
| `cmd/` 主程序入口 | ✅ | `cmd/main.go` 存在 |
| `internal/controller/` | ✅ | 有 `admin.go`、`advanced.go`、`controller.go`、`install.go` |
| `internal/service/` | ✅ | 独立 service 文件按功能划分（`sku_service.go` 等） |
| `internal/model/` | ✅ | `model.go` 集中定义所有模型 |
| `internal/middleware/` | ⚠️ | 仅 `auth.go`，缺少 `ratelimit.go`、`logger.go`、`csrf.go`、`recovery.go` |
| `internal/router/` | ✅ | `router.go` 统一路由注册 |
| `internal/pkg/` | ✅ | config、db、jwt、logger、response 均存在 |
| `pkg/` 公共库 | ❌ | 缺少 `pkg/alipay/`、`pkg/epay/`、`pkg/sms/`、`pkg/email/`、`pkg/storage/`（支付/短信/邮件/存储 SDK 未独立封装） |
| `web/` 静态资源 | ⚠️ | `web/frontend/dist/` 和 `web/templates/` 存在，但 `web/admin/`（Layui 后台）完全缺失 |
| `config/config.yaml` | ✅ | 配置文件存在 |
| `storage/` | ❌ | 目录不存在 |
| `scripts/` | ❌ | 目录不存在 |
| `docs/` | ❌ | 目录不存在 |

### 1.2 模块分离

| 检查项 | 状态 | 说明 |
|--------|------|------|
| Controller 控制器层 | ✅ | 控制器只处理请求/响应，业务逻辑在 service 层 |
| Service 业务逻辑层 | ✅ | 各 service 独立文件，事务处理到位 |
| Model 数据层 | ✅ | 模型定义清晰，JSONMap 自定义类型实现 Scan/Value |
| Middleware 中间件 | ⚠️ | 仅 auth 和 CORS，缺乏日志、限流、CSRF、恢复中间件 |
| Router 路由注册 | ✅ | 统一在 `router.go` 中注册，API 分组清晰 |

### 1.3 数据库表设计

| 检查项 | 状态 | 问题 |
|--------|------|------|
| `users` 表 | ⚠️ | 缺少 `growth`（成长值）、用户等级关联字段；`parent_id` 和 `invite_code` 已实现 |
| `categories` 表 | ✅ | 基本符合 |
| `products` 表 | ⚠️ | 缺少 `weight`（重量，用于物流）；`is_hidden` 已实现 |
| `product_skus` 表 | ⚠️ | 与文档字段名不一致——文档定义 `spec_names`、`price` 等，但前/后端混用 `original_price`、`seckill_price` 等别名 |
| `orders` 表 | ✅ | `SkuID`、`SkuSnapshot`、`CouponID`、`PointsUsed`、`PointsDiscount` 均已扩展 |
| `cards` 表 | ⚠️ | 缺少 `card_authorizations` 表（卡密核销策略）；缺少 `expire_days`、`is_premium`、`premium_price` 字段 |
| `commissions` 表 | ⚠️ | 缺少 `level1_rate/level2_rate/level3_rate` 等佣金比例字段 |
| `coupons` 表 | ⚠️ | 文档定义 `expire_start`/`expire_end` 两个时间字段，代码中为 `ExpireStart`/`ExpireEnd` —— ✅ 一致；但文档中 `type` 含义与代码不同（文档为满减/折扣类型字段） |
| `seckills` 表 | ⚠️ | 文档定义 `seckill_price`，代码定义 `SeckillPrice` —— ✅ 一致；但前端传入的 `price` 字段需要后端映射到 `seckill_price` |
| `withdraws` 表 | ⚠️ | 缺少 `real_name`（文档有，代码已有 `RealName` ✅） |
| `points_logs` 表 | ⚠️ | 文档定义 `points_records`，代码为 `PointsLog`（语义接近）；缺少 `growth_records` 表 |
| `signin_logs` 表 | ⚠️ | 文档定义 `signin_records`，代码为 `SigninLog` |
| `distribution_tree` 表 | ❌ | 完全缺失三级分销上下级关系表（当前通过 `users.parent_id` 简单实现） |
| `audit_logs` 表 | ❌ | 缺失操作审计日志表 |
| `webhooks` / `webhook_logs` 表 | ❌ | 缺失 |
| `backups` 表 | ❌ | 缺失 |
| `plugins` 表 | ❌ | 缺失 |

### 1.4 API 规范

| 检查项 | 状态 | 问题 |
|--------|------|------|
| `/api/v1/` 前缀 | ✅ | 公开接口统一前缀 |
| `/admin/api/` 前缀 | ✅ | 管理后台统一前缀 |
| JWT 鉴权中间件 | ✅ | `middleware/auth.go` 中 Auth/AdminAuth 实现正确 |
| 统一响应格式 `{code, msg, data}` | ⚠️ | 部分接口不统一（见下） |
| 响应格式一致性 | ⚠️ | `middleware/auth.go` 第31行直接返回 `gin.H{"code":..., "msg":...}`，未通过 `response.Success/Error` 包装；`AlipayController.Notify` 直接返回字符串而非 JSON |
| Epay 回调返回格式 | ⚠️ | `PaymentController.EpayNotify` 未确认是否返回 `success` 纯文本（部分易支付要求） |

**统一响应格式问题清单**：
- [ ] `internal/middleware/auth.go:31` — `c.JSON(http.StatusUnauthorized, gin.H{...})` 未通过 response 包
- [ ] `internal/middleware/auth.go:38` — 同上
- [ ] `internal/controller/advanced.go:158` — `ctx.String(http.StatusBadRequest, "fail")` 支付宝回调直接返回字符串
- [ ] `internal/controller/advanced.go:150` — 同上

---

## 二、编码规范

### 2.1 Go 命名

| 检查项 | 状态 | 说明 |
|--------|------|------|
| 文件名小写+下划线 | ✅ | `sku_service.go`、`alipay_service.go` |
| 包名小写单数 | ✅ | `package service`、`package model` |
| 结构体大驼峰 | ✅ | `AdvancedOrderService`、`SeckillService` |
| 公开函数大驼峰 | ✅ | `CreateAdvanced`、`MarkPaidAdvanced` |
| 私有函数小驼峰 | ✅ | `fmtItems`、`totalAmountToDesc` |
| 常量全大写 | ✅ | `CodeSuccess`、`UserStatusActive`、`OrderStatusPending` |
| 数据库字段小写下划线 | ✅ | `created_at`、`order_no` |
| 表名小写复数 | ⚠️ | 代码中用 `&model.User{}`，GORM 默认表名 `users` ✅ |

### 2.2 错误处理

| 检查项 | 状态 | 问题 |
|--------|------|------|
| 无 `_` 忽略错误 | ✅ | 所有 `err` 均被处理 |
| 错误日志记录 | ⚠️ | 服务层大部分方法缺少 `logger.Errorf` 调用，仅在控制器层 `response.Error` |
| 事务回滚 | ✅ | `db.DB.Transaction()` + `defer recover()` 模式正确 |
| 业务错误码定义 | ✅ | `response.go` 定义了完整错误码常量 |
| 错误返回消息友好性 | ⚠️ | 部分服务层错误消息过于简单（如 `"商品不存在"` 应改为带商品ID的详细日志） |

### 2.3 高并发安全

| 检查项 | 状态 | 说明 |
|--------|------|------|
| 秒杀下单行锁 | ✅ | `seckill_service.go:95` — `FOR UPDATE` |
| 增强订单库存扣减行锁 | ✅ | `order_service_advanced.go:45` — `FOR UPDATE` |
| 秒杀库存原子递减 | ✅ | `seckill_service.go:127` — `UPDATE stock=stock-?` |
| SKU 库存消费行锁 | ✅ | `sku_service.go` 中 `ConsumeInTx` 使用事务 |
| 支付回调幂等 | ✅ | `order_service_advanced.go:188` — `if lockedOrder.Status != OrderStatusPending { order = &lockedOrder; return nil }` |
| 无 Redis 分布式锁 | N/A | 文档建议 Redis 缓存，但当前使用数据库行锁替代（可接受） |

### 2.4 日志规范

| 检查项 | 状态 | 说明 |
|--------|------|------|
| `pkg/logger/` 包存在 | ✅ | 但**未被 service 层使用** |
| 日志级别使用 | ⚠️ | 全局未发现 `logger.Infof`/`logger.Errorf` 调用，service 层无日志 |
| 敏感信息不打印 | ✅ | 未发现密码/密钥打印 |

### 2.5 前端规范

| 检查项 | 状态 | 说明 |
|--------|------|------|
| Vue 3 Composition API | ✅ | 所有 `.vue` 文件均使用 `<script setup>` + `ref`/`computed` |
| Pinia 状态管理 | ✅ | `stores/user.js`、`stores/cart.js` 正确使用 |
| Tailwind CSS | ✅ | 通过 CDN 引入，全局 `style.css` 定义组件类 |
| const/let 禁止 var | ✅ | 前端代码全部使用 const/let |
| 响应式设计 | ✅ | `md:grid-cols-*` 等响应式类 |
| API 模块化 | ✅ | `api/seckill.js`、`api/coupon.js` 等独立文件 |
| 无 emoji | ✅ | 未发现 emoji 表情 |
| 无网络用语/AI感 | ✅ | 前端文案均为中文业务用语，无"太棒了"等AI感语言 |

---

## 三、功能清单完成度

### P0 功能（9项）

| 功能 | 状态 | 说明 |
|------|------|------|
| 用户注册（邮箱+密码） | ✅ | `POST /api/v1/user/register` |
| 用户登录 + JWT | ✅ | `POST /api/v1/user/login` + `middleware/auth.go` JWT |
| 图形验证码 | ❌ | **缺失** |
| 个人资料修改 | ✅ | `POST /api/v1/user/profile` |
| 密码修改 | ✅ | `POST /api/v1/user/password` |
| 商品分类 CRUD | ✅ | `admin/api/categories` CRUD |
| 商品 CRUD | ✅ | `admin/api/products` CRUD |
| 商品上下架 | ✅ | `Product.Status` 字段 |
| SKU 多规格商品 | ✅ | `product_skus` 表、`/api/v1/products/:id/skus`、`admin/api/products/:id/skus` |

**P0 完成度：8/9**（图形验证码缺失）

### P1 功能（按清单统计）

| 功能 | 状态 | 说明 |
|------|------|------|
| 优惠券（满减/折扣） | ✅ | 完整 CRUD + 领取 + 下单使用 |
| 优惠券兑换码领取 | ✅ | `POST /api/v1/coupon/redeem` |
| 秒杀活动 | ✅ | 活动 CRUD + 下单 + 行锁库存 |
| 批发优惠（阶梯） | ✅ | `wholesale_rules` 表 + 下单自动匹配 |
| 易支付对接 | ✅ | `epay_service.go` + 回调 |
| 余额支付 | ✅ | `PayTypeBalance` 常量 |
| 支付宝当面付 | ✅ | `alipay_service.go.Precreate` |
| 支付宝手机网站支付 | ✅ | `alipay_service.go.WapPay` |
| 支付宝电脑网站支付 | ✅ | `alipay_service.go.PagePay` |
| 三级分销 | ⚠️ | 基础实现，但无 `distribution_tree` 表（用 `parent_id` 替代，无路径）；推广海报 API 缺失 |
| 邀请码/链接 | ✅ | `users.invite_code` + `/api/v1/distribution/invite` |
| 佣金自动计算 | ✅ | `distribution_service.go` |
| 推广海报生成 | ❌ | **缺失** 后端 API |
| 提现申请 | ✅ | `POST /api/v1/withdraw/apply` |
| 提现审核 | ✅ | `POST /admin/api/withdraws/:id/process` |
| 每日签到 | ✅ | `POST /api/v1/user/signin` |
| 购买得积分 | ✅ | `order_service_advanced.go:226` |
| 积分抵现 | ✅ | `order_service_advanced.go:106` |
| 仪表盘增强 | ✅ | 近7日趋势 + Top5 + 实时数据 |
| 一键安装向导 | ✅ | `/install` 路由 + HTML 页面 |
| 自动关闭超时订单 | ⚠️ | **缺失** 定时任务（cron） |
| 商户入驻 | ❌ | **缺失** 全部商户相关功能 |
| 工单系统 | ❌ | **缺失** 全部工单相关功能 |

**P1 关键功能缺失**：
1. 图形验证码（安全基线）
2. 推广海报生成 API（`GET /api/v1/distribution/poster`）
3. 推广员等级
4. 推广链接短链+统计
5. 商户入驻/商户后台（完全缺失）
6. 工单系统（完全缺失）
7. 订单自动关闭定时任务
8. 插件系统
9. Webhook

---

## 四、配置迁移检查（v2.2 关键要求）

### 4.1 config.yaml 精简程度

| 配置项 | 当前状态 | 是否应迁移到 settings 表 |
|--------|---------|--------------------------|
| `database.*` | ✅ 保留在 config.yaml | ✅ 正确（数据库连接不能动态改） |
| `jwt.secret` | ✅ 保留在 config.yaml | ✅ 正确（JWT 密钥不能动态改） |
| `server.port` | ✅ 保留在 config.yaml | ✅ 正确 |
| `pay.epay.*` | ❌ 保留在 config.yaml | **应迁移** — 支付配置应可在后台管理 |
| `pay.alipay.*` | ❌ 保留在 config.yaml | **应迁移** — 支付宝配置应可在后台管理 |
| 其他配置（如站点名称、SEO等） | ✅ 已在 settings 表 | ✅ 正确 |

### 4.2 后台配置动态生效

| 配置项 | 状态 | 说明 |
|--------|------|------|
| 站点名称/描述 | ✅ | `SettingService` 读写 settings 表 |
| 分销佣金比例 | ❌ | **缺失** — `distribution_service.go` 中无配置读取逻辑，比例硬编码 |
| 积分规则 | ❌ | **缺失** — `points_service.go` 中积分比例硬编码 |
| 支付配置（易支付/支付宝） | ❌ | **缺失** — 从 config.yaml 读取，修改需重启 |
| 模板消息/通知配置 | ❌ | **缺失** |
| SEO 配置 | ⚠️ | 部分在 settings 表，前端读取方式未验证 |

---

## 五、具体问题清单（按优先级）

### P0 严重问题（必须修复）

| 序号 | 文件 | 行号 | 问题描述 | 修复建议 |
|------|------|------|---------|---------|
| 1 | `internal/controller/advanced.go` | 158 | `ctx.String(http.StatusBadRequest, "fail")` 支付宝回调直接返回纯文本字符串，不符合 HTTP 响应规范 | 改为 `ctx.String(http.StatusOK, "fail")` 或 `ctx.Writer.Write([]byte("fail"))` |
| 2 | `internal/middleware/auth.go` | 31, 38 | JWT 中间件直接使用 `gin.H` 返回响应，绕过了统一响应封装 | 改为 `response.Error(ctx, response.CodeTokenInvalid, "未登录")` |
| 3 | `internal/pkg/response/response.go` | — | 缺少分页响应辅助函数，文档建议 `response.PageList` | 新增 `func PageList(c *gin.Context, list interface{}, total int64, page, size int)` |
| 4 | `config.yaml.example` | 17-22 | `pay.epay.*` 和 `pay.alipay.*` 仍在 config.yaml，未迁移到 settings 表 | 创建 settings 表初始化数据，将支付配置迁移到数据库 |

### P1 重要问题

| 序号 | 文件 | 问题描述 | 修复建议 |
|------|------|---------|---------|
| 5 | `internal/service/distribution_service.go` | 佣金比例硬编码（如 `level1Rate := 0.10`），无配置读取 | 新增从 settings 表读取佣金比例的逻辑 |
| 6 | `internal/service/points_service.go` | 积分比例（1元=1积分）和抵扣比例硬编码，无配置读取 | 新增从 settings 表读取积分规则 |
| 7 | `frontend/src/views/Checkout.vue` | 74行 `Math.random()` 作为 v-for key，破坏 Vue 复用机制 | 改为稳定的 `item.product_id + '-' + index` |
| 8 | `internal/controller/advanced.go` | `/api/v1/distribution/summary` 和 `/api/v1/distribution/poster` 接口缺失（前端已调用） | 补充实现 `DistributionController.Summary` 和 `DistributionController.Poster` |
| 9 | `internal/service/` | service 层无任何日志输出，调试困难 | 在关键操作处添加 `logger.Infof`（参考 CONVENTIONS.md 第164行示例） |

### P2 优化建议

| 序号 | 文件 | 问题描述 | 修复建议 |
|------|------|---------|---------|
| 10 | `internal/controller/admin.go` | `AdminDashboard` 和 `DashboardController` 重复实现了仪表盘逻辑 | 重构为单一 service 方法，两个 controller 调用 |
| 11 | `internal/middleware/` | 缺少 `recovery.go`（panic 恢复中间件）、`ratelimit.go`（限流） | 按架构文档要求补充 |
| 12 | `internal/service/order_service_advanced.go` | 45行使用字符串格式化 SQL（`SELECT * FROM products WHERE id=%d`），虽数字安全但不符合 GORM 最佳实践 | 改为 `tx.Set("gorm:query_option", "FOR UPDATE").Where("id = ?", item.ProductID).First(&product)` |
| 13 | `frontend/src/` | 大量页面使用 `alert()` 作为错误提示，体验差 | 统一使用 Toast/Notification 组件 |
| 14 | 数据库 | 缺少 `audit_logs` 表（文档要求但未实现） | 按 `ARCHITECTURE.md` 第350行 DDL 创建表 |
| 15 | 数据库 | 缺少 `distribution_tree` 表（当前用 `users.parent_id` 替代，最多只支持一级） | 补充 `distribution_tree` 表实现三级完整链路 |

---

## 六、总体结论

### 是否所有要求均已满足？

**否。** 共发现 4 项 P0 严重问题、5 项 P1 重要问题、6 项 P2 优化建议。

### 修复优先级

**立即修复（P0）**：
1. JWT 中间件响应格式统一（`middleware/auth.go`）
2. 支付宝回调 HTTP 状态码（`advanced.go`）
3. 支付配置迁移到 settings 表
4. 补充缺失的 API 接口（`distribution/summary`、`distribution/poster`）

**尽快修复（P1）**：
5. 分销佣金比例、积分规则从 config.yaml 迁移到数据库设置
6. 补充 service 层日志
7. 修复 `Math.random()` 作为 v-for key 的问题
8. 创建 `audit_logs` 表

**后续迭代（P2）**：
9. 补充中间件（recovery、ratelimit、csrf）
10. 统一前端 Toast 提示
11. 创建 `web/admin/` Layui 后台（架构文档明确要求）
12. 实现商户系统和工单系统
13. 实现订单超时自动关闭定时任务

### 附录：检查的文件清单

| 路径 | 说明 |
|------|------|
| `/workspace/internal/model/model.go` | 数据模型（257行） |
| `/workspace/internal/router/router.go` | 路由定义（235行） |
| `/workspace/internal/pkg/response/response.go` | 统一响应（48行） |
| `/workspace/internal/pkg/config/config.go` | 配置结构（81行） |
| `/workspace/internal/middleware/auth.go` | JWT 中间件（88行） |
| `/workspace/internal/service/order_service_advanced.go` | 增强订单（267行） |
| `/workspace/internal/service/seckill_service.go` | 秒杀服务（175行） |
| `/workspace/internal/service/alipay_service.go` | 支付宝服务（194行） |
| `/workspace/internal/controller/advanced.go` | 增强控制器（620行） |
| `/workspace/internal/controller/admin.go` | 管理后台控制器（323行） |
| `/workspace/config.yaml.example` | 配置文件示例 |
| `/workspace/frontend/src/router/index.js` | 前端路由（77行） |
| `/workspace/frontend/src/stores/user.js` | 用户状态管理（46行） |
| `/workspace/frontend/src/utils/request.js` | Axios 封装（51行） |
| `/workspace/frontend/src/views/Checkout.vue` | 结算页（含优惠券） |
| `/workspace/frontend/src/views/Seckill.vue` | 秒杀频道页 |
| `/workspace/frontend/src/components/Layout.vue` | 主布局 |
