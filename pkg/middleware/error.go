package middleware

import (
	"net/http"

	"github.com/szpinc/image-sync/pkg/types"

	"github.com/gin-gonic/gin"
)

func errResp(message string) types.Resp {
	return types.Resp{
		Code:    http.StatusInternalServerError,
		Message: message,
		Data:    nil,
	}
}

func okResp(data interface{}) types.Resp {
	return types.Resp{
		Code:    http.StatusOK,
		Message: "ok",
		Data:    data,
	}
}

// ExceptionHandlerFunc 异常处理函数
type ExceptionHandlerFunc func(c *gin.Context) (data any, err error)

// Wrapper 中间件
func Wrapper(handlerFunc ExceptionHandlerFunc) func(c *gin.Context) {
	return func(c *gin.Context) {
		data, err := handlerFunc(c)
		if err != nil {
			c.JSON(http.StatusOK, errResp(err.Error()))
			return
		}
		c.JSON(http.StatusOK, okResp(data))
	}
}
