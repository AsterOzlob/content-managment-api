package models

import (
	"strings"

	"gorm.io/gorm"
)

// PermissionBits - тип для битовой маски прав
type PermissionBits int

// Константы для прав доступа (битовая маска)
const (
	PermissionRead   PermissionBits                                        = 1 << iota // 0001 (1)
	PermissionWrite                                                                    // 0010 (2)
	PermissionDelete                                                                   // 0100 (4)
	PermissionAll    = PermissionRead | PermissionWrite | PermissionDelete             // 0111 (7)
)

// Permission - модель для хранения прав в базе данных
type Permission struct {
	gorm.Model
	RoleID   uint           `gorm:"not null;index"`
	Resource string         `gorm:"not null;size:255;index"`
	Rights   PermissionBits `gorm:"not null"`
}

// HasPermission проверяет, есть ли конкретное право в маске
func (p PermissionBits) HasPermission(permission PermissionBits) bool {
	return p&permission == permission
}

// AddPermission добавляет право к маске
func (p *PermissionBits) AddPermission(permission PermissionBits) {
	*p |= permission
}

// RemovePermission удаляет право из маски
func (p *PermissionBits) RemovePermission(permission PermissionBits) {
	*p &^= permission
}

// String возвращает текстовое представление прав
func (p PermissionBits) String() string {
	var rights []string
	if p.HasPermission(PermissionRead) {
		rights = append(rights, "read")
	}
	if p.HasPermission(PermissionWrite) {
		rights = append(rights, "write")
	}
	if p.HasPermission(PermissionDelete) {
		rights = append(rights, "delete")
	}
	return strings.Join(rights, "|")
}
