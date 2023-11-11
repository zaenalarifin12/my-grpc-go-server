package main

import (
	mygrpc "github.com/zaenalarifin12/my-grpc-go-server/internal/adapter/grpc"
	app "github.com/zaenalarifin12/my-grpc-go-server/internal/application"
	"log"
)

func main() {
	log.SetFlags(0)
	log.SetOutput(log.Writer())

	hs := &app.HelloService{}

	grpcAdapter := mygrpc.NewGrpcAdapter(hs, 9000)
	grpcAdapter.Run()
}
