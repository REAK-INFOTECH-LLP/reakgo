package controllers

import (
	"reakgo/models"
	"reakgo/utility"
)

type Env struct {
	authentication interface {
		GetUserByColumn(column string, value string) (models.Authentication, error)
		ForgotPassword(id int32) (string, error)
		TokenVerify(token string, newPassword string) (bool, error)
		TwoFactorAuthAdd(secret string, userId int) (bool, error)
		CheckTwoFactorRegistration(userId int32) string
	}
	formAddView interface {
		Add(name string, address string)
		View() ([]models.FormAddView, error)
	}
}

var Db *Env

func init() {
	Db = &Env{
		authentication: models.AuthenticationModel{DB: utility.Db},
		formAddView:    models.FormAddViewModel{DB: utility.Db},
	}
}
