package main

import (
	"fmt"
	"strings"

	"github.com/DazFather/brush"
)

type theme[cSet brush.ColorType] struct {
	std                        brush.Brush[cSet]
	dir, branch, flow, command brush.Brush[cSet]
	out, err, errDesc          brush.Brush[cSet]
}

var defaultTheme = theme[brush.ANSIColor]{
	std:     brush.New(brush.White, nil),
	dir:     brush.New(brush.White, brush.UseColor(brush.Blue)),
	branch:  brush.New(brush.White, brush.UseColor(brush.BrightBlue)),
	flow:    brush.New(brush.Green, nil),
	command: brush.New(brush.BrightMagenta, nil),
	out:     brush.New(brush.BrightYellow, nil),
	err:     brush.New(brush.Black, brush.UseColor(brush.BrightRed)),
	errDesc: brush.New(brush.BrightWhite, brush.UseColor(brush.Red)),
}

func (t theme[cSet]) printCommand(branch, command string, args ...string) {
	if len(args) > 0 {
		command += " " + strings.Join(args, " ")
	}

	fmt.Print(t.std.Embed(
		"• ",
		t.command.Paint("git ", command),
		": ",
	))
}

func (t theme[cSet]) printCursor(branch, location string) {
	fmt.Print(t.std.Embed(
		t.dir.Paint(" ", location, " "),
		t.branch.Paint(" ", branch, " "),
		" ► ",
	))
}

func (t theme[cSet]) printError(err any) {
	fmt.Printf("%s%s\n",
		t.errDesc.Paint(" ERROR "),
		t.err.Paint(" - ", err),
	)
}

func (t theme[cSet]) printOut(output string) {
	t.out.Print(output, " \n")
	// t.std.Print("\n")
}

func (t theme[cSet]) printFlowStart(command string) {
	fmt.Print(t.std.Embed(
		t.flow.Paint(" » "),
		"executing ",
		t.command.Paint("git+ ", command),
		" command flow \n",
	))
}

func (t theme[cSet]) printFlowEnd(command string) {
	fmt.Print(t.std.Embed(
		t.flow.Paint(" ✓ "),
		t.command.Paint("git+ ", command),
		" executed successfully \n",
	))
}

func (t theme[cSet]) printHelp(output string) {
	fmt.Println(t.std.Embed(
		"gitp aka git+ is a cli that facilitate you when using git commands\n\n",
		t.dir.Paint(" Flows "), " are a list of commands that gets executed one after the other, for common tasks:\n",
		" • ", t.command.Paint("update"), " (status > stash* > fetch > pull > stash pop*): ",
		"update your branch with possible incoming remote changes\n",
		" • ", t.command.Paint("fork <branch-name>"), " (status > stash* > fetch > pull > checkout -b <branch> > push --set-upstream origin <branch> > stash pop*):",
		"update current branch and creates a new one from current with given name and sets remote upstream\n",
		" • ", t.command.Paint("undo [commit|branch|merge|stash|upstream] <args...>"), ": has different effects depending on input\n",
		"\t commit (reset HEAD~1 <args...>): reset last commit preserving changes locally by default\n",
		"\t merge (merge abort <args...>): abort current merge\n",
		"\t stash (stash pop <args...>): reapply last stashed item and remove it from the stack\n",
		"\t upstream (--unset-upstream <args...>): disable remote tracking from a branch\n",
		"\t branch: it remove given branch or if missing the current one, it deletes also remote\n",
		t.dir.Paint(" Terminal "), " An interactive git command line that will constantly ask for new gitp+ flows or git commands.\n",
		"To use it simply launch this program using 'terminal', '--terminal' or '-terminal' as the only argument.\n",
		"To escape just insert a blank line\n\n",
		t.dir.Paint(" Git Help "), t.out.Paint(" ", output, "\n"),
	))
}

func ThemeWrapper[cSet brush.ColorType](t theme[cSet]) *tui[cSet] {
	return &tui[cSet]{theme: t}
}

func (t tui[cSet]) Cursor() {
	t.printCursor(t.branch, t.directory)
}

func (t tui[cSet]) ShowError(err any) {
	t.printError(err)
}