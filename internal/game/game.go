package game

import (
	"fmt"
	"time"

	"github.com/eiannone/keyboard"
	Const "github.com/hawschiat/go-snake/internal/constants"
)

type coordinate struct {
	X int
	Y int
}
type snakePart struct {
	coordinate coordinate
	nextPart   *snakePart
}
type gameState struct {
	fruitCoordinate coordinate
	speed           int
	snakeHead       *snakePart
}

func initializeGame(state *gameState, width int, height int) {
	startingPoint := coordinate{
		(width - 2) / 2,
		(height - 2) / 2,
	}

	head := snakePart{startingPoint, nil}
	currentPart := &head
	for i := 1; i < Const.InitialLength; i++ {
		currentPoint := currentPart.coordinate
		nextPoint := coordinate{currentPoint.X, currentPoint.Y + 1}
		// Assign next body part
		currentPart.nextPart = &snakePart{nextPoint, nil}
		currentPart = currentPart.nextPart
	}

	state.fruitCoordinate = coordinate{startingPoint.X, startingPoint.Y - 2}
	state.snakeHead = &head
	state.speed = Const.InitialSpeed
}

// LaunchGame initializes and launches the Snake game
func LaunchGame(keysEvents <-chan keyboard.KeyEvent, width int, height int) {
	gameOver := make(chan bool)
	gamePause := make(chan bool)
	direction := make(chan string)

	fmt.Print("\033[?25l") // Hides cursor

	defer func() {
		close(gameOver)
		fmt.Print("\033[?25h") // Shows cursor
		LaunchMenu(keysEvents)
	}()

	state := gameState{}
	initializeGame(&state, width, height)

	go listenGameControl(keysEvents, direction, gamePause, gameOver)

	go func(gamePause chan bool, gameOver chan bool) {
		defer func() {
			close(gamePause)
		}()

		paused := false
		for {
			event := <-keysEvents
			if event.Err != nil {
				panic(event.Err)
			}
			if event.Key == keyboard.KeyEsc {
				paused = !paused
				gamePause <- paused
			}
			if event.Key == keyboard.KeyEnter {
				gameOver <- true
				return
			}
		}
	}(gamePause, gameOver)

	paused := false

	for {
		drawGame(&state, width, height)

		select {
		case paused = <-gamePause:
		case <-gameOver:
			return
		default:
			break
		}

		if paused {
			showInGameMenu(keysEvents, width, height)
		}
		time.Sleep(time.Second / Const.Fps)
	}
}
