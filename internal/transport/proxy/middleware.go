package proxy

import (
	"bytes"
	"io"
	"strconv"

	"encoding/json"
	"net/http"
	"net/http/httputil"

	"github.com/handmade-jewelry/auth-service/internal/jwt"
	routeService "github.com/handmade-jewelry/auth-service/internal/service/route"
	"github.com/handmade-jewelry/auth-service/internal/utils/errors"
	"github.com/handmade-jewelry/auth-service/internal/utils/logger"
)

type AuthMiddleware struct {
	routeService *routeService.Service
	jwtService   *jwt.Service
}

func NewAuthMiddleware(
	routeService *routeService.Service,
	jwtService *jwt.Service,
) *AuthMiddleware {
	return &AuthMiddleware{
		routeService: routeService,
		jwtService:   jwtService,
	}
}

func (a *AuthMiddleware) CheckAccess(_ http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		ctx := req.Context()
		route, err := a.checkAuth(ctx, req)
		if err != nil {
			errors.WriteHTTPError(rw, err)
			return
		}

		proxy := &httputil.ReverseProxy{
			Director: func(req *http.Request) {
				req.URL.Scheme = route.Scheme
				req.URL.Host = route.Host
				req.Host = route.Host
				req.URL.Path = route.ServicePath
			},
			ModifyResponse: modifyResponse,
			ErrorHandler: func(rw http.ResponseWriter, req *http.Request, err error) {
				logger.ErrorWithFields("proxy errors", err, "path", req.URL.Path)
				errors.WriteHTTPError(rw, errors.BadGatewayError())
			},
		}

		proxy.ServeHTTP(rw, req)
	})
}

func modifyResponse(resp *http.Response) error {
	if resp.StatusCode >= 400 {
		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		resp.Body.Close()

		logger.ErrorWithFields("proxy errors", nil, "body", string(bodyBytes))

		errStr := errors.Error(http.StatusText(resp.StatusCode), resp.StatusCode)
		newBody, err := json.Marshal(errStr)
		if err != nil {
			return err
		}

		resp.Body = io.NopCloser(bytes.NewReader(newBody))
		resp.ContentLength = int64(len(newBody))
		resp.Header.Set("Content-Length", strconv.Itoa(len(newBody)))
		resp.Header.Set("Content-Type", "application/json; charset=utf-8")
	}
	return nil
}
