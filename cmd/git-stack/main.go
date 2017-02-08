package main

import "github.com/abhinav/git-fu/cli"

func main() {
	cli.Main(
		cli.ShortDesc("Make stacked GitHub PRs easier."),
		&cli.Command{
			Name:      "land",
			ShortDesc: "Lands a stack of GitHub PRs.",
			Build:     newLandCommand,
		},
	)
}
