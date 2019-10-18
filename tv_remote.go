package main

import (
	"fmt"
	"os/exec"
	"time"
)

const (
	disconnectedError = "exit status 255"
	androidTvAddress  = "192.168.1.99:5555"
)

func executeKeyPress(keyCode *string) string {
	executeFunc := runCommandFunction(*keyCode)
	isDone, msg := executeFunc()
	for !isDone {
		fmt.Printf("Execution failed with message [%v]\n", msg)
		if msg == disconnectedError {
			reconnect()
			time.Sleep(500 * time.Millisecond)
		}
		isDone, msg = executeFunc()
	}
	fmt.Printf("Finished command [%v] execution isSuccess=%v msg=[%v]\n", *keyCode, isDone, msg)
	return msg
}

func runCommandFunction(keyCode string) func() (bool, string) {
	retryCount := 0
	return func() (bool, string) {
		if retryCount > 1 {
			return true, "Exhausted retry attempts"
		}
		retryCount++
		err := executeKeyEvent(keyCode)
		if err != nil {
			return false, err.Error()
		}
		return true, "Success"
	}
}

func executeKeyEvent(keyCode string) error {
	cmd := exec.Command("adb", "shell", "input", "keyevent", keyCode)
	return cmd.Run()
}

func reconnect() {
	disconnectCmd := exec.Command("adb", "disconnect")
	disconnectCmd.Run()

	connectCmd := exec.Command("adb", "connect", androidTvAddress)
	if err := connectCmd.Run(); err != nil {
		fmt.Printf("Failed to connect to Android TV [%v]", err.Error())
	} else {
		fmt.Println("Connected to Android TV")
	}
}
