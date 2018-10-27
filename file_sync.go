package main

import (
	"strings"
	"bufio"
	"io"
	"fmt"
	"sinago/Error"
	"sinago/fswatch"
	"filesync/boot"
)

func main() {

	var end       chan bool
	var watchList []string

	var (
		viper     = boot.InitConfig()
		rsync     = boot.InitRsync()
		hk        = boot.InitKeeper()
	)
	boot.InitWatcher()

	for _, v := range viper.GetStringSlice("Rsync.pathList") {
		watchList = append(watchList, viper.GetString("Rsync.baseDir")+v)
	}

	var cmd = fswatch.NewFswatcher(watchList...)

	hk.Start()
	rsync.StartWork()

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		Error.CheckErr(err)
	}
	cmd.Start()
	reader := bufio.NewReader(stdout)
	fmt.Println("File watcher inited and start successful")
	for {
		line, err2 := reader.ReadString('\n')
		if err2 != nil || io.EOF == err2 {
			break
		}
		go func(inline string) {
			lineArr := strings.Split(inline, " ") // has an empty
			inline = lineArr[5]
			lineArr = strings.Split(inline, "/")
			inline = strings.Join(lineArr[3:], "/")
			aliveIPs := hk.List(viper.GetInt("HealthKeeper.status.ready"))
			for k, _ := range aliveIPs {
				rsync.Collect(rsync.NewTask(inline,k))
			}
		}(line)
	}
	cmd.Wait()
	<-end
}


