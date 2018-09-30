package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/awslabs/goformation"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/go-connections/nat"

	"golang.org/x/net/context"
)

// Image docker for localstack
var dockerImage = "localstack/localstack"

// Service supported
var serviceSns = serviceNamePort{name: "sns", port: "4575/tcp"}
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

	// Get file as parameter and parse it
	templateFile := flag.String("f", "", "cloudformation template file")

	flag.Parse()

	if len(*templateFile) == 0 {
		flag.Usage()
		os.Exit(1)
	}

	if _, err := os.Stat(*templateFile); os.IsNotExist(err) {
		fmt.Fprintf(os.Stderr, "%s not found\n", *templateFile)
		os.Exit(2)
	}

	template, err := goformation.Open(*templateFile)
	if err != nil {
		log.Fatalf("There was an error processing the template: %s", err)
	}

	// Create docker client and download image if it's needed
	docker, err := createDockerClient()
	pullImage(ctx, dockerImage, docker, err)

	// Get all service inside cloudformation template provided
	services := getAllServices(template)

	// Create docker container with all service found in your cloudformation template
	response := createContainer(ctx, docker, dockerImage, services)

	// Start docker container
	startContainer(ctx, docker, response.resp, response.err)
	for i := range services {
		fmt.Println(services[i].name, " service created")
	}

	// Wait until services is up
	time.Sleep(20 * time.Second)

	// Create all resources found in your cloudformation template that it support
	createResources(template)
}
