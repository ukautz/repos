// +build

package main

import (
	"github.com/ukautz/repos/commands"
	"github.com/ukautz/repos/common/debug"
	"github.com/ukautz/repos/list"
	"gopkg.in/ukautz/clif.v0"
	"os"
	"path/filepath"
)

var (
	Version = "1.0.0"
)

func addDefaultOptions(cli *clif.Cli) {
	// use path in --store to provide list.List in injection container so that
	// all commands can access
	storeDefault := filepath.Join(os.Getenv("HOME"), ".repos.json")
	storeOpt := clif.NewOption("store", "s", "Path JSON file storing the registered repos", storeDefault, true, false).
		SetEnv("REPOS_STORE").
		SetParse(func(name, value string) (string, error) {
		lst := list.New(value)
		cli.Register(lst)
		if err := lst.Refresh(); err != nil {
			return value, err
		} else {
			return value, nil
		}
	})
	verboseOpt := clif.NewOption("verbose", "v", "Log verbose messages to STDERR", "", false, true).
		IsFlag().
		SetEnv("REPOS_DEBUG").
		SetParse(func(name, value string) (string, error) {
		if value == "true" {
			if debug.DebugLevel < debug.DEBUG3 {
				debug.DebugLevel++
			}
		} else {
			debug.DebugLevel = 0
		}
		return value, nil
	})
	cli.AddDefaultOptions(storeOpt, verboseOpt)
}

func main() {
	clif.DefaultStyles = clif.SunburnStyles
	cli := clif.New("repos", Version, "Keep track of local repositories")
	addDefaultOptions(cli)

	// add all commands and run the cli
	for _, cb := range commands.Commands {
		cli.Add(cb())
	}
	cli.Run()
}
