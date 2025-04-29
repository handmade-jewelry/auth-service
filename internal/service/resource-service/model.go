package resource_service

import "time"

type Resource struct {
	ID        int
	ServiceID int
	Path      string
	Method    string
	Roles     []string
	IsActive  bool
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAT *time.Time
}
