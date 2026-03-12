package cmd

import (
    "os"
    "fmt"
	"log"
    "golang.org/x/term"
)

func PromptSelect[T any](options []T, renderItem func (item T) string) T {
    selected := 0

    oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
    if err != nil {
        log.Fatal(err)
    }
    defer term.Restore(int(os.Stdin.Fd()), oldState)

    render := func() {
        // Clear all previously rendered lines
        for i := 0; i < len(options); i++ {
            fmt.Print("\033[2K\033[1A") // clear line, move up
        }
        fmt.Print("\033[2K") // clear current line

        for i, option := range options {
            if i == selected {
                fmt.Println(BoldCyan + "› " + renderItem(option) + Reset)
            } else {
                fmt.Println(Dim + "  " + renderItem(option) + Reset)
            }
        }
    }

    // Initial render
	for i, option := range options {
		if i == selected {
			fmt.Println(BoldCyan + "› " + renderItem(option) + Reset)
		} else {
			fmt.Println(Dim + "  " + renderItem(option) + Reset)
		}
	}

    buf := make([]byte, 3)
    for {
        n, _ := os.Stdin.Read(buf)

        if n == 3 && buf[0] == 27 && buf[1] == 91 {
            switch buf[2] {
            case 65: // up arrow
                if selected > 0 {
                    selected--
                }
            case 66: // down arrow
                if selected < len(options)-1 {
                    selected++
                }
            }
        } else if n == 1 {
            switch buf[0] {
            case 13: // enter
                return options[selected % len(options)]
            case 3: // ctrl+c
                term.Restore(int(os.Stdin.Fd()), oldState)
                os.Exit(0)
            }
        }

        render()
    }
}
