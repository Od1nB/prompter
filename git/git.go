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
	Trash     = "🗑️"
)

type Git struct {
	Dirty    bool
	branch   string
	statuses []Porcelain
}

func New() (Git, error) {
	var g Git

	branch, err := branchCMD.Output()
	if err == nil && !strings.Contains("fatal", string(branch)) {
		g.branch = strings.TrimSpace(string(branch))
	}

	res, err := statusCMD.Output()
	if err != nil {
		return Git{}, err
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

func (g Git) String() string {
	var prompt string
	if g.branch != "" {
		prompt += color.Paint(color.BrightBlue, fmt.Sprintf("🔀%s ", g.branch))
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
