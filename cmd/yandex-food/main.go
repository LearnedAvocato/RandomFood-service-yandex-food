package main

import (
	"log"
	"net"
	yandex_food "yandex-food/internal/app/yandex-food"
	desc "yandex-food/pkg/api/yandex-food"

	"google.golang.org/grpc"
)

type Config struct{}

func main() {
	lis, err := net.Listen("tcp", "0.0.0.0:50001")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	yandexFood := yandex_food.NewYandexFood()
	desc.RegisterYandexFoodServer(grpcServer, yandexFood)

	log.Println("Service yandex-food started")

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}

}
