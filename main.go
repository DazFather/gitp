package main

import "os"

func main() {
	if len(os.Args) <= 1 {
		defaultTheme.printError("Invalid given arguments")
		return
	}

	var (
		term = DefaultTui()
		err  error
	)

	switch os.Args[1] {
	case "terminal", "-terminal", "--terminal":
		if term == nil {
			defaultTheme.printError("Cannot create gitp+ terminal, missing git on project or system")
			return
		}
		err = term.InteractiveGitp("")
	case "-h", "--help", "help":
		if term == nil {
			term = ThemeWrapper(defaultTheme)
		}
		err = term.Gitp(os.Args[1], os.Args[2:]...)
	default:
		if term == nil {
			defaultTheme.printError("Cannot execute any command, missing git on project or system")
			return
		}
		term.Cursor()
		err = term.Gitp(os.Args[1], os.Args[2:]...)
	}

	if err != nil {
		term.ShowError(err)
	}
}
