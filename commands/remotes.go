package commands

import (
	"github.com/ukautz/repos/common"
	"gopkg.in/ukautz/clif.v0"
)

func cmdRemotes() *clif.Command {
	cb := func(out clif.Output, lst *common.List) {
		watches := lst.List()
		out.Printf("Found <headline>%d<reset> watches\n", len(watches))
		for _, watch := range watches {
			out.Printf(" <info>%s<reset>\n", watch.Name)
			if remotes, err := watch.Watch.Remotes(); err != nil {
				out.Printf("  <error>%s<reset>\n\n", err)
			} else {
				for _, remote := range remotes {
					out.Printf("  %s\n", remote)
				}
				out.Printf("\n")
			}
		}
	}

	return clif.NewCommand("remotes", "List all registerd watches", cb)
}

func init() {
	Commands = append(Commands, cmdRemotes)
}
