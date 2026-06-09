package controller

import (
	"chenze-faka/internal/pkg/config"
	"chenze-faka/internal/pkg/jwt"
	"chenze-faka/internal/pkg/response"
	"chenze-faka/internal/service"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// ========== Frontend User ==========

type UserController struct{}

func NewUserController() *UserController { return &UserController{} }

func (c *UserController) Register(ctx *gin.Context) {
	var req struct {
		Email    string `json:"email" form:"email"`
		Password string `json:"password" form:"password"`
		Nickname string `json:"nickname" form:"nickname"`
	}
	if err := ctx.ShouldBind(&req); err != nil {
		response.Error(ctx, response.CodeParamError, "参数错误")
		return
	}
	if req.Email == "" || len(req.Password) < 6 {
		response.Error(ctx, response.CodeParamError, "邮箱或密码格式错误")
		return
	}
	u, err := service.NewUserService().Register(req.Email, req.Password, req.Nickname)
	if err != nil {
		response.Error(ctx, response.CodeParamError, err.Error())
		return
	}
	response.Success(ctx, u)
}

func (c *UserController) Login(ctx *gin.Context) {
	var req struct {
		Email    string `json:"email" form:"email"`
		Password string `json:"password" form:"password"`
	}
	if err := ctx.ShouldBind(&req); err != nil {
		response.Error(ctx, response.CodeParamError, "参数错误")
		return
	}
	u, err := service.NewUserService().Login(req.Email, req.Password)
	if err != nil {
		response.Error(ctx, response.CodePasswordError, err.Error())
		return
	}
	token, _ := jwt.Generate(u.ID, u.Email, u.IsAdmin, config.AppConfig.JWT.Secret, config.AppConfig.JWT.Expire)
	response.Success(ctx, gin.H{
		"user":  u,
		"token": token,
	})
}

func (c *UserController) Profile(ctx *gin.Context) {
	uid, _ := ctx.Get("user_id")
	u, err := service.NewUserService().GetByID(uid.(uint))
	if err != nil {
		response.Error(ctx, response.CodeUserNotFound, "用户不存在")
		return
	}
	response.Success(ctx, u)
}

func (c *UserController) UpdateProfile(ctx *gin.Context) {
	uid, _ := ctx.Get("user_id")
	var req struct {
		Nickname string `json:"nickname" form:"nickname"`
	}
	ctx.ShouldBind(&req)
	service.NewUserService().UpdateProfile(uid.(uint), req.Nickname)
	response.OK(ctx)
}

func (c *UserController) ChangePassword(ctx *gin.Context) {
	uid, _ := ctx.Get("user_id")
	var req struct {
		OldPwd string `json:"old_password" form:"old_password"`
		NewPwd string `json:"new_password" form:"new_password"`
	}
	ctx.ShouldBind(&req)
	err := service.NewUserService().ChangePassword(uid.(uint), req.OldPwd, req.NewPwd)
	if err != nil {
		response.Error(ctx, response.CodeParamError, err.Error())
		return
	}
	response.OK(ctx)
}

// ========== Product (Frontend) ==========

type ProductController struct{}

func NewProductController() *ProductController { return &ProductController{} }

func (c *ProductController) List(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(ctx.DefaultQuery("size", "20"))
	catID, _ := strconv.Atoi(ctx.DefaultQuery("category_id", "0"))
	keyword := ctx.DefaultQuery("keyword", "")
	total, list, err := service.NewProductService().List(page, size, uint(catID), keyword)
	if err != nil {
		response.Error(ctx, response.CodeServerError, err.Error())
		return
	}
	response.Success(ctx, gin.H{"total": total, "list": list, "page": page, "size": size})
}

func (c *ProductController) Detail(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.Param("id"))
	p, err := service.NewProductService().GetByID(uint(id))
	if err != nil {
		response.Error(ctx, response.CodeProductNotFound, "商品不存在")
		return
	}
	response.Success(ctx, p)
}

func (c *ProductController) Categories(ctx *gin.Context) {
	list, _ := service.NewCategoryService().All()
	response.Success(ctx, list)
}

