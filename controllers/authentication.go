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
	responce := utility.AjaxResponce{Status: "failure", Message: "Incorrect credentials, Please try again."}
	if r.Method == "POST" {
		// Check for any form parsing error
		err := r.ParseForm()
		if err != nil {
			log.Println("form parsing failed")
		} else {
			// Parsing form went fine, Now we can access all the values
			email := r.FormValue("email")
			// Backend validation for password and confirmPassword
			row, err := Db.authentication.GetUserByColumn("email", email)
			log.Println(email, row)
			message := "Credentials didn't match, Please try again."
			if err != nil {
				// In case of MYSQL issues or no results are returned
				responce.Message = message
				utility.AddFlash("error", message, w, r)

			} else {
				match := bcrypt.CompareHashAndPassword([]byte(row.Password), []byte(r.FormValue("password")))
				if match != nil {
					utility.AddFlash("error", message, w, r)
				} else {
					// Password match has been a success
					utility.SessionSet(w, r, utility.Session{Key: "id", Value: row.Id})
					utility.SessionSet(w, r, utility.Session{Key: "email", Value: row.Email})
					utility.SessionSet(w, r, utility.Session{Key: "type", Value: "user"})
					utility.AddFlash("success", "Success !, Logged in.", w, r)
					if r.FormValue("rememberMe") == "yes" {
						utility.Store.Options = &sessions.Options{
							MaxAge: 60 * 1,
						}
					}
					token := Db.authentication.CheckTwoFactorRegistration(row.Id)
					secure2FAkey := utility.Session{Key: "2faSecret", Value: token}
					utility.SessionSet(w, r, secure2FAkey)
					if utility.IsCurlApiRequest(r) {
						if token != "" {
							responce.Status = "success"
							responce.Payload = secure2FAkey
							responce.Message = "Login Success!, Verify Two facter auth"
						}
					} else {
						if token != "" {
							utility.RedirectTo(w, r, "verify-2fa")
						} else {
							utility.RedirectTo(w, r, "dashboard")
						}
					}
				}
				log.Println(responce)
			}
		}
		utility.RenderTemplate(w, r, "login", responce)
	}
	// only once call this function because repeate ui multiple time
}

func VerifyTwoFa(w http.ResponseWriter, r *http.Request) {
	responce := utility.AjaxResponce{Status: "failure", Message: "Token not match verify, Please try again."}
	if r.Method == "POST" {
		err := r.ParseForm()
		if err != nil {
			log.Println(err)
		} else {
			twoFaVerify := r.FormValue("twoFaVerify")
			secret := fmt.Sprintf("%v", utility.SessionGet(r, "2faSecret"))
			if totp.Validate(twoFaVerify, secret) {
				utility.SessionSet(w, r, utility.Session{Key: "islogedin", Value: true})
				if utility.IsCurlApiRequest(r) {
					responce.Status = "success"
					responce.Payload = ""
					responce.Message = "Login Success!, Verify Two facter auth"
				} else {
					utility.RedirectTo(w, r, "dashboard")
				}
			} else {
				if utility.IsCurlApiRequest(r) {

				} else {
					utility.RedirectTo(w, r, "login")
				}
			}
		}
		utility.RenderTemplate(w, r, "verifyTwoFa", "demo")
	}
}

func RegisterTwoFa(w http.ResponseWriter, r *http.Request) {

	if r.Method == "GET" {
		email := utility.SessionGet(r, "email")
		key, _ := totp.Generate(totp.GenerateOpts{
			Issuer:      os.Getenv("TOTP_ISSUER"),
			AccountName: fmt.Sprintf("%v", email),
		})
		utility.SessionSet(w, r, utility.Session{Key: "totpSecret", Value: key.Secret()})

		img, _ := key.Image(400, 400)

		var buf bytes.Buffer
		png.Encode(&buf, img)

		data := b64.StdEncoding.EncodeToString([]byte(buf.String()))
		data = "data:image/png;base64," + data

		utility.RenderTemplate(w, r, "twoFactorRegister", data)
	} else if r.Method == "POST" {
		err := r.ParseForm()
		if err != nil {

		} else {
			verifyToken := r.FormValue("challengeCode")
			validationResult := totp.Validate(verifyToken, fmt.Sprintf("%v", utility.SessionGet(r, "totpSecret")))
			if validationResult {
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
