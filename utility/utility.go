package utility

import (
	"errors"
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
	primaryKeyField, err := PrimaryKeyIdentifier(structure)
	if err != nil {
		return err
	}
	// Build the query using the determined primary key field.
	query := fmt.Sprintf("SELECT * FROM %s ORDER BY %s LIMIT 1", tableName, primaryKeyField)
	// Execute the query and scan the result into the provided structure.
	errFromDb := Db.Get(structure, query)
	return errFromDb
}

//structure variable is having the &struct.please pass the pointer of the structure not the variable
func FindLast(tableName string, structure interface{}) error {
	primaryKeyField, err := PrimaryKeyIdentifier(structure)
	if err != nil {
		return err
	}
	query := fmt.Sprintf("SELECT * FROM %s ORDER BY  %s DESC LIMIT 1", tableName, primaryKeyField)
	errFromDb := Db.Get(structure, query)
	return errFromDb
}

//data varibale need to have the three compulsory keys (tablename,columnname,columnvalue),structure variable is having the &[]struct
// func Find(data map[string]interface{}, structure interface{}) error {
// 	// Check if the required keys exist in the 'data' map.
// 	tableName, tableNameExists := data["tablename"].(string)
// 	columnName, columnNameExists := data["columnname"].(string)
// 	columnValue, columnValueExists := data["columnvalue"]
// 	sortColumn, sortColumnExists := data["sortcolumn"].(string)
// 	sortValue, sortValueExists := data["sortvalue"].(string)

// 	// Check if the required keys are missing and handle the error condition.
// 	if !tableNameExists || !columnNameExists || !columnValueExists {
// 		return errors.New("missing required keys in map[string]interface{} any of this tablename,columnname,columnvalue")
// 	}
// 	if !sortColumnExists {
// 		// Assuming structut is a pointer to a slice of structs.
// 		sliceValue := reflect.ValueOf(structure).Elem()
// 		log.Println(reflect.ValueOf(structure).Elem().Type())
// 		// Call the PrimaryKeyIdentifier function with the first element.
// 		primaryKeyField, err := PrimaryKeyIdentifier(controllers.MyStruct{})
// 		if err != nil {
// 			return err
// 		}
// 		log.Println(primaryKeyField)
// 		// Check if the slice is not empty.
// 		if sliceValue.Len() > 0 {
// 			// Get the first element from the slice.
// 			// firstElement := sliceValue.Index(0).Addr().Interface()

// 			// Call the PrimaryKeyIdentifier function with the first element.
// 			primaryKeyField, err := PrimaryKeyIdentifier(controllers.MyStruct{})
// 			if err != nil {
// 				return err
// 			}
// 			log.Println(primaryKeyField)
// 			sortColumn = primaryKeyField
// 		}
// 	}
// 	// 	primaryKeyField, err := PrimaryKeyIdentifier(structure)
// 	// 	if err != nil {
// 	// 		return err
// 	// 	}
// 	// 	sortColumn = primaryKeyField
// 	// }
// 	if !sortValueExists {
// 		sortValue = "ASC"
// 	}
// 	query := fmt.Sprintf("SELECT * FROM %s WHERE %s = ? ORDER BY %s %s ;", tableName, columnName, sortColumn, sortValue)
// 	log.Println(query)
// 	// Execute the query and scan the results into the provided structure slice.
// 	err := Db.Select(structure, query, columnValue)
// 	return err
// }

//provide the tablenName and the primarykey which is id.
func Delete(data map[string]interface{}) (bool, error) {
	tableName, tableNameExists := data["tablename"].(string)
	columnName, columnNameExists := data["columnname"].(string)
	columnValue, columnValueExists := data["columnvalue"]
	// Check if the required keys are missing and handle the error condition.
	if !tableNameExists || !columnNameExists || !columnValueExists {
		return false, errors.New("missing required keys in map[string]interface{} any of this tablename,columnname,columnvalue")
	}
	query := fmt.Sprintf("DELETE FROM %s WHERE %s= :value", tableName, columnName)
	result, err := Db.NamedExec(query, map[string]interface{}{"value": columnValue})
	if err != nil {
		// return false, err
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

		// Get the database column name from the struct tag, if available.
		columnName := field.Tag.Get("db")

		if columnName == "" {
			// If no db tag is specified, use the field name as the column name.
			columnName = field.Name
		}

		columns = append(columns, columnName)
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
func Update(tableName string, structure interface{}) error {
	var primaryKeyColumnName string
	var setValues []string
	// Get the type and value of the structure.
	valueType := reflect.TypeOf(structure)
	value := reflect.ValueOf(structure)
	// Build the SET clause for the UPDATE statement dynamically based on the struct.
	values := make([]interface{}, 0)

	for i := 0; i < valueType.NumField(); i++ {
		field := valueType.Field(i)
		valueField := value.Field(i)
		// Get the database column name from the struct tag, if available.
		columnName := field.Tag.Get("db")

		if columnName == "" {
			// If no db tag is specified, use the field name as the column name.
			columnName = field.Name
		}
		// fieldName := field.Name
		primaryName := field.Tag.Get("primarykey")
		// Exclude the primary key column from the update.
		if primaryName != "true" {
			setValues = append(setValues, fmt.Sprintf("%s = ?", columnName))
			values = append(values, valueField.Interface())
		}
		if primaryName == "true" {
			primaryKeyColumnName = field.Name
			values = append(values)
		}
	}
	setClause := strings.Join(setValues, ", ")

	// Check if the primary key column name is empty.
	if primaryKeyColumnName == "" {
		return errors.New("primary key column name not found")
	}

	// Build the SQL UPDATE statement.
	updateQuery := fmt.Sprintf("UPDATE %s SET %s WHERE %s = ?",
		tableName, setClause, primaryKeyColumnName)

	// Get the primary key value.
	primaryKeyField := value.FieldByName(primaryKeyColumnName)
	if !primaryKeyField.IsValid() {
		return errors.New("primary key field not found")
	}
	primaryKeyValue := primaryKeyField.Interface()
	values = append(values, primaryKeyValue)
	// Execute the dynamic UPDATE query.
	result, err := Db.Exec(updateQuery, values...)
	if err != nil {
		return err
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return errors.New("no rows were updated")
	}

	return nil
}
func PrimaryKeyIdentifier(structure interface{}) (string, error) {
	// Get the value of the structure (should be a pointer to a struct).
	structValue := reflect.ValueOf(structure)

	// Ensure that the provided structure is a pointer to a struct.
	if structValue.Kind() != reflect.Ptr || structValue.Elem().Kind() != reflect.Struct {
		return "", errors.New("invalid interface, must be a pointer to a struct")
	}

	// Get the type and value of the structure.
	structType := structValue.Elem().Type()
	primaryKeyField := ""

	// Iterate through the struct fields to find the primary key field.
	for i := 0; i < structType.NumField(); i++ {
		field := structType.Field(i)
		primaryKeyTag := field.Tag.Get("primarykey")
		if primaryKeyTag == "true" {
			primaryKeyField = field.Name
			return primaryKeyField, nil
		}
	}
	return "", errors.New("primary key field not found in the struct")
}
