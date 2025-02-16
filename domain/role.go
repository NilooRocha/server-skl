package domain

type Role struct {
	Name        string   `json:"name"`
	Permissions []string `json:"permissions"`
}

const (
	AdminRole     = "admin"
	ModeratorRole = "moderator"
	UserRole      = "user"
)

var Roles = map[string]Role{
	AdminRole: {
		Name:        AdminRole,
		Permissions: []string{"create_user", "read_user", "read_users", "update_user", "delete_user"},
	},
	ModeratorRole: {
		Name:        ModeratorRole,
		Permissions: []string{"read_user", "read_users", "update_user"},
	},
	UserRole: {
		Name:        UserRole,
		Permissions: []string{"read_user", "update_user"},
	},
}
