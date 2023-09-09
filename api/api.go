package api

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
)

type Error struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func (e Error) Error() string {
	return fmt.Sprintf("code: %s, message: %s", e.Code, e.Message)
}

func sendErr(ctx context.Context, w http.ResponseWriter, statusCode int, err error) {
	var httpErr Error
	if !errors.As(err, &httpErr) {
		httpErr = Error{
			Code:    "unknown_error",
			Message: "An unexpected error happened",
		}
	}
	if statusCode >= 500 {
		slog.Error("unable to process request", "error", err.Error(), "status_code", statusCode)
	}

	sendJSON(ctx, w, statusCode, httpErr)
}

func sendJSON(ctx context.Context, w http.ResponseWriter, statusCode int, body interface{}) {
	const jsonContentType = "application/json; charset=utf-8"

	w.Header().Set("Content-Type", jsonContentType)
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(body); err != nil {
		slog.Error("unable to encode response", "error", err.Error())
	}
}
