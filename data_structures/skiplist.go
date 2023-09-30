package data_structures

import (
	"bytes"
	"fmt"
	"math/rand"

	"gdarruda.me/todydbgo/base_types"
)

type KeyNotFoundError struct {
	Key []byte
}

func (knf *KeyNotFoundError) Error() string {
	return fmt.Sprintf("key (%v) not found", knf.Key)
}

type Node struct {
	key   []byte
	value []byte
	verb  base_types.Verb
	nexts []*Node
}

type SkipList struct {
	heads  []*Node
	levels int
}

func GetLevel() int {

	level := 1

	for {
		if rand.Intn(2) != 1 {
			break
		}
		level += 1
	}

	return level
}

func (sl *SkipList) Print() {

	for i := sl.levels - 1; i >= 0; i-- {

		step := sl.heads[i]

		for {

			if step == nil {
				fmt.Printf("X\n")
				break
			}

			fmt.Printf("%v --> ", step.key)
			step = step.nexts[i]
		}

	}

}

func (sl *SkipList) Get(key []byte) ([]byte, error) {

	level := sl.levels - 1
	node := sl.heads[level]
	befores := sl.heads

	for {

		if node == nil {
			return nil, &KeyNotFoundError{Key: key}
		}

		if bytes.Equal(key, node.key) {
			return node.value, nil
		}

		if bytes.Compare(key, node.key) == -1 {

			level -= 1

			if level < 0 {
				return nil, &KeyNotFoundError{Key: key}
			}

			node = befores[level]

		}

		if bytes.Compare(key, node.key) == 1 {

			next_node := node.nexts[level]

			if next_node == nil {

				level -= 1

				if level < 0 {
					return nil, &KeyNotFoundError{Key: key}
				}

			} else {
				befores = node.nexts
				node = next_node
			}

		}

	}
}

func (sl *SkipList) Insert(key []byte, value []byte) int {

	level := GetLevel()

	for {
		if level <= sl.levels {
			break
		}
		sl.levels += 1
		sl.heads = append(sl.heads, nil)
	}

	newNode := Node{
		key,
		value,
		base_types.PUT,
		make([]*Node, level)}

	for i := level - 1; i >= 0; i-- {

		if sl.heads[i] == nil {
			sl.heads[i] = &newNode
			continue
		}

		n := sl.heads[i]
		var b *Node

		for {

			if n == nil {
				b.nexts[i] = &newNode
				break
			}

			if bytes.Compare(key, n.key) == -1 {

				newNode.nexts[i] = n

				if b == nil {
					sl.heads[i] = &newNode
				} else {
					b.nexts[i] = &newNode
				}

				break
			}

			if bytes.Compare(key, n.key) == 1 {
				b = n
				n = n.nexts[i]
			}
		}
	}

	return level

}

func NewSkipList() SkipList {
	return SkipList{nil, 0}
}
