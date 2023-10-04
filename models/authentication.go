package models

import (
	"database/sql"
	"fmt"
	"log"
	"reakgo/utility"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type Authentication struct {
	Id             int32
	Email          string
	Password       string
	Token          string
	TokenTimestamp int64 `db:"tokenTimestamp"`
}

type TwoFactor struct {
	UserId int32  `db:"userId"`
	Secret string `db:"secret"`
}

func (auth Authentication) GetUserByEmail(email string) (Authentication, error) {
	var selectedRow Authentication

	err := utility.Db.Get(&selectedRow, "SELECT * FROM authentication WHERE email = ?", email)

	return selectedRow, err
}

func (auth Authentication) ForgotPassword(id int32) (string, error) {
	Token, err := utility.GenerateRandomString(60)
	if err != nil {
		log.Println("Random String Generator Failed")
	}
	TokenTimestamp := time.Now().Unix()
	query, err := utility.Db.Prepare("UPDATE authentication SET Token = ?, TokenTimestamp = ? WHERE id = ?")
	if err != nil {
		log.Println("MySQL Query Failed")
	}
	_, err = query.Exec(Token, TokenTimestamp, id)
	if err != nil {
		log.Println(err)
	} else {
		// pass
	}
	return Token, err
}

func (auth Authentication) TokenVerify(token string, newPassword string) (bool, error) {
	var selectedRow Authentication

	rows := utility.Db.QueryRow("SELECT * FROM authentication WHERE token = ?", token)
	err := rows.Scan(&selectedRow.Id, &selectedRow.Email, &selectedRow.Password, &selectedRow.Token, &selectedRow.TokenTimestamp)
	if err != nil {
		log.Println(err)
		return true, err
	}
	if (selectedRow.TokenTimestamp + 360000) > time.Now().Unix() {
		_, err := auth.ChangePassword(newPassword, selectedRow.Id)
		if err != nil {
			return true, err
		} else {
			return false, err
		}
	}
	return false, err
}

func (auth Authentication) ChangePassword(newPassword string, id int32) (bool, error) {
	query, err := utility.Db.Prepare("UPDATE authentication SET password = ? WHERE id = ?")
	if err != nil {
		log.Println("MySQL Query Failed")
	}
	newPasswordHash, err := bcrypt.GenerateFromPassword([]byte(newPassword), 10)
	if err != nil {
		log.Println(err)
		return true, err
	}
	_, err = query.Exec(newPasswordHash, id)
	if err != nil {
		log.Println(err)
		return true, err
	} else {
		return false, err
	}
}

func (auth Authentication) TwoFactorAuthAdd(secret string, userId int) (bool, error) {
	_, err := utility.Db.NamedExec("INSERT INTO twoFactor (userId, secret) VALUES(:id, :2faSecret) ON DUPLICATE KEY UPDATE secret=:2faSecret", map[string]interface{}{"2faSecret": secret, "id": userId})
	if err != nil {
		return false, err
	} else {
		return true, err
	}
}

func (auth Authentication) CheckTwoFactorRegistration(userId int32) string {
	twoFactor := TwoFactor{}
	utility.Db.Get(&twoFactor, "SELECT * FROM twoFactor WHERE userId = ?", userId)
	return twoFactor.Secret
}

func (auth Authentication) GetAllAuthRecords() ([]Authentication, error) {
	var allAuthenticationRows []Authentication

	// SQL query to select all rows from the Abc table
	query := "SELECT * FROM authentication"

	// Execute the query and scan the results into the Abc slice
	err := utility.Db.Select(&allAuthenticationRows, query)
	if err != nil {
		return nil, err
	}

	return allAuthenticationRows, nil
}

func (auth Authentication) GetAuthenticationByToken(token string) (*Authentication, error) {

	var authO Authentication

	// SQL query to select a row by Token
	query := "SELECT * FROM authentication WHERE Token = ? LIMIT 1"

	// Execute the query and scan the result into the Authentication struct
	err := utility.Db.Get(&authO, query, token)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("Authentication not found for Token: %s", token)
		}
		return nil, err
	}
	return &authO, nil
}
