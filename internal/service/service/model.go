package service

import "time"

// todo rename
type ServiceModel struct {
	Id        int
	Name      string
	IsActive  bool
	Host      string
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}
