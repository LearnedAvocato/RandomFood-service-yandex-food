package yandex_food

import desc "yandex-food/pkg/api/yandex-food"

type Implementation struct {
	desc.UnimplementedYandexFoodServer
}

func NewYandexFood() *Implementation {
	return &Implementation{}
}
