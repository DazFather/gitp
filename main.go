package main

import (
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

	if IsFlag(os.Args[1], "terminal") {
		err = term.InteractiveGitp("", len(os.Args) > 2 && IsFlag(os.Args[2], "keep-alive"))
	} else {
		term.Cursor(strings.Join(os.Args[1:], " ") + "\n")
		err = term.Gitp(os.Args[1], os.Args[2:]...)
	}

	if err != nil {
		term.ShowError(err)
	}
}
