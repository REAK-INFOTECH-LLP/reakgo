package controllers

import (
	"net/http"
	"reakgo/utility"
)

func BaseIndex(w http.ResponseWriter, r *http.Request) {
    utility.RenderTemplate(w, r, "index", nil)
}

func Dashboard(w http.ResponseWriter, r *http.Request) {
    utility.RenderTemplate(w, r, "dashboard", nil)
}
