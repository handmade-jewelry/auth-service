package route

type Route struct {
	Host             string
	PublicPath       string
	ServicePath      string
	Method           string
	Scheme           string
	Roles            []string
	CheckAccessToken bool
	CheckRoles       bool
}
