package model

import "time"

type Task struct {
	Id        int       `json:"id" example:"1"`
	Owner     User      `json:"user"`
	Name      string    `json:"name" example:"string"`
	CreatedAt time.Time `json:"created_at" example:"2024-07-09T18:15:32.579945Z"`
	UpdatedAt time.Time `json:"updated_at" example:"2024-07-09T18:15:32.579945Z"`
	IsActive  bool      `json:"is_active" example:"true"`
	Duration  int       `json:"duration" example:"120"`
} // @name Task
