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
)

func TestReadFile(t *testing.T) {
	testString := `2009/01/23 01:23:23 Hello world this is a log message.`
	d, _ := time.Parse("02/01/2006 15:04:05", "2009/01/23 01:23:23")
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
	
	wg.Add(1)
	ReadFile(file.Name(), &wg, testChannel, "generic", false)

	val := <-testChannel

	if !reflect.DeepEqual(expectedLog, val) {
		fmt.Println("Res: ", val)
		fmt.Println("Exp: ", expectedLog)
		t.Fatal("Failed.")
	}
}
