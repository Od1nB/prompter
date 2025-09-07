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
