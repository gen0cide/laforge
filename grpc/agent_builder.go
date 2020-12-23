package main

import (
	"fmt"
	"log"
	"os/exec"
)

func main(){
	agentID := "12"
	serverAddress := "18.232.181.162:50051"
	binaryName := "testbinary"

	command:="go build -ldflags=\" -X 'main.clientID="+agentID+"' -X 'main.address="+serverAddress+"'\" -o "+binaryName+" github.com/gen0cide/laforge/grpc/agent"
	cmd := exec.Command("bash","-c",command)
	stdoutStderr, err := cmd.CombinedOutput()
	cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s\n", stdoutStderr)
}
