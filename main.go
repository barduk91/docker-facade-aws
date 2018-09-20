package main

import (
	"fmt"
	"log"
	"time"

	"github.com/awslabs/goformation"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/go-connections/nat"

	"golang.org/x/net/context"
)

// Image docker for localstack
var DockerImage = "localstack/localstack"

// Service supported
var serviceDynamo = serviceNamePort{name: "dynamodb", port: "4569/tcp"}
var serviceS3 = serviceNamePort{name: "S3", port: "4572/tcp"}
var serviceSqs = serviceNamePort{name: "sqs", port: "4576/tcp"}
var serviceUI = serviceNamePort{name: "", port: "8080/tcp"}

// Type for response when create docker container
type response struct {
	resp container.ContainerCreateCreatedBody
	err  error
}

// Type to handle service name and port for aws service
type serviceNamePort struct {
	name string
	port nat.Port
}

func main() {
	ctx := context.Background()

	docker, err := createDockerClient()
	pullImage(ctx, DockerImage, docker, err)

	template, err := goformation.Open("s3example.yaml")
	if err != nil {
		log.Fatalf("There was an error processing the template: %s", err)
	}

	services := getAllServices(template)

	response := createOneContainer(ctx, docker, DockerImage, services)

	startContainer(ctx, docker, response.resp, response.err)
	for i := range services {
		fmt.Println(services[i].name, " service running!!")
	}

	// Wait until services is up
	time.Sleep(20 * time.Second)

	// Bucket creation
	s3buckets := template.GetAllAWSS3BucketResources()
	for _, bucketName := range s3buckets {
		createBucket(bucketName.BucketName, serviceS3)
	}
}
