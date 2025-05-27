package main

import (
	"context"
	"fmt"
	"time"

	"github.com/tajogii/goWatch/internal/pkg/dto"
	roomservice "github.com/tajogii/goWatch/internal/room-service"
	"github.com/tajogii/goWatch/pkg/cache"
	"github.com/tajogii/goWatch/pkg/httpserver"
	"github.com/tajogii/goWatch/pkg/storage"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	pgConf := storage.PgConf{
		User:              "user",
		Password:          "12345",
		Host:              "localhost",
		DBname:            "postgres",
		MaxConns:          5,
		MinConns:          0,
		MaxConnLifetime:   time.Hour,
		MaxConnIdleTime:   time.Minute * 30,
		HealthCheckPeriod: time.Minute,
	}

	roomDb, err := storage.NewPgStorage(ctx, &pgConf)
	if err != nil {
		fmt.Print(err)
	}
	defer roomDb.Close()

	roomStorage := roomservice.NewRoomServiceStorage(roomDb)
	roomCache := cache.NewCache[dto.RoomDto](time.Hour)
	roomService := roomservice.NewRoomService(roomStorage, roomCache)

	roomHandler := roomservice.NewHandler(roomService)

	httpServer := httpserver.NewHttpServer(roomHandler)

	if httpServer.Listen(":80") != nil {

	}

}
