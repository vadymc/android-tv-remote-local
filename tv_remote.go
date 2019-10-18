package main

import (
	"fmt"
	"os/exec"
	"time"
)

var androidTvAddress string

func initRemoteTvConnection(remoteAddress string) {
	if remoteAddress == "" {
		panic("Specify Android TV address and port")
	}
	androidTvAddress = remoteAddress
	reconnect()
}

func executeLiteralCommand(command string) string {
	var response string
	switch command {
	case "VOLUME_UP":
		response = executeKeyPress("24")
	case "VOLUME_DOWN":
		response = executeKeyPress("25")
	case "MUTE":
		response = executeKeyPressNTimes("25", 15)
	}
	if response != "" {
		fmt.Printf("Finished command [%v] execution. response=%v\n", command, response)
		return response
	}
	return fmt.Sprintf("Command %v is not supported", command)
}

func executeKeyPressNTimes(keyCode string, n int) string {
	var msg string
	for i := 0; i < n; i++ {
		msg = executeKeyPress(keyCode)
	}
	return msg
}

func executeKeyPress(keyCode string) string {
	executeFunc := runCommandFunction(keyCode)
	isDone, msg := executeFunc()
	for !isDone {
		fmt.Printf("Execution failed with message [%v]\n", msg)
		reconnect()
		time.Sleep(500 * time.Millisecond)
		isDone, msg = executeFunc()
	}
	return msg
}

func runCommandFunction(keyCode string) func() (bool, string) {
	retryCount := 0
	return func() (bool, string) {
		if retryCount > 1 {
			return true, "Exhausted retry attempts"
		}
		retryCount++
		err := executeCommand(keyCode)
		if err != nil {
			return false, err.Error()
		}
		return true, "Success"
	}
}

func executeCommand(keyCode string) error {
	cmd := exec.Command("adb", "shell", "input", "keyevent", keyCode)
	return cmd.Run()
}

func reconnect() {
	disconnectCmd := exec.Command("adb", "disconnect")
	disconnectCmd.Run()

	connectCmd := exec.Command("adb", "connect", androidTvAddress)
	if err := connectCmd.Run(); err != nil {
		fmt.Printf("Failed to connect to Android TV [%v]\n", err.Error())
	} else {
		fmt.Println("Connected to Android TV\n")
	}
}
