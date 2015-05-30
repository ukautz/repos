package commands

import (
	"github.com/ukautz/cli"
	"path/filepath"
	"encoding/json"
	"io/ioutil"
)

func runAdd(c *cli.Cli, o *cli.Command) {
	if dir := o.Argument("directory").String(); dir == "" {
		c.Output.Die("Missing directory")
	} else if typ := o.Option("type").String(); typ == "" {
		c.Output.Die("Invalid empty --type")
	} else if abs, err := filepath.Abs(dir); err != nil {
		c.Output.Die("Could not get absolute file path from \"%s\": %s", dir, err)
	} else if idx, store, err := readIndex(o); err != nil {
		c.Output.Die(err.Error())
	} else {
		name := o.Argument("name").String()
		if name == "" {
			name = filepath.Base(abs)
		}
		if err := idx.Add(name, abs, typ); err != nil {
			c.Output.Die("Failed to add \"%s\" in \"%s\": %s", name, dir, err)
		} else if raw, err := json.MarshalIndent(idx, "", "  "); err != nil {
			c.Output.Die("Failed to marshal index into JSON: %s", err)
		} else if err := ioutil.WriteFile(store, raw, 0600); err != nil {
			c.Output.Die("Failed to write index to \"%s\": %s", store, err)
		} else {
			c.Output.Printf("Successfully added directory \"%s\" as \"%s\"\n", abs, name)
		}
	}
}

func init() {
	cmd := cli.NewCommand("add", "Add directory to list of repositories", runAdd)
	addCommandDefaults(cmd)
	cmd.NewArgument("directory", "Path to directory", ".", true, false)
	cmd.NewArgument("name", "Name the repo is stored under", "", false, false)
	cmd.NewOption("type", "t", "Type of repository", "git", false, false)
	Commands = append(Commands, cmd)
}
