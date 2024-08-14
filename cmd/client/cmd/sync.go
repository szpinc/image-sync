package cmd

import (
	"errors"
	"strings"

	"github.com/spf13/cobra"
	"hua-cloud.com/tools/image-sync/internal/client"
	"hua-cloud.com/tools/image-sync/internal/types"
)

var (
	server            string
	username          string
	password          string
	registryUserName  string
	registryPassword  string
	targetHost        string
	targetPort        int
	targetComposeFile string
	targetAppName     string
	deploy            bool
)

// syncCmd represents the sync command
var syncCmd = &cobra.Command{
	Use:   "sync",
	Short: "sync from source registry to target registry",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return errors.New("sync image args required")
		}

		image := args[0]

		repo, name, tag, app := ParseDockerImage(image)

		if targetAppName == "" {
			targetAppName = app
		}

		c := client.NewClient(&types.ClientConfig{
			Version: "v1",
			Server: types.Server{
				Address:  server,
				Username: username,
				Password: password,
			},
			Registry: types.RegistryConfig{
				Url:      "http://" + repo,
				Username: registryUserName,
				Password: registryPassword,
			},
		})

		err := c.Copy(name, tag)

		if err != nil {
			return err
		}

		if !deploy {
			return nil
		}

		return c.Deploy(types.CmdRequest{
			Host:              targetHost,
			Port:              targetPort,
			App:               targetAppName,
			Repository:        name,
			Tag:               tag,
			DockerComposeFile: targetComposeFile,
			Deploy:            deploy,
		})
	},
}

func init() {
	rootCmd.AddCommand(syncCmd)

	syncCmd.Flags().StringVarP(&server, "server", "s", "", "target registry proxy server")
	syncCmd.Flags().StringVarP(&username, "username", "u", "", "target registry proxy username")
	syncCmd.Flags().StringVarP(&password, "password", "p", "", "target registry proxy password")
	syncCmd.Flags().StringVarP(&registryUserName, "src-registry-username", "", "", "source registry username")
	syncCmd.Flags().StringVarP(&registryPassword, "src-registry-password", "", "", "source registry password")
	syncCmd.Flags().StringVarP(&targetHost, "dest-deploy-host", "", "", "target deploy host")
	syncCmd.Flags().IntVarP(&targetPort, "dest-deploy-port", "", 22, "target deploy port")
	syncCmd.Flags().StringVarP(&targetComposeFile, "dest-deploy-compose-file", "", "", "target deploy docker compose file")
	syncCmd.Flags().StringVarP(&targetAppName, "application", "", "", "target app name")
	syncCmd.Flags().BoolVarP(&deploy, "deploy", "d", false, "deploy image to target registry")

}

// ParseDockerImage takes a Docker image string and returns the repository and image name.
func ParseDockerImage(image string) (string, string, string, string) {
	var repository, imageName, tag, app string

	// Check if the image contains a tag
	if strings.Contains(image, ":") {
		idx := strings.LastIndex(image, ":")
		tag = image[idx+len(":"):]
		image = image[:idx]
	}

	// Check if the image contains a repository
	if strings.Contains(image, "/") {
		parts := strings.Split(image, "/")
		repository = parts[0]
		imageName = strings.Join(parts[1:], "/")
		app = parts[len(parts)-1]
	} else {
		// If there's no repository, assume the image is from the default Docker Hub library
		repository = "library"
		imageName = image
	}

	return repository, imageName, tag, app
}
