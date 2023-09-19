package controllers

import (
    "net/http"
    "log"

    "reakgo/utility"
)

func AddForm(w http.ResponseWriter, r *http.Request){
    if (r.Method == "POST") {
        err := r.ParseForm()
        if (err != nil){
            log.Println(err)
        } else {
            // Logic
            name := r.FormValue("name")
            address := r.FormValue("address")

            log.Println(name, address)
            Db.formAddView.Add(name, address)
        }
    }
    utility.RenderTemplate(w, r, "addForm", nil) 
}


func ViewForm(w http.ResponseWriter, r *http.Request){
    result, err := Db.formAddView.View()
    if (err != nil){
        log.Println(err)
    }
    utility.RenderTemplate(w, r, "viewForm", result)
}
