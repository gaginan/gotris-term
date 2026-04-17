package term

import (
	"io"
	"sync"
	"unicode/utf8"

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
	mu      sync.RWMutex
	term    Term
	corners corners
}
type corners struct {
	topLeft     cursor
	topRight    cursor
	bottomLeft  cursor
	bottomRight cursor
}

var (
	boardTopLeft = cursor{row: 1, col: 1} // Top-left corner of the board
)

const (
	statsLinesRow = 1
	statsLevelRow = 2
	previewTopRow = 3
)

var (
	rightPanelGap   = utf8.RuneCountInString(SpaceCell)
	boardCellWidth = utf8.RuneCountInString(DefaultCell) + utf8.RuneCountInString(SpaceCell)
)

func (r *render) Update(state gotris.GameState) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.term.Reset()
	// Update corners based on current board size and position
	r.updateCorners(state)
	// Update the board, preview, and stats
	r.updateBoard(state)
	r.updatePreview(state)
	r.updateStats(state)
	r.term.Flush()
	r.term.MoveTo(r.corners.topLeft)
}

func (r *render) updateBoard(state gotris.GameState) {
	var grid = state.Board.Combine(state.Current.Grid, state.Current.Location)
	r.term.MoveTo(r.corners.topLeft)
	for _, line := range grid {
		for _, state := range line {
			r.term.Write(CellWithState(state))
			r.term.Write(SpaceCell)
		}
		r.term.NewLine()
	}
}

func (r *render) updatePreview(state gotris.GameState) {
	if len(state.Next) == 0 || len(state.Next[0]) == 0 {
		return
	}
	next := state.Next[0]
	grid := gotris.NewGrid(5, 5) // 5x5 empty grid for I as max size
	center := gotris.Location{X: (5 - len(next[0])) / 2, Y: (5 - len(next)) / 2}
	shape := grid.Combine(next, center)
	var cur = cursor{row: previewTopRow, col: r.corners.topRight.col + rightPanelGap}
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

func (r *render) updateStats(state gotris.GameState) {
	statsCol := r.corners.topRight.col + rightPanelGap
	linesTopLeft := cursor{row: statsLinesRow, col: statsCol}
	levelTopLeft := cursor{row: statsLevelRow, col: statsCol}

	r.term.MoveTo(linesTopLeft)
	r.term.Write(DarkGray.Sprintf(LinesText, state.Lines))
	r.term.MoveTo(levelTopLeft)
	r.term.Write(DarkGray.Sprintf(LevelText, state.Level))
}

func (r *render) updateCorners(state gotris.GameState) {
	if len(state.Board) == 0 || len(state.Board[0]) == 0 {
		r.corners = corners{}
		return
	}
	rows, cols := state.Board.Size()
	r.corners = corners{
		topLeft:     boardTopLeft,
		topRight:    cursor{row: boardTopLeft.row, col: boardTopLeft.col + cols*boardCellWidth},
		bottomLeft:  cursor{row: boardTopLeft.row + rows, col: boardTopLeft.col},
		bottomRight: cursor{row: boardTopLeft.row + rows, col: boardTopLeft.col + cols*boardCellWidth},
	}
}

func (r *render) Clear() {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.term.Clear()
	r.corners = corners{}
}
