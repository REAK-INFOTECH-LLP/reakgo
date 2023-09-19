package utility

import (
	"encoding/json"
	"errors"
	"log"
	"os"
	"reflect"

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

func CheckACL(w http.ResponseWriter, r *http.Request, minLevel int) bool {
	userType := SessionGet(r, "type")
	var level int = 0
	switch userType {
	case "user":
		level = 1
	case "admin":
		level = 2
	default:
		level = 0
	}
	if level >= minLevel {
		return true
	} else {
		RedirectTo(w, r, "forbidden")
		return false
	}
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
	session.AddFlash(flash, "message")
	err = session.Save(r, w)
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

func RenderTemplate(w http.ResponseWriter, r *http.Request, template string, data interface{}) {
	session, _ := Store.Get(r, os.Getenv("SESSION_NAME"))
	tmplData := make(map[string]interface{})
	tmplData["data"] = data
	tmplData["flash"] = viewFlash(w, r)
	tmplData["session"] = session.Values["email"]
	View.ExecuteTemplate(w, template, tmplData)
}

func ParseDataFromRequest(r *http.Request) (map[string]interface{}, error) {
	if os.Getenv("APP_IS") == "monolith" {
		var jsonDataMap map[string]interface{}
		result := make(map[string]interface{})

		// Read JSON data from the HTTP request body
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&jsonDataMap)
		if err != nil {
			return result, err
		}

		// Close the request body to prevent resource leaks
		defer r.Body.Close()
		// Iterate through the JSON values
		for key, value := range jsonDataMap {
			result[key] = value
		}

		return result, nil
	} else if os.Getenv("APP_IS") == "microservice" {
		formData := make(map[string]interface{})
		err := r.ParseForm()
		if err != nil {
			return formData, err
		}

		// Iterate through the form values
		for key, values := range r.Form {
			// If there's only one value for the key, store it directly
			if len(values) == 1 {
				formData[key] = values[0]
			} else {
				// If there are multiple values, store them as a slice
				formData[key] = values
			}
		}

		return formData, nil
	}
	return nil, errors.New("The value of the environment variable APP_IS is invalid or Not set. It should be [monolith or microservice]. Please update the environment variable and try again.")
}

func StrictParseDataFromRequest(r *http.Request, structure interface{}) error {
	if os.Getenv("APP_IS") == "monolith" {
		err := json.NewDecoder(r.Body).Decode(structure)
		if err != nil {
			return err
		}
		return err

	} else if os.Getenv("APP_IS") == "microservice" {

		err := r.ParseForm()
		if err != nil {
			return err
		}

		// Use reflection to set field values based on form data
		structValue := reflect.ValueOf(structure)
		if structValue.Kind() != reflect.Ptr || structValue.Elem().Kind() != reflect.Struct {
			return errors.New("invalid structure, must be a pointer to a struct")
		}
		structElem := structValue.Elem()
		structType := structElem.Type()
		log.Println("structtype=", structType)
		for key, values := range r.Form {
			field := structElem.FieldByName(key)
			if !field.IsValid() {
				continue // Skip fields that don't exist in the structure
			}

			// Handle fields with different types (e.g., slice or single value)
			if len(values) == 1 {
				log.Println(reflect.ValueOf(values[0]), "value")
				log.Println(field.Type(), "feild")
				log.Println(reflect.ValueOf(values[0]).Type(), "type")
				value := reflect.ValueOf(values[0])
				if value.Type().ConvertibleTo(field.Type()) {
					field.Set(value.Convert(field.Type()))
				}
			} else {
				// Handle fields with multiple values as a slice (if field is a slice)
				if field.Kind() == reflect.Slice {
					sliceType := field.Type().Elem()
					slice := reflect.MakeSlice(field.Type(), len(values), len(values))

					for i, v := range values {
						elemValue := reflect.ValueOf(v)
						if elemValue.Type().ConvertibleTo(sliceType) {
							slice.Index(i).Set(elemValue.Convert(sliceType))
						}
					}

					field.Set(slice)
				}
			}
		}
		return err
	}
	return errors.New("The value of the environment variable APP_IS is invalid or Not set. It should be [monolith or microservice]. Please update the environment variable and try again.")
}

// func StrictParseDataFromRequest(r *http.Request, structure struct{}) (struct{}, error) {
// 	if os.Getenv("APP_IS") == "monolith" {
// 		err := json.NewDecoder(r.Body).Decode(&structure)
// 		if err != nil {
// 			return structure, err
// 		}
// 		return structure, err

// 	} else if os.Getenv("APP_IS") == "microservice" {
// 		formData := make(map[string]interface{})
// 		err := r.ParseForm()
// 		if err != nil {
// 			return structure, err
// 		}

// 		// Iterate through the form values
// 		for key, values := range r.Form {
// 			// If there's only one value for the key, store it directly
// 			if len(values) == 1 {
// 				formData[key] = values[0]
// 			} else {
// 				// If there are multiple values, store them as a slice
// 				formData[key] = values
// 			}
// 		}

// 		return structure, nil
// 	}
// 	return nil, errors.New("The value of the environment variable APP_IS is invalid or Not set. It should be [monolith or microservice]. Please update the environment variable and try again.")
// }
