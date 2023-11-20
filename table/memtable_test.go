package table

import (
	"bytes"
	"encoding/gob"
	"os"
	"path/filepath"
	"testing"

	"gdarruda.me/todydbgo/base_types"
	"gdarruda.me/todydbgo/wal"
)

func TestCreate(t *testing.T) {

	name := "example_table"
	table := Create(name, nil)

	if table.name != name {
		t.Fatalf("Expected table name '%v', got '%v' instead", name, table.name)
	}

	table.wal.Delete()

}

func Setup() Table {
	name := "example_table"
	return Create(name, nil)
}

func Cleanup(f *os.File, table Table) {
	f.Close()
	table.wal.Delete()

	files, _ := filepath.Glob("example_table_*")

	for _, f := range files {
		os.Remove(f)
	}
}

func TestPut(t *testing.T) {

	table := Setup()

	table.Put([]byte("1"), []byte("a"))
	table.Put([]byte("2"), []byte("b"))

	node, _ := table.content.Get([]byte("1"))

	if !bytes.Equal(node.Value, []byte("a")) {
		t.Fatalf("Expected value a, got '%v' instead", node.Value)
	}

	node, _ = table.content.Get([]byte("2"))

	if !bytes.Equal(node.Value, []byte("b")) {
		t.Fatalf("Expected value b, got '%v' instead", node.Value)
	}

	f, _ := os.Open(table.wal.File.Name())
	dec := gob.NewDecoder(f)

	for _, element := range [2]string{"1", "2"} {

		var v wal.Record
		err := dec.Decode(&v)

		if err != nil {
			t.Fatalf("Error loading WAL: %v", err)
		}

		if !bytes.Equal(v.Key, []byte(element)) {
			t.Fatalf("Expected key '%v', got '%v' instead", v.Key, element)
		}

		if v.Verb != base_types.PUT {
			t.Fatalf("Expected verb PUT, got '%v' instead", v.Verb)
		}
	}

	t.Cleanup(func() { Cleanup(f, table) })

}

func TestMerge(t *testing.T) {

	table := Setup()

	table.Put([]byte("1"), []byte("a"))
	table.Merge([]byte("1"), []byte("new value"))

	node, _ := table.content.Get([]byte("1"))

	if !bytes.Equal(node.Value, []byte("new value")) {
		t.Fatalf(
			"Expected value '%v', got '%v' instead",
			node.Value,
			[]byte("new value"))
	}

	f, _ := os.Open(table.wal.File.Name())
	dec := gob.NewDecoder(f)

	var put, merge wal.Record

	dec.Decode(&put)

	if !bytes.Equal(put.Value, []byte("a")) {
		t.Fatalf("Put log should contain a, got '%v' instead", put.Value)
	}

	dec.Decode(&merge)

	if !bytes.Equal(merge.Value, []byte("new value")) {
		t.Fatalf("Put log should contain b, got '%v' instead", merge.Value)
	}

	t.Cleanup(func() { Cleanup(f, table) })
}

func TestDelete(t *testing.T) {

	table := Setup()

	table.Put([]byte("1"), []byte("a"))
	node, _ := table.Delete([]byte("1"))

	if node.Verb != base_types.DEL {
		t.Fatalf("Verb should be deleted, got %v instead", node.Verb)
	}

	f, _ := os.Open(table.wal.File.Name())
	dec := gob.NewDecoder(f)

	var delete wal.Record

	dec.Decode(&delete)
	dec.Decode(&delete)

	if delete.Verb != base_types.DEL {
		t.Fatalf("Put log should change verb do DEL, got '%v' instead", delete.Verb)
	}

	t.Cleanup(func() { Cleanup(f, table) })

}
