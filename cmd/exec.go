/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"encoding/json"
	"fmt"
	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/spf13/cobra"
	"os"
	"os/exec"
)

// execCmd represents the exec command
var execCmd = &cobra.Command{
	Use:   "exec",
	Short: "Выполняет команду",
	Long:  `Выполняет команду и выводит её результат`,
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		executeCommand(args)
	},
}

func executeCommand(args []string) {
	commandString := args[0]

	shell, flag := getShellAndFlag(os.PathSeparator)

	command := exec.Command(shell, flag, commandString)

	output, err := command.CombinedOutput()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to execute command: %v\n", err)
		return
	}

	command.Stdout = os.Stdout
	command.Stderr = os.Stderr

	event := cloudevents.NewEvent()
	event.SetData(cloudevents.ApplicationJSON, map[string]string{"stdin": commandString, "stdout": string(output)})

	bytes, err := json.Marshal(event)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to marshal Cloud Event")
		return
	}

	fmt.Println(string(bytes))
}

func getShellAndFlag(separator rune) (string, string) {
	var shell, flag string
	if separator == '\\' {
		shell = "cmd"
		flag = "/C"
	} else {
		shell = "sh"
		flag = "-c"
	}

	return shell, flag
}

func init() {
	rootCmd.AddCommand(execCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// execCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// execCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
