package path

import (
	"os/exec"
	"strings"

	"github.com/Od1nB/prompter/color"
)

var (
	pwdCMD = exec.Command("pwd")
)

type Path struct {
	Color color.Color
	Str   string
}

func New() (Path, error) {
	r, err := pwdCMD.Output()
	if err != nil {
		return Path{}, err
	}
	var path = Path{Color: color.BrightMagenta, Str: string(r)}
	splits := strings.Split(strings.TrimSpace(string(r)), "/")
	if len(splits) >= 3 && strings.ToLower(splits[1]) == "users" {
		path.Str = "~"
		if len(splits) > 3 {
			path.Str += "/" + strings.Join(splits[3:], "/")
		}
	}
	return path, nil
}

func (p Path) String() string {
	return color.Paint(p.Color, p.Str)
}
