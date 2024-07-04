package server

import (
	"github.com/gin-gonic/gin"
	"github.com/opencontainers/go-digest"
)

var imageServer *ImageServer

func uploads(c *gin.Context) (data any, err error) {
	digestHash := c.Param("digest")
	repository := c.Param("repository")

	dst, err := digest.Parse(digestHash)

	if err != nil {
		return nil, err
	}

	if err := imageServer.UploadBlob(dst, repository, c.Request.Body); err != nil {
		return nil, err
	}

	return nil, nil
}

func blobExists(c *gin.Context) (data any, err error) {
	digestHash := c.Param("digest")
	repository := c.Param("repository")

	dst, err := digest.Parse(digestHash)

	if err != nil {
		return nil, err
	}

	exists, err := imageServer.HasBlob(repository, dst)

	if err != nil {
		return nil, err
	}

	return exists, nil
}
