package main

import (
	"github.com/alexpana/udeploy/udeployd/process"
	"time"
)

func main() {
	daemonContext := process.NewContext(process.Settings{
		PollInterval: 5 * time.Second,
	})

	process.StartDaemon(daemonContext)

	//noinspection GoInfiniteFor
	for {
		time.Sleep(100 * time.Second)
	}
}
