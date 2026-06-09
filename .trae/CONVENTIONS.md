# 晨泽发卡系统 - 编码规范

## 一、Go 代码规范

### 1.1 命名规范

| 类型 | 规范 | 示例 |
|------|------|------|
| 文件名 | 小写 + 下划线 | `user_service.go` |
| 包名 | 小写，单数 | `package model` |
| 结构体 | 大驼峰（公开） | `type UserService struct` |
| 接口 | 大驼峰，以 er 结尾 | `type Paymenter interface` |
| 公开函数/方法 | 大驼峰 | `func CreateOrder()` |
| 私有函数/方法 | 小驼峰 | `func validateParams()` |
| 常量 | 全大写 + 下划线 | `const MAX_RETRY_COUNT = 3` |
| 变量 | 小驼峰 | `var userName string` |
| 数据库字段 | 小写 + 下划线 | `created_at`, `user_id` |
| 表名 | 小写复数 | `users`, `orders` |

### 1.2 代码组织结构

```go
// 标准文件结构顺序
package service

import (
    // 标准库
    "context"
    "fmt"

    // 第三方库
    "github.com/gin-gonic/gin"
    "gorm.io/gorm"

    // 项目内部包
    "chenze-faka/internal/model"
    "chenze-faka/pkg/logger"
)

// 常量定义
const (
    DefaultPageSize = 20
)

// 变量定义
var (
    ErrUserNotFound = errors.New("user not found")
)

// 类型定义
type UserService struct {
    db *gorm.DB
}

// 构造函数
func NewUserService(db *gorm.DB) *UserService {
    return &UserService{db: db}
}

// 公开方法
func (s *UserService) Create(req *CreateUserRequest) (*model.User, error) {
    // 实现
}

// 私有方法
func (s *UserService) validateEmail(email string) bool {
    // 实现
}
```

### 1.3 错误处理规范

```go
// 错误必须处理，禁止使用 _ 忽略
result, err := doSomething()
if err != nil {
    logger.Errorf("do something failed: %v", err)
    return nil, fmt.Errorf("do something: %w", err)
}

// 业务错误使用自定义错误码
const (
    CodeSuccess          = 0
    CodeParamError       = 1001
    CodeUserNotFound     = 1002
    CodePasswordError    = 1003
    CodeTokenInvalid     = 1004
    CodePermissionDenied = 1005
    CodeOrderNotFound    = 2001
    CodeInsufficientStock = 2002
    CodePaymentFailed    = 3001
)

// 统一响应结构
type Response struct {
    Code int         `json:"code"`
    Msg  string      `json:"msg"`
    Data interface{} `json:"data,omitempty"`
}

// 成功响应
func Success(c *gin.Context, data interface{}) {
    c.JSON(http.StatusOK, Response{
        Code: CodeSuccess,
        Msg:  "success",
        Data: data,
    })
}

// 错误响应
func Error(c *gin.Context, code int, msg string) {
    c.JSON(http.StatusOK, Response{
        Code: code,
        Msg:  msg,
    })
}
```

### 1.4 数据库操作规范

```go
// 使用参数化查询（GORM 默认安全）
db.Where("user_id = ?", userID).Find(&orders)

// 事务处理
func (s *OrderService) CreateOrder(req *CreateOrderRequest) error {
    tx := s.db.Begin()
    defer func() {
        if r := recover(); r != nil {
            tx.Rollback()
        }
    }()

    // 扣减库存
    if err := tx.Model(&model.Product{}).
        Where("id = ? AND stock >= ?", req.ProductID, req.Quantity).
        Update("stock", gorm.Expr("stock - ?", req.Quantity)).Error; err != nil {
        tx.Rollback()
        return err
    }

    // 创建订单
    order := &model.Order{
        OrderNo: generateOrderNo(),
        UserID:  req.UserID,
        Amount:  req.Amount,
        Status:  model.OrderStatusPending,
    }
    if err := tx.Create(order).Error; err != nil {
        tx.Rollback()
        return err
    }

    return tx.Commit().Error
}
```

### 1.5 日志规范

```go
import "chenze-faka/pkg/logger"

// 日志级别使用
logger.Debugf("调试信息: %v", data)      // 开发环境
logger.Infof("用户登录成功: %d", userID)  // 关键操作
logger.Warnf("库存不足: product_id=%d", id) // 警告
logger.Errorf("支付回调失败: %v", err)   // 错误
logger.Fatalf("致命错误: %v", err)       // 程序退出

// 禁止在日志中打印敏感信息
// ❌ 错误示例
logger.Infof("用户密码: %s", password)
logger.Infof("支付密钥: %s", apiKey)

// ✅ 正确示例
logger.Infof("用户登录: user_id=%d", userID)
logger.Infof("支付回调: order_no=%s", orderNo)
```

### 1.6 并发安全规范

