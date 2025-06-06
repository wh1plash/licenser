package types

import "time"

type App struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"createdAt"`
	Until     time.Time `json:"until"`
}

type AppParams struct {
	Name string `json:"app" validate:"required,min=2"`
}

func NewAppFromParams(params AppParams) (*App, error) {
	return &App{
		Name:      params.Name,
		CreatedAt: time.Now(),
		Until:     time.Now().AddDate(0, 1, 0),
	}, nil
}
