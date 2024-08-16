### About

- This project generates a single `.sql` file (PostgreSQL database) for the provided `.csv` files in `./data` directory.
- It also generates a basic standalone CRUD API for these tables.
- I am currently focused on adding more features to the generated app.
- Please read this doc thoroughly to understand how it works.

![App Flow](./data/CSV_App.png)

#### Step 0: Remove existing data and copy CSV files

- Clone the project
- Delete the existing `./app` directory
- Remove the files in `./data` directory and copy your CSV files here.

```bash
rm -r app
rm data/*
# copy csv files in ./data
```

#### Step 1: Schema Generation

```bash
go build . && ./CSV_App schema
```

- This generates an initial schema in `data/schema.json` by analyzing the CSV files `./data` directory

- In this step, all the tables, columns and datatypes are identified. Also, special constraints can also be set (read ahead to understand).

- Always successful unless:

  - an error is encountered while parsing `.csv` files
  - duplicate table name or column name name is found
  - empty table name or column name is found

- No data validation is performed

- Constraints in `data/schema.json` should be reviewed

- Special constraints like `Min`, `Max`, `Enums` & `Default` are always required to be set manually

- Checkout the sample schema file [here](./data/test.json)

- Table Names & Column Names should NOT be changed as they are **_sanitized_**

> [!NOTE]
> Leading & Trailing spaces are trimmed in sanitization.
> Any sequence of non-alphanumeric character is replaced by a single underscore.
> E.g. `!'Table Name'!` is transformed into `_Table_Name_`
> Hence, two tables or their columns can't have the same sanitized name

#### Step 2: SQL File Generation

```bash
go build . && ./CSV_App sql
```

- This generates the SQL file in `data/db.sql` for the given CSV files
- CSV Data is validated against the constraints provided in `schema.json`
- Checkout the sample SQL file [here](./data/db.sql)

#### Step 3: Execute SQL file

Run following commands to create database and execute the `db.sql` file

```bash
psql -h localhost -U postgres -c 'CREATE DATABASE "DB_Name"'
psql -h localhost -U postgres -d "DB_Name" -f data/db.sql
```

#### Step 4: Create Application

```bash
go build . && ./CSV_App app
```

- This generates a standalone app in `./app` directory
- Checkout the sample app here [here](./app/)

### Data Types

- integer
- float
- boolean
- text
- date
- time
- datetime (with timezone)
- array of these primitive types

> [!NOTE]
> Checkout [TypeTest.csv](./data/TypeTest.csv) to understand data formats in CSV files.
> If a column doesn't have any value, its datatype is marked as empty in `schema.json` and it should be set manually.

### Column Constraints

#### Special Constraint Codes:

| Code | Constraint  |
| ---- | ----------- |
| N    | Not Null    |
| U    | Unique      |
| P    | Primary Key |
| F    | Foreign Key |

#### Notes:

- The constraints should be mentioned before the column name
- Constraint code & column name should be separated by a colon (:)
- Presence of 'P' overrides all the other constraints, P also sets the Not Null & Unique
- Extra unwanted characters in constraints are ignored
- Eg: `P:columnName`, `U:columnName`, `NU:columnName`, `FN:columnName`, `FNU:columnName`

> [!NOTE]
> If a colon is present in a column label, the left part will be treated as constraint

#### Primary Key

- Composite Primary Keys are not allowed
- If no primary key is provided, a column named `__ID` is added by default in SQL file and the generated app but `schema.json`remains unmodified
- All datatypes are allowed to be primary key

#### Foreign Key Mapping

If a column is marked as a foreign key, the referenced column is mapped with the following idea:

- if any other table has a column with the same name and datatype, it is marked as referenced column
- if multiple such columns exist, any one is choosen
- if no match is found, the `ForeignTable` & `ForeignColumn` fields are left with value `__`

#### Min & Max

- Should always be mentioned in strings
- Should be empty for boolean values or if they aren't required

- Data is validated by value for integer, float, date, time, datetime
- E.g. `Min:"3"` and `Max:"10"` for `integer` mean `3 >= value <= 10`

- Data is validated by length for strings (should be a +ve integer)
- E.g. `Min:"3"` and `Max:"10"` for `text` mean `3 >= length(value) <= 10`

- For array: "array_length,individual_value_constraint"
- E.g. `Min: "2,3"` and `Max: "5,10"` for `integer[]` implies `2 >= len(arr) <= 5` and `3 >= each_element <= 10`
- E.g. `Min: "2,3"` and `Max: "4,10"` for `text[]` implies `2 >= len(arr) <= 4` and `3 >= length(each_element) <= 10`

#### Enums

- Enums should be an array that specifies the allowed values for the column
- Each indiviual value in Enums should satisfy the individual `Min` and `Max` constraints
- For array columns, each values_arr[i] should be present in enums array

#### Default

- Default value should satisy the min, max and enums constraint if they're present
- Primary Key & Unique columns shouldn't have a default value

### CRUD App

- It is a standalone app, it runs on its own.

- Remember to modify the .env file before building or starting the app.

- Each API response has following structure in JSON format:

  | Field   | Description                                               |
  | ------- | --------------------------------------------------------- |
  | success | boolean flag                                              |
  | message | string message                                            |
  | data    | array for GET all, object for Single GET, null for others |

- For each table, there are five API routes:

  - `POST /tableName` - insert a single row
  - `GET /tableName` - get all rows
  - `GET /tableNameByPK` - get single row by primary key
  - `PUT /tableName` - update single row by primary key
  - `DELETE /tableName` - delete single row by primary key

- Checkout the sample Postman collection for sample tables [here](https://documenter.getpostman.com/view/25403102/2sA3s7iUZc)

#### 1. POST /tableName

- The request body should have JSON data
- The value for fallback primary key `__ID`, in case of no primary key in table schema, will be added automatically by PostgreSQL
- Data is validated against constraint by PostgreSQL, not at the application level as of now

#### 2. GET /tableName

- All the table data is returned in JSON format
- Foreign Key columns data is looked-up from the referenced table
- Data filtering & sorting is supported by query params

- **Data Filtering:**

  - For data filtering, the keys should have the same name as column names in schema

    ```
    ?int_field=1,2,3&bool_field=true&date_field=2024-01-25,2024-02-15

    ?float_arr=1.2,2.6&time_arr=15:26:59,07:56:20

    ?str_field=text1&str_field=text2&str_field=text3

    ?str_arr=text1&str_arr=text2&str_arr=text3
    ```

  - OR clause is applied for different values of a particular column
  - AND clause is applied among different columns
  - For array fields, if the array contains any of the passed values, the condition is true
  - For string fields including string array fields, the values should be passed in separate pairs

#### 3. GET /tableNameByPK

- It selects a single row based on the primary key received as query param `?id=value`
- The key should always be `id` regardless of the primary key column name in the table
- For array columns except string arrays, send values as `localhost:8080/tableByPK?id=1,2,3`
- For string columns, send values as individual string elements as `localhost:8080/tableByPK?id=text1&id=text2&id=text3`

#### 4. PUT /tableName

- It updates a single row based on the primary key received as query param `?id=value`
- Query param should be passed in the same manner as `GET /tableNameByPK` route
- The request body should be the same as the POST request but the fallback primary key `__ID` is required even if it is not modified

#### 5. DELETE /tableName

- It deletes a single row based on the primary key received as query param `?id=value`
- Query param should be passed in the same manner as `GET /tableNameByPK` route
