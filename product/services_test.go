package product

import (
	"testing"
)

func TestSaveProduct(t *testing.T) {
	service := NewInMemoryProductService()
	product := Product{Guid: "g12312-0", ProductName: "Test", Weight: 100, Calories: 200, Proteins: 1.1, Fats: 2.2, Carbohydrates: 3.3}
	err := service.SaveProduct(product, "u001")
	if err != nil {
		t.Errorf("Error while saving product: %v", err)
	}
}

func TestGetProduct(t *testing.T) {
	service := NewInMemoryProductService()
	product := Product{Guid: "g12312-0", ProductName: "Test", Weight: 100, Calories: 200, Proteins: 1.1, Fats: 2.2, Carbohydrates: 3.3}
	service.SaveProduct(product, "u001")

	retrievedProduct, err := service.GetProduct("g12312-0")
	if err != nil {
		t.Errorf("Error while getting product: %v", err)
		return
	}

	if retrievedProduct.ProductName != "Test" {
		t.Errorf("Expected product name to be 'Test', got '%v'", retrievedProduct.ProductName)
		return
	}
}

func TestGetAllProducts(t *testing.T) {
	service := NewInMemoryProductService()
	product1 := Product{Guid: "g12312-0", ProductName: "Test1", Weight: 100, Calories: 200, Proteins: 1.1, Fats: 2.2, Carbohydrates: 3.3}
	product2 := Product{Guid: "g12312-1", ProductName: "Test2", Weight: 200, Calories: 300, Proteins: 2.2, Fats: 3.3, Carbohydrates: 4.4}
	service.SaveProduct(product1, "u001")
	service.SaveProduct(product2, "u001")

	products, err := service.GetAllProducts("u001")
	if err != nil {
		t.Errorf("Error while getting all products: %v", err)
		return
	}

	if len(products) != 2 {
		t.Errorf("Expected 2 products, got %v", len(products))
		return
	}
}

func TestDeleteProduct(t *testing.T) {
	service := NewInMemoryProductService()
	product := Product{Guid: "g12312-0", ProductName: "Test", Weight: 100, Calories: 200, Proteins: 1.1, Fats: 2.2, Carbohydrates: 3.3}
	service.SaveProduct(product, "u001")

	err := service.DeleteProduct("g12312-0")
	if err != nil {
		t.Errorf("Error while deleting product: %v", err)
		return
	}

	_, err = service.GetProduct("12345")
	if err == nil {
		t.Errorf("Expected error when getting deleted product, got nil")
		return
	}
}
