package game

import (
	"fmt"
	"strings"

	"github.com/eiannone/keyboard"
	Util "github.com/hawschiat/go-snake/internal/utilities"
	"github.com/inancgumus/screen"
)

const heightOffset = 2

func drawGame(state *gameState, width int, height int) {
	screen.MoveTopLeft()
	fmt.Println(fmt.Sprintf("Score: %-10d", state.score), fmt.Sprintf("%[1]*s", width-10, "Press \033[1mESC\033[0m for menu."))
	fmt.Println(strings.Repeat("-", width))
	for row := 2; row < height-1; row++ {
		fmt.Println("|" + strings.Repeat(" ", width-2) + "|")
	}
	fmt.Println(strings.Repeat("-", width))

	// Draw fruit
	fruitCoord := state.fruitCoordinate
	fmt.Print(fmt.Sprintf("\033[%d;%dH", fruitCoord.Y+heightOffset, fruitCoord.X))
	fmt.Print("\u2591")

	// Draw snake body
	currentPart := state.snakeHead
	for currentPart != nil {
		coord := currentPart.coordinate
		// Move cursor to the coordinate
		fmt.Print(fmt.Sprintf("\033[%d;%dH", coord.Y+heightOffset, coord.X))
		fmt.Print("\u2588")
		currentPart = currentPart.nextPart
	}
	screen.MoveTopLeft()
}

func displayGameOverBox(keysEvents <-chan keyboard.KeyEvent, width int, height int, finalScore int) {
	boxWidth := width * 50 / 100
	boxHeight := Util.Max(height*30/100, 20)
	startCoordinate := coordinate{
		(width / 2) - (boxWidth / 2),
		(height / 2) - (boxHeight / 2),
	}

	// Move cursor to the coordinate
	fmt.Print(fmt.Sprintf("\033[%d;%dH", startCoordinate.Y, startCoordinate.X))
	offset := 1

	fmt.Print(strings.Repeat("\u2594", boxWidth))
	for offset < boxHeight-1 {
		if offset == 2 {
			printInGameMenuText(false, "   ___                                  ", boxWidth, startCoordinate.Y+offset, startCoordinate.X)
			offset++
			printInGameMenuText(false, "  / __|__ _ _ __  ___   _____ _____ _ _ ", boxWidth, startCoordinate.Y+offset, startCoordinate.X)
			offset++
			printInGameMenuText(false, " | (_ / _` | '  \\/ -_) / _ \\ V / -_) '_|", boxWidth, startCoordinate.Y+offset, startCoordinate.X)
			offset++
			printInGameMenuText(false, "  \\___\\__,_|_|_|_\\___| \\___/\\_/\\___|_|  ", boxWidth, startCoordinate.Y+offset, startCoordinate.X)
		} else if offset == boxHeight-5 {
			printInGameMenuText(false, fmt.Sprintf("Your score is: \033[1m%d\033[0m", finalScore), boxWidth+8, startCoordinate.Y+offset, startCoordinate.X)
		} else if offset == boxHeight-2 {
			printInGameMenuText(false, "Press \033[1mEnter\033[0m to exit.", boxWidth+8, startCoordinate.Y+offset, startCoordinate.X)
		} else {
			fmt.Print(fmt.Sprintf("\033[%d;%dH", startCoordinate.Y+offset, startCoordinate.X))
			fmt.Print("\u258c" + strings.Repeat(" ", boxWidth-2) + "\u2590")
		}
		offset++
	}
	fmt.Print(fmt.Sprintf("\033[%d;%dH", startCoordinate.Y+offset, startCoordinate.X))
	fmt.Print(strings.Repeat("\u2581", boxWidth))

	for {
		event, ok := <-keysEvents
		if !ok {
			return
		}
		if event.Err != nil {
			panic(event.Err)
		}

		if event.Key == keyboard.KeyEnter {
			fmt.Print("\033[?25h") // Shows cursor
			LaunchMenu(keysEvents)
			return
		}
	}
}
