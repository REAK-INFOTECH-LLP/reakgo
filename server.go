package main

import (
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
    "encoding/gob"

	"reakgo/router"
	"reakgo/utility"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/sessions"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
)

func init() {
    // Set log configuration
    log.SetFlags(log.LstdFlags | log.Lshortfile)
    var err error
    err = godotenv.Load()
    if err != nil {
        log.Println(".env file wasn't found, looking at env variables")
    }
    motd()
    // Read Config
    utility.Db, err = sqlx.Open("mysql", os.Getenv("DB_USER")+":"+os.Getenv("DB_PASSWORD")+"@/"+os.Getenv("DB_NAME"))
    if err != nil {
        log.Println("Wowza !, We didn't find the DB or you forgot to setup the env variables")
        panic(err)
    }
    utility.Store = sessions.NewFilesystemStore("",[]byte(os.Getenv("SESSION_KEY")))
    utility.Store.Options = &sessions.Options{
        Path: "/",
        MaxAge: 60*1,
        HttpOnly: true,
    }
    utility.View = cacheTemplates()
    // See "Important settings" section.
    utility.Db.SetConnMaxLifetime(time.Minute * 3)
    utility.Db.SetMaxOpenConns(10)
    utility.Db.SetMaxIdleConns(10)

    gob.Register(utility.Flash{})

}

func main() {
    http.HandleFunc("/", handler)
    // Serve static assets
    http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("./assets"))))
    
    /*
    router := mux.NewRouter()
    router.HandleFunc("/", controllers.BaseIndex)
    router.HandleFunc("/login", controllers.Login)
    router.HandleFunc("/dashboard", controllers.Dashboard)
    */

    log.Fatal(http.ListenAndServe(":"+os.Getenv("WEB_PORT"), nil))
}

func cacheTemplates() *template.Template {

    funcMap := template.FuncMap{
        // Only to be used for SAFE attributes, SAFE = Computer Generated and not USER DRIVEN
        "attr":func(s string) template.HTMLAttr{
            return template.HTMLAttr(s)
        },
        // Only to be used for SAFE HTML, SAFE = Computer Generated and not USER DRIVEN
        "safe": func(s string) template.HTML {
            return template.HTML(s)
         },
        // Only to be used for SAFE URLs, SAFE = Computer Generated and not USER DRIVEN
         "safeURL": func(s string) template.URL {
            return template.URL(s)
         },
    }

    templ := template.New("")
    templ.Funcs(funcMap)
    err := filepath.Walk("./templates", func(path string, info os.FileInfo, err error) error {
        if strings.Contains(path, ".html") {
            _, err = templ.ParseFiles(path)
            if err != nil {
                log.Println(err)
            }
        }

        return err
    })

    if err != nil {
        panic(err)
    }

    return templ
}

func handler(w http.ResponseWriter, r *http.Request){
    router.Routes(w, r)
}

func motd(){
    logo := `
______ _____  ___   _   __
| ___ \  ___|/ _ \ | | / /
| |_/ / |__ / /_\ \| |/ / 
|    /|  __||  _  ||    \ 
| |\ \| |___| | | || |\  \
\_| \_\____/\_| |_/\_| \_/
                          
----------------------------
Application should now be accessible on port `+os.Getenv("WEB_PORT")+`

`
    log.Println(logo)
}
