package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/abhinav/git-fu/cli"
	"github.com/abhinav/git-fu/editor"
	"github.com/abhinav/git-fu/git"
	"github.com/abhinav/git-fu/pr"

	"github.com/jessevdk/go-flags"
)

type landCmd struct {
	Editor string `long:"editor" env:"EDITOR" default:"vi" value-name:"EDITOR" description:"Editor to use for interactively editing commit messages."`
	Args   struct {
		// TODO: Auto guess base
		Base string `positional-arg-name:"BASE" required:"yes" description:"Base branch against which this PR was made."`
		Head string `positional-arg-name:"HEAD" description:"Name of the branch at the top of the stack."`
	} `positional-args:"yes" required:"yes"`

	getConfig cli.ConfigBuilder
}

func newLandCommand(cbuild cli.ConfigBuilder) flags.Commander {
	return &landCmd{getConfig: cbuild}
}

func (l *landCmd) Execute(args []string) error {
	cfg, err := l.getConfig()
	if err != nil {
		return err
	}

	repo := cfg.Repo()
	base := l.Args.Base
	head := l.Args.Head
	if head == "" {
		out, err := git.Output("rev-parse", "--abbrev-ref", "HEAD")
		if err != nil {
			return fmt.Errorf("Could not determine current branch: %v", err)
		}
		head = strings.TrimSpace(out)
	}

	editor, err := editor.Pick(l.Editor)
	if err != nil {
		return fmt.Errorf("Could not determine editor: %v", err)
	}

	prs, _, err := cfg.GitHub().PullRequests.List(repo.Owner, repo.Name, nil)
	if err != nil {
		return err
	}

	if err := pr.SortByLandingOrder(base, head, prs); err != nil {
		return err
	}

	lander := &pr.MessageEditLander{
		Editor: editor,
		Lander: &pr.SquashLander{
			GitHubClient: cfg.GitHub(),
		},
	}

	for _, pull := range prs {
		url := pr.URL(pull)
		log.Println("Landing", url)
		if err := lander.Land(pull); err != nil {
			return fmt.Errorf("could not land %v: %v", url, err)
		}
	}

	return nil
}
