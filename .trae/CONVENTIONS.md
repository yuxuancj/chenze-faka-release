# 晨泽发卡 - 编码规范

## 命名规范
- 文件名：小写下划线（user_service.go）
- 结构体：大驼峰（UserService）
- 方法：大驼峰公开、小驼峰私有（CreateOrder / validateOrder）
- 常量：全大写下划线（MAX_RETRY_COUNT）

## 错误处理
- 所有 error 必须处理，不允许使用 `_` 忽略
- 业务错误使用自定义错误码（如 1001=用户不存在）
- API 返回格式统一为 `{code, msg, data}`

## 数据库规范
- 表名：小写复数（users, orders）
- 字段：下划线命名（created_at, user_id）
- 必须包含 created_at、updated_at 字段

## Git 提交规范
- 格式：`<type>: <subject>`
- type 类型：feat（新功能）、fix（修复）、docs（文档）、refactor（重构）、test（测试）
- 示例：`feat: 新增支付宝当面付支付接口`
