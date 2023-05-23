package yandex_food

import (
	"context"
	"log"
	desc "yandex-food/pkg/api/yandex-food"
)

func (i *Implementation) GetRandomFood(ctx context.Context, req *desc.FoodRequest) (*desc.FoodResponse, error) {
	log.Println("GetRandomFood called")

	cardsNum := req.GetCardsNum()
	latitude := req.GetLatitude()
	longitude := req.GetLongitude()
	getTags := req.GetGetTags()
	selectedTags := req.GetSelectedTags()
	log.Println("GetRandomFood request received")
	return getRandomFood(int(cardsNum), float64(latitude), float64(longitude), getTags, selectedTags)
}
