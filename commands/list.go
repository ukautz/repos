package commands

import (
	"github.com/ukautz/cli"
)

func runList(c *cli.Cli, o *cli.Command) {
	if idx, _, err := readIndex(o); err != nil {
		c.Output.Die(err.Error())
	} else {
		c.Output.Printf("Found %d repos:\n", len(idx.Repos))
		for _, repo := range idx.Repos {
			if hdl := repo.Handler(); hdl != nil {
				if state, err := hdl.State(repo); err != nil {
					c.Output.Printf("%-30s: %s (%s)\n", repo.Name, state, err)
				} else {
					c.Output.Printf("%-30s: %s\n", repo.Name, state)
				}
			} else {
				c.Output.Printf("%-30s: Unknown type %s\n", repo.Name, repo.Type)
			}
		}
	}
}

func init() {
	cmd := cli.NewCommand("list", "List registered repos", runList)
	addCommandDefaults(cmd)
	Commands = append(Commands, cmd)
}
