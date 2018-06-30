package main

import (
	"syscall"
	"net"
	"log"
	"text/tabwriter"
	"os"
	"fmt"
	"crypto/sha512"
	"encoding/hex"
	"github.com/alexpana/udeploy/models"
	"github.com/alexpana/udeploy/process"
	"bufio"
)

func handleConnection(conn net.Conn) {
	var buff []byte
	buff = make([]byte, 1024, 1024)

	for {
		conn.Read(buff);
		log.Println(buff)
	}

}

func runDaemon() {
	syscall.Unlink("/tmp/udeploy.sock")
	log.Println("Starting udeploy daemon server")
	ln, _ := net.Listen("unix", "/tmp/udeploy.sock")

	for {
		var conn, err = ln.Accept()
		if err != nil {
			log.Println("Connection error")
		} else {
			go handleConnection(conn)
		}

	}
}

func PrintProjects(projects []models.Project) {

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)

	fmt.Fprintln(w, "ID\tCONFIG\tSTATUS\t")
	for _, project := range projects {
		fmt.Fprintln(w, project.Id[:9]+"\t"+project.ConfigPath[:40]+"..."+"\t"+fmt.Sprintf("%s", project.Status)+"\t")
	}

	w.Flush()
}

func Sum512(s string) string {
	sum512 := sha512.Sum512([]byte(s))
	return hex.EncodeToString(sum512[:64])
}

func processInput(context *process.Context) {
	reader := bufio.NewReader(os.Stdin)
	const prompt = ":> "

	for {
		fmt.Print(prompt)
		input, _ := reader.ReadString('\n')
		context.Commands <- input
	}
}

func main() {

	//projects := []models.Project{
	//	{
	//		Id:         Sum512("1"),
	//		ConfigPath: "/Users/apana/go/src/github.com/apana/udeploy/deploy.yml",
	//		Status:     models.RUNNING,
	//	},
	//	{
	//		Id:         Sum512("2"),
	//		ConfigPath: "/Users/apana/go/src/github.com/apana/udeploy/deploy.yml",
	//		Status:     models.RUNNING,
	//	},
	//	{
	//		Id:         Sum512("3"),
	//		ConfigPath: "/Users/apana/go/src/github.com/apana/udeploy/deploy.yml",
	//		Status:     models.STOPPED,
	//	},
	//}
	//
	//config := models.ReadConfig("/Users/apana/go/src/github.com/apana/udeploy/deploy.yml")
	//fmt.Printf("%v\n", config)
	//
	//PrintProjects(projects)

	//process.Spawn(models.Config{})

	daemonContext := process.Context{
		Commands:     make(chan string, 10),
		//WaitForClose: make(chan bool),
	}

	go process.StartDaemon(&daemonContext)
	go processInput(&daemonContext)

	fmt.Println("Waiting")
	<-daemonContext.WaitForClose
	fmt.Println("Done")
}
