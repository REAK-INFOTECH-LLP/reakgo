package controllers

import (
	"reakgo/models"
	"reakgo/utility"
)

type Env struct {
    authentication interface {
        GetUserByEmail(email string) (models.Authentication, error)
        ForgotPassword(id int32) (string, error)
        TokenVerify(token string, newPassword string) (bool, error)
    }
    formAddView interface {
        Add(name string, address string)
        View () ([]models.FormAddView, error)
    }
}

var Db *Env

func init(){
    Db = &Env{
        authentication: models.AuthenticationModel{DB: utility.Db},
        formAddView: models.FormAddViewModel{DB: utility.Db},
    }
}
