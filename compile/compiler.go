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
	case "python2.7":
		res = compilePython(in.Body)
	case "nodejs":
		res = compileNodeJS(in.Body)
	case "java17":
		class := getMainClass(in.Body)
		res = compileJava17(in.Body, class)
	case "go":
		res = compileGolang(in.Body)
	}
	return &Output{Result: res}, nil
}
