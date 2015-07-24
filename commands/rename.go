package commands

import (
	"fmt"
	"github.com/ukautz/repos/common"
	"gopkg.in/ukautz/clif.v0"
)

func cmdRename() *clif.Command {
	cb := func(c *clif.Command, out clif.Output, lst *common.List) error {
		oldName := c.Argument("old-name").String()
		newName := c.Argument("new-name").String()
		if path := lst.Get(oldName); path == "" {
			return fmt.Errorf("No repo with name \"%s\" found")
		} else if existingPath := lst.Get(newName); existingPath != "" {
			return fmt.Errorf("Repo with name \"%s\" already exists: %s", newName, existingPath)
		} else {
			lst.Remove(oldName)
			if _, err := lst.Add(newName, path); err != nil {
				return fmt.Errorf("Failed to re-add repo in \"%s\" under new name %s: %s", path, newName, err)
			} else if err = lst.Persist(); err != nil {
				return fmt.Errorf("Failed to persist repos: %s", err)
			}
			out.Printf("Renamed <info>%s<reset> to <info>%s<reset>\n", oldName, newName)
			return nil
		}
	}

	return clif.NewCommand("rename", "Rename registered repo", cb).
		NewArgument("old-name", "Old name of the repo", "", true, false).
		NewArgument("new-name", "New name of the repo", "", true, false)

}

func init() {
	Commands = append(Commands, cmdRename)
}
