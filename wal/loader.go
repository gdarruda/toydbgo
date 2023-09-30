package wal

import (
	"encoding/gob"
	"io"
	"log"
	"os"
)

func LoadLog(
	file *os.File,
	apply func(Record)) {

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

		apply(v)
	}

}
