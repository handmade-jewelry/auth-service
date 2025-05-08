package service

import "time"

type ServiceModel struct {
	Id        int
	Name      string
	IsActive  bool
	Host      string
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}
