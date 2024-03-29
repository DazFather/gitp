package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
)

type tui struct {
	branch, directory string
	theme
}

func NewTui(t theme) (term *tui, err error) {
	term = &tui{theme: t}

	if err = term.refreshBranch(); err != nil {
		term.branch = "[no git repo]"
		t.printWarning("Cannot find git repositoy for this project, use 'init' or 'clone' to create a new one")
		if term.directory, err = os.Getwd(); err != nil {
			return nil, err
		}
	} else {
		term.directory, err = git("rev-parse", "--show-toplevel")
	}

	if err == nil {
		term.directory = filepath.Base(term.directory)
	}

	return
}

func DefaultTui() (tui, error) {
	term, err := NewTui(defaultTheme)
	if term == nil {
		panic(fmt.Sprint("Cannot create gitp+ terminal: ", err))
	}

	return *term, err
}

func (t *tui) Gitp(command string, args ...string) (err error) {
	switch command {
	case "git", "gitp", "git+":
		if len(args) > 0 {
			return t.Gitp(args[0], args[1:]...)
		}
		t.printError("No given arguments, use 'gitp help' to learn more'")
	case "help", "-h", "--help":
		err = t.help(command, args...)
	case "undo":
		err = t.undo(command, args...)
	case "update":
		err = t.update(command, args...)
	case "fork":
		err = t.fork(command, args...)
	case "align":
		err = t.align(command, args...)
	case "init", "clone":
		if err = t.execute(command, args...); err == nil {
			if err = t.refreshBranch(); err == nil {
				t.directory, err = git("rev-parse", "--show-toplevel")
			}
		}
	case "checkout", "switch":
		if err = t.execute(command, args...); err == nil {
			err = t.refreshBranch()
		}
	default:
		err = t.execute(command, args...)
	}
	return
}

func (t tui) InteractiveGitp(escapeSeq string, keepAlive bool) error {
	var (
		rgx     = regexp.MustCompile(`".+"|[^\s]+`)
		scanner = bufio.NewScanner(os.Stdin)
	)

	t.Cursor("")
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == escapeSeq {
			return nil
		}

		args := rgx.FindAllString(line, -1)
		if err := t.Gitp(args[0], args[1:]...); err != nil {
			if !keepAlive {
				return err
			}
			t.ShowError(err)
		}
		t.Cursor("")
	}

	return scanner.Err()
}

func (t tui) executeStash(exec func() error) error {
	t.printCommand(t.branch, "status", "--p=v1")
	out, err := git("status", "--p=v1")
	t.printOut(out)

	if err != nil {
		return err
	}

	if out == "" { // Nothing to stash
		err = exec()
	} else if err = t.execute("stash"); err == nil {
		if err = exec(); err == nil {
			err = t.execute("stash", "pop")
		}
	}

	return err
}

func (t tui) executeAll(commands [][]string) (err error) {
	for _, command := range commands {
		switch len(command) {
		case 0:
			continue
		case 1:
			err = t.execute(command[0])
		default:
			err = t.execute(command[0], command[1:]...)
		}

		if err != nil {
			break
		}
	}
	return
}

func (t tui) executeFlow(flowName string, stash bool, commands [][]string) error {
	var err error

	t.printFlowStart(flowName)

	if stash {
		err = t.executeStash(func() error { return t.executeAll(commands) })
	} else {
		err = t.executeAll(commands)
	}

	if err == nil {
		t.printFlowEnd(flowName)
	}
	return err
}

func (t tui) execute(command string, args ...string) error {
	t.printCommand(t.branch, command, args...)
	out, err := git(command, args...)
	if err == nil {
		t.printOut(out)
	}

	return err
}

func (t *tui) refreshBranch() (err error) {
	if t.branch, err = git("branch", "--show-current"); err != nil {
		return
	}

	if t.branch == "" {
		if hash, err := git("rev-parse", "--short", "HEAD"); err == nil {
			t.branch = hash + " [detached]"
		}
	}
	return
}

func git(command string, args ...string) (stdout string, err error) {
	cmdArgs := make([]string, len(args)+1)
	cmdArgs[0] = command
	for i := range args {
		cmdArgs[i+1] = args[i]
	}

	out, err := exec.Command("git", cmdArgs...).CombinedOutput()
	if out != nil {
		stdout = strings.TrimSpace(string(out))
		if err != nil {
			err = wrapError("["+strings.TrimPrefix(err.Error(), "exit status ")+"] "+stdout, err)
		}
	}
	return
}
