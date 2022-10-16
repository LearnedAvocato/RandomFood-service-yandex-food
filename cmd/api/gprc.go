package main

import (
	"context"
	"fmt"
	"log"
	"net"
	proto "yandex-food/proto/generated"

	"google.golang.org/grpc"
)

type YandexFoodServer struct {
	proto.UnimplementedYandexFoodServer
}

func (f *YandexFoodServer) GetRandomFood(ctx context.Context, req *proto.FoodRequest) (*proto.FoodResponse, error) {
	cardsNum := req.GetCardsNum()
	latitude := req.GetLatitude()
	longitude := req.GetLongitude()
	getTags := req.GetGetTags()
	selectedTags := req.GetSelectedTags()
	log.Println("GetRandomFood request received")
	return GetRandomFood(int(cardsNum), float64(latitude), float64(longitude), getTags, selectedTags)
}

func (app *Config) gPRCListen() {
	lis, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%s", gRpcPort))
	if err != nil {
		log.Fatalf("Failed to listen for gRPC: %v", err)
	}

	s := grpc.NewServer()

	proto.RegisterYandexFoodServer(s, &YandexFoodServer{})
	log.Println("gRPC server started on port", gRpcPort)

	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to listen for gRPC: %v", err)
	}
}
