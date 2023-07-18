package config

import (
	"sync"

	"github.com/Atgoogat/openmensarobot/domain"
	"github.com/Atgoogat/openmensarobot/openmensa"
)

var (
	mealService     *domain.MealService
	mealServiceOnce sync.Once
)

func GetMealService(openmensaApi openmensa.OpenmensaApi) *domain.MealService {
	mealServiceOnce.Do(func() {
		mealService = domain.NewMealService(openmensaApi)
	})
	return mealService
}
