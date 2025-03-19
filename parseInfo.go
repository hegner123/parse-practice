package main

import (
	"fmt"
	"strings"
	"time"
)

//2025-03-18T15:04:05Z INFO Server started on port 8080

type infoLog struct {
	timestamp time.Time
	level     logLevel
	message   string
}

func (i infoLog) String() string {
	var entry strings.Builder
	entry.WriteString(fmt.Sprintf("Time: %s\n", i.timestamp.String()))
	entry.WriteString(fmt.Sprintf("Level: %s\n", string(i.level)))
	entry.WriteString(fmt.Sprintf("Message: %s\n", i.message))
	return entry.String()
}

func parseInfoLog(file string) (*[]Log, error) {
	var report []Log
	for line := range strings.Lines(file) {
		if line == "\n" {
			continue
		}
		var entry infoLog = infoLog{}
		var message []string
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
			if i == 1 {
				entry.level = logLevel(t)
			}
			message = append(message, t)

		}
		entry.message = strings.Join(message, " ")
		report = append(report, entry)
	}
	return &report, nil
}
