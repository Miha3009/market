package model

type CreateResponse struct {
	Id int `json:"id"`
}

type SelectResponse struct {
	Products []Product `json:"products"`
	Count    int       `json:"count"`
}
