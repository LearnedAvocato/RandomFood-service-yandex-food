package main

import (
	"context"
	"log"
	"net"
	"yandex-food/internal/app/repository"
	"yandex-food/internal/app/service"
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

	ctx := context.Background()
	repo, err := repository.NewRepository(ctx)
	if err != nil {
		log.Panicf("failed to init repository: %v", err)
	}
	defer repo.Close(ctx)

	yandexFoodService := service.NewYandexFoodService(repo)
	app := yandex_food.NewImplementation(yandexFoodService)

	desc.RegisterYandexFoodServer(grpcServer, app)

	log.Println("Service yandex-food started")

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
