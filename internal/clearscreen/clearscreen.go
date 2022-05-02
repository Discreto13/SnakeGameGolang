package clearscreen

import (
	"os"
	"os/exec"
	"runtime"
)

var clearScreen map[string]func()

func init() {
	clearScreen = make(map[string]func())
	clearScreen["linux"] = func() {
		cmd := exec.Command("clear")
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
	clearScreen["windows"] = func() {
		cmd := exec.Command("cmd", "/c", "cls")
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
}

func ClearScreen() {
	value, ok := clearScreen[runtime.GOOS]
	if ok {
		value()
	} else {
		panic("Unknow OS")
	}
}
