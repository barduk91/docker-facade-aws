package main

import (
	"fmt"
	"os/exec"

	"github.com/awslabs/goformation/cloudformation"
)

func getAllServices(template *cloudformation.Template) map[int]serviceNamePort {
	var allServices = make(map[int]serviceNamePort)
	serviceNumber := 0

	s3service := template.GetAllAWSS3BucketResources()
	sqsService := template.GetAllAWSSQSQueueResources()
	dynamoDbService := template.GetAllAWSDynamoDBTableResources()

	// UI to see all infrastructure
	allServices[serviceNumber] = serviceUI
	serviceNumber++

	if len(s3service) != 0 {
		allServices[serviceNumber] = serviceS3
		serviceNumber++
	}
	if len(sqsService) != 0 {
		allServices[serviceNumber] = serviceSqs
		serviceNumber++
	}

	if len(dynamoDbService) != 0 {
		allServices[serviceNumber] = serviceDynamo
		serviceNumber++
	}

	return allServices
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
