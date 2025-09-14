package term

import (
	"io"
	"sync"

	"github.com/gaginan/gotris"
)

const (
	LinesText = "Lines: %v"
	LevelText = "Level: %v"
)

// NewRender creates a new terminal-based renderer that implements gotris.Renderer.
func NewRender(w io.Writer) gotris.Renderer {
	return &render{
		term: newTerm(w),
	}
}

// render implements gotris.Renderer for terminal output.
type render struct {
	mu   sync.RWMutex
	term Term
}

var (
	boardTopLeft   = cursor{row: 1, col: 1}  // Top-left corner of the board
	linesTopLeft   = cursor{row: 1, col: 22} // Top-left corner of the stats
	levelTopLeft   = cursor{row: 2, col: 22} // Top-left corner of the stats
	previewTopLeft = cursor{row: 3, col: 22} // Top-left corner of the next piece preview
)

func (r *render) Update(state gotris.GameState) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if len(state.Board) == 0 {
		return
	}
	r.term.Reset()
	// Draw board
	// Combine current piece with the board
	grid := state.Board.Combine(state.Current.Grid, state.Current.Location)
	r.updateBoard(grid)
	if len(state.Next) > 0 && len(state.Next[0]) > 0 {
		r.updatePreview(state.Next[0])
	}
	// Draw stats
	r.term.MoveTo(linesTopLeft)
	r.term.Write(DarkGray.Sprintf(LinesText, state.Lines))
	r.term.MoveTo(levelTopLeft)
	r.term.Write(DarkGray.Sprintf(LevelText, state.Level))
	r.term.Flush()
	r.term.MoveTo(boardTopLeft)
}

func (r *render) updateBoard(grid gotris.Grid) {
	var cur = boardTopLeft
	r.term.MoveTo(cur)
	for _, line := range grid {
		for _, state := range line {
			r.term.Write(CellWithState(state))
			r.term.Write(SpaceCell)
		}
		r.term.NewLine()
	}
}

func (r *render) updatePreview(next gotris.Grid) {
	grid := gotris.NewGrid(5, 5) // 5x5 empty grid for I as max size
	center := gotris.Location{X: (5 - len(next[0])) / 2, Y: (5 - len(next)) / 2}
	shape := grid.Combine(next, center)
	var cur = previewTopLeft
	r.term.MoveTo(cur)
	for _, cols := range shape {
		for _, state := range cols {
			r.term.Write(CellWithState(state))
			r.term.Write(SpaceCell)
		}
		cur.row += 1
		r.term.MoveTo(cur)
	}
}

func (r *render) Clear() {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.term.Clear()
}
