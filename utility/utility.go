package utility

import (
	"encoding/json"
	"log"
	"os"

	//"log"
	//"fmt"
	"html/template"
	"net/http"

	"github.com/gorilla/sessions"
	"github.com/jmoiron/sqlx"
)

// Template Pool
var View *template.Template

// Session Store
var Store *sessions.FilesystemStore

// DB Connections
var Db *sqlx.DB

type Session struct {
	Key   string
	Value interface{}
}

type Flash struct {
	Type    string
	Message string
}

// {
// 	"status": "failure",
// 	"message": "Incorrect credentials, Please try again.",
// 	"payload": {},
//   }
type AjaxResponce struct {
	Status  string
	Message string
	Payload interface{}
}

func RedirectTo(w http.ResponseWriter, r *http.Request, path string) {
	http.Redirect(w, r, os.Getenv("APP_URL")+"/"+path, http.StatusFound)
}

func SessionSet(w http.ResponseWriter, r *http.Request, data Session) {
	session, _ := Store.Get(r, os.Getenv("SESSION_NAME"))
	// Set some session values.
	session.Values[data.Key] = data.Value
	// Save it before we write to the response/return from the handler.
	err := session.Save(r, w)
	if err != nil {
		log.Println(err)
	}
}

func SessionGet(r *http.Request, key string) interface{} {
	session, _ := Store.Get(r, os.Getenv("SESSION_NAME"))
	// Set some session values.
	return session.Values[key]
}

// func fetchSession(r *http.Request) map[interface{}]interface{} {
// 	session, _ := Store.Get(r, os.Getenv("SESSION_NAME"))
// 	return session.Values
// }

// check value exist in array or not
// func stringInSlice(a string, list []string) bool {
// 	for _, b := range list {
// 		if b == a {
// 			return true
// 		}
// 	}
// 	return false
// }
func CheckACL(w http.ResponseWriter, r *http.Request, useracl []string) bool {
	// userTypeSession := fmt.Sprintf("%v", SessionGet(r, "type"))
	// if stringInSlice(userTypeSession, useracl) || stringInSlice(userTypeSession, useracl)  {
	return true
	// } else {
	// 	RedirectTo(w, r, "forbidden")
	// 	return false
	// }
}

func AddFlash(flavour string, message string, w http.ResponseWriter, r *http.Request) {
	session, err := Store.Get(r, os.Getenv("SESSION_NAME"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	//flash := make(map[string]string)
	//flash["Flavour"] = flavour
	//flash["Message"] = message
	flash := Flash{
		Type:    flavour,
		Message: message,
	}
	err = session.Save(r, w)
	session.AddFlash(flash, "message")
	log.Println("flash", flash)
	if err != nil {
		log.Println(err)
	}
}

func viewFlash(w http.ResponseWriter, r *http.Request) interface{} {

	session, err := Store.Get(r, os.Getenv("SESSION_NAME"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	fm := session.Flashes("message")
	if fm == nil {
		return nil
	}
	session.Save(r, w)
	return fm
}

/* if is http@header["reak-api"] return true otherwise false*/
func IsCurlApiRequest(r *http.Request) bool {
	return r.Header.Get("reak-api") == "true"
}

/* if isCurlAPI w.Write json otherwise ExcuteTemplate() */
func RenderTemplate(w http.ResponseWriter, r *http.Request, template string, data interface{}) {
	tmplData := make(map[string]interface{})
	tmplData["data"] = data
	tmplData["flash"] = viewFlash(w, r)
	// tmplData["session"] = fetchSession(r)
	if IsCurlApiRequest(r) {
		jsonresponce, err := json.Marshal(tmplData)
		if err != nil {
			log.Println(err)
		}
		w.Write([]byte(jsonresponce))
	} else {
		View.ExecuteTemplate(w, template, tmplData)
	}
}

func IsCurrectAdminOrUser(w http.ResponseWriter, r *http.Request, typeOfUsr string) bool {
	return true
}
