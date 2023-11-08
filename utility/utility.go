package utility

import (
	"crypto/rand"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math/big"
	"os"
	"reflect"
	"strconv"

	//"log"
	//"fmt"
	"html/template"
	"net/http"

	"github.com/allegro/bigcache/v3"
	"github.com/gorilla/csrf"
	"github.com/gorilla/sessions"
	"github.com/jmoiron/sqlx"
)

// Template Pool
var View *template.Template

// Session Store
var Store *sessions.FilesystemStore

// DB Connections
var Db *sqlx.DB

// Cache
var Cache *bigcache.BigCache

// CSRF
var CSRF func(http.Handler) http.Handler

type Session struct {
	Key   string
	Value interface{}
}

type Flash struct {
	Type    string
	Message string
}

type Helper interface {
	GenerateRandomString(n int) (string, error)
	RedirectTo(w http.ResponseWriter, r *http.Request, path string)
	SessionGet(r *http.Request, key string) interface{}
	SessionSet(w http.ResponseWriter, r *http.Request, data Session)
	AddFlash(flavour string, message string, w http.ResponseWriter, r *http.Request)
	viewFlash(w http.ResponseWriter, r *http.Request) interface{}
	RenderTemplate(w http.ResponseWriter, r *http.Request, template string, data interface{})
	ParseDataFromPostRequestToMap(r *http.Request) (map[string]interface{}, error)
	ParseDataFromJsonToMap(r *http.Request) (map[string]interface{}, error)
	StrictParseDataFromJson(r *http.Request, structure interface{}) error
	StrictParseDataFromPostRequest(r *http.Request, structure interface{}) error
	RenderJsonResponse(w http.ResponseWriter, r *http.Request, data interface{})
	RenderTemplateData(w http.ResponseWriter, r *http.Request, template string, data interface{})
	StringInArray(target string, arr []string) bool
	ReturnUserDetails(r *http.Request, user interface{}) error
}

func GenerateRandomString(n int) (string, error) {
	const letters = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz-"
	ret := make([]byte, n)
	for i := 0; i < n; i++ {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(letters))))
		if err != nil {
			return "", err
		}
		ret[i] = letters[num.Int64()]
	}

	return string(ret), nil
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
	tmplData["csrf"] = csrf.TemplateField(r)
	View.ExecuteTemplate(w, template, tmplData)
}

func ParseDataFromPostRequestToMap(r *http.Request) (map[string]interface{}, error) {
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

func ParseDataFromJsonToMap(r *http.Request) (map[string]interface{}, error) {
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
}

func StrictParseDataFromJson(r *http.Request, structure interface{}) error {
	err := json.NewDecoder(r.Body).Decode(structure)
	if err != nil {
		return err
	}
	return err

}

func StrictParseDataFromPostRequest(r *http.Request, structure interface{}) error {
	err := r.ParseForm()
	if err != nil {
		return err
	}
	// Use reflection to set field values based on form data
	structValue := reflect.ValueOf(structure)
	if structValue.Kind() != reflect.Ptr || structValue.Elem().Kind() != reflect.Struct {
		return errors.New("invalid argument: 'interface{}' must be a pointer to a struct")
	}
	structElem := structValue.Elem()
	for key, values := range r.Form {
		field := structElem.FieldByName(key)
		if !field.IsValid() {
			// Skip fields that don't exist in the structure
			continue
		}
		// Handle fields with different types (e.g., slice or single value)
		if len(values) == 1 {
			value := values[0]
			//conversion of the data came according to the feilds values that are used in the struct
			switch field.Kind() {
			case reflect.Int:
				intValue, err := strconv.Atoi(value)
				if err != nil {
					return fmt.Errorf("Error in converting string to int: %v", err)
				}
				field.SetInt(int64(intValue))
			case reflect.Int8:
				int8Value, err := strconv.ParseInt(value, 10, 8)
				if err != nil {
					return fmt.Errorf("Error in converting string to int8: %v", err)
				}
				field.SetInt(int64(int8(int8Value)))
			case reflect.Int16:
				int16Value, err := strconv.ParseInt(value, 10, 16)
				if err != nil {
					return fmt.Errorf("Error in converting string to int16: %v", err)
				}
				field.SetInt(int64(int16(int16Value)))
			case reflect.Int32:
				int32Value, err := strconv.ParseInt(value, 10, 32)
				if err != nil {
					return fmt.Errorf("Error in converting string to int32: %v", err)
				}
				field.SetInt(int64(int32(int32Value)))
			case reflect.Int64:
				int64Value, err := strconv.ParseInt(value, 10, 64)
				if err != nil {
					return fmt.Errorf("Error in converting string to int64: %v", err)
				}
				field.SetInt(int64Value)
			case reflect.Float32:
				float32Value, err := strconv.ParseFloat(value, 32)
				if err != nil {
					return fmt.Errorf("Error in converting string to float32: %v", err)
				}
				field.SetFloat(float64(float32(float32Value)))
			case reflect.Float64:
				float64Value, err := strconv.ParseFloat(value, 64)
				if err != nil {
					return fmt.Errorf("Error in converting string to float64: %v", err)
				}
				field.SetFloat(float64(float64Value))
			case reflect.Bool:
				boolValue, err := strconv.ParseBool(value)
				if err != nil {
					return fmt.Errorf("Error in converting string to bool: %v", err)
				}
				field.SetBool(boolValue)
			case reflect.String:
				field.SetString(value)
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
func RenderJsonResponse(w http.ResponseWriter, r *http.Request, data interface{}) {
	jsonresponce, err := json.Marshal(data)
	if err != nil {
		log.Println(err)
	}
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE")
	w.Write([]byte(jsonresponce))
}

func RenderTemplateData(w http.ResponseWriter, r *http.Request, template string, data interface{}) {
	session, _ := Store.Get(r, os.Getenv("SESSION_NAME"))
	tmplData := make(map[string]interface{})
	tmplData["data"] = data
	tmplData["flash"] = viewFlash(w, r)
	tmplData["session"] = session.Values["email"]
	View.ExecuteTemplate(w, template, tmplData)
}

func StringInArray(target string, arr []string) bool {
	// Can change to slices.Contain if we're targetting 1.21+
	for _, s := range arr {
		if s == target {
			return true
		}
	}
	return false
}
