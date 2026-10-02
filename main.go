package main

import (
	"fmt"
	"net/http"

	"buy/handlers"
)

func main() {
	http.HandleFunc("/", handlers.LoginHandler)
	http.HandleFunc("/login", handlers.LoginHandler)

	http.HandleFunc("/buyer", handlers.BuyerHandler)
	http.HandleFunc("/seller", handlers.SellerHandler)

	http.HandleFunc("/order", handlers.OrderHandler)
	http.HandleFunc("/order/accept", handlers.AcceptOrderHandler)

	fmt.Println("Server running on http://localhost:8080")

	err := http.ListenAndServe(":8080", nil)

	if err != nil {
		fmt.Println("Server error:", err)
	}
}
