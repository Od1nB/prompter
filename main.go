package main

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/Od1nB/prompter/git"
	"github.com/Od1nB/prompter/path"
)

func main() {
	var prompt string
	path, err := path.New()
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
