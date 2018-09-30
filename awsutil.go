package main

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"

	"github.com/awslabs/goformation/cloudformation"
)

// Get all type of service inside cloudformation template
func getAllServices(template *cloudformation.Template) map[int]serviceNamePort {
	var allServices = make(map[int]serviceNamePort)
	serviceNumber := 0

	s3service := template.GetAllAWSS3BucketResources()
	sqsService := template.GetAllAWSSQSQueueResources()
	snsService := template.GetAllAWSSNSTopicResources()

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

	if len(snsService) != 0 {
		allServices[serviceNumber] = serviceSns
		serviceNumber++
	}

	return allServices
}

// Create all resources for service supported
func createResources(template *cloudformation.Template) {
	createResourcesForSQS(template)
	createResourcesForS3(template)
	createResourcesForSNS(template)
}

// Create all resources available for SQS
func createResourcesForSQS(template *cloudformation.Template) {
	sqsQueues := template.GetAllAWSSQSQueueResources()
	for _, queueName := range sqsQueues {
		createQueue(queueName.QueueName, serviceSqs)
	}
}

// Create all resources available for S3
func createResourcesForS3(template *cloudformation.Template) {
	s3buckets := template.GetAllAWSS3BucketResources()
	for _, bucketName := range s3buckets {
		createBucket(bucketName.BucketName, serviceS3)
	}
}

// Create all resources available for SNS
func createResourcesForSNS(template *cloudformation.Template) {
	topicsSns := template.GetAllAWSSNSTopicResources()
	for _, topicName := range topicsSns {
		createTopic(topicName.TopicName, serviceSns)
	}

	// Subscribe to SNS Topic if needed
	getSnsSubscriptions(template)
}

// Create an empty bucket
func createBucket(bucketName string, serviceS3 serviceNamePort) {
	cmd := exec.Command("/bin/bash", "-c", "aws s3api create-bucket --endpoint-url=http://localhost:"+serviceS3.port.Port()+" --bucket "+bucketName)

	_, errCommand := cmd.Output()

	if errCommand != nil {
		fmt.Println(errCommand.Error())
		fmt.Println("AWS cli should be installed and configured to create resources")
		return
	}
	fmt.Printf("Bucket %q successfully created\n", bucketName)
}

// Get all subscriptions to sns inside cloudformation template
func getSnsSubscriptions(template *cloudformation.Template) {
	snsSubscriptions := template.GetAllAWSSNSSubscriptionResources()
	for _, subscription := range snsSubscriptions {
		protocol := subscription.Protocol
		var subscriptor string
		topic := getTopicArn(subscription.TopicArn)
		if strings.Contains(protocol, "sqs") {
			subscriptor = "http://localhost:" + serviceSqs.port.Port()
		} else {
			subscriptor = subscription.Endpoint
		}
		createSubscription(protocol, subscriptor, topic)
	}
}

// Get arn for given topic sns name. Should be create before instead it return empty
func getTopicArn(topicName string) string {
	cmd := exec.Command("/bin/bash", "-c", "aws sns list-topics --endpoint-url=http://localhost:"+serviceSns.port.Port())

	topicCloudFormation, errCommand := cmd.Output()

	var dat map[string][]map[string]string

	if err := json.Unmarshal(topicCloudFormation, &dat); err != nil {
		panic(err)
	}

	listTopics := dat["Topics"]

	for _, topic := range listTopics {
		candidate := topic["TopicArn"]
		if strings.Contains(candidate, topicName) {
			return candidate
		}
	}

	if errCommand != nil {
		fmt.Println(errCommand.Error())
		fmt.Println("AWS cli should be installed and configured to create resources")
		return ""
	}

	return ""
}

// Subscribe to an sns topic
// Protocol available: sqs, http, https
func createSubscription(protocol string, subscriptor string, topic string) {
	cmd := exec.Command("/bin/bash", "-c", "aws sns subscribe --endpoint-url=http://localhost:"+serviceSns.port.Port()+" --topic-arn "+topic+" --protocol "+protocol+" --notification-endpoint "+subscriptor)

	subscriptionResponse, errCommand := cmd.Output()

	if errCommand != nil {
		fmt.Println(errCommand.Error())
		fmt.Println("AWS cli should be installed and configured to create resources")
		return
	}

	var token map[string]string

	if err := json.Unmarshal(subscriptionResponse, &token); err != nil {
		panic(err)
	}

	confirmSubscription(topic, token["SubscriptionArn"])
}

// Confirm subscription to an sns topic
func confirmSubscription(topicName string, token string) {
	cmd := exec.Command("/bin/bash", "-c", "aws sns confirm-subscription --endpoint-url=http://localhost:"+serviceSns.port.Port()+" --topic-arn  "+topicName+" --token "+token)

	_, errCommand := cmd.Output()

	if errCommand != nil {
		fmt.Println(errCommand.Error())
		fmt.Println("AWS cli should be installed and configured to create resources")
		return
	}

	fmt.Printf("Subscription to %q successfully created\n", topicName)
}

// Create an sns topic
func createTopic(topicName string, serviceSns serviceNamePort) {
	cmd := exec.Command("/bin/bash", "-c", "aws sns create-topic  --endpoint-url=http://localhost:"+serviceSns.port.Port()+" --name "+topicName)

	_, errCommand := cmd.Output()

	if errCommand != nil {
		fmt.Println(errCommand.Error())
		fmt.Println("AWS cli should be installed and configured to create resources")
		return
	}
	fmt.Printf("Topic %q successfully created\n", topicName)
}

// Create a sqs default queue
func createQueue(queueName string, serviceSqs serviceNamePort) {
	cmd := exec.Command("/bin/bash", "-c", "aws sqs create-queue  --endpoint-url=http://localhost:"+serviceSqs.port.Port()+" --queue-name "+queueName)

	_, errCommand := cmd.Output()

	if errCommand != nil {
		fmt.Println(errCommand.Error())
		fmt.Println("AWS cli should be installed and configured to create resources")
		return
	}
	fmt.Printf("SQS queue %q successfully created\n", queueName)
}
