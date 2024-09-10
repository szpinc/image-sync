package server

import (
	"io"

	"github.com/docker/distribution/manifest/schema2"
	"github.com/gin-gonic/gin"
	"github.com/heroku/docker-registry-client/registry"
	"github.com/opencontainers/go-digest"
	"github.com/szpinc/image-sync/internal/config"
	"github.com/szpinc/image-sync/internal/util"
)

type ImageServer struct {
	Config   *config.ServerConfig
	Registry *registry.Registry
	engine   *gin.Engine
	Log      *util.Logger
}

func NewImageServer(serverConfig *config.ServerConfig) *ImageServer {

	server := ImageServer{Config: serverConfig}

	r, err := registry.New(serverConfig.RegistryConfig.Url, serverConfig.RegistryConfig.Username, serverConfig.RegistryConfig.Password)

	if err != nil {
		panic(err)
	}
	server.Registry = r
	server.engine = gin.New()
	server.Log = util.BuildLogger(serverConfig.LogConfig.Level)
	return &server
}

func (s *ImageServer) Start() {

	s.Log.Info("Server starting at %s", s.Config.Addr)

	// 初始化路由
	s.InitRouters()

	imageServer = s

	err := s.engine.Run(s.Config.Addr)
	if err != nil {
		panic(err)
	}
}

func (s *ImageServer) Stop() {}

func (s *ImageServer) UploadBlob(digest digest.Digest, repository string, data io.ReadCloser) error {
	return s.Registry.UploadBlob(repository, digest, data)
}

func (s *ImageServer) GetManifest(repository string, tag string) (*schema2.DeserializedManifest, error) {
	return s.Registry.ManifestV2(repository, tag)
}

func (s *ImageServer) HasBlob(repository string, digest digest.Digest) (bool, error) {
	return s.Registry.HasBlob(repository, digest)
}

func (s *ImageServer) PushManifest(repository string, tag string, manifest *schema2.DeserializedManifest) error {
	return s.Registry.PutManifest(repository, tag, manifest)
}
