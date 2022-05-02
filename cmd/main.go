package main

import (
	sg "SnakeGameGolang/internal/snakegame"
)

func main() {
	snakeGame := sg.SnakeGame{}
	snakeGame.Init(20, 50, false)
	snakeGame.Run()
}
