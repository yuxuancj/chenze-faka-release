// Package router 提供 HTTP 路由配置。
// Vue 3 SPA 通过 embed.FS 提供静态服务，安装向导独立嵌入。
package router

import (
	"chenze-faka/internal/controller"
	"chenze-faka/internal/middleware"
	"chenze-faka/internal/pkg/config"
	"chenze-faka/web"
	"io/fs"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

func Setup(r *gin.Engine) {
	// 未安装时，跳过 JWT secret 验证（此时 config 可能不完整）
	var jwtSecret string
	if config.AppConfig != nil {
		jwtSecret = config.AppConfig.JWT.Secret
	}

	// 安装检测中间件：未安装则强制跳转到 /install
	r.Use(installGuard())

	api := r.Group("/api/v1")
	api.Use(middleware.CORS())
	{
		userCtrl := controller.NewUserController()
		productCtrl := controller.NewProductController()
		orderCtrl := controller.NewOrderController()
		payCtrl := controller.NewPaymentController()

		api.POST("/user/register", userCtrl.Register)
		api.POST("/user/login", userCtrl.Login)
		api.GET("/products", productCtrl.List)
		api.GET("/products/:id", productCtrl.Detail)
		api.GET("/categories", productCtrl.Categories)
		api.POST("/pay/epay/notify", payCtrl.EpayNotify)
		api.GET("/pay/epay/notify", payCtrl.EpayNotify)
		api.GET("/pay/epay/return", payCtrl.EpayReturn)

		auth := api.Group("")
		auth.Use(middleware.Auth(jwtSecret))
		{
			auth.GET("/user/profile", userCtrl.Profile)
			auth.POST("/user/profile", userCtrl.UpdateProfile)
			auth.POST("/user/password", userCtrl.ChangePassword)
			auth.POST("/orders", orderCtrl.Create)
			auth.GET("/orders", orderCtrl.List)
			auth.GET("/orders/:order_no", orderCtrl.Detail)
			auth.POST("/pay", payCtrl.Pay)
		}
	}

	adminAPI := r.Group("/admin/api")
	if jwtSecret != "" {
		adminAPI.Use(middleware.AdminAuth(jwtSecret))
	}
	{
		admin := controller.NewAdminController()
		adminAPI.GET("/dashboard", admin.Dashboard)
		adminAPI.GET("/products", admin.ProductList)
		adminAPI.POST("/products", admin.ProductCreate)
		adminAPI.PUT("/products/:id", admin.ProductUpdate)
		adminAPI.DELETE("/products/:id", admin.ProductDelete)
		adminAPI.GET("/cards", admin.CardList)
		adminAPI.POST("/cards/import", admin.CardImport)
		adminAPI.GET("/categories", admin.CategoryList)
		adminAPI.POST("/categories", admin.CategoryCreate)
		adminAPI.PUT("/categories/:id", admin.CategoryUpdate)
		adminAPI.DELETE("/categories/:id", admin.CategoryDelete)
		adminAPI.GET("/orders", admin.OrderList)
		adminAPI.GET("/orders/:id", admin.OrderDetail)
		adminAPI.GET("/users", admin.UserList)
		adminAPI.PUT("/users/:id", admin.UserUpdate)
		adminAPI.GET("/settings", admin.SettingsGet)
		adminAPI.POST("/settings", admin.SettingsSet)
	}

	setupInstallAPI(r)
	setupStaticFiles(r)
}

// installGuard 中间件：检查是否已安装。
// 若 install.lock 不存在，且请求路径不是 /install 或 /install/api，则重定向到 /install。
func installGuard() gin.HandlerFunc {
	return func(c *gin.Context) {
		path := c.Request.URL.Path

		// 安装相关路径放行
		if path == "/install" || strings.HasPrefix(path, "/install/api") {
			c.Next()
			return
		}
		// API 路径放行（让后端 API 处理）
		if strings.HasPrefix(path, "/api/") || strings.HasPrefix(path, "/admin/api") {
			c.Next()
			return
		}

		// 静态资源（带扩展名）放行
		if strings.Contains(path, ".") {
			c.Next()
			return
		}

		// 检查是否已安装
		if _, err := os.Stat("install.lock"); os.IsNotExist(err) {
			// 未安装，重定向到安装页面
			if path != "/install" {
				c.Redirect(http.StatusFound, "/install")
				c.Abort()
				return
			}
		}
		c.Next()
	}
}

// setupStaticFiles 将 Vue 3 SPA 构建产物通过 embed.FS 提供静态服务。
func setupStaticFiles(r *gin.Engine) {
	subFS, err := fs.Sub(web.StaticFiles, "frontend/dist")
	if err != nil {
		r.NoRoute(func(c *gin.Context) {
			c.String(http.StatusOK, "前端资源尚未构建，请先在 frontend 目录执行 npm run build。")
		})
		return
	}

	indexHTML, readErr := fs.ReadFile(subFS, "index.html")
	if readErr != nil {
		r.NoRoute(func(c *gin.Context) {
			c.String(http.StatusOK, "index.html 不存在，请先构建前端。")
		})
		return
	}

	fileServer := http.FileServer(http.FS(subFS))

	// 独立处理 /install 路径，从嵌入文件返回安装页面
	r.GET("/install", func(c *gin.Context) {
		installHTML, err := web.InstallPage.ReadFile("install.html")
		if err != nil {
			c.String(http.StatusInternalServerError, "安装页面加载失败")
			return
		}
		c.Header("Content-Type", "text/html; charset=utf-8")
		c.String(http.StatusOK, "%s", installHTML)
	})

	r.Use(func(c *gin.Context) {
		path := c.Request.URL.Path
		if strings.HasPrefix(path, "/api/") || strings.HasPrefix(path, "/install/api") ||
			strings.HasPrefix(path, "/admin/api") || path == "/install" {
			c.Next()
			return
		}
		if strings.Contains(path, ".") {
			c.Writer.Header().Set("Cache-Control", "public, max-age=31536000, immutable")
			fileServer.ServeHTTP(c.Writer, c.Request)
			c.Abort()
			return
		}
		c.Writer.Header().Set("Content-Type", "text/html; charset=utf-8")
		c.Writer.WriteHeader(http.StatusOK)
		c.Writer.Write(indexHTML)
		c.Abort()
	})
}

// setupInstallAPI 注册安装向导 API 路由
func setupInstallAPI(r *gin.Engine) {
	install := controller.NewInstallController()
	installAPI := r.Group("/install/api")
	{
		installAPI.GET("/env", install.EnvCheck)
		installAPI.POST("/install", install.Install)
	}
}
