package color

import "fmt"

type Color string

var (
	Reset         Color = "\033[0m"
	Green         Color = "\x1b[32m"
	Magenta       Color = "\033[35m"
	Red           Color = "\033[31m"
	Yellow        Color = "\033[33m"
	BrightBlue    Color = "\x1b[94m"
	BrightMagenta Color = "\x1b[95m"
	BrightRed     Color = "\x1b[91m"
)

func (c Color) String() string {
	return string(c)
}

func Paint(c Color, s string) string {
	return fmt.Sprintf("%s%s%s", c, s, Reset)
}
