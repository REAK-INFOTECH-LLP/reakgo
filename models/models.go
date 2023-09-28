package models

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log"
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
	allRows, err := Authentication{}.GetAllAuthRecords()
	if err != nil {
		log.Println(err)
	} else {
		for _, val := range allRows {
			jsonData, err := json.Marshal(val)
			if err != nil {
				log.Println("Error encoding JSON:", err)
				break
			}
			utility.Cache.Set(val.Token, jsonData)
		}
	}
}

func VerifyToken() {
	if entry, err := utility.Cache.Get("123"); err == nil {
		log.Println(string(entry))
	} else {
		// Pull Record from DB and add to Cache
		data, err := Authentication.GetAuthenticationByToken(Authentication{}, "123")
		log.Println(data, err)
	}
}
