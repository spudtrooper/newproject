package main

import (
	"flag"

	"github.com/pkg/errors"
	"github.com/spudtrooper/goutil/check"
	"github.com/spudtrooper/newproject/newproject"
)

var (
	outdir   = flag.String("outdir", "", "output directory")
	username = flag.String("username", "spudtrooper", "github username")
	name     = flag.String("name", "", "name of the new project")
)

func realMain() error {
	if *name == "" {
		return errors.Errorf("--name required")
	}
	if *username == "" {
		return errors.Errorf("--username required")
	}
	if err := newproject.Main(*name, *username, newproject.Outdir(*outdir)); err != nil {
		return err
	}
	return nil
}

func main() {
	flag.Parse()
	check.Err(realMain())
}