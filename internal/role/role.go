package role

type RoleType uint32

const (
	Root RoleType = iota + 1
	User
)
