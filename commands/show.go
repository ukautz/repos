package commands

import (
	"fmt"
	"github.com/ukautz/repos/common"
	"gopkg.in/ukautz/clif.v0"
)

func cmdShow() *clif.Command {
	cb := func(out clif.Output, lst *common.List) {
		watches := lst.List()
		out.Printf("Found <headline>%d<reset> watches\n", len(watches))
		max := 0
		for _, watch := range watches {
			if l := len(watch.Name); l > max {
				max = l
			}
		}
		m := fmt.Sprintf("%d", max)
		for _, watch := range watches {
			out.Printf(" <info>%- "+m+"s<reset> ", watch.Name)
			if watch.Error != nil {
				out.Printf("<error>%s<reset>")
			} else {
				out.Printf("<important>%s<reset>", watch.Type)
			}
			out.Printf(" @ %s\n", watch.Path)
		}
	}

	return clif.NewCommand("show", "List all registerd watches", cb)
}

func init() {
	Commands = append(Commands, cmdShow)
}
