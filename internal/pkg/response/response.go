package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

const (
	CodeSuccess          = 0
	CodeParamError       = 1001
	CodeUserNotFound     = 1002
	CodePasswordError    = 1003
	CodeTokenInvalid     = 1004
	CodePermissionDenied = 1005
	CodeUserExists       = 1006
	CodeCaptchaError     = 1007
	CodeOrderNotFound    = 2001
	CodeInsufficientStock = 2002
	CodeProductNotFound  = 2003
	CodePaymentFailed    = 3001
	CodeServerError      = 5000
)

type Response struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data,omitempty"`
}

func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code: CodeSuccess,
		Msg:  "success",
		Data: data,
	})
}

func Error(c *gin.Context, code int, msg string) {
	c.JSON(http.StatusOK, Response{
		Code: code,
		Msg:  msg,
	})
}

func OK(c *gin.Context) {
	Success(c, nil)
}
