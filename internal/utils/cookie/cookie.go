package cookie

import (
	"fmt"
	"net/http"
)

const (
	AccessTokenName  = "access_token"
	RefreshTokenName = "refresh_token"
)

func GetCookie(req *http.Request, name string) (string, error) {
	cookie, err := req.Cookie(name)
	if err != nil {
		return "", err
	}

	if cookie.Value == "" {
		return "", fmt.Errorf("cookie %q is empty", name)
	}

	return cookie.Value, nil
}
