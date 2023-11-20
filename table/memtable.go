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
}

func Create(
	name string,
	merge_function func([]byte, []byte) []byte) Table {

	if merge_function == nil {
		merge_function = func(old []byte, new []byte) []byte { return new }
	}

	list := data_structures.NewSkipList(
		bytes.Compare,
		func(v []byte) string { return string(v[:]) })

	table := Table{
		name,
		list,
		wal.NewLog(name),
		merge_function}

	return table

}

func (t *Table) Put(
	key []byte,
	value []byte) {

	t.wal.Append(wal.NewRecord(key, value, base_types.PUT))
	t.content.Insert(key, value)

}

func (t *Table) Merge(
	key []byte,
	value []byte) (*data_structures.Node, error) {

	node, err := t.content.Get(key)

	if err != nil {
		return nil, err
	}

	t.wal.Append(wal.NewRecord(key, value, base_types.MERGE))

	node.Value = t.merge_function(node.Value, value)

	return node, err

}

func (t *Table) Delete(
	key []byte) (*data_structures.Node, error) {

	node, err := t.content.Get(key)

	if err != nil {
		return nil, err
	}

	t.wal.Append(wal.NewRecord(node.Key, nil, base_types.DEL))

	node.Verb = base_types.DEL

	return node, err

}
