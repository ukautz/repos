package commands

import (
	"github.com/ukautz/cli"
)

func runRemove(c *cli.Cli, o *cli.Command) {
	if idx, store, err := readIndex(o); err != nil {
		c.Output.Die(err.Error())
	} else {
		name := o.Argument("name").String()
		if !idx.Remove(name, "") {
			c.Output.Die("No repo named \"%s\" registered", name)
		} else if err = storeIndex(idx, store); err != nil {
			c.Output.Die(err.Error())
		} else {
			c.Output.Printf("Successfully removed repo \"%s\"\n", name)
		}
	}
}

func init() {
	cmd := cli.NewCommand("remove", "Remove a named repository from index", runRemove)
	addCommandDefaults(cmd)
	cmd.NewArgument("name", "Name of the to-be-removed repo", "", true, false)
	Commands = append(Commands, cmd)
}
