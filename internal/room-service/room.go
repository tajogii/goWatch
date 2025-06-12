package roomservice

import "github.com/google/uuid"

type RoomDto struct {
	Id       uuid.UUID
	Size     uint
	password string
}

func (r *RoomDto) isSizeValid() bool {
	return r.Size != 0
}

func (r *RoomDto) isPasswordValid() bool {
	return r.password != ""
}

func (r *RoomDto) setId(id uuid.UUID) {
	r.Id = id
}
