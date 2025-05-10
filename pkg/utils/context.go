// utils/context.go

package utils

import (
	"errors"

	"github.com/gin-gonic/gin"
)

func GetUserIDFromContext(ctx *gin.Context) (uint, error) {
	userID, exists := ctx.Get("userID")
	if !exists {
		return 0, errors.New("user id not found in context")
	}

	parsedUserID, ok := userID.(uint)
	if !ok {
		return 0, errors.New("invalid user id type")
	}

	return parsedUserID, nil
}

func GetUserRolesFromContext(ctx *gin.Context) ([]string, error) {
	userRoles, exists := ctx.Get("userRoles")
	if !exists {
		return nil, errors.New("user roles not found in context")
	}

	roles, ok := userRoles.([]string)
	if !ok {
		return nil, errors.New("invalid user roles type")
	}

	return roles, nil
}
