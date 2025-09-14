package term

import (
	"fmt"
	"io"
	"strings"
)

type cursor struct {
	row int
	col int
}

const (
	MoveTo = "\033[%d;%dH" // Move cursor to (row, col) start with 1
)

// Term represents a terminal interface for rendering and control.
type Term interface {
	// Reset clears the internal buffer.
	Reset()
	// Write writes data to the buffer.
	Write(a ...any) (n int, err error)
	// Writef writes formatted data to the buffer.
	Writef(format string, a ...any) (n int, err error)
	// MoveTo moves the cursor to the specified position in the buffer.
	MoveTo(c cursor) (n int, err error)
	// NewLine writes a newline to the buffer.
	NewLine() (n int, err error)
	// ResetCursor moves the cursor to the top-left in the buffer.
	ResetCursor()
	// Flush writes the buffer to the terminal output.
	Flush() (n int, err error)
	// Clear clears the terminal screen and buffer.
	Clear()
}

func newTerm(w io.Writer) Term {
	return &term{
		writer: w,
		buffer: &strings.Builder{},
	}
}

type term struct {
	writer io.Writer
	buffer *strings.Builder
}

func (t *term) Reset() {
	t.buffer.Reset()
}

func (t *term) Write(a ...any) (n int, err error) {
	return fmt.Fprint(t.buffer, a...)
}

func (t *term) Writef(format string, a ...any) (n int, err error) {
	return fmt.Fprintf(t.buffer, format, a...)
}

func (t *term) MoveTo(c cursor) (n int, err error) {
	return fmt.Fprintf(t.buffer, MoveTo, c.row, c.col)
}

func (t *term) NewLine() (n int, err error) {
	return fmt.Fprint(t.buffer, "\n")
}

func (t *term) ResetCursor() {
	t.MoveTo(cursor{1, 1})
}

func (t *term) Flush() (n int, err error) {
	n, err = fmt.Fprint(t.writer, t.buffer.String())
	t.Reset()
	return
}

func (t *term) Clear() {
	fmt.Fprint(t.writer, "\033[2J") // Clear screen
	t.Reset()
}

// HideCursor hides the terminal cursor.
func HideCursor(w io.Writer) {
	fmt.Fprint(w, "\033[?25l") // Hide cursor
}

// ShowCursor shows the terminal cursor.
func ShowCursor(w io.Writer) {
	fmt.Fprint(w, "\033[?25h") // Show cursor
}

// ClearScreen clears the terminal screen and moves the cursor to the top-left.
func ClearScreen(w io.Writer) {
	fmt.Fprint(w, "\033[2J") // Clear screen
	fmt.Fprint(w, "\033[H")  // Move cursor to top-left
}
