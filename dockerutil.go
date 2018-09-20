package main

import (
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	"golang.org/x/net/context"
)

func createDockerClient() (*client.Client, error) {
	cli, err := client.NewEnvClient()
	if err != nil {
		panic(err)
	}

	return cli, err
}

func pullImage(ctx context.Context, imageName string, docker *client.Client, err error) {
	_, err = docker.ImagePull(ctx, imageName, types.ImagePullOptions{})
	if err != nil {
		panic(err)
	}
}

func startContainer(ctx context.Context, docker *client.Client, resp container.ContainerCreateCreatedBody, err error) {
	if err := docker.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		panic(err)
	}
}

func createOneContainer(ctx context.Context, docker *client.Client, imageName string, services map[int]serviceNamePort) response {
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
		panic(errCreate)
	}

	responseCreate := response{resp: respCreate, err: errCreate}

	return responseCreate
}
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

func exposedPorts(services map[int]serviceNamePort) nat.PortSet {
	var exposedPorts = make(map[nat.Port]struct{})
	for _, service := range services {
		exposedPorts[service.port] = struct{}{}
	}
	return exposedPorts
}
