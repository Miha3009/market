package model

type Product struct {
	Id          int    `json:"id"`
	Price       int    `json:"price"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Avaliable   bool   `json:"avaliable"`
}
