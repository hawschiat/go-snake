package game

import (
	"fmt"
	"strings"
	"time"

	Const "github.com/hawschiat/go-snake/internal/constants"
	Util "github.com/hawschiat/go-snake/internal/utilities"

	"github.com/eiannone/keyboard"
	"github.com/inancgumus/screen"
)

func printMenuText(active bool, s string, width int) {
	if active {
		fmt.Println("|", Util.Center(fmt.Sprintf("\033[30;47m %s \033[0m", s), width+8), "|")
	} else {
		fmt.Println("|", Util.Center(s, width-4), "|")
	}
}

func printMenu(
	keysEvents <-chan keyboard.KeyEvent,
	width int,
	height int,
	menuIndex chan int,
	exitGame chan bool,
	shouldLaunchGame chan bool,
) {
	titleRow := 3
	menuRow := height / 2

	fmt.Print("\033[?25l") // Hides cursor
	screen.Clear()

	defer func() {
		fmt.Print("\033[?25h") // Shows cursor
	}()

	willExit := false
	selectedIndex := 0
	ok := true

	for {
		select {
		case launch := <-shouldLaunchGame:
			if launch {
				LaunchGame(keysEvents, width, height)
				return
			}
		case selectedIndex, ok = <-menuIndex:
		case willExit = <-exitGame:
		default:
		}

		if !ok || willExit {
			screen.MoveTopLeft()

			fmt.Println(fmt.Sprintf("%[1]*s", -width, "_____  _      __    _      _     __       ____  ___   ___       ___   _      __    _     _   _      __    __ "))
			fmt.Println(fmt.Sprintf("%[1]*s", -width, " | |  | |_|  / /\\  | |\\ | | |_/ ( (`     | |_  / / \\ | |_)     | |_) | |    / /\\  \\ \\_/ | | | |\\ | / /`_ |_/"))
			fmt.Println(fmt.Sprintf("%[1]*s", -width, " |_|  |_| | /_/--\\ |_| \\| |_| \\ _)_)     |_|   \\_\\_/ |_| \\     |_|   |_|__ /_/--\\  |_|  |_| |_| \\| \\_\\_/ (_) "))
			fmt.Println(strings.Repeat(" ", width))

			return
		}

		screen.MoveTopLeft()
		fmt.Println(strings.Repeat("-", width))
		for i := 1; i < height-2; i++ {
			if i == titleRow {
				fmt.Println("|", Util.Center(" _______  _______         _______  _        _______  _        _______ ", width-4), "|")
				fmt.Println("|", Util.Center("(  ____ \\(  ___  )       (  ____ \\( (    /|(  ___  )| \\    /\\(  ____ \\", width-4), "|")
				fmt.Println("|", Util.Center("| (    \\/| (   ) |       | (    \\/|  \\  ( || (   ) ||  \\  / /| (    \\/", width-4), "|")
				fmt.Println("|", Util.Center("| |      | |   | | _____ | (_____ |   \\ | || (___) ||  (_/ / | (__    ", width-4), "|")
				fmt.Println("|", Util.Center("| | ____ | |   | |(_____)(_____  )| (\\ \\) ||  ___  ||   _ (  |  __)   ", width-4), "|")
				fmt.Println("|", Util.Center("| | \\_  )| |   | |             ) || | \\   || (   ) ||  ( \\ \\ | (      ", width-4), "|")
				fmt.Println("|", Util.Center("| (___) || (___) |       /\\____) || )  \\  || )   ( ||  /  \\ \\| (____/\\", width-4), "|")
				fmt.Println("|", Util.Center("(_______)(_______)       \\_______)|/    )_)|/     \\||_/    \\/(_______/", width-4), "|")
				i += 7
			} else if i == menuRow {
				printMenuText(selectedIndex == 0, "Start Game", width)
				printMenuText(selectedIndex == 1, "Info", width)
				printMenuText(selectedIndex == 2, "Quit", width)
				i += 2
			} else {
				fmt.Println("|" + strings.Repeat(" ", width-2) + "|")
			}
		}
		fmt.Println("|" + fmt.Sprintf("%[1]*s", width-2, "v1.0 by SC Haw ") + "|")
		fmt.Println(strings.Repeat("-", width))
		time.Sleep(time.Second / Const.Fps)
	}
}

func listenNavigate(
	keysEvents <-chan keyboard.KeyEvent,
	menuIndex chan int,
	exitGame chan bool,
	shouldLaunchGame chan bool,
) {
	defer func() {
		close(menuIndex)
		close(exitGame)
		close(shouldLaunchGame)
	}()

	index := 0
	menuIndex <- index

	for {
		event := <-keysEvents
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
			if index < 2 {
				index++
				menuIndex <- index
			}
		case keyboard.KeyEnter:
			switch index {
			case 0:
				shouldLaunchGame <- true
				return
			case 2:
				exitGame <- true
				return
			}
		case keyboard.KeyCtrlC:
		case keyboard.KeyEsc:
			return
		}
	}
}

// LaunchMenu prints out the main menu and wait for user's command
func LaunchMenu(keysEvents <-chan keyboard.KeyEvent) {
	w, h := screen.Size()

	if w < Const.MinWidth || h < Const.MinHeight {
		fmt.Println()
		fmt.Print("The console window is too small for the game! ")
		fmt.Printf("Please ensure the window is at least %dx%d and try again!\n", Const.MinWidth, Const.MinHeight)
		fmt.Println()
		return
	}

	w = Util.Max(w, Const.MinWidth) - 1
	h = Util.Max(h, Const.MinHeight) - 1

	menuIndex := make(chan int)
	exitGame := make(chan bool)
	shouldLaunchGame := make(chan bool)

	go listenNavigate(keysEvents, menuIndex, exitGame, shouldLaunchGame)
	printMenu(keysEvents, w, h, menuIndex, exitGame, shouldLaunchGame)
}
