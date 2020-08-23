package game

import (
	"fmt"
	"strings"

	"github.com/inancgumus/screen"
)

func sameCoordinate(a coordinate, b coordinate) bool {
	return a.X == b.X && a.Y == b.Y
}

func drawGame(state *gameState, width int, height int) {
	screen.MoveTopLeft()
	fmt.Println(strings.Repeat("-", width))
	for row := 1; row < height-1; row++ {
		fmt.Println("|" + strings.Repeat(" ", width-2) + "|")
	}
	fmt.Println(strings.Repeat("-", width))

	// Draw fruit
	fruitCoord := state.fruitCoordinate
	fmt.Print(fmt.Sprintf("\033[%d;%dH", fruitCoord.Y, fruitCoord.X))
	fmt.Print("\u2591")

	// Draw snake body
	currentPart := state.snakeHead
	for currentPart != nil {
		coord := currentPart.coordinate
		// Move cursor to the coordinate
		fmt.Print(fmt.Sprintf("\033[%d;%dH", coord.Y, coord.X))
		fmt.Print("\u2588")
		currentPart = currentPart.nextPart
	}
	screen.MoveTopLeft()
}