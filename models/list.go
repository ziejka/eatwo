package models

type List struct {
	Name   string
	UserID string
}

type ListRecord struct {
	ID uint
	List
}

type Item struct {
	Value  string
	ListID string
}

type ItemRecord struct {
	ID uint
	Item
}
