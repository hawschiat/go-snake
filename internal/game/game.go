package game

import (
	"fmt"
	"math/rand"
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
	score           int
	direction       string
	snakeHead       *snakePart
}

func sameCoordinate(a coordinate, b coordinate) bool {
	return a.X == b.X && a.Y == b.Y
}

func advanceGame(state *gameState, gameOver chan bool, width int, height int) {
	fruitEaten := false
	currentHeadCoord := state.snakeHead.coordinate
	newHeadCoord := coordinate{}

	switch state.direction {
	case "up":
		newHeadCoord.X = currentHeadCoord.X
		newHeadCoord.Y = currentHeadCoord.Y - 1
	case "down":
		newHeadCoord.X = currentHeadCoord.X
		newHeadCoord.Y = currentHeadCoord.Y + 1
	case "left":
		newHeadCoord.X = currentHeadCoord.X - 1
		newHeadCoord.Y = currentHeadCoord.Y
	case "right":
		newHeadCoord.X = currentHeadCoord.X + 1
		newHeadCoord.Y = currentHeadCoord.Y
	default:
		panic("I don't know where the snake is heading")
	}

	if newHeadCoord.X < 0 || newHeadCoord.Y < 0 ||
		newHeadCoord.X > width || newHeadCoord.Y > height {
		gameOver <- true
		return
	}

	if sameCoordinate(newHeadCoord, state.fruitCoordinate) {
		currentHead := state.snakeHead
		state.snakeHead = &snakePart{newHeadCoord, currentHead}
		fruitEaten = true
	} else {
		// Move each blocks towards the next coordinates
		nextCoord := newHeadCoord
		currentPart := state.snakeHead
		for currentPart != nil {
			replacedCoord := currentPart.coordinate
			currentPart.coordinate = nextCoord
			nextCoord = replacedCoord
			currentPart = currentPart.nextPart
		}
	}

	if fruitEaten {
		// Add points
		state.score += 10

		// Generate new location for fruit
		newCoord := coordinate{}
		locationDecided := false

		for !locationDecided {
			newCoord = coordinate{
				rand.Intn(width - 1),
				rand.Intn(height - 1),
			}
			locationDecided = true

			// Ensure that the fruit is not placed on the snake
			currentPart := state.snakeHead
			for currentPart != nil {
				locationDecided = !sameCoordinate(newCoord, currentPart.coordinate)
				currentPart = currentPart.nextPart
			}
		}

		state.fruitCoordinate = newCoord
	}
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
	state.direction = "up"
	state.score = 0
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
		defer close(gamePause)

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
	advanceCountdown := 0

	for {
		select {
		case paused = <-gamePause:
		case dir := <-direction:
			if state.direction == "up" || state.direction == "down" {
				if dir == "left" || dir == "right" {
					state.direction = dir
				}
			} else if state.direction == "left" || state.direction == "right" {
				if dir == "up" || dir == "down" {
					state.direction = dir
				}
			} else {
				panic("I don't know where the snake is heading")
			}
		case <-gameOver:
			return
		default:
			break
		}

		drawGame(&state, width, height)
		if paused {
			showInGameMenu(keysEvents, width, height)
		}

		if advanceCountdown == 0 {
			advanceGame(&state, gameOver, width, height)
			advanceCountdown = 60
		} else {
			advanceCountdown -= state.speed
		}

		time.Sleep(time.Second / Const.Fps)
	}
}
