package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type User struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	Email     string    `gorm:"uniqueIndex;not null" json:"email"`
	Password  string    `gorm:"not null" json:"-"`
	Role      string    `gorm:"not null;default:'reader'" json:"role"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	u.ID = uuid.New()
	return
}

type Book struct {
	ID          uuid.UUID  `gorm:"type:uuid;primaryKey" json:"id"`
	Title       string     `gorm:"not null" json:"title" validate:"required"`
	Author      string     `json:"author" validate:"required"`
	Description string     `json:"description"`
	PublishedAt *time.Time `json:"published_at"`
	Available   bool       `gorm:"default:true" json:"available"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (b *Book) BeforeCreate(tx *gorm.DB) (err error) {
	b.ID = uuid.New()
	return
}

type Borrow struct {
	ID         uuid.UUID  `gorm:"type:uuid;primaryKey" json:"id"`
	UserID     uuid.UUID  `gorm:"type:uuid;index" json:"user_id"`
	BookID     uuid.UUID  `gorm:"type:uuid;index" json:"book_id"`
	BorrowedAt time.Time  `json:"borrowed_at"`
	ReturnedAt *time.Time `json:"returned_at"`
}

func (b *Borrow) BeforeCreate(tx *gorm.DB) (err error) {
	b.ID = uuid.New()
	return
}
