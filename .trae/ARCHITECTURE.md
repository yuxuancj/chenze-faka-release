# 晨泽发卡 - 系统架构文档

## 技术栈
- 后端：Go 1.21+ / Gin / GORM
- 数据库：MySQL 5.7+
- 缓存：Redis（可选）
- 前端：后台 Layui，前台响应式 HTML/CSS

## 目录结构
chenze-faka/
├── cmd/              # 主程序入口
├── internal/         # 内部核心代码
│   ├── controller/   # API控制器
│   ├── model/        # 数据模型（GORM）
│   ├── service/      # 业务逻辑层
│   ├── middleware/   # 中间件（JWT、限流、日志）
│   └── router/       # 路由定义
├── pkg/              # 公共工具包
└── web/              # 静态资源

## 模块边界定义
- 用户模块（user）：负责注册、登录、个人中心、等级、积分
- 商品模块（product）：负责商品CRUD、SKU、分类、库存
- 订单模块（order）：负责订单创建、状态流转、支付回调
- 支付模块（payment）：负责各支付渠道对接（易支付/支付宝/PAYJS/USDT）
- 分销模块（distribution）：负责三级分销、佣金计算、提现
- 工单模块（ticket）：负责售后工单
- 插件模块（plugin）：负责插件市场、钩子机制

## 核心约束
- 各模块之间通过 Service 层调用，禁止跨模块直接操作数据库
- 所有 API 必须使用 JWT 认证（除登录/注册/公开接口）
- 数据库操作必须在事务中保证一致性（订单、库存扣减）
