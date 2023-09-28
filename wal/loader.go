package wal

import (
	"encoding/gob"
	"io"
	"log"
	"os"
)

func LoadLog[V any](
	file *os.File,
	apply func(V)) {

	dec := gob.NewDecoder(file)

	for {

		var v V
		err := dec.Decode(&v)

		if err == io.EOF {
			break
		}

		if err != nil {
			log.Fatal("decode error:", err)
		}

		apply(v)
	}

}
