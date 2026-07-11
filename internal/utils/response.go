package utils

import (
	"encoding/json"
	"log"
	"net/http"

	"be-ecommerce/internal/customerror"
)

// SuccessResponse structure for API success
type SuccessResponse struct {
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// ErrorResponse structure for API errors
type ErrorResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Error   string `json:"error,omitempty"`
}

// RespondSuccess sends a JSON success response
func RespondSuccess(w http.ResponseWriter, statusCode int, message string, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(SuccessResponse{
		Status:  "success",
		Message: message,
		Data:    data,
	})
}

// RespondError sends a JSON error response
func RespondError(w http.ResponseWriter, statusCode int, message string, err string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(ErrorResponse{
		Status:  "error",
		Message: message,
		Error:   err,
	})
}

// HandleError automatically translates our AppError to HTTP response
func HandleError(w http.ResponseWriter, err error) {
	// Check if the error is our custom AppError
	if appErr, ok := err.(*customerror.AppError); ok {
		switch appErr.Type {
		case customerror.TypeNotFound:
			RespondError(w, http.StatusNotFound, appErr.Message, "not_found")
		case customerror.TypeConflict:
			RespondError(w, http.StatusConflict, appErr.Message, "conflict")
		case customerror.TypeBadRequest:
			RespondError(w, http.StatusBadRequest, appErr.Message, "bad_request")
		case customerror.TypeUnauthorized:
			RespondError(w, http.StatusUnauthorized, appErr.Message, "unauthorized")
		default:
			// Internal Server Error
			log.Printf("Internal Server Error: %v", appErr.Err)
			RespondError(w, http.StatusInternalServerError, "Internal Server Error", "internal_server_error")
		}
		return
	}

	// Fallback for standard/unknown errors
	log.Printf("Unknown Error: %v", err)
	RespondError(w, http.StatusInternalServerError, "Internal Server Error", "internal_server_error")
}
