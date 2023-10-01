package wal

import (
	"encoding/gob"
	"io"
	"log"
	"os"
)

func (wal *WAL) Load(
	file *os.File,
	apply func(Record, *WAL)) {

	dec := gob.NewDecoder(file)

	for {

		var v Record
		err := dec.Decode(&v)

		if err == io.EOF {
			break
		}

		if err != nil {
			log.Fatal("decode error:", err)
		}

		if v.Hash != MD5KeyValue(v.Key, v.Value) {
			log.Fatal("loaded values and hash don't match:", err)
		}

		apply(v, wal)

	}
}
