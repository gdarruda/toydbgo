package data_structures

import (
	"bytes"
	"strconv"
	"testing"
)

func TestNewList(t *testing.T) {

	empty_list := NewSkipList[int]()

	if empty_list.heads != nil {
		t.Fatalf("Expected empty list: %v", empty_list.heads)
	}

	if empty_list.levels != 0 {
		t.Fatalf("Expected 0 levels, found: %v", empty_list.levels)
	}

}

func TestInsert(t *testing.T) {

	list := NewSkipList[int]()

	list.Insert(13, []byte("13"))
	list.Insert(10, []byte("10"))
	list.Insert(15, []byte("15"))
	list.Insert(12, []byte("12"))

	list.Print()

	node := list.heads[0]

	for _, k := range [4]int{10, 12, 13, 15} {

		if node.key != k {
			t.Fatalf("Expected %v value on base list, found: %v", k, node)
		}

		node = node.nexts[0]
	}

}

func TestGet(t *testing.T) {

	list := NewSkipList[int]()

	for i := 100; i >= 0; i-- {
		list.Insert(i, []byte(strconv.FormatInt(int64(i), 10)))
	}

	list.Print()

	for i := 1; i <= 100; i++ {

		value, err := list.Get(i)
		expected := []byte(strconv.FormatInt(int64(i), 10))

		if err != nil {
			t.Fatalf("Error should be found, got %v", expected)
		}

		if !bytes.Equal(value, expected) {
			t.Fatalf("Value should be equal for %v, got %v", value, expected)
		}
	}

	value, err := list.Get(101)

	if value != nil {
		t.Fatalf("Value 101 wasn't added, can't be retrieved")
	}

	if err.Error() != (&KeyNotFoundError[int]{101}).Error() {
		t.Fatalf("Value 101 wasn't added, should cause an error")
	}

	value, err = list.Get(-1)

	if value != nil {
		t.Fatalf("Value -1 wasn't added, can't be retrieved")
	}

	if err.Error() != (&KeyNotFoundError[int]{-1}).Error() {
		t.Fatalf("Value -1 wasn't added, should cause an error")
	}
}
