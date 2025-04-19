package models

import "gorm.io/gorm"

type Role struct {
	gorm.Model
	Name        string `gorm:"unique;not null;size:64"`
	Description string
	Users       []User       `gorm:"many2many:user_roles;"`
	Permissions []Permission `gorm:"foreignKey:RoleID"`
}

// Can проверяет, есть ли у роли указанное право для ресурса
func (r *Role) Can(resource string, permission PermissionBits) bool {
	for _, perm := range r.Permissions {
		if perm.Resource == resource && perm.Rights.HasPermission(permission) {
			return true
		}
	}
	return false
}

// AddPermission добавляет новое право
func (r *Role) AddPermission(resource string, permission PermissionBits) error {
	for i, perm := range r.Permissions {
		if perm.Resource == resource {
			r.Permissions[i].Rights.AddPermission(permission)
			return nil
		}
	}

	r.Permissions = append(r.Permissions, Permission{
		RoleID:   r.ID,
		Resource: resource,
		Rights:   permission,
	})
	return nil
}
