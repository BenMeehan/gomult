package main

import (
	"net"
	"net/http"

	log "github.com/sirupsen/logrus"

	"github.com/benmeehan111/gomult/compile"

	"google.golang.org/grpc"
)

func startHttpServer() {
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/", fs)

	log.Info("To use the built in JS test client visit localhost:3000")

	err := http.ListenAndServe(":3000", nil)
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	lis, err := net.Listen("tcp", ":9000")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	c := compile.Compiler{}
	compile.RegisterCompileServiceServer(grpcServer, &c)

	log.Info("Compiler gRPC Server running on on PORT :9000")
	go startHttpServer()

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %s", err)
	}
}
