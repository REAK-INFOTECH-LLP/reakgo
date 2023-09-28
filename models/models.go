package models

import (
	"database/sql"
	"errors"
	"fmt"
	"reakgo/utility"
)

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

func GenerateCache() {

}

func VerifyToken() {
	if entry, err := utility.Cache.Get("my-unique-key"); err == nil {
		fmt.Println(string(entry))
	} else {
		// Pull Record from DB and add to Cache
	}
}
