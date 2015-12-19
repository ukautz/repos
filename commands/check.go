package commands

import (
	"github.com/ukautz/repos/common"
	. "github.com/ukautz/repos/common/debug"
	"gopkg.in/ukautz/clif.v1"
	"sync"
	"fmt"
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
		var pbs clif.ProgressBarPool

		wg.Add(1)
		if DebugLevel == DEBUG0 {
			pbs = out.ProgressBars()
			style := clif.CloneProgressBarStyle(clif.ProgressBarStyleUtf8)
			style.Count = clif.PROGRESS_BAR_ADDON_PREPEND
			style.Elapsed = clif.PROGRESS_BAR_ADDON_PREPEND
			style.Estimate = clif.PROGRESS_BAR_ADDON_APPEND
			style.Percentage = clif.PROGRESS_BAR_ADDON_OFF
			pbs.Style(style)
			pbs.Start()
			pb, _ := pbs.Init("repos", len(repos))
			go func() {
				defer wg.Done()
				for range progress {
					pb.Increment()
				}
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
		if pbs != nil {
			<-pbs.Finish()
		}

		any := false
		if len(reposWithError) > 0 {
			any = true
			out.Printf("\n- - -\n\n Found <headline>%d<reset> with <subline>errors<reset>\n\n", len(reposWithError))
			table := out.Table([]string{"Name", "Path", "Error"})
			for _, repo := range reposWithError {
				table.AddRow([]string{repo.Name, repo.Path, repo.Error.Error()})
				//out.Printf("  <info>%s<reset>: <error>%s<reset>\n", repo.Name, repo.Error)
			}
			fmt.Println(table.Render())
		}
		if len(reposWithLocalChanges) > 0 {
			any = true
			out.Printf("\n- - -\n\n Found <headline>%d<reset> with <subline>local changes<reset>\n", len(reposWithLocalChanges))
			out.Printf("  <debug>Eg uncommited changes<reset>\n\n")
			table := out.Table([]string{"Name", "Type", "Path"})
			for _, repo := range reposWithLocalChanges {
				table.AddRow([]string{repo.Name, repo.Type, repo.Path})
				//out.Printf("  <info>%s<reset> (%s): <important>%s<reset>\n", repo.Name, repo.Type, repo.Path)
			}
			fmt.Println(table.Render())
		}
		if len(reposAheadOfRemote) > 0 {
			any = true
			out.Printf("\n- - -\n\n Found <headline>%d<reset> which are <subline>ahead of local<reset>\n", len(reposAheadOfRemote))
			out.Printf("  <debug>Eg remote has commits which are not merged into local<reset>\n\n")
			table := out.Table([]string{"Name", "Type", "Path"})
			for _, repo := range reposAheadOfRemote {
				table.AddRow([]string{repo.Name, repo.Type, repo.Path})
				//out.Printf("  <info>%s<reset> (%s): <important>%s<reset>\n", repo.Name, repo.Type, repo.Path)
			}
			fmt.Println(table.Render())
		}
		if len(reposBehindOfRemote) > 0 {
			any = true
			out.Printf("\n- - -\n\n Found <headline>%d<reset> which are <subline>behind of local<reset>\n", len(reposBehindOfRemote))
			out.Printf("  <debug>Eg local has commits which are not merged with (at least one) remote<reset>\n\n")
			table := out.Table([]string{"Name", "Type", "Path"})
			for _, repo := range reposBehindOfRemote {
				table.AddRow([]string{repo.Name, repo.Type, repo.Path})
				//out.Printf("  <info>%s<reset> (%s): <important>%s<reset>\n", repo.Name, repo.Type, repo.Path)
			}
			fmt.Println(table.Render())
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
