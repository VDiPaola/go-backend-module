package module_models

import "github.com/dgrijalva/jwt-go/v4"

// GoogleClaims -
type GoogleClaims struct {
	Email         string `json:"email"`
	EmailVerified bool   `json:"email_verified"`
	FirstName     string `json:"given_name"`
	LastName      string `json:"family_name"`
	jwt.StandardClaims
}

type GoogleResponse struct {
	JWT string `json:"jwt"`
}
