package wal

import (
	"os"
	"testing"

	"gdarruda.me/todydbgo/base_types"
)

func TestLoadLog(t *testing.T) {

	name := "table_example"
	new_log, enc := NewLog(name)

	log_name := new_log.Name()

	AppendBinary[Record](enc, Record{[]byte("1"), []byte("a"), base_types.PUT})
	AppendBinary[Record](enc, Record{[]byte("2"), []byte("b"), base_types.MERGE})

	new_log.Close()

	log, err := os.Open(log_name)

	if err != nil {
		t.Fatal(err)
	}

	LoadLog[Record](log, func(value Record) {})
	log.Close()

	os.Remove(log.Name())
}
