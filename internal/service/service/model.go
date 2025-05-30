package service

import "time"

const serviceTable = "service"

type ServiceEntity struct {
	ID        int64
	Name      string
	IsActive  bool
	Host      string
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}

type ServiceDTO struct {
	ID       int64
	Name     string
	IsActive bool
	Host     string
}
