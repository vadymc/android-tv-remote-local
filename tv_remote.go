package main

import (
	"fmt"
	"log"
	"os/exec"
	"strconv"
	"strings"
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
	inputParams := strings.Split(command, " ")
	var response string
	switch inputParams[0] {
	case "VOLUME_UP":
		response = executeKeyPressNTimes("24", getKeyPressCount(inputParams))
	case "VOLUME_DOWN":
		response = executeKeyPressNTimes("25", getKeyPressCount(inputParams))
	case "MUTE":
		response = executeKeyPressNTimes("25", 15)
	case "POWER":
		response = executeKeyPress("26")
	case "SET_VOLUME":
		response = setVolumeToLevel(parseString(inputParams[1]))
	}
	if response != "" {
		fmt.Printf("Finished command [%v] execution. response=%v\n", command, response)
		return response
	}
	return fmt.Sprintf("Command %v is not supported", command)
}

func setVolumeToLevel(volumeLevel int) string {
	executeKeyPressNTimes("25", 15)                 //mute
	return executeKeyPressNTimes("24", volumeLevel) //raise volume to requested level
}

func executeKeyPressNTimes(keyCode string, n int) string {
	adbParam := keyCode
	for i := 0; i < n; i++ {
		adbParam += " " + keyCode
	}
	return executeKeyPress(adbParam)
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
		if retryCount > 4 {
			return true, "Exhausted retry attempts"
		}
		retryCount++
		err := executeAdbCommand(keyCode)
		if err != nil {
			return false, err.Error()
		}
		return true, "Success"
	}
}

func executeAdbCommand(keyCode string) error {
	cmd := exec.Command("adb", "shell", "input keyevent", keyCode)
	return cmd.Run()
}

func reconnect() {
	disconnectCmd := exec.Command("adb", "disconnect")
	disconnectCmd.Run()

	connectCmd := exec.Command("adb", "connect", androidTvAddress)
	if err := connectCmd.Run(); err != nil {
		fmt.Printf("Failed to connect to Android TV [%v]\n", err.Error())
	} else {
		fmt.Println("Connected to Android TV")
	}
}

func getKeyPressCount(inputParams []string) int {
	keyPressCount := 1
	if len(inputParams) > 1 {
		keyPressCount = parseString(inputParams[1])
	}
	return keyPressCount
}

func parseString(s string) int {
	if s == "" {
		return 0
	}
	i, err := strconv.Atoi(s)
	if err != nil {
		log.Fatalf("Failed to parse string [%v] to int. Error %v", s, err)
		return 0
	}
	return i
}
