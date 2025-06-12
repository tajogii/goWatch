package main

import (
	"context"
	"fmt"

	"github.com/tajogii/goWatch/cmd/room-service/config"
	roomservice "github.com/tajogii/goWatch/internal/room-service"
	"github.com/tajogii/goWatch/pkg/cache"
	"github.com/tajogii/goWatch/pkg/httpserver"
	"github.com/tajogii/goWatch/pkg/storage"
	"go.uber.org/zap"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		panic(err)
	}

	logger, _ := zap.NewProduction()
	defer func() {
		logger.Sync()
	}()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	roomDb, err := storage.NewPgStorage(ctx, cfg.DB.RoomDb)
	if err != nil {
		logger.Panic("failed connect to roomDb", zap.Error(err))
	}
	defer roomDb.Close()

	roomStorage := roomservice.NewRoomServiceStorage(roomDb)
	roomCache := cache.NewCache[roomservice.RoomDto](cfg.Cache.RoomCache.Ttl)
	roomService := roomservice.NewRoomService(roomStorage, roomCache)

	roomHandler := roomservice.NewHandler(roomService)

	httpServer := httpserver.NewHttpServer(logger, roomHandler).Listen(fmt.Sprintf(":%d", cfg.General.PublicPort))

	if httpServer != nil {
		logger.Panic("failed to start http server", zap.Error(err))
	}

}
