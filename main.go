package main

import (
	"fmt"
	"os/exec"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"

	"golang.org/x/net/context"
)

func createDockerClient() (*client.Client, error) {
	cli, err := client.NewClientWithOpts(client.WithVersion("1.37"))
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

func createContainer(ctx context.Context, docker *client.Client, imageName string, containerPort nat.Port, serviceName string) response {

	env := []string{
		"HOSTNAME=localhost",
		"SERVICE=" + serviceName,
	}

	config := &container.Config{
		Image:    imageName,
		Hostname: "localhost",
		ExposedPorts: nat.PortSet{
			containerPort: struct{}{},
			"8080/tcp":    struct{}{},
		},
		Env: env,
	}

	hostConfig := &container.HostConfig{
		PortBindings: nat.PortMap{
			containerPort: []nat.PortBinding{
				{
					HostIP:   "0.0.0.0",
					HostPort: containerPort.Port(),
				},
			},
			"8080/tcp": []nat.PortBinding{
				{
					HostIP:   "0.0.0.0",
					HostPort: "8080",
				},
			},
		},
	}

	respCreate, errCreate := docker.ContainerCreate(ctx, config, hostConfig, nil, "")
	if errCreate != nil {
		panic(errCreate)
	}

	responseCreate := response{resp: respCreate, err: errCreate}

	return responseCreate
}

func startContainer(ctx context.Context, docker *client.Client, resp container.ContainerCreateCreatedBody, err error) {
	if err := docker.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		panic(err)
	}
}

type response struct {
	resp container.ContainerCreateCreatedBody
	err  error
}

type serviceNamePort struct {
	name string
	port nat.Port
}

const DockerImage = "localstack/localstack"

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

func main() {
	ctx := context.Background()

	docker, err := createDockerClient()
	pullImage(ctx, DockerImage, docker, err)

	serviceDynamo := serviceNamePort{name: "dynamodb", port: "4569/tcp"}
	serviceS3 := serviceNamePort{name: "S3", port: "4572/tcp"}
	serviceUI := serviceNamePort{name: "", port: "8080/tcp"}

	services := map[int]serviceNamePort{
		0: serviceDynamo,
		1: serviceS3,
		2: serviceUI,
	}

	response := createOneContainer(ctx, docker, DockerImage, services)

	startContainer(ctx, docker, response.resp, response.err)
	for i := range services {
		fmt.Println(services[i].name, " service running!!")
	}

	// Wait until services is accesible
	time.Sleep(20 * time.Second)

	// Bucket creation
	bucketName := "test"
	createBucket(bucketName, serviceS3)
}

func createBucket(bucketName string, serviceS3 serviceNamePort) {
	cmd := exec.Command("/bin/bash", "-c", "aws s3api create-bucket --endpoint-url=http://localhost:"+serviceS3.port.Port()+" --bucket "+bucketName)

	_, errCommand := cmd.Output()

	if errCommand != nil {
		fmt.Println(errCommand.Error())
		return
	}
	fmt.Printf("Bucket %q successfully created\n", bucketName)
}
