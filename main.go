package main

import (
	"bytes"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	initRemoteTvConnection(os.Args[1])
	go startSqs()
	startRest()
}

func startRest() {
	router := gin.Default()
	router.POST("/v1/events/", runCommand)
	router.Run(":11002")
}

func runCommand(c *gin.Context) {
	bodyBytes := getBody(c)
	body := string(bodyBytes)
	msg := executeLiteralCommand(body)
	c.String(http.StatusOK, msg)
}

func getBody(c *gin.Context) []byte {
	buf := new(bytes.Buffer)
	buf.ReadFrom(c.Request.Body)
	return buf.Bytes()
}
