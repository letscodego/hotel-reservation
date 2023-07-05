package api

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/lets-goo/hotel-reservation/db"
	"github.com/lets-goo/hotel-reservation/types"
)

type HotelHandler struct {
	store *db.Store
}

func NewHotelHandler(store *db.Store) *HotelHandler {
	return &HotelHandler{
		store: store,
	}
}

func (h *HotelHandler) HandleGetHotel(c *fiber.Ctx) error {
	id := c.Params("id")
	hotel, err := h.store.Hotel.GetHotelByID(c.Context(), id)
	if err != nil {
		return err
	}
	return c.JSON(hotel)
}

type HotelQueryParams struct {
	Rooms  bool
	Rating int
}

func (h *HotelHandler) HandleGetHotelRooms(c *fiber.Ctx) error {
	id := c.Params("id")

	rooms, err := h.store.Room.GetRooms(c.Context(), id)
	if err != nil {
		return err
	}
	return c.JSON(rooms)
}

func (h *HotelHandler) HandleGetHotels(c *fiber.Ctx) error {
	// var qparams HotelQueryParams
	// if err := c.QueryParser(&qparams); err != nil {
	// 	return err
	// }
	hotels, err := h.store.Hotel.GetHotels(c.Context())
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
	insertedHotel, err := h.store.Hotel.CreateHotel(c.Context(), hotel)
	if err != nil {
		return err
	}
	return c.JSON(insertedHotel)
}

func (h *HotelHandler) HandleDeleteHotel(c *fiber.Ctx) error {
	id := c.Params("id")
	err := h.store.Hotel.DeleteHotel(c.Context(), id)
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
	res, err := h.store.Hotel.UpdateHotel(c.Context(), hotelID, params)
	if err != nil {
		return err
	}
	return c.JSON(map[string]string{"updated": hotelID, "rowCountUpdated": fmt.Sprint(res)})
}
