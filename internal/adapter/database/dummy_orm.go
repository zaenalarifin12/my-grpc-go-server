package database

import (
	"github.com/google/uuid"
	"time"
)

type DummyOrm struct {
	UserId    uuid.UUID `gorm:"primaryKey"`
	UserName  string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (DummyOrm) TableName() string {
	return "dummy"
}
