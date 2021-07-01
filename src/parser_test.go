package src

import (
	"testing"
)


func TestParseGeneric(t *testing.T) {
	testString := `2009/01/23 01:23:23 Hello world this is a log message.`
	
	parsedLog, _ := ParseLog("", testString)

	if parsedLog.Message != testString {
		t.Fatal("Failed.")
	}
}

func TestParseEmpty(t *testing.T) {
	expectedLog := LogStruct{}

	if expectedLog.Message != "" {
		t.Fatal("Failed.")
	}
}
