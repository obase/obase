package cvendor

import (
	"github.com/obase/obase/kits"
	"os"
	"os/exec"
	"path/filepath"
)

func Process(args ...string) {
	var version string
	var basedir string
	if len(args) > 1 {
		basedir = args[1]
	} else {
		basedir, _ = os.Getwd()
	}
	if len(args) > 0 {
		version = args[0]
	} else {
		version = "latest"
	}

	if !kits.IsExists(basedir) {
		os.MkdirAll(basedir, os.ModePerm)
	}
	process(basedir, version)
}

func process(dir string, version string) {

	kits.Infof("write main.go, waiting.......")
	writeFile(filepath.Join(dir, "main.go"), main_go)
	kits.Infof("write go.mod, waiting.......")
	writeFile(filepath.Join(dir, "go.mod"), kits.GetTpl(go_mod, map[string]interface{}{"Version": version}))
	os.Setenv("GO111MODULE", "on")
	kits.Infof("exec go mod tidy, waiting.......")
	command("go", "mod", "tidy")
	kits.Infof("exec go mod vendor, waiting.......")
	command("go", "mod", "vendor")

	oldobasedir := filepath.Join(dir, "vendor", "github.com", "obase")
	newobasedir := filepath.Join(dir, "github.com", "obase")
	oldvendordir := filepath.Join(dir, "vendor")
	newvendordir := filepath.Join(dir, "github.com", "obase", "vendor")

	kits.Infof("move github.com/obase to vendor, waiting.......")
	parent := filepath.Dir(newobasedir)
	if !kits.IsExists(parent) {
		os.MkdirAll(parent, os.ModePerm)
	}

	if err := os.Rename(oldobasedir, newobasedir); err != nil {
		panic(err)
	}
	kits.Infof("move vendor to github.com/obase, waiting.......")
	parent = filepath.Dir(newvendordir)
	if err := os.Rename(oldvendordir, newvendordir); err != nil {
		panic(err)
	}

	kits.Infof("process successfully, please enjoyed!")
}

func writeFile(path string, content string) error {
	maingo, err := os.OpenFile(path, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0666)
	if err != nil {
		return err
	}
	defer maingo.Close()

	maingo.WriteString(content)
	return nil
}

func command(command string, args ...string) error {
	cmd := exec.Command(command, args...)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	return cmd.Run()
}

func Usage() string {
	return "cvendor [dir], 在dir生成obase的vendor(自含全部依赖)"
}
