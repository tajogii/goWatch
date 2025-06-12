package roomservice

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/tajogii/goWatch/pkg/cache"
	logm "github.com/tajogii/goWatch/pkg/logger"
)

var errNotFound = errors.New("room not found")
var errZeroSize = errors.New("zero size")
var errIncorrectPassword = errors.New("incorrect password")

type Storage interface {
	GetAllRooms(ctx context.Context, offset int) (*[]RoomDto, error)
	GetRoomById(ctx context.Context, id uuid.UUID) (*RoomDto, error)
	CreateRoom(ctx context.Context, room *RoomDto) (*RoomDto, error)
}

type RoomService struct {
	storage Storage
	cache   cache.ICashe[RoomDto]
}

func NewRoomService(storage Storage, cache cache.ICashe[RoomDto]) *RoomService {
	return &RoomService{
		storage: storage,
		cache:   cache,
	}
}

func (s *RoomService) GetRoomById(ctx context.Context, id uuid.UUID) (*RoomDto, error) {
	logger := logm.GetLogger(ctx)
	v, ok := s.cache.Get(id.String())
	if ok {
		return &v, nil
	}
	logger.Warn(fmt.Sprintf("failed to get room with id from cache: %s", id))

	room, err := s.storage.GetRoomById(ctx, id)
	if err != nil {
		return &RoomDto{}, err
	}

	go func() {
		s.cache.Set(id.String(), *room)
		logger.Info(fmt.Sprintf("set to cache room with id: %s", id))
	}()

	return room, nil
}

func (s *RoomService) CreateRoom(ctx context.Context, r *RoomDto) (*RoomDto, error) {
	logger := logm.GetLogger(ctx)
	if !r.isSizeValid() {
		return &RoomDto{}, errZeroSize
	}
	if !r.isPasswordValid() {
		return &RoomDto{}, errIncorrectPassword
	}
	room, err := s.storage.CreateRoom(ctx, r)
	if err != nil {
		return &RoomDto{}, err
	}

	go func() {
		s.cache.Set(room.Id.String(), *r)
		logger.Info(fmt.Sprintf("set to cache room with id: %s", room.Id.String()))
	}()

	return room, nil
}
