package data_structures

import (
	"fmt"
	"math/rand"

	"gdarruda.me/todydbgo/base_types"
	"golang.org/x/exp/constraints"
)

type KeyNotFoundError[K constraints.Ordered] struct {
	Key K
}

func (knf *KeyNotFoundError[K]) Error() string {
	return fmt.Sprintf("key (%v) not found", knf.Key)
}

type Node[K constraints.Ordered] struct {
	key   K
	value []byte
	verb  base_types.Verb
	nexts []*Node[K]
}

type SkipList[K constraints.Ordered] struct {
	heads  []*Node[K]
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

func (this *SkipList[K]) Print() {

	for i := this.levels - 1; i >= 0; i-- {

		step := this.heads[i]

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

func (this *SkipList[K]) Get(key K) ([]byte, error) {

	level := this.levels - 1
	node := this.heads[level]
	befores := this.heads

	for {

		if node == nil {
			return nil, &KeyNotFoundError[K]{Key: key}
		}

		if node.key == key {
			return node.value, nil
		}

		if key < node.key {

			level -= 1

			if level < 0 {
				return nil, &KeyNotFoundError[K]{Key: key}
			}

			node = befores[level]

		}

		if key > node.key {

			next_node := node.nexts[level]

			if next_node == nil {

				level -= 1

				if level < 0 {
					return nil, &KeyNotFoundError[K]{Key: key}
				}

			} else {
				befores = node.nexts
				node = next_node
			}

		}

	}
}

func (this *SkipList[K]) Insert(key K, value []byte) int {

	level := GetLevel()

	for {
		if level <= this.levels {
			break
		}
		this.levels += 1
		this.heads = append(this.heads, nil)
	}

	newNode := Node[K]{
		key,
		value,
		base_types.PUT,
		make([]*Node[K], level)}

	for i := level - 1; i >= 0; i-- {

		if this.heads[i] == nil {
			this.heads[i] = &newNode
			continue
		}

		n := this.heads[i]
		var b *Node[K]

		for {

			if n == nil {
				b.nexts[i] = &newNode
				break
			}

			if key < n.key {

				newNode.nexts[i] = n

				if b == nil {
					this.heads[i] = &newNode
				} else {
					b.nexts[i] = &newNode
				}

				break
			}

			if key > n.key {
				b = n
				n = n.nexts[i]
			}
		}
	}

	return level

}

func NewSkipList[K constraints.Ordered]() SkipList[K] {
	return SkipList[K]{nil, 0}
}
