package lib

import (
	"fmt"
	"os/exec"

	"github.com/fatih/color"
)

func Play(title string, refr string, videoLink string) {
	color.Magenta("Playing...")
	cmd := exec.Command("vlc", "--meta-title", title, "--meta-url", refr, videoLink)
	err := cmd.Start()
	if err != nil {
		panic(err)
	}
	fmt.Println("PID:", cmd.Process.Pid)
}
