package main

import (
	"bytes"
	"iate/product"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestProductHandler(t *testing.T) {
	data := []byte(`
		{
			"guid": "12345",
			"product_name": "Рис",
			"weight": 200,
			"calories": 340,
			"proteins": 4.2,
			"fats": 0.6,
			"carbohydrates": 43.5,
			"created_at": "2022-01-01T00:00:00Z"
		}
	`)

	req, err := http.NewRequest("POST", "/products", bytes.NewBuffer(data))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	service := product.NewInMemoryProductService()
	service.SaveProduct(product.Product{Guid: "12345"})
	handler := productHandlerUpdate(service)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v.", status, http.StatusOK)
		return
	}

	updatedProduct, _ := service.GetProduct("12345")
	if updatedProduct.ProductName != "Рис" {
		t.Errorf("handler does not update entity")
		return
	}
}

func TestProductsHandlerGetAll(t *testing.T) {
	service := product.NewInMemoryProductService()
	createdAt, _ := time.Parse(time.RFC3339, "2022-01-01T00:00:00Z")
	product := product.Product{Guid: "122345", ProductName: "Test", Weight: 100, Calories: 200, Proteins: 1.1, Fats: 2.2, Carbohydrates: 3.3, CreatedAt: createdAt}
	service.SaveProduct(product)

	req, err := http.NewRequest("GET", "/products/getAll", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := productsHandlerGetAll(service)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	r := rr.Body.String()
	expected := `[{"product_name":"Test","weight":100,"calories":200,"proteins":1.1,"fats":2.2,"carbohydrates":3.3,"created_at":"2022-01-01T00:00:00Z","guid":"122345"}]
`
	if r != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}
}
