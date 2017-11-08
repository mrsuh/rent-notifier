package model

type Message struct {
	ChatId  int
	ChatIds []int
	Text    string
	IsBulk  bool
}
