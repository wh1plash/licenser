package api

import (
	"licenser/server/store"
	"time"

	"github.com/gofiber/fiber/v2"
)

type AppParams struct {
	Name string `json:"app"`
}

type AppHandler struct {
	AppStore store.AppStore
}

func NewAppHandler(s store.AppStore) *AppHandler {
	return &AppHandler{
		AppStore: s,
	}
}

func (h AppHandler) HandleGetApp(c *fiber.Ctx) error {
	var params AppParams
	if err := c.BodyParser(&params); err != nil {
		return ErrBadRequest()
	}

	res, err := h.AppStore.GetApp(c.Context(), params.Name)
	if err != nil {
		return err
	}
	return c.JSON(res)
}

func (h AppHandler) HandleInsertApp(c *fiber.Ctx) error {
	var params AppParams
	if err := c.BodyParser(&params); err != nil {
		return ErrBadRequest()
	}

	app, err := NewAppFromParams(params)
	if err != nil {
		return err
	}

	insApp, err := h.AppStore.InsertApp(c.Context(), app)
	if err != nil {
		return err
	}
	return c.JSON(insApp)
}

func NewAppFromParams(params AppParams) (*store.App, error) {
	return &store.App{
		Name:      params.Name,
		CreatedAt: time.Now(),
		Until:     time.Now().AddDate(0, 1, 0),
	}, nil
}

func (e Error) Error() string {
	return e.Message
}

func NewError(code int, err string) Error {
	return Error{
		Code:    code,
		Message: err,
	}
}

type Error struct {
	Code    int    `json:"code"`
	Message string `json:"error"`
}

func ErrBadRequest() Error {
	return Error{
		Code:    fiber.StatusBadRequest,
		Message: "invalid JSON request",
	}
}
