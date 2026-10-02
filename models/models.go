package models

type Order struct {
	ID       int
	Customer string
	Item     string
	Quantity int
	Status   string
}