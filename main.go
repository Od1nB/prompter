package main

import (
	"flag"
	"fmt"
	"log/slog"
	"os"

	"github.com/Od1nB/prompter/git"
	"github.com/Od1nB/prompter/path"
)

var opts = []path.Option{}

func init() {
	maxLen := flag.Int("max", 40, "set the max amount of chars the first prompt line should be")
	flag.Parse()

	if maxLen != nil {
		opts = append(opts, path.WithMaxLen(*maxLen))
	}
}

func main() {
	var prompt string
	path, err := path.New(opts...)
	if err != nil {
		fmt.Print("⚡")
	}
	prompt += path.String()
	if git.InRepo() {
		g, err := git.New()
		if err != nil {
			slog.Error("err parsing git ", "err", err)
			os.Exit(2)
		}
		prompt += " " + g.String()
	}
	prompt += "\n⚡"

	fmt.Print(prompt)
}
