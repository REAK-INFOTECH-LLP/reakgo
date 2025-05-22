# reakgo


Simple Framework to quickly build webapps in GoLang

## Running the App

Install gin

`go install github.com/codegangsta/gin@latest`

and run

`go run server.go`

## Struct srtickly according to db schema

`MyStruct` is a struct representing a database table schema`(required)`, designed to work seamlessly with an Object-Relational Mapping (ORM) system. It strictly adheres to the database schema, ensuring that only fields corresponding to the database columns are included in the struct. Additionally, it uses struct tags, including `primaryKey:"true"`, to provide metadata for the ORM's functionality.

### Fields:

- `Id` (int64): An integer field representing the primary key of the table. The `primaryKey:"true"` tag indicates that this field is the primary key of the table.

- `Name` (string): A string field representing a name in the database schema.

- `Age` (int): An integer field representing an age value in the database schema.

- `Email` (string): A string field representing an email address in the database schema.

- `PhoneNumber` (int64): An integer field representing a phone number in the database schema.

### Primary Key:

The `Id` field of the `MyStruct` struct is marked as the primary key of the table using the `primaryKey:"true"` tag. This tag is essential for ORM systems to correctly identify the primary key of the table and perform operations like retrieval, updates, and deletes based on primary key values.

### ORM Usage:

By adhering strictly to the database schema and utilizing struct tags like `primaryKey:"true"`, `MyStruct` is optimized for use with ORM systems. You can easily use it for seamless integration between your Go application and your relational database, simplifying data manipulation and retrieval operations.

### Example Usage:

```go
// This is how an example struct will look like which is according to db schema.
type MyStruct struct {
	Id          int64  `json:"id" db:"id" primarykey:"true" `
	Name        string `json:"name" db:"name" `
	Age         int    `json:"age" db:"age"`
	Email       string `json:"email" db:"email"`
	PhoneNumber int64  `json:"phone_number" db:"phone_number"`
}

+---------------+--------------+------+-----+---------+----------------+
| Field         | Type         | Null | Key | Default | Extra          |
+---------------+--------------+------+-----+---------+----------------+
| id            | int64        | NO   | PRI | NULL    | auto_increment |
| name          | varchar(255) | YES  |     | NULL    |                |
| age           | int          | YES  |     | NULL    |                |
| email         | varchar(255) | YES  |     | NULL    |                |
| phone_number  | int64        | YES  |     | NULL    |                |
+---------------+--------------+------+-----+---------+----------------+
```
# Creating Strict Database Schema-Compliant Structures

To work efficiently with a database schema, it is advisable to create dedicated structures for each database table. These structures should strictly follow the schema defined in the database. The primary principle is to have no additional fields or keys in these structures apart from those defined in the database schema.

## The Importance of Strict Compliance

Strict compliance with the database schema offers numerous advantages:

1. **Seamless ORM Integration:** These compliant structures are optimized for integration with Object-Relational Mapping (ORM) systems. ORM systems rely on the structure of these objects to perform various operations like data retrieval, updates, and deletions.

2. **Reduced Complexity:** By adhering strictly to the database schema, you simplify the interaction between your Go application and the database. There's no need to handle additional fields that are not part of the schema.

## Creating Structures

For each table in your database schema, create a corresponding structure in your Go code. Ensure that the fields in these structures mirror the columns in the respective tables. Additionally, you can use struct tags to provide metadata about the fields, such as indicating the primary key.

## Extending or Adding Keys

If you need to extend the functionality or add extra keys to a structure, create a new structure to accommodate these changes. This new structure should embed the original database schema-compliant structure. In this way, you maintain separation between the strict schema structure and the extended version.


## FindFirst Function

The `FindFirst` function retrieves the first record from a specified database table based on the primary key field of a provided structure. This function is particularly useful when you want to retrieve the row with the smallest primary key value from a table.

### Parameters:

