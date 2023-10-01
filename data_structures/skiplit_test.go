package data_structures

import (
	"bytes"
	"fmt"
	"strconv"
	"testing"
)

func TestNewList(t *testing.T) {

	empty_list := NewSkipList(
		bytes.Compare,
		func(v []byte) string { return string(v[:]) })

	if empty_list.heads != nil {
		t.Fatalf("Expected empty list: %v", empty_list.heads)
	}

	if empty_list.levels != 0 {
		t.Fatalf("Expected 0 levels, found: %v", empty_list.levels)
	}

}

func TestInsert(t *testing.T) {

	list := NewSkipList(
		bytes.Compare,
		func(v []byte) string { return string(v[:]) })

	list.Insert([]byte("13"), []byte("13"))
	list.Insert([]byte("10"), []byte("10"))
	list.Insert([]byte("15"), []byte("15"))
	list.Insert([]byte("12"), []byte("12"))

	list.Print()

	node := list.heads[0]

	for _, k := range [4]string{"10", "12", "13", "15"} {

		if !bytes.Equal(node.Key, []byte(k)) {
			t.Fatalf("Expected %v value on base list, found: %v", k, node)
		}

		node = node.nexts[0]
	}

}

func TestGet(t *testing.T) {

	list := NewSkipList(
		bytes.Compare,
		func(v []byte) string { return string(v[:]) })

	for i := 100; i >= 0; i-- {

		list.Insert(
			[]byte(fmt.Sprintf("%03d", i)),
			[]byte(strconv.FormatInt(int64(i), 10)))
	}

	list.Print()

	for i := 1; i <= 100; i++ {

		node, err := list.Get([]byte(fmt.Sprintf("%03d", i)))
		expected := []byte(strconv.FormatInt(int64(i), 10))

		if err != nil {
			t.Fatalf("Error should be found, got %v", expected)
		}

		if !bytes.Equal(node.Value, expected) {
			t.Fatalf("Value should be equal for %v, got %v", node.Value, expected)
		}
	}

	not_present_key := []byte(fmt.Sprintf("%03d", 101))

	value, err := list.Get(not_present_key)

	if value != nil {
		t.Fatalf("Value 101 wasn't added, shouldn't be retrieved")
	}

	if err.Error() != (&KeyNotFoundError{"101"}).Error() {
		t.Fatalf("Value 101 wasn't added, should cause an error")
	}

}
