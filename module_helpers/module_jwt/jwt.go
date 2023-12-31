package module_jwt

import (
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go/v4"
	"github.com/gofiber/fiber/v2"

	"context"
	"fmt"

	"google.golang.org/api/oauth2/v2"
	"google.golang.org/api/option"
)

func CheckToken(c *fiber.Ctx, jwtSecret string) (*jwt.Token, error) {
	//get jwt cookie
	cookie := c.Cookies("jwt")

	//checks token is valid
	token, err := jwt.ParseWithClaims(cookie,
		&jwt.StandardClaims{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(jwtSecret), nil
		})

	if err != nil {
		return nil, err
	}

	return token, nil
}

func GetClaims(t *jwt.Token) *jwt.StandardClaims {
	return t.Claims.(*jwt.StandardClaims)
}

func SetToken(c *fiber.Ctx, id uint, jwtSecret string, expiresIn time.Duration) error {
	//JWT token
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		Issuer:    strconv.FormatUint(uint64(id), 10),
		ExpiresAt: jwt.At(time.Now().Add(expiresIn)),
		IssuedAt:  jwt.Now(),
	})

	token, err := claims.SignedString([]byte(jwtSecret))

	if err != nil {
		return err
	}

	SetCookie(c, token, expiresIn)

	return nil
}

func SetCookie(c *fiber.Ctx, token string, expiresIn time.Duration) error {
	//set cookie
	cookie := &fiber.Cookie{
		Name:     "jwt",
		Value:    token,
		Expires:  time.Now().Add(expiresIn),
		HTTPOnly: true,
	}
	c.Cookie(cookie)

	return nil
}

//GOOGLE JWT

func SetTokenGoogle(c *fiber.Ctx, token *oauth2.Tokeninfo, jwtSecret string) error {
	//JWT token
	expiresDuration := time.Duration(token.ExpiresIn)
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		Issuer:    token.UserId,
		ExpiresAt: jwt.At(time.Now().Add(expiresDuration)),
		IssuedAt:  jwt.Now(),
	})

	signedToken, err := claims.SignedString([]byte(jwtSecret))

	if err != nil {
		return err
	}

	SetCookie(c, signedToken, expiresDuration)

	return nil
}

func VerifyIDToken(idToken string, googleClientID string) (*oauth2.Tokeninfo, error) {
	ctx := context.Background()
	oauth2Service, err := oauth2.NewService(ctx, option.WithoutAuthentication())
	if err != nil {
		return nil, err
	}

	tokenInfoCall := oauth2Service.Tokeninfo()
	tokenInfoCall.IdToken(idToken)
	tokenInfo, err := tokenInfoCall.Do()
	if err != nil {
		return nil, err
	}

	if tokenInfo.Audience != googleClientID {
		return nil, fmt.Errorf("mismatching client ID, expected %s but got %s", googleClientID, tokenInfo.Audience)
	}

	return tokenInfo, nil
}
