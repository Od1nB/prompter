package color

import "fmt"

type Color string

var (
	Reset               Color = "\033[0m"
	Green               Color = "\x1b[32m"
	Magenta             Color = "\033[35m"
	Red                 Color = "\033[31m"
	Yellow              Color = "\033[33m"
	Cyan                Color = "\033[36m"
	BrightBlue          Color = "\x1b[94m"
	BrightMagenta       Color = "\x1b[95m"
	BrightRed           Color = "\x1b[91m"
	BoldCyan            Color = "\033[1;36m"
	BoldHighIntenseCyan Color = "\033[1;96m"
)

func (c Color) String() string {
	return string(c)
}

func Paint(c Color, s string) string {
	return fmt.Sprintf("%s%s%s", c, s, Reset)
}
