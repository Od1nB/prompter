package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/Od1nB/prompter/container"
	"github.com/Od1nB/prompter/git"
	"github.com/Od1nB/prompter/hostname"
	"github.com/Od1nB/prompter/path"

	"golang.org/x/term"
)

const containerEmoji = "🐋"

var (
	showContainer = flag.Bool("showcontainer", false, "display the "+containerEmoji+" emoji at the start of prompt if set")
	showHostname  = flag.Bool("hostname", false, "display the hostname before the path")
	showPath      = flag.Bool("showpath", true, "show the 'pwd'")
	showGit       = flag.Bool("showgit", true, "show the current branch if in a git repo")
)

func init() {
	flag.Parse()
}

func main() {
	p := New()
	fmt.Println(p.out)
}

type promptParam interface {
	fmt.Stringer
	Reduce() (int, bool)
}

type prompt struct {
	width  int
	lenght int
	outs   []promptParam
	out    string
	def    string
	errs   []error
}

func New() prompt {
	p := prompt{
		def:  "⚡",
		outs: make([]promptParam, 0, 4),
		errs: make([]error, 0),
	}

	w, _, err := term.GetSize(int(os.Stdin.Fd()))
	if err != nil {
		p.errs = append(p.errs, err)
	} else {
		p.width = w
	}

	c := container.New(containerEmoji, *showContainer)
	if c != nil {
		p.outs = append(p.outs, c)
		p.lenght += c.Len()
	}

	hostnameInd := -1
	hostNamePrompt, err := hostname.New(*showHostname)
	if err != nil {
		p.errs = append(p.errs, fmt.Errorf("hostname err: %w", err))
	}
	if hostNamePrompt != nil {
		p.outs = append(p.outs, hostNamePrompt)
		hostnameInd = len(p.outs) - 1
		p.lenght += hostNamePrompt.Len()
	}

	pathInd := -1
	pathPrompt, err := path.New(*showHostname, path.WithShow(*showPath))
	if err != nil {
		p.errs = append(p.errs, fmt.Errorf("path err: %w", err))
	}
	if pathPrompt != nil {
		p.outs = append(p.outs, pathPrompt)
		pathInd = len(p.outs) - 1
		p.lenght += pathPrompt.Len()
	}

	if git.InRepo() {
		g, err := git.New(*showGit)
		if err != nil {
			p.errs = append(p.errs, fmt.Errorf("git err: %w", err))
		}
		if g != nil {
			p.outs = append(p.outs, g)
			p.lenght += g.Len()
		}
	}

	switch {
	case len(p.errs) != 0:
		p.out = fmt.Sprintf("got errors when creating prompt: %v\n"+p.def, p.errs)
	case len(p.outs) != 0:
		p.populateOut(hostnameInd, pathInd)
		fallthrough
	default:
		p.out += "\n" + p.def
	}

	return p
}

func (p *prompt) reduceIfTooWide(host, path int) {
	if p.lenght <= p.width {
		return
	}

	reduceHost := func() (int, bool) {
		if host == -1 {
			return 0, false
		}
		return p.outs[host].Reduce()
	}

	reducePath := func() (int, bool) {
		if path == -1 {
			return 0, false
		}
		return p.outs[path].Reduce()
	}

	canChange := true
	for p.lenght > p.width && canChange {
		h, canHost := reduceHost()
		p.lenght = p.lenght - h
		if p.lenght <= p.width {
			return
		}

		pi, canPath := reducePath()
		p.lenght = p.lenght - pi
		if p.lenght <= p.width {
			return
		}
		canChange = canHost || canPath
	}
}

func (p *prompt) populateOut(host, path int) {
	p.reduceIfTooWide(host, path)
	for _, pr := range p.outs {
		p.out += pr.String()
	}
}
