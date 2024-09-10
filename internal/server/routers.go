package server

import (
	"github.com/gin-gonic/gin"
	"github.com/szpinc/image-sync/internal/middleware"
)

func (s *ImageServer) InitRouters() {

	v1 := s.engine.Group("/api/v1", gin.BasicAuth(s.Config.Accounts))

	{
		blob := v1.Group("blob")
		{
			// 上传blob
			blob.PUT("uploads", middleware.Wrapper(uploads))
			// 校验
			blob.GET("exists", middleware.Wrapper(blobExists))
		}

		manifest := v1.Group("manifest")
		{
			// 推送manifest
			manifest.POST("push", middleware.Wrapper(pushManifest))
		}

		deploy := v1.Group("deploy")

		{
			deploy.POST("", middleware.Wrapper(exec))
		}
	}
}
