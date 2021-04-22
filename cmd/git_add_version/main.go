package main

import (
	"os/exec"
	"strings"
	"fmt"
	"github.com/cheneylew/gotools/tool"
)

func main() {
	exec.Command("git", "add", "*").CombinedOutput()
	exec.Command("git", "commit", "*", "-m", "++").CombinedOutput()
	exec.Command("git", "push").CombinedOutput()

	//计算新tag
	tagOut, _ := exec.Command("git", "tag").CombinedOutput()
	version := strings.Split(strings.Trim(string(tagOut), " \n"), "\n")
	maxVersion := version[len(version)-1]
	maxVersion = strings.Trim(maxVersion, "v")
	vs := tool.ArrMap(func(a string) int {return tool.ToInt(a)}, strings.Split(maxVersion, ".")).([]int)
	for i:=len(vs)-1; i>=0; i-- {
		if vs[i] == 9 && vs[i-1] == 9 {
			vs[i] = 0
			vs[i-1] = 0
			vs[i-2] += 1
			break
		} else if vs[i] == 9 {
			vs[i] = 0
			vs[i-1]+=1
			break
		} else {
			vs[i] += 1
			break
		}
	}
	newVersion := fmt.Sprintf("v%s",strings.Join(tool.ArrMap(func(a int) string{return fmt.Sprintf("%v", a)},vs).([]string), "."))
	fmt.Println("新版本为:",newVersion)
	exec.Command("git", "tag", newVersion).CombinedOutput()
	exec.Command("git", "push").CombinedOutput()
	exec.Command("git", "push", "origin", "--tags").CombinedOutput()
}