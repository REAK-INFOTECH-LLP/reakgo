package controllers

import (
    "reakgo/models"
    "reakgo/utility"
)

type Env struct {
    authentication interface {
        All() ([]models.Authentication, error)
    }
}

var Db *Env

func init(){
    Db = &Env{
        authentication: models.AuthenticationModel{DB: utility.Db},
    }
}
