package grpc

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	resl "github.com/zaenalarifin12/grpc-course/protogen/go/resiliency"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"io"
	"log"
	"strconv"
	"time"
)

func dummyRequestMetadata(ctx context.Context) {
	if requestMetadata, ok := metadata.FromIncomingContext(ctx); ok {
		log.Println("Request metadata :")
		for k, v := range requestMetadata {
			log.Printf(" %v : %v \n", k, v)
		}
	} else {
		log.Println("Request metadata not found")
	}
}

func dummyResponseMetadata() metadata.MD {
	md := map[string]string{
		"grpc-server-time":     fmt.Sprint(time.Now().Format("15:04:05")),
		"grpc-server-location": "Jakarta, Indonesia",
		"grpc-response-uuid":   uuid.New().String(),
	}

	return metadata.New(md)
}

func (a *GrpcAdapter) UnaryResiliencyWithMetadata(ctx context.Context, req *resl.ResiliencyRequest) (*resl.ResiliencyResponse, error) {
	log.Println("UnaryResiliency called")
	str, sts := a.resiliencyService.GenerateResiliency(req.MinDelaySecond, req.MaxDelaySecond, req.StatusCodes)

	// read request metdata
	dummyRequestMetadata(ctx)

	if errStatus := generateErrStatus(sts); errStatus != nil {
		return nil, errStatus
	}

	// read response metdata
	grpc.SendHeader(ctx, dummyResponseMetadata())

	return &resl.ResiliencyResponse{DummyString: str}, nil
}

func (a *GrpcAdapter) StreamingResiliencyWithMetadata(req *resl.ResiliencyRequest, stream resl.ResiliencyWithMetadataService_StreamingResiliencyWithMetadataServer) error {
	log.Println("Server StreamingResiliency invoked")

	ctx := stream.Context()

	dummyRequestMetadata(ctx)

	err := stream.SendHeader(dummyResponseMetadata())

	if err != nil {
		log.Println("Error while sending response metadata : ", err)
	}

	for {
		select {
		case <-ctx.Done():
			log.Println("Client cancelled request")
			return nil
		default:
			str, sts := a.resiliencyService.GenerateResiliency(req.MinDelaySecond, req.MaxDelaySecond, req.StatusCodes)

			if errStatus := generateErrStatus(sts); errStatus != nil {
				return errStatus
			}

			stream.Send(&resl.ResiliencyResponse{DummyString: str})
		}
	}

}

func (a *GrpcAdapter) ClientStreamingResiliencyWithMetadata(stream resl.ResiliencyWithMetadataService_ClientStreamingResiliencyWithMetadataServer) error {
	log.Println("ClientStreamingResiliency was invoked")

	i := 0

	for {
		req, err := stream.Recv()

		if err == io.EOF {
			res := resl.ResiliencyResponse{DummyString: fmt.Sprintf("Received %v request from client", strconv.Itoa(i))}

			err := stream.SendHeader(dummyResponseMetadata())
			if err != nil {
				log.Println("Error while sending response metadata : ", err)
			}
			return stream.SendAndClose(&res)
		}

		ctx := stream.Context()
		dummyRequestMetadata(ctx)

		if req != nil {
			_, sts := a.resiliencyService.GenerateResiliency(req.MinDelaySecond, req.MaxDelaySecond, req.StatusCodes)

			if errStatus := generateErrStatus(sts); errStatus != nil {
				return errStatus
			}
		}

		i = i + 1
	}
}

func (a *GrpcAdapter) BiDirectionalResiliencyWithMetadata(stream resl.ResiliencyWithMetadataService_BiDirectionalResiliencyWithMetadataServer) error {
	log.Println("BiDirectionalResiliency called")

	ctx := stream.Context()

	err := stream.SendHeader(dummyResponseMetadata())

	if err != nil {
		log.Println("Error while sending response metadata : ", err)
	}

	for {
		select {
		case <-ctx.Done():
			log.Println("Client cancelled request")
			return nil
		default:
			req, err := stream.Recv()

			if err == io.EOF {
				return nil
			}

			if err != nil {
				log.Fatalln("Error while reading from client : ", err)
			}

			dummyRequestMetadata(ctx)

			str, sts := a.resiliencyService.GenerateResiliency(req.MinDelaySecond, req.MaxDelaySecond, req.StatusCodes)

			if errStatus := generateErrStatus(sts); errStatus != nil {
				return errStatus
			}

			err = stream.Send(&resl.ResiliencyResponse{DummyString: str})

			if err != nil {
				log.Fatalln("Error while sending response to client : ", err)
			}
		}
	}
}
