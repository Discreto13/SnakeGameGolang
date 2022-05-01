package snakegame

import (
	clearscreen "SnakeGameGolang/internal/clearscreen"
	"fmt"
	"math/rand"
	"sync"
	"time"

	"github.com/eiannone/keyboard"
)

type cell int8

const (
	cellEmpty cell = iota
	cellFood
	cellSnakeHead
	cellSnakeTail
)

// Board structure
type board struct {
	hight uint8
	width uint8
	matrix [][]cell
}

// Fill board matrix with zero values
func (b *board) clean() {
	for i := range b.matrix {
		for j := range b.matrix[i] {
			b.matrix[i][j] = cellEmpty
		}
	}
}

// Board initializaion
func (b *board) init(h uint8, w uint8) {
	b.hight = h
	b.width = w
	b.matrix = make([][]cell, h)
	for i := range b.matrix {
		b.matrix[i] = make([]cell, w)
	}
}

type direction int8

// Stores coordinates
type vertex struct {
	x, y uint8
}

const (
	directionUp direction = iota
	directionRight
	directionDown
	directionLeft
)

// Main snake game structure
type SnakeGame struct {
	snake []vertex
	food vertex
	ateFood bool
	mutex sync.Mutex
	direction direction
	board board
	gameOver bool
}

// Initialization
func (game *SnakeGame) Init(h uint8, w uint8) {
	game.board.init(h,w)
	game.direction = directionLeft
	game.snake = []vertex{{w/2,h/2},{w/2,h/2+1}}
	game.gameOver = false
	game.generateFood()
}

// Print board matrix
func (game *SnakeGame) PrintBoard() {
	game.mutex.Lock()
	defer game.mutex.Unlock()

	clearscreen.ClearScreen()
	for hight := range game.board.matrix {
		for widht := range game.board.matrix[hight] {
			switch game.board.matrix[hight][widht] {
			case cellEmpty:
				fmt.Print("_")
			case cellFood:
				fmt.Print("$")
			case cellSnakeHead:
				fmt.Print("%")
			case cellSnakeTail:
				fmt.Print("*")
			}
		}
		fmt.Println()
	}
}

// Update internal board-matrix with actual snake and food coordinates
func (game *SnakeGame) refreshBoard() {
	game.mutex.Lock()
	defer game.mutex.Unlock()

	game.board.clean()
	for i,v := range game.snake {
		if i == 0 {
			game.board.matrix[v.y][v.x] = cellSnakeHead
		} else {
			game.board.matrix[v.y][v.x] = cellSnakeTail
		}
	}

	game.board.matrix[game.food.y][game.food.x] = cellFood
}

// Calculate and update the internal board matrix
func (game *SnakeGame) calculateIteration() {
	game.mutex.Lock()
	defer game.mutex.Unlock()

	// Move snake and grow if food eaten
	tailEnd := len(game.snake)-1
	for i := tailEnd; i >= 0; i-- {
		if i == tailEnd && game.ateFood {
			game.snake = append(game.snake, game.snake[tailEnd])
			game.ateFood = false
		}
		if (i > 0) {
			game.snake[i] = game.snake[i-1]
		}
	}

	// Check if faced with the border
	switch game.direction {
	case directionUp:
		if game.snake[0].y == 0 {
			game.gameOver = true
			return
		}
		game.snake[0].y -= 1
	case directionRight:
		if game.snake[0].x == game.board.width-1 {
			game.gameOver = true
			return
		}
		game.snake[0].x += 1
	case directionDown:
		if game.snake[0].y == game.board.hight-1 {
			game.gameOver = true
			return
		}
		game.snake[0].y += 1
	case directionLeft:
		if game.snake[0].x == 0 {
			game.gameOver = true
			return
		}
		game.snake[0].x -= 1
	}

	// Check if faced with ourself
	for _,tail := range game.snake[1:] {
		if tail == game.snake[0] {
			game.gameOver = true
			return
		}
	}

	// Check if ate the food
	if game.snake[0] == game.food {
		game.ateFood = true
		game.generateFood()
	}
}

// Re-generate food coordinates
func (game *SnakeGame) generateFood() {
	var v vertex
	for {
		v = vertex{
			x: (uint8)(rand.Intn(int(game.board.width-1))),
			y: (uint8)(rand.Intn(int(game.board.hight-1))),
		}

		// Regenerate if food created "in snake"
		inSnake := false
		for _,e := range game.snake {
			if e == v {
				inSnake = true
				break
			}
		}
		if !inSnake {
			break
		}
	}
	game.food = v
}

// Thread-safe direction change
func (game *SnakeGame) changeDirection(newDirection direction) {
	game.mutex.Lock()
	defer game.mutex.Unlock()

	if len(game.snake) == 1 {
		game.direction = newDirection
		return
	}

	switch newDirection {
	case directionUp:
		if (game.snake[1].x == game.snake[0].x && game.snake[1].y == game.snake[0].y-1) {
			break
		}
		game.direction = newDirection
	case directionRight:
		if (game.snake[1].y == game.snake[0].y && game.snake[1].x == game.snake[0].x+1) {
			break
		}
		game.direction = newDirection
	case directionDown:
		if (game.snake[1].x == game.snake[0].x && game.snake[1].y == game.snake[0].y+1) {
			break
		}
		game.direction = newDirection
	case directionLeft:
		if (game.snake[1].y == game.snake[0].y && game.snake[1].x == game.snake[0].x-1) {
			break
		}
		game.direction = newDirection
	}
}

// Run key-handler thread
func (game *SnakeGame) runController() {
	go func() {
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
				game.changeDirection(directionUp)
			case keyboard.KeyArrowRight:
				game.changeDirection(directionRight)
			case keyboard.KeyArrowDown:
				game.changeDirection(directionDown)
			case keyboard.KeyArrowLeft:
				game.changeDirection(directionLeft)
			}

			if event.Key == keyboard.KeyEsc {
				break
			}
		}
	}()
}

// Display game over screen
func printGameOver() {
	clearscreen.ClearScreen()
	fmt.Println("<< GAME_OVER >>")
}

// Run main loop
func (game *SnakeGame) Run() {
	game.runController()

	for {
		game.calculateIteration()
		if game.gameOver {
			printGameOver()
			return
		}

		game.refreshBoard()
		game.PrintBoard()
		time.Sleep(time.Second/4)
	}
}
