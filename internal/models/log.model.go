package models

import "time"

type CommonLog struct {
	TimeStamp   time.Time `json:"timestamp"`
	Service     string    `json:"service"`
	RequestId   string    `json:"request_id"`
	LoggingType string    `json:"logging_type"`
	LogLevel    string    `json:"log_level"`
	Message     string    `json:"message"`
}

type RequestLog struct {
	Timestamp   time.Time     `json:"timestamp"`
	Service     string        `json:"service"`
	RequestId   string        `json:"request_id"`
	LoggingType string        `json:"logging_type"`
	Hostname    string        `json:"hostname"`
	Uri         string        `json:"uri"`
	Method      string        `json:"method"`
	Path        string        `json:"path"`
	Start       time.Time     `json:"start"`
	End         time.Time     `json:"end"`
	Headers     string        `json:"headers"`
	Request     string        `json:"request"`
	Response    string        `json:"response"`
	Status      int           `json:"status"`
	Latency     time.Duration `json:"latency"`
}
