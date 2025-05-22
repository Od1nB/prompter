package git

type Status rune

var (
	UnModified Status = ' '
	Modified   Status = 'M'
	FileType   Status = 'T'
	Added      Status = 'A'
	Deleted    Status = 'D'
	Renamed    Status = 'R'
	Copied     Status = 'C'
	Updated    Status = 'U'
	UnTracked  Status = '?'
	Illegal    Status
)

func ConvStatus(inp rune) Status {
	switch inp {
	case rune(UnModified):
		return UnModified
	case rune(Modified):
		return Modified
	case rune(FileType):
		return FileType
	case rune(Added):
		return Added
	case rune(Deleted):
		return Deleted
	case rune(Renamed):
		return Renamed
	case rune(Copied):
		return Copied
	case rune(Updated):
		return Updated
	case rune(UnTracked):
		return UnTracked
	default:
		return Illegal
	}
}

type Porcelain struct {
	X, Y     Status
	Meaning  string
	Filename string
}

func ConvPorcelain(inp string) Porcelain {
	return Porcelain{
		X:        ConvStatus(rune(inp[0])),
		Y:        ConvStatus(rune(inp[1])),
		Filename: inp[3:],
	}
}

func (p Porcelain) Staged() bool {
	for _, s := range []Status{p.X} {
		if s == Added || s == Modified || s == Deleted || s == Renamed || s == Copied {
			return true
		}
	}
	return false
}
