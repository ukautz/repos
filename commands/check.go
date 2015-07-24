package commands

import (
	"github.com/cheggaaa/pb"
	"github.com/ukautz/repos/common"
	. "github.com/ukautz/repos/common/debug"
	"gopkg.in/ukautz/clif.v0"
	"sync"
)

func cmdCheck() *clif.Command {
	cb := func(c *clif.Command, out clif.Output, lst *common.List) error {
		repos, err := reduceWithRepoFilters(c, lst.List())
		if err != nil {
			return err
		} else if len(repos) == 0 {
			out.Printf("<warn>No repos found<reset>\n")
			return nil
		}

		// starting now
		out.Printf("Checking <headline>%d<reset> repos\n", len(repos))
		reposWithError := []*common.Info{}
		reposWithLocalChanges := []*common.Info{}
		reposAheadOfRemote := []*common.Info{}
		reposBehindOfRemote := []*common.Info{}
		var wg sync.WaitGroup
		mux := new(sync.Mutex)
		total := len(repos)
		count := 0
		progress := make(chan string)

		wg.Add(1)
		if DebugLevel == DEBUG0 {
			pb := pb.StartNew(len(repos))
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
			for _, repo := range repos {
				Debug(DEBUG1, "Checking repo %s", repo.Name)
				wgCheck.Add(1)
				go func(repo *common.Info) {
					defer wgCheck.Done()
					var add *[]*common.Info
					if repo.Error != nil {
						add = &reposWithError
					} else if changes, err := repo.Repo.Changes(); err != nil {
						repo.Error = err
						add = &reposWithError
					} else if changes {
						add = &reposWithLocalChanges
					} else if synced, err := repo.Repo.Synced(); err != nil {
						repo.Error = err
						add = &reposWithError
					} else if synced == common.SYNC_STATE_AHEAD {
						add = &reposAheadOfRemote
					} else if synced == common.SYNC_STATE_BEHIND {
						add = &reposBehindOfRemote
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
		if len(reposWithError) > 0 {
			any = true
			out.Printf("\n Found <headline>%d<reset> with <subline>errros<reset>\n", len(reposWithError))
			for _, repo := range reposWithError {
				out.Printf("  <info>%s<reset>: <error>%s<reset>\n", repo.Name, repo.Error)
			}
		}
		if len(reposWithLocalChanges) > 0 {
			any = true
			out.Printf("\n Found <headline>%d<reset> with <subline>local changes<reset>\n", len(reposWithLocalChanges))
			out.Printf("  <debug>Eg uncommited changes<reset>\n")
			for _, repo := range reposWithLocalChanges {
				out.Printf("  <info>%s<reset> (%s): <important>%s<reset>\n", repo.Name, repo.Type, repo.Path)
			}
		}
		if len(reposAheadOfRemote) > 0 {
			any = true
			out.Printf("\n Found <headline>%d<reset> which are <subline>ahead of local<reset>\n", len(reposAheadOfRemote))
			out.Printf("  <debug>Eg remote has commits which are not merged into local<reset>\n")
			for _, repo := range reposAheadOfRemote {
				out.Printf("  <info>%s<reset> (%s): <important>%s<reset>\n", repo.Name, repo.Type, repo.Path)
			}
		}
		if len(reposBehindOfRemote) > 0 {
			any = true
			out.Printf("\n Found <headline>%d<reset> which are <subline>behind of local<reset>\n", len(reposBehindOfRemote))
			out.Printf("  <debug>Eg local has commits which are not merged with (at least one) remote<reset>\n")
			for _, repo := range reposBehindOfRemote {
				out.Printf("  <info>%s<reset> (%s): <important>%s<reset>\n", repo.Name, repo.Type, repo.Path)
			}
		}
		if !any {
			out.Printf(" <success>All is in sync!<reset>\n")
		}

		return nil
	}

	return addRepoFilterOptions(clif.NewCommand("check", "Check all registered repos", cb))
}

func init() {
	Commands = append(Commands, cmdCheck)
}
