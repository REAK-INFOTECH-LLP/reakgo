package controllers

import (
	"log"
	"net/http"
	"reakgo/models"
	"reakgo/utility"
)

func BaseIndex(w http.ResponseWriter, r *http.Request) {
	name := []string{"Test1", "Test2"}
	utility.RenderTemplate(w, r, "index", name)
}

func BaseAPI(w http.ResponseWriter, r *http.Request) {
	log.Println(r.Header.Get("tokenPayload"))
	err := models.VerifyToken(r)
	if err != nil {
		// Redirect to 403
	} else {
		name := []string{"Test1", "Test2"}
		utility.RenderTemplate(w, r, "index", name)
	}
}

func Dashboard(w http.ResponseWriter, r *http.Request) {
	utility.RenderTemplate(w, r, "dashboard", nil)
}
