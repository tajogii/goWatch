package roomservice

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/tajogii/goWatch/internal/pkg/dto"
	"github.com/tajogii/goWatch/pkg/cache"
	logm "github.com/tajogii/goWatch/pkg/logger"
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

func (s *RoomService) GetRoomById(ctx context.Context, id uuid.UUID) (*dto.RoomDto, error) {
	logger := logm.GetLogger(ctx)
	v, ok := s.cache.Get(id.String())
	if ok {
		return &v, nil
	}
	logger.Warn(fmt.Sprintf("failed to get room with id from cache: %s", id))

	room, err := s.storage.GetRoomById(ctx, id)
	if err != nil {
		return &dto.RoomDto{}, err
	}

	go func() {
		s.cache.Set(id.String(), *room)
		logger.Info(fmt.Sprintf("set to cache room with id: %s", id))
	}()

	return room, nil

}
