package main

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"
)

//2025-03-18T15:06:12Z DEBUG HTTP request received: method=GET, path=/api/items, status=200, duration=150ms

type DebugMethod string

const (
	DEBUGPOST   DebugMethod = "POST"
	DEBUGGET    DebugMethod = "GET"
	DEBUGPUT    DebugMethod = "PUT"
	DEBUGDELETE DebugMethod = "DELETE"
)

type debugLog struct {
	Timestamp time.Time     `json:"timestamp"`
	Level     logLevel      `json:"level"`
	Method    DebugMethod   `json:"method"`
	Path      string        `json:"path"`
	Status    int64         `json:"status"`
	Duration  time.Duration `json:"duration"`
}

func (d debugLog) String() string {
	var entry strings.Builder
	entry.WriteString(fmt.Sprintf("Time: %s\n", d.Timestamp.String()))
	entry.WriteString(fmt.Sprintf("Level: %s\n", string(d.Level)))
	entry.WriteString(fmt.Sprintf("Method: %s\n", d.Method))
	entry.WriteString(fmt.Sprintf("Path: %s\n", d.Path))
	entry.WriteString(fmt.Sprintf("Status: %d\n", d.Status))
	entry.WriteString(fmt.Sprintf("Duration: %s\n", d.Duration.String()))
	return entry.String()
}

func (d debugLog) JSON() {
	js, err := json.Marshal(d)
    if err!=nil {
        log.Fatal(err)
    }
    fmt.Println(string(js))
}

func parseDebugLog(file string) (*[]Log, error) {
	var report []Log
	for line := range strings.Lines(file) {
		if line == "\n" {
			continue
		}
		var entry debugLog = debugLog{}
		kvI := strings.Index(line, ": ")
		if kvI == -1 {
			return nil, fmt.Errorf("Couldn't parse log for keyvalues")
		}
		splitKv := strings.Split(line[kvI+1:], ",")
		kvP := make(map[string]string)
		for _, kv := range splitKv {
			kvPair := strings.Split(kv, "=")
			tPair1 := strings.TrimSpace(kvPair[0])
			tPair2 := strings.TrimSpace(kvPair[1])
			kvP[tPair1] = tPair2
		}
		entry.Method = DebugMethod(kvP["method"])
		entry.Path = kvP["path"]
		s := kvP["status"]
		statusI, err := strconv.ParseInt(s, 10, 64)
		if err != nil {
			return nil, err
		}
		entry.Status = statusI
		td := kvP["duration"]
		dur, err := time.ParseDuration(td)
		if err != nil {
			return nil, err
		}
		entry.Duration = dur
		tokens := strings.Split(line, " ")
		for i, tk := range tokens {
			t := strings.TrimSpace(tk)
			if i == 0 {
				ts, err := time.Parse("2006-01-02T15:04:03Z", t)
				if err != nil {
					return nil, err
				}
				entry.Timestamp = ts
				continue
			}
			if logLevel(t) == INFO || logLevel(t) == ERROR || logLevel(t) == DEBUG {

				entry.Level = logLevel(t)
			}
		}
		report = append(report, entry)
	}
	return &report, nil
}
