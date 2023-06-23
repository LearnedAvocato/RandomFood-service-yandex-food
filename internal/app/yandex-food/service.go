package yandex_food

import (
	"yandex-food/internal/app/service"
	desc "yandex-food/pkg/api/yandex-food"
)

type Implementation struct {
	desc.UnimplementedYandexFoodServer

	yandexFoodService service.IYandexFoodService
}

func NewImplementation(yandexFoodService service.IYandexFoodService) *Implementation {
	return &Implementation{yandexFoodService: yandexFoodService}
}
