package models

import (
	"database/sql"
	"errors"
	"fmt"
	"reakgo/utility"
	"reflect"
	"strings"
)

var ColumnMap = make(map[string]string)
var (
	// ErrCode is a config or an internal error
	ErrCode = errors.New("Case statement in code is not correct.")
	// ErrNoResult is a not results error
	ErrNoResult = errors.New("Result not found.")
	// ErrUnavailable is a database not available error
	ErrUnavailable = errors.New("Database is unavailable.")
	// ErrUnauthorized is a permissions violation
	ErrUnauthorized = errors.New("User does not have permission to perform this operation.")
)

// standardizeErrors returns the same error regardless of the database used
func standardizeError(err error) error {
	if err == sql.ErrNoRows {
		return ErrNoResult
	}

	return err
}

//This is an example of a structured format that is compatible with Object-Relational Mapping (ORM) systems.
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
	errFromDb := utility.Db.Get(structure, query)
	return errFromDb
}

//structure variable is having the &struct.please pass the pointer of the structure not the variable
func FindLast(tableName string, structure interface{}) error {
	primaryKeyField, err := PrimaryKeyIdentifier(structure)
	if err != nil {
		return err
	}
	query := fmt.Sprintf("SELECT * FROM %s ORDER BY  %s DESC LIMIT 1", tableName, primaryKeyField)
	errFromDb := utility.Db.Get(structure, query)
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
	colInterface, colExists := data["showcolumn"]

	// Check if the required keys are missing and handle the error condition.
	if !tableNameExists || !columnNameExists || !columnValueExists {
		return errors.New("missing required keys any of this tablename,columnname,columnvalue")
	}
	//checking the columnName provided is matching to the database columnName or not
	_, exist := ColumnMap[columnName]
	if !exist {
		return errors.New("columnn does not exist .Please check the columnname")
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
	//checking the sorting columnName provided is matching to the database columnName or not.
	_, ok := ColumnMap[sortColumn]
	if !ok {
		return errors.New("sort columnn does not exist .Please check the columnname")
	}

	if !sortValueExists {
		sortValue = "ASC"
	}
	// Check if 'col' key exists and is a list of column names
	var columns []string
	if colExists {
		col, isStringSlice := colInterface.([]string)
		if isStringSlice {
			for _, colName := range col {
				// Check if the value exists in the ColumnMap.
				if _, ok := ColumnMap[colName]; !ok {
					return errors.New("Column '" + colName + "' does not exist. Please check the column name.")
				}
			}
			columns = col
		} else {
			return errors.New("showcolumn key having the wrong datatype .It should be []string")
		}
	} else {
		columns = append(columns, "*")
	}

	query := fmt.Sprintf("SELECT %s FROM %s WHERE %s = ? ORDER BY %s %s ;", strings.Join(columns, ", "), tableName, columnName, sortColumn, sortValue)

	// Execute the query and scan the results into the provided structure slice.
	err := utility.Db.Select(structure, query, columnValue)
	return err
}

//provide the tablenName and the primarykey which is id.
func Delete(data map[string]interface{}) (bool, error) {
	tableName, tableNameExists := data["tablename"].(string)
	columnName, columnNameExists := data["columnname"].(string)
	columnValue, columnValueExists := data["columnvalue"]
	// Check if the required keys are missing and handle the error condition.
	if !tableNameExists || !columnNameExists || !columnValueExists {
		return false, errors.New("missing required keys any of this tablename,columnname,columnvalue")
	}
	//checking the columnName provided is matching to the database columnName or not
	_, exist := ColumnMap[columnName]
	if !exist {
		return false, errors.New("columnn does not exist .Please check the columnname")
	}
	query := fmt.Sprintf("DELETE FROM %s WHERE %s= :value", tableName, columnName)
	result, err := utility.Db.NamedExec(query, map[string]interface{}{"value": columnValue})
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
	_, err := utility.Db.Exec(insertQuery, values...)
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
	result, err := utility.Db.Exec(updateQuery, values...)
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
func ListTables() ([]string, error) {
	rows, err := utility.Db.Query("SHOW TABLES")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tables []string
	for rows.Next() {
		var table string
		if err := rows.Scan(&table); err != nil {
			return nil, err
		}
		tables = append(tables, table)
	}

	return tables, nil
}
func ListColumns(table string) ([]string, error) {
	query := fmt.Sprintf("SELECT column_name FROM information_schema.columns WHERE table_name = '%s'", table)
	rows, err := utility.Db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var columns []string
	for rows.Next() {

		var column string
		if err := rows.Scan(&column); err != nil {
			return nil, err
		}
		columns = append(columns, column)
		// Add the column and its column to the map
		ColumnMap[column] = column

	}

	return columns, nil
}
