package auth

import (
	"errors"
	"fmt"
	"git.foxminded.ua/foxstudent106092/user-management/config"
	"git.foxminded.ua/foxstudent106092/user-management/internal/business/model"
	"github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"slices"
	"time"
)

type UsersWithRoleJwtClaims struct {
	Name string `json:"name"`
	Role string `json:"role"`
	jwt.RegisteredClaims
}

func GenerateJWTToken(u *model.User, secretKey interface{}) (string, error) {
	claims := &UsersWithRoleJwtClaims{
		u.Username,
		u.Role,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(secretKey)
}

func getKeyFunc(cfgAuth *config.Auth) jwt.Keyfunc {
	return func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(cfgAuth.SecretKey), nil
	}
}

func GetParseTokenFunc(cfgAuth *config.Auth, accessibleRoles []string) func(ctx echo.Context, auth string) (interface{}, error) {
	return func(ctx echo.Context, auth string) (interface{}, error) {
		keyFunc := getKeyFunc(cfgAuth)

		token, err := jwt.Parse(auth, keyFunc)
		if err != nil {
			return nil, err
		}

		if !token.Valid {
			return nil, errors.New("provided token is invalid")
		}

		claimsMap, ok := token.Claims.(jwt.MapClaims)

		if !ok {
			return nil, errors.New("failed to map claims")
		}

		claimsUsername := claimsMap["name"].(string)
		claimsRole := claimsMap["role"].(string)

		if !slices.Contains(accessibleRoles, claimsRole) {
			return nil, errors.New("no access to this source")
		}

		ctx.Set("username", claimsUsername)
		ctx.Set("role", claimsRole)

		return token, nil
	}
}

func GetTokenConfig(cfgAuth *config.Auth, accessibleRoles []string) echojwt.Config {
	tokenConfig := echojwt.Config{
		NewClaimsFunc: func(c echo.Context) jwt.Claims {
			return new(UsersWithRoleJwtClaims)
		},
		SigningKey:     []byte(cfgAuth.SecretKey),
		ParseTokenFunc: GetParseTokenFunc(cfgAuth, accessibleRoles),
	}

	return tokenConfig
}
