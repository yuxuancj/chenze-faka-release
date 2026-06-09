package router

import (
	"chenze-faka/internal/controller"
	"chenze-faka/internal/middleware"
	"chenze-faka/internal/pkg/config"
	"chenze-faka/internal/service"
	"chenze-faka/web"
	"html/template"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Setup(r *gin.Engine) {
	tmpl, err := template.ParseFS(web.StaticFiles, "templates/**/*.html")
	if err != nil {
		panic("parse templates from embed: " + err.Error())
	}
	r.SetHTMLTemplate(tmpl)

	jwtSecret := config.AppConfig.JWT.Secret

	// CORS for API
	api := r.Group("/api/v1")
	api.Use(middleware.CORS())
	{
		userCtrl := controller.NewUserController()
		productCtrl := controller.NewProductController()
		orderCtrl := controller.NewOrderController()
		payCtrl := controller.NewPaymentController()

		// 公开接口
		api.POST("/user/register", userCtrl.Register)
		api.POST("/user/login", userCtrl.Login)
		api.GET("/products", productCtrl.List)
		api.GET("/products/:id", productCtrl.Detail)
		api.GET("/categories", productCtrl.Categories)

		// 支付回调（必须公开）
		api.POST("/pay/epay/notify", payCtrl.EpayNotify)
		api.GET("/pay/epay/notify", payCtrl.EpayNotify)
		api.GET("/pay/epay/return", payCtrl.EpayReturn)

		// 需要登录
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

	// 后台管理 API
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

		adminAPI.GET("/settings", admin.SettingsGet)
		adminAPI.POST("/settings", admin.SettingsSet)
	}

	setupPages(r)
}

func setupPages(r *gin.Engine) {
	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "templates/index.html", commonData(c))
	})
	r.GET("/products", func(c *gin.Context) {
		c.HTML(http.StatusOK, "templates/product/list.html", commonData(c))
	})
	r.GET("/product/:id", func(c *gin.Context) {
		data := commonData(c)
		data["product_id"] = c.Param("id")
		c.HTML(http.StatusOK, "templates/product/detail.html", data)
	})
	r.GET("/user/login", func(c *gin.Context) {
		c.HTML(http.StatusOK, "templates/user/login.html", commonData(c))
	})
	r.GET("/user/register", func(c *gin.Context) {
		c.HTML(http.StatusOK, "templates/user/register.html", commonData(c))
	})
	r.GET("/user/profile", func(c *gin.Context) {
		c.HTML(http.StatusOK, "templates/user/profile.html", commonData(c))
	})
	r.GET("/user/orders", func(c *gin.Context) {
		c.HTML(http.StatusOK, "templates/user/orders.html", commonData(c))
	})
	r.GET("/order/:order_no", func(c *gin.Context) {
		data := commonData(c)
		data["order_no"] = c.Param("order_no")
		c.HTML(http.StatusOK, "templates/order/detail.html", data)
	})

	admin := r.Group("/admin")
	{
		admin.GET("/login", func(c *gin.Context) {
			c.HTML(http.StatusOK, "templates/admin/login.html", commonData(c))
		})
		admin.GET("", func(c *gin.Context) {
			c.Redirect(http.StatusFound, "/admin/")
		})
		admin.GET("/", func(c *gin.Context) {
			c.HTML(http.StatusOK, "templates/admin/dashboard.html", commonData(c))
		})
		admin.GET("/products", func(c *gin.Context) {
			c.HTML(http.StatusOK, "templates/admin/products.html", commonData(c))
		})
		admin.GET("/product/new", func(c *gin.Context) {
			c.HTML(http.StatusOK, "templates/admin/product_edit.html", commonData(c))
		})
		admin.GET("/product/:id", func(c *gin.Context) {
			data := commonData(c)
			data["product_id"] = c.Param("id")
			c.HTML(http.StatusOK, "templates/admin/product_edit.html", data)
		})
		admin.GET("/cards", func(c *gin.Context) {
			c.HTML(http.StatusOK, "templates/admin/cards.html", commonData(c))
		})
		admin.GET("/categories", func(c *gin.Context) {
			c.HTML(http.StatusOK, "templates/admin/categories.html", commonData(c))
		})
		admin.GET("/orders", func(c *gin.Context) {
			c.HTML(http.StatusOK, "templates/admin/orders.html", commonData(c))
		})
		admin.GET("/users", func(c *gin.Context) {
			c.HTML(http.StatusOK, "templates/admin/users.html", commonData(c))
		})
	}

	// 安装向导（嵌入二进制）
	r.GET("/install", func(c *gin.Context) {
		c.HTML(http.StatusOK, "templates/install/index.html", commonData(c))
	})

	// 404
	r.NoRoute(func(c *gin.Context) {
		c.HTML(http.StatusNotFound, "templates/404.html", commonData(c))
	})
}

func commonData(c *gin.Context) map[string]interface{} {
	ss := service.NewSettingService()
	return map[string]interface{}{
		"site_name": ss.Get("site_name", "晨泽发卡"),
		"site_desc": ss.Get("site_desc", ""),
	}
}
