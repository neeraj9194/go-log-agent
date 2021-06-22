package src

import (
	"fmt"
	"regexp"
	"time"
	// "github.com/influxdata/go-syslog/v3/rfc3164"
)

type LogStruct struct {
	// common attributes
	Host      string    `json:"host"`
	Level     string    `json:"level"`
	Timestamp time.Time `json:"timestamp"`
	Message   string    `json:"message"`
	Service   string    `json:"service"`

	// The structure can be extended to support diffrent types of services like HTTP, DB etc.
	Http HTTP `json:"http"`
}

type HTTP struct {
	URL      string `json:"url"`
	ClientIP string `json:"client_ip"`
	Version  string `json:"version"`
}

// ===========Nginx============

func parseNginxLog(logLine string) LogStruct {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("Errror while parsing: %v", logLine)
		}
	}()
	ipaddress_reg := regexp.MustCompile(`([0-9]{1,3}\.){3}[0-9]{1,3}|(([0-9a-fA-F]{1,4}:){7,7}[0-9a-fA-F]{1,4}|([0-9a-fA-F]{1,4}:){1,7}:|([0-9a-fA-F]{1,4}:){1,6}:[0-9a-fA-F]{1,4}|([0-9a-fA-F]{1,4}:){1,5}(:[0-9a-fA-F]{1,4}){1,2}|([0-9a-fA-F]{1,4}:){1,4}(:[0-9a-fA-F]{1,4}){1,3}|([0-9a-fA-F]{1,4}:){1,3}(:[0-9a-fA-F]{1,4}){1,4}|([0-9a-fA-F]{1,4}:){1,2}(:[0-9a-fA-F]{1,4}){1,5}|[0-9a-fA-F]{1,4}:((:[0-9a-fA-F]{1,4}){1,6})|:((:[0-9a-fA-F]{1,4}){1,7}|:)|fe80:(:[0-9a-fA-F]{0,4}){0,4}%[0-9a-zA-Z]{1,}|::(ffff(:0{1,4}){0,1}:){0,1}((25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9])\.){3,3}(25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9])|([0-9a-fA-F]{1,4}:){1,4}:((25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9])\.){3,3}(25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9]))`)
	ipaddress_data := ipaddress_reg.FindAllString(logLine, -1)[0]

	datetime_reg := regexp.MustCompile(`\[\d{1,2}\/\w{3}\/\d{1,4}(:[0-9]{1,2}){3} \+([0-9]){1,4}\]`)
	datetime_data := datetime_reg.FindAllString(logLine, -1)[0]
	parsed_time, _ := time.Parse("[02/Jan/2006:15:04:05 -0700]", datetime_data)

	url_reg := regexp.MustCompile(`"\w+\s[^\s]+`)
	url_data := url_reg.FindAllString(logLine, -1)[0]

	http_version_reg := regexp.MustCompile(`HTTP\/\d.\d"`)
	http_version_data := http_version_reg.FindAllString(logLine, -1)[0]

	response_and_byte_reg := regexp.MustCompile(`([0-9]{1,3}) \d+`)
	response_and_byte_data := response_and_byte_reg.FindAllString(logLine, -1)[0]

	return LogStruct{
		"",
		"",
		parsed_time,
		response_and_byte_data,
		"nginx",
		HTTP{
			url_data,
			ipaddress_data,
			http_version_data,
		},
	}

}

func parseGenericLog(logLine string) LogStruct {
	// matches Golang like log : 2009/01/23 01:23:23 message...
	re := regexp.MustCompile(`(?P<datetime>\d{4}/\d{2}/\d{2} \d{2}:\d{2}:\d{2}) (?P<message>\w+)`)
	match := re.FindStringSubmatch(logLine)
	result := make(map[string]string)
	for i, name := range re.SubexpNames() {
		if i != 0 && name != "" {
			result[name] = match[i]
		}
	}
	parsedTime, _ := time.Parse("Jan 02 15:04:05", result["datetime"])

	return LogStruct{
		"",
		"",
		parsedTime,
		result["message"],
		"generic",
		HTTP{},
	}

}

// ===========Syslog============
func parseSyslog(logLine string) LogStruct {
	// matches Golang like log : 2009/01/23 01:23:23 message...
	re := regexp.MustCompile(`(?P<datetime>[A-Z][a-z][a-z]\s{1,2}\d{1,2}\s\d{2}[:]\d{2}[:]\d{2})\s(?P<machinename>[\w][\w\d\.@-]*)\s(?P<message>.*)`)
	match := re.FindStringSubmatch(logLine)
	result := make(map[string]string)
	for i, name := range re.SubexpNames() {
		if i != 0 && name != "" {
			result[name] = match[i]
		}
	}
	parsedTime, _ := time.Parse("Jan 02 15:04:05", result["datetime"])
	currentYear := time.Now().Year()
	// No year present in logs
	parsedTime = parsedTime.AddDate(currentYear, 0, 0)

	return LogStruct{
		"",
		"",
		parsedTime,
		result["message"],
		"generic",
		HTTP{},
	}

}

func ParseLog(service string, log string) LogStruct {

	switch service {
	case "syslog":
		return parseSyslog(log)
	case "nginx":
		return parseNginxLog(log)
	default:
		return parseGenericLog(log)
	}

}
