// Package git has functions for creating a string about the git context from the current working dir
package git

import (
	"bufio"
	"bytes"
	"fmt"
	"os/exec"
	"strconv"
	"strings"

	"github.com/Od1nB/prompter/color"
)

var (
	isRepoCMD = exec.Command("git", "rev-parse", "--is-inside-work-tree")
	statusCMD = exec.Command("git", "status", "--porcelain=v1")
	branchCMD = exec.Command("git", "branch", "--show-current")
	tagCMD    = exec.Command("git", "describe", "--tags", "--exact-match", "HEAD")
	commitCMD = exec.Command("git", "rev-parse", "--short", "HEAD")
	Trash     = "🗑️"
)

type Git struct {
	Dirty    bool
	location string
	statuses []Porcelain
}

func New(show bool) (*Git, error) {
	if !show {
		return nil, nil
	}
	g := new(Git)
	g.location = location()
	res, err := statusCMD.Output()
	if err != nil {
		return nil, err
	}

	g.statuses = parseLines(res)
	if len(g.statuses) > 0 {
		g.Dirty = true
	}

	return g, nil
}

func InRepo() bool {
	resp, err := isRepoCMD.Output()
	if err != nil {
		return false
	}
	b, _ := strconv.ParseBool(strings.TrimSpace(string(resp)))
	return b
}

func (g Git) Len() int {
	return len(g.String())
}

func (g Git) Reduce() (int, bool) {
	return 0, false
}

func (g Git) String() string {
	var prompt string
	if g.location != "" {
		prompt += g.location
	}
	if g.Dirty {
		prompt += Trash
		c := statusColor(g.statuses)
		prompt += color.Paint(c, "±"+fmt.Sprintf("%d", len(g.statuses)))
	} else {
		prompt += "✨"
	}

	return prompt
}

func statusColor(ss []Porcelain) color.Color {
	var c color.Color
	if len(ss) >= 1 {
		c = color.Yellow
	}

	var added int
	for _, s := range ss {
		if s.X == UnTracked || s.Y == UnTracked {
			return color.Red
		}
		if s.Staged() {
			added++
		}
	}
	if len(ss) == added {
		c = color.Green
	}
	return c
}

func parseLines(b []byte) []Porcelain {
	var res []Porcelain
	reader := bufio.NewScanner(bytes.NewReader(b))

	for reader.Scan() {
		l := reader.Text()
		if len(l) == 0 || len(l) < 4 {
			continue
		}
		res = append(res, ConvPorcelain(l))
	}
	return res
}

func location() string {
	if res, err := tagCMD.Output(); err == nil &&
		len(res) != 0 &&
		!strings.Contains("fatal:", string(res)) {
		return color.Paint(color.Yellow,
			strings.TrimSpace("🏷️"+string(res)))
	}

	if res, err := branchCMD.Output(); err == nil &&
		len(res) != 0 &&
		!strings.Contains("fatal:", string(res)) {
		return color.Paint(color.BrightBlue,
			strings.TrimSpace("🔀"+string(res)))
	}

	if res, err := commitCMD.Output(); err == nil &&
		len(res) != 0 &&
		!strings.Contains("fatal:", string(res)) {
		return color.Paint(color.BrightRed,
			strings.TrimSpace("⁉️"+string(res)))
	}

	return ""
}
