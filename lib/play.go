package lib

import (
	"fmt"
	"os/exec"

	"github.com/fatih/color"
)

func Play(title string, refr string, videoLink string) {
	color.Magenta("Playing...")
	cmd := exec.Command("iina", videoLink, "--no-stdin", "--keep-running", "--mpv-referrer="+refr, "--mpv-force-media-title="+title)
	err := cmd.Start()
	if err != nil {
		panic(err)
	}
	fmt.Println("PID:", cmd.Process.Pid)
}
