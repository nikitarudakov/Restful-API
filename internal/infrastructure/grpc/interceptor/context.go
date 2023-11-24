package interceptor

import (
	"context"
	"errors"
	"google.golang.org/grpc/metadata"
)

func parseContext(ctx context.Context, keys ...string) ([]string, error) {
	var headerValues []string

	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, errors.New("error parsing incoming context")
	}

	for _, key := range keys {
		header := md.Get(key)

		if header != nil {
			headerValues = append(headerValues, header[0])
		}
	}

	if len(keys) != len(headerValues) {
		return nil, errors.New("some headers are missing")
	}

	return headerValues, nil
}
