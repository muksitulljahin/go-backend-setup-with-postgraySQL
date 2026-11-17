package enum

type Type string

const (
	SuperAdmin Type = "super_admin"
	Admin      Type = "admin"
	Support    Type = "support"
	Client     Type = "client"
)
