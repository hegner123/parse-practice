package main

import (
	"flag"
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
	JSON()
}

var FILE *string = flag.String("f", "error.log", "specifiy file")
var JSON *bool = flag.Bool("json", false, "output json")

func main() {
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Usage: parselog [options]\n")
		flag.PrintDefaults()
	}
	flag.Parse()


	var fileName string
	fileName = "log/error.log"
	if FILE != nil {
		fileName = "log/" + *FILE
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
		if *JSON {
			r.JSON()
		} else {
			fmt.Println(r.String())
		}
	}
}

