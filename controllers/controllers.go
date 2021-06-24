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
    profile interface {
        Fetch() (models.Profile, error)
    }
}

var Db *Env

func init(){
    Db = &Env{
        authentication: models.AuthenticationModel{DB: utility.Db},
        profile: models.ProfileModel{DB: utility.Db},
    }
}
