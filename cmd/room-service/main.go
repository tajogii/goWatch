package main

import (
	"context"
	"fmt"

	"github.com/tajogii/goWatch/cmd/room-service/config"
	"github.com/tajogii/goWatch/internal/pkg/dto"
	roomservice "github.com/tajogii/goWatch/internal/room-service"
	"github.com/tajogii/goWatch/pkg/cache"
	"github.com/tajogii/goWatch/pkg/httpserver"
	"github.com/tajogii/goWatch/pkg/storage"
)

func main() {
	cfg, err := config.LoadConfig()

	if err != nil {
		panic(err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	roomDb, err := storage.NewPgStorage(ctx, cfg.DB.RoomDb)
	if err != nil {
		fmt.Print(err)
	}
	defer roomDb.Close()

	roomStorage := roomservice.NewRoomServiceStorage(roomDb)
	roomCache := cache.NewCache[dto.RoomDto](cfg.Cache.RoomCache.Ttl)
	roomService := roomservice.NewRoomService(roomStorage, roomCache)

	roomHandler := roomservice.NewHandler(roomService)

	httpServer := httpserver.NewHttpServer(roomHandler)

	if httpServer.Listen(fmt.Sprintf(":%d", cfg.General.PublicPort)) != nil {

	}

}
