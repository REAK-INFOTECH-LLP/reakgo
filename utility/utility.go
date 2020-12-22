package utility

import (
    "net/http"
    "html/template"
    "database/sql"
    "github.com/gorilla/sessions"
    "fmt"
)

// Template Pool
var View *template.Template
// Session Store
var Store *sessions.CookieStore
// DB Connections
var Db *sql.DB

type Session struct {
    Key string
    Value string
}

func RedirectTo(w http.ResponseWriter, r *http.Request) {
    fmt.Println("Redirect")
    http.Redirect(w, r, "https://google.com", 302)
}

func SessionSet(w http.ResponseWriter, r *http.Request, data Session) {
    session, _ := Store.Get(r, "session-name")
    // Set some session values.
    session.Values[data.Key] = data.Value
    // Save it before we write to the response/return from the handler.
    err := session.Save(r, w)
    fmt.Println(err)
}

func SessionGet(r *http.Request, key string) interface{} {
    session, _ := Store.Get(r, "session-name")
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
