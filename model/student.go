package model

import (
	"time"
)

type Student struct {
	ID            uint
	Name          string
	StudentTypeId uint
	GuardianName  string
	CreatedAt     time.Time
	UpdatedAt     time.Time
}
