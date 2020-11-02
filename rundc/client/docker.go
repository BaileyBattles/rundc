package client

import (
	"archive/tar"
	"bufio"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"github.com/opencontainers/go-digest"
	v1 "github.com/opencontainers/image-spec/specs-go/v1"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"rundc/pkg/log"
)

const (
	registry = "https://registry-1.docker.io"
	tage     = "latest"
)
type DockerHubClient struct {

}

func (self *DockerHubClient) PullImage(imageName string) (v1.Image, error) {
	fmt.Printf("Pulling %s from Docker\n", imageName)
	log.Info("pulling manifest")
	manifest := getManifest(imageName)
	log.Info("pulling image")
	image := getImage(imageName, manifest.Config.Digest)
	log.Info("creating fs")
	createFs(imageName)
	log.Info("pulling layers")
	for _, layer := range manifest.Layers {
		pullLayer(imageName, layer.Digest)
	}
	log.Info(image.OS)
	return image, nil
}

func createFs(imageName string) {
	err := os.Mkdir(imageName, os.ModeDir)
	if err != nil {
		log.ErrorAndExit(err.Error())
	}
	err = os.Mkdir(fmt.Sprintf("%s/rootfs", imageName), os.ModeDir)
	if err != nil {
		log.ErrorAndExit(err.Error())
	}
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

func pullLayer(imageName string, digest digest.Digest) {
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
	filename := fmt.Sprintf("%s/layer.tar", imageName)
	f, err := os.Create(filename)
	if err != nil {
		log.ErrorAndExit(err.Error())
	}
	defer f.Close()
	for {
		buffer := make([]byte, 8192)
		_, err := io.ReadFull(reader, buffer)
		if err != nil {
			if err == io.EOF || err == io.ErrUnexpectedEOF{
				f.Write(buffer)
				break
			}
			log.ErrorAndExit(err.Error())
		}
		f.Write(buffer)
	}
	decompressLayer("alpine/rootfs", filename)
	err = os.Remove(filename)

}

func decompressLayer(dst string, filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	gzr, err := gzip.NewReader(file)
	defer gzr.Close()
	tr := tar.NewReader(gzr)
	for {
		header, err := tr.Next()
		switch {
		case err == io.EOF:
			return nil
		case err != nil:
			return err
		case header == nil:
			continue
		}
		target := filepath.Join(dst, header.Name)
		switch header.Typeflag {
		case tar.TypeDir:
			if _, err := os.Stat(target); err != nil {
				if err := os.MkdirAll(target, 0755); err != nil {
					return err
				}
			}
		case tar.TypeReg:
			f, err := os.OpenFile(target, os.O_CREATE|os.O_RDWR, os.FileMode(header.Mode))
			if err != nil {
				return err
			}
			if _, err := io.Copy(f, tr); err != nil {
				return err
			}
			f.Close()
		}
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