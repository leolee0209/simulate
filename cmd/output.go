package cmd

import (
	"runtime"
	"os"
	"os/exec"
)

func CallClear() {
    // Cross-platform "hack" using system commands
	value, ok := map[string]string{
		"windows": "cls",
		"linux":   "clear",
		"darwin":  "clear",
	}[runtime.GOOS]
    
	if ok {
		cmd := exec.Command(value)
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
}

