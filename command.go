package main

import (
	"errors"
	"strings"
)

func (t tui[xSet]) executeFlow(flowName string, stash bool, commands [][]string) error {
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

func (t *tui[cSet]) help(command string, args ...string) error {
	out, err := git(command, args...)
	if err == nil {
		t.printHelp(out)
	}
	return err
}

func (t *tui[cSet]) update(command string, args ...string) error {
	return t.executeFlow(command, true, [][]string{{"fetch"}, {"pull"}})
}

func (t *tui[cSet]) fork(command string, args ...string) error {
	if len(args) != 1 {
		return errors.New("Invalid given argument, usage: fork <new-branch-name>")
	}

	err := t.executeFlow(command, true, [][]string{
		{"fetch"},
		{"pull"},
		{"checkout", "-b", args[0]},
		{"push", "--set-upstream", "origin", args[0]},
	})

	if err == nil {
		t.branch = args[0]
	}
	return err
}

func (t *tui[cSet]) undo(command string, args ...string) error {
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
	return errors.New("Unrecognize give argument, usage: undo [commit|branch|merge|stash|upstream]")
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

func prepend[T any](item T, slice []T) []T {
	return append([]T{item}, slice...)
}
