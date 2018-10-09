package main

import (
	"reflect"
	"testing"

	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	"golang.org/x/net/context"
)

var ctx = context.Background()
var docker, err = createDockerClient()

func Test_pullImage(t *testing.T) {
	type args struct {
		ctx       context.Context
		imageName string
		docker    *client.Client
		err       error
	}
	tests := []struct {
		name string
		args args
	}{
		{"error: Image doesn't exist", args{ctx, "random", docker, err}},
		{"success: Image correct", args{ctx, "alpine", docker, err}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pullImage(tt.args.ctx, tt.args.imageName, tt.args.docker, tt.args.err)
		})
	}
}

// func Test_createContainer(t *testing.T) {
// 	var uiPort, _ = nat.NewPort("tcp", "8080")

// 	type args struct {
// 		ctx       context.Context
// 		docker    *client.Client
// 		imageName string
// 		services  map[int]serviceNamePort
// 	}
// 	tests := []struct {
// 		name string
// 		args args
// 		want response
// 	}{
// 		// {"error: Image doesn't exist",
// 		// 	args{
// 		// 		ctx, docker, "random", nil,
// 		// 	},
// 		// 	response{container.ContainerCreateCreatedBody{}, nil},
// 		// },
// 		{"success: Image correct",
// 			args{
// 				ctx, docker, "localstack/localstack", map[int]serviceNamePort{0: serviceNamePort{"", uiPort}},
// 			},
// 			// response{container.ContainerCreateCreatedBody{}, nil},
// 			response{mock.Arguments, nil},
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			if got := createContainer(tt.args.ctx, tt.args.docker, tt.args.imageName, tt.args.services); !reflect.DeepEqual(got, tt.want) {
// 				t.Errorf("createContainer() = %v, want %v", got, tt.want)
// 			}
// 		})
// 	}
// }

func Test_portBindings(t *testing.T) {
	var s3Port, _ = nat.NewPort("tcp", "4572")
	var uiPort, _ = nat.NewPort("tcp", "8080")

	type args struct {
		services map[int]serviceNamePort
	}
	tests := []struct {
		name string
		args args
		want nat.PortMap
	}{
		{"none service",
			args{map[int]serviceNamePort{
				0: {"", uiPort},
			}},
			map[nat.Port][]nat.PortBinding{
				uiPort: []nat.PortBinding{nat.PortBinding{"0.0.0.0", uiPort.Port()}},
			},
		},
		{"S3 service",
			args{map[int]serviceNamePort{
				0: {"", uiPort},
				1: {"S3", s3Port},
			}},
			map[nat.Port][]nat.PortBinding{
				uiPort: []nat.PortBinding{nat.PortBinding{"0.0.0.0", uiPort.Port()}},
				s3Port: []nat.PortBinding{nat.PortBinding{"0.0.0.0", s3Port.Port()}},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := portBindings(tt.args.services); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("portBindings() = %v, want %v", got, tt.want)
			}
		})
	}
}

// func Test_exposedPorts(t *testing.T) {
// 	type args struct {
// 		services map[int]serviceNamePort
// 	}
// 	tests := []struct {
// 		name string
// 		args args
// 		want nat.PortSet
// 	}{
// 		// TODO: Add test cases.
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			if got := exposedPorts(tt.args.services); !reflect.DeepEqual(got, tt.want) {
// 				t.Errorf("exposedPorts() = %v, want %v", got, tt.want)
// 			}
// 		})
// 	}
// }
