
# gotris-term

English | [简体中文](./README.zh-CN.md)

A demonstration project showing how to implement a custom renderer for gotris. gotris-term uses text-based rendering to build a terminal mini-game based on gotris, illustrating how to extend gotris with your own visual output.

# How to Build and Play

## Build the Binary

To build the gotris-term binary, run the following command in the project root:

```sh
go build -o gotris-term
```

This will create an executable named `gotris-term` in the current directory.

## Play the Game

To start the game, run:

```sh
./gotris-term
```

Controls:

- Arrow keys: Move and rotate pieces
- Space: Hard drop
- Esc or Ctrl+C: Quit the game
