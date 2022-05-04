package main

import (
	cls "SnakeGameGolang/internal/clearscreen"
	sg "SnakeGameGolang/internal/snakegame"
	"fmt"

	"github.com/eiannone/keyboard"
)

var (
	displayFunc sg.DisplayFunc = func(board [][]sg.Cell, score int) {
		cls.ClearScreen()
		fmt.Printf("\t<Score: %d>\n", score)
		for hight := range board {
			for widht := range board[hight] {
				switch board[hight][widht] {
				case sg.CellEmpty:
					fmt.Print("_")
				case sg.CellFood:
					fmt.Print("$")
				case sg.CellSnakeHead:
					fmt.Print("%")
				case sg.CellSnakeTail:
					fmt.Print("*")
				}
			}
			fmt.Println()
		}
	}

	keyHandlerFunc sg.KeyHandlerFunc = func(quit chan bool, turn chan sg.Direction) {
		{
			keysEvents, err := keyboard.GetKeys(10)
			if err != nil {
				panic(err)
			}
			defer func() {
				_ = keyboard.Close()
			}()

			for {
				event := <-keysEvents
				if event.Err != nil {
					panic(event.Err)
				}

				switch event.Key {
				case keyboard.KeyArrowUp:
					turn <- sg.DirectionUp
				case keyboard.KeyArrowRight:
					turn <- sg.DirectionRight
				case keyboard.KeyArrowDown:
					turn <- sg.DirectionDown
				case keyboard.KeyArrowLeft:
					turn <- sg.DirectionLeft

				case keyboard.KeyEsc:
					quit <- true
					return
				default:
				}
			}
		}
	}
)

func main() {
	snakeGame := sg.SnakeGame{}
	snakeGame.Init(15, 15, false, displayFunc, keyHandlerFunc)
	score := snakeGame.Run()

	// Game over
	cls.ClearScreen()
	fmt.Printf("<< Score: %d >>\n", score)
}
