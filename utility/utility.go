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

type MyStruct struct {
	Id          int64  `json:"id" db:"id" primarykey:"true" `
	Name        string `json:"name" db:"name" `
	Age         int    `json:"age" db:"age"`
	Email       string `json:"email" db:"email"`
	PhoneNumber int64  `json:"phone_number" db:"phone_number"`
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
func Find(data map[string]interface{}, structure interface{}) error {
	// Check if the required keys exist in the 'data' map.
	tableName, tableNameExists := data["tablename"].(string)
	columnName, columnNameExists := data["columnname"].(string)
	columnValue, columnValueExists := data["columnvalue"]
	sortColumn, sortColumnExists := data["sortcolumn"].(string)
	sortValue, sortValueExists := data["sortvalue"].(string)

	// Check if the required keys are missing and handle the error condition.
	if !tableNameExists || !columnNameExists || !columnValueExists {
		return errors.New("missing required keys in map[string]interface{} any of this tablename,columnname,columnvalue")
	}

	// Ensure that 'structure' is a pointer to a slice of structs.
	structureValue := reflect.ValueOf(structure)
	if structureValue.Kind() != reflect.Ptr || structureValue.Elem().Kind() != reflect.Slice {
		return errors.New("interface{} must be a pointer to a slice of structs")
	}

	if !sortColumnExists {
		// Assuming structure is a pointer to a slice of structs.
		valueType := structureValue.Elem().Type().Elem()
		emptyStruct := reflect.New(valueType).Interface()
		// Call the PrimaryKeyIdentifier function with the first element.
		primaryKeyField, err := PrimaryKeyIdentifier(emptyStruct)
		if err != nil {
			return err
		} else {
			sortColumn = primaryKeyField
		}
	}

	if !sortValueExists {
		sortValue = "ASC"
	}

	query := fmt.Sprintf("SELECT * FROM %s WHERE %s = ? ORDER BY %s %s ;", tableName, columnName, sortColumn, sortValue)

	// Execute the query and scan the results into the provided structure slice.
	err := Db.Select(structure, query, columnValue)
	return err
}

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
