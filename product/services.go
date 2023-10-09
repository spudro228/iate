package product

import "fmt"

type InMemoryProductService struct {
	products map[string]Product
}

func NewInMemoryProductService() *InMemoryProductService {
	return &InMemoryProductService{
		products: make(map[string]Product),
	}
}

func (s *InMemoryProductService) SaveProduct(product Product) error {
	s.products[product.Guid] = product
	return nil
}

func (s *InMemoryProductService) GetProduct(id string) (*Product, error) {
	product, ok := s.products[id]
	if !ok {
		return nil, fmt.Errorf("Product not found")
	}
	return &product, nil
}

func (s *InMemoryProductService) GetAllProducts() ([]Product, error) {
	values := make([]Product, 0, len(s.products))
	for _, v := range s.products {
		values = append(values, v)
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
	_, ok := s.products[productObj.Guid]
	if !ok {
		return fmt.Errorf("Product not found")
	}

	s.products[productObj.Guid] = productObj
	return nil
}
