package data_structures

import (
	"fmt"
	"math/rand"

	"gdarruda.me/todydbgo/base_types"
)

type KeyNotFoundError struct {
	Key string
}

func (knf *KeyNotFoundError) Error() string {
	return fmt.Sprintf("key (%v) not found", knf.Key)
}

type Node struct {
	Key   []byte
	Value []byte
	Verb  base_types.Verb
	nexts []*Node
}

type SkipList struct {
	Heads     []*Node
	Levels    int
	compare   func([]byte, []byte) int
	formatKey func([]byte) string
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

	for i := sl.Levels - 1; i >= 0; i-- {

		step := sl.Heads[i]

		for {

			if step == nil {
				fmt.Printf("X\n")
				break
			}

			fmt.Printf("%v --> ", sl.formatKey(step.Key))
			step = step.nexts[i]
		}

	}

}

func (sl *SkipList) Get(key []byte) (*Node, error) {

	level := sl.Levels - 1
	node := sl.Heads[level]
	befores := sl.Heads

	for {

		if node == nil {
			return nil, &KeyNotFoundError{Key: sl.formatKey(key)}
		}

		if sl.compare(key, node.Key) == 0 {
			return node, nil
		}

		if sl.compare(key, node.Key) == -1 {

			level -= 1

			if level < 0 {
				return nil, &KeyNotFoundError{Key: sl.formatKey(key)}
			}

			node = befores[level]

		}

		if sl.compare(key, node.Key) == 1 {

			next_node := node.nexts[level]

			if next_node == nil {

				level -= 1

				if level < 0 {
					return nil, &KeyNotFoundError{Key: sl.formatKey(key)}
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
		if level <= sl.Levels {
			break
		}
		sl.Levels += 1
		sl.Heads = append(sl.Heads, nil)
	}

	newNode := Node{
		key,
		value,
		base_types.PUT,
		make([]*Node, level)}

	for i := level - 1; i >= 0; i-- {

		if sl.Heads[i] == nil {
			sl.Heads[i] = &newNode
			continue
		}

		n := sl.Heads[i]
		var b *Node

		for {

			if n == nil {
				b.nexts[i] = &newNode
				break
			}

			if sl.compare(key, n.Key) == -1 {

				newNode.nexts[i] = n

				if b == nil {
					sl.Heads[i] = &newNode
				} else {
					b.nexts[i] = &newNode
				}

				break
			}

			if sl.compare(key, n.Key) == 1 {
				b = n
				n = n.nexts[i]
			}
		}
	}

	return level

}

func NewSkipList(compare func([]byte, []byte) int, formatKey func([]byte) string) SkipList {

	return SkipList{nil,
		0,
		compare,
		formatKey}
}
