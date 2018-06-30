package process

import (
	"time"
	"fmt"
	"strings"
)

type RunningProcess struct {
	Pid int
}

type Context struct {
	Commands chan string

	// TODO: find better sync method
	WaitForClose chan bool
}

func ProcessCommand(command string, context *Context) {
	switch strings.Trim(command, "\n") {
	case "stop":
		fmt.Println("Writing to WaitForClose")
		fmt.Println("%v", context.WaitForClose)
		context.WaitForClose <- true
		fmt.Println("Done Writing")
	}
}

func StartDaemon(context *Context) {
	for {
		//processes, _ := ps.Processes()
		//for _, p := range processes {
		//fmt.Println(p.Pid())
		//}

		select {
		case command := <-context.Commands:
			ProcessCommand(command, context)
		default:
		}

		// sleep 10 seconds
		time.Sleep(3 * time.Second)
	}
}
