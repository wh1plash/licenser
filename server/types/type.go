package types

import "time"

type App struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"createdAt"`
	Until     time.Time `json:"until"`
}
