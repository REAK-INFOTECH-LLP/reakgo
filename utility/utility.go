package utility

import (
	"fmt"
	"log"
	"os"
	"reflect"
	"strings"

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

//structure variable is having the &struct.please pass the pointer of the structure not the variable
func FindFirst(tableName string, structure interface{}) error {
	query := fmt.Sprintf("SELECT * FROM %s ORDER BY id LIMIT 1", tableName)
	err := Db.Get(structure, query)
	return err
}

//structure variable is having the &struct.please pass the pointer of the structure not the variable
func FindLast(tableName string, structure interface{}) error {
	query := fmt.Sprintf("SELECT * FROM %s ORDER BY id DESC LIMIT 1", tableName)
	err := Db.Get(structure, query)
	return err
}

//structure variable is having the &[]struct
func Find(tableName string, columnName string, columnValue interface{}, structure interface{}) error {
	query := fmt.Sprintf("SELECT * FROM %s WHERE %s= ?", tableName, columnName)
	// Execute the query and scan the results into the provided structure slice.
	err := Db.Select(structure, query, columnValue)
	return err
}

//provide the tablenName and the primarykey which is id.
func Delete(tableName string, id int) (bool, error) {
	query := fmt.Sprintf("DELETE * FROM %s WHERE id= =:id", tableName)
	result, err := Db.NamedExec(query, map[string]interface{}{"id": id})
	if err != nil {
		log.Println(err)
	} else {
		Rowefffect, _ := result.RowsAffected()
		return Rowefffect > 0, err
	}
	return false, err

}

// InsertDynamic inserts a new record into the specified table based on the struct values.
func Insert(tableName string, dataStruct interface{}) error {
	// Get the type and value of the dataStruct.
	valueType := reflect.TypeOf(dataStruct)
	value := reflect.ValueOf(dataStruct)

	// Build the INSERT INTO statement dynamically based on the struct and table name.
	var columns []string
	var placeholders []string
	values := make([]interface{}, 0)

	for i := 0; i < valueType.NumField(); i++ {
		field := valueType.Field(i)
		valueField := value.Field(i)

		columns = append(columns, field.Name)
		placeholders = append(placeholders, "?")
		values = append(values, valueField.Interface())
	}

	insertQuery := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)", tableName, strings.Join(columns, ", "), strings.Join(placeholders, ", "))

	// Execute the dynamic INSERT query.
	_, err := Db.Exec(insertQuery, values...)
	return err
}

// UpdateDynamic updates a record in the specified table based on the struct values.
// It excludes the primary key column from the update.
func update(tableName string, dataStruct interface{}, primaryKeyColumn string) error {
	// Get the type and value of the dataStruct.
	valueType := reflect.TypeOf(dataStruct)
	value := reflect.ValueOf(dataStruct)

	// Build the SET clause for the UPDATE statement dynamically based on the struct.
	var setValues []string
	values := make([]interface{}, 0)

	for i := 0; i < valueType.NumField(); i++ {
		field := valueType.Field(i)
		valueField := value.Field(i)
		fieldName := field.Name

		// Exclude the primary key column from the update.
		if fieldName != primaryKeyColumn {
			setValues = append(setValues, fmt.Sprintf("%s = ?", fieldName))
			values = append(values, valueField.Interface())
		}
	}

	setClause := strings.Join(setValues, ", ")

	// Build the SQL UPDATE statement.
	updateQuery := fmt.Sprintf("UPDATE %s SET %s WHERE %s = ?",
		tableName, setClause, primaryKeyColumn)

	// Get the primary key value.
	primaryKeyValue := value.FieldByName(primaryKeyColumn).Interface()
	values = append(values, primaryKeyValue)

	// Execute the dynamic UPDATE query.
	result, err := Db.Exec(updateQuery, values...)
	if err != nil {
		return err
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("no rows were updated")
	}

	return nil
}
