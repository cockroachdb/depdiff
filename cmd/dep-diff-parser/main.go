package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/BurntSushi/toml"
	"github.com/dt/glide-diff-parser/lockfile"
)

const lockfileName = "Gopkg.lock"

func parseDepTOML(content []byte) (lockfile.Versions, error) {
	l := struct {
		Projects []struct {
			Name     string
			Revision string
		}
	}{}
	err := toml.Unmarshal(content, &l)
	if err != nil {
		return nil, err
	}
	m := make(map[string]string, len(l.Projects))
	for _, dep := range l.Projects {
		m[dep.Name] = dep.Revision
	}
	return m, nil
}

func main() {
	verbose := flag.Bool("v", false, "print a detailed summary of added, removed and changed dependencies")
	flag.Usage = lockfile.Usage(lockfileName)
	flag.Parse()

	args := flag.Args()
	if len(args) > 2 {
		flag.Usage()
		os.Exit(-1)
	}
	if err := lockfile.SummarizeDiff(args, *verbose, lockfileName, parseDepTOML); err != nil {
		fmt.Fprintf(os.Stderr, "%+v", err)
		os.Exit(1)
	}
}
