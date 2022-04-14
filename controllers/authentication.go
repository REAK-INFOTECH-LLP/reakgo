package controllers

import (
	"log"
	"net/http"
	"os"
    "bytes"
    "image/png"
//    b64 "encoding/base64"
	
    "reakgo/utility"

	"github.com/gorilla/sessions"
	"github.com/pquerna/otp/totp"
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
                        if(r.FormValue("rememberMe") == "yes"){
                            utility.Store.Options = &sessions.Options{
                                MaxAge: 60*1,
                            }
                        }
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


func TwoFaChallenge(w http.ResponseWriter, r *http.Request) {
    key, _ := totp.Generate(totp.GenerateOpts{
		Issuer: "Example.com",
		AccountName: "alice@example.com",
    })
    log.Println(key)
    log.Println(key.Secret())
    utility.RenderTemplate(w, r, "twoFactorAuthCode", "demo")
}


func RegisterTwoFa(w http.ResponseWriter, r *http.Request){
    log.Println("Register Two FA")
    key, _ := totp.Generate(totp.GenerateOpts{
		Issuer: os.Getenv("TOTP_ISSUER"),
		AccountName: "alice@example.com",
    })
    log.Println(key)
    log.Println(key.Secret())
    img, _ := key.Image(400,400)

    var buf bytes.Buffer
    png.Encode(&buf, img)
    //buf = b64.StdEncoding.EncodeToString([]byte(buf))

    utility.RenderTemplate(w, r, "twoFactorAuthCode", buf)
}
