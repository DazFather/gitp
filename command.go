package main

import (
	"errors"
	"strings"
)

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
		var (
			branch     string
			preConfirm bool
		)

		switch len(args) {
		case 1:
			branch = t.branch
		case 2:
			if args[1] == "--confirm" {
				preConfirm = true
				branch = t.branch
			} else {
				branch = args[1]
			}
		case 3:
			ind := 0
			if args[1] == "--confirm" {
				ind = 2
			} else if args[2] == "--confirm" {
				ind = 1
			}
			if ind != 0 {
				preConfirm = true
				branch = args[ind]
				break
			}
			fallthrough
		default:
			return errors.New("Invald given arguments, usage: undo branch <branch-name> [-confirm]")
		}

		flowName := "undo branch " + branch
		t.printFlowStart(flowName)

		err := t.executeStash(func() error { return t.removeBranch(branch, preConfirm) })
		if err == nil {
			t.printFlowEnd(flowName)
		}
		return err
	case "merge":
		return t.execute("merge", prepend("abort", args[1:])...)
	case "stash":
		return t.execute("stash", prepend("pop", args[1:])...)
	case "upstream":
		return t.execute("branch", prepend("--unset-upstream", args[1:])...)
	case "add", "stage":
		return t.execute("restore", prepend("--staged", args[1:])...)
	}
	return errors.New("Unrecognize give argument, usage: undo [commit|branch|merge|stash|upstream]")
}

func (t *tui[cSet]) removeBranch(branch string, isConfirmed bool) error {
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
	if t.branch == branch {
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

	// Delete remote branch
	if isRemoteBranch {
		if !isConfirmed {
			if isConfirmed, err = t.confirmRemoveBranch(branch); isConfirmed {
				err = t.execute("push", "origin", "--delete", branch)
			}
		} else {
			err = t.execute("push", "origin", "--delete", branch)
		}

		if err != nil {
			return err
		}
	}

	// Delete local branch
	if err = t.execute("branch", "-rd", "origin/"+branch); err != nil {
		t.printOut(err.Error())
		err = t.execute("branch", "-D", branch)
	}

	// Cleanup removed branch from list
	if err == nil {
		err = t.execute("fetch", "--prune")
	}
	return err
}

func prepend[T any](item T, slice []T) []T {
	return append([]T{item}, slice...)
}
