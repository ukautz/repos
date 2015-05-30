package commands

import (
	"encoding/json"
	"fmt"
	"github.com/ukautz/cli"
	"github.com/ukautz/repos/index"
	"io/ioutil"
	"os"
)

func addCommandDefaults(cmd *cli.Command) *cli.Command {
	home := os.Getenv("HOME")
	return cmd.NewOption("store", "s", "Location of the index file", home+string(os.PathSeparator)+".repos.json", true, false)
}

func readIndex(o *cli.Command) (*index.Index, string, error) {
	idx := index.NewIndex()
	if store := o.Option("store").String(); store == "" {
		return nil, "", fmt.Errorf("Missing --store location")
	} else if raw, err := ioutil.ReadFile(store); err != nil {
		return nil, "", fmt.Errorf("Failed to open index store \"%s\": %s", store, err)
	} else if err = json.Unmarshal(raw, idx); err != nil {
		return nil, "", fmt.Errorf("Could not parse index in \"%s\": %s", store, err)
	} else {
		return idx, store, nil
	}
}

func storeIndex(idx *index.Index, store string) error {
	if raw, err := json.MarshalIndent(idx, "", "  "); err != nil {
		return fmt.Errorf("Failed to marshal index into JSON: %s", err)
	} else if err := ioutil.WriteFile(store, raw, 0600); err != nil {
		return fmt.Errorf("Failed to write index to \"%s\": %s", store, err)
	} else {
		return nil
	}
}

var (
	Commands = make([]*cli.Command, 0)
)
