package commands

import (
	"fmt"
	"github.com/ukautz/repos/common"
	"gopkg.in/ukautz/clif.v1"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
)

func scan(dir string, dirs chan<- string, out clif.Output, maxDepth, depth int) error {
	//defer out.Printf("-1 %s\n", dir)
	//out.Printf("+1 %s\n", dir)
	dirs <- dir
	var wg sync.WaitGroup
	if fhs, err := ioutil.ReadDir(dir); err != nil {
		return err
	} else {
		for _, f := range fhs {
			if f.IsDir() {
				abs := filepath.Join(dir, f.Name())
				if depth+1 < maxDepth {
					wg.Add(1)
					go func() {
						defer wg.Done()
						scan(abs, dirs, out, maxDepth, depth+1)
					}()
				}
			}
		}
		wg.Wait()
		return nil
	}
}

func cmdScan() *clif.Command {
	cb := func(c *clif.Command, in clif.Input, out clif.Output, lst *common.List) error {
		include := regexp.MustCompile(c.Option("include").String())
		var exclude *regexp.Regexp
		if excludeStr := c.Option("exclude").String(); excludeStr != "" {
			exclude = regexp.MustCompile(excludeStr)
		}
		added := 0
		exit := false
		mux := new(sync.Mutex)
		c.Cli.SetOnInterrupt(func() error {
			exit = true
			out.Printf("\n")
			if added > 0 {
				mux.Unlock()
				mux.Lock()
				out.Printf("\n<important>Aborting<reset>. Just got one question:\n")
				if in.Confirm(fmt.Sprintf("You've added %d repos. Shall I persist?", added)) {
					lst.Persist()
				}
			}
			os.Exit(0)
			return nil
		})

		prefix := c.Option("prefix").String()
		dir := c.Argument("directory").String()
		maxDepth := c.Option("max-depth").Int() + 1
		root, err := filepath.Abs(dir)
		if err != nil {
			return err
		}
		out.Printf("Scanning <headline>%s<reset> for repositories\n", root)
		dirs := make(chan string)
		var wg sync.WaitGroup
		wg.Add(1)
		go func() {
			defer wg.Done()
			defer close(dirs)
			scan(root, dirs, out, maxDepth, 0)
		}()

		for dir := range dirs {
			if exit {
				break
			}
			if addRepo(lst, dir, prefix, include, exclude, out, in, mux) {
				added++
			}
		}

		return lst.Persist()
	}

	return clif.NewCommand("scan", "Scan a directory for repositories", cb).
		SetDescription(strings.Join([]string{
		"Scan a directory (recursively) for repostories. For each found repository a name,",
		"derived from the folder name, will be suggested automatically.",
		"",
	}, "\n")).
		NewArgument("directory", "Directory to scan", ".", true, false).
		NewOption("max-depth", "d", "Max depth to scan (1 = all sub folders, 2 = also subfolders within these, ..) ", "1", false, false).
		NewOption("include", "i", "Regular expression which names must match to be considered", ".", false, false).
		NewOption("exclude", "e", "Regular expression which names must not match to be considered", "", false, false).
		NewOption("prefix", "p", "Prefix added to all suggested names", "", false, false)
}

func addRepo(lst *common.List, dir, prefix string, include, exclude *regexp.Regexp, out clif.Output, in clif.Input, mux *sync.Mutex) bool {
	mux.Lock()
	defer mux.Unlock()
	if lst.Watched(dir) == "" {
		out.Printf("Considering directory <info>%s<reset>\n", dir)
		base := strings.ToLower(filepath.Base(dir))
		name := prefix + base
		if _, err := common.NewRepo(dir, name); err != nil {
			out.Printf("  <warn>%s<reset>\n\n", err)
			return false
		}
		if !include.MatchString(name) {
			out.Printf("  <debug>Does not match include filter: \"%s\"<reset>\n\n", name)
			return false
		} else if exclude != nil && exclude.MatchString(name) {
			out.Printf("  <debug>Does match exclude filter: \"%s\"<reset>\n\n", name)
			return false
		}

		// find name
		cnt := 1
		for lst.Get(name) != "" {
			name = base + fmt.Sprintf("%d", cnt)
			cnt++
		}
		checkRx := regexp.MustCompile(`^(y(?:es)?|no?|r(?:ename)?)$`)
		add := in.AskRegex(fmt.Sprintf("  Add as \"%s\"? <info>(y=yes, n=no, r=rename)", name), checkRx)
		if add == "n" || add == "no" {
			out.Printf("  <debug>Not added<reset>\n\n")
			return false
		} else if add == "r" || add == "rename" {
			for {
				if as := in.Ask(fmt.Sprintf("Name (%s):", name), clif.InputEmptyOk); as != "" {
					if p := lst.Get(as); p != "" {
						out.Printf("  <warn>Name \"%s\" already in use for directory \"%s\".\n  Choose something else!<reset>\n\n", as, p)
					} else {
						name = as
						break
					}
				}
			}
		}
		if _, err := lst.Add(name, dir); err != nil {
			out.Printf("  <error>Failed to add: %s<reset>\n\n", err)
			return false
		} else {
			out.Printf("  <success>Added directory <info>%s<success> as <info>%s<reset>\n\n", dir, name)
			return true
		}
	} else {
		out.Printf("<debug>%s is already watched\n<reset>\n\n", dir)
		return false
	}
}

func init() {
	Commands = append(Commands, cmdScan)
}
