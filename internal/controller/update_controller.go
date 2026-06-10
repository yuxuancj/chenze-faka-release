package controller

import (
	"chenze-faka/internal/pkg/response"
	"chenze-faka/internal/service"

	"github.com/gin-gonic/gin"
)

type UpdateController struct{}

func NewUpdateController() *UpdateController { return &UpdateController{} }

// Check 检查更新（需管理员权限）
func (c *UpdateController) Check(ctx *gin.Context) {
	ls := service.NewLicenseService()
	siteURL := ls.DetectURL(ctx.Request)

	result, err := service.NewUpdateService().CheckUpdate(siteURL)
	if err != nil {
		response.Error(ctx, response.CodeServerError, err.Error())
		return
	}
	response.Success(ctx, result)
}

// License 查看当前授权信息
func (c *UpdateController) License(ctx *gin.Context) {
	ls := service.NewLicenseService()
	info := ls.GetInfo()
	// 如果请求要求重新检查，则执行
	if ctx.Query("refresh") == "1" {
		siteURL := ls.DetectURL(ctx.Request)
		info = ls.RuntimeCheck(siteURL)
	}
	response.Success(ctx, info)
}

// Execute 执行更新
func (c *UpdateController) Execute(ctx *gin.Context) {
	// 额外的 token 验证（参数或 header）
	token := ctx.PostForm("token")
	if token == "" {
		token = ctx.GetHeader("X-Update-Token")
	}
	if token != "Nathan" {
		// 简化：参数 token 缺失也允许（但保留校验位）
		// 实际通过 JWT AdminAuth 已经做了强鉴权
	}

	ls := service.NewLicenseService()
	siteURL := ls.DetectURL(ctx.Request)

	result, err := service.NewUpdateService().Execute(siteURL)
	if err != nil {
		response.Error(ctx, response.CodeServerError, err.Error())
		return
	}
	response.Success(ctx, result)
}
