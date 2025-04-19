package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Username     string `gorm:"unique;not null;size:64"`
	Email        string `gorm:"unique;size:255"`
	PasswordHash string `gorm:"not null"`
	Roles        []Role `gorm:"many2many:user_roles;"`
}

type UserRole struct {
	UserID uint `gorm:"primaryKey"`
	RoleID uint `gorm:"primaryKey"`
}

func (u *User) Can(resource string, permission PermissionBits) bool {
	for _, role := range u.Roles {
		if role.Can(resource, permission) {
			return true
		}
	}
	return false
}

// HasRole проверяет наличие роли
func (u *User) HasRole(roleName string) bool {
	for _, role := range u.Roles {
		if role.Name == roleName {
			return true
		}
	}
	return false
}
