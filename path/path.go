package path

import (
	"fmt"
	"os/exec"
	"slices"
	"strings"

	"github.com/Od1nB/prompter/color"
)

var pwdCMD = exec.Command("pwd")

type Path struct {
	Options options
	splits  []string
	Color   color.Color
	Str     string
}

type Option func(o *options)

type options struct {
	maxLen int
}

func WithMaxLen(i int) Option {
	return func(o *options) {
		o.maxLen = i
	}
}

func New(opts ...Option) (Path, error) {
	pwd, err := pwdCMD.Output()
	if err != nil {
		return Path{}, err
	}

	o := &options{
		maxLen: 40,
	}

	for _, f := range opts {
		f(o)
	}
	tmpSplits := slices.DeleteFunc(strings.Split(strings.TrimSpace(string(pwd)), "/"), emptyOrNone)
	path := Path{
		Color:   color.BrightMagenta,
		splits:  make([]string, 0, len(tmpSplits)),
		Options: *o,
		// Str: string(r)
	}

	switch {
	// OSx /users/
	case len(tmpSplits) >= 3 && strings.ToLower(tmpSplits[1]) == "users":
		path.splits = append(path.splits, "~")
		if len(tmpSplits) > 3 {
			path.splits = append(path.splits, tmpSplits[3:]...)
		}
	case tmpSplits[0] == "home":
		path.splits = append(path.splits, "~")
		if len(tmpSplits) > 2 {
			path.splits = append(path.splits, tmpSplits[2:]...)
		}
	}

	change := true
	for path.len() > path.Options.maxLen && len(path.splits) > 1 && change {
		path.splits, change = path.removeFromBeginning()
	}

	return path, nil
}

func (p Path) String() string {
	p.splits = slices.DeleteFunc(p.splits, emptyOrNone)
	str := p.splits[0] + "/"
	if len(p.splits) > 1 {
		str = strings.Join(p.splits, "/") + "/"
	}
	return color.Paint(p.Color, str)
}

func (p Path) len() int {
	i := 0
	for _, v := range p.splits {
		i += len(v) + 1
	}
	return i
}

func (p *Path) removeFromBeginning() ([]string, bool) {
	cpArr := slices.DeleteFunc(p.splits, emptyOrNone)
	for ind, v := range cpArr {
		fmt.Println(v, len(cpArr), ind+1)
		switch {
		case v == "..":
			inds, ok := findTwoOccurrences(cpArr, "..")
			if ok {
				cpArr[inds[0]] = "..."
				return slices.Delete(cpArr, inds[1], inds[1]+1), true
			}
			fallthrough
		case v == "~" || v == ".." || ind+1 == len(cpArr):
			continue
		case v == "..." && ind+1 == len(cpArr)-1:
			return cpArr, false
		case v == "..." && ind+1 != len(cpArr)-1:
			return slices.Delete(cpArr, ind+1, ind+2), true
		default:
			cpArr[ind] = ".."
			inds, ok := findTwoOccurrences(cpArr, "..")
			if ok {
				cpArr[inds[0]] = "..."
				return slices.Delete(cpArr, inds[1], inds[1]+1), true
			}
			return cpArr, true
		}
	}
	return cpArr, false
}

var emptyOrNone = func(s string) bool {
	return len(s) == 0 ||
		s == "" ||
		strings.TrimSpace(s) == "" ||
		len(strings.TrimSpace(s)) == 0
}

func findTwoOccurrences(s []string, target string) ([]int, bool) {
	indices := make([]int, 0, 2)
	for ind, elem := range s {
		if elem == target {
			indices = append(indices, ind)
			if len(indices) > 2 {
				return indices, true
			}
		}
	}

	if len(indices) == 2 {
		return indices, true
	}

	return nil, false
}
