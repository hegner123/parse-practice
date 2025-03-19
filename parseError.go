package main

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"
)

//2025-03-20T12:15:01Z ERROR Database connection failed: timeout error (db=users, retry=3)

type errorLog struct {
	Timestamp time.Time `json:"timestamp"`
	Level     logLevel  `json:"level"`
	Message   string    `json:"message"`
	DB        string    `json:"db"`
	Retry     int64     `json:"retry"`
}

func (e errorLog) String() string {
	var entry strings.Builder
	entry.WriteString(fmt.Sprintln("Time:", e.Timestamp))
	entry.WriteString(fmt.Sprintf("Message: %s\n", e.Message))
	entry.WriteString(fmt.Sprintf("Level: %s\n", e.Level))
	entry.WriteString(fmt.Sprintf("DB: %s\n", e.DB))
	entry.WriteString(fmt.Sprintf("Retry: %d\n\n", e.Retry))
	return entry.String()
}

func (e errorLog) JSON() {
	js, err := json.Marshal(e)
    if err!=nil {
        log.Fatal(err)
    }
    fmt.Println(string(js))
}

func parseErrorLog(file string) (*[]Log, error) {
	var report []Log
	for line := range strings.Lines(file) {
		if line == "\n" {
			continue
		}
		keyValue := make(map[string]string)
		kvStart := strings.Index(line, "(")
		if kvStart == -1 {
			return nil, fmt.Errorf("line %v, %v\n", line, "kvStart couldn't be found")
		}
		kvEnd := strings.Index(line, ")")
		if kvEnd == -1 {
			return nil, fmt.Errorf("line %v, %v\n", line, "kvEnd couldn't be found")
		}
		kvPair := line[kvStart+1 : kvEnd]
		kvs := strings.SplitSeq(kvPair, ",")
		for p := range kvs {
			token := strings.TrimSpace(p)
			kv := strings.Split(token, "=")
			keyValue[kv[0]] = kv[1]
		}

		tokens := strings.Split(line, " ")
		var entry errorLog = errorLog{}
		entry.DB = keyValue["db"]
		i, err := strconv.ParseInt(keyValue["retry"], 10, 64)
		if err != nil {
			return nil, err
		}
		entry.Retry = i
		var msg []string
		for i, t := range tokens {
			t := strings.TrimSpace(t)
			if i == 0 {
				if t == "" {
					continue
				}
				ts, err := time.Parse("2006-01-02T15:04:03Z", t)
				if err != nil {
					return nil, err
				}
				entry.Timestamp = ts
				continue
			}
			if i == 1 {
				entry.Level = logLevel(t)
				continue
			}
			if !strings.Contains(t, "(") && !strings.Contains(t, ")") {
				msg = append(msg, t)
			}
		}
		entry.Message = strings.Join(msg, " ")
		report = append(report, entry)
	}
	return &report, nil
}
