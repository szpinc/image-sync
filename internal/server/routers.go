package server

import (
	"github.com/gin-gonic/gin"
	"hua-cloud.com/tools/image-sync/internal/middleware"
)

func InitRouters(engine *gin.Engine) {
	v1 := engine.Group("/api/v1")

	{
		blob := v1.Group("blob")
		// 上传blob
		blob.PUT("uploads", middleware.Wrapper(uploads))
		// 校验
		blob.GET(":repository/:digest", middleware.Wrapper(blobExists))
	}
}
