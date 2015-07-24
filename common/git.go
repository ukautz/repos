package common

import (
	"bufio"
	"bytes"
	"fmt"
	. "github.com/ukautz/repos/common/debug"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
)

// Git is a watch of a Git repository
type (
	Git struct {
		name string
		path string
	}

	gitRemote struct {
		name     string
		url      string
		pushable bool
	}
)

func (this *Git) Changes() (bool, error) {
	if lines, err := this.exec("status", "--porcelain"); err != nil {
		return false, err
	} else if len(lines) == 0 {
		return false, nil
	} else {
		return true, nil
	}
}

func (this *Git) Remotes() ([]string, error) {
	if remotes, err := this.remotes(); err != nil {
		return nil, err
	} else {
		names := make([]string, len(remotes))
		for i, remote := range remotes {
			names[i] = remote.name
		}
		return names, nil
	}
}

func (this *Git) Synced() (SyncStateNum, error) {
	if branches, err := this.branches(); err != nil {
		return SYNC_STATE_FAIL, err
	} else if remotes, err := this.remotes(); err != nil {
		return SYNC_STATE_FAIL, err
	} else {
		for _, remote := range remotes {
			if err := this.fetch(remote.name); err != nil {
				return SYNC_STATE_FAIL, err
			} else if remoteBranches, err := this.remoteBranches(remote.name); err != nil {
				return SYNC_STATE_FAIL, err
			} else {
				for _, branch := range branches {
					if _, ok := remoteBranches[branch]; !ok {
						continue
					}
					if lines, err := this.diff(remote.name, branch); err != nil {
						return SYNC_STATE_FAIL, err
					} else if len(lines) > 0 {
						if ahead, err := this.isAhead(remote.name, branch); err != nil {
							return SYNC_STATE_FAIL, err
						} else if ahead {
							return SYNC_STATE_AHEAD, nil
						} else {
							return SYNC_STATE_BEHIND, nil
						}
					}
				}
			}
		}
		return SYNC_STATE_SAME, nil
	}
}

func (this *Git) States() ([]*SyncState, error) {
	if branches, err := this.branches(); err != nil {
		return nil, err
	} else if remotes, err := this.remotes(); err != nil {
		return nil, err
	} else {
		states := make([]*SyncState, 0)
		for _, remote := range remotes {
			if err := this.fetch(remote.name); err != nil {
				return nil, err
			} else if remoteBranches, err := this.remoteBranches(remote.name); err != nil {
				return nil, err
			} else {
				for _, branch := range branches {
					if _, ok := remoteBranches[branch]; !ok {
						states = append(states, &SyncState{
							Remote: remote.name,
							Branch: branch,
							State:  SYNC_STATE_MISSING,
						})
					} else if lines, err := this.diff(remote.name, branch); err != nil {
						states = append(states, &SyncState{
							Remote: remote.name,
							Branch: branch,
							State:  SYNC_STATE_FAIL,
							Error:  err,
						})
					} else if len(lines) > 0 {
						if !remote.pushable {

						} else if behind, err := this.isAhead(remote.name, branch); err != nil {
							states = append(states, &SyncState{
								Remote: remote.name,
								Branch: branch,
								State:  SYNC_STATE_FAIL,
								Error:  err,
							})
						} else if behind {
							states = append(states, &SyncState{
								Remote: remote.name,
								Branch: branch,
								State:  SYNC_STATE_BEHIND,
							})
						} else {
							states = append(states, &SyncState{
								Remote: remote.name,
								Branch: branch,
								State:  SYNC_STATE_AHEAD,
							})
						}
					} else {
						states = append(states, &SyncState{
							Remote: remote.name,
							Branch: branch,
							State:  SYNC_STATE_SAME,
						})
					}
				}
			}
		}
		return states, nil
	}
}

func (this *Git) Type() string {
	return "Git"
}

func (this *Git) Updates() (bool, error) {
	return false, nil
}

// fetch fetches remote
func (this *Git) fetch(remote string) error {
	_, err := this.exec("fetch", remote)
	return err
}

