package entities

import "time"

type Reaction struct {
	ID        int64
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
}
