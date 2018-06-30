package process

import (
	"github.com/alexpana/udeploy/models"
	"os/exec"
	"os"
	"log"
)



func Spawn(config models.Config) {
	cmd := exec.Command("bash", "testdata\\dummy.sh")
	cmd.Stdout = os.Stdout
	err := cmd.Start()
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Just ran subprocess %d, exiting\n", cmd.Process.Pid)
	cmd.Wait()
}
