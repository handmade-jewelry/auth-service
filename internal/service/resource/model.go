package resource

import "time"

type Scheme string

type Method string

const (
	resourceTable = "resource"

	HTTPSScheme = Scheme("HTTPS")

	GetMethod    = Method("GET")
	PostMethod   = Method("POST")
	PutMethod    = Method("PUT")
	DELETEMethod = Method("DELETE")
)

type Resource struct {
	ID               int64
	ServiceID        int64
	PublicPath       string
	ServicePath      string
	Method           Method
	Scheme           Scheme
	Roles            []string
	IsActive         bool
	CheckAccessToken bool
	CheckRoles       bool
	CreatedAt        time.Time
	UpdatedAt        time.Time
	DeletedAt        *time.Time
}

type ResourceDTO struct {
	ServiceID        int64
	PublicPath       string
	ServicePath      string
	Method           string
	Scheme           string
	Roles            []string
	IsActive         bool
	CheckAccessToken bool
	CheckRoles       bool
}
