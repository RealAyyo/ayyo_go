package grpc

import (
	"context"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"
)

func loggingInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	start := time.Now()

	p, ok := peer.FromContext(ctx)
	if !ok {
		log.Println("Could not get peer from context")
		return nil, status.Errorf(codes.Internal, "could not get peer from context")
	}

	clientIP := p.Addr.String()
	resp, err := handler(ctx, req)

	log.Printf("Request - IP:%s Method:%s\tDuration:%s\tError:%v\n", clientIP, info.FullMethod, time.Since(start), err)
	return resp, err
}
