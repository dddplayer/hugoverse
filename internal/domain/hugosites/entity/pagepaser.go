package entity

// Result holds the parse result.
type Result interface {
	// Iterator returns a new Iterator positioned at the beginning of the parse tree.
	Iterator() *Iterator
	// Input returns the input to Parse.
	Input() []byte
}

// An Iterator has methods to iterate a parsed page with support going back
// if needed.
type Iterator struct {
	items   Items
	lastPos int // position of the last item returned by nextItem
}

type Items []Item

type lowHigh struct {
	Low  int
	High int
}

type Item struct {
	Type ItemType
	Err  error

	// The common case is a single segment.
	low  int
	high int

	// This is the uncommon case.
	segments []lowHigh

	// Used for validation.
	firstByte byte

	isString bool
}

type ItemType int
