package client

import (
	"fmt"
	"google.golang.org/grpc"
)

func CreateServiceClient[T any](addr string, creationFun func(cc grpc.ClientConnInterface) T, opts ...grpc.DialOption) (T, error) {
	conn, err := grpc.NewClient(fmt.Sprintf("dynamyc:///%s", addr), opts...)
	if err != nil {
		return any(nil), nil
	}

	return creationFun(conn), nil
}
