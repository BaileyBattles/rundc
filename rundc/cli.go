package rundc

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"rundc/pkg/log"

	digest "github.com/opencontainers/go-digest"
	v1 "github.com/opencontainers/image-spec/specs-go/v1"
)

type Cli struct {
}

func (c *Cli) Main(args []string) {
	switch os.Args[1] {
	case "pull":
		pull(os.Args[2])
	default:
		log.LogErrorAndExit("Bad Arguments")
	}
	fmt.Println(args)
}

//dockerResolver.Resolve in containerd/remotes/docker/resolver.go
func pull(imageName string) {
	manifest := getManifest(imageName)
	image := getImage(imageName, manifest.Config.Digest)
	fmt.Println(image.OS)
}

func getManifest(imageName string) v1.Manifest {
	req, err := getRequestWithAuthHeaders(fmt.Sprintf("https://registry-1.docker.io/v2/library/%s/manifests/latest", imageName))
	if err != nil {
		panic(err)
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.LogErrorAndExit(err.Error())
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	var manifest v1.Manifest
	err = json.Unmarshal(body, &manifest)
	if err != nil {
		log.LogErrorAndExit(err.Error())
	}
	return manifest
}

func getImage(imageName string, digest digest.Digest) v1.Image {
	req, err := getRequestWithAuthHeaders(fmt.Sprintf("https://registry-1.docker.io/v2/library/%s/blobs/%s", imageName, digest))
	if err != nil {
		panic(err)
	}
	client := &http.Client{}

	resp, _ := client.Do(req)
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	var image v1.Image
	err = json.Unmarshal(body, &image)
	if err != nil {
		log.LogErrorAndExit(err.Error())
	}
	return image
}

func getRequestWithAuthHeaders(endpoint string) (*http.Request, error) {
	resp, err := http.Get("https://auth.docker.io/token?service=registry.docker.io&scope=repository:library/alpine:pull")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var tokenResponse map[string]interface{}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(body, &tokenResponse)
	if err != nil {
		return nil, err
	}

	auth_header := fmt.Sprintf("Bearer %s", tokenResponse["token"])

	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", auth_header)
	req.Header.Add("Accept", "application/vnd.docker.distribution.manifest.v2+json")
	return req, nil
}

func printResponse(resp *http.Response) {
	var dataResponse map[string]interface{}
	x, _ := ioutil.ReadAll(resp.Body)
	json.Unmarshal(x, &dataResponse)
	fmt.Println(dataResponse)
	fmt.Print("\n\n\n")
	fmt.Println(resp.Header)
	fmt.Print("\n\n\n")
}
