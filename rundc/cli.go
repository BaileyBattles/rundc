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

func pull(imageName string) {
	resp, err := http.Get("https://auth.docker.io/token?service=registry.docker.io&scope=repository:library/alpine:pull")
	if err != nil {
		fmt.Println(err)
	}
	var response map[string]string
	body, _ := ioutil.ReadAll(resp.Body)
	json.Unmarshal(body, &response)
	fmt.Println(response["token"])

	req, _ := http.NewRequest("GET", "https://registry-1.docker.io/v2/library/alpine/manifests/latest", nil)
	auth_header := fmt.Sprintf("Bearer %s", response["token"])
	req.Header.Add("Authorization", auth_header)
	req.Header.Add("Accept", "application/vnd.docker.distribution.manifest.v2+json")
	client := &http.Client{}
	resp, _ = client.Do(req)
	fmt.Println(resp)
}
