package controllers

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"reakgo/models"
	"reakgo/utility"
)

var Utility utility.Helper

var (
	// ErrCode is a config or an internal error
	ErrCode = errors.New("Case statement in code is not correct.")
	// ErrNoResult is a not results error
	ErrNoResult = errors.New("Result not found.")
	// ErrUnavailable is a database not available error
	ErrUnavailable = errors.New("Database is unavailable.")
	// ErrUnauthorized is a permissions violation
	ErrUnauthorized = errors.New("User does not have permission to perform this operation.")
)

// standardizeErrors returns the same error regardless of the database used
func standardizeError(err error) error {
	if err == sql.ErrNoRows {
		return ErrNoResult
	}

	return err
}

func CheckACL(w http.ResponseWriter, r *http.Request, allowedAccess []string) bool {
	// Check if Token is provided else we continue with Session based auth management
	apiToken := r.Header.Get("Authorization")
	if apiToken == "" {
		// API Token not found, Switching to session based auth
		userType := Utility.SessionGet(r, "type")
		if userType == nil {
			userType = "guest"
		}
		if !Utility.StringInArray(fmt.Sprintf("%v", userType), allowedAccess) {
			Utility.RedirectTo(w, r, "forbidden")
			return false
		}
	} else {
		// Token based Auth
		err := models.VerifyToken(r)
		if err != nil {
			http.Error(w, "403 Forbidden", http.StatusForbidden)
			return false
		}
	}
	return true
}
