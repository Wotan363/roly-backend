package database

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Email     string    `gorm:"unique;not null"`
	Password  string    `gorm:"not null"`
	CreatedAt time.Time
}

type Role struct {
	ID           uuid.UUID  `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	UserID       *uuid.UUID `gorm:"type:uuid"` // null is for default roles
	Name         string     `gorm:"unique;not null"`
	SystemPrompt string
	CreatedAt    time.Time
}

type Chat struct {
	ID        uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	UserID    uuid.UUID `gorm:"type:uuid;not null"`
	Title     string
	CreatedAt time.Time
}

type Message struct {
	ID             uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	ChatID         uuid.UUID `gorm:"type:uuid;not null"`
	SenderRole     string
	Content        string
	CreatedAt      time.Time
	RoleSnapshotID uuid.UUID `gorm:"type:uuid;not null"`
}

type RoleSnapshot struct {
	ID           uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	RoleID       uuid.UUID `gorm:"type:uuid;not null"`
	Name         string
	SystemPrompt string
	CreatedAt    time.Time
}
