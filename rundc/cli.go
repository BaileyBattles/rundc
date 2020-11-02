package rundc

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"rundc/pkg/log"
	"rundc/rundc/cmd"
)

type Cli struct {
}

func (c *Cli) Main(args []string) {
	switch os.Args[1] {
	case "pull":
		cmd.Pull(os.Args[2])
	default:
		log.ErrorAndExit("Bad Arguments")
	}
	log.Info(args)
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
