package game

import (
	"github.com/eiannone/keyboard"
)

func listenGameControl(
	keysEvents <-chan keyboard.KeyEvent,
	direction chan string,
	gamePause <-chan bool,
	gameOver chan bool,
) {
	defer func() {
		close(direction)
	}()

	paused := false

	dir := "up"
	direction <- dir

	for {
		select {
		case paused = <-gamePause:
		case event, ok := <-keysEvents:
			if !ok {
				return
			}
			if event.Err != nil {
				panic(event.Err)
			}
			// Handles navigation if game is not paused
			if !paused {
				// Determine action based on input
				switch event.Key {
				case keyboard.KeyArrowUp:
					if dir != "up" {
						dir = "up"
						direction <- dir
					}
				case keyboard.KeyArrowDown:
					if dir != "down" {
						dir = "down"
						direction <- dir
					}
				case keyboard.KeyArrowLeft:
					if dir != "left" {
						dir = "left"
						direction <- dir
					}
				case keyboard.KeyArrowRight:
					if dir != "right" {
						dir = "right"
						direction <- dir
					}
				}
			}
		case <-gameOver:
			return
		default:
		}
	}
}
