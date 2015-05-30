package commands

import (
	"fmt"
	"github.com/ukautz/cli"
	"sync"
	"github.com/ukautz/repos/index"
)

var (
	states   map[string]string
	stateMux = new(sync.Mutex)
)

func setState(name, state string) {
	stateMux.Lock()
	defer stateMux.Unlock()
	states[name] = state
}

func runList(c *cli.Cli, o *cli.Command) {
	if idx, _, err := readIndex(o); err != nil {
		c.Output.Die(err.Error())
	} else {
		c.Output.Printf("Gathering states of %d repos\n", len(idx.Repos))
		var wg sync.WaitGroup
		states = make(map[string]string)
		for _, repo := range idx.Repos {
			wg.Add(1)
			go func(repo *index.Repo) {
				defer wg.Done()
				if hdl := repo.Handler(); hdl != nil {
					if state, err := hdl.State(repo); err != nil {
						setState(repo.Name, fmt.Sprintf("%s (%s)", state, err))
					} else {
						setState(repo.Name, string(state))
					}
				} else {
					setState(repo.Name, fmt.Sprintf("(unknown type \"%s\")", repo.Type))
				}
			}(repo)
		}
		wg.Wait()
		c.Output.Printf("--\n")
		for name, state := range states {
			c.Output.Printf("%-30s %s\n", name, state)
		}
	}
}

func init() {
	cmd := cli.NewCommand("list", "List registered repos", runList)
	addCommandDefaults(cmd)
	Commands = append(Commands, cmd)
}
