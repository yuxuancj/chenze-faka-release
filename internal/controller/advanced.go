package controller

import (
	"chenze-faka/internal/model"
	"chenze-faka/internal/pkg/response"
	"chenze-faka/internal/service"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// ========== SKU 控制器 ==========

type SkuController struct{}

func NewSkuController() *SkuController { return &SkuController{} }

// GetByProduct 前台获取商品SKU
func (c *SkuController) GetByProduct(ctx *gin.Context) {
	productID, _ := strconv.Atoi(ctx.Param("id"))
	list, err := service.NewSkuService().GetSkusByProduct(uint(productID))
	if err != nil {
		response.Error(ctx, response.CodeServerError, err.Error())
		return
	}
	response.Success(ctx, list)
}

// AdminList 后台获取SKU
func (c *SkuController) AdminList(ctx *gin.Context) {
	productID, _ := strconv.Atoi(ctx.Param("id"))
	list, err := service.NewSkuService().GetSkusByProduct(uint(productID))
	if err != nil {
		response.Error(ctx, response.CodeServerError, err.Error())
		return
	}
	response.Success(ctx, list)
}

// AdminBatchUpdate 后台批量更新SKU
func (c *SkuController) AdminBatchUpdate(ctx *gin.Context) {
	productID, _ := strconv.Atoi(ctx.Param("id"))
	var req struct {
		Skus []struct {
			ID        uint                   `json:"id"`
			SkuCode   string                 `json:"sku_code"`
			SpecNames map[string]interface{} `json:"spec_names"`
			Price     float64                `json:"price"`
			Stock     int                    `json:"stock"`
			Image     string                 `json:"image"`
			Weight    float64                `json:"weight"`
			Status    int                    `json:"status"`
		} `json:"skus"`
	}
	ctx.ShouldBindJSON(&req)
	// 转换为 model.ProductSku
	skus := make([]model.ProductSku, 0, len(req.Skus))
	for _, s := range req.Skus {
		skus = append(skus, model.ProductSku{
			ID:        s.ID,
			ProductID: uint(productID),
			SkuCode:   s.SkuCode,
			SpecNames: s.SpecNames,
			Price:     s.Price,
			Stock:     s.Stock,
			Image:     s.Image,
			Weight:    s.Weight,
			Status:    s.Status,
		})
	}
	if err := service.NewSkuService().ReplaceSkus(uint(productID), skus); err != nil {
		response.Error(ctx, response.CodeServerError, err.Error())
		return
	}
	response.OK(ctx)
}

// ========== 支付宝支付控制器 ==========

type AlipayController struct{}

func NewAlipayController() *AlipayController { return &AlipayController{} }

// Pay 发起支付宝支付
func (c *AlipayController) Pay(ctx *gin.Context) {
	var req struct {
		OrderNo string `json:"order_no" form:"order_no"`
		Method  string `json:"method" form:"method"` // qr, wap, page
	}
	ctx.ShouldBind(&req)
	if req.OrderNo == "" {
		response.Error(ctx, response.CodeParamError, "订单号不能为空")
		return
	}
	order, err := service.NewOrderService().GetByOrderNo(req.OrderNo)
	if err != nil {
		response.Error(ctx, response.CodeOrderNotFound, "订单不存在")
		return
	}
	if order.Status != model.OrderStatusPending {
		response.Success(ctx, gin.H{"paid": true, "status": order.Status})
		return
	}
	alipaySvc := service.NewAlipayService()
	payReq := service.TradeCreateRequest{
		OutTradeNo:  order.OrderNo,
		TotalAmount: order.Amount,
		Subject:     "商品订单",
		Method:      req.Method,
	}
	switch req.Method {
	case "qr":
		qrContent, err := alipaySvc.Precreate(payReq)
		if err != nil {
			response.Error(ctx, response.CodePaymentFailed, err.Error())
			return
		}
		response.Success(ctx, gin.H{"qr_content": qrContent, "pay_type": "alipay_qr"})
	case "wap":
		url, err := alipaySvc.WapPay(payReq, ctx.Request.Referer()+"?order="+order.OrderNo)
		if err != nil {
			response.Error(ctx, response.CodePaymentFailed, err.Error())
			return
		}
		response.Success(ctx, gin.H{"pay_url": url, "pay_type": "alipay_wap"})
	default: // page
		url, err := alipaySvc.PagePay(payReq, ctx.Request.Referer()+"?order="+order.OrderNo)
		if err != nil {
			response.Error(ctx, response.CodePaymentFailed, err.Error())
			return
		}
		response.Success(ctx, gin.H{"pay_url": url, "pay_type": "alipay_page"})
	}
}

// Notify 支付宝异步回调
func (c *AlipayController) Notify(ctx *gin.Context) {
	_ = ctx.Request.ParseForm()
	form := ctx.Request.PostForm
	alipaySvc := service.NewAlipayService()
	orderNo, _, tradeStatus, ok := alipaySvc.VerifyNotify(form)
	if !ok {
		ctx.String(http.StatusOK, "fail")
		return
	}
	if tradeStatus != "" && !strings.Contains(tradeStatus, "SUCCESS") && !strings.Contains(tradeStatus, "TRADE_FINISHED") {
		ctx.String(http.StatusOK, "success")
		return
	}
	_, err := service.NewAdvancedOrderService().MarkPaidAdvanced(orderNo, "alipay")
	if err != nil {
		ctx.String(http.StatusInternalServerError, "fail")
		return
	}
	ctx.String(http.StatusOK, "success")
}

// ========== 优惠券控制器 ==========

type CouponController struct{}

func NewCouponController() *CouponController { return &CouponController{} }

func (c *CouponController) List(ctx *gin.Context) {
	list, err := service.NewCouponService().Available()
	if err != nil {
		response.Error(ctx, response.CodeServerError, err.Error())
		return
	}
	response.Success(ctx, list)
}

func (c *CouponController) Claim(ctx *gin.Context) {
	uid, _ := ctx.Get("user_id")
	var req struct {
		CouponID uint `json:"coupon_id" form:"coupon_id"`
	}
	ctx.ShouldBind(&req)
	if err := service.NewCouponService().Claim(uid.(uint), req.CouponID); err != nil {
		response.Error(ctx, response.CodeParamError, err.Error())
		return
	}
	response.OK(ctx)
}

func (c *CouponController) UserCoupons(ctx *gin.Context) {
	uid, _ := ctx.Get("user_id")
	status, _ := strconv.Atoi(ctx.DefaultQuery("status", "-1"))
	list, err := service.NewCouponService().UserCoupons(uid.(uint), status)
	if err != nil {
		response.Error(ctx, response.CodeServerError, err.Error())
		return
	}
	response.Success(ctx, list)
}

func (c *CouponController) AdminList(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(ctx.DefaultQuery("size", "20"))
	total, list, err := service.NewCouponService().List(page, size)
	if err != nil {
		response.Error(ctx, response.CodeServerError, err.Error())
		return
	}
	response.Success(ctx, gin.H{"total": total, "list": list, "page": page, "size": size})
}

func (c *CouponController) AdminCreate(ctx *gin.Context) {
	var req model.Coupon
	ctx.ShouldBindJSON(&req)
	if err := service.NewCouponService().Create(&req); err != nil {
		response.Error(ctx, response.CodeServerError, err.Error())
		return
	}
	response.Success(ctx, req)
}

func (c *CouponController) AdminUpdate(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.Param("id"))
	var req map[string]interface{}
	ctx.ShouldBindJSON(&req)
	if err := service.NewCouponService().Update(uint(id), req); err != nil {
		response.Error(ctx, response.CodeServerError, err.Error())
		return
	}
	response.OK(ctx)
}

