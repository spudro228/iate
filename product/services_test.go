package product

import (
	"testing"
)

func TestSaveProduct(t *testing.T) {
	service := NewInMemoryProductService()
	product := Product{ProductName: "Test", Weight: 100, Calories: 200, Proteins: 1.1, Fats: 2.2, Carbohydrates: 3.3}
	err := service.SaveProduct(product)
	if err != nil {
		t.Errorf("Error while saving product: %v", err)
	}
}

func TestGetProduct(t *testing.T) {
	service := NewInMemoryProductService()
	product := Product{ProductName: "Test", Weight: 100, Calories: 200, Proteins: 1.1, Fats: 2.2, Carbohydrates: 3.3}
	service.SaveProduct(product)

	retrievedProduct, err := service.GetProduct("Test")
	if err != nil {
		t.Errorf("Error while getting product: %v", err)
	}

	if retrievedProduct.ProductName != "Test" {
		t.Errorf("Expected product name to be 'Test', got '%v'", retrievedProduct.ProductName)
	}
}

func TestGetAllProducts(t *testing.T) {
	service := NewInMemoryProductService()
	product1 := Product{ProductName: "Test1", Weight: 100, Calories: 200, Proteins: 1.1, Fats: 2.2, Carbohydrates: 3.3}
	product2 := Product{ProductName: "Test2", Weight: 200, Calories: 300, Proteins: 2.2, Fats: 3.3, Carbohydrates: 4.4}
	service.SaveProduct(product1)
	service.SaveProduct(product2)

	products, err := service.GetAllProducts()
	if err != nil {
		t.Errorf("Error while getting all products: %v", err)
	}

	if len(products) != 2 {
		t.Errorf("Expected 2 products, got %v", len(products))
	}
}

func TestDeleteProduct(t *testing.T) {
	service := NewInMemoryProductService()
	product := Product{Guid: "12345", ProductName: "Test", Weight: 100, Calories: 200, Proteins: 1.1, Fats: 2.2, Carbohydrates: 3.3}
	service.SaveProduct(product)

	err := service.DeleteProduct("12345")
	if err != nil {
		t.Errorf("Error while deleting product: %v", err)
	}

	_, err = service.GetProduct("12345")
	if err == nil {
		t.Errorf("Expected error when getting deleted product, got nil")
	}
}
