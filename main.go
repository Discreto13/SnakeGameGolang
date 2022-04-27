package main

import (
	"fmt"
	"os"
	cleanscreen "snake/vendor/clearscreen"
	"time"
)

var height, width int

func init(){
	fmt.Println("Welcome to Snake game!")
	fmt.Println("Please, specify size of board.")
	fmt.Print("Height: ")
	_,err:= fmt.Scanln(&height)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(-1)
	}

	fmt.Print("Width: ")
	_,err= fmt.Scanln(&width)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(-1)
	}
}

func main() {
	fmt.Println("ABC")
	time.Sleep(time.Second * 3)
	cleanscreen.ClearScreen()
	fmt.Println("CBA")
	time.Sleep(time.Second * 3)
}
