package model

type SelectResponse struct {
	Products []Product `json:"products"`
	Count    int       `json:"count"`
}
