package process

import (
	"strings"
	"github.com/alexpana/udeploy/udeployd/deployment"
	"github.com/mitchellh/go-ps"
	"text/tabwriter"
	"fmt"
	"time"
	"syscall"
	"log"
	"net"
	"io"
	"bytes"
	"os"
	"os/user"
	"io/ioutil"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

type Settings struct {
	PollInterval time.Duration
}

type Context struct {
	Settings  Settings         // deploy udeployd settings
	Requests  chan string      // channel for receiving commands
	Replies   chan string      // channel for returning command results
	processes []RunningProcess // list of managed processes

	StateDir string
}

func NewContext(settings Settings) *Context {
	context := new(Context)
	context.Settings = settings
	context.Requests = make(chan string)
	context.processes = make([]RunningProcess, 0, 10)

	currentUser, err := user.Current()
	if err != nil {
		panic("Could not find current currentUser")
	}

	context.StateDir = currentUser.HomeDir + "/.udeploy"
	return context
}

type RunningProcess struct {
	Pid       int               // the running process id
	StartTime time.Time         // time this process has been started
	Config    deployment.Config // the deployment config
}

func handleRequest(command string, context *Context, writer io.Writer) {
	switch strings.Trim(command, "\n") {
	case "ps":
		handlePs(context, writer)
	case "register":
		handleRegister(context, writer)
	case "stop":
		handleStop(context, writer)
	}
}

func handlePs(context *Context, writer io.Writer) {
	w := tabwriter.NewWriter(writer, 0, 0, 3, ' ', 0)

	fmt.Fprintln(w, "NAME\tPID\tCONFIG\tUPTIME\t")
	for _, project := range context.processes {
		fmt.Fprintf(w,
			"%s\t%d\t%s\t%v\t\n",
			project.Config.Name,
			project.Pid,
			project.Config.Path[:60]+"...",
			time.Now().Sub(project.StartTime))
	}

	w.Flush()
}

func handleRegister(context *Context, writer io.Writer) {
	fmt.Fprintln(writer, "Not implemented")
}

func handleStop(context *Context, writer io.Writer) {
	fmt.Fprintln(writer, "Not implemented")
}

func isProcessRunning(pid int) bool {
	processes, _ := ps.Processes()
	for _, p := range processes {
		if p.Pid() == pid {
			return true
		}
	}

	return false
}

func restoreState(context *Context) {
	deploymentsFile := context.StateDir + ("/registered_deployments")

	if _, err := os.Stat(deploymentsFile); err == nil {
		dat, err := ioutil.ReadFile(deploymentsFile)
		check(err)

		deploymentConfigurations := strings.Split(string(dat), "\n")
		for _, deploymentConfig := range deploymentConfigurations {
			config, err := deployment.ReadConfig(deploymentConfig)
			if err != nil {
				log.Println("Could not find deployment file " + deploymentConfig)
			} else {
				registerDeployment(context, config)
			}
		}
	}
}
func registerDeployment(context *Context, config deployment.Config) {
	log.Println("Registering a new deployment called " + config.Name)
	context.processes = append(context.processes, RunningProcess{
		Pid:    -1,
		Config: config,
	})
}

func supervisor(context *Context) {
	log.Println("Starting process supervisor.")
	for {
		for i := range context.processes {
			process := &context.processes[i]

			if process.Pid != -1 && !isProcessRunning(process.Pid) {
				// Mark process as stopped
				process.Pid = -1
			}

			// Revive / Start processes
			if process.Pid == -1 {
				spawn := Spawn(process.Config)
				process.Pid = spawn.Process.Pid
				process.StartTime = time.Now()
				go Wait(spawn)
			}
		}

		time.Sleep(context.Settings.PollInterval)
	}
}

func requestHandler(context *Context) {
	syscall.Unlink("/tmp/udeploy.sock")
	log.Println("Starting request handler. Listening on UNIX socket /tmp/udeploy.sock")
	ln, _ := net.Listen("unix", "/tmp/udeploy.sock")

	for {
		var conn, err = ln.Accept()
		if err != nil {
			log.Println("Connection error")
		} else {
			log.Printf("New connection from %s", conn.RemoteAddr())
			go handleConnection(context, conn)
		}
	}
}

func handleConnection(context *Context, conn net.Conn) {
	inputBuffer := make([]byte, 1024, 1024)

	for {
		// clear the buffer
		for i := range inputBuffer {
			inputBuffer[i] = 0
		}

		_, err := conn.Read(inputBuffer)
		if err != nil {
			log.Printf("Connection closed %s", conn.RemoteAddr())
			return
		}

		requestString := strings.Trim(string(inputBuffer), "\x00\n")

		log.Printf("Received '%s' from %s", requestString, conn.RemoteAddr())

		outputBytes := bytes.NewBuffer(make([]byte, 1024))

		handleRequest(requestString, context, outputBytes)

		conn.Write(outputBytes.Bytes())
	}
}

func prepareSystem(context *Context) {
	err := os.MkdirAll(context.StateDir, 0700)
	if err != nil {
		log.Printf("Could not create state directory: %s\n", err)
	}
}

func StartDaemon(context *Context) {
	prepareSystem(context)
	restoreState(context)

	go supervisor(context)
	go requestHandler(context)
}
