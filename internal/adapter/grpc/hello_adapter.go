package grpc

import (
	"context"
	"github.com/zaenalarifin12/grpc-course/protogen/go/hello"
)

func (a *GrpcAdapter) SayHello(ctx context.Context, req *hello.HelloRequest) (*hello.HelloResponse, error) {
	greet := a.helloService.GenerateHello(req.Name)

	return &hello.HelloResponse{Greet: greet}, nil
}
