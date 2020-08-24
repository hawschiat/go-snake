package main

import (
	"github.com/eiannone/keyboard"
	Game "github.com/hawschiat/go-snake/internal/game"
)

func main() {
	keysEvents, err := keyboard.GetKeys(1)
	if err != nil {
		panic(err)
	}

	defer func() {
		_ = keyboard.Close()
	}()

	Game.LaunchMenu(keysEvents)
}
