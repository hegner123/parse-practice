package main

import (
	"fmt"
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
	timestamp time.Time
	level     logLevel
	method    DebugMethod
	path      string
	status    int64
	duration  time.Duration
}

func (d debugLog) String() string {
	var entry strings.Builder
	entry.WriteString(fmt.Sprintf("Time: %s\n", d.timestamp.String()))
	entry.WriteString(fmt.Sprintf("Level: %s\n", string(d.level)))
	entry.WriteString(fmt.Sprintf("Method: %s\n", d.method))
	entry.WriteString(fmt.Sprintf("Path: %s\n", d.path))
	entry.WriteString(fmt.Sprintf("Status: %d\n", d.status))
	entry.WriteString(fmt.Sprintf("Duration: %s\n", d.duration.String()))
	return entry.String()
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
		entry.method = DebugMethod(kvP["method"])
		entry.path = kvP["path"]
		s := kvP["status"]
		statusI, err := strconv.ParseInt(s, 10, 64)
		if err != nil {
			return nil, err
		}
		entry.status = statusI
		td := kvP["duration"]
		dur, err := time.ParseDuration(td)
		if err != nil {
			return nil, err
		}
		entry.duration = dur
		tokens := strings.Split(line, " ")
		for i, tk := range tokens {
			t := strings.TrimSpace(tk)
			if i == 0 {
				ts, err := time.Parse("2006-01-02T15:04:03Z", t)
				if err != nil {
					return nil, err
				}
				entry.timestamp = ts
				continue
            }
			if logLevel(t) == INFO || logLevel(t) == ERROR || logLevel(t) == DEBUG {

				entry.level = logLevel(t)
			}
		}
		report = append(report, entry)
	}
	return &report, nil
}
