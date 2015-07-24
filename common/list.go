package common

import (
	"encoding/json"
	"fmt"
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
		repos map[string]string
	}

	// Info represents full information about a single repo
	Info struct {
		Name, Path, Type string
		Error            error
		Repo             Repo
	}
)

// NewList constructs new List instance
func NewList(path string) *List {
	return &List{
		path:  path,
		repos: make(map[string]string),
	}
}

// Add includes given named repo under path to watches
func (this *List) Add(name, path string) (Repo, error) {
	if w, err := NewRepo(path, name); err != nil {
		return nil, err
	} else {
		this.repos[name] = path
		return w, nil
	}
}

// Get returns path of registered repo or empty string
func (this *List) Get(name string) string {
	if path, ok := this.repos[name]; ok {
		return path
	} else {
		return ""
	}
}

// Info returns Repo
func (this *List) Info(name string) (*Info, error) {
	if path, ok := this.repos[name]; ok {
		if repo, err := NewRepo(path, name); err != nil {
			return nil, err
		} else {
			return &Info{
				Name: name,
				Path: path,
				Type: repo.Type(),
				Repo: repo,
			}, nil
		}
	} else {
		return nil, fmt.Errorf("Repo not found")
	}
}

// List returns list of all registered repos.
func (this *List) List() []*Info {
	names := []string{}
	for name, _ := range this.repos {
		names = append(names, name)
	}
	sort.Strings(names)
	named := make([]*Info, len(names))
	for i, name := range names {
		named[i] = &Info{
			Name: name,
			Path: this.repos[name],
		}
		if watch, err := NewRepo(this.repos[name], name); err != nil {
			named[i].Type = "UNDEF"
			named[i].Error = err
		} else {
			named[i].Type = watch.Type()
			named[i].Repo = watch
		}
	}
	return named
}

// Persist writes watched repos to storage
func (this *List) Persist() error {
	if raw, err := json.MarshalIndent(this.repos, "", "  "); err != nil {
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
			this.repos = m
			return nil
		} else {
			return err
		}
	} else if err = json.Unmarshal(raw, &m); err != nil {
		return err
	} else {
		this.repos = m
		return nil
	}
}

// Remove watched repository by name
func (this *List) Remove(name string) bool {
	if _, ok := this.repos[name]; ok {
		delete(this.repos, name)
		return true
	} else {
		return false
	}
}

// Watch returns name if path is already watched and empty string if it's not
func (this *List) Watched(path string) string {
	for name, watch := range this.repos {
		if watch == path {
			return name
		}
	}
	return ""
}
