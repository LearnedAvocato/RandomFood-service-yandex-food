package service

import (
	"context"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"net/url"
	"strings"
	"sync"
	"time"

	"yandex-food/internal/pkg/datastruct"
	"yandex-food/internal/pkg/network"

	"github.com/Jeffail/gabs"
)

type IYandexFoodService interface {
	GetRandomFood(ctx context.Context, cardsNum int, latitude, longitude float64, getTags bool, selectedTags []string) ([]*datastruct.FoodCard, error)
}

type IRepository interface {
	LogRequest(ctx context.Context, log *datastruct.RequestLog) error
}

type yandexFoodService struct {
	repo IRepository
}

const (
	apiUrl                  = "https://eda.yandex/api/v2"
	baseRestarauntUrl       = "https://eda.yandex.ru/api/v2/catalog"
	baseUrl                 = "https://eda.yandex.ru"
	baseRestarauntPublicUrl = "https://eda.yandex/restaurant/"
	searchByTextUrl         = "https://eda.yandex.ru/eats/v1/full-text-search/v1/search"
)

const (
	maxCardsNum = 30
)

func NewYandexFoodService(repo IRepository) IYandexFoodService {
	return &yandexFoodService{repo: repo}
}

func (yfs *yandexFoodService) GetRandomFood(ctx context.Context, cardsNum int, latitude, longitude float64, getTags bool, selectedTags []string) ([]*datastruct.FoodCard, error) {
	if cardsNum > maxCardsNum {
		cardsNum = maxCardsNum
		log.Printf("got cardsNum %d, set it to %d\n", cardsNum, maxCardsNum)
	}

	requestStartTime := time.Now()
	log.Printf("User coordinates: {%f, %f}\n", latitude, longitude)

	restarauntDataArr, tags, err := getRestaraunts(latitude, longitude, cardsNum, getTags, selectedTags)
	if err != nil {
		log.Printf("failed to get restaraunts: %v", err)
		return nil, err
	}
	if !getTags || tags == nil {
		tags = make([]string, 0)
	}

	// Take cardsNum food cards
	// 1 restaraunt - X card

	cardsPerRestaraunt := cardsNum / len(restarauntDataArr)
	log.Printf("Request %d cards per restaraunt, %d in total\n", cardsPerRestaraunt, cardsPerRestaraunt*len(restarauntDataArr))

	foodCards := make([]*datastruct.FoodCard, 0, cardsNum)
	var wg sync.WaitGroup
	wg.Add(len(restarauntDataArr))
	for _, data := range restarauntDataArr {
		go func(data datastruct.Restaraunt, foodCards *[]*datastruct.FoodCard) {
			defer wg.Done()

			menu, err := getRestarauntMenu(data.Slug)
			if err != nil {
				log.Printf("failed to get restaraunt menu for restaraunt with id %s and slug %s: %v", data.ID, data.Slug, err)
				return
			}

			dishes, err := extractRandomDishes(menu, cardsPerRestaraunt)
			if err != nil {
				log.Printf("failed to get dishes for restaraunt with id %s and slug %s: %v", data.ID, data.Slug, err)
				return
			}

			for _, dish := range dishes {
				if err := checkDishData(dish); err != nil {
					log.Printf("invalid dish data: %v", err)
					continue
				}

				*foodCards = append(*foodCards, createFoodCard(dish, data))
			}
		}(data, &foodCards)
	}
	wg.Wait()
	log.Printf("got %d cards when %d were requested", len(foodCards), cardsNum)

	foodCardsPermuted := make([]*datastruct.FoodCard, 0, len(foodCards))
	for _, i := range rand.Perm(len(foodCards)) {
		foodCardsPermuted = append(foodCardsPermuted, foodCards[i])
	}

	err = yfs.repo.LogRequest(ctx, &datastruct.RequestLog{RequestedCardsNum: int64(cardsNum),
		GotCardsNum:        int64(len(foodCards)),
		Longitude:          longitude,
		Latitude:           latitude,
		UsedRestarauntsNum: int64(len(restarauntDataArr)),
		CardsPerRestaraunt: int64(cardsPerRestaraunt),
		CreatedAt:          requestStartTime})

	return foodCardsPermuted, err
}

