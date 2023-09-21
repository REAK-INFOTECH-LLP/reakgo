# reakgo
Simple Framework to quickly build webapps in GoLang


## FindFirst Function

The FindFirst function retrieves the first record from a specified database table based on the primary key field of a provided structure. This function is particularly useful when you want to retrieve the row with the smallest primary key value from a table.
### Parameters:
    tableName (string): The name of the database table from which to retrieve the record.
    structure (interface{}): A pointer to a struct (e.g., &MyStruct{}) representing the structure into which the retrieved data will be scanned. The primary key field of the struct is used to determine the record to retrieve.

### Return Value:
    error:
    An error is returned if any of the following conditions are met:
    The provided structure is not a pointer to a struct.
    The primary key field is not found in the struct or is missing the "primarykey" tag.
    There is an issue with executing the SQL query or scanning the result.

## FindLast Function

The FindLast function retrieves the Last record from a specified database table based on the primary key field of a provided structure. This function is particularly useful when you want to retrieve the row with the Largest primary key value from a table.
### Parameters:
    tableName (string): The name of the database table from which to retrieve the record.
    structure (interface{}): A pointer to a struct (e.g., &MyStruct{}) representing the structure into which the retrieved data will be scanned. The primary key field of the struct is used to determine the record to retrieve.

### Return Value:
    error:
    An error is returned if any of the following conditions are met:
    The provided structure is not a pointer to a struct.
    The primary key field is not found in the struct or is missing the "primarykey" tag.
    There is an issue with executing the SQL query or scanning the result.
