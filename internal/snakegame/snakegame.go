package snakegame

import (
	clearscreen "SnakeGameGolang/internal/clearscreen"
	"fmt"
	"math/rand"
	"sync"
	"time"

	"github.com/eiannone/keyboard"
)

type Cell int8

const (
	Empty Cell = iota
	Food
	SnakeHead
	SnakeTail
)

type Board struct {
	hight uint8
	width uint8
	matrix [][]Cell
}

func (b *Board) clean() {
	for i := range b.matrix {
		for j := range b.matrix[i] {
			b.matrix[i][j] = Empty
		}
	}
}

func (b *Board) init(h uint8, w uint8) {
	b.hight = h
	b.width = w
	b.matrix = make([][]Cell, h)
	for i := range b.matrix {
		b.matrix[i] = make([]Cell, w)
	}
}

type Direction int8

type Vertex struct {
	x, y uint8
}

const (
	Up Direction = iota
	Right
	Down
	Left
)

type SnakeGame struct {
	snake []Vertex
	food Vertex
	ateFood bool
	mutex sync.Mutex
	direction Direction
	board Board
	gameOver bool
}

func (game *SnakeGame) Init(h uint8, w uint8) {
	game.board.init(h,w)
	game.direction = Left
	// game.snake = []Vertex{{w/2,h/2},{w/2,h/2+1},{w/2,h/2+2},{w/2-1,h/2+2},{w/2-2,h/2+2},{w/2-3,h/2+2},{w/2-4,h/2+2},{w/2-5,h/2+2},{w/2-6,h/2+2} }
	game.snake = []Vertex{{w/2,h/2}}
	game.gameOver = false
	game.generateFood()
}

func (game *SnakeGame) PrintBoard() {
	game.mutex.Lock()
	defer game.mutex.Unlock()

	clearscreen.ClearScreen()
	for hight := range game.board.matrix {
		for widht := range game.board.matrix[hight] {
			switch game.board.matrix[hight][widht] {
			case Empty:
				fmt.Print("_")
			case Food:
				fmt.Print("$")
			case SnakeHead:
				fmt.Print("%")
			case SnakeTail:
				fmt.Print("*")
			}
		}
		fmt.Println()
	}
}

func (game *SnakeGame) refreshBoard() {
	game.mutex.Lock()
	defer game.mutex.Unlock()

	game.board.clean()
	for i,v := range game.snake {
		if i == 0 {
			game.board.matrix[v.y][v.x] = SnakeHead
		} else {
			game.board.matrix[v.y][v.x] = SnakeTail
		}
	}

	game.board.matrix[game.food.y][game.food.x] = Food
}

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
	case Up:
		if game.snake[0].y == 0 {
			game.gameOver = true
			return
		}
		game.snake[0].y -= 1
	case Right:
		if game.snake[0].x == game.board.width-1 {
			game.gameOver = true
			return
		}
		game.snake[0].x += 1
	case Down:
		if game.snake[0].y == game.board.hight-1 {
			game.gameOver = true
			return
		}
		game.snake[0].y += 1
	case Left:
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

func (game *SnakeGame) generateFood() {
	var v Vertex
	for {
		v = Vertex{
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
func (game *SnakeGame) changeDirection(newDirection Direction) {
	game.mutex.Lock()
	defer game.mutex.Unlock()

	if len(game.snake) == 1 {
		game.direction = newDirection
		return
	}

	switch newDirection {
	case Up:
		if (game.snake[1].x == game.snake[0].x && game.snake[1].y == game.snake[0].y-1) {
			break
		}
		game.direction = newDirection
	case Right:
		if (game.snake[1].y == game.snake[0].y && game.snake[1].x == game.snake[0].x+1) {
			break
		}
		game.direction = newDirection
	case Down:
		if (game.snake[1].x == game.snake[0].x && game.snake[1].y == game.snake[0].y+1) {
			break
		}
		game.direction = newDirection
	case Left:
		if (game.snake[1].y == game.snake[0].y && game.snake[1].x == game.snake[0].x-1) {
			break
		}
		game.direction = newDirection
	}
}

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
				game.changeDirection(Up)
			case keyboard.KeyArrowRight:
				game.changeDirection(Right)
			case keyboard.KeyArrowDown:
				game.changeDirection(Down)
			case keyboard.KeyArrowLeft:
				game.changeDirection(Left)
			}

			if event.Key == keyboard.KeyEsc {
				break
			}
		}
	}()
}

func printGameOver() {
	clearscreen.ClearScreen()
	fmt.Println("<< GAME_OVER >>")
}

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
