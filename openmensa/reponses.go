package openmensa

const (
	PRICE_STUDENT   = "students"
	PRICE_EMPLOYEES = "employees"
	PRICE_PUPILS    = "pupils"
	PRICE_OTHERS    = "others"
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
	Id       int                `json:"id"`
	Name     string             `json:"name"`
	Notes    []string           `json:"notes"`
	Prices   map[string]float32 `json:"prices"`
	Category string             `json:"category"`
}
