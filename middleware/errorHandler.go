package middleware

import (
	"encoding/json"
	"log"
	"net/http"
	"runtime/debug"

	"github.com/IsraeliPS/GoMongo/config"
)

// ErrorResponse defines the structure for error responses
type ErrorResponse struct {
    Code    int    `json:"code"`
    Message string `json:"message"`
}

// ErrorHandler is a middleware function that handles errors
func ErrorHandler(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        defer func() {
            if err := recover(); err != nil {
                config.Logger.Error("Unhandled error: ", err, "\n", string(debug.Stack()))
                
                log.Printf("Recovered from panic: %v", err)
                var code int
                var message string

                switch t := err.(type) {
                case string:
                    message = t
                    code = http.StatusInternalServerError
                case error:
                    message = t.Error()
                    code = http.StatusInternalServerError
                default:
                    message = "Unknown error"
                    code = http.StatusInternalServerError
                }

                w.Header().Set("Content-Type", "application/json")
                w.WriteHeader(code)
                json.NewEncoder(w).Encode(ErrorResponse{Code: code, Message: message})
            }
        }()
        next.ServeHTTP(w, r)
    })
}