// branches returns list of local branches
func (this *Git) branches() ([]string, error) {
	if lines, err := this.exec("branch", "--no-color"); err != nil {
		return nil, err
	} else {
		branches := make([]string, len(lines))
		rx := regexp.MustCompile(`^\*?\s+`)
		for i, line := range lines {
			branches[i] = rx.ReplaceAllString(line, "")
		}
		return branches, nil
	}
}

// diff returns list of files which have changed between local and remote
func (this *Git) diff(remote, branch string) ([]string, error) {
	return this.exec("diff", "--name-only", fmt.Sprintf("%s/%s", remote, branch))
}

// isAhead returns whether REMOTE branch is ahead of LOCAL branch
func (this *Git) isAhead(remote, branch string) (bool, error) {
	if lines, err := this.exec("push", "--dry-run", remote, branch); lines == nil || len(lines) == 0 {
		if err != nil {
			return false, err
		} else {
			return false, fmt.Errorf("Oopsie.. something went very wrong")
		}
	} else {
		for _, line := range lines {
			if strings.Index(line, "[rejected]") != -1 {
				return true, nil
			} else if strings.Index(line, "the tip of your current branch is behind") != -1 {
				return true, nil
			}
		}
		return false, nil
	}
}

// remotes returns list of remote repos
func (this *Git) remotes() ([]*gitRemote, error) {
	if lines, err := this.exec("remote", "-v"); err != nil {
		return nil, err
	} else {
		rxSplit := regexp.MustCompile(`\s+`)
		rxPush := regexp.MustCompile(`git@`)
		rxGithub := regexp.MustCompile(`github.com`)
		remotes := make([]*gitRemote, 0)
		seen := make(map[string]bool)
		for _, line := range lines {
			p := rxSplit.Split(line, 3)
			if len(p) == 3 && !seen[p[0]] {
				seen[p[0]] = true
				pushable := true
				if rxGithub.MatchString(p[1]) {
					pushable = rxPush.MatchString(p[1])
				}
				remotes = append(remotes, &gitRemote{
					name:     p[0],
					url:      p[1],
					pushable: pushable,
				})
			}
		}
		return remotes, nil
	}
}

// remotes returns list of remote repos
func (this *Git) remoteBranches(remote string) (map[string]bool, error) {
	if lines, err := this.exec("branch", "-a"); err != nil {
		return nil, err
	} else {
		rx := regexp.MustCompile(`^\s*remotes/(.+)$`)
		branches := map[string]bool{}
		prefix := remote + "/"
		for _, line := range lines {
			if rx.MatchString(line) {
				p := rx.FindStringSubmatch(line)
				if strings.Index(p[1], prefix) == 0 {
					branches[p[1][len(prefix):]] = true
				}
			}
		}
		return branches, nil
	}
}

// remotes returns list of remote repos
func (this *Git) exec(args ...string) ([]string, error) {
	Debug(DEBUG2, "Git exec [%s: %s]: %s", this.name, this.path, strings.Join(args, " "))
	cmd := exec.Command("git", args...)
	cmd.Dir = this.path
	errOut := bytes.NewBuffer(nil)
	stdOut := bytes.NewBuffer(nil)
	cmd.Stderr = errOut
	cmd.Stdout = stdOut
	cmd.Stdin = os.Stdin
	err := cmd.Run()
	lines := []string{}
	for n, buf := range map[string]*bytes.Buffer{"err": errOut, "out": stdOut} {
		scn := bufio.NewScanner(buf)
		for scn.Scan() {
			if line := scn.Text(); line != "" {
				lines = append(lines, line)
				Debug(DEBUG3, " %s: %s", n, line)
			}
		}
	}
	return lines, err
}

func init() {
	watches = append(watches, func(path, name string) (Watch, error) {
		git := filepath.Join(path, ".git")
		if stat, err := os.Stat(git); err != nil {
			if os.IsNotExist(err) {
				return nil, nil
			} else {
				return nil, err
			}
		} else if !stat.IsDir() {
			return nil, fmt.Errorf("Found \".git\" in \"%s\", but it is not a directory", path)
		} else {
			return &Git{
				path: path,
				name: name,
			}, nil
		}
	})
}
