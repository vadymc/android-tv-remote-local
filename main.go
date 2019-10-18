package main

import (
	"bytes"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	go startSqs()
	initRemoteTvConnection(os.Args[1])
	startWeb()
}

func startWeb() {
	router := gin.Default()
	router.POST("/v1/events/", runCommand)
	router.Run(":11002")
}

func runCommand(c *gin.Context) {
	body := string(getBody(c))
	msg := executeKeyPress(&body)
	c.String(http.StatusOK, msg)
}

func getBody(c *gin.Context) []byte {
	buf := new(bytes.Buffer)
	buf.ReadFrom(c.Request.Body)
	return buf.Bytes()
}
