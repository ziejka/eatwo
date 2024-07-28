package models

type CheckList struct {
	Name   string
	UserID string
}

type CheckListRecord struct {
	ID uint
	CheckList
}

type CheckListItem struct {
	Value  string
	ListID uint
}

type CheckListItemRecord struct {
	ID uint
	CheckListItem
}

type ListWithItems struct {
	CheckListRecord
	Items []CheckListItemRecord
}
