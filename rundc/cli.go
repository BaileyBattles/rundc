package rundc

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"rundc/pkg/log"
	"strings"

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
		log.ErrorAndExit("Bad Arguments")
	}
	log.Info(args)
}

//dockerResolver.Resolve in containerd/remotes/docker/resolver.go
func pull(imageName string) {
	manifest := getManifest(imageName)
	image := getImage(imageName, manifest.Config.Digest)
	createFs(imageName)
	for _, layer := range manifest.Layers {
		pullLayer(imageName, layer.Digest)
	}
	log.Info(image.OS)
}

func createFs(imageName string) {
	err := os.Mkdir(imageName, os.ModeDir)
	if err != nil {
		log.ErrorAndExit(err.Error())
	}
}

func getManifest(imageName string) v1.Manifest {
	req, err := getRequestWithAuthHeaders(fmt.Sprintf("https://registry-1.docker.io/v2/library/%s/manifests/latest", imageName))
	if err != nil {
		panic(err)
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.ErrorAndExit(err.Error())
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	var manifest v1.Manifest
	err = json.Unmarshal(body, &manifest)
	if err != nil {
		log.ErrorAndExit(err.Error())
	}
	return manifest
}

func getImage(imageName string, digest digest.Digest) v1.Image {
	req, err := getRequestWithAuthHeaders(fmt.Sprintf("https://registry-1.docker.io/v2/library/%s/blobs/%s", imageName, digest))
	if err != nil {
		panic(err)
	}
	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		log.ErrorAndExit(err.Error())
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	var image v1.Image
	err = json.Unmarshal(body, &image)
	if err != nil {
		log.ErrorAndExit(err.Error())
	}
	return image
}

func pullLayer(imageName string, digest digest.Digest) {
	fmt.Println(digest)
	req, err := getRequestWithAuthHeaders(fmt.Sprintf("https://registry-1.docker.io/v2/library/%s/blobs/%s", imageName, digest))
	if err != nil {
		log.ErrorAndExit(err.Error())
	}
	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		log.ErrorAndExit(err.Error())
	}
	reader := bufio.NewReader(resp.Body)
	f, err := os.Create(fmt.Sprintf("%s/layer.tar", imageName))
	if err != nil {
		log.ErrorAndExit(err.Error())
	}
	defer f.Close()
	for {
		buffer := make([]byte, 8192)
		_, err := io.ReadFull(reader, buffer)
		if err != nil {
			if strings.Contains(err.Error(), "EOF") {
				f.Write(buffer)
				break
			}
			log.ErrorAndExit(err.Error())
		}
		f.Write(buffer)
	}
	fmt.Println("here")

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
	log.Info(dataResponse)
	fmt.Print("\n\n\n")
	log.Info(resp.Header)
	fmt.Print("\n\n\n")
}
