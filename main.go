package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"

	cloudevents "github.com/cloudevents/sdk-go/v2"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: ./iodocument-gen '<command>'")
		os.Exit(1)
	}

	cmdStr := os.Args[1]
	stdout, err := execCmd(cmdStr)
	if err != nil {
		fmt.Printf("Error executing command: %v\n", err)
		os.Exit(1)
	}

	event := cloudevents.NewEvent()
	event.SetSource("example/uri")
	event.SetType("example.type")
	err = event.SetData(cloudevents.ApplicationJSON, map[string]string{"stdin": cmdStr, "stdout": stdout})
	if err != nil {
		fmt.Printf("Failed to set cloud event data: %v\n", err)
		os.Exit(1)
	}

	eventJSON, err := json.Marshal(event)
	if err != nil {
		fmt.Printf("Error marshaling Cloud Event: %v\n", err)
		os.Exit(1)
	}

	fmt.Println(string(eventJSON))
}

func execCmd(cmdStr string) (string, error) {
	var cmd *exec.Cmd

	// Определяем команду оболочки в зависимости от операционной системы
	if isWindows() {
		cmd = exec.Command("cmd.exe", "/C", cmdStr)
	} else {
		cmd = exec.Command("/bin/sh", "-c", cmdStr)
	}

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		return "", fmt.Errorf("execution failed: %v and %v", stderr.String(), err)
	}

	return stdout.String(), nil
}

func isWindows() bool {
	return os.PathSeparator == '\\' && os.PathListSeparator == ';'
}
