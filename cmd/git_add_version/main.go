package main

import (
	"os/exec"
)

func main() {
	exec.Command("git", "add", "*").CombinedOutput()
	exec.Command("git", "commit", "*", "-m", "'++'").CombinedOutput()
}