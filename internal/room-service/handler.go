package roomservice

import (
	"context"
	"errors"
	"fmt"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
	"github.com/tajogii/goWatch/pkg/httpserver"
)

type roomService interface {
	GetRoomById(ctx context.Context, id uuid.UUID) (*RoomDto, error)
	CreateRoom(ctx context.Context, room *RoomDto) (*RoomDto, error)
}

type Handler struct {
	roomService roomService
}

type roomCro struct {
	Size     uint   `json:"size"`
	Password string `json:"password"`
}

type room struct {
	Id   uuid.UUID `json:"id"`
	Size uint      `json:"size"`
}

func NewHandler(roomService roomService) *Handler {
	return &Handler{
		roomService: roomService,
	}
}

func (h *Handler) RegisterRoutes(api fiber.Router) {
	g := api.Group("/room")
	g.Get("/:room", h.getRoomById)
	g.Post("/", h.createRoom)
}

func (h *Handler) getRoomById(c fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("room"))

	if err != nil {
		return c.Status(400).JSON(httpserver.Error{
			Message:     fmt.Sprintf("invalid room id: %s", c.Params("room")),
			Description: "room id must be uuid",
		})
	}

	room, err := h.roomService.GetRoomById(c.Context(), id)
	if err != nil {
		if errors.Is(err, errNotFound) {
			return c.Status(404).JSON(httpserver.Error{
				Message:     "room not found",
				Description: fmt.Sprintf("room with id: %s, doesnt exist", c.Params("room")),
			})
		}
		return err
	}

	return c.JSON(mapRoomRes(room))
}

func (h *Handler) createRoom(c fiber.Ctx) error {
	r := new(roomCro)
	if err := c.Bind().Body(r); err != nil {
		return c.Status(400).JSON(httpserver.Error{
			Message:     "invalid request",
			Description: "size must be positive integer and password string",
		})
	}

	room, err := h.roomService.CreateRoom(c.Context(), mapRoom(r))
	if err != nil {
		if errors.Is(err, errZeroSize) {
			return c.Status(400).JSON(httpserver.Error{
				Message:     "room size is zero",
				Description: "room size must be bigger then 0",
			})
		}
		if errors.Is(err, errIncorrectPassword) {
			return c.Status(400).JSON(httpserver.Error{
				Message:     "incorrect password",
				Description: "incorrect password",
			})
		}
		return err
	}

	return c.JSON(mapRoomRes(room))
}

func mapRoom(room *roomCro) *RoomDto {
	return &RoomDto{
		Size:     room.Size,
		password: room.Password,
	}
}

func mapRoomRes(r *RoomDto) *room {
	return &room{
		Size: r.Size,
		Id:   r.Id,
	}
}
