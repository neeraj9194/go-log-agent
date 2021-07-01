package src

import (
	"os"
)

type LogStruct struct {
	// common attributes
	Host      string    `json:"host"`
	Message   string    `json:"message"`
	Service   string    `json:"service"`
}


func ParseLog(service string, log string) (LogStruct, error) {
	hostname, err := os.Hostname()
	return LogStruct{
		hostname,
		log,
		service,
	}, err
}
