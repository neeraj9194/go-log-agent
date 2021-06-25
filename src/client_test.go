package src

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"reflect"
	"sync"
	"testing"
	"time"

	"github.com/neeraj9194/go-log-agent/config"
)

func TestReadFile(t *testing.T) {
	testString := `2009/01/23 01:23:23 Hello world this is a log message.`
	d, _ := time.Parse("2006/01/02 15:04:05", "2009/01/23 01:23:23")
	expectedLog := LogStruct{
		"",
		"",
		d,
		"Hello world this is a log message.",
		"generic",
		HTTP{},
	}

	// Create a temp file
	file, err := ioutil.TempFile(".", "logfile")
	if err != nil {
		log.Fatal(err)
	}
	defer os.Remove(file.Name())

	file.WriteString(fmt.Sprintf("%v\n", testString))
	
	testChannel := make(chan LogStruct, 100)
	var wg sync.WaitGroup

	conf := config.Config{
		FilePath: file.Name(),   
		ServiceName: "generic",
		ServerURL: "string",
	}
	
	wg.Add(1)
	ReadFile(conf, &wg, testChannel, false)

	val := <-testChannel

	if !reflect.DeepEqual(expectedLog, val) {
		fmt.Println("Res: ", val)
		fmt.Println("Exp: ", expectedLog)
		t.Fatal("Failed.")
	}
}
