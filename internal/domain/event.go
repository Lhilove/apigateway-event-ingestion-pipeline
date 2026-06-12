package domain

import "time"

type Event struct {
	Source    string                 `json:"source" binding:"required"`
	Host      string                 `json:"host" binding:"required"`
	EventType string                 `json:"event_type" binding:"required"`
	Severity  string                 `json:"severity" binding:"required"`
	Timestamp time.Time              `json:"timestamp" binding:"required"`
	Message   string                 `json:"message"`
	Metadata  map[string]interface{} `json:"metadata"`
}
