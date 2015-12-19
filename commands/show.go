package commands

import (
	"fmt"
	"github.com/ukautz/repos/common"
	"gopkg.in/ukautz/clif.v1"
)

func cmdShow() *clif.Command {
	cb := func(out clif.Output, lst *common.List) {
		watches := lst.List()
		out.Printf("Found <headline>%d<reset> watches\n", len(watches))

		rowsWith := make([][]string, 0)
		rowsWithout := make([][]string, 0)
		errs := 0
		for _, watch := range watches {
			row := []string{watch.Name, watch.Type, watch.Path, ""}
			if watch.Error != nil {
				errs ++
				row[3] = watch.Error.Error()
			}
			rowsWith = append(rowsWith, row)
			rowsWithout = append(rowsWithout, row[0:3])
		}

		var table *clif.Table
		if errs > 0 {
			table = out.Table([]string{"Name", "Type", "Path", "Error"})
			table.AddRows(rowsWith)
		} else {
			table = out.Table([]string{"Name", "Type", "Path"})
			table.AddRows(rowsWithout)
		}
		fmt.Println(table.Render())
	}

	return clif.NewCommand("show", "Show all registered repos", cb)
}

func init() {
	Commands = append(Commands, cmdShow)
}
