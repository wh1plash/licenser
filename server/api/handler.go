package api

import (
	"fmt"
	"licenser/server/store"
	"licenser/server/types"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type AppParams struct {
	Name string `json:"app" validate:"required,min=2"`
}

type AppHandler struct {
	AppStore store.AppStore
}

func NewAppHandler(s store.AppStore) *AppHandler {
	return &AppHandler{
		AppStore: s,
	}
}

func (h AppHandler) HandleGetAppList(c *fiber.Ctx) error {
	res, err := h.AppStore.GetAppList(c.Context())
	if err != nil {
		return err
	}

	return c.JSON(res)
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

	validate := validator.New()
	if err := validate.Struct(params); err != nil {
		errs := err.(validator.ValidationErrors)
		errors := make(map[string]string)
		for _, e := range errs {
			errors[e.Field()] = fmt.Sprintf("failed on '%s' tag", e.Tag())
		}
		Err := NewValidationError(errors)
		return c.Status(Err.Status).JSON(Err)
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

func NewAppFromParams(params AppParams) (*types.App, error) {
	return &types.App{
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

type ValidationError struct {
	Status int               `json:"status"`
	Errors map[string]string `json:"errors"`
}

func (e ValidationError) Error() string {
	return "validation failed"
}

func NewValidationError(errors map[string]string) ValidationError {
	return ValidationError{
		Status: fiber.StatusUnprocessableEntity,
		Errors: errors,
	}
}
