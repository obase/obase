package cvendor

import (
	"flag"
	"github.com/obase/obase/kits"
	"os"
	"os/exec"
	"path/filepath"
)

func cwd() string {
	cwd, _ := os.Getwd()
	return cwd
}

func Process() {
	var (
		version string
		basedir string
		use4jx3 bool
	)
	flag.StringVar(&version, "version", "latest", "specified the version of cvendor")
	flag.StringVar(&basedir, "basedir", cwd(), "specified the basedir of cvendor")
	flag.BoolVar(&use4jx3, "use4jx3", false, "specified used for jx3 vendor")
	flag.Parse()

	if !kits.IsExists(basedir) {
		os.MkdirAll(basedir, os.ModePerm)
	}
	process(basedir, version)
	// 事后处理
	if use4jx3 {
		processForJX3(basedir)
	}
}

func processForJX3(basedir string) {
	if err := os.RemoveAll(filepath.Join(basedir, "github.com", "obase", "vendor", "github.com", "gin-gonic")); err != nil {
		panic(err)
	}
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
