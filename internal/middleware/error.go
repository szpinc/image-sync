package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// 返回结构体
type resp struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data"`
}

func errResp(message string) resp {
	return resp{
		Code:    http.StatusInternalServerError,
		Message: message,
		Data:    nil,
	}
}

func okResp(data interface{}) resp {
	return resp{
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