func (c *CouponController) AdminDelete(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.Param("id"))
	if err := service.NewCouponService().Delete(uint(id)); err != nil {
		response.Error(ctx, response.CodeServerError, err.Error())
		return
	}
	response.OK(ctx)
}

// ========== 秒杀控制器 ==========

type SeckillController struct{}

func NewSeckillController() *SeckillController { return &SeckillController{} }

func (c *SeckillController) List(ctx *gin.Context) {
	list, err := service.NewSeckillService().ActiveList()
	if err != nil {
		response.Error(ctx, response.CodeServerError, err.Error())
		return
	}
	response.Success(ctx, gin.H{"list": list})
}

func (c *SeckillController) Order(ctx *gin.Context) {
	uid, _ := ctx.Get("user_id")
	var req struct {
		SeckillID uint   `json:"seckill_id" form:"seckill_id"`
		Quantity  int    `json:"quantity" form:"quantity"`
		Email     string `json:"email" form:"email"`
		PayType   string `json:"pay_type" form:"pay_type"`
		Remark    string `json:"remark" form:"remark"`
	}
	ctx.ShouldBind(&req)
	if req.Quantity <= 0 {
		req.Quantity = 1
	}
	order, err := service.NewSeckillService().Order(req.SeckillID, uid.(uint), req.Quantity, req.Email, req.PayType, req.Remark)
	if err != nil {
		response.Error(ctx, response.CodeParamError, err.Error())
		return
	}
	response.Success(ctx, order)
}

