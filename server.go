package main

import (
    "net/http"
    "html/template"
    "path/filepath"
    "log"
    "os"
    "strings"
    "time"
    
    "reakgo/router"
    "reakgo/utility"
    
    "github.com/gorilla/sessions"
    "github.com/joho/godotenv"
    _ "github.com/go-sql-driver/mysql"
    "github.com/jmoiron/sqlx"
)

func init() {
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
    utility.Store = sessions.NewCookieStore([]byte(os.Getenv("SESSION_KEY")))
    utility.View = cacheTemplates()
    // See "Important settings" section.
    utility.Db.SetConnMaxLifetime(time.Minute * 3)
    utility.Db.SetMaxOpenConns(10)
    utility.Db.SetMaxIdleConns(10)

}

func main() {
    http.HandleFunc("/", handler)
    // Serve static assets
    http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("./assets"))))
    log.Fatal(http.ListenAndServe(":"+os.Getenv("WEB_PORT"), nil))
}

func cacheTemplates() *template.Template {
    templ := template.New("")
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
