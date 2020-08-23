package game

import (
	"fmt"
	"strings"

	"github.com/inancgumus/screen"
)

func drawGame(state *gameState, width int, height int) {
	screen.MoveTopLeft()
	fmt.Println(fmt.Sprintf("Score: %-10d", 0), fmt.Sprintf("%[1]*s", width-10, "Press \033[1mESC\033[0m for menu."))
	fmt.Println(strings.Repeat("-", width))
	for row := 2; row < height-1; row++ {
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
