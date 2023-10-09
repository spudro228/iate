package product

import "time"

type Product struct {
	ProductName   string    `json:"product_name"`
	Weight        int       `json:"weight"`
	Calories      int       `json:"calories"`
	Proteins      float64   `json:"proteins"`
	Fats          float64   `json:"fats"`
	Carbohydrates float64   `json:"carbohydrates"`
	CreatedAt     time.Time `json:"created_at"`
	Guid          string    `json:"guid"`
}
