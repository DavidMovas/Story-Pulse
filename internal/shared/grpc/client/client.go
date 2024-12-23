package client

import (
	"fmt"
	"google.golang.org/grpc"
)

func CreateServiceClient[T any](addr string, creationFun func(cc grpc.ClientConnInterface) T, opts ...grpc.DialOption) (T, error) {
	conn, err := grpc.NewClient(fmt.Sprintf("dynamic:///%s", addr), opts...)
	if err != nil {
		fmt.Printf("grpc.NewClient failed: %v\n", err)
		return *new(T), nil
	}

	return creationFun(conn), nil
}
