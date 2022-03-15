package utility

import (
    "os"
    "log"
    //"log"
    //"fmt"
    "net/http"
    "html/template"
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
    Key string
    Value string
}

type Flash struct {
    Type string
    Message string
}

func RedirectTo(w http.ResponseWriter, r *http.Request, path string){
    http.Redirect(w, r, os.Getenv("APP_URL")+"/"+path, http.StatusFound)
}

func SessionSet(w http.ResponseWriter, r *http.Request, data Session) {
    session, _ := Store.Get(r, os.Getenv("SESSION_NAME"))
    // Set some session values.
    session.Values[data.Key] = data.Value
    // Save it before we write to the response/return from the handler.
    err := session.Save(r, w)
    if(err != nil){
        log.Println(err)
    }
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
        RedirectTo(w, r, "forbidden")
        return false
    }
}

func AddFlash(flavour string, message string, w http.ResponseWriter, r *http.Request){
    session, err := Store.Get(r, os.Getenv("SESSION_NAME"))
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
    }
    //flash := make(map[string]string)
    //flash["Flavour"] = flavour
    //flash["Message"] = message
    flash := Flash {
        Type: flavour,
        Message: message,
    }
    session.AddFlash(flash, "message")
    err = session.Save(r, w)
    if (err != nil){
        log.Println(err)
    }
}

func viewFlash(w http.ResponseWriter, r *http.Request) interface{}{
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

func RenderTemplate(w http.ResponseWriter, r *http.Request, template string, data interface{}){
    session, _ := Store.Get(r, os.Getenv("SESSION_NAME"))
    tmplData := make(map[string]interface{})
    tmplData["data"] = data
    tmplData["flash"] = viewFlash(w, r)
    tmplData["session"] = session.Values["email"]
    View.ExecuteTemplate(w, template, tmplData)
}
