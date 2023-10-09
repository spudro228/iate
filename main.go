package main

import (
	"encoding/json"
	"iate/product"
	"io/ioutil"
	"log"
	"net/http"
)

func productHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusInternalServerError)
		return
	}

	var products []product.Product
	err = json.Unmarshal(body, &products)
	if err != nil {
		http.Error(w, "Error decoding JSON", http.StatusBadRequest)
		return
	}

	for _, product := range products {
		log.Printf("%+v\n", product)
	}

	w.WriteHeader(http.StatusOK)
}

func productsHandlerGetAll(service *product.InMemoryProductService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		products, err := service.GetAllProducts()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(products)
	}
}

func main() {
	productService := product.NewInMemoryProductService()

	http.HandleFunc("/products", productHandler)
	http.HandleFunc("/products/getAll", productsHandlerGetAll(productService))
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
