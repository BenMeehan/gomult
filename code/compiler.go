package code

import (
	"context"
	"log"
)

type Compiler struct {
}

func (c *Compiler) Compile(ctx context.Context, in *Input) (*Output, error) {
	log.Printf("Receive message body from client: %s", in.Body)
	var res string
	switch in.Lang{
	case "python3":
		res=compilePython3(in.Body)
	}
	return &Output{Result: res}, nil
}
