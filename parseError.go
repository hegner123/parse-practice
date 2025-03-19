package main

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

//2025-03-20T12:15:01Z ERROR Database connection failed: timeout error (db=users, retry=3)

type errorLog struct {
	timestamp time.Time
	level     logLevel
	message   string
	db        string
	retry     int64
}

func (e errorLog) String() string {
	var entry strings.Builder
	entry.WriteString(fmt.Sprintln("Time:", e.timestamp))
	entry.WriteString(fmt.Sprintf("Message: %s\n", e.message))
	entry.WriteString(fmt.Sprintf("Level: %s\n", e.level))
	entry.WriteString(fmt.Sprintf("DB: %s\n", e.db))
	entry.WriteString(fmt.Sprintf("Retry: %d\n\n", e.retry))
	return entry.String()
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
		entry.db = keyValue["db"]
		i, err := strconv.ParseInt(keyValue["retry"], 10, 64)
		if err != nil {
			return nil, err
		}
		entry.retry = i
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
				entry.timestamp = ts
				continue
			}
			if i == 1 {
				entry.level = logLevel(t)
				continue
			}
			if !strings.Contains(t, "(") && !strings.Contains(t, ")") {
				msg = append(msg, t)
			}
		}
		entry.message = strings.Join(msg, " ")
		report = append(report, entry)
	}
	return &report, nil
}
