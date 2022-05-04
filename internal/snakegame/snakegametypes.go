package snakegame

type Cell int8

const (
	CellEmpty Cell = iota
	CellFood
	CellSnakeHead
	CellSnakeTail
)

type Direction int8

const (
	DirectionUp Direction = iota
	DirectionRight
	DirectionDown
	DirectionLeft
)

type DisplayFunc func(board [][]Cell, score int)
type KeyHandlerFunc func(quit chan bool, turn chan Direction)

// Board structure
type board struct {
	hight  uint8
	width  uint8
	matrix [][]Cell
}

// Stores coordinates
type vertex struct {
	x, y uint8
}