func (c *SeckillController) AdminList(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(ctx.DefaultQuery("size", "20"))
	total, list, err := service.NewSeckillService().List(page, size)
	if err != nil {
		response.Error(ctx, response.CodeServerError, err.Error())
		return
	}
	response.Success(ctx, gin.H{"total": total, "list": list, "page": page, "size": size})
}

func (c *SeckillController) AdminCreate(ctx *gin.Context) {
	var req model.Seckill
	ctx.ShouldBindJSON(&req)
	if err := service.NewSeckillService().Create(&req); err != nil {
		response.Error(ctx, response.CodeServerError, err.Error())
		return
	}
	response.Success(ctx, req)
}

func (c *SeckillController) AdminUpdate(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.Param("id"))
	var req map[string]interface{}
	ctx.ShouldBindJSON(&req)
	if err := service.NewSeckillService().Update(uint(id), req); err != nil {
		response.Error(ctx, response.CodeServerError, err.Error())
		return
	}
	response.OK(ctx)
}

func (c *SeckillController) AdminDelete(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.Param("id"))
	if err := service.NewSeckillService().Delete(uint(id)); err != nil {
		response.Error(ctx, response.CodeServerError, err.Error())
		return
	}
	response.OK(ctx)
}

// ========== 批发优惠控制器 ==========

type WholesaleController struct{}

func NewWholesaleController() *WholesaleController { return &WholesaleController{} }

func (c *WholesaleController) List(ctx *gin.Context) {
	productID, _ := strconv.Atoi(ctx.DefaultQuery("product_id", "0"))
	list, err := service.NewWholesaleService().List(uint(productID))
	if err != nil {
		response.Error(ctx, response.CodeServerError, err.Error())
		return
	}
	response.Success(ctx, list)
}

func (c *WholesaleController) AdminCreate(ctx *gin.Context) {
	var req model.WholesaleRule
	ctx.ShouldBindJSON(&req)
	if err := service.NewWholesaleService().Create(&req); err != nil {
		response.Error(ctx, response.CodeServerError, err.Error())
		return
	}
	response.Success(ctx, req)
}

func (c *WholesaleController) AdminUpdate(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.Param("id"))
	var req map[string]interface{}
	ctx.ShouldBindJSON(&req)
	if err := service.NewWholesaleService().Update(uint(id), req); err != nil {
		response.Error(ctx, response.CodeServerError, err.Error())
		return
	}
	response.OK(ctx)
}

func (c *WholesaleController) AdminDelete(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.Param("id"))
	if err := service.NewWholesaleService().Delete(uint(id)); err != nil {
		response.Error(ctx, response.CodeServerError, err.Error())
		return
	}
	response.OK(ctx)
}

// ========== 分销控制器 ==========

type DistributionController struct{}

func NewDistributionController() *DistributionController { return &DistributionController{} }

func (c *DistributionController) GetSummary(ctx *gin.Context) {
	uid, _ := ctx.Get("user_id")
	summary, err := service.NewDistributionService().GetSummary(uid.(uint))
	if err != nil {
		response.Error(ctx, response.CodeServerError, err.Error())
		return
	}
	response.Success(ctx, summary)
}

func (c *DistributionController) GetInviteCode(ctx *gin.Context) {
	uid, _ := ctx.Get("user_id")
	code, err := service.NewDistributionService().GetInviteCode(uid.(uint))
	if err != nil {
		response.Error(ctx, response.CodeServerError, err.Error())
		return
	}
	response.Success(ctx, gin.H{"invite_code": code})
}

func (c *DistributionController) GetCommissions(ctx *gin.Context) {
	uid, _ := ctx.Get("user_id")
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(ctx.DefaultQuery("size", "20"))
	total, list, err := service.NewDistributionService().GetCommissions(uid.(uint), page, size)
	if err != nil {
		response.Error(ctx, response.CodeServerError, err.Error())
		return
	}
	response.Success(ctx, gin.H{"total": total, "list": list, "page": page, "size": size})
}

func (c *DistributionController) GetTeam(ctx *gin.Context) {
	uid, _ := ctx.Get("user_id")
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(ctx.DefaultQuery("size", "20"))
	total, list, err := service.NewDistributionService().GetTeam(uid.(uint), page, size)
	if err != nil {
		response.Error(ctx, response.CodeServerError, err.Error())
		return
	}
	response.Success(ctx, gin.H{"total": total, "list": list, "page": page, "size": size})
}

