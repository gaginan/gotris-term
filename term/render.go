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
	boardTopLeft = cursor{row: 1, col: 1} // Top-left corner of the board
)

const (
	cellRenderWidth = 2
	statsColGap     = 1

	linesRow   = 1
	levelRow   = 2
	previewRow = 3

	defaultStatsCol = 22
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
	boardRightTop := boardTopRight(state.Board)
	statsCol := boardRightTop.col + statsColGap
	previewTopLeft := cursor{row: previewRow, col: statsCol}
	if len(state.Next) > 0 && len(state.Next[0]) > 0 {
		r.updatePreview(state.Next[0], previewTopLeft)
	}
	// Draw stats
	linesTopLeft := cursor{row: linesRow, col: statsCol}
	levelTopLeft := cursor{row: levelRow, col: statsCol}

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

func (r *render) updatePreview(next gotris.Grid, at cursor) {
	grid := gotris.NewGrid(5, 5) // 5x5 empty grid for I as max size
	center := gotris.Location{X: (5 - len(next[0])) / 2, Y: (5 - len(next)) / 2}
	shape := grid.Combine(next, center)
	var cur = at
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

func boardTopRight(board gotris.Grid) cursor {
	if len(board) == 0 || len(board[0]) == 0 {
		return cursor{row: boardTopLeft.row, col: defaultStatsCol - statsColGap}
	}
	_, cols := board.Size()
	return cursor{row: boardTopLeft.row, col: boardTopLeft.col + cols*cellRenderWidth}
}

func (r *render) Clear() {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.term.Clear()
}
