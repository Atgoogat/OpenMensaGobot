package domain

import (
	"testing"

	"github.com/Atgoogat/openmensarobot/openmensa"
	"github.com/stretchr/testify/assert"
)

var meals = []openmensa.CanteenMeal{
	{
		Id:    0,
		Name:  "Reis",
		Notes: nil,
		Prices: map[openmensa.PriceType]float32{
			openmensa.PRICE_STUDENT: 3.00,
		},
		Category: "Linie 1",
	},
	{
		Id:    1,
		Name:  "Fisch",
		Notes: nil,
		Prices: map[openmensa.PriceType]float32{
			openmensa.PRICE_STUDENT: 3.00,
		},
		Category: "Linie 1",
	},
	{
		Id:       2,
		Name:     "Nudeln",
		Notes:    []string{"Sind eh immer die besten!"},
		Prices:   map[openmensa.PriceType]float32{},
		Category: "Nudeln",
	},
}

func TestFormatSingleMsg(t *testing.T) {
	msg := mealsToMsg(meals[:1], openmensa.PRICE_STUDENT)

	assert.Equal(t, `Linie 1

Reis
3,00€`,
		msg)

}

func TestFormatMultiMsg(t *testing.T) {
	msg := mealsToMsg(meals, openmensa.PRICE_STUDENT)

	assert.Equal(t, `Linie 1

Reis
3,00€

Fisch
3,00€

Nudeln

Nudeln
(Sind eh immer die besten!)
Keine Preisinformation`, msg)

}
