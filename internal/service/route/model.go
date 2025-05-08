package route

type Route struct {
	Host             string
	Path             string
	Method           string
	Scheme           string
	Roles            []string
	CheckAccessToken bool
	CheckRoles       bool
}
