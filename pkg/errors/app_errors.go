package apperrors

// Ошибки сервера
const (
	ErrInternalServerError = "internal server error"
)

// Ошибки, связанные с комментариями
const (
	ErrCommentNotFound = "comment not found"
)

// Ошибки, связанные со статьями
const (
	ErrArticleNotFound  = "article not found"
	ErrInvalidArticleID = "invalid article ID"
)

// Ошибки, связанные с пользователями
const (
	ErrAccessDenied         = "access denied: you are not the owner or don't have required role"
	ErrUserNotAuthenticated = "user not authenticated"
	ErrUserNotFound         = "user not found"
	ErrInvalidUserID        = "invalid user ID"
)

// Ошибки, связанные с аутентификацией
const (
	ErrInvalidCredentials     = "invalid credentials"
	ErrFailedToCreateUser     = "failed to create user"
	ErrFailedToGenerateTokens = "failed to generate tokens"
	ErrUserAlreadyExists      = "user with this email already exists"
)

// Ошибки, связанные с ролями
const (
	ErrRoleNotFound        = "role not found"
	ErrInvalidRoleID       = "invalid role ID"
	ErrUserRolesNotFound   = "user roles not found"
	ErrFailedToAssignRole  = "failed to assign role"
	ErrRoleAlreadyAssigned = "role is already assigned to the user"
)

// Ошибки, связанные с файлами
const (
	ErrInvalidFile         = "failed to retrieve file"
	ErrFileSizeExceeded    = "file size exceeds allowed limit"
	ErrUnsupportedFileType = "unsupported file type"
	ErrFailedToSaveFile    = "failed to save file"
	ErrMediaNotFound       = "media not found"
	ErrInvalidMediaID      = "invalid media ID"
)