func (c *DistributionController) ApplyWithdraw(ctx *gin.Context) {
	uid, _ := ctx.Get("user_id")
	var req struct {
		Amount      float64 `json:"amount" form:"amount"`
		AccountType string  `json:"account_type" form:"account_type"`
		Account     string  `json:"account" form:"account"`
		RealName    string  `json:"real_name" form:"real_name"`
	}
	ctx.ShouldBind(&req)
	w, err := service.NewDistributionService().ApplyWithdraw(uid.(uint), req.Amount, req.AccountType, req.Account, req.RealName)
	if err != nil {
		response.Error(ctx, response.CodeParamError, err.Error())
		return
	}
	response.Success(ctx, w)
}

func (c *DistributionController) MyWithdraws(ctx *gin.Context) {
	uid, _ := ctx.Get("user_id")
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(ctx.DefaultQuery("size", "20"))
	total, list, err := service.NewDistributionService().ListWithdraws(uid.(uint), page, size)
	if err != nil {
		response.Error(ctx, response.CodeServerError, err.Error())
		return
	}
	response.Success(ctx, gin.H{"total": total, "list": list, "page": page, "size": size})
}

func (c *DistributionController) AdminListWithdraws(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(ctx.DefaultQuery("size", "20"))
	total, list, err := service.NewDistributionService().ListWithdraws(0, page, size)
	if err != nil {
		response.Error(ctx, response.CodeServerError, err.Error())
		return
	}
	response.Success(ctx, gin.H{"total": total, "list": list, "page": page, "size": size})
}

func (c *DistributionController) AdminProcessWithdraw(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.Param("id"))
	var req struct {
		Approved bool   `json:"approved" form:"approved"`
		Remark   string `json:"remark" form:"remark"`
	}
	ctx.ShouldBind(&req)
	if err := service.NewDistributionService().ProcessWithdraw(uint(id), req.Approved, req.Remark); err != nil {
		response.Error(ctx, response.CodeParamError, err.Error())
		return
	}
	response.OK(ctx)
}

func (c *DistributionController) Poster(ctx *gin.Context) {
	var req struct {
		InviteCode string `json:"invite_code" form:"invite_code"`
	}
	ctx.ShouldBind(&req)
	if req.InviteCode == "" {
		uid, ok := ctx.Get("user_id")
		if ok {
			var err error
			req.InviteCode, err = service.NewDistributionService().GetInviteCode(uid.(uint))
			if err != nil {
				response.Error(ctx, response.CodeServerError, err.Error())
				return
			}
		}
	}
	url, err := service.NewDistributionService().GeneratePoster(req.InviteCode)
	if err != nil {
		response.Error(ctx, response.CodeServerError, err.Error())
		return
	}
	response.Success(ctx, gin.H{"poster_url": url})
}

// ========== 积分与签到控制器 ==========

type PointsController struct{}

func NewPointsController() *PointsController { return &PointsController{} }

func (c *PointsController) SignIn(ctx *gin.Context) {
	uid, _ := ctx.Get("user_id")
	points, days, err := service.NewPointsService().SignIn(uid.(uint))
	if err != nil {
		response.Error(ctx, response.CodeParamError, err.Error())
		return
	}
	response.Success(ctx, gin.H{"points": points, "continuous_days": days})
}

func (c *PointsController) SignInStatus(ctx *gin.Context) {
	uid, _ := ctx.Get("user_id")
	can := service.NewPointsService().CanSignIn(uid.(uint))
	days := service.NewPointsService().GetContinuousDays(uid.(uint))
	balance, _ := service.NewPointsService().GetBalance(uid.(uint))
	response.Success(ctx, gin.H{"can_signin": can, "continuous_days": days, "balance": balance})
}

func (c *PointsController) Balance(ctx *gin.Context) {
	uid, _ := ctx.Get("user_id")
	balance, err := service.NewPointsService().GetBalance(uid.(uint))
	if err != nil {
		response.Error(ctx, response.CodeServerError, err.Error())
		return
	}
	response.Success(ctx, gin.H{"balance": balance})
}

func (c *PointsController) Logs(ctx *gin.Context) {
	uid, _ := ctx.Get("user_id")
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(ctx.DefaultQuery("size", "20"))
	total, list, err := service.NewPointsService().GetLogs(uid.(uint), page, size)
	if err != nil {
		response.Error(ctx, response.CodeServerError, err.Error())
		return
	}
	response.Success(ctx, gin.H{"total": total, "list": list, "page": page, "size": size})
}

