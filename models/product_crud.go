package models

type DeleteProductById struct {
	Guid string `json:"guid"`
}

type ProductUpdateModel struct {
	ProductName   string  `json:"product_name"`
	Weight        float64 `json:"weight"`
	Calories      float64 `json:"calories"`
	Proteins      float64 `json:"proteins"`
	Fats          float64 `json:"fats"`
	Carbohydrates float64 `json:"carbohydrates"`
	Guid          string  `json:"guid"`
}

type UpdateProductById struct {
	UserId  string             `json:"user_id"`
	Product ProductUpdateModel `json:"product"`
}

type GetAllProductsForUser struct {
	UserId string `json:"user_id"`
	Today  bool   `json:"today"`
}
