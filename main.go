package main

import (
	"fmt"
	"github.com/obase/obase/cvendor"
	"github.com/obase/obase/kits"
	"os"
	"strings"
	"time"
)

func init() {
	// 添加子命令
	all = append(all, &cmd{"cvendor", cvendor.Process, cvendor.Usage})
}

func main() {

	if len(os.Args) > 1 && os.Args[1] != "help" {
		cmd := os.Args[1]
		os.Args = os.Args[1:]
		for _, c := range all {
			if c.Name == cmd {
				kits.Infof("%v %v start\n", c.Name, os.Args[2:])
				start := time.Now().UnixNano()
				c.Process()
				end := time.Now().UnixNano()
				kits.Infof("%v %v finish, used time(ms): %v\n", c.Name, os.Args[2:], (end-start)/1000000)
				return
			}
		}
	}
	fmt.Fprintln(os.Stderr, usage())
}

func usage() string {
	sb := new(strings.Builder)
	sb.WriteString(fmt.Sprintf("Usage: %s <command> [args...]\n", os.Args[0]))
	for _, c := range all {
		sb.WriteString("\n命令: ")
		sb.WriteString(c.Name)
		sb.WriteString(", 用法: ")
		sb.WriteString(c.Usage())
		sb.WriteString("\n")
	}
	return sb.String()
}

type cmd struct {
	Name    string
	Process func()
	Usage   func() string
}

var all []*cmd
