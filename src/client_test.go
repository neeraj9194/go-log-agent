package src

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"sync"
	"testing"

	"github.com/neeraj9194/go-log-agent/config"
)

func TestReadFile(t *testing.T) {
	testString := `2009/01/23 01:23:23 Hello world this is a log message.`

	// Create a temp file
	file, err := ioutil.TempFile(".", "logfile")
	if err != nil {
		log.Fatal(err)
	}
	defer os.Remove(file.Name())

	file.WriteString(fmt.Sprintf("%v\n", testString))
	
	testChannel := make(chan LogStruct, 100)
	var wg sync.WaitGroup
	watcher := config.Watcher {
		FilePath: file.Name(),   
		ServiceName: "generic",
	}
	
	wg.Add(1)
	ReadFile(watcher, &wg, testChannel, false)

	val := <-testChannel
	if val.Message != testString {
		t.Fatal("Failed.")
	}
	
}
