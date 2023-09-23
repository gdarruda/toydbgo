package base_types

type Verb int16

const (
	PUT Verb = iota
	MERGE
	DEL
)
