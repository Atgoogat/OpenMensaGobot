package openmensa

type PriceType string

const (
	PRICE_STUDENT   PriceType = "students"
	PRICE_EMPLOYEES PriceType = "employees"
	PRICE_PUPILS    PriceType = "pupils"
	PRICE_OTHERS    PriceType = "others"
	PRICE_NONE      PriceType = "none"
)

type Canteen struct {
	Id          int       `json:"id"`
	Name        string    `json:"name"`
	City        string    `json:"city"`
	Address     string    `json:"address"`
	Coordinates []float64 `json:"coordinates"`
}

type CanteenDay struct {
	Date   string `json:"date"`
	Closed bool   `json:"closed"`
}

type CanteenMeal struct {
	Id       int                   `json:"id"`
	Name     string                `json:"name"`
	Notes    []string              `json:"notes"`
	Prices   map[PriceType]float32 `json:"prices"`
	Category string                `json:"category"`
}
