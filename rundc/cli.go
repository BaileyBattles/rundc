package rundc

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"rundc/pkg/log"
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
	req, err := getRequestWithAuthHeaders("https://registry-1.docker.io/v2/library/alpine/manifests/latest")
	if err != nil {
		panic(err)
	}
	client := &http.Client{}
	resp, _ := client.Do(req)
	printResponse(resp)

	resp.Body.Close()
	//fmt.Sprintf("https://registry-1.docker.io/v2/library/alpine/blobs/%s", resp[config][digest])
	req, err = getRequestWithAuthHeaders("https://registry-1.docker.io/v2/library/alpine/blobs/sha256:d6e46aa2470df1d32034c6707c8041158b652f38d2a9ae3d7ad7e7532d22ebe0")
	if err != nil {
		panic(err)
	}
	resp, _ = client.Do(req)

	printResponse(resp)
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
