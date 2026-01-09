package utils

import (
	"os"
	"os/exec"
	"runtime"
)

func CallClear() {
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		// Windows: use "cmd /c cls"
		cmd = exec.Command("cmd", "/c", "cls")
	} else {
		// Unix-like systems (Linux, macOS, etc.): use "clear"
		cmd = exec.Command("clear")
	}
	cmd.Stdout = os.Stdout
	cmd.Run()
}
