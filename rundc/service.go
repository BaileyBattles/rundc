package rundc

import (
	"fmt"
	"syscall"

	"rundc/pkg/log"
)

type ContainerService struct {
}

func (self *ContainerService) StartContainer(c *Container) {

}

func startContainerProcess() {
	err := syscall.Sethostname([]byte("container"))
	if err != nil {
		log.LogErrorAndExit(fmt.Sprintf("failed to set hostname: %s", err.Error()))
	}
}
