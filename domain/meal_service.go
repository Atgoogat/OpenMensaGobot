package domain

import (
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/Atgoogat/openmensarobot/openmensa"
)

const (
	MealLimitPerMessage = 50 
)

var (
	ErrCouldNotFetch = errors.New("could not fetch meals from openmensa api")
)

type MealService struct {
	openmensaApi openmensa.OpenmensaApi
	formatter    []TextFormatter
}

func NewMealService(openmensaApi openmensa.OpenmensaApi, formatter ...TextFormatter) MealService {
	return MealService{
		openmensaApi: openmensaApi,
		formatter:    formatter,
	}
}

type TextFormatter interface {
	Format(text string) (string, error)
}

func (s MealService) GetFormatedMeals(mensaID int, date time.Time, priceType openmensa.PriceType) (string, error) {
	meals, err := s.openmensaApi.ListMealsForADay(mensaID, date, 1, MealLimitPerMessage)
	if err != nil {
		return "", errors.Join(ErrCouldNotFetch, err)
	}

	t := mealsToMsg(meals, priceType)
	return s.applyFormatterChain(t), nil
}

func (s MealService) applyFormatterChain(text string) string {
	for _, f := range s.formatter {
		newText, err := f.Format(text)
		if err != nil {
			log.Printf("error while applying formatter: %v", err)
		} else {
			text = newText
		}
	}
	return text
}

func mealsToMsg(meals []openmensa.CanteenMeal, priceType openmensa.PriceType) string {
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
			catMsg = append(catMsg, "<b>"+category+"</b>")

			for _, m := range meals {
				catMsg = append(catMsg, "", m.Name)
				var filteredNotes []string
				for _, note := range m.Notes {
					note = strings.Trim(note, " \t\n")
					filteredNotes = append(filteredNotes, note)
				}

				notes := strings.Join(filteredNotes, ", ")
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
