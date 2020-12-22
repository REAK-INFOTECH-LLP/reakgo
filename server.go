package main

import (
    "net/http"
    "html/template"
    "path/filepath"
    "log"
    "os"
    "strings"
    "reakgo/router"
    "reakgo/utility"
    "github.com/gorilla/sessions"
    "database/sql"
    _ "github.com/go-sql-driver/mysql"
    "time"
)

func init() {
    // Needs to use OS.Getenv()
    utility.Store = sessions.NewCookieStore([]byte("test"))
    utility.View = cacheTemplates()
    var err error
    utility.Db, err = sql.Open("mysql", "reak:reak@/reakgo")
    if err != nil {
        panic(err)
    }
    // See "Important settings" section.
    utility.Db.SetConnMaxLifetime(time.Minute * 3)
    utility.Db.SetMaxOpenConns(10)
    utility.Db.SetMaxIdleConns(10)

}

func main() {
    http.HandleFunc("/", handler)
    // Serve static assets
    http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("./assets"))))
    log.Fatal(http.ListenAndServe(":4000", nil))
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
