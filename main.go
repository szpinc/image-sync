package main

import (
	"hua-cloud.com/tools/image-sync/internal/config"
	"hua-cloud.com/tools/image-sync/internal/server"
)

func main() {

	imageServer := server.NewImageServer(&config.ServerConfig{
		Addr: ":23333",
		RegistryConfig: config.RegistryConfig{
			Url: "http://127.0.0.1:50000",
		},
		LogConfig: config.LogConfig{
			Level: "debug",
		},
	})

	imageServer.Start()

	//src, err := registry.New("http://harbor.hy-zw.com", "admin", "Hyzs(%23@7&8")
	//
	//if err != nil {
	//	panic(err)
	//}
	//
	//target, err := registry.New("https://smartum.sz.gov.cn/park/registry", "", "")
	//
	//if err != nil {
	//	panic(err)
	//}
	//
	////repositories, err := target.Repositories()
	//
	//manifest, err := src.ManifestV2("park/file-upload-service", "2.5.0")
	//
	//if err != nil {
	//	panic(err)
	//}
	//
	//for _, layer := range manifest.Layers {
	//	println(layer.Digest.String())
	//	blob, err := src.DownloadBlob("park/file-upload-service", layer.Digest)
	//	if err != nil {
	//		panic(err)
	//	}
	//	err = UploadBlob("park/file-upload-service", target, layer.Digest, blob)
	//
	//	if err != nil {
	//		panic(err)
	//	}
	//}

	//target.PutManifest()

}

//func initiateUpload(repository string, registry *registry.Registry) (*url.URL, error) {
//	initiateURL := fmt.Sprintf("%s/v2/%s/blobs/uploads/", registry.URL, repository)
//	registry.Logf("registry.blob.initiate-upload url=%s repository=%s", initiateURL, repository)
//
//	resp, err := registry.Client.Post(initiateURL, "application/octet-stream", nil)
//	if resp != nil {
//		defer resp.Body.Close()
//	}
//	if err != nil {
//		return nil, err
//	}
//
//	location := resp.Header.Get("Location")
//	locationURL, err := url.Parse(location)
//	if err != nil {
//		return nil, err
//	}
//	registry.Logf("Location: %s", locationURL.String())
//	return locationURL, nil
//}
//
//func getUrl(uri string, pathTemplate string, args ...interface{}) string {
//	pathSuffix := fmt.Sprintf(pathTemplate, args...)
//	url := fmt.Sprintf("%s%s", uri, pathSuffix)
//	return url
//}
//
//func UploadBlob(repository string, registry *registry.Registry, digest digest.Digest, content io.Reader) error {
//
//	registryUri, _ := url.ParseRequestURI(registry.URL)
//
//	uploadURL, err := initiateUpload(repository, registry)
//	if err != nil {
//		return err
//	}
//
//	q := uploadURL.Query()
//	q.Set("digest", digest.String())
//	uploadURL.RawQuery = q.Encode()
//
//	registryUri = registryUri.JoinPath(uploadURL.Path)
//	registryUri.RawQuery = uploadURL.RawQuery
//
//	registry.Logf("registry.blob.upload url=%s repository=%s digest=%s", registryUri, repository, digest)
//
//	upload, err := http.NewRequest("PUT", registryUri.String(), content)
//	if err != nil {
//		return err
//	}
//	upload.Header.Set("Content-Type", "application/octet-stream")
//
//	_, err = registry.Client.Do(upload)
//	return err
