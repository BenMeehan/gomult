package main

import (
	"net"

	log "github.com/sirupsen/logrus"

	"github.com/benmeehan/gomult/compile"

	"google.golang.org/grpc"
)

func main() {
	lis, err := net.Listen("tcp", ":9000")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	c := compile.Compiler{}
	compile.RegisterCompileServiceServer(grpcServer, &c)

	log.Info("Compiler gRPC Server running on on PORT :9000")

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %s", err)
	}
}
