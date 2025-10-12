// Package path defines functions for making a path string. The default is using the pwd command
package path

import (
	"os/exec"
	"slices"
	"strings"

	"github.com/Od1nB/prompter/color"
)

var pwdCMD = exec.Command("pwd")

type Path struct {
	Show           bool
	splits         []string
	Color          color.Color
	HostnamePrefix bool
}

type Option func(o *Path)

func WithShow(b bool) Option {
	return func(o *Path) {
		o.Show = b
	}
}

func New(withHostname bool, opts ...Option) (*Path, error) {
	path := &Path{
		Color:          color.BrightMagenta,
		Show:           false,
		HostnamePrefix: withHostname,
	}
	for _, f := range opts {
		f(path)
	}

	if !path.Show {
		return nil, nil
	}

	pwd, err := pwdCMD.Output()
	if err != nil {
		return nil, err
	}
	tmpSplits := slices.DeleteFunc(strings.Split(strings.TrimSpace(string(pwd)), "/"), emptyOrNone)
	path.splits = make([]string, 0, len(tmpSplits))

	switch {
	// on / slice is empty
	case len(tmpSplits) == 0:
		path.splits = append(path.splits, "/")
	// OSx /users/
	case len(tmpSplits) >= 2 && strings.ToLower(tmpSplits[0]) == "users":
		path.splits = append(path.splits, "~")
		if len(tmpSplits) > 2 {
			path.splits = append(path.splits, tmpSplits[2:]...)
		}
	case tmpSplits[0] == "home":
		path.splits = append(path.splits, "~")
		if len(tmpSplits) > 2 {
			path.splits = append(path.splits, tmpSplits[2:]...)
		}
	case len(tmpSplits) >= 2 && strings.ToLower(tmpSplits[0]) == "var" && strings.ToLower(tmpSplits[1]) == "home":
		path.splits = append(path.splits, "~")
		if len(tmpSplits) > 3 {
			path.splits = append(path.splits, tmpSplits[3:]...)
		}
	default:
		path.splits = tmpSplits
	}

	return path, nil
}

func (p *Path) Reduce() (int, bool) {
	var canReduce bool
	before := p.Len()
	p.splits, canReduce = p.removeFromBeginning()

	return before - p.Len(), canReduce
}

func (p Path) String() string {
	p.splits = slices.DeleteFunc(p.splits, emptyOrNone)
	var pathPrefix string
	if p.HostnamePrefix {
		pathPrefix += color.Paint(color.Cyan, "@")
	}
	switch {
	case len(p.splits) == 1 && p.splits[0] != "/":
		return pathPrefix + color.Paint(p.Color, p.splits[0]+"/")
	case len(p.splits) == 1:
		return pathPrefix + color.Paint(p.Color, p.splits[0])
	case len(p.splits) > 1:
		return pathPrefix + color.Paint(p.Color, strings.Join(p.splits, "/")+"/")
	default:
		return pathPrefix + color.Paint(p.Color, "/")
	}
}

func (p Path) Len() int {
	i := 0
	for _, v := range p.splits {
		i += len(v) + 1
	}
	return i
}

func (p *Path) removeFromBeginning() ([]string, bool) {
	cpArr := slices.DeleteFunc(p.splits, emptyOrNone)
	for ind, v := range cpArr {
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
