package wal

import (
	"os"
	"strings"
	"testing"

	"gdarruda.me/todydbgo/base_types"
)

func TestNewLog(t *testing.T) {

	name := "table_example"
	file, _ := NewLog(name)

	if !strings.HasPrefix(file.Name(), name+"_") {
		t.Fatalf("Expected filename starting with log_")
	}

	file.Close()
	os.Remove(file.Name())

}

func TestAppendBinary(t *testing.T) {

	name := "table_example"
	log, enc := NewLog(name)

	AppendBinary[Record](enc, Record{[]byte("1"), []byte("a"), base_types.PUT})

	log.Close()
	os.Remove(log.Name())
}
