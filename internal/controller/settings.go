package controller

import (
	"chenze-faka/config"
	"chenze-faka/internal/model"
	"chenze-faka/internal/pkg/db"
	"chenze-faka/internal/pkg/response"
	"chenze-faka/internal/service"
	"database/sql"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type SettingController struct{}

func NewSettingController() *SettingController {
	return &SettingController{}
}

func (s *SettingController) GetSettings(c *gin.Context) {
	group := c.Param("group")
	
	ss := service.NewSettingService()
	settings := make(map[string]interface{})

	switch group {
	case "basic":
		settings["site_name"] = ss.GetString("site_name", "晨泽发卡")
		settings["site_logo"] = ss.GetString("site_logo", "/static/logo.png")
		settings["site_favicon"] = ss.GetString("site_favicon", "")
		settings["site_description"] = ss.GetString("site_description", "专业的数字商品发卡平台")
		settings["seo_title"] = ss.GetString("seo_title", "")
		settings["seo_keywords"] = ss.GetString("seo_keywords", "")
		settings["seo_description"] = ss.GetString("seo_description", "")
		settings["footer_html"] = ss.GetString("footer_html", "")
		settings["site_closed"] = ss.GetBool("site_closed", false)
		
	case "payment":
		settings["epay_enabled"] = ss.GetBool("epay_enabled", false)
		settings["epay_url"] = ss.GetString("epay_url", "")
		settings["epay_mchid"] = ss.GetString("epay_mchid", "")
		settings["epay_key"] = ss.GetString("epay_key", "")
		
		settings["alipay_enabled"] = ss.GetBool("alipay_enabled", false)
		settings["alipay_appid"] = ss.GetString("alipay_appid", "")
		settings["alipay_private_key"] = ss.GetString("alipay_private_key", "")
		settings["alipay_public_key"] = ss.GetString("alipay_public_key", "")
		settings["alipay_sandbox"] = ss.GetBool("alipay_sandbox", false)
		
		settings["wechat_enabled"] = ss.GetBool("wechat_enabled", false)
		settings["wechat_appid"] = ss.GetString("wechat_appid", "")
		settings["wechat_mchid"] = ss.GetString("wechat_mchid", "")
		settings["wechat_key"] = ss.GetString("wechat_key", "")
		
		settings["balance_enabled"] = ss.GetBool("balance_enabled", true)
		settings["balance_recharge_rules"] = ss.GetString("balance_recharge_rules", "100:10")
		
	case "mail":
		settings["smtp_host"] = ss.GetString("smtp_host", "")
		settings["smtp_port"] = ss.GetInt("smtp_port", 465)
		settings["smtp_email"] = ss.GetString("smtp_email", "")
		settings["smtp_name"] = ss.GetString("smtp_name", "")
		settings["smtp_password"] = ss.GetString("smtp_password", "")
		settings["smtp_encryption"] = ss.GetString("smtp_encryption", "SSL")
		
	case "distribution":
		settings["distrib_level1_rate"] = ss.GetInt("distrib_level1_rate", 10)
		settings["distrib_level2_rate"] = ss.GetInt("distrib_level2_rate", 5)
		settings["distrib_level3_rate"] = ss.GetInt("distrib_level3_rate", 2)
		settings["distrib_min_withdraw"] = ss.GetFloat("distrib_min_withdraw", 10)
		settings["distrib_withdraw_fee_rate"] = ss.GetInt("distrib_withdraw_fee_rate", 1)
		settings["distrib_audit_required"] = ss.GetBool("distrib_audit_required", true)
		settings["distrib_short_domain"] = ss.GetString("distrib_short_domain", "")
		
	case "points":
		settings["points_per_yuan"] = ss.GetInt("points_per_yuan", 10)
		settings["points_deduct_ratio"] = ss.GetInt("points_deduct_ratio", 100)
		settings["points_max_deduct_percent"] = ss.GetInt("points_max_deduct_percent", 30)
		settings["signin_points"] = ss.GetInt("signin_points", 10)
		settings["signin_continuous_reward"] = ss.GetString("signin_continuous_reward", `{"3":5,"7":10}`)
		
	case "seckill":
		settings["seckill_order_timeout"] = ss.GetInt("seckill_order_timeout", 5)
		settings["seckill_rate_limit"] = ss.GetInt("seckill_rate_limit", 100)
		
	case "security":
		settings["captcha_enabled"] = ss.GetBool("captcha_enabled", true)
		settings["ip_rate_limit"] = ss.GetInt("ip_rate_limit", 60)
		settings["csrf_enabled"] = ss.GetBool("csrf_enabled", true)
		settings["two_factor_enabled"] = ss.GetBool("two_factor_enabled", false)
		
	case "other":
		settings["order_timeout_minutes"] = ss.GetInt("order_timeout_minutes", 30)
		settings["auto_backup_enabled"] = ss.GetBool("auto_backup_enabled", false)
		settings["backup_time"] = ss.GetString("backup_time", "02:00")
		
	case "template":
		settings["frontend_template"] = ss.GetString("frontend_template", "default")
		settings["backend_template"] = ss.GetString("backend_template", "default")
		
	default:
		response.Error(c, response.CodeParamError, "未知的设置分组")
		return
	}

	response.Success(c, settings)
}

func (s *SettingController) SaveSettings(c *gin.Context) {
	group := c.Param("group")
	var data map[string]interface{}
	c.ShouldBindJSON(&data)

	ss := service.NewSettingService()
	err := ss.BatchSet(data)
	if err != nil {
		response.Error(c, response.CodeServerError, err.Error())
		return
	}

	response.OK(c)
}

type DashboardController struct{}

func NewDashboardController() *DashboardController {
	return &DashboardController{}
}

func (dc *DashboardController) GetDashboard(c *gin.Context) {
	now := time.Now()
	todayStart := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	yesterdayStart := todayStart.AddDate(0, 0, -1)
	weekStart := todayStart.AddDate(0, 0, -int(now.Weekday()))

	var todayOrders, yesterdayOrders, weekOrders int64
	var todaySales, yesterdaySales, weekSales float64

	db := service.GetDB()
	
	db.Model(&model.Order{}).Where("created_at >= ?", todayStart).Count(&todayOrders)
	db.Model(&model.Order{}).Where("created_at >= ? AND created_at < ?", yesterdayStart, todayStart).Count(&yesterdayOrders)
	db.Model(&model.Order{}).Where("created_at >= ?", weekStart).Count(&weekOrders)
	
	db.Model(&model.Order{}).Where("created_at >= ?", todayStart).Select("SUM(total_amount)").Scan(&todaySales)
	db.Model(&model.Order{}).Where("created_at >= ? AND created_at < ?", yesterdayStart, todayStart).Select("SUM(total_amount)").Scan(&yesterdaySales)
	db.Model(&model.Order{}).Where("created_at >= ?", weekStart).Select("SUM(total_amount)").Scan(&weekSales)

	var todayNewUsers, totalUsers, totalProducts, totalOrders int64
	var totalSales float64

	db.Model(&model.User{}).Where("created_at >= ?", todayStart).Count(&todayNewUsers)
	db.Model(&model.User{}).Count(&totalUsers)
	db.Model(&model.Product{}).Count(&totalProducts)
	db.Model(&model.Order{}).Count(&totalOrders)
	db.Model(&model.Order{}).Select("SUM(total_amount)").Scan(&totalSales)

	last7DaysOrders := make([]int, 7)
	last7DaysSales := make([]float64, 7)
	for i := 6; i >= 0; i-- {
		dayStart := todayStart.AddDate(0, 0, -i)
		dayEnd := dayStart.AddDate(0, 0, 1)
		var count int64
		var amount float64
		db.Model(&model.Order{}).Where("created_at >= ? AND created_at < ?", dayStart, dayEnd).Count(&count)
		db.Model(&model.Order{}).Where("created_at >= ? AND created_at < ?", dayStart, dayEnd).Select("SUM(total_amount)").Scan(&amount)
		last7DaysOrders[6-i] = int(count)
		last7DaysSales[6-i] = amount
	}

	type productRank struct {
		Name   string  `json:"name"`
		Sales  int     `json:"sales"`
		Amount float64 `json:"amount"`
	}
	var productRankList []productRank
	db.Raw(`SELECT p.name, SUM(oi.quantity) as sales, SUM(oi.price * oi.quantity) as amount 
			FROM order_items oi 
			LEFT JOIN products p ON oi.product_id = p.id 
			GROUP BY oi.product_id 
			ORDER BY sales DESC 
			LIMIT 10`).Scan(&productRankList)

	type paymentRatio struct {
		Method  string  `json:"method"`
		Percent float64 `json:"percent"`
	}
	var paymentRatioList []paymentRatio
	db.Raw(`SELECT pay_type as method, COUNT(*) as count 
			FROM orders 
			WHERE pay_type IS NOT NULL 
			GROUP BY pay_type`).Scan(&paymentRatioList)
	
	totalPayCount := 0
	for _, pr := range paymentRatioList {
		totalPayCount += int(pr.Percent)
	}
	if totalPayCount > 0 {
		for i := range paymentRatioList {
			paymentRatioList[i].Percent = float64(paymentRatioList[i].Percent) / float64(totalPayCount) * 100
		}
	}

	var pendingOrders, pendingWithdraws int64
	db.Model(&model.Order{}).Where("status = ?", "pending").Count(&pendingOrders)
	
	var lowStockCount int64
	db.Model(&model.Product{}).Where("stock < ?", 10).Count(&lowStockCount)

	var goVersion string
	dbVersion := "MySQL"
	
	resp := gin.H{
		"today_orders":      todayOrders,
		"today_sales":       todaySales,
		"today_new_users":   todayNewUsers,
		"yesterday_orders":  yesterdayOrders,
		"yesterday_sales":   yesterdaySales,
		"week_orders":       weekOrders,
		"week_sales":        weekSales,
		"total_users":       totalUsers,
		"total_products":    totalProducts,
		"total_orders":      totalOrders,
		"total_sales":       totalSales,
		"last_7_days_orders": last7DaysOrders,
		"last_7_days_sales":  last7DaysSales,
		"product_rank":       productRankList,
		"payment_ratio":      paymentRatioList,
		"version":            config.Version,
		"go_version":         goVersion,
		"db_version":         dbVersion,
		"server_time":        time.Now().Format("2006-01-02 15:04:05"),
		"pending_orders":     pendingOrders,
		"pending_withdraws":  pendingWithdraws,
		"low_stock_count":    lowStockCount,
	}

	response.Success(c, resp)
}

func (s *SettingController) TestEmail(c *gin.Context) {
	var data struct {
		ToEmail string `json:"to_email"`
	}
	c.ShouldBindJSON(&data)
	
	ss := service.NewSettingService()
	err := service.NewEmailService().SendTestEmail(data.ToEmail, ss)
	if err != nil {
		response.Error(c, response.CodeServerError, err.Error())
		return
	}
	response.OK(c)
}