package main

import (
	"os/exec"
	"testing"
)

var cmdEmpty = exec.Command("")

func Test_createTopic(t *testing.T) {
	type args struct {
		topicName string
		cmd       *exec.Cmd
	}
	tests := []struct {
		name string
		args args
	}{
		{"none topic created",
			args{"", cmdEmpty},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			createTopic(tt.args.topicName, tt.args.cmd)
		})
	}
}

func Test_createQueue(t *testing.T) {
	type args struct {
		queueName string
		cmd       *exec.Cmd
	}
	tests := []struct {
		name string
		args args
	}{
		{"none queue created",
			args{"", cmdEmpty},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			createQueue(tt.args.queueName, tt.args.cmd)
		})
	}
}

func Test_confirmSubscription(t *testing.T) {
	type args struct {
		topicName string
		cmd       *exec.Cmd
	}
	tests := []struct {
		name string
		args args
	}{
		{"none subscription confirmed created",
			args{"", cmdEmpty},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			confirmSubscription(tt.args.topicName, tt.args.cmd)
		})
	}
}

func Test_getTopicArn(t *testing.T) {
	var cmdFail = exec.Command("/bin/bash", "-c", "echo test")

	type args struct {
		topicName string
		cmd       *exec.Cmd
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"none topicarc can retrieve",
			args{"", cmdEmpty},
			"",
		},
		{"fail unmarshall",
			args{"", cmdFail},
			"",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getTopicArn(tt.args.topicName, tt.args.cmd); got != tt.want {
				t.Errorf("getTopicArn() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_createBucket(t *testing.T) {
	type args struct {
		bucketName string
		cmd        *exec.Cmd
	}
	tests := []struct {
		name string
		args args
	}{
		{"none bucket created",
			args{"", cmdEmpty},
		}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			createBucket(tt.args.bucketName, tt.args.cmd)
		})
	}
}
