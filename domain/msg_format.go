package domain

import (
	"fmt"
	"strings"

	"github.com/Atgoogat/openmensarobot/openmensa"
)

func MealsToMsg(meals []openmensa.CanteenMeal, priceType openmensa.PriceType) string {
	categories := make(map[string][]openmensa.CanteenMeal)
	for _, m := range meals {
		categories[m.Category] = append(categories[m.Category], m)
	}

	categoriesDone := make(map[string]struct{})

	var msg []string

	// process categorys in order of meals
	for _, meal := range meals {
		if _, ok := categoriesDone[meal.Category]; !ok {
			category := meal.Category
			meals := categories[category]

			var catMsg []string
			catMsg = append(catMsg, category)

			for _, m := range meals {
				catMsg = append(catMsg, "", m.Name)

				notes := strings.Join(m.Notes, ", ")
				if notes != "" {
					catMsg = append(catMsg, "("+notes+")")
				}

				if priceType != openmensa.PRICE_NONE {
					price, ok := m.Prices[priceType]
					if ok {
						catMsg = append(catMsg, strings.Replace(fmt.Sprintf("%.2f€", price), ".", ",", 1))
					} else {
						catMsg = append(catMsg, "Keine Preisinformation")
					}
				}
			}

			msg = append(msg, strings.Join(catMsg, "\n"))
			// mark category as done
			categoriesDone[category] = struct{}{}
		}
	}

	return strings.Join(msg, "\n\n")
}