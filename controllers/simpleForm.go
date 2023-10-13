package controllers

import (
	"log"
	"net/http"
	"reakgo/models"
	"reakgo/utility"
)

func AddForm(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		err := r.ParseForm()
		if err != nil {
			log.Println(err)
		} else {
			// Logic
			name := r.FormValue("name")
			address := r.FormValue("address")

			log.Println(name, address)
			models.FormAddView{}.Add(name, address)
		}
	}
	utility.RenderTemplate(w, r, "addForm", nil)
}

func ViewForm(w http.ResponseWriter, r *http.Request) {
	result, err := models.FormAddView{}.View()
	if err != nil {
		log.Println(err)
	}
	utility.RenderTemplate(w, r, "viewForm", result)
}
func Test(w http.ResponseWriter, r *http.Request) {
	log.Println("in controllers")
}
