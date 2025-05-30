package resource

import "time"

type Resource struct {
	ID               int
	ServiceID        int
	PublicPath       string
	ServicePath      string
	Method           string
	Scheme           string
	Roles            []string
	IsActive         bool
	CheckAccessToken bool
	CheckRoles       bool
	CreatedAt        time.Time
	UpdatedAt        time.Time
	DeletedAt        *time.Time
}
