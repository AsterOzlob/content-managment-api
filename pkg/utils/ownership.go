package utils

// IsOwner проверяет, является ли пользователь владельцем ресурса.
// Если роль пользователя — модератор или администратор, доступ разрешён без проверки владельца.
func IsOwner(resourceOwnerID uint, userID uint, roles []string) bool {
	// Проверяем, есть ли у пользователя роли, позволяющие игнорировать проверку владельца
	for _, role := range roles {
		if role == "moderator" || role == "admin" {
			return true
		}
	}

	// Проверяем, является ли пользователь владельцем
	return resourceOwnerID == userID
}
