package index

import (
	"bytes"
	"os/exec"
)

type Git struct {
}

func (this *Git) State(repo *Repo) (State, error) {
	cmd := exec.Command("git", "status", "--porcelain")
	cmd.Dir = repo.Directory
	out := bytes.NewBuffer(nil)
	cmd.Stdout = out
	if err := cmd.Run(); err != nil {
		return UNKNOWN, err
	} else if str := out.String(); str == "" {
		return UNCHANGED, nil
	} else {
		return CHANGED, nil
	}
}

func init() {
	handlers["git"] = &Git{}
}
