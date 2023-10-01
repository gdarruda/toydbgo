package wal

import (
	"os"
	"strings"
	"testing"

	"gdarruda.me/todydbgo/base_types"
)

func TestNewLog(t *testing.T) {

	name := "table_example"
	wal := NewLog(name)

	if !strings.HasPrefix(wal.file.Name(), name+"_") {
		t.Fatalf("Expected filename starting with log_")
	}

	wal.file.Close()
	os.Remove(wal.file.Name())

}

func TestAppend(t *testing.T) {

	name := "table_example"
	wal := NewLog(name)

	wal.Append(NewRecord(
		[]byte("1"),
		[]byte("a"),
		base_types.PUT))

	wal.file.Close()
	os.Remove(wal.file.Name())
}
