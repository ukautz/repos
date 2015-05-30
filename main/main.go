package main
import (
	"github.com/ukautz/cli"
	"github.com/ukautz/repos/commands"
	"github.com/ukautz/repos/index"
)

var (
	Version = "1.0.0"
)

func main() {
	app := cli.New("Repos", Version, "Keep track of local repositories")

	idx := index.NewIndex()
	app.Register(idx)

	for _, cmd := range commands.Commands {
		app.Add(cmd)
	}

	app.Run()
}
