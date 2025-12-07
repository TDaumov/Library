package dto

import "time"

type RegisterDTO struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
	Role     string `json:"role" validate:"required,oneof=admin reader"`
}

type LoginDTO struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type CreateBookDTO struct {
	Title       string     `json:"title" validate:"required"`
	Author      string     `json:"author" validate:"required"`
	Description string     `json:"description"`
	PublishedAt *time.Time `json:"published_at" validate:"omitempty"`
}

type UpdateBookDTO struct {
	Title       *string    `json:"title"`
	Author      *string    `json:"author"`
	Description *string    `json:"description"`
	PublishedAt *time.Time `json:"published_at"`
}
