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

func main() {
	ctx := context.Background()

	docker, err := createDockerClient()
	pullImage(ctx, DockerImage, docker, err)

	// serviceDynamo := serviceNamePort{name: "dynamodb", port: "4569/tcp"}
	serviceS3 := serviceNamePort{name: "S3", port: "4572/tcp"}
	// serviceMonitor := serviceNamePort{name: "", port: "8080/tcp"}

	services := map[int]serviceNamePort{
		// 0: serviceDynamo,
		1: serviceS3,
		// 2: serviceMonitor,
	}

	resp := make(map[int]response)

	for k, s := range services {
		resp[k] = createContainer(ctx, docker, DockerImage, s.port, s.name)
		fmt.Println(services[k].name, " container created!!")
	}

	for k := range resp {
		startContainer(ctx, docker, resp[k].resp, resp[k].err)
		fmt.Println(services[k].name, " container running!!")
	}

	// Wait until services is accesible
	time.Sleep(20 * time.Second)

	// Bucket creation
	bucketName := "test"
	createBucket(bucketName)
}

func createBucket(bucketName string) {
	cmd := exec.Command("/bin/bash", "-c", "aws s3api create-bucket --endpoint-url=http://localhost:4572 --bucket "+bucketName)

	_, errCommand := cmd.Output()

	if errCommand != nil {
		fmt.Println(errCommand.Error())
		return
	}
	fmt.Printf("Bucket %q successfully created\n", bucketName)
}
