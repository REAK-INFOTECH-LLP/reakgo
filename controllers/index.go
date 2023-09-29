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
	data, err := models.VerifyToken(r, false)
	if err != nil {
		// Redirect to 403
	} else {
		log.Println(data)
		name := []string{"Test1", "Test2"}
		utility.RenderTemplate(w, r, "index", name)
	}
}

func Dashboard(w http.ResponseWriter, r *http.Request) {
	utility.RenderTemplate(w, r, "dashboard", nil)
}
