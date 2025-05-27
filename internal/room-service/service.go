package roomservice

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/tajogii/goWatch/internal/pkg/dto"
	"github.com/tajogii/goWatch/pkg/cache"
)

type Storage interface {
	GetAllRooms(ctx context.Context, offset int) (*[]dto.RoomDto, error)
	GetRoomById(ctx context.Context, id uuid.UUID) (*dto.RoomDto, error)
	CreateRoom(ctx context.Context, room *RoomCro) (uuid.UUID, error)
}

type RoomService struct {
	storage Storage
	cache   cache.ICashe[dto.RoomDto]
}

func NewRoomService(storage Storage, cache cache.ICashe[dto.RoomDto]) *RoomService {
	return &RoomService{
		storage: storage,
		cache:   cache,
	}
}

func (s *RoomService) GetRoomById(id uuid.UUID) (*dto.RoomDto, error) {
	v, ok := s.cache.Get(id.String())
	if ok {
		return &v, nil
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*1)
	defer cancel()
	room, err := s.storage.GetRoomById(ctx, id)
	if err != nil {
		return &dto.RoomDto{}, fmt.Errorf("cannot get room by id = %s: %v ", id.String(), err)
	}

	go func() {
		s.cache.Set(id.String(), *room)
	}()

	return room, nil

}
