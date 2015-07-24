package watch

import "fmt"

// Watch is a repository of a certain kind
type (
	Watch interface {

		// Changes checks if there are local changes in the repo, which are not committed
		Changes() (bool, error)

		// Remotes returns list of remote (URLs) of repository
		Remotes() ([]string, error)

		// Synced checks if all local changes are synchronized with all remotes
		Synced() (SyncStateNum, error)

		// Synced checks if all local changes are synchronized with all remotes
		States() ([]*SyncState, error)

		// Type returns name of the kind of repository this watch represents
		Type() string

		// Updates checks if there are remote updates for the repo
		Updates() (bool, error)
	}

	// SyncState describes state of a single (remote) branch compared to local
	SyncState struct {
		Remote string
		Branch string
		State  SyncStateNum
		Error  error
	}
	SyncStateNum int
)

const (
	// failed sync state
	SYNC_STATE_FAIL SyncStateNum = 1 + iota

	// local repo is in sync with remote
	SYNC_STATE_SAME

	// local repo is behind remote (we need to pull)
	SYNC_STATE_BEHIND

	// local repo is ahead of remote (we need to push)
	SYNC_STATE_AHEAD

	// branch exists locally but not on remote
	SYNC_STATE_MISSING
)

// watches holds checkers/constructors of specific watch implementations
var watches = make([]func(path, name string) (Watch, error), 0)

// Factory tries to create a new watch from given path. If there is no repo
// found under path (or there is no implementation for the repo kind) an error
// is returned
func Factory(path, name string) (Watch, error) {
	for _, check := range watches {
		if watch, err := check(path, name); err != nil {
			return nil, err
		} else if watch != nil {
			return watch, nil
		}
	}
	return nil, fmt.Errorf("No implementation found to watch %s", path)
}
