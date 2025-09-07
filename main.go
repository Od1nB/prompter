package main

import (
	"flag"
	"fmt"

	"github.com/Od1nB/prompter/color"
	"github.com/Od1nB/prompter/git"
	"github.com/Od1nB/prompter/hostname"
	"github.com/Od1nB/prompter/path"
)

var containerEmoji = "🐋"

var (
	opts          = []path.Option{}
	maxLen        = flag.Int("max", 40, "set the max amount of chars the first prompt line should be")
	showContainer = flag.Bool("showcontainer", false, "display the "+containerEmoji+" emoji at the start of prompt if set")
	showHostname  = flag.Bool("hostname", false, "display the hostname before the path")
	showPath      = flag.Bool("showpath", true, "show the 'pwd'")
	showGit       = flag.Bool("showgit", true, "show the current branch if in a git repo")
)

func init() {
	flag.Parse()
	if maxLen != nil {
		opts = append(opts, path.WithMaxLen(*maxLen))
	}
}

func main() {
	p := New()
	fmt.Println(p.out)
}

type prompt struct {
	out  string
	def  string
	errs []error
}

type printPrompt func() (string, error)

var (
	containerPrompt printPrompt = func() (string, error) {
		if !*showContainer {
			return "", nil
		}
		return containerEmoji, nil
	}
	hostnamePrompt printPrompt = func() (string, error) {
		if !*showHostname {
			return "", nil
		}
		hn, err := hostname.New()
		if err != nil {
			return "", fmt.Errorf("hostname err: %w", err)
		}
		return hn.String(), nil
	}
	pathPrompt printPrompt = func() (string, error) {
		if !*showPath {
			return "", nil
		}
		p, err := path.New(opts...)
		if err != nil {
			return "", fmt.Errorf("path err: %w", err)
		}

		if *showHostname {
			return color.Paint(color.Cyan, "@") + p.String(), nil
		}
		return p.String(), nil
	}
	gitPrompt printPrompt = func() (string, error) {
		if !*showGit || !git.InRepo() {
			return "", nil
		}
		g, err := git.New()
		if err != nil {
			return "", fmt.Errorf("git err: %w", err)
		}
		return g.String(), nil
	}
)

func New() prompt {
	p := prompt{
		def:  "⚡",
		errs: make([]error, 0),
	}

	for _, pr := range []printPrompt{containerPrompt, hostnamePrompt, pathPrompt, gitPrompt} {
		str, err := pr()
		if err != nil {
			p.errs = append(p.errs, err)
		}
		p.out += str
	}

	switch {
	case len(p.errs) != 0:
		p.out = fmt.Sprintf("got errors when creating prompt: %v\n"+p.def, p.errs)
	case len(p.out) == 0:
		p.out = p.def
	default:
		p.out += "\n" + p.def
	}

	return p
}
