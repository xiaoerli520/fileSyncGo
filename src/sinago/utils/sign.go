package utils

import (
	"os"
	"os/signal"
	"fmt"
)

func SignalListen(osSignal os.Signal, outStr string, callback func()) {
	c := make(chan os.Signal)
	signal.Notify(c, osSignal)
	for {
		s := <-c
		callback()
		fmt.Println(outStr, s)
	}
}

