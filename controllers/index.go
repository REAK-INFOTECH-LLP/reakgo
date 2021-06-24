package controllers

import (
	"net/http"
	"reakgo/utility"
)

func BaseIndex(w http.ResponseWriter, r *http.Request) {
    data, err := Db.profile.Fetch()
    if err != nil {
        utility.AddFlash("danger", "Hello World !", w, r)
    }
    utility.AddFlash("danger", "Hello World !", w, r)
    utility.RenderTemplate(w, r, "profile", data)
}

func Login(w http.ResponseWriter, r *http.Request) {
    utility.View.ExecuteTemplate(w, "login", nil)
}

func Edit(w http.ResponseWriter, r *http.Request) {
    utility.View.ExecuteTemplate(w, "editProfile", nil)
}
