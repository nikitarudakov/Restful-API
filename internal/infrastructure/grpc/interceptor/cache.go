package interceptor

import (
	"context"
	"git.foxminded.ua/foxstudent106092/user-management/config"
	"git.foxminded.ua/foxstudent106092/user-management/internal/infrastructure/datastore/cache"
	"git.foxminded.ua/foxstudent106092/user-management/internal/presenter/grpc/user"
	"git.foxminded.ua/foxstudent106092/user-management/internal/presenter/grpc/vote"
	"git.foxminded.ua/foxstudent106092/user-management/tools/timeparse"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"net/http"
)

func GetServerCacheInterceptor(cfgCache *config.Cache, db *cache.Database) grpc.UnaryServerInterceptor {
	return func(ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler) (interface{}, error) {

		if isAuthLoginOrRegisterMethod(info.FullMethod) {
			h, err := handler(ctx, req)
			return h, err
		}

		expiration, err := timeparse.ParseExpirationTime(cfgCache.ExpirationQuan, cfgCache.ExpirationUnit)
		if err != nil {
			log.Error().Err(err).
				Msg("invalid expiration time settings, cache timeout set to default 1 min")
		}

		values, err := parseContext(ctx, "cache", "target")
		if err != nil {
			log.Warn().Err(err).Msg("error parsing cache from headers")

			h, err := handler(ctx, req)
			return h, err
		}

		cacheMethod, target := values[0], values[1]

		if cacheMethod != http.MethodGet {
			h, err := handler(ctx, req)
			return h, err
		}

		cacheKey := info.FullMethod + "/" + target
		dest := setDestination(info.FullMethod)

		// try to get cache
		if err = db.GetCache(cacheKey, &dest); err == nil {
			return dest, nil
		} else {
			log.Warn().Err(err).Msg("cache was not retrieved")
		}

		// if no cache saved then run handler
		h, err := handler(ctx, req)

		// save result to cache
		if err = db.SetCache(cacheKey, h, expiration); err != nil {
			log.Warn().Str("service", "caching").Err(err).Send()
		}

		return h, err
	}
}

func isAuthLoginOrRegisterMethod(method string) bool {
	if method == "/auth.AuthController/Login" || method == "/auth.AuthController/Register" {
		return true
	}

	return false
}

func setDestination(method string) interface{} {
	var dest interface{}
	switch method {
	case "/vote.VoteController/GetRating":
		dest = &vote.Rating{}
	case "/user.UserController/ListProfiles":
		dest = &user.Profiles{}
	}

	return dest
}
