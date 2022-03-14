package controllers

import (
    "net/http"
    "log"
    "reakgo/utility"
    "golang.org/x/crypto/bcrypt"
)

func Login(w http.ResponseWriter, r *http.Request) {
    if (r.Method == "POST"){
        // Check for any form parsing error
        err := r.ParseForm()
        if (err != nil) {
            log.Println("form parsing failed")
            utility.RenderTemplate(w, r, "login", "demo")
        } else {
            // Parsing form went fine, Now we can access all the values
            email := r.FormValue("email")
            password := r.FormValue("password")
            confirmPassword := r.FormValue("confirmPassword")
            //rememberMe := r.FormValue("rememberMe")

            // No need to check for empty values as DB Authentication will take care of it

            // Backend validation for password and confirmPassword
            if(confirmPassword == password){
                row, err := Db.authentication.GetUserByEmail(email)
                if(err != nil){
                    // In case of MYSQL issues or no results are returned    
                    utility.AddFlash("error","Credentials didn't match, Please try again.", w, r)
                    utility.RenderTemplate(w, r, "login", "demo")
                } else {
                    match := bcrypt.CompareHashAndPassword([]byte(row.Password), []byte(r.FormValue("password")))
                    if(match != nil){
                        utility.AddFlash("error","Credentials didn't match, Please try again.", w, r)
                        utility.RenderTemplate(w, r, "login", "demo")
                    } else {
                        // Password match has been a success
                        utility.SessionSet(w, r, utility.Session{Key:"email", Value: row.Email})
                        utility.SessionSet(w, r, utility.Session{Key:"type", Value: "user"})
                        utility.AddFlash("success","Success !, Logged in.", w, r)
                        utility.RedirectTo(w, r, "dashboard")
                        //utility.RenderTemplate(w, r, "login", "demo")
                    }
                }
            }
        }
    } else {
        utility.RenderTemplate(w, r, "login", "demo")
    }
}
