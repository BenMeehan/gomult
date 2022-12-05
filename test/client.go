package main

import (
	"log"

	"github.com/benmeehan/gomult/compile"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

func main() {

	var conn *grpc.ClientConn
	conn, err := grpc.Dial(":9000", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %s", err)
	}
	defer conn.Close()

	c := compile.NewCompileServiceClient(conn)

	response, err := c.Compile(context.Background(), &compile.Input{Lang: "python3", Body: `print("hello")`})
	if err != nil {
		log.Fatalf("Error when calling SayHello: %s", err)
	}
	log.Printf("Response from server: %s", response.Result)

}
