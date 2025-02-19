package tokens

import (
	"github.com/blinkinglight/pocketbase-mysql/core"
	"github.com/blinkinglight/pocketbase-mysql/models"
	"github.com/blinkinglight/pocketbase-mysql/tools/security"
	"github.com/golang-jwt/jwt/v4"
)

// NewAdminAuthToken generates and returns a new admin authentication token.
func NewAdminAuthToken(app core.App, admin *models.Admin) (string, error) {
	return security.NewToken(
		jwt.MapClaims{"id": admin.Id, "type": "admin"},
		(admin.TokenKey + app.Settings().AdminAuthToken.Secret),
		app.Settings().AdminAuthToken.Duration,
	)
}

// NewAdminResetPasswordToken generates and returns a new admin password reset request token.
func NewAdminResetPasswordToken(app core.App, admin *models.Admin) (string, error) {
	return security.NewToken(
		jwt.MapClaims{"id": admin.Id, "type": "admin", "email": admin.Email},
		(admin.TokenKey + app.Settings().AdminPasswordResetToken.Secret),
		app.Settings().AdminPasswordResetToken.Duration,
	)
}