func getRestaraunts(latitude, longitude float64, num int, getTags bool, selectedTags []string) ([]datastruct.Restaraunt, []string, error) {
	log.Println("request restaraunts list")
	u, _ := url.JoinPath(apiUrl, "catalog")
	urlStruct, _ := url.Parse(u)
	q := urlStruct.Query()
	q.Set("latitude", fmt.Sprintf("%f", latitude))
	q.Set("longitude", fmt.Sprintf("%f", longitude))
	urlStruct.RawQuery = q.Encode()

	payload := gabs.New()
	payload.SetP(latitude, "location.latitude")
	payload.SetP(longitude, "location.longitude")
	payload.SetP("еда", "text")

	data, err := network.DoPostRequest(searchByTextUrl, payload.String())
	if err != nil {
		return nil, nil, err
	}

	data, err = data.Search("blocks").ArrayElement(0)
	if err != nil {
		return nil, nil, err
	}

	data = data.Search("payload")
	total, err := data.ArrayCount()
	if err != nil {
		return nil, nil, err
	}

	tags := make([]string, 0)

	respCardsNum := num
	if total < num {
		respCardsNum = total
	}
	log.Printf("%d cards requested, %d cards got from api, return %d cards", num, total, respCardsNum)
	restaraunts := make([]datastruct.Restaraunt, 0, respCardsNum)

	for _, i := range rand.Perm(total)[:respCardsNum] {
		el, err := data.ArrayElement(i)
		if err != nil {
			continue
		}

		restaraunts = append(restaraunts, extractRestarauntData(el))
	}
	return restaraunts, tags, err
}

func getRestarauntMenu(slug string) (*gabs.Container, error) {
	log.Println("request restaraunt menu")
	u, _ := url.JoinPath(buildRestarauntUrl(slug), "menu")
	urlStruct, _ := url.Parse(u)
	//q := urlStruct.Query()
	//q.Set("data", "products")
	//urlStruct.RawQuery = q.Encode()

	return network.DoGetRequest(urlStruct.String())
}

func buildRestarauntUrl(slug string) string {
	u, _ := url.JoinPath(baseRestarauntUrl, slug)
	return u
}

func extractRestarauntData(data *gabs.Container) datastruct.Restaraunt {
	id := 0 //data.Search("place", "id").Data().(float64)
	name, _ := data.Search("title").Data().(string)
	slug, _ := data.Search("slug").Data().(string)

	return datastruct.Restaraunt{
		ID:   fmt.Sprintf("%d", int64(id)),
		Name: name,
		Url:  buildRestarauntUrl(slug),
		Slug: slug}
}

func fixDescription(desc *string) {
	if desc == nil {
		*desc = ""
	}

	i := strings.Index(*desc, "<br>")
	if i >= 0 {
		*desc = (*desc)[:i]
	}
}

func createFoodCard(data *gabs.Container, resData datastruct.Restaraunt) *datastruct.FoodCard {
	foodCard := datastruct.FoodCard{}
	foodCard.Name, _ = data.Search("name").Data().(string)
	foodCard.ImageUrl, _ = data.Search("picture", "uri").Data().(string)
	foodCard.ImageUrl = baseUrl + foodCard.ImageUrl
	foodCard.Description, _ = data.Search("description").Data().(string)
	fixDescription(&foodCard.Description)
	foodCard.Id = fmt.Sprintf("%d", int64(data.Search("id").Data().(float64)))
	foodCard.RestarauntName = resData.Name
	foodCard.Price = float32(data.Search("price").Data().(float64))
	foodCard.RestarauntUrl = baseRestarauntPublicUrl + resData.Slug

	return &foodCard
}

// todo: process categories and filters
func extractRandomDishes(menu *gabs.Container, dishNum int) ([]*gabs.Container, error) {
	dishesData := make([]*gabs.Container, 0, dishNum)
	for i := 0; i < dishNum; i++ {
		categories, err := menu.Search("payload", "categories").Children()
		if err != nil {
			continue
		}
		if len(categories) == 0 {
			continue
		}
		randomCat := categories[rand.Intn(len(categories))]

		dishes, err := randomCat.Search("items").Children()
		if err != nil {
			continue
		}
		if len(dishes) == 0 {
			continue
		}

		dishesData = append(dishesData, dishes[rand.Intn(len(dishes))])
	}
	return dishesData, nil
}

func checkDishData(data *gabs.Container) error {
	if !data.Exists("picture", "uri") {
		return errors.New("dish has no image")
	}
	return nil
}
