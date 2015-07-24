package list

import (
	"encoding/json"
	"github.com/ukautz/repos/watch"
	"io/ioutil"
	"os"
	"sort"
)

type (

	// List contains all watched repos
	List struct {

		// path is where the list is persisted (JSON file)
		path string

		// repos contain a (name => path) map of all watched repos
		watches map[string]string
	}

	// Repo represents a single named repo, which is in a directory and of a type (eg Git)
	Repo struct {
		Name, Path, Type string
		Error            error
		Watch            watch.Watch
	}
)

// New constructs new List instance
func New(path string) *List {
	return &List{
		path:    path,
		watches: make(map[string]string),
	}
}

// Add includes given named repo under path to watches
func (this *List) Add(name, path string) (watch.Watch, error) {
	if w, err := watch.Factory(path, name); err != nil {
		return nil, err
	} else {
		this.watches[name] = path
		return w, nil
	}
}

// Get returns path of registered repo or empty string
func (this *List) Get(name string) string {
	if path, ok := this.watches[name]; ok {
		return path
	} else {
		return ""
	}
}

// Get returns path of registered repo or empty string
func (this *List) Info(name string) (watch.Watch, error) {
	if path, ok := this.watches[name]; ok {
		return path
	} else {
		return ""
	}
}

// List returns list of all registered watches.
func (this *List) List() []*Repo {
	names := []string{}
	for name, _ := range this.watches {
		names = append(names, name)
	}
	sort.Strings(names)
	named := make([]*Repo, len(names))
	for i, name := range names {
		named[i] = &Repo{
			Name: name,
			Path: this.watches[name],
		}
		if watch, err := watch.Factory(this.watches[name], name); err != nil {
			named[i].Type = "UNDEF"
			named[i].Error = err
		} else {
			named[i].Type = watch.Type()
			named[i].Watch = watch
		}
	}
	return named
}

// Persist writes watched repos to storage
func (this *List) Persist() error {
	if raw, err := json.MarshalIndent(this.watches, "", "  "); err != nil {
		return err
	} else if err = ioutil.WriteFile(this.path, raw, 0600); err != nil {
		return err
	} else {
		return nil
	}
}

// Refresh reads watch list from previously persisted storage. Does not error
// if storage does not exist.
func (this *List) Refresh() error {
	m := make(map[string]string)
	if raw, err := ioutil.ReadFile(this.path); err != nil {
		if os.IsNotExist(err) {
			this.watches = m
			return nil
		} else {
			return err
		}
	} else if err = json.Unmarshal(raw, &m); err != nil {
		return err
	} else {
		this.watches = m
		return nil
	}
}

// Remove watched repository by name
func (this *List) Remove(name string) bool {
	if _, ok := this.watches[name]; ok {
		delete(this.watches, name)
		return true
	} else {
		return false
	}
}

// Watch returns name if path is already watched and empty string if it's not
func (this *List) Watched(path string) string {
	for name, watch := range this.watches {
		if watch == path {
			return name
		}
	}
	return ""
}
