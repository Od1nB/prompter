package main

import (
	"flag"
	"fmt"
	"log/slog"
	"os"

	"github.com/Od1nB/prompter/git"
	"github.com/Od1nB/prompter/hostname"
	"github.com/Od1nB/prompter/path"
)

var containerEmoji = "🐋"

var (
	opts            = []path.Option{}
	maxLen          = flag.Int("max", 40, "set the max amount of chars the first prompt line should be")
	showContainer   = flag.Bool("showcontainer", false, "display the "+containerEmoji+" emoji at the start of prompt if set")
	displayHostname = flag.Bool("hostname", false, "display the hostname before the path")
	showPath        = flag.Bool("showpath", true, "show the 'pwd'")
	showGit         = flag.Bool("showgit", true, "show the current branch if in a git repo")
)

func init() {
	flag.Parse()
	if maxLen != nil {
		opts = append(opts, path.WithMaxLen(*maxLen))
	}
}

func main() {
	var prompt string
	if *showContainer {
		prompt += containerEmoji
	}

	if *displayHostname {
		hn, err := hostname.New()
		if err != nil {
			fmt.Println("⚡")
			os.Exit(2)
		}
		prompt += hn.String()
	}

	if *showPath {
		path, err := path.New(opts...)
		if err != nil {
			fmt.Print("⚡")
			os.Exit(2)
		}
		prompt += path.String()
	}

	if *showGit {
		if git.InRepo() {
			g, err := git.New()
			if err != nil {
				slog.Error("err parsing git ", "err", err)
				os.Exit(2)
			}
			prompt += " " + g.String()
		}
	}
	prompt += "\n⚡"

	fmt.Print(prompt)
}
