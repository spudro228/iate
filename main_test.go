package main

import (
	"bytes"
	"iate/product"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestProductsHandlerGetAll(t *testing.T) {
	service := product.NewInMemoryProductService()
	createdAt, _ := time.Parse(time.RFC3339, "2023-01-01T00:00:12.842427065Z")
	createdAtCustom := product.CustomDate{Time: createdAt}
	product := product.Product{Guid: "122345", ProductName: "Test", Weight: 100, Calories: 200, Proteins: 1.1, Fats: 2.2, Carbohydrates: 3.3, CreatedAt: createdAtCustom}
	service.SaveProduct(product, "u0001")

	data := []byte(`
		{
			"user_id": "u0001",
			"today": false
		}
	`)

	req, err := http.NewRequest("GET", "/products/getAll", bytes.NewBuffer(data))
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
	expected := `[{"product_name":"Test","weight":100,"calories":200,"proteins":1.1,"fats":2.2,"carbohydrates":3.3,"created_at":"2023-01-01T00:00:12Z","guid":"122345"}]
`
	if r != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}
}
