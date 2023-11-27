package controllers

import (
	"net/http"
)

func BaseIndex(w http.ResponseWriter, r *http.Request) {
	name := []string{"Test1", "Test2"}
	Helper.RenderTemplate(w, r, "index", name)
}

func Dashboard(w http.ResponseWriter, r *http.Request) {
	Helper.RenderTemplate(w, r, "dashboard", nil)
}
