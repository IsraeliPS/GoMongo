package middleware

import (
	"encoding/json"
	"net/http"

	"github.com/IsraeliPS/GoMongo/config"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
)

var validate *validator.Validate

func init() {
    validate = validator.New()
}

// ValidationMiddleware validates the request body against a given struct
func ValidationMiddleware(schema interface{}) mux.MiddlewareFunc {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            if err := json.NewDecoder(r.Body).Decode(&schema); err != nil {
                config.Logger.Error("Invalid request body: ", err)
                http.Error(w, "Invalid request body", http.StatusBadRequest)
                return
            }

            if err := validate.Struct(schema); err != nil {
                if _, ok := err.(*validator.InvalidValidationError); ok {
                    config.Logger.Error("Validation error: ", err)
                    http.Error(w, err.Error(), http.StatusInternalServerError)
                    return
                }

                validationErrors := err.(validator.ValidationErrors)
                errorMessages := make(map[string]string)
                for _, err := range validationErrors {
                    errorMessages[err.Field()] = err.Tag()
                }

                config.Logger.Error("Validation errors: ", validationErrors)
                w.Header().Set("Content-Type", "application/json")
                w.WriteHeader(http.StatusBadRequest)
                json.NewEncoder(w).Encode(map[string]interface{}{
                    "errors": errorMessages,
                })
                return
            }

            next.ServeHTTP(w, r)
        })
    }
}
