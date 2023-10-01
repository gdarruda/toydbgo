package wal

import (
	"os"
	"testing"

	"gdarruda.me/todydbgo/base_types"
)

func TestLoadLog(t *testing.T) {

	name := "table_example"
	wal := NewLog(name)

	log_name := wal.file.Name()

	wal.Append(NewRecord([]byte("1"), []byte("a"), base_types.PUT))
	wal.Append(NewRecord([]byte("1"), []byte("b"), base_types.MERGE))

	wal.file.Close()

	log, err := os.Open(log_name)

	if err != nil {
		t.Fatal(err)
	}

	wal.Load(log, func(value Record, wal *WAL) {})

	log.Close()

	os.Remove(log.Name())
}
