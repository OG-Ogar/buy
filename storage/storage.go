package storage

import "buy/models"

var Orders []models.Order

func AddOrder(order models.Order) {
	Orders = append(Orders, order)
}

func GetOrders() []models.Order {
	return Orders
}

func AcceptOrder(id int) bool {
	for index := range Orders {
		if Orders[index].ID == id {
			Orders[index].Status = "Accepted"
			return true
		}
	}

	return false
}
