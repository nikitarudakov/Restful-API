package interceptor

import (
	"context"
	"errors"
	"git.foxminded.ua/foxstudent106092/user-management/config"
	"git.foxminded.ua/foxstudent106092/user-management/internal/infrastructure/auth"
	"github.com/golang-jwt/jwt/v5"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"slices"
)

func GetServerAuthInterceptor(cfgAuth *config.Auth) grpc.UnaryServerInterceptor {
	defaultRoles := []string{"user", "vote", "moderator"}

	return func(ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler) (interface{}, error) {

		var accessibleRoles []string

		// Skip authorize when GetJWT is requested
		switch info.FullMethod {
		case "/user.UserController/DeleteUserAndProfile":
			accessibleRoles = []string{"vote", "moderator"}
		case "/user.UserController/ListProfiles":
			accessibleRoles = []string{"vote", "moderator"}
		default:
			accessibleRoles = defaultRoles
		}

		if info.FullMethod != "/auth.AuthController/Login" && info.FullMethod != "/auth.AuthController/Register" {
			if err := authorize(ctx, cfgAuth, accessibleRoles); err != nil {
				return nil, err
			}
		}

		// Calls the handler
		h, err := handler(ctx, req)

		return h, err
	}
}

func authorize(ctx context.Context, cfgAuth *config.Auth, accessibleRoles []string) error {
	values, err := parseContext(ctx, "authorization", "target")
	if err != nil {
		return status.Errorf(codes.Unauthenticated, err.Error())
	}

	token, target := values[0], values[1]

	// validateToken function validates the token
	if err = validateToken(token, target, cfgAuth, accessibleRoles); err != nil {
		return status.Errorf(codes.Unauthenticated, err.Error())
	}

	return nil
}

func validateToken(token string, requestTarget string,
	cfgAuth *config.Auth, accessibleRoles []string) error {
	keyFunc := auth.GetKeyFunc(cfgAuth)

	parsedToken, err := jwt.Parse(token, keyFunc)
	if err != nil {
		return errors.New("error with parsing jwt token")
	}

	if !parsedToken.Valid {
		return errors.New("provided token is invalid")
	}

	claimsMap, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok {
		return errors.New("failed to map claims")
	}

	claimsUsername := claimsMap["name"].(string)
	claimsRole := claimsMap["role"].(string)

	if claimsRole != "vote" && requestTarget != "" && requestTarget != claimsUsername {
		return errors.New("no access to this source")
	}

	if !slices.Contains(accessibleRoles, claimsRole) {
		return errors.New("no access to this source")
	}

	return nil
}
