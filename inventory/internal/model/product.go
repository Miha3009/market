package model

type Product struct {
	Id       int
	Count    int
	Reserved int
}

type ReserveRequest struct {
	ProductId int
}
