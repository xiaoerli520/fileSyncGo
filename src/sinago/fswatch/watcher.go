package fswatch

import (
	"runtime"
	"os/exec"
	"fmt"
)

var watcherOSX = []string{"-t","-l","0.1" , "-f", "'%a %b %d %T %Y'", "-r", "-v", "-x", "--event", "IsDir", "--event", "Created", "--event", "Updated"}
var watcherLinux = []string{"-e", "modify", "-e", "create", "-e", "move",  "--format", "'%T %w%f %Xe %e '", "-m", "--exclude", "\\.([0-9a-zA-Z]{6}|tmp|sw[a-zA-Z0-9]*)$","-r", "--timefmt", "'%a %b %d %T %Y'"}


// get a cmd object to watch the directories changes
func NewFswatcher(src ...string) (cmd *exec.Cmd) {

	switch runtime.GOOS {
	case "darwin":
		commandParam := append(watcherOSX, src...)
		cmd = exec.Command("fswatch", commandParam...)
	case "linux":
		commandParam := append(watcherLinux, src...)
		cmd = exec.Command("inotifywait", commandParam...)
		fmt.Println(commandParam)
	}
	return cmd
}