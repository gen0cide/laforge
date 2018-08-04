package main

import (
	"os"

	"github.com/davecgh/go-spew/spew"
	"github.com/gen0cide/laforge"
)

func main() {
	loader := laforge.NewLoader()

	if len(os.Args) < 2 {
		panic("need at least one argument as a laforge config file")
	}

	for _, c := range os.Args[1:] {
		err := loader.ParseConfigFile(c)
		if err != nil {
			panic(err)
		}
	}

	config, err := loader.Bind()
	if err != nil {
		panic(err)
	}

	_ = config
	spew.Dump(config)
}
