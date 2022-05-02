package snakegame

import (
	"math/rand"
	"time"
)

// Main snake game structure
type SnakeGame struct {
	board board
	food  vertex
	snake []vertex

	moveDirection Direction
	turnDirection chan Direction

	ateFood      bool
	borderKiller bool
	gameOver     bool
	quit         chan bool

	controller   ControllerFunc
	displayBoard DisplayBoardFunc
}

// Initialization
func (game *SnakeGame) Init(boardHight uint8, boardWidth uint8, borderKiller bool, displayBoard DisplayBoardFunc, controller ControllerFunc) {
	if boardHight > 100 || boardHight < 0 || boardWidth > 100 || boardWidth < 0 {
		panic("Expected board size [0-100]:[0-100] was not satisfied")
	}
	game.board.init(boardHight, boardWidth)

	if displayBoard == nil || controller == nil {
		panic("Control and display methods is not initialized")
	}
	game.controller = controller
	game.displayBoard = displayBoard

	game.turnDirection = make(chan Direction, 10)
	game.quit = make(chan bool, 1)
	game.moveDirection = DirectionUp
	game.snake = []vertex{{boardWidth / 2, boardHight / 2}}
	game.borderKiller = borderKiller
	game.generateFood()
}

// Run main loop
func (game *SnakeGame) Run() int {
	game.runControllerThread()

	for {
		game.calculateIteration()
		if game.isQuit() || game.gameOver {
			return len(game.snake) - 1
		}

		game.refreshBoard()
		game.printBoard()
		time.Sleep(time.Second / 5)
	}
}

// Fill board matrix with zero values
func (b *board) clean() {
	for i := range b.matrix {
		for j := range b.matrix[i] {
			b.matrix[i][j] = CellEmpty
		}
	}
}

// Board initializaion
func (b *board) init(boardHight uint8, boardWidth uint8) {
	b.hight = boardHight
	b.width = boardWidth
	b.matrix = make([][]Cell, boardHight)
	for i := range b.matrix {
		b.matrix[i] = make([]Cell, boardWidth)
	}
}

// Print board matrix
func (game *SnakeGame) printBoard() {
	if game.controller == nil {
		panic("Display method is not initialized")
	}
	game.displayBoard(game.board.matrix)
}

// Run key-handler thread
func (game *SnakeGame) runControllerThread() {
	if game.controller == nil {
		panic("Controller method is not initialized")
	}
	go game.controller(game.quit, game.turnDirection)
}

// Update internal board-matrix with actual snake and food coordinates
func (game *SnakeGame) refreshBoard() {
	game.board.clean()
	for i, v := range game.snake {
		if i == 0 {
			game.board.matrix[v.y][v.x] = CellSnakeHead
		} else {
			game.board.matrix[v.y][v.x] = CellSnakeTail
		}
	}

	game.board.matrix[game.food.y][game.food.x] = CellFood
}

// Calculate and update the internal board matrix
func (game *SnakeGame) calculateIteration() {
	// Move snake body and grow if food eaten
	tailEnd := len(game.snake) - 1
	for i := tailEnd; i >= 0; i-- {
		if i == tailEnd && game.ateFood {
			game.snake = append(game.snake, game.snake[tailEnd])
			game.ateFood = false
		}
		if i > 0 {
			game.snake[i] = game.snake[i-1]
		}
	}

	// Update direction and check if faced with the border
	game.updateDirection()
	game.moveSnakeHead()

	// Check if faced with ourself
	for _, tail := range game.snake[1:] {
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
			x: (uint8)(rand.Intn(int(game.board.width - 1))),
			y: (uint8)(rand.Intn(int(game.board.hight - 1))),
		}

		// Regenerate if food created "in snake"
		inSnake := false
		for _, e := range game.snake {
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

// Check and change direction according to signal
func (game *SnakeGame) updateDirection() {
	var newDirection Direction

	for {
		select {
		case newDirection = <-game.turnDirection:
		default:
			return
		}

		// Ignore inapplicable turn triggers
		if newDirection == game.moveDirection ||
			(newDirection == DirectionUp && game.moveDirection == DirectionDown) ||
			(newDirection == DirectionRight && game.moveDirection == DirectionLeft) ||
			(newDirection == DirectionDown && game.moveDirection == DirectionUp) ||
			(newDirection == DirectionLeft && game.moveDirection == DirectionRight) {
			continue
		}

		game.moveDirection = newDirection
		break
	}
}

// Move snake head and handle border interaction
func (game *SnakeGame) moveSnakeHead() {
	switch game.moveDirection {
	case DirectionUp:
		if game.snake[0].y != 0 {
			game.snake[0].y -= 1
			break
		}

		if game.borderKiller {
			game.gameOver = true
			return
		}
		game.snake[0].y = game.board.hight - 1

	case DirectionRight:
		if game.snake[0].x != game.board.width-1 {
			game.snake[0].x += 1
			break
		}

		if game.borderKiller {
			game.gameOver = true
			return
		}
		game.snake[0].x = 0

	case DirectionDown:
		if game.snake[0].y != game.board.hight-1 {
			game.snake[0].y += 1
			break
		}

		if game.borderKiller {
			game.gameOver = true
			return
		}
		game.snake[0].y = 0
	case DirectionLeft:
		if game.snake[0].x != 0 {
			game.snake[0].x -= 1
			break
		}

		if game.borderKiller {
			game.gameOver = true
			return
		}
		game.snake[0].x = game.board.width - 1
	}
}

// Exit initiation
func (game *SnakeGame) isQuit() bool {
	select {
	case <-game.quit:
		return true
	default:
		return false
	}
}
