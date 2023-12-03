package jwt

import (
	"boilerplate/backend/database"
	"boilerplate/backend/models"
	"os"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go/v4"
	"github.com/gofiber/fiber/v2"

	"context"
	"fmt"

	"google.golang.org/api/oauth2/v2"
	"google.golang.org/api/option"
)

func CheckToken(c *fiber.Ctx) (*jwt.Token, error) {
	jwtSecret := os.Getenv("JWT_SECRET")
	//get jwt cookie
	cookie := c.Cookies("jwt")

	//this checks that token is valid idk how it works tho
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

func SetToken(c *fiber.Ctx, id uint) error {
	jwtSecret := os.Getenv("JWT_SECRET")
	//JWT token
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		Issuer:    strconv.FormatUint(uint64(id), 10),
		ExpiresAt: jwt.At(time.Now().Add(time.Hour * 24)),
		IssuedAt:  jwt.Now(),
	})

	token, err := claims.SignedString([]byte(jwtSecret))

	if err != nil {
		return err
	}

	SetCookie(c, token)

	return nil
}

func SetCookie(c *fiber.Ctx, token string) error {
	//set cookie
	cookie := &fiber.Cookie{
		Name:     "jwt",
		Value:    token,
		Expires:  time.Now().Add(time.Hour * 24),
		HTTPOnly: true,
	}
	c.Cookie(cookie)

	return nil
}

//GOOGLE JWT

func SetTokenGoogle(c *fiber.Ctx, token *oauth2.Tokeninfo) error {
	jwtSecret := os.Getenv("JWT_SECRET")
	//JWT token
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		Issuer:    token.UserId,
		ExpiresAt: jwt.At(time.Now().Add(time.Second * time.Duration(token.ExpiresIn))),
		IssuedAt:  jwt.Now(),
	})

	signedToken, err := claims.SignedString([]byte(jwtSecret))

	if err != nil {
		return err
	}

	SetCookie(c, signedToken)

	return nil
}

func VerifyIDToken(idToken string) (*oauth2.Tokeninfo, error) {
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

	clientID := os.Getenv("GOOGLE_CLIENT_ID")
	if tokenInfo.Audience != clientID {
		return nil, fmt.Errorf("mismatching client ID, expected %s but got %s", clientID, tokenInfo.Audience)
	}

	return tokenInfo, nil
}

func GetUserFromToken(c *fiber.Ctx) (models.User, error) {
	//check token
	token, err := CheckToken(c)

	if err != nil {
		return models.User{}, err
	}

	//get claim in correct format
	claims := GetClaims(token)

	var user models.User

	//get user from claim
	database.Connection.Where("id = ?", claims.Issuer).First(&user)

	return user, nil
}