// ========== 增强版订单控制器（支持多商品/SKU/优惠券/积分） ==========

type AdvancedOrderController struct{}

func NewAdvancedOrderController() *AdvancedOrderController { return &AdvancedOrderController{} }

func (c *AdvancedOrderController) Create(ctx *gin.Context) {
	uid, ok := ctx.Get("user_id")
	var userID uint
	if ok {
		userID = uid.(uint)
	}
	var req struct {
		Items []struct {
			ProductID uint `json:"product_id"`
			SkuID     uint `json:"sku_id"`
			Quantity  int  `json:"quantity"`
		} `json:"items"`
		CouponID  uint   `json:"coupon_id"`
		UsePoints int    `json:"use_points"`
		PayType   string `json:"pay_type"`
		Email     string `json:"email"`
		Remark    string `json:"remark"`
	}
	ctx.ShouldBind(&req)
	if len(req.Items) == 0 {
		response.Error(ctx, response.CodeParamError, "订单商品为空")
		return
	}
	items := make([]service.OrderItem, 0, len(req.Items))
	for _, it := range req.Items {
		items = append(items, service.OrderItem{
			ProductID: it.ProductID,
			SkuID:     it.SkuID,
			Quantity:  it.Quantity,
		})
	}
	order, err := service.NewAdvancedOrderService().CreateAdvanced(userID, items, req.CouponID, req.UsePoints, req.Email, req.PayType, req.Remark)
	if err != nil {
		response.Error(ctx, response.CodeParamError, err.Error())
		return
	}
	response.Success(ctx, order)
}

func (c *AdvancedOrderController) MarkPaid(ctx *gin.Context) {
	var req struct {
		OrderNo string `json:"order_no" form:"order_no"`
		PayType string `json:"pay_type" form:"pay_type"`
	}
	ctx.ShouldBind(&req)
	if req.PayType == "" {
		req.PayType = "alipay"
	}
	order, err := service.NewAdvancedOrderService().MarkPaidAdvanced(req.OrderNo, req.PayType)
	if err != nil {
		response.Error(ctx, response.CodeServerError, err.Error())
		return
	}
	response.Success(ctx, order)
}

// ========== 仪表盘增强 ==========

type DashboardController struct{}

func NewDashboardController() *DashboardController { return &DashboardController{} }

func (c *DashboardController) AdminDashboard(ctx *gin.Context) {
	dbInst := getDB()
	if dbInst == nil {
		response.Error(ctx, response.CodeServerError, "数据库未连接")
		return
	}
	var userCount, productCount, orderCount int64
	var totalAmount float64
	dbInst.Model(&model.User{}).Count(&userCount)
	dbInst.Model(&model.Product{}).Count(&productCount)
	dbInst.Model(&model.Order{}).Count(&orderCount)
	dbInst.Model(&model.Order{}).Where("status IN ?", []int{1, 2}).Select("IFNULL(SUM(amount),0)").Scan(&totalAmount)

	// 今日数据
	type dailyStat struct {
		Date  string  `json:"date"`
		Count int64   `json:"count"`
		Total float64 `json:"total"`
	}
	var last7Days []dailyStat
	dbInst.Model(&model.Order{}).Where("status IN ?", []int{1, 2}).
		Select("DATE(created_at) as date, COUNT(*) as count, IFNULL(SUM(amount),0) as total").
		Group("DATE(created_at)").Order("date desc").Limit(7).Scan(&last7Days)

	// 商品销售TOP5
	type productRank struct {
		Name   string  `json:"name"`
		Sales  int     `json:"sales"`
		Amount float64 `json:"amount"`
	}
	var topProducts []productRank
	dbInst.Model(&model.Order{}).Where("status IN ?", []int{1, 2}).
		Select("product_snapshot as name, SUM(quantity) as sales, IFNULL(SUM(amount),0) as amount").
		Group("product_id, product_snapshot").Order("amount desc").Limit(5).Scan(&topProducts)

	response.Success(ctx, gin.H{
		"user_count":     userCount,
		"product_count":  productCount,
		"order_count":    orderCount,
		"total_amount":   totalAmount,
		"last_7_days":    last7Days,
		"top_products":   topProducts,
	})
}

func getDB() *gorm.DB {
	// 使用 service 层获取数据库连接（通过内部包引用）
	return service.GetDB()
}
