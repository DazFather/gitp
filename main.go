package main

import "os"

func main() {
	if len(os.Args) <= 1 {
		defaultTheme.ShowNoArgsError()
		return
	}

	var term, err = DefaultTui()

	switch os.Args[1] {
	case "terminal", "-terminal", "--terminal":
		err = term.InteractiveGitp("")
	default:
		term.Cursor("\n")
		err = term.Gitp(os.Args[1], os.Args[2:]...)
	}

	if err != nil {
		term.ShowError(err)
	}
}
