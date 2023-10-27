package auth

import (
	"encoding/base64"
	"errors"
	"strings"
)

const authPrefix = "Basic"

func ReadAuthHeader(h map[string][]string) (string, error) {
	var authCred string

	if auth, ok := h["Authorization"]; ok {
		if len(auth) > 0 {
			authCred = auth[0]
		} else {
			return "", errors.New("authorization header is missing")
		}

	} else {
		return authCred, errors.New("authorization header is missing")
	}

	return authCred, nil
}

func DecodeBasicAuthCred(a string) (*[]string, error) {
	c := make([]string, 2)

	if a != "" && strings.HasPrefix(a, authPrefix) {
		encodedA := strings.Split(a, " ")[1]
		decodedA, err := base64.StdEncoding.DecodeString(encodedA)
		if err != nil {
			return nil, errors.New("authorization has been unsuccessful")
		}

		aSlice := strings.Split(string(decodedA), ":")
		if len(aSlice) > 1 {
			c[0], c[1] = aSlice[0], aSlice[1]

			return &c, nil
		}
	}

	return nil, errors.New("authorization has been unsuccessful")

}