```go
// 使用 sync.Mutex 保护共享变量
type Counter struct {
    mu    sync.Mutex
    value int
}

func (c *Counter) Inc() {
    c.mu.Lock()
    defer c.mu.Unlock()
    c.value++
}

// 使用 atomic 包处理简单计数器
var counter int64
atomic.AddInt64(&counter, 1)

// 使用 channel 进行 goroutine 通信
func worker(jobs <-chan Job, results chan<- Result) {
    for job := range jobs {
        results <- process(job)
    }
}

// 使用 context 控制 goroutine 生命周期
func (s *Service) Process(ctx context.Context) error {
    ticker := time.NewTicker(5 * time.Second)
    defer ticker.Stop()

    for {
        select {
        case <-ctx.Done():
            return ctx.Err()
        case <-ticker.C:
            // 执行定时任务
        }
    }
}
```

### 1.7 注释规范

```go
// Package service 提供业务逻辑层实现
package service

// UserService 用户服务，处理用户注册、登录、信息管理
type UserService struct {
    db *gorm.DB
}

// Create 创建新用户
// 参数:
//   - req: 创建用户请求，包含邮箱、密码、昵称
// 返回:
//   - *model.User: 创建成功的用户对象
//   - error: 错误信息（邮箱已存在、参数错误等）
func (s *UserService) Create(req *CreateUserRequest) (*model.User, error) {
    // 实现
}

// 私有函数使用单行注释
// validateEmail 验证邮箱格式
func validateEmail(email string) bool {
    // 正则验证
}
```

## 二、前端规范

### 2.1 HTML/CSS 规范

```html
<!-- 使用语义化标签 -->
<header></header>
<nav></nav>
<main></main>
<aside></aside>
<footer></footer>

<!-- 类名使用中划线 -->
<div class="product-card">
    <div class="product-title"></div>
    <div class="product-price"></div>
</div>

<!-- 响应式设计 -->
<style>
    /* 移动端优先 */
    .container {
        padding: 0 12px;
    }

    /* 平板 */
    @media (min-width: 768px) {
        .container {
            max-width: 720px;
            margin: 0 auto;
        }
    }

    /* PC */
    @media (min-width: 1024px) {
        .container {
            max-width: 960px;
        }
    }
</style>
```

### 2.2 JavaScript 规范

```javascript
// 使用 const/let，禁止 var
const API_BASE = '/api/v1';
let currentPage = 1;

// 使用箭头函数
const fetchData = async (params) => {
    try {
        const response = await axios.get(`${API_BASE}/products`, { params });
        return response.data;
    } catch (error) {
        console.error('请求失败:', error);
        throw error;
    }
};

// 使用解构赋值
const { user, token } = response.data;

// 模板字符串
const message = `用户 ${user.name} 登录成功`;

// 防抖/节流
const debouncedSearch = debounce(searchHandler, 300);
```

### 2.3 Layui 规范（后台）

```javascript
// 表格渲染
layui.use(['table', 'form'], function() {
    var table = layui.table;
    var form = layui.form;

    table.render({
        elem: '#dataTable',
        url: '/admin/api/products',
        cols: [[
            {field: 'id', title: 'ID', width: 80},
            {field: 'name', title: '商品名称'},
            {field: 'price', title: '价格', width: 120},
            {field: 'status', title: '状态', width: 100, templet: '#statusTpl'},
            {title: '操作', width: 150, toolbar: '#actionBar'}
        ]],
        page: true
    });
});
```

## 三、Git 提交规范

### 3.1 提交信息格式

```text
<type>(<scope>): <subject>

<body>

<footer>
```

### 3.2 Type 类型

| 类型 | 说明 | 示例 |
|------|------|------|
| feat | 新功能 | feat(order): 新增订单导出功能 |
| fix | Bug 修复 | fix(payment): 修复支付宝回调验签失败 |
| docs | 文档更新 | docs: 更新部署文档 |
| style | 代码格式（不影响功能） | style: 格式化代码缩进 |
| refactor | 重构 | refactor(user): 重构用户模块 |
| test | 测试相关 | test: 添加分销模块单元测试 |
| chore | 构建/工具变动 | chore: 更新依赖版本 |

### 3.3 提交示例

```text
feat(product): 实现 SKU 多规格商品功能

- 新增 product_skus 表
- 商品编辑页支持动态添加规格行
- 前台商品详情支持规格选择

Closes #123
```

## 四、Git 分支规范

```text
main        # 生产分支，只接受 release 合并
develop     # 开发分支
feature/*   # 功能分支，从 develop 切出
fix/*       # 修复分支，从 develop 或 main 切出
release/*   # 发布分支
```

## 五、版本号规范

遵循语义化版本：`主版本.次版本.修订号`

- **主版本**: 不兼容的 API 修改
- **次版本**: 向下兼容的功能性新增
- **修订号**: 向下兼容的问题修正

示例：`v2.1.3`

---

文档版本: v1.0
最后更新: 2026-06-09
