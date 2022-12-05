package compile

import (
	"context"
)

type Compiler struct {
}

func (c *Compiler) Compile(ctx context.Context, in *Input) (*Output, error) {
	var res string
	switch in.Lang {
	case "python3":
		res = compilePython3(in.Body)
	}
	return &Output{Result: res}, nil
}
