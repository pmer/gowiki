package main

const (
	LocNone = iota
	LocInbox
	LocToday
	LocUpcoming
	LocAnytime
	LocSomeday
	LocLogbook
	LocTrash
)

const (
	Unchecked = iota
	CheckedOff
	CrossedOut
)

type Location struct {
}

type Todo struct {
	id       int
	Title    string
	Notes    string
	Location int
	Checked  int
}
