package controllers

import (
	"net/http"
	"reakgo/utility"
)

func BaseIndex(w http.ResponseWriter, r *http.Request) {
	name := []string{"Test1", "Test2"}
	utility.RenderTemplate(w, r, "index", name)
}

func Dashboard(w http.ResponseWriter, r *http.Request) {
	utility.RenderTemplate(w, r, "dashboard", nil)
}
