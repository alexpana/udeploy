package process

import (
	"github.com/alexpana/udeploy/udeployd/deployment"
	"os/exec"
	"os"
	"log"
	"path"
)

func Spawn(config deployment.Config) *exec.Cmd {
	cmd := exec.Command(config.RunCommand, config.RunArgs...)
	cmd.Stdout = os.Stdout
	cmd.Dir = path.Dir(config.Path)
	err := cmd.Start()
	if err != nil {
		log.Fatal(err)
	}

	return cmd
}

func Wait(cmd *exec.Cmd) {
	log.Printf("Running process with PID %d\n", cmd.Process.Pid)
	err := cmd.Wait()
	if err != nil {
		log.Printf("Processed %v finished with error: %s\n", cmd.Process.Pid, err.Error())
	}
}
