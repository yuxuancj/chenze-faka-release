// Package router 提供 HTTP 路由配置。
// Vue 3 SPA 通过 embed.FS 提供静态服务，所有非 API 路径返回 index.html。
package router

import (
	"chenze-faka/internal/controller"
	"chenze-faka/internal/middleware"
	"chenze-faka/internal/pkg/config"
	"chenze-faka/web"
	"io/fs"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func Setup(r *gin.Engine) {
	jwtSecret := config.AppConfig.JWT.Secret

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
	adminAPI.Use(middleware.AdminAuth(jwtSecret))
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

// setupStaticFiles 将 Vue 3 SPA 构建产物通过 embed.FS 提供静态服务。
// 对于包含扩展名的请求（如 /assets/app.js），直接返回静态文件；
// 其他所有路径（API 和 /install/api 除外）返回 index.html。
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

	r.Use(func(c *gin.Context) {
		path := c.Request.URL.Path
		if strings.HasPrefix(path, "/api/") || strings.HasPrefix(path, "/install/api") ||
			strings.HasPrefix(path, "/admin/api") {
			c.Next()
			return
		}
		// 路径包含扩展名，作为静态资源处理
		if strings.Contains(path, ".") {
			c.Writer.Header().Set("Cache-Control", "public, max-age=31536000, immutable")
			fileServer.ServeHTTP(c.Writer, c.Request)
			c.Abort()
			return
		}
		// SPA 路由，返回 index.html
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
