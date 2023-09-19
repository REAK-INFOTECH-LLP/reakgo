package controllers

import (
	"bytes"
	b64 "encoding/base64"
	"fmt"
	"image/png"
	"log"
	"net/http"
	"os"
	"strconv"

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
                        utility.SessionSet(w, r, utility.Session{Key:"id", Value: row.Id})
                        utility.SessionSet(w, r, utility.Session{Key:"email", Value: row.Email})
                        utility.SessionSet(w, r, utility.Session{Key:"type", Value: "user"})
                        utility.AddFlash("success","Success !, Logged in.", w, r)
                        if(r.FormValue("rememberMe") == "yes"){
                            utility.Store.Options = &sessions.Options{
                                MaxAge: 60*1,
                            }
                        }
                        token := Db.authentication.CheckTwoFactorRegistration(row.Id)
                        utility.SessionSet(w, r, utility.Session{Key:"2faSecret", Value: token})

                        if(token != ""){
                            utility.RedirectTo(w, r, "verify-2fa")
                        } else {
                            utility.RedirectTo(w, r, "dashboard")
                        }
                        //utility.RenderTemplate(w, r, "login", "demo")
                    }
                }
            }
        }
    } else {
        utility.RenderTemplate(w, r, "login", "demo")
    }
}


func VerifyTwoFa(w http.ResponseWriter, r *http.Request) {
    if(r.Method == "GET") {
        utility.RenderTemplate(w, r, "verifyTwoFa", "demo")
    } else if (r.Method == "POST") {
        err := r.ParseForm()
        if(err!=nil){
            log.Println(err)
        } else {
            twoFaVerify := r.FormValue("twoFaVerify")
            secret := fmt.Sprintf("%v", utility.SessionGet(r, "2faSecret")) 
            if(totp.Validate(twoFaVerify, secret)){
                utility.RedirectTo(w, r, "dashboard")
            } else {
                utility.RedirectTo(w, r, "login")
            }
        }
    }
}


func RegisterTwoFa(w http.ResponseWriter, r *http.Request){
    
    if(r.Method == "GET"){
        email := utility.SessionGet(r, "email")
        key, _ := totp.Generate(totp.GenerateOpts{
            Issuer: os.Getenv("TOTP_ISSUER"),
            AccountName: fmt.Sprintf("%v", email),
        })
        utility.SessionSet(w, r, utility.Session{Key:"totpSecret", Value: key.Secret()})
        
        img, _ := key.Image(400,400)

        var buf bytes.Buffer
        png.Encode(&buf, img)

        data := b64.StdEncoding.EncodeToString([]byte(buf.String()))
        data = "data:image/png;base64,"+data

        utility.RenderTemplate(w, r, "twoFactorRegister", data)
    } else if (r.Method == "POST") {
        err := r.ParseForm()
        if (err != nil){

        } else {
            verifyToken := r.FormValue("challengeCode")
            validationResult := totp.Validate(verifyToken, fmt.Sprintf("%v", utility.SessionGet(r, "totpSecret")))
            if(validationResult){
                secret := fmt.Sprintf("%v", utility.SessionGet(r, "totpSecret"))
                userId := fmt.Sprintf("%v", utility.SessionGet(r, "id")) 
                intUserId, _ := strconv.Atoi(userId)
                Db.authentication.TwoFactorAuthAdd(secret, intUserId)
                utility.RenderTemplate(w, r, "successTwoFactor", nil)
            } else {
                // Show Error Page
                utility.RenderTemplate(w, r, "failureTwoFactor", nil)
            }
        }
    }
}
