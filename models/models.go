package models

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"reakgo/utility"
	"strings"
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

func VerifyToken(r *http.Request) error {
	authToken := r.Header.Get("Authorization")
	userToken := strings.Split(authToken, " ")
	if len(userToken) != 2 || strings.ToLower(userToken[0]) != "bearer" {
		return fmt.Errorf("invalid authorization header format")
	}

	if entry, err := utility.Cache.Get(userToken[1]); err == nil {

		// Set JSON Payload to the header so the users can use the same
		r.Header.Add("tokenPayload", string(entry))
		return err
	} else {
		// Pull Record from DB and add to Cache

		// PS : Adding DB failsafe opens up a DDoS security issue that people can keep trying with random tokens
		// and crash the server easily by blocking DB pool connections

		data, err := Authentication.GetAuthenticationByToken(Authentication{}, userToken[1])
		if err != nil {
			return err
		}
		jsonData, err := json.Marshal(data)
		if err == nil {
			// Rehydrate if we got the JSON conversion done
			// Fails would be rare, but if it happens kind of defeat the purpose as JSON unmarshall would also crash
			utility.Cache.Set(data.Token, jsonData)
			r.Header.Add("tokenPayload", string(jsonData))
			return err
		} else {
			return err
		}
	}
}

func jsonStringToAuthentication(jsonStr string) (*Authentication, error) {
	var auth Authentication

	// Unmarshal the JSON string into the Authentication struct
	if err := json.Unmarshal([]byte(jsonStr), &auth); err != nil {
		return nil, err
	}

	return &auth, nil
}