- `tableName` (string): The name of the database table from which to retrieve the record.
- `structure` (interface{}): A pointer to a struct (e.g., &MyStruct{}) representing the structure into which the retrieved data will be scanned. The primary key field of the struct is used to determine the record to retrieve.

### Return Value:

- `error`: An error is returned if any of the following conditions are met:
  - The provided structure is not a pointer to a struct.
  - The primary key field is not found in the struct or is missing the "primarykey" tag.
  - There is an issue with executing the SQL query or scanning the result.

### Functionality:

- The `FindFirst` function queries the specified database table to retrieve the first record based on the primary key field of the provided structure.

- It checks if the `structure` argument is a pointer to a struct. If not, it returns an error.

- It identifies the primary key field in the provided struct using the "primarykey" tag.

- It constructs an SQL query to retrieve the first record from the specified `tableName` based on the primary key.

- The query results are scanned into the provided `structure`, which should be a pointer to a struct.

### Example Usage:

```go
var result MyStruct
err := FindFirst("your_table", &result)
if err != nil {
    // Handle the error
}
```

## FindLast Function

The `FindLast` function retrieves the last record from a specified database table based on the primary key field of a provided structure. This function is particularly useful when you want to retrieve the row with the largest primary key value from a table.

### Parameters:

- `tableName` (string): The name of the database table from which to retrieve the record.
- `structure` (interface{}): A pointer to a struct (e.g., &MyStruct{}) representing the structure into which the retrieved data will be scanned. The primary key field of the struct is used to determine the record to retrieve.

### Return Value:

- `error`: An error is returned if any of the following conditions are met:
  - The provided structure is not a pointer to a struct.
  - The primary key field is not found in the struct or is missing the "primarykey" tag.
  - There is an issue with executing the SQL query or scanning the result.

### Functionality:

- The `FindLast` function queries the specified database table to retrieve the last record based on the primary key field of the provided structure.

- It checks if the `structure` argument is a pointer to a struct. If not, it returns an error.

- It identifies the primary key field in the provided struct using the "primarykey" tag.

- It constructs an SQL query to retrieve the last record from the specified `tableName` based on the primary key.

- The query results are scanned into the provided `structure`, which should be a pointer to a struct.

### Example Usage:

```go
var result MyStruct
err := FindLast("your_table", &result)
if err != nil {
    // Handle the error
}
```
## Find Function

The `Find` function queries a database table based on the provided criteria and scans the result into a slice of structs.

### Parameters:

- `data` (`map[string]interface{}`): A map containing the query criteria.
  - `tablename` (string, required): The name of the database table to query.
  - `columnname` (string, required): The name of the column to filter on.
  - `columnvalue` (interface{}, required): The value to filter the `columnname` by.
  - `sortcolumn` (string, optional): The name of the column to sort the results by. Defaults to the primary key of the table defined in the structure.
  - `sortvalue` (string, optional): The sorting order, which can be "ASC" (ascending) or "DESC" (descending). Defaults to "ASC".
  - `showcolumn`([]string,optional):The column which data you want to retrive only by default its *.

### Return Value:

- `error`: An error indicating success or failure. Returns `nil` on success.

### Functionality:

- It checks if the required keys (`tablename`, `columnname`, `columnvalue`) are present in the `data` map. If any of these keys are missing, it returns an error.

- It validates that the `structure` argument is a pointer to a slice of structs. If the validation fails, it returns an error.

- If the `sortcolumn` key is not provided in the `data` map, it automatically determines the `sortcolumn` based on the primary key of the table defined in the structure.

- If the `sortvalue` key is not provided in the `data` map, it defaults to "ASC" (ascending) for sorting.

- It constructs a SQL query based on the provided criteria and executes the query on the database.

- The query retrieves rows from the specified `tablename` where the value in the `columnname` matches `columnvalue`. The results are sorted by `sortcolumn` in the specified order (`sortvalue`).

- The query results are scanned into the provided `structure`, which should be a pointer to a slice of structs.

### Example Usage:

```go
data := map[string]interface{}{
    "tablename":   "your_table",
    "columnname":  "name",
    "columnvalue": "John",
    "sortcolumn":  "age",
    "sortvalue":   "DESC",
    "showcolumn":[]string{"name","age"},
}

var result []YourStruct
err := Find(data, &result)
if err != nil {
    // Handle the error
}
```
# Insert Function

The `Insert` function simplifies the process of adding a new record to a database table.

## Parameters:

- `tablename` (string): The name of the database table where you want to insert the record.

- `dataStruct` (interface{}): This parameter should be a pre-declared struct with fields that match the attributes present in your database table.

## Return Value:

- `error`: The `Insert` function returns an error in the following situations:

  1. The provided `tablename` does not exist in the database.

  2. The `dataStruct`  was not strickly according to db schema.

  3. An error occurs while executing the SQL query.

If you encounter one of the above two errors after reading this documentation, perhaps it's time for a coffee break!
### Example Usage:

```go
tablename=   "your_table"
type MyStruct struct {
	Id          int64  `json:"id" db:"id" primarykey:"true" `
	Name        string `json:"name" db:"name" `
	Age         int    `json:"age" db:"age"`
	Email       string `json:"email" db:"email"`
	PhoneNumber int64  `json:"phone_number" db:"phone_number"`
}
var result MyStruct
    result.Name= "name"
    result.Age= 14
    result.Email  "email@google.com"
    result.Phonenumber 1234567890

err := Insert(tablename, &result)
if err != nil {
    // Handle the error
}
```
# Delete Function

The `Delete` function allows you to delete records from a database table based on specified criteria.

## Parameters:

- `data` (map[string]interface{}): A map containing the following required keys:
    - `tablename` (string): The name of the database table from which you want to delete records.
    - `columnname` (string): The name of the column to use as the filter criterion.
    - `columnvalue` (interface{}): The value to compare against the `columnname` for filtering.

## Return Value:

- `bool`: The `Delete` function returns `true` if one or more records were successfully deleted, and `false` if no records were deleted.

- `error`: An error is returned in the following situations:
    - Any of the required keys (`tablename`, `columnname`, `columnvalue`) is missing in the `data` map.
    - An error occurs while executing the SQL query.

The `Delete` function is useful for removing records from a database table based on specific conditions.

### Example Usage:

```go
data := map[string]interface{}{
    "tablename":   "your_table",
    "columnname":  "name",
    "columnvalue": "John",
}

var result []YourStruct
err := Delete(data, &result)
if err != nil {
    // Handle the error
}
```
# Update Function

The `Update` function allows you to modify a row in a specified database table based on a unique identifier (primary key) that you provide in your struct.

## Parameters:

- `tablename` (string): The name of the database table where you want to update the record.

- `structure` (interface{}): The `structure` parameter represents a pre-declared struct containing fields that correspond to the columns in your database table. You need to pass in the struct with the desired updates you want to apply to that row.

    * You must ensure that the unique identifier (primary key) value is provided within the struct. If you're not sure how to define a field as a primary key in the struct, refer to the documentation for struct tags.

## Return Value:

- `error`: An error is returned only in the following situations:
    1. You specify a `tablename` that does not exist.
    2. The value of the field defined as the primary key is either empty or invalid (non-existent).

The `Update` function simplifies the process of modifying a specific row in a database table based on the provided unique identifier.
### Example Usage:

```go
tablename=   "your_table"
type MyStruct struct {
	Id          int64  `json:"id" db:"id" primarykey:"true" `
	Name        string `json:"name" db:"name" `
	Age         int    `json:"age" db:"age"`
	Email       string `json:"email" db:"email"`
	PhoneNumber int64  `json:"phone_number" db:"phone_number"`
}
var result MyStruct
    result.Name= "name"
    result.Age= 14
    result.Email  "email@google.com"
    result.Phonenumber 1234567890

err := Update(tablename, &result)
if err != nil {
    // Handle the error
}
```
