package wal

import (
	"crypto/md5"
	"encoding/gob"
	"fmt"
	"log"
	"os"

	"gdarruda.me/todydbgo/base_types"
	"github.com/google/uuid"
)

type Record struct {
	Key   []byte
	Value []byte
	Hash  [16]byte
	Verb  base_types.Verb
}

type WAL struct {
	File *os.File
	enc  *gob.Encoder
}

func MD5KeyValue(key []byte, value []byte) [16]byte {
	return md5.Sum(append(value, key[:]...))
}

func NewRecord(
	key []byte,
	value []byte,
	verb base_types.Verb) Record {

	return Record{
		key,
		value,
		MD5KeyValue(key, value),
		verb}

}

func NewLog(table_name string) WAL {

	f, err := os.Create(fmt.Sprintf("%v_%v.tdb", table_name, uuid.NewString()))

	if err != nil {
		log.Fatal(err)
	}

	return WAL{f, gob.NewEncoder(f)}

}

func (wal *WAL) Append(content Record) error {
	return wal.enc.Encode(content)
}

func (wal *WAL) Delete() {
	wal.File.Close()
	os.Remove(wal.File.Name())
}
