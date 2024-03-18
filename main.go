package main

import (
	"errors"
	"os"
	"strings"
)

func main() {
	if len(os.Args) <= 1 {
		defaultTheme.ShowNoArgsError()
		return
	}

	var term, err = DefaultTui()
	if err != nil {
		term.ShowWarning(err)
	}

	if !IsFlag(os.Args[1], "terminal") {
		term.Cursor(strings.Join(os.Args[1:], " ") + "\n")
		err = term.Gitp(os.Args[1], os.Args[2:]...)
	} else if nArgs := len(os.Args); nArgs == 2 {
		err = term.InteractiveGitp("", false)
	} else if name, _, valid := ParseFlag(os.Args[2]); !valid || name != "keep-alive" {
		err = errors.New("Invalid argument: the only flag allowed for 'gitp terminal' is: --keep-alive")
	} else {
		if nArgs > 3 {
			term.ShowWarning("Too many arguments, ignoring unnecessary ones")
		}
		term.InteractiveGitp("", true)
	}

	if err != nil {
		term.ShowError(err)
	}
}
