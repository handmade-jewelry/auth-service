package service_service

import "time"

type Service struct {
	Id        int
	Name      string
	IsActive  bool
	Host      string
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}
