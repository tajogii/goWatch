package roomservice

import (
	"context"
	"errors"
	"fmt"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
	"github.com/tajogii/goWatch/internal/pkg/dto"
	"github.com/tajogii/goWatch/pkg/httpserver"
)

type roomService interface {
	GetRoomById(ctx context.Context, id uuid.UUID) (*dto.RoomDto, error)
}

type Handler struct {
	roomService roomService
}

func NewHandler(roomService roomService) *Handler {
	return &Handler{
		roomService: roomService,
	}
}

func (h *Handler) RegisterRoutes(api fiber.Router) {
	g := api.Group("/room")
	g.Get("/:room", h.getRoomById)
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

	return c.JSON(room)
}
