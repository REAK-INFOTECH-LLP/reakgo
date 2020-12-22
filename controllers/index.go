package controllers

import (
	"fmt"
	"log"
	"net/http"
	"reakgo/utility"
)

func BaseIndex(w http.ResponseWriter, r *http.Request) {
	//session, _ := main.Store.Get(r, "session-name")
	//data,_ := models.ModelDB.authentication.All()
	fmt.Println(Db.authentication.All())
	if r.Method == "GET" {
		data := struct {
			Title string
		}{
			Title: "Demo From Controller",
		}
		utility.View.ExecuteTemplate(w, "index", data)
	} else {
		// In case you want to parse multipart data,
		// Number is the size of form data needed to be saved on RAM rest goes on disk
		//r.ParseMultipartForm(0)
		err := r.ParseForm()
		if err != nil {
			log.Println(err)
		}
		if r.Form.Get("key") == "password" {
			d := utility.Session{Key: "type", Value: "admin"}
			utility.SessionSet(w, r, d)
			utility.RedirectTo(w, r)
		}
	}
}

func Test(w http.ResponseWriter, r *http.Request) {
	d := utility.SessionGet(r, "secret")
	fmt.Println(d)
	utility.View.ExecuteTemplate(w, "test", nil)
}

func Teams(w http.ResponseWriter, r *http.Request) {
	d := utility.SessionGet(r, "secret")
	fmt.Println(d)
	utility.View.ExecuteTemplate(w, "teams", nil)
}

func Manage(w http.ResponseWriter, r *http.Request) {
	d := utility.SessionGet(r, "secret")
	fmt.Println(d)
	utility.View.ExecuteTemplate(w, "manage", nil)
}

func Stats(w http.ResponseWriter, r *http.Request) {
	d := utility.SessionGet(r, "secret")
	fmt.Println(d)
	utility.View.ExecuteTemplate(w, "stats", nil)
}