// ========== Order (Frontend) ==========

type OrderController struct{}

func NewOrderController() *OrderController { return &OrderController{} }

func (c *OrderController) Create(ctx *gin.Context) {
	var req struct {
		ProductID uint   `json:"product_id" form:"product_id"`
		Quantity  int    `json:"quantity" form:"quantity"`
		PayType   string `json:"pay_type" form:"pay_type"`
		Email     string `json:"email" form:"email"`
		Remark    string `json:"remark" form:"remark"`
	}
	ctx.ShouldBind(&req)
	if req.Quantity <= 0 {
		req.Quantity = 1
	}
	uid, ok := ctx.Get("user_id")
	var userID uint
	if ok {
		userID = uid.(uint)
	}
	order, err := service.NewOrderService().Create(userID, req.ProductID, req.Quantity, req.Email, req.PayType, req.Remark)
	if err != nil {
		response.Error(ctx, response.CodeParamError, err.Error())
		return
	}
	response.Success(ctx, order)
}

func (c *OrderController) List(ctx *gin.Context) {
	uid, _ := ctx.Get("user_id")
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(ctx.DefaultQuery("size", "20"))
	total, list, _ := service.NewOrderService().ListByUser(uid.(uint), page, size)
	response.Success(ctx, gin.H{"total": total, "list": list, "page": page, "size": size})
}

func (c *OrderController) Detail(ctx *gin.Context) {
	orderNo := ctx.Param("order_no")
	order, err := service.NewOrderService().GetByOrderNo(orderNo)
	if err != nil {
		response.Error(ctx, response.CodeOrderNotFound, "订单不存在")
		return
	}
	cards, _ := service.NewOrderService().GetOrderCards(order.ID)
	response.Success(ctx, gin.H{"order": order, "cards": cards})
}

// ========== Payment ==========

type PaymentController struct{}

func NewPaymentController() *PaymentController { return &PaymentController{} }

func (c *PaymentController) Pay(ctx *gin.Context) {
	var req struct {
		OrderNo string `json:"order_no" form:"order_no"`
	}
	ctx.ShouldBind(&req)
	order, err := service.NewOrderService().GetByOrderNo(req.OrderNo)
	if err != nil {
		response.Error(ctx, response.CodeOrderNotFound, "订单不存在")
		return
	}
	if order.Status != 0 {
		response.Success(ctx, gin.H{"paid": true, "status": order.Status})
		return
	}
	epaySvc := service.NewEpayService()
	scheme := "http://"
	if ctx.Request.TLS != nil {
		scheme = "https://"
	}
	base := scheme + ctx.Request.Host
	returnURL := base + "/order/" + order.OrderNo
	notifyURL := base + "/api/v1/pay/epay/notify"
	payURL := epaySvc.BuildRedirectURL(order.OrderNo, order.Amount, order.ProductSnapshot, returnURL, notifyURL)
	response.Success(ctx, gin.H{"pay_url": payURL, "order_no": order.OrderNo})
}

func (c *PaymentController) EpayNotify(ctx *gin.Context) {
	form := ctx.Request.URL.Query()
	_ = ctx.Request.ParseForm()
	for k, v := range ctx.Request.PostForm {
		form[k] = v
	}
	epaySvc := service.NewEpayService()
	orderNo, _, ok := epaySvc.VerifyNotify(form)
	if !ok {
		ctx.String(http.StatusBadRequest, "fail")
		return
	}
	_, err := service.NewOrderService().MarkPaid(orderNo, "epay")
	if err != nil {
		ctx.String(http.StatusInternalServerError, "fail")
		return
	}
	ctx.String(http.StatusOK, "success")
}

func (c *PaymentController) EpayReturn(ctx *gin.Context) {
	orderNo := ctx.Query("out_trade_no")
	ctx.Redirect(http.StatusFound, "/order/"+orderNo)
}

// ========== Common helper ==========

func getUID(ctx *gin.Context) uint {
	uid, ok := ctx.Get("user_id")
	if !ok {
		return 0
	}
	return uid.(uint)
}

// helper to get current time
func now() time.Time { return time.Now() }
