package wal

import (
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
	Verb  base_types.Verb
}

// func GetMD5HashWithSum(value []byte) string {
// 	hash := md5.Sum(value)
// 	return hex.EncodeToString(hash[:])
// }

func NewLog(table_name string) (*os.File, *gob.Encoder) {

	f, err := os.Create(fmt.Sprintf("%v_%v.txt", table_name, uuid.NewString()))

	if err != nil {
		log.Fatal(err)
	}

	return f, gob.NewEncoder(f)

}

func AppendBinary[V any](enc *gob.Encoder, content V) error {
	return enc.Encode(content)
}
