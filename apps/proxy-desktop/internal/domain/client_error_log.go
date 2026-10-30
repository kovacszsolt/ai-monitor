package domain

import "time"

type ClientErrorLogEntry struct {
	ID        string
	Level     string
	Source    string
	Message   string
	Detail    string
	CreatedAt time.Time
}
