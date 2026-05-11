package model

type Message struct {
	ID  int    `json:"id" db:"id"`
	Msg string `json:"msg" db:"msg"`
}
