package api

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/lets-goo/hotel-reservation/db"
	"github.com/lets-goo/hotel-reservation/types"
)

type HotelHandler struct {
	hotelStore db.HotelStore
	roomStore  db.RoomStore
}

func NewHotelHandler(hotelStore db.HotelStore, roomStore db.RoomStore) *HotelHandler {
	return &HotelHandler{
		hotelStore: hotelStore,
		roomStore:  roomStore,
	}
}

func (h *HotelHandler) HandleGetHotel(c *fiber.Ctx) error {
	id := c.Params("id")
	hotel, err := h.hotelStore.GetHotelByID(c.Context(), id)
	if err != nil {
		return err
	}
	return c.JSON(hotel)
}

type HotelQueryParams struct {
	Rooms  bool
	Rating int
}

func (h *HotelHandler) HandleGetHotels(c *fiber.Ctx) error {
	var qparams HotelQueryParams
	if err := c.QueryParser(&qparams); err != nil {
		return err
	}
	hotels, err := h.hotelStore.GetHotels(c.Context())
	if err != nil {
		return err
	}
	return c.JSON(hotels)
}

func (h *HotelHandler) HandlePostHotel(c *fiber.Ctx) error {
	var params types.Hotel
	if err := c.BodyParser(&params); err != nil {
		return err
	}
	if errors := params.Validate(); len(errors) > 0 {
		return c.JSON(errors)
	}
	hotel, err := types.NewHotelFromParams(params)
	if err != nil {
		return err
	}
	insertedHotel, err := h.hotelStore.CreateHotel(c.Context(), hotel)
	if err != nil {
		return err
	}
	return c.JSON(insertedHotel)
}

func (h *HotelHandler) HandleDeleteHotel(c *fiber.Ctx) error {
	id := c.Params("id")
	err := h.hotelStore.DeleteHotel(c.Context(), id)
	if err != nil {
		return err
	}
	return c.JSON(map[string]string{"deleted": id})
}

func (h *HotelHandler) HandlePutHotel(c *fiber.Ctx) error {
	var (
		params  types.UpdateHotelParams
		hotelID = c.Params("id")
	)
	if err := c.BodyParser(&params); err != nil {
		return err
	}
	res, err := h.hotelStore.UpdateHotel(c.Context(), hotelID, params)
	if err != nil {
		return err
	}
	return c.JSON(map[string]string{"updated": hotelID, "rowCountUpdated": fmt.Sprint(res)})
}
