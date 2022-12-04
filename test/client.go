package main

import (
	"log"

	"golang.org/x/net/context"
	"google.golang.org/grpc"

	"benmeehan111/go-multicompiler/code"
)

func main() {

	var conn *grpc.ClientConn
	conn, err := grpc.Dial(":9000", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %s", err)
	}
	defer conn.Close()

	c := code.NewCompileServiceClient(conn)

	response, err := c.Compile(context.Background(), &code.Input{Lang: "python3", Body: `print("hello")`})
	if err != nil {
		log.Fatalf("Error when calling SayHello: %s", err)
	}
	log.Printf("Response from server: %s", response.Result)

}
