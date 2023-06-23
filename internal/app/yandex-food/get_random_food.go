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
	foodCards, err := i.yandexFoodService.GetRandomFood(ctx, int(cardsNum), float64(latitude), float64(longitude), getTags, selectedTags)
	if err != nil {
		return &desc.FoodResponse{
			Succeed: false,
		}, err
	}

	res := &desc.FoodResponse{
		Succeed: true,
	}

	for _, foodCard := range foodCards {
		res.FoodCards = append(res.FoodCards, &desc.FoodCard{
			Name:           foodCard.Name,
			ImageUrl:       foodCard.ImageUrl,
			Description:    foodCard.Description,
			Id:             foodCard.Id,
			RestarauntName: foodCard.RestarauntName,
			Price:          foodCard.Price,
			RestarauntUrl:  foodCard.RestarauntUrl,
		})
	}

	return res, nil
}
