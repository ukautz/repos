package commands

import (
	. "github.com/ukautz/repos/common/debug"
	"github.com/ukautz/repos/common"
	"gopkg.in/ukautz/clif.v0"
	"regexp"
	"sync"
	"github.com/cheggaaa/pb"
	"fmt"
)

func cmdCheck() *clif.Command {
	cb := func(c *clif.Command, out clif.Output, lst *common.List) error {

		// determine repo name pattern
		includeRx := regexp.MustCompile(`.`)
		if pattern := c.Option("include").String(); pattern != "" {
			var err error
			if includeRx, err = regexp.Compile(pattern); err != nil {
				return fmt.Errorf("Failed to compile include pattern: %s", err)
			}
		}
		var excludeRx *regexp.Regexp
		if pattern := c.Option("exclude").String(); pattern != "" {
			var err error
			if excludeRx, err = regexp.Compile(pattern); err != nil {
				return fmt.Errorf("Failed to compile exclude pattern: %s", err)
			}
		}

		// reduce list of repos
		all := lst.List()
		watches := make([]*common.Repo, 0)
		for _, repo := range all {
			if !includeRx.MatchString(repo.Name) {
				Debug(DEBUG3, "Exclude repo %s since not matching include", repo.Name)
				continue
			} else if excludeRx != nil && excludeRx.MatchString(repo.Name) {
				Debug(DEBUG3, "Exclude repo %s since match exclude", repo.Name)
				continue
			}
			watches = append(watches, repo)
		}
		if len(watches) == 0 {
			out.Printf("<warn>No repos found<reset>\n")
			return nil
		}

		// starting now
		out.Printf("Checking <headline>%d<reset> watches\n", len(watches))
		errs := []*common.Repo{}
		changed := []*common.Repo{}
		ahead := []*common.Repo{}
		behind := []*common.Repo{}
		var wg sync.WaitGroup
		mux := new(sync.Mutex)
		total := len(watches)
		count := 0
		progress := make(chan string)

		wg.Add(1)
		if DebugLevel == DEBUG0 {
			pb := pb.StartNew(len(watches))
			go func() {
				defer wg.Done()
				for range progress {
					pb.Increment()
				}
				pb.Finish()
			}()
		} else {
			go func() {
				defer wg.Done()
				for range progress {
					// ...
				}
			}()
		}

		go func() {
			defer close(progress)
			var wgCheck sync.WaitGroup
			for _, repo := range watches {
				Debug(DEBUG1, "Checking repo %s", repo.Name)
				wgCheck.Add(1)
				go func(repo *common.Repo) {
					defer wgCheck.Done()
					var add *[]*common.Repo
					if repo.Error != nil {
						add = &errs
					} else if changes, err := repo.Watch.Changes(); err != nil {
						repo.Error = err
						add = &errs
					} else if changes {
						add = &changed
					} else if synced, err := repo.Watch.Synced(); err != nil {
						repo.Error = err
						add = &errs
					} else if synced == common.SYNC_STATE_AHEAD {
						add = &ahead
					} else if synced == common.SYNC_STATE_BEHIND {
						add = &behind
					}
					mux.Lock()
					defer mux.Unlock()
					if add != nil {
						*add = append(*add, repo)
						count++
						Debug(DEBUG1, "Done: Repo %s changed (%d of %d)", repo.Name, count, total)
					} else {
						count++
						Debug(DEBUG1, "Done: Repo %s unchanged (%d of %d)", repo.Name, count, total)
					}
					progress <- repo.Name
				}(repo)
			}
			wgCheck.Wait()
		}()
		wg.Wait()

		any := false
		if len(errs) > 0 {
			any = true
			out.Printf("\n Found <headline>%d<reset> with <subline>errros<reset>\n", len(errs))
			for _, watch := range errs {
				out.Printf("  <info>%s<reset>: <error>%s<reset>\n", watch.Name, watch.Error)
			}
		}
		if len(changed) > 0 {
			any = true
			out.Printf("\n Found <headline>%d<reset> with <subline>local changes<reset>\n", len(changed))
			out.Printf("  <debug>Eg uncommited changes<reset>\n")
			for _, watch := range changed {
				out.Printf("  <info>%s<reset> (%s): <warn>%s<reset>\n", watch.Name, watch.Type, watch.Path)
			}
		}
		if len(ahead) > 0 {
			any = true
			out.Printf("\n Found <headline>%d<reset> which are <subline>ahead of local<reset>\n", len(ahead))
			out.Printf("  <debug>Remote has commits which are not merged into local<reset>\n")
			for _, watch := range ahead {
				out.Printf("  <info>%s<reset> (%s): <warn>%s<reset>\n", watch.Name, watch.Type, watch.Path)
			}
		}
		if len(behind) > 0 {
			any = true
			out.Printf("\n Found <headline>%d<reset> which are <subline>behind of local<reset>\n", len(behind))
			out.Printf("  <debug>Local has commits which are not merged with (at least one) remote<reset>\n")
			for _, watch := range behind {
				out.Printf("  <info>%s<reset> (%s): <warn>%s<reset>\n", watch.Name, watch.Type, watch.Path)
			}
		}
		if !any {
			out.Printf(" <success>All is in sync!<reset>\n")
		}

		return nil
	}

	return clif.NewCommand("check", "Check all registered repos", cb).
		NewOption("include", "i", "Regular expression for repos to include, eg '^foo'", "", false, true).
		NewOption("exclude", "e", "Regular expression for repos to exclude, eg 'bar$'", "", false, true)
}

func init() {
	Commands = append(Commands, cmdCheck)
}
