package main

import (
	"flag"

	"github.com/spudtrooper/foo/foo"
	"github.com/spudtrooper/goutil/check"
)

func realMain() error {
	if err := foo.Main(); err != nil {
		return err
	}
	return nil
}

func main() {
	flag.Parse()
	check.Err(realMain())
}
