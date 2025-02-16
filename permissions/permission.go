package permissions

import "server/domain"

func Can(userRole domain.Role, action, resource string) bool {
	for _, permission := range userRole.Permissions {
		if permission == action+"_"+resource {
			return true
		}
	}
	return false
}
