package commands

import (
	"fmt"
	"github.com/ukautz/repos/common"
	"gopkg.in/ukautz/clif.v1"
	"path/filepath"
)

func cmdRemove() *clif.Command {
	cb := func(c *clif.Command, out clif.Output, lst *common.List) error {
		name := c.Argument("name").String()
		if lst.Get(name) == "" {
			if abs, err := filepath.Abs(name); err != nil {
				return fmt.Errorf("Found no named watch for %s, so tried directory but could not get abs path: %s\n", name, err)
			} else if n := lst.Watched(abs); n == "" {
				return fmt.Errorf("Found no named watch nor any path for %s\n", name)
			} else {
				name = n
			}
		}
		lst.Remove(name)
		out.Printf("<success>Removed %s from watch list<reset>\n", name)
		return lst.Persist()
	}

	return clif.NewCommand("remove", "Remove a registered repository", cb).
		NewArgument("name", "Name or directory of a repo which shall be removed from watch", "", true, false)

}

func init() {
	Commands = append(Commands, cmdRemove)
}
