package product

import "fmt"

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
func (s *InMemoryProductService) GetAllProducts(userId UserId) ([]Product, error) {
	values := make([]Product, 0, len(s.products))
	for _, v := range s.products {
		if v.UserId == userId {
			values = append(values, v.Product)
		}
	}
	return values, nil
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
