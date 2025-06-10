package errors

import (
	"encoding/json"
	"net/http"
	"strings"
)

type Details struct {
	Field string `json:"field,omitempty"`
	Error string `json:"errors"`
}

type HTTPError struct {
	Code    int       `json:"code"`
	Message string    `json:"message"`
	Details []Details `json:"details,omitempty"`
}

func (e *HTTPError) Error() string {
	return e.Message
}

func Error(msg string, statusCode int, details ...string) *HTTPError {
	var fields []Details

	for _, d := range details {
		parts := strings.SplitN(d, ":", 2)
		if len(parts) == 2 {
			fields = append(fields, Details{
				Field: parts[0],
				Error: parts[1],
			})
		}
	}

	return &HTTPError{
		Code:    statusCode,
		Message: msg,
		Details: fields,
	}
}

func WriteHTTPError(rw http.ResponseWriter, err *HTTPError) {
	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(err.Code)

	_ = json.NewEncoder(rw).Encode(struct {
		Error HTTPError `json:"error"`
	}{Error: *err})
}

func UnauthorizedError(details ...string) *HTTPError {
	return Error("Unauthorized", http.StatusUnauthorized, details...)
}

func NotFoundError(msg string, details ...string) *HTTPError {
	return Error(msg+" not found", http.StatusNotFound, details...)
}

func BadGatewayError(details ...string) *HTTPError {
	return Error("bad gateway", http.StatusBadGateway, details...)
}

func InternalError(details ...string) *HTTPError {
	return Error("internal error", http.StatusInternalServerError, details...)
}

func BadRequestError(msg string, details ...string) *HTTPError {
	return Error(msg, http.StatusBadRequest, details...)
}
