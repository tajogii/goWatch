package dto

import "github.com/google/uuid"

type RoomDto struct {
	Id   uuid.UUID
	Size int8
}
