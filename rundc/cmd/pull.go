package cmd

import "rundc/rundc/client"

//dockerResolver.Resolve in containerd/remotes/docker/resolver.go
func Pull(imageName string) {
	dockerClient := client.DockerHubClient{}
	dockerClient.PullImage(imageName)
}

