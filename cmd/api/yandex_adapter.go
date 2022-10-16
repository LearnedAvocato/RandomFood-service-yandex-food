package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"
	proto "yandex-food/proto/generated"

	"github.com/Jeffail/gabs"
)

const (
	apiUrl            = "https://eda.yandex/api/v2"
	baseRestarauntUrl = "https://eda.yandex.ru/api/v2/catalog"
	baseUrl           = "https://eda.yandex.ru"
)

type void struct{}

var setMember void

type restarauntData struct {
	id   string
	url  string
	name string
	slug string
}

func doGetRequest(url string) (*gabs.Container, error) {
	client := &http.Client{
		Timeout: 1 * time.Second,
	}

	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)

	jsonParsed, err := gabs.ParseJSON(body)
	if err != nil {
		return nil, err
	}

	return jsonParsed, err
}

func extractTags(data *gabs.Container) ([]string, error) {
	tagsSet := make(map[string]void)
	tags := make([]string, 0, 100)

	dataLen, err := data.ArrayCount()
	if err != nil {
		return nil, err
	}

	//var wg sync.WaitGroup
	//wg.Add(dataLen)

	for i := 0; i < dataLen; i++ {
		//go func(*map[string]void, *gabs.Container, int) {
		el, err := data.ArrayElement(i)
		if err != nil {
			continue
		}

		elementTags := el.Search("place", "tags")

		tagsLen, err := elementTags.ArrayCount()
		if err != nil {
			continue
		}

		for i := 0; i < tagsLen; i++ {
			tag, err := elementTags.ArrayElement(i)
			if err != nil {
				continue
			}
			tagStr, ok := tag.Search("name").Data().(string)
			if !ok {
				continue
			}
			tagStr = strings.ToLower(tagStr)

			if _, ok := tagsSet[tagStr]; !ok {
				tagsSet[tagStr] = setMember
				tags = append(tags, tagStr)
			}
		}
		//wg.Done()

		//}(&tags, data, i)
	}
	//wg.Wait()

	return tags, nil
}

func getRestaraunts(latitude, longitude float64, num int, getTags bool, selectedTags []string) ([]restarauntData, []string, error) {
	log.Println("request restaraunts list")
	u, _ := url.JoinPath(apiUrl, "catalog")
	urlStruct, _ := url.Parse(u)
	q := urlStruct.Query()
	q.Set("latitude", fmt.Sprintf("%f", latitude))
	q.Set("longitude", fmt.Sprintf("%f", longitude))
	urlStruct.RawQuery = q.Encode()
	data, err := doGetRequest(urlStruct.String())
	if err != nil {
		return nil, nil, err
	}
	data = data.Search("payload", "foundPlaces")
	total, err := data.ArrayCount()
	if err != nil {
		return nil, nil, err
	}

	tags := make([]string, 0)
	if getTags {
		tags, err = extractTags(data)
		if err == nil {
		} else {
			log.Printf("failed to extract tags: %v", err)
		}
	}

	restaraunts := make([]restarauntData, 0, num)

	for _, i := range rand.Perm(total)[:num] {
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

	return doGetRequest(urlStruct.String())
}

func buildRestarauntUrl(slug string) string {
	u, _ := url.JoinPath(baseRestarauntUrl, slug)
	return u
}

func extractRestarauntData(data *gabs.Container) restarauntData {
	id, _ := data.Search("place", "id").Data().(float64)
	name, _ := data.Search("place", "name").Data().(string)
	slug, _ := data.Search("place", "slug").Data().(string)

	return restarauntData{
		id:   fmt.Sprintf("%d", int64(id)),
		name: name,
		url:  buildRestarauntUrl(slug),
		slug: slug}
}

func createFoodCard(data *gabs.Container, resData restarauntData) *proto.FoodCard {
	foodCard := proto.FoodCard{}
	foodCard.Name, _ = data.Search("name").Data().(string)
	foodCard.ImageUrl, _ = data.Search("picture", "uri").Data().(string)
	foodCard.ImageUrl = baseUrl + foodCard.ImageUrl
	foodCard.Description, _ = data.Search("description").Data().(string)
	foodCard.Id = fmt.Sprintf("%d", int64(data.Search("id").Data().(float64)))
	foodCard.RestarauntName = resData.name
	foodCard.Price = float32(data.Search("price").Data().(float64))
	foodCard.RestarauntUrl = resData.url

	return &foodCard
}

// todo: process categories and filters
func extractRandomDish(menu *gabs.Container) (*gabs.Container, error) {
	categories, err := menu.Search("payload", "categories").Children()
	if err != nil {
		return nil, err
	}
	if len(categories) == 0 {
		return nil, errors.New("empty categories")
	}
	randomCat := categories[rand.Intn(len(categories))]

	dishes, err := randomCat.Search("items").Children()
	if err != nil {
		return nil, err
	}
	if len(dishes) == 0 {
		return nil, errors.New("empty dishes")
	}
	return dishes[rand.Intn(len(dishes))], nil
}

func checkDishData(data *gabs.Container) error {
	if !data.Exists("picture", "uri") {
		return errors.New("dish has no image")
	}
	return nil
}

func GetRandomFood(cardsNum int, latitude, longitude float64, getTags bool, selectedTags []string) (*proto.FoodResponse, error) {
	restarauntDataArr, tags, err := getRestaraunts(latitude, longitude, cardsNum, getTags, selectedTags)
	if err != nil {
		log.Printf("failed to get restaraunts: %v", err)
		return &proto.FoodResponse{Succeed: false}, err
	}
	if !getTags || tags == nil {
		tags = make([]string, 0)
	}

	// Take cardsNum food cards
	// 1 restaraunt - 1 card
	foodCards := make([]*proto.FoodCard, 0, cardsNum)
	var wg sync.WaitGroup
	wg.Add(cardsNum)
	for _, data := range restarauntDataArr {

		go func(data restarauntData, foodCards *[]*proto.FoodCard) {
			defer wg.Done()

			menu, err := getRestarauntMenu(data.slug)
			if err != nil {
				log.Printf("failed to get restaraunt menu for restaraunt with id %s and slug %s: %v", data.id, data.slug, err)
				return
			}

			dish, err := extractRandomDish(menu)
			//log.Println(dish.String())
			if err != nil {
				log.Printf("failed to get dishes for restaraunt with id %s and slug %s: %v", data.id, data.slug, err)
				return
			}

			if err := checkDishData(dish); err != nil {
				log.Printf("invalid dish data: %v", err)
				return
			}

			*foodCards = append(*foodCards, createFoodCard(dish, data))
		}(data, &foodCards)
	}
	wg.Wait()
	log.Printf("got %d cards when %d were requested", len(foodCards), cardsNum)

	return &proto.FoodResponse{
		Succeed:       true,
		FoodCards:     foodCards,
		AvailableTags: tags,
	}, nil
}
