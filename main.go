package main

import (
	"benmeehan111/go-multicompiler/code"
	"log"
	"net"
	"net/http"

	"google.golang.org/grpc"
)

func startHttpServer() {
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/", fs)

	log.Print("To use the built in JS client visit localhost:3000")
	err := http.ListenAndServe(":3000", nil)
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	go startHttpServer()
	lis, err := net.Listen("tcp", ":9000")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()

	c := code.Compiler{}
	code.RegisterCompileServiceServer(grpcServer, &c)

	log.Print("gRPC Server running on on PORT 9000")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %s", err)
	}
}
