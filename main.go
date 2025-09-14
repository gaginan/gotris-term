package main

import (
	"context"
	"fmt"
	"os"

	"github.com/eiannone/keyboard"
	"github.com/gaginan/gotris"
	"github.com/gaginan/gotris-term/term"
)

// main is the entry point for the gotris terminal game.
func main() {
	//panic recover
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("Recovered from panic: %v\n", r)
		}
	}()
	term.HideCursor(os.Stdout)
	defer term.ShowCursor(os.Stdout)
	term.ClearScreen(os.Stdout)
	controls := make(chan gotris.Control)
	ctx, cancel := context.WithCancel(context.Background())
	var game = gotris.New(ctx, term.NewRender(os.Stdout), controls)
	go accept(cancel, controls)
	game.Run(ctx)
}

// accept listens for keyboard input and sends game controls.
func accept(cancel context.CancelFunc, controls chan<- gotris.Control) {
	if err := keyboard.Open(); err != nil {
		fmt.Printf("keyboard init error: %v\n", err)
		return
	}
	defer keyboard.Close()
	for {
		_, k, err := keyboard.GetKey()
		if err != nil {
			continue
		}
		switch k {
		case keyboard.KeyEsc, keyboard.KeyCtrlC:
			cancel()
			return
		case keyboard.KeySpace:
			controls <- gotris.HardDrop()
		case keyboard.KeyArrowUp:
			controls <- gotris.Rotate(gotris.RotateRight)
		case keyboard.KeyArrowDown:
			controls <- gotris.Move(gotris.Down)
		case keyboard.KeyArrowLeft:
			controls <- gotris.Move(gotris.Left)
		case keyboard.KeyArrowRight:
			controls <- gotris.Move(gotris.Right)
		}
	}
}
