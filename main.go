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
		if err == nil {
			err = term.InteractiveGitp("")
		}
	case "init", "clone", "-h", "--help", "help":
		err = term.Gitp(os.Args[1], os.Args[2:]...)
	default:
		if err == nil {
			term.Cursor()
			err = term.Gitp(os.Args[1], os.Args[2:]...)
		}
	}

	if err != nil {
		term.ShowError(err)
	}
}
