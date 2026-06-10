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
	"time"

	"github.com/gin-gonic/gin"
)

func Setup(r *gin.Engine) {
	var jwtSecret string
	if config.AppConfig != nil {
		jwtSecret = config.AppConfig.JWT.Secret
	}

	r.Use(installGuard())
	r.Use(middleware.Recovery())

	api := r.Group("/api/v1")
	api.Use(middleware.CORS())
	{
		userCtrl := controller.NewUserController()
		productCtrl := controller.NewProductController()
		orderCtrl := controller.NewOrderController()
		payCtrl := controller.NewPaymentController()

		skuCtrl := controller.NewSkuController()
		alipayCtrl := controller.NewAlipayController()
		couponCtrl := controller.NewCouponController()
		seckillCtrl := controller.NewSeckillController()
		wholesaleCtrl := controller.NewWholesaleController()
		distCtrl := controller.NewDistributionController()
		pointsCtrl := controller.NewPointsController()
		advOrderCtrl := controller.NewAdvancedOrderController()

		// 公开接口
		api.POST("/user/register", userCtrl.Register)
		api.POST("/user/login", userCtrl.Login)
		api.GET("/products", productCtrl.List)
		api.GET("/products/:id", productCtrl.Detail)
		api.GET("/products/:id/skus", skuCtrl.GetByProduct)
		api.GET("/categories", productCtrl.Categories)
		api.GET("/pay/epay/notify", payCtrl.EpayNotify)
		api.POST("/pay/epay/notify", payCtrl.EpayNotify)
		api.GET("/pay/epay/return", payCtrl.EpayReturn)
		api.POST("/pay/alipay/notify", alipayCtrl.Notify)
		api.GET("/seckill/activities", seckillCtrl.List)
		api.GET("/wholesale/rules", wholesaleCtrl.List)
		api.GET("/coupons", couponCtrl.List)

		// 需要登录的接口
		auth := api.Group("")
		auth.Use(middleware.Auth(jwtSecret))
		{
			auth.GET("/user/profile", userCtrl.Profile)
			auth.POST("/user/profile", userCtrl.UpdateProfile)
			auth.POST("/user/password", userCtrl.ChangePassword)
			auth.POST("/orders", middleware.RateLimit(5, time.Minute), orderCtrl.Create)
			auth.POST("/orders/advanced", middleware.RateLimit(5, time.Minute), advOrderCtrl.Create)
			auth.POST("/orders/markpaid", advOrderCtrl.MarkPaid)
			auth.GET("/orders", orderCtrl.List)
			auth.GET("/orders/:order_no", orderCtrl.Detail)
			auth.POST("/pay", middleware.RateLimit(10, time.Minute), payCtrl.Pay)
			auth.POST("/pay/alipay", middleware.RateLimit(10, time.Minute), alipayCtrl.Pay)

			// 优惠券
			auth.POST("/coupon/redeem", couponCtrl.Claim)
			auth.GET("/user/coupons", couponCtrl.UserCoupons)

			// 秒杀下单
			auth.POST("/seckill/order", seckillCtrl.Order)

			// 分销
			auth.GET("/distribution/summary", distCtrl.GetSummary)
			auth.GET("/distribution/invite", distCtrl.GetInviteCode)
			auth.GET("/distribution/commissions", distCtrl.GetCommissions)
			auth.GET("/distribution/team", distCtrl.GetTeam)
			auth.GET("/distribution/poster", distCtrl.Poster)
			auth.POST("/withdraw/apply", distCtrl.ApplyWithdraw)
			auth.GET("/withdraw/list", distCtrl.MyWithdraws)

			// 积分签到
			auth.POST("/user/signin", pointsCtrl.SignIn)
			auth.GET("/user/signin/status", pointsCtrl.SignInStatus)
			auth.GET("/user/points", pointsCtrl.Balance)
			auth.GET("/user/points/logs", pointsCtrl.Logs)
		}
	}

	// 管理后台API
	adminAPI := r.Group("/admin/api")
	if jwtSecret != "" {
		adminAPI.Use(middleware.AdminAuth(jwtSecret))
	}
	{
		admin := controller.NewAdminController()
		skuCtrl := controller.NewSkuController()
		couponCtrl := controller.NewCouponController()
		seckillCtrl := controller.NewSeckillController()
		wholesaleCtrl := controller.NewWholesaleController()
		distCtrl := controller.NewDistributionController()
		dashCtrl := controller.NewDashboardController()

		adminAPI.GET("/dashboard", dashCtrl.AdminDashboard)
		adminAPI.GET("/products", admin.ProductList)
		adminAPI.POST("/products", admin.ProductCreate)
		adminAPI.PUT("/products/:id", admin.ProductUpdate)
		adminAPI.DELETE("/products/:id", admin.ProductDelete)
		adminAPI.GET("/products/:id/skus", skuCtrl.AdminList)
		adminAPI.POST("/products/:id/skus", skuCtrl.AdminBatchUpdate)
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

		// 优惠券管理
		adminAPI.GET("/coupons", couponCtrl.AdminList)
		adminAPI.POST("/coupons", couponCtrl.AdminCreate)
		adminAPI.PUT("/coupons/:id", couponCtrl.AdminUpdate)
		adminAPI.DELETE("/coupons/:id", couponCtrl.AdminDelete)

		// 秒杀管理
		adminAPI.GET("/seckills", seckillCtrl.AdminList)
		adminAPI.POST("/seckills", seckillCtrl.AdminCreate)
		adminAPI.PUT("/seckills/:id", seckillCtrl.AdminUpdate)
		adminAPI.DELETE("/seckills/:id", seckillCtrl.AdminDelete)

		// 批发规则
		adminAPI.POST("/wholesale/rules", wholesaleCtrl.AdminCreate)
		adminAPI.PUT("/wholesale/rules/:id", wholesaleCtrl.AdminUpdate)
		adminAPI.DELETE("/wholesale/rules/:id", wholesaleCtrl.AdminDelete)

		// 提现审核
		adminAPI.GET("/withdraws", distCtrl.AdminListWithdraws)
		adminAPI.POST("/withdraws/:id/process", distCtrl.AdminProcessWithdraw)
	}

	setupInstallAPI(r)
	setupStaticFiles(r)
}

func installGuard() gin.HandlerFunc {
	return func(c *gin.Context) {
		path := c.Request.URL.Path
		if path == "/install" || strings.HasPrefix(path, "/install/api") {
			c.Next()
			return
		}
		if strings.HasPrefix(path, "/api/") || strings.HasPrefix(path, "/admin/api") {
			c.Next()
			return
		}
		if strings.Contains(path, ".") {
			c.Next()
			return
		}
		if _, err := os.Stat("install.lock"); os.IsNotExist(err) {
			if path != "/install" {
				c.Redirect(http.StatusFound, "/install")
				c.Abort()
				return
			}
		}
		c.Next()
	}
}

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

func setupInstallAPI(r *gin.Engine) {
	install := controller.NewInstallController()
	installAPI := r.Group("/install/api")
	{
		installAPI.GET("/env", install.EnvCheck)
		installAPI.POST("/install", install.Install)
	}
}
