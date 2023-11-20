package table

import (
	"bytes"

	"gdarruda.me/todydbgo/base_types"
	"gdarruda.me/todydbgo/data_structures"
	"gdarruda.me/todydbgo/wal"
)

type Table struct {
	name           string
	content        data_structures.SkipList
	wal            wal.WAL
	merge_function func([]byte, []byte) []byte
	size           int
	max_size       int
}

func Create(
	name string,
	merge_function func([]byte, []byte) []byte,
	max_size int) Table {

	if merge_function == nil {
		merge_function = func(old []byte, new []byte) []byte { return new }
	}

	if max_size == 0 {
		max_size = 64_000_000
	}

	list := data_structures.NewSkipList(
		bytes.Compare,
		func(v []byte) string { return string(v[:]) })

	table := Table{
		name,
		list,
		wal.NewLog(name),
		merge_function,
		0,
		max_size}

	return table

}

func (t *Table) Put(
	key []byte,
	value []byte) {

	t.wal.Append(wal.NewRecord(key, value, base_types.PUT))
	t.content.Insert(key, value)
	t.size += len(key) + len(value)

}

func (t *Table) Merge(
	key []byte,
	value []byte) (*data_structures.Node, error) {

	node, err := t.content.Get(key)
	value_size := len(node.Value)

	if err != nil {
		return nil, err
	}

	t.wal.Append(wal.NewRecord(key, value, base_types.MERGE))
	new_value := t.merge_function(node.Value, value)
	t.size += len(new_value) - value_size

	node.Value = new_value

	return node, err

}

func (t *Table) Delete(
	key []byte) (*data_structures.Node, error) {

	node, err := t.content.Get(key)
	value_size := len(node.Value)

	if err != nil {
		return nil, err
	}

	t.wal.Append(wal.NewRecord(node.Key, nil, base_types.DEL))
	t.size -= value_size

	node.Value = nil
	node.Verb = base_types.DEL

	return node, err

}
