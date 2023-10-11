package product

import (
	"fmt"
	"time"
)

type UserId string

type ProductDataModel struct {
	UserId  UserId
	Product Product
}
type InMemoryProductService struct {
	products map[string]ProductDataModel
}

func NewInMemoryProductService() *InMemoryProductService {
	return &InMemoryProductService{
		products: map[string]ProductDataModel{},
	}
}

func (s *InMemoryProductService) SaveProduct(product Product, userId UserId) error {
	if product.Guid == "" {
		return fmt.Errorf("Empty guid")
	}

	s.products[product.Guid] = ProductDataModel{UserId: userId, Product: product}
	return nil
}

func (s *InMemoryProductService) GetProduct(id string) (*Product, error) {
	product, ok := s.products[id]
	if !ok {
		return nil, fmt.Errorf("Product not found")
	}
	return &product.Product, nil
}

// GetAllProducts todo: add filter "today"
func (s *InMemoryProductService) GetAllProducts(userId UserId, today bool) ([]Product, error) {
	values := make([]Product, 0, len(s.products))
	for _, v := range s.products {
		if v.UserId == userId {
			values = append(values, v.Product)
		}
	}

	if today {
		todayProducts := make([]Product, 0, len(values))
		for _, value := range values {
			if isSameDay(value.CreatedAt, time.Now()) {
				todayProducts = append(todayProducts, value)
			}
		}

		return todayProducts, nil
	}

	return values, nil
}

func isSameDay(t1, t2 time.Time) bool {
	return t1.Year() == t2.Year() && t1.Month() == t2.Month() && t1.Day() == t2.Day()
}

func (s *InMemoryProductService) DeleteProduct(id string) error {
	_, ok := s.products[id]
	if !ok {
		return fmt.Errorf("Product not found")
	}
	delete(s.products, id)
	return nil
}

func (s *InMemoryProductService) UpdateProduct(productObj Product) error {
	p, ok := s.products[productObj.Guid]
	if !ok {
		return fmt.Errorf("Product not found")
	}

	p.Product.ProductName = productObj.ProductName
	p.Product.Calories = productObj.Calories
	p.Product.Fats = productObj.Fats
	p.Product.Carbohydrates = productObj.Carbohydrates
	p.Product.Proteins = productObj.Proteins
	p.Product.Weight = productObj.Weight

	s.products[productObj.Guid] = p
	return nil
}
