package model

import "time"

type User struct {
	Id   int
	Name string
}

type Item struct {
	Id    int
	Price int
	Count int
}

type Order struct {
	Id    int
	Buyer User
	Date  time.Time
	Items []Item
}
