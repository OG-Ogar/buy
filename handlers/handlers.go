package handlers

import (
	"html/template"
	"net/http"
	"strconv"

	"buy/models"
	"buy/storage"
)

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		tmpl, err := template.ParseFiles("templates/login.html")
		if err != nil {
			http.Error(w, "Unable to load login page: "+err.Error(), http.StatusInternalServerError)
			return
		}

		err = tmpl.Execute(w, nil)
		if err != nil {
			http.Error(w, "Unable to display login page: "+err.Error(), http.StatusInternalServerError)
		}

		return
	}

	role := r.FormValue("role")

	if role == "buyer" {
		http.Redirect(w, r, "/buyer", http.StatusSeeOther)
		return
	}

	if role == "seller" {
		http.Redirect(w, r, "/seller", http.StatusSeeOther)
		return
	}

	http.Error(w, "Invalid role", http.StatusBadRequest)
}

func BuyerHandler(w http.ResponseWriter, r *http.Request) {
	// Get the buyer's name from the cookie.
	cookie, err := r.Cookie("buyer_name")

	buyerName := ""

	if err == nil {
		buyerName = cookie.Value
	}

	// Get all orders from storage.
	orders := storage.GetOrders()

	// Keep only orders belonging to this buyer.
	var buyerOrders []models.Order

	for _, order := range orders {
		if order.Customer == buyerName {
			buyerOrders = append(buyerOrders, order)
		}
	}

	tmpl, err := template.ParseFiles("templates/buyer.html")
	if err != nil {
		http.Error(w, "Unable to load buyer dashboard", http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, buyerOrders)
	if err != nil {
		http.Error(w, "Unable to display buyer dashboard", http.StatusInternalServerError)
	}
}

func SellerHandler(w http.ResponseWriter, r *http.Request) {
	orders := storage.GetOrders()

	tmpl, err := template.ParseFiles("templates/seller.html")
	if err != nil {
		http.Error(w, "Unable to load seller dashboard", http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, orders)
	if err != nil {
		http.Error(w, "Unable to display seller dashboard", http.StatusInternalServerError)
	}
}

func OrderHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	customer := r.FormValue("customer")
	item := r.FormValue("item")
	quantityText := r.FormValue("quantity")

	quantity, err := strconv.Atoi(quantityText)
	if err != nil || quantity < 1 {
		http.Error(w, "Invalid quantity", http.StatusBadRequest)
		return
	}

	order := models.Order{
		ID:       len(storage.GetOrders()) + 1,
		Customer: customer,
		Item:     item,
		Quantity: quantity,
		Status:   "Pending",
	}

	// Save the order.
	storage.AddOrder(order)

	// Remember which buyer placed the order.
	http.SetCookie(w, &http.Cookie{
		Name:  "buyer_name",
		Value: customer,
		Path:  "/",
	})

	http.Redirect(w, r, "/buyer", http.StatusSeeOther)
}

func AcceptOrderHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	idText := r.FormValue("id")

	id, err := strconv.Atoi(idText)
	if err != nil {
		http.Error(w, "Invalid order ID", http.StatusBadRequest)
		return
	}

	accepted := storage.AcceptOrder(id)

	if !accepted {
		http.Error(w, "Order not found", http.StatusNotFound)
		return
	}

	http.Redirect(w, r, "/seller", http.StatusSeeOther)
}
