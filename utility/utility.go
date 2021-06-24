package utility

import (
    "os"
    "fmt"
    "log"
    "net/http"
    "html/template"
    "github.com/gorilla/sessions"
    "github.com/jmoiron/sqlx"
)

// Template Pool
var View *template.Template
// Session Store
var Store *sessions.CookieStore
// DB Connections
var Db *sqlx.DB

type Session struct {
    Key string
    Value string
}

func RedirectTo(w http.ResponseWriter, r *http.Request) {
    fmt.Println("Redirect")
    http.Redirect(w, r, "https://google.com", 302)
}

func SessionSet(w http.ResponseWriter, r *http.Request, data Session) {
    session, _ := Store.Get(r, os.Getenv("SESSION_NAME"))
    // Set some session values.
    session.Values[data.Key] = data.Value
    // Save it before we write to the response/return from the handler.
    err := session.Save(r, w)
    fmt.Println(err)
}

func SessionGet(r *http.Request, key string) interface{} {
    session, _ := Store.Get(r, os.Getenv("SESSION_NAME"))
    // Set some session values.
    return session.Values[key]
}


func CheckACL(w http.ResponseWriter, r *http.Request, minLevel int) bool {
    userType := SessionGet(r, "type")
    var level int = 0
    switch(userType){
    case "user":
        level = 1
    case "admin":
        level = 2
    default:
        level = 0
    }
    if(level >= minLevel){
        return true
    } else {
        RedirectTo(w, r)
        return false
    }
}

func AddFlash(flavour string, message string, w http.ResponseWriter, r *http.Request){
    session, err := Store.Get(r, "flash-session")
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
    }
    flash := make(map[string]string)
    flash["Flavour"] = flavour
    flash["Message"] = message
    session.AddFlash(flash, "message")
    session.Save(r, w)
}

func viewFlash(w http.ResponseWriter, r *http.Request) interface{}{
  session, err := Store.Get(r, "flash-session")
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

func RenderTemplate(w http.ResponseWriter, r *http.Request, template string, data interface{}){
    tmplData := make(map[string]interface{})
    tmplData["data"] = data
    tmplData["flash"] = viewFlash(w, r)
    View.ExecuteTemplate(w, template, tmplData)
    log.Println(tmplData)
}
