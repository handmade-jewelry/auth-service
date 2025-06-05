package errors

import (
	"encoding/json"
	"net/http"
	"strings"
)

type FieldError struct {
	Field string `json:"field,omitempty"`
	Error string `json:"error"`
}

type ErrorResponse struct {
	Code    int          `json:"code"`
	Message string       `json:"message"`
	Details []FieldError `json:"details,omitempty"`
}

func Error(msg string, statusCode int, details ...string) *ErrorResponse {
	var fields []FieldError

	for _, d := range details {
		parts := strings.SplitN(d, ":", 2)
		if len(parts) == 2 {
			fields = append(fields, FieldError{
				Field: parts[0],
				Error: parts[1],
			})
		}
	}

	return &ErrorResponse{
		Code:    statusCode,
		Message: msg,
		Details: fields,
	}
}

func UnauthorizedError(details ...string) *ErrorResponse {
	return Error("Unauthorized", http.StatusUnauthorized, details...)
}

func NotFoundError(msg string, details ...string) *ErrorResponse {
	return Error(msg, http.StatusNotFound, details...)
}

func WriteHTTPError(rw http.ResponseWriter, err *ErrorResponse) {
	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(err.Code)

	_ = json.NewEncoder(rw).Encode(struct {
		Error ErrorResponse `json:"error"`
	}{Error: *err})
}
