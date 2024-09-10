package server

import (
	"fmt"
	"os"
	"path"
	"time"

	"github.com/docker/distribution/manifest/schema2"
	"github.com/gin-gonic/gin"
	"github.com/opencontainers/go-digest"
	"github.com/szpinc/image-sync/pkg/types"
	"golang.org/x/crypto/ssh"
)

var imageServer *ImageServer

func uploads(c *gin.Context) (data any, err error) {
	digestHash := c.Query("digest")
	repository := c.Query("repository")

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
	digestHash := c.Query("digest")
	repository := c.Query("repository")

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

func pushManifest(c *gin.Context) (data any, err error) {
	repository := c.Query("repository")
	tag := c.Query("tag")

	manifest := schema2.DeserializedManifest{}
	err = c.ShouldBind(&manifest)
	if err != nil {
		return nil, err
	}
	return nil, imageServer.PushManifest(repository, tag, &manifest)
}

func exec(c *gin.Context) (data any, err error) {

	cmdRequest := types.CmdRequest{}
	if err = c.ShouldBind(&cmdRequest); err != nil {
		return nil, err
	}

	sshClient, err := SSHConnect("root", cmdRequest.Host, cmdRequest.Port)

	if err != nil {
		return nil, err
	}

	defer func(sshClient *ssh.Client) {
		_ = sshClient.Close()
	}(sshClient)

	var result []string
	var output string

	// 更新docker compose文件镜像
	output, err = execCmd(sshClient, fmt.Sprintf("sed -i s#%s.*#%s#g %s", cmdRequest.Repository, cmdRequest.Repository+":"+cmdRequest.Tag, cmdRequest.DockerComposeFile))

	result = append(result, output)

	if err != nil {
		return output, err
	}

	if cmdRequest.Deploy {
		output, err = execCmd(sshClient, fmt.Sprintf("docker-compose -f %s up -d %s", cmdRequest.DockerComposeFile, cmdRequest.App))
		if err != nil {
			return output, err
		}
		result = append(result, output)
		imageServer.Log.Info("应用更新成功!")
	}

	return result, nil
}

// 执行ssh命令
func execCmd(client *ssh.Client, cmd string) (string, error) {
	session, err := client.NewSession()

	if err != nil {
		return "", err
	}

	defer func(session *ssh.Session) {
		_ = session.Close()
	}(session)

	output, err := session.CombinedOutput(cmd)

	if err != nil {
		return "", err
	}

	return string(output), nil
}

func SSHConnect(user, host string, port int) (*ssh.Client, error) {
	var (
		addr         string
		clientConfig *ssh.ClientConfig
		client       *ssh.Client
		err          error
	)

	homePath, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}
	key, err := os.ReadFile(path.Join(homePath, ".ssh", "id_rsa"))
	if err != nil {
		return nil, err
	}
	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		return nil, err
	}

	clientConfig = &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
		Timeout:         30 * time.Second,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	// connect to ssh
	addr = fmt.Sprintf("%s:%d", host, port)

	if client, err = ssh.Dial("tcp", addr, clientConfig); err != nil {
		return nil, err
	}

	return client, nil
}
