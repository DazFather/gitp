package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/DazFather/brush"
)

type tui[cSet brush.ColorType] struct {
	branch, directory string
	theme[cSet]
}

func NewTui[cSet brush.ColorType](t theme[cSet]) (term *tui[cSet]) {
	term = &tui[cSet]{theme: t}
	if err := term.refreshInfo(); err != nil {
		t.printError(err)
		return nil
	}
	return
}

func DefaultTui() *tui[brush.ANSIColor] {
	return NewTui(defaultTheme)
}

func (t *tui[cSet]) Gitp(command string, args ...string) error {
	switch command {
	case "git":
		return t.Gitp(args[0], args[1:]...)
	case "-h", "--help", "help":
		out, err := git(command, args...)
		if err == nil {
			t.printHelp(out)
		}
		return err
	case "undo":
		if len(args) == 0 {
			return errors.New("Invalid given argument, usage: undo [commit|branch|merge|stash|upstream]")
		}

		switch args[0] {
		case "commit":
			return t.execute("reset", prepend("HEAD~1", args[1:])...)
		case "branch":
			if len(args) == 2 {
				return t.removeBranch(args[1])
			}
			return t.removeBranch("")
		case "merge":
			return t.execute("merge", prepend("abort", args[1:])...)
		case "stash":
			return t.execute("stash", prepend("pop", args[1:])...)
		case "upstream":
			return t.execute("branch", prepend("--unset-upstream", args[1:])...)
		}
	case "update":
		return t.executeFlow(command, true, flow{{"fetch"}, {"pull"}})
	case "fork":
		if len(args) != 1 {
			return errors.New("Invalid given argument, usage: fork <branch-name>")
		}

		err := t.executeFlow(command, true, flow{
			{"fetch"},
			{"pull"},
			{"checkout", "-b", args[0]},
			{"push", "--set-upstream", "origin", args[0]},
		})
		if err != nil {
			t.branch = args[0]
		}
	case "checkout", "switch":
		if len(args) == 0 {
			return errors.New("Invalid given argument: missing branch name to switch to")
		}

		if err := t.execute(command, args...); err != nil {
			return err
		}
		t.branch = args[0]
	default:
		return t.execute(command, args...)
	}
	return nil
}

func (t tui[cSet]) InteractiveGitp(escapeSeq string) error {
	var (
		rgx     = regexp.MustCompile(`".+"|[^\s]+`)
		scanner = bufio.NewScanner(os.Stdin)
	)

	t.Cursor()
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == escapeSeq {
			return nil
		}

		args := rgx.FindAllString(line, -1)
		if err := t.Gitp(args[0], args[1:]...); err != nil {
			return err
		}
		t.Cursor()
	}

	return scanner.Err()
}

func (t tui[cSet]) execute(command string, args ...string) error {
	t.printCommand(t.branch, command, args...)
	out, err := git(command, args...)
	if err == nil {
		t.printOut(out)
	}

	return err
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
			err = fmt.Errorf("[%s] %s", strings.TrimPrefix(err.Error(), "exit status "), stdout)
		}
	}
	return
}

func (t *tui[cSet]) refreshInfo() error {
	branch, err := git("branch", "--show-current")
	if err != nil {
		return err
	}

	if branch == "" {
		if branch, err = git("rev-parse", "--short", "HEAD"); err != nil {
			return err
		}
		branch += " [detached]"
	}
	t.branch = branch

	dir, err := git("rev-parse", "--show-toplevel")
	if err == nil {
		t.directory = filepath.Base(dir)
	}

	return err
}

func (t *tui[cSet]) removeBranch(branch string) error {
	var isCurrent bool
	if branch == "" {
		branch = t.branch
		isCurrent = true
	}

	flowName := "undo branch " + branch
	t.printFlowStart(flowName)

	// Check if branch to delete has been pushed on remote
	t.printCommand(t.branch, "ls-remote", "--exit-code", "--heads", "origin", branch)
	out, err := git("ls-remote", "--exit-code", "--heads", "origin", branch)
	t.printOut(out)
	if err != nil {
		if n := strings.Trim(err.Error(), "[ ]"); n != "1" && n != "2" {
			return err
		}
	}
	isRemoteBranch := out != ""

	// If the branch to delete is current checkout to a detached HEAD
	if isCurrent {
		t.printCommand(t.branch, "rev-parse", "--short", "HEAD")
		if out, err = git("rev-parse", "--short", "HEAD"); err != nil {
			return err
		}
		t.printOut(out)
		if err = t.execute("checkout", "--detach", out); err != nil {
			return err
		}
		t.branch = out + " [detached]"
	}

	// Delete branch
	if isRemoteBranch {
		err = t.execute("push", "origin", "--delete", branch)
	} else if err = t.execute("branch", "-rd", "origin/"+branch); err != nil {
		t.printOut(err.Error())
		err = t.execute("branch", "-D", branch)
	}

	// Cleanup removed branch from list
	if err == nil {
		if err = t.execute("fetch", "--prune"); err == nil {
			t.printFlowEnd(flowName)
		}
	}
	return err
}

type flow [][]string

func (t tui[xSet]) executeFlow(flowName string, stash bool, commands flow) error {
	var err error

	t.printFlowStart(flowName)
	if stash {
		t.printCommand(t.branch, "status", "--p=v1")
		out, err := git("status", "--p=v1")
		t.printOut(out)
		if stash = out != ""; stash {
			err = t.execute("stash")
		}
		if err != nil {
			return err
		}
	}

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
			return err
		}
	}

	if stash {
		err = t.execute("stash", "pop")
	}

	if err == nil {
		t.printFlowEnd(flowName)
	}
	return err
}

func prepend[T any](item T, slice []T) []T {
	return append([]T{item}, slice...)
}
