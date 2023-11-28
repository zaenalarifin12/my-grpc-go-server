package grpc

import (
	"fmt"
	"github.com/zaenalarifin12/grpc-course/protogen/go/bank"
	"github.com/zaenalarifin12/grpc-course/protogen/go/hello"
	resl "github.com/zaenalarifin12/grpc-course/protogen/go/resiliency"
	"github.com/zaenalarifin12/my-grpc-go-server/internal/port"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"log"
	"net"
)

type GrpcAdapter struct {
	helloService      port.HelloServicePort
	bankService       port.BankServicePort
	resiliencyService port.ResiliencyServicePort
	grpcPort          int
	server            *grpc.Server
	hello.HelloServiceServer
	bank.BankServiceServer
	resl.ResiliencyServiceServer
	resl.ResiliencyWithMetadataServiceServer
}

func NewGrpcAdapter(helloServer port.HelloServicePort, bankService port.BankServicePort, resiliencyService port.ResiliencyServicePort, grpcPort int) *GrpcAdapter {
	return &GrpcAdapter{
		helloService:      helloServer,
		bankService:       bankService,
		grpcPort:          grpcPort,
		resiliencyService: resiliencyService,
	}
}

func (a *GrpcAdapter) Run() {
	var err error

	listen, err := net.Listen("tcp", fmt.Sprintf(":%d", a.grpcPort))

	if err != nil {
		log.Fatalf("Failed to listen on port %d: %v\n", a.grpcPort, err)
	}

	log.Printf("server listening on port %v", a.grpcPort)

	cred, err := credentials.NewServerTLSFromFile("ssl/server.crt", "ssl/server.pem")

	if err != nil {
		log.Fatalln("Can't create server credentials :", err)
	}

	grpcServer := grpc.NewServer(
		grpc.Creds(cred),
		//grpc.ChainUnaryInterceptor(
		//	interceptor.LogUnaryServerInterceptor(),
		//	interceptor.BasicUnaryServerInterceptor(),
		//),
		//grpc.ChainStreamInterceptor(
		//	interceptor.LogStreamServerInterceptor(),
		//	interceptor.BasicStreamServerInterceptor(),
		//),
	)
	a.server = grpcServer

	hello.RegisterHelloServiceServer(grpcServer, a)
	bank.RegisterBankServiceServer(grpcServer, a)
	resl.RegisterResiliencyServiceServer(grpcServer, a)
	resl.RegisterResiliencyWithMetadataServiceServer(grpcServer, a)

	if err := grpcServer.Serve(listen); err != nil {
		log.Fatalf("Failed to serve gRPC on port %d : %v\n", a.grpcPort, err)
	}
}

func (a *GrpcAdapter) Stop() {
	a.server.Stop()
}
