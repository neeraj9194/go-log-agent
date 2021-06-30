package src

import (
	"fmt"
	"reflect"
	"testing"
	"time"
)

func TestParseNginxLog(t *testing.T) {
	testString := `93.180.71.3 - - [17/May/2015:08:05:32 +0000] "GET /downloads/product_1 HTTP/1.1" 304 0 "-" "Debian APT-HTTP/1.3 (0.8.16~exp12ubuntu10.21)"`
	d, _ := time.Parse("[02/Jan/2006:15:04:05 -0700]", "[17/May/2015:08:05:32 +0000]")
	expectedLog := LogStruct{
		"",
		"",
		d,
		"304 0",
		"nginx",
		HTTP{
			"\"GET /downloads/product_1",
			"93.180.71.3",
			"HTTP/1.1\"",
		},
	}
	parsedLog, _ := ParseLog("nginx", testString)

	if !reflect.DeepEqual(expectedLog, parsedLog) {
		fmt.Println("Res: ", parsedLog)
		fmt.Println("Exp: ", expectedLog)
		t.Fatal("Failed.")
	}
}

func TestParseSyslog(t *testing.T) {
	testString := `Oct 11 22:14:15 mymachine su: 'su root' failed for lonvick on /dev/pts/8`
	currentYear := time.Now().Year()
	d, _ := time.Parse("02/01/2006:15:04:05", fmt.Sprintf("11/10/%v:22:14:15", currentYear))
	expectedLog := LogStruct{
		"mymachine",
		"",
		d,
		"su: 'su root' failed for lonvick on /dev/pts/8",
		"syslog",
		HTTP{},
	}
	parsedLog, _ := ParseLog("syslog", testString)

	if !reflect.DeepEqual(expectedLog, parsedLog) {
		fmt.Println("Res: ", parsedLog)
		fmt.Println("Exp: ", expectedLog)
		t.Fatal("Failed.")
	}
}

func TestParseGeneric(t *testing.T) {
	testString := `2009/01/23 01:23:23 Hello world this is a log message.`
	d, _ := time.Parse("2006/01/02 15:04:05", "2009/01/23 01:23:23")
	expectedLog := LogStruct{
		"",
		"",
		d,
		"Hello world this is a log message.",
		"",
		HTTP{},
	}
	parsedLog, _ := ParseLog("", testString)

	if !reflect.DeepEqual(expectedLog, parsedLog) {
		fmt.Println("Res: ", parsedLog)
		fmt.Println("Exp: ", expectedLog)
		t.Fatal("Failed.")
	}
}

func TestParseEmpty(t *testing.T) {
	expectedLog := LogStruct{}
	parsedLog, _ := ParseLog("", "")

	if !reflect.DeepEqual(expectedLog, parsedLog) {
		fmt.Println("Res: ", parsedLog)
		fmt.Println("Exp: ", expectedLog)
		t.Fatal("Failed.")
	}
}
