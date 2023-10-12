package product

import (
	"fmt"
	"time"
)

const layout = "2021-01-03T02:30:00+01:00"

type Product struct {
	ProductName   string     `json:"product_name"`
	Weight        float64    `json:"weight"`
	Calories      float64    `json:"calories"`
	Proteins      float64    `json:"proteins"`
	Fats          float64    `json:"fats"`
	Carbohydrates float64    `json:"carbohydrates"`
	CreatedAt     CustomDate `json:"created_at"`
	Guid          string     `json:"guid"`
}

type CustomDate struct {
	time.Time
}

func (t CustomDate) MarshalJSON() ([]byte, error) {
	if t.Time.IsZero() {
		return nil, nil
	}
	return []byte(fmt.Sprintf(`"%s"`, t.Time.Format(time.RFC3339))), nil
}
