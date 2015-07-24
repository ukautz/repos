package commands
import (
	"github.com/ukautz/repos/common"
	. "github.com/ukautz/repos/common/debug"
	"regexp"
	"fmt"
	"gopkg.in/ukautz/clif.v0"
)

func addRepoFilterOptions(c *clif.Command) *clif.Command {
	return c.NewOption("include", "i", "Regular expression for repos to include, eg '^foo'", "", false, true).
		NewOption("exclude", "e", "Regular expression for repos to exclude, eg 'bar$'", "", false, true)
}

func reduceWithRepoFilters(c *clif.Command, allRepos []*common.Info) ([]*common.Info, error) {

	// determine repo name pattern
	includeRx := regexp.MustCompile(`.`)
	if pattern := c.Option("include").String(); pattern != "" {
		var err error
		if includeRx, err = regexp.Compile(pattern); err != nil {
			return nil, fmt.Errorf("Failed to compile include pattern: %s", err)
		}
	}
	var excludeRx *regexp.Regexp
	if pattern := c.Option("exclude").String(); pattern != "" {
		var err error
		if excludeRx, err = regexp.Compile(pattern); err != nil {
			return nil, fmt.Errorf("Failed to compile exclude pattern: %s", err)
		}
	}

	repos := make([]*common.Info, 0)
	for _, repo := range allRepos {
		if !includeRx.MatchString(repo.Name) {
			Debug(DEBUG3, "Exclude repo %s since not matching include", repo.Name)
			continue
		} else if excludeRx != nil && excludeRx.MatchString(repo.Name) {
			Debug(DEBUG3, "Exclude repo %s since match exclude", repo.Name)
			continue
		}
		repos = append(repos, repo)
	}

	return repos, nil
}