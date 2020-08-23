package menu

import (
	"fmt"
	"strings"
	"time"

	Const "github.com/hawschiat/go-snake/internal/constants"
	Game "github.com/hawschiat/go-snake/internal/game"

	"github.com/eiannone/keyboard"
	"github.com/inancgumus/screen"
)

func max(a int, b int) int {
	if a > b {
		return a
	}
	return b
}

func center(s string, w int) string {
	// https://stackoverflow.com/questions/41133006/how-to-fmt-printprint-this-on-the-center
	return fmt.Sprintf("%[1]*s", -w, fmt.Sprintf("%[1]*s", (w+len(s))/2, s))
}

func printMenuText(active bool, s string, width int) {
	if active {
		fmt.Println("|", center(fmt.Sprintf("\033[30;47m %s \033[0m", s), width+8), "|")
	} else {
		fmt.Println("|", center(s, width-4), "|")
	}
}

func printMenu(width int, height int, menuIndex chan int) {
	titleRow := 3
	menuRow := height / 2

	fmt.Print("\033[?25l") // Hides cursor
	screen.Clear()

	defer func() {
		screen.MoveTopLeft()

		fmt.Println(center("_____  _      __    _      _     __       ____  ___   ___       ___   _      __    _     _   _      __    __ ", width))
		fmt.Println(center(" | |  | |_|  / /\\  | |\\ | | |_/ ( (`     | |_  / / \\ | |_)     | |_) | |    / /\\  \\ \\_/ | | | |\\ | / /`_ |_/", width))
		fmt.Println(center(" |_|  |_| | /_/--\\ |_| \\| |_| \\ _)_)     |_|   \\_\\_/ |_| \\     |_|   |_|__ /_/--\\  |_|  |_| |_| \\| \\_\\_/ (_) ", width))
		fmt.Println(strings.Repeat(" ", width))

		fmt.Print("\033[?25h") // Shows cursor
	}()

	for {
		selectedIndex, ok := <-menuIndex
		if !ok {
			return
		}
		screen.MoveTopLeft()
		fmt.Println(strings.Repeat("-", width))
		for i := 1; i < height-1; i++ {
			if i == titleRow {
				fmt.Println("|", center(" _______  _______         _______  _        _______  _        _______ ", width-4), "|")
				fmt.Println("|", center("(  ____ \\(  ___  )       (  ____ \\( (    /|(  ___  )| \\    /\\(  ____ \\", width-4), "|")
				fmt.Println("|", center("| (    \\/| (   ) |       | (    \\/|  \\  ( || (   ) ||  \\  / /| (    \\/", width-4), "|")
				fmt.Println("|", center("| |      | |   | | _____ | (_____ |   \\ | || (___) ||  (_/ / | (__    ", width-4), "|")
				fmt.Println("|", center("| | ____ | |   | |(_____)(_____  )| (\\ \\) ||  ___  ||   _ (  |  __)   ", width-4), "|")
				fmt.Println("|", center("| | \\_  )| |   | |             ) || | \\   || (   ) ||  ( \\ \\ | (      ", width-4), "|")
				fmt.Println("|", center("| (___) || (___) |       /\\____) || )  \\  || )   ( ||  /  \\ \\| (____/\\", width-4), "|")
				fmt.Println("|", center("(_______)(_______)       \\_______)|/    )_)|/     \\||_/    \\/(_______/", width-4), "|")
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
		fmt.Println(strings.Repeat("-", width))
		time.Sleep(time.Second / Const.Fps)
	}
}

func listenNavigate(launchCallback func(), menuIndex chan int) {
	if err := keyboard.Open(); err != nil {
		panic(err)
	}
	defer func() {
		_ = keyboard.Close()
		close(menuIndex)
	}()

	index := 0
	menuIndex <- index

	for {
		_, key, err := keyboard.GetKey()
		if err != nil {
			panic(err)
		}
		// Determine action based on input
		switch key {
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
				launchCallback()
			case 2:
				return
			}
		case keyboard.KeyCtrlC:
		case keyboard.KeyEsc:
			return
		}
	}
}

// LaunchMenu prints out the main menu and wait for user's command
func LaunchMenu() {
	w, h := screen.Size()

	if w < Const.MinWidth || h < Const.MinHeight {
		fmt.Println()
		fmt.Print("The console window is too small for the game! ")
		fmt.Printf("Please ensure the window is at least %dx%d and try again!\n", Const.MinWidth, Const.MinHeight)
		fmt.Println()
		return
	}

	w = max(w, Const.MinWidth)
	h = max(h, Const.MinHeight) - 1

	menuIndex := make(chan int)

	go listenNavigate(func() { Game.LaunchGame(w, h) }, menuIndex)
	printMenu(w, h, menuIndex)
}
