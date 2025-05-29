package user

type RoleName string

//const (
//	AdminRole    = RoleName("ADMIN_ROLE")
//	CustomerRole = RoleName("CUSTOMER_ROLE")
//	SellerRole   = RoleName("SELLER_ROLE")
//)

type Role struct {
	ID   int
	Name RoleName
	//CreatedAt time.Time
	//UpdatedAt time.Time
	//DeletedAT time.Time todo нужно ли
}

type UserWithRoles struct {
	UserID int64
	Roles  []string
}
