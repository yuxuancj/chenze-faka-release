package controller

import (
	"chenze-faka/internal/model"
	"chenze-faka/internal/pkg/config"
	"chenze-faka/internal/pkg/db"
	"chenze-faka/internal/pkg/jwt"
	"chenze-faka/internal/pkg/response"
	"chenze-faka/internal/service"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

type AdminController struct{}

func NewAdminController() *AdminController { return &AdminController{} }

func (a *AdminController) Login(ctx *gin.Context) {
	var req struct {
		Email    string `json:"email" form:"email"`
		Password string `json:"password" form:"password"`
	}
	ctx.ShouldBind(&req)
	u, err := service.NewUserService().Login(req.Email, req.Password)
	if err != nil || !u.IsAdmin {
		response.Error(ctx, response.CodePasswordError, "管理员账号或密码错误")
		return
	}
	token, _ := jwt.Generate(u.ID, u.Email, u.IsAdmin, config.AppConfig.JWT.Secret, config.AppConfig.JWT.Expire)
	response.Success(ctx, gin.H{"user": u, "token": token})
}

func (a *AdminController) Dashboard(ctx *gin.Context) {
	var userCount, productCount, orderCount int64
	var totalAmount float64
	if db.DB != nil {
		db.DB.Model(&model.User{}).Count(&userCount)
		db.DB.Model(&model.Product{}).Count(&productCount)
		db.DB.Model(&model.Order{}).Count(&orderCount)
		db.DB.Model(&model.Order{}).Where("status IN ?", []int{1, 2}).Select("IFNULL(SUM(amount),0)").Scan(&totalAmount)
	}
	response.Success(ctx, gin.H{
		"user_count":    userCount,
		"product_count": productCount,
		"order_count":   orderCount,
		"total_amount":   totalAmount,
	})
}

func (a *AdminController) ProductList(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(ctx.DefaultQuery("size", "20"))
	keyword := ctx.DefaultQuery("keyword", "")
	total, list, err := service.NewProductService().AdminList(page, size, keyword)
	if err != nil {
		response.Error(ctx, response.CodeServerError, err.Error())
		return
	}
	response.Success(ctx, gin.H{"total": total, "list": list, "page": page, "size": size})
}

func (a *AdminController) ProductCreate(ctx *gin.Context) {
	var req struct {
		CategoryID  uint    `json:"category_id" form:"category_id"`
		Name        string  `json:"name" form:"name"`
		Description string  `json:"description" form:"description"`
		Price      float64 `json:"price" form:"price"`
		Stock      int     `json:"stock" form:"stock"`
		Image      string  `json:"image" form:"image"`
		Type       string  `json:"type" form:"type"`
		Status     int     `json:"status" form:"status"`
	}
	ctx.ShouldBind(&req)
	if req.Name == "" || req.Price <= 0 {
		response.Error(ctx, response.CodeParamError, "商品名称或价格无效")
		return
	}
	p := &model.Product{
		CategoryID:  req.CategoryID,
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		Stock:       req.Stock,
		Image:       req.Image,
		Type:        req.Type,
		Status:      req.Status,
	}
	if p.Status == 0 {
		p.Status = 1
	}
	if p.Type == "" {
		p.Type = "card"
	}
	if err := service.NewProductService().Create(p); err != nil {
		response.Error(ctx, response.CodeServerError, err.Error())
		return
	}
	response.Success(ctx, p)
}

func (a *AdminController) ProductUpdate(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.Param("id"))
	var req map[string]interface{}
	ctx.ShouldBindJSON(&req)
	if err := service.NewProductService().Update(uint(id), req); err != nil {
		response.Error(ctx, response.CodeServerError, err.Error())
		return
	}
	response.OK(ctx)
}

func (a *AdminController) ProductDelete(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.Param("id"))
	if err := service.NewProductService().Delete(uint(id)); err != nil {
		response.Error(ctx, response.CodeServerError, err.Error())
		return
	}
	response.OK(ctx)
}

