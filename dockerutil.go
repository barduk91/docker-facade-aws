package main

import (
	"fmt"
	"io"
	"os"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	"golang.org/x/net/context"
)

// Create a docker client
func createDockerClient() (*client.Client, error) {
	cli, err := client.NewEnvClient()
	if err != nil {
		panic(err)
	}

	return cli, err
}

// Pull image from docker registry. It wait till image is downloaded
func pullImage(ctx context.Context, imageName string, docker *client.Client, err error) {
	var image io.ReadCloser
	image, err = docker.ImagePull(ctx, imageName, types.ImagePullOptions{})
	if err != nil {
		fmt.Println(err)
		return
	}
	io.Copy(os.Stdout, image)

}

// Start container previously created
func startContainer(ctx context.Context, docker *client.Client, resp container.ContainerCreateCreatedBody, err error) {
	if err := docker.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		fmt.Println(err)
		fmt.Println("Can't start docker container")
		return
	}
}

// Create container with service needed
func createContainer(ctx context.Context, docker *client.Client, imageName string, services map[int]serviceNamePort) response {
	servicesName := ""
	for _, service := range services {
		servicesName += service.name
	}

	env := []string{
		"HOSTNAME=localhost",
		"SERVICE=" + servicesName,
	}

	var servicePort nat.Port
	for _, service := range services {
		servicePort += service.port + ": struct{}{},"
	}

	config := &container.Config{
		Image:        imageName,
		Hostname:     "localhost",
		ExposedPorts: exposedPorts(services),
		Env:          env,
	}

	hostConfig := &container.HostConfig{
		PortBindings: portBindings(services),
	}

	respCreate, errCreate := docker.ContainerCreate(ctx, config, hostConfig, nil, "")
	if errCreate != nil {
		fmt.Println(errCreate.Error())
		fmt.Println("Can not create a container")
		os.Exit(-1)
	}

	responseCreate := response{resp: respCreate, err: errCreate}

	return responseCreate
}

// Bind port between localhost and docker container
func portBindings(services map[int]serviceNamePort) nat.PortMap {
	var portBindings = make(map[nat.Port][]nat.PortBinding)

	for _, service := range services {
		portBindings[service.port] = []nat.PortBinding{
			{
				HostIP:   "0.0.0.0",
				HostPort: service.port.Port(),
			},
		}
	}
	return portBindings
}

// Expose ports needed inside your docker container
func exposedPorts(services map[int]serviceNamePort) nat.PortSet {
	var exposedPorts = make(map[nat.Port]struct{})
	for _, service := range services {
		exposedPorts[service.port] = struct{}{}
	}
	return exposedPorts
}
