package index

import "fmt"

type State string

const (
	UNKNOWN   State = "unknown"
	CHANGED   State = "changed"
	UNCHANGED State = "unchanged"
)

var (
	handlers = make(map[string]Handler)
)

type Handler interface {
	State(repo *Repo) (State, error)
}

type Repo struct {
	Name      string
	Directory string
	Type      string
}

func (this *Repo) Handler() Handler {
	if h, ok := handlers[this.Type]; ok {
		return h
	} else {
		return nil
	}
}

type Index struct {
	Repos []*Repo
}

func NewIndex() *Index {
	return &Index{
		Repos: make([]*Repo, 0),
	}
}

func (this *Index) Add(name, dir, typ string) error {
	if exist := this.Get(name, dir); exist != nil {
		return fmt.Errorf("Repo with name \"%s\" or directory \"%s\" already registered", name, dir)
	} else {
		this.Repos = append(this.Repos, &Repo{
			Name:      name,
			Directory: dir,
			Type:      typ,
		})
		return nil
	}
}

func (this *Index) Get(name, dir string) *Repo {
	for _, repo := range this.Repos {
		if repo.Name == name {
			return repo
		} else if repo.Directory == dir {
			return repo
		}
	}
	return nil
}

func (this *Index) List() []*Repo {
	return this.Repos
}