func (a *AdminController) CardList(ctx *gin.Context) {
	productID, _ := strconv.Atoi(ctx.Query("product_id"))
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(ctx.DefaultQuery("size", "50"))
	total, list, err := service.NewCardService().ListByProduct(uint(productID), page, size)
	if err != nil {
		response.Error(ctx, response.CodeServerError, err.Error())
		return
	}
	response.Success(ctx, gin.H{"total": total, "list": list, "page": page, "size": size})
}

func (a *AdminController) CardImport(ctx *gin.Context) {
	productID, _ := strconv.Atoi(ctx.PostForm("product_id"))
	cardsStr := ctx.PostForm("cards")
	if productID == 0 || cardsStr == "" {
		response.Error(ctx, response.CodeParamError, "参数错误")
		return
	}
	cards := splitLines(cardsStr)
	count, err := service.NewCardService().Import(uint(productID), cards)
	if err != nil {
		response.Error(ctx, response.CodeServerError, err.Error())
		return
	}
	response.Success(ctx, gin.H{"count": count})
}

func (a *AdminController) CategoryList(ctx *gin.Context) {
	list, _ := service.NewCategoryService().All()
	response.Success(ctx, list)
}

func (a *AdminController) CategoryCreate(ctx *gin.Context) {
	var req struct {
		Name     string `json:"name" form:"name"`
		ParentID uint   `json:"parent_id" form:"parent_id"`
		Sort     int    `json:"sort" form:"sort"`
	}
	ctx.ShouldBind(&req)
	c, err := service.NewCategoryService().Create(req.Name, req.ParentID, req.Sort)
	if err != nil {
		response.Error(ctx, response.CodeServerError, err.Error())
		return
	}
	response.Success(ctx, c)
}

func (a *AdminController) CategoryUpdate(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.Param("id"))
	var req struct {
		Name string `json:"name" form:"name"`
		Sort int    `json:"sort" form:"sort"`
	}
	ctx.ShouldBind(&req)
	if err := service.NewCategoryService().Update(uint(id), req.Name, req.Sort); err != nil {
		response.Error(ctx, response.CodeServerError, err.Error())
		return
	}
	response.OK(ctx)
}

func (a *AdminController) CategoryDelete(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.Param("id"))
	if err := service.NewCategoryService().Delete(uint(id)); err != nil {
		response.Error(ctx, response.CodeServerError, err.Error())
		return
	}
	response.OK(ctx)
}

func (a *AdminController) OrderList(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(ctx.DefaultQuery("size", "20"))
	keyword := ctx.DefaultQuery("keyword", "")
	total, list, err := service.NewOrderService().AdminList(page, size, keyword)
	if err != nil {
		response.Error(ctx, response.CodeServerError, err.Error())
		return
	}
	response.Success(ctx, gin.H{"total": total, "list": list, "page": page, "size": size})
}

func (a *AdminController) OrderDetail(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.Param("id"))
	order, err := service.NewOrderService().GetByID(uint(id))
	if err != nil {
		response.Error(ctx, response.CodeOrderNotFound, "订单不存在")
		return
	}
	cards, _ := service.NewOrderService().GetOrderCards(order.ID)
	response.Success(ctx, gin.H{"order": order, "cards": cards})
}

func (a *AdminController) UserList(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(ctx.DefaultQuery("size", "20"))
	total, list, err := service.NewUserService().List(page, size)
	if err != nil {
		response.Error(ctx, response.CodeServerError, err.Error())
		return
	}
	response.Success(ctx, gin.H{"total": total, "list": list, "page": page, "size": size})
}

func (a *AdminController) SettingsGet(ctx *gin.Context) {
	ss := service.NewSettingService()
	response.Success(ctx, gin.H{
		"site_name": ss.Get("site_name", "晨泽发卡"),
		"site_desc": ss.Get("site_desc", ""),
	})
}

func (a *AdminController) SettingsSet(ctx *gin.Context) {
	var req map[string]string
	ctx.ShouldBindJSON(&req)
	ss := service.NewSettingService()
	for k, v := range req {
		_ = ss.Set(k, v)
	}
	response.OK(ctx)
}

func splitLines(s string) []string {
	lines := strings.Split(s, "\n")
	out := make([]string, 0, len(lines))
	for _, l := range lines {
		l = strings.TrimSpace(l)
		if l != "" {
			out = append(out, l)
		}
	}
	return out
}
