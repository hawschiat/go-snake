package game

import (
	"context"
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
	snakeLength     int
}

func sameCoordinate(a coordinate, b coordinate) bool {
	return a.X == b.X && a.Y == b.Y
}

func snakeBittenItself(head *snakePart) bool {
	currentCheck := head.nextPart
	for currentCheck != nil {
		if sameCoordinate(head.coordinate, currentCheck.coordinate) {
			return true
		}
		currentCheck = currentCheck.nextPart
	}
	return false
}

func advanceGame(state *gameState, gameOver *bool, width int, height int) {
	scoreMultiplier := 1
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

	if newHeadCoord.X <= 0 || newHeadCoord.Y <= 0 ||
		newHeadCoord.X >= width || newHeadCoord.Y >= height {
		*gameOver = true
		return
	}

	if sameCoordinate(newHeadCoord, state.fruitCoordinate) {
		if state.snakeLength >= Const.MaxLength {
			scoreMultiplier++
		} else {
			currentHead := state.snakeHead
			state.snakeHead = &snakePart{newHeadCoord, currentHead}
			state.snakeLength++
		}
		// Adjust speed
		if state.snakeLength > 5 {
			state.speed++
		}
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

	if snakeBittenItself(state.snakeHead) {
		*gameOver = true
		return
	}

	if fruitEaten {
		// Add points
		state.score += 10 * scoreMultiplier * (state.speed / 5)

		// Generate new location for fruit
		newCoord := coordinate{}
		locationDecided := false

		for !locationDecided {
			newCoord = coordinate{
				(width * 25 / 100) + rand.Intn(width*60/100),
				(height * 25 / 100) + rand.Intn(height*60/100),
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

	state.fruitCoordinate = coordinate{startingPoint.X, startingPoint.Y - 5}
	state.snakeHead = &head
	state.snakeLength = Const.InitialLength
	state.speed = Const.InitialSpeed
	state.direction = "up"
	state.score = 0
}

// LaunchGame initializes and launches the Snake game
func LaunchGame(keysEvents <-chan keyboard.KeyEvent, width int, height int) {
	gamePause := make(chan bool)
	direction := make(chan string)

	fmt.Print("\033[?25l") // Hides cursor

	state := gameState{}
	initializeGame(&state, width, height-2)

	ctx, cancel := context.WithCancel(context.Background())

	go listenGameControl(ctx, keysEvents, direction, gamePause)

	defer func() {
		cancel()
		displayGameOverBox(keysEvents, width, height, state.score)
	}()

	over := false
	paused := false
	advanceCountdown := 0
	hasChangedDirection := false // Safety mechanism to prevent spamming direction command

	for {
		select {
		case paused = <-gamePause:
		case dir := <-direction:
			if !hasChangedDirection {
				if state.direction == "up" || state.direction == "down" {
					if dir == "left" || dir == "right" {
						state.direction = dir
						hasChangedDirection = true
					}
				} else if state.direction == "left" || state.direction == "right" {
					if dir == "up" || dir == "down" {
						state.direction = dir
						hasChangedDirection = true
					}
				} else {
					panic("I don't know where the snake is heading")
				}
			}
		default:
		}

		if advanceCountdown <= 0 {
			advanceGame(&state, &over, width, height)
			advanceCountdown = Const.Fps
			hasChangedDirection = false
		} else {
			advanceCountdown -= state.speed
		}

		drawGame(&state, width, height)
		if paused {
			showInGameMenu(keysEvents, width, height, &over)
			paused = false
			gamePause <- false
		}

		if over {
			return
		}

		time.Sleep(time.Second / Const.Fps)
	}
}
