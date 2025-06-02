package roomservice

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/google/uuid"
	"github.com/tajogii/goWatch/internal/pkg/dto"
	logm "github.com/tajogii/goWatch/pkg/logger"
	"github.com/tajogii/goWatch/pkg/storage"
	"go.uber.org/zap"
)

type RoomStorage struct {
	store *storage.PostgresDB
}

type RoomCro struct {
	Size     int8
	Password string
}

var errNotFound = errors.New("room not found")

const selectroomquery = `SELECT id, size FROM room`
const insertroomquery = `INSERT INTO room (size, password) VALUES ($1,$2)`

func NewRoomServiceStorage(store *storage.PostgresDB) *RoomStorage {
	return &RoomStorage{
		store: store,
	}
}

func (s *RoomStorage) GetAllRooms(ctx context.Context, offset int) (*[]dto.RoomDto, error) {
	var query strings.Builder
	query.WriteString(selectroomquery)
	query.WriteString(" LIMIT $1")
	query.WriteString(" OFFSET $2")

	rows, err := s.store.Query(ctx, query.String(), 100, offset)
	if err != nil {
		log.Fatalf("Query failed: %v", err)
	}
	defer rows.Close()
	rooms := make([]dto.RoomDto, 0, 100)

	for rows.Next() {
		var room dto.RoomDto
		err = rows.Scan(&room.Id, &room.Size)
		if err != nil {
			return &[]dto.RoomDto{}, errNotFound
		}
		rooms = append(rooms, room)
	}

	if err := rows.Err(); err != nil {
		return &[]dto.RoomDto{}, fmt.Errorf("failed to scan row: %v", err)
	}

	return &rooms, nil
}

func (s *RoomStorage) GetRoomById(ctx context.Context, id uuid.UUID) (*dto.RoomDto, error) {
	logger := logm.GetLogger(ctx)

	var query strings.Builder
	query.WriteString(selectroomquery)
	query.WriteString(" WHERE id = $1")

	logger.Info("get room from db query",
		zap.String("sql", query.String()),
		zap.String("param", id.String()),
	)
	row := s.store.QueryRow(ctx, query.String(), id)
	var room dto.RoomDto

	if err := row.Scan(&room.Id, &room.Size); err != nil {
		logger.Warn(fmt.Sprintf("failed to get room with id: %s", id))
		return &dto.RoomDto{}, errNotFound
	}

	logger.Info("get room from db",
		zap.Any("room", room),
	)

	return &room, nil
}

func (s *RoomStorage) CreateRoom(ctx context.Context, room *RoomCro) (uuid.UUID, error) {
	var query strings.Builder
	query.WriteString(insertroomquery)
	query.WriteString(" RETURNING id")

	row := s.store.QueryRow(ctx, query.String(), room.Size, room.Password)
	var id uuid.UUID

	if err := row.Scan(&id); err != nil {
		return uuid.UUID{}, fmt.Errorf("failed to scan row: %v", err)
	}

	return id, nil

}
