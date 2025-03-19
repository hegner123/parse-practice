package main

import (
	"fmt"
	"log"
	"os"
	"strings"
)

type fileType string
type logLevel string

const (
	Error fileType = "error"
	Debug fileType = "debug"
	Info  fileType = "info"
	ERROR logLevel = "ERROR"
	INFO  logLevel = "INFO"
	DEBUG logLevel = "DEBUG"
)

type Log interface {
	String() string
}

func main() {
	var fileName string
	if len(os.Args) > 1 {
		fileName = "log/" + os.Args[1]
	} else {
		fileName = "log/error.log"
	}
	file, err := os.ReadFile(fileName)
	if err != nil {
		fmt.Printf("couldn't read %s\n%v\n", fileName, err)
		os.Exit(1)
	}

	name := strings.TrimSuffix(fileName, ".log")
    n := strings.TrimPrefix(name, "log/")
	var report *[]Log
	switch fileType(n) {
	case Error:
		report, err = parseErrorLog(string(file))
		if err != nil {
			log.Fatal(err)
			return
		}
	case Info:
		report, err = parseInfoLog(string(file))
		if err != nil {
			log.Fatal(err)
			return
		}
	case Debug:
		report, err = parseDebugLog(string(file))
		if err != nil {
			log.Fatal(err)
			return
		}
	}
    reportSlice := *report

	for _, r := range reportSlice {
		fmt.Println(r.String())
	}
}
