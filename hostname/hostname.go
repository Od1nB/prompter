// Package hostname can create a string based on the hostname where the functions are ran
package hostname

import (
	"os/exec"
	"strings"

	"github.com/Od1nB/prompter/color"
)

type Host struct {
	color color.Color
	Str   string
}

var hostNameCMD = exec.Command("hostname")

func New() (Host, error) {
	res, err := hostNameCMD.Output()
	if err != nil {
		return Host{}, err
	}

	return Host{
		color: color.Cyan,
		Str:   strings.TrimSpace(string(res)),
	}, nil
}

func (h Host) String() string {
	return color.Paint(h.color, h.Str)
}

func (h Host) Len() int {
	return len(h.Str)
}

func (h *Host) Reduce() (int, bool) {
	if h.Len() <= 1 {
		return 0, false
	}

	prev := h.Len()
	h.Str = h.Str[:len(h.Str)/2]
	return prev - h.Len(), true
}
