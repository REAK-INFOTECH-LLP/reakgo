package main

import (
	"context"
	"encoding/gob"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"reakgo/models"
	"reakgo/router"
	"reakgo/utility"

	"github.com/allegro/bigcache/v3"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/csrf"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
)

func init() {

    db_user := os.Getenv("DB_USER")
    if(db_user == "") {
        log.Fatal("Missing Env value DB_USER")
    }
    
    db_password := os.Getenv("DB_PASSWORD")
    if(db_user == "") {
        log.Fatal("Missing Env value DB_PASSWORD")
    }
    
    db_name := os.Getenv("DB_NAME")
    if(db_user == "") {
        log.Fatal("Missing Env value DB_NAME")
    }

    session_key := os.Getenv("SESSION_KEY")
    if(db_user == "") {
        log.Fatal("Missing Env value SESSION_KEY")
    }

	// Set log configuration
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	var err error
	err = godotenv.Load()
	if err != nil {
		log.Println(".env file wasn't found, looking at env variables")
	}
	motd()
	// Read Config
	utility.Db, err = sqlx.Open("mysql", db_user+":"+db_password+"@/"+db_name)
	if err != nil {
		log.Println("Wowza !, We didn't find the DB or you forgot to setup the env variables")
		panic(err)
	}
	utility.Store = sessions.NewFilesystemStore("", []byte(session_key))
	utility.Store.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   60 * 1,
		HttpOnly: true,
	}
	utility.View = cacheTemplates()
	// See "Important settings" section.
	utility.Db.SetConnMaxLifetime(time.Minute * 3)
	utility.Db.SetMaxOpenConns(10)
	utility.Db.SetMaxIdleConns(10)
	columnNameReciever()

	gob.Register(utility.Flash{})

}

func main() {
	
    csrf_secret_key := os.Getenv("CSRF_SECRET_KEY")
    if(csrf_secret_key == "") {
        log.Fatal("Missing Env value CSRF_SECRET_KEY")
    }

    app_is := os.Getenv("APP_IS")
    if(app_is == "") {
        log.Fatal("Missing Env value APP_IS")
    }

    web_port := os.Getenv("WEB_PORT")
    if(app_is == "") {
        log.Fatal("Missing Env value WEB_PORT")
    }


    // Initialize Caching
	cacheInit()
	// Generate cache as a go routine as to not halt operation,
	// Cache fail-safe is already implemented so will fetch from DB incase the cache is not populated
	go models.GenerateCache()

	utility.CSRF = csrf.Protect([]byte(csrf_secret_key))

	mux := mux.NewRouter()

	mux.PathPrefix("/").HandlerFunc(handler)

	// Serve static assets
	staticHandler := http.StripPrefix("/assets/", http.FileServer(http.Dir("./assets/")))
	mux.PathPrefix("/assets/").Handler(staticHandler)

	if app_is == "monolith" {
		log.Fatal(http.ListenAndServe(":"+web_port, utility.CSRF(mux)))
	} else if app_is == "microservice" {
		log.Fatal(http.ListenAndServe(":"+web_port, mux))
	}
}

func cacheTemplates() *template.Template {

	funcMap := template.FuncMap{
		// Only to be used for SAFE attributes, SAFE = Computer Generated and not USER DRIVEN
		"attr": func(s string) template.HTMLAttr {
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

func handler(w http.ResponseWriter, r *http.Request) {
	router.Routes(w, r)
}

func columnNameReciever() {
	// Get a list of tables in the database
	tables, err := models.ListTables()
	if err != nil {
		log.Println(err)
		return
	}

	// Iterate through the tables and write column names to the output file
	for _, table := range tables {
		_, err := models.ListColumns(table)
		if err != nil {
			log.Println(err)
			return
		}
	}
}

func motd() {
	logo := `
______ _____  ___   _   __
| ___ \  ___|/ _ \ | | / /
| |_/ / |__ / /_\ \| |/ / 
|    /|  __||  _  ||    \ 
| |\ \| |___| | | || |\  \
\_| \_\____/\_| |_/\_| \_/
                          
----------------------------
Application should now be accessible on port ` + os.Getenv("WEB_PORT") + `

`
	log.Println(logo)
}

func cacheInit() {

	token_cache_size, err := strconv.Atoi(os.Getenv("TOKEN_CACHE_SIZE"))
	if err != nil {
		// Set Standard Size in case we can't convert to int, or value is missing
		token_cache_size = 500
	}
	config := bigcache.Config{
		// number of shards (must be a power of 2)
		Shards: 1024,

		// time after which entry can be evicted
		LifeWindow: 10 * time.Minute,

		// Interval between removing expired entries (clean up).
		// If set to <= 0 then no action is performed.
		// Setting to < 1 second is counterproductive — bigcache has a one second resolution.
		CleanWindow: 5 * time.Minute,

		// rps * lifeWindow, used only in initial memory allocation
		MaxEntriesInWindow: 1000 * 10 * 60,

		// max entry size in bytes, used only in initial memory allocation
		MaxEntrySize: 500,

		// prints information about additional memory allocation
		Verbose: true,

		// cache will not allocate more memory than this limit, value in MB
		// if value is reached then the oldest entries can be overridden for the new ones
		// 0 value means no size limit
		HardMaxCacheSize: token_cache_size,

		// callback fired when the oldest entry is removed because of its expiration time or no space left
		// for the new entry, or because delete was called. A bitmask representing the reason will be returned.
		// Default value is nil which means no callback and it prevents from unwrapping the oldest entry.
		OnRemove: nil,

		// OnRemoveWithReason is a callback fired when the oldest entry is removed because of its expiration time or no space left
		// for the new entry, or because delete was called. A constant representing the reason will be passed through.
		// Default value is nil which means no callback and it prevents from unwrapping the oldest entry.
		// Ignored if OnRemove is specified.
		OnRemoveWithReason: nil,
	}

	utility.Cache, err = bigcache.New(context.Background(), config)
	if err != nil {
		log.Fatal(err)
	}
}
