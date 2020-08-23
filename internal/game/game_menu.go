package game

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/eiannone/keyboard"
	Const "github.com/hawschiat/go-snake/internal/constants"
	Util "github.com/hawschiat/go-snake/internal/utilities"
)

func listenMenuControl(
	ctx context.Context,
	keysEvents <-chan keyboard.KeyEvent,
	menuIndex chan<- int,
	menuCommand chan<- string,
) {
	index := 0

	for {
		select {
		case <-ctx.Done():
			return
		case event, ok := <-keysEvents:
			if !ok {
				return
			}
			if event.Err != nil {
				panic(event.Err)
			}
			// Determine action based on input
			switch event.Key {
			case keyboard.KeyArrowUp:
				if index > 0 {
					index--
					menuIndex <- index
				}
			case keyboard.KeyArrowDown:
				if index < 1 {
					index++
					menuIndex <- index
				}
			case keyboard.KeyEnter:
				switch index {
				case 0:
					menuCommand <- "resume"
				case 1:
					menuCommand <- "quit"
					return
				}
			}
		}
	}
}

func printInGameMenuText(active bool, s string, width int, y int, x int) {
	fmt.Print(fmt.Sprintf("\033[%d;%dH", y, x))
	if active {
		fmt.Print("\u258c", Util.Center(fmt.Sprintf("\033[30;47m %s \033[0m", s), width+10), "\u2590")
	} else {
		fmt.Print("\u258c", Util.Center(s, width-2), "\u2590")
	}
}

func showInGameMenu(keysEvents <-chan keyboard.KeyEvent, width int, height int, gameOver *bool) {
	menuIndex := make(chan int)
	menuCommand := make(chan string)

	menuWidth := width * 30 / 100
	menuHeight := height * 50 / 100
	startCoordinate := coordinate{
		(width / 2) - (menuWidth / 2),
		(height / 2) - (menuHeight / 2),
	}
	titleRow := menuHeight * 30 / 100
	startingRow := menuHeight * 50 / 100

	selectedIndex := 0
	ok := true
	ctx, cancel := context.WithCancel(context.Background())

	go listenMenuControl(ctx, keysEvents, menuIndex, menuCommand)

	defer cancel()

	for {
		select {
		case selectedIndex, ok = <-menuIndex:
		case command := <-menuCommand:
			switch command {
			case "resume":
				return
			case "quit":
				*gameOver = true
				return
			}
		default:
		}

		if !ok {
			return
		}

		// Move cursor to the coordinate
		fmt.Print(fmt.Sprintf("\033[%d;%dH", startCoordinate.Y, startCoordinate.X))
		offset := 1

		fmt.Print(strings.Repeat("\u2594", menuWidth))
		for offset < menuHeight-1 {
			if offset == titleRow {
				printInGameMenuText(false, "\033[1mGame Paused\033[0m", menuWidth+8, startCoordinate.Y+offset, startCoordinate.X)
			} else if offset == startingRow {
				printInGameMenuText(selectedIndex == 0, "Resume", menuWidth, startCoordinate.Y+offset, startCoordinate.X)
				offset++
				printInGameMenuText(selectedIndex == 1, "Quit", menuWidth, startCoordinate.Y+offset, startCoordinate.X)
			} else {
				fmt.Print(fmt.Sprintf("\033[%d;%dH", startCoordinate.Y+offset, startCoordinate.X))
				fmt.Print("\u258c" + strings.Repeat(" ", menuWidth-2) + "\u2590")
			}
			offset++
		}
		fmt.Print(fmt.Sprintf("\033[%d;%dH", startCoordinate.Y+offset, startCoordinate.X))
		fmt.Print(strings.Repeat("\u2581", menuWidth))
		time.Sleep(time.Second / Const.Fps)
	}
}
