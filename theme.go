package main

import (
	"errors"
	"fmt"
	"strings"

	"github.com/DazFather/brush"
)

type theme interface {
	printCursor(branch, location, suffix string)
	printFlowStart(command string)
	printFlowEnd(command string)
	printCommand(branch, command string, args ...string)
	printOut(output string)
	printWarning(warning any)
	printError(err any)
	printHelp(gitHelp string)

	confirmRemoveBranch(branch string) (bool, error)
	ShowNoArgsError()
}

type palette[cSet brush.ColorType] struct {
	std, out                     brush.Brush[cSet]
	dir, branch, flow, command   brush.Brush[cSet]
	err, errDesc, warn, warnDesc brush.Brush[cSet]
}

var defaultTheme theme = palette[brush.ANSIColor]{
	std:      brush.New(brush.White, nil),
	dir:      brush.New(brush.White, brush.UseColor(brush.Blue)),
	branch:   brush.New(brush.White, brush.UseColor(brush.BrightBlue)),
	flow:     brush.New(brush.Green, nil),
	command:  brush.New(brush.BrightMagenta, nil),
	out:      brush.New(brush.BrightYellow, nil),
	err:      brush.New(brush.Black, brush.UseColor(brush.BrightRed)),
	errDesc:  brush.New(brush.BrightWhite, brush.UseColor(brush.Red)),
	warn:     brush.New(brush.BrightWhite, brush.UseColor(brush.Yellow)),
	warnDesc: brush.New(brush.Black, brush.UseColor(brush.BrightYellow)),
}

func (p palette[cSet]) printCommand(branch, command string, args ...string) {
	if len(args) > 0 {
		command += " " + strings.Join(args, " ")
	}

	p.printMessage(
		"• ",
		p.command.Paint("git ", command),
		": ",
	)
}

func (p palette[cSet]) printCursor(branch, location, suffix string) {
	p.printMessage(
		p.dir.Paint(" ", location, " "),
		p.branch.Paint(" ", branch, " "),
		" ► ",
		suffix,
	)
}

func (p palette[cSet]) printError(err any) {
	fmt.Printf("%s%s\n",
		p.errDesc.Paint(" ERROR "),
		p.err.Paint(" - ", err, " "),
	)
}

func (p palette[cSet]) printOut(output string) {
	p.out.Println(output)
}

func (p palette[cSet]) printWarning(warning any) {
	fmt.Printf("%s%s\n",
		p.warnDesc.Paint(" ! "),
		p.warn.Paint(" ", warning, " "),
	)
}

func (p palette[cSet]) printFlowStart(command string) {
	p.printMessage(
		p.flow.Paint(" » "),
		"executing ",
		p.command.Paint("git+ ", command),
		" command flow \n",
	)
}

func (p palette[cSet]) printFlowEnd(command string) {
	p.printMessage(
		p.flow.Paint(" ✓ "),
		p.command.Paint("git+ ", command),
		" executed successfully \n",
	)
}

func (p palette[cSet]) printHelp(output string) {
	p.printMessage(
		"gitp aka git+ is a cli that facilitate you when using git commands\n\n",
		p.dir.Paint(" Flows "), " are a list of commands that gets executed one after the other, for common tasks:\n",
		" • ", p.command.Paint("update"), " (status > stash* > fetch > pull > stash pop*): ",
		"update your branch with possible incoming remote changes\n",
		" • ", p.command.Paint("fork <branch-name>"), " (status > stash* > fetch > pull > checkout -b <branch> > push --set-upstream origin <branch> > stash pop*):",
		"update current branch and creates a new one from current with given name and sets remote upstream\n",
		" • ", p.command.Paint("undo [commit|branch|fork|merge|stash|upstream|add|stage] <args...>"), ": has different effects depending on input\n",
		"\t commit (reset HEAD~1 <args...>): reset last commit preserving changes locally by default\n",
		"\t merge (merge abort <args...>): abort current merge\n",
		"\t stash (stash pop <args...>): reapply last stashed item and remove it from the stack\n",
		"\t upstream (branch --unset-upstream <args...>): disable remote tracking from a branch\n",
		"\t add, stage: (restore --staged <args...>): remove matching files from stage\n",
		"\t branch: remove given branch, if missing the current one, it deletes also remote after a confirm, pass '--confirm' to skip\n",
		"\t fork (undo branch <branch-name> --confirm): a simple alias pre-confirmed to integrate better with fork flow\n",
		" • ", p.command.Paint("align <reference-branch>"), ": update current and reference branches and merge reference into current\n",
		p.dir.Paint(" Terminal "), " An interactive git command line that will constantly ask for new gitp+ flows or git commands.\n",
		"To use it simply launch this program using 'terminal', '--terminal' or '-terminal' as first argument.\n",
		"By default if a command result into an error the interactive terminal will stop, if you want to override this behaviour you can use --keep-alive.\n",
		"To escape just insert a blank line\n\n",
		p.dir.Paint(" Git Help "), p.out.Paint(" ", output, "\n"),
	)
}

func (p palette[cSet]) printMessage(values ...any) {
	fmt.Print(p.std.Embed(values...))
}

func (p palette[cSet]) confirmRemoveBranch(branch string) (bool, error) {
	var choice string

	p.printWarning("This action is not reversible")
	p.printMessage("\nConfirm: delete branch ", p.branch.Paint(" ", branch, " "), " also from remote?\n",
		"\t[", brush.Paint(brush.White, brush.UseColor(brush.Green), " Yes "), "]",
		" | ", brush.Paint(brush.White, brush.UseColor(brush.Red), " No "),
		" ► ",
	)

	if _, err := fmt.Scanln(&choice); err != nil {
		return false, err
	}

	choice = strings.TrimSpace(choice)
	switch strings.ToUpper(choice) {
	case "", "Y", "YES":
		return true, nil
	case "N", "NO":
		return false, nil
	}

	return false, errors.New("Invalid input '" + choice + "'")
}

func (p palette[cSet]) ShowNoArgsError() {
	p.printError("No given arguments, use 'gitp help' to learn more'")
}

func (t tui) Cursor(suffix string) {
	t.printCursor(t.branch, t.directory, suffix)
}

func (t tui) ShowError(err any) {
	t.printError(err)
}

func (t tui) ShowWarning(err any) {
	t.printWarning(err)
}
