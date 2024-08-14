package cmd

import "testing"

func TestParseDockerImage(t *testing.T) {
	image, s, s2, s3 := ParseDockerImage("127.0.0.1:50000/mzj_hyzs/screen-service-homeless:v1.0.0")

	println(image, s, s2, s3)
}
