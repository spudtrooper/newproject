package newproject

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/spudtrooper/goutil/io"
	"github.com/spudtrooper/goutil/or"
)

func Main(name, username string, inOpts ...Option) error {
	opts := MakeOptions(inOpts...)

	outdir := or.String(opts.Outdir(), name)
	pkg := name

	rootDir, err := io.MkdirAll(outdir)
	if err != nil {
		return err
	}
	pkgDir, err := io.MkdirAll(rootDir, pkg)
	if err != nil {
		return err
	}
	scriptsDir, err := io.MkdirAll(rootDir, "scripts")
	if err != nil {
		return err
	}

	main, err := writeFile(`	
package main

import (
	"flag"

	"github.com/spudtrooper/goutil/check"
	"github.com/{{.Username}}/{{.Pkg}}/{{.Pkg}}"
)

func realMain() error {
	if err := {{.Pkg}}.Main(); err != nil {
		return err
	}
	return nil
}

func main() {
	flag.Parse()
	check.Err(realMain())
}
`, struct {
		Pkg      string
		Username string
	}{
		Pkg:      pkg,
		Username: username,
	}, rootDir, "main.go")
	if err != nil {
		return err
	}

	lib, err := writeFile(`	
package {{.Pkg}}

import (
	"log"
)

func Main() error {
	log.Println("TODO: Implement!")
	return nil
}
`, struct {
		Pkg string
	}{
		Pkg: pkg,
	}, pkgDir, pkg+".go")
	if err != nil {
		return err
	}

	if _, err := writeFile(`
# {{.Pkg}}

TODO
	`, struct {
		Pkg string
	}{
		Pkg: pkg,
	}, rootDir, "README.md"); err != nil {
		return err
	}

	if _, err := writeFile(`
#!/bin/sh

msg="$@"
if [[ -z "$msg" ]]; then
	msg="update $(date)"
fi
git add .
git commit -m "$msg"
open /Applications/GitHub\ Desktop.app	
		`, struct {
	}{}, scriptsDir, "commit.sh"); err != nil {
		return err
	}

	if err := run(rootDir, "go", "mod", "init", fmt.Sprintf("github.com/%s/%s", username, name)); err != nil {
		return err
	}
	if err := run(rootDir, "go", "mod", "tidy"); err != nil {
		return err
	}

	relMain, err := filepath.Rel(rootDir, main)
	if err != nil {
		return err
	}
	if err := run(rootDir, "go", "fmt", relMain); err != nil {
		return err
	}

	relLib, err := filepath.Rel(rootDir, lib)
	if err != nil {
		return err
	}
	if err := run(rootDir, "go", "fmt", relLib); err != nil {
		return err
	}

	if err := run(rootDir, "go", "run", relMain); err != nil {
		return err
	}

	return nil
}

func writeFile(t string, data interface{}, dir string, outFileName string) (string, error) {
	b, err := renderTemplate(t, outFileName, data)
	if err != nil {
		return "", err
	}
	outFile := path.Join(dir, outFileName)
	if err := ioutil.WriteFile(outFile, b, 7055); err != nil {
		return "", err
	}

	log.Printf("wrote to %s", outFile)
	return outFile, nil
}

func run(dir, command string, args ...string) error {
	log.Printf("running from %s: %s %s", dir, command, strings.Join(args, " "))
	cmd := exec.Command(command, args...)
	cmd.Dir = dir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return err
	}
	return nil
}

func renderTemplate(t string, name string, data interface{}) ([]byte, error) {
	tmpl, err := template.New(name).Parse(strings.TrimSpace(t))
	if err != nil {
		return nil, err
	}
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
