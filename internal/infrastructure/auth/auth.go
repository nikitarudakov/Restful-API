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

func DecodeBasicAuthCred(headerVal string) (*[]string, error) {
	c := make([]string, 2)

	if headerVal != "" && strings.HasPrefix(headerVal, authPrefix) {
		encodedHeaderVal := strings.Split(headerVal, " ")[1]
		decodedHeaderVal, err := base64.StdEncoding.DecodeString(encodedHeaderVal)
		if err != nil {
			return nil, errors.New("authorization has been unsuccessful")
		}

		credSlice := strings.Split(string(decodedHeaderVal), ":")
		if len(credSlice) > 1 {
			c[0], c[1] = credSlice[0], credSlice[1]

			return &c, nil
		}
	}

	return nil, errors.New("authorization has been unsuccessful")

}
