# Git+
[![Go Report Card](https://goreportcard.com/badge/github.com/DazFather/gitp)](https://goreportcard.com/report/github.com/DazFather/gitp) 


A cli tool to enhance your git experience

## Features
> Why even bother with another git thingy?

If you are looking for something very hight-level, this is not for you. 
But if you are looking for some handy utilities and a bit of colors in your terminal you're welcome.

### 📟 Interactive terminal for lazy people
Tired of writing "git" all the time before your commands? 
Now you can use `gitp terminal` to keep sending git commands without the need to type "git" all the time.
And before you ask: to exit simply enter an empty input

### 🎁 Pre-made flows and commands
- 🔄 **`update`**: Update your branch with possible incoming remote changes
  > Equivalent of:
  > ```shell
  > git status
  > git stash     # If there are changes
  > git fetch
  > git pull
  > git stash pop # If there were changes
  > ```
- ⤴️ **`fork <branch-name>`**: forks your current branch and creates a new one with given name and sets remote upstream
  > Equivalent of:
  > ```shell
  > git status
  > git stash     # If there are changes
  > git fetch
  > git pull
  > git checkout -b <branch>
  > git push --set-upstream origin <branch>
  > git stash pop # If there were changes
  > ```
- 🛟 **`undo [commit|branch|merge|stash|upstream] <args...>`**: The undo-button you wish you had earlier, has different effects depending on the input
  - commit (`reset HEAD~1 <args...>`): reset last commit preserving changes locally by default
  - merge (`merge abort <args...>`): abort current merge
  - stash (`stash pop <args...>`): reapply last stashed item and remove it from the stack
  - upstream (`--unset-upstream <args...>`): disable remote tracking from a branch
  - branch: it remove given branch or if missing the current one, it deletes also from remote

### ✨ Cool looking
_(And soon customizable)_ Nothing too fancy, is still a terminal application after all, 
but it uses the [brush](https://github.com/DazFather/brush) library to bring back 
some colors on your boring terminal to help you guide on the current branch, 
project or command in use. Try it and see for yourself

### 🪟 Transparent I/O
All commands that are being fed to git with of course the related output will be shown on screen,
 so you know what is happening all the time


## Usage
Siply use it as you normally would use git.
 > `gitp <command> <agruments...>` 

If you want to keep input new commands, you can use the terminal feature. 
Then you can directly input new commands without writing `git` or `gitp` all the time 
 > `gitp terminal` or `gitp --terminal`

Enjoy!
