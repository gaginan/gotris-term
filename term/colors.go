package term

import (
	"fmt"

	"github.com/gaginan/gotris"
)

// Color represents an ANSI color code for terminal output.
type Color string

// ANSI colors per Tetris guideline colors
const (
	DarkGray Color = "\033[38;2;76;72;69m"
	Red      Color = "\033[91m"       // Z
	Orange   Color = "\033[38;5;208m" // L
	Yellow   Color = "\033[93m"       // O
	Green    Color = "\033[92m"       // S
	Cyan     Color = "\033[96m"       // I
	Blue     Color = "\033[94m"       // J
	Purple   Color = "\033[95m"       // T
	Gray     Color = "\033[90m"
	White    Color = "\033[97m"
	Reset    Color = "\033[0m"
)

// Sprintf returns a formatted string with the color applied.
func (c Color) Sprintf(format string, a ...any) string {
	return fmt.Sprintf(string(c)+format+string(Reset), a...)
}

const (
	SpaceCell   = " "
	ColorCell   = "%s■\033[0m"
	DefaultCell = "■"
)

var stateColors = map[gotris.State]Color{
	gotris.Empty:  DarkGray,
	gotris.Red:    Red,
	gotris.Orange: Orange,
	gotris.Yellow: Yellow,
	gotris.Green:  Green,
	gotris.Cyan:   Cyan,
	gotris.Blue:   Blue,
	gotris.Purple: Purple,
	gotris.Gray:   Gray,
}

// CellWithState returns a colored cell string for the given game state.
func CellWithState(s gotris.State) string {
	if color, ok := stateColors[s]; ok {
		return color.Sprintf(ColorCell, color)
	}
	return DefaultCell
}
