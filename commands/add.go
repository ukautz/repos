package commands

import (
	"fmt"
	"github.com/ukautz/repos/common"
	"gopkg.in/ukautz/clif.v1"
	"path/filepath"
)

func cmdAdd() *clif.Command {
	cb := func(c *clif.Command, in clif.Input, out clif.Output, lst *common.List) error {
		name := c.Argument("name").String()
		directory := c.Argument("directory").String()
		abs, err := filepath.Abs(directory)
		if err != nil {
			return err
		}
		out.Printf("Adding repository <headline>%s<reset> as <subline>%s<reset>\n", abs, name)
		if p := lst.Get(name); p != "" {
			out.Printf("<warn>There is a watch \"%s\" witch watches \"%s\"<reset>\n", name, p)
			if !in.Confirm("<query>Overwrite?<reset> ") {
				out.Printf("  Not overwriting. Stop.\n")
				return nil
			} else {
				lst.Remove(name)
			}
		}
		if n := lst.Watched(abs); n != "" {
			return fmt.Errorf("Directory \"%s\" is already watched (%s)", abs, n)
		} else if w, err := lst.Add(name, abs); err != nil {
			return err
		} else {
			out.Printf("  Type: <info>%s<reset>\n", w.Type())
			return lst.Persist()
		}
	}

	return clif.NewCommand("add", "Add a new repository to the watch list", cb).
		NewArgument("name", "Name of the repo, so you can remember what it was", "", true, false).
		NewArgument("directory", "Path to directory of the repo. Defaults to current directory.", ".", true, false)

}

func init() {
	Commands = append(Commands, cmdAdd)
}
