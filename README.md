### About
This project generates a single `.sql` file (PostgreSQL database) for the provided `.csv` files in `./data` directory.

My more ambitious aim is to automatically create a CRUP app for these tables.

#### Step 1: Schema Generation

```bash
go build . && ./CSV_App schema
```

- This generates initial schema in `data/schema.json` by analyzing the CSV data
- Always successful unless an error is encountered while parsing `.csv` files
- No data validation is performed
- Schema constraints should be reviewed
- Special constraints like `Min`, `Max`, `Enums` & `Default` are always required to be added manually
- DB Name should be also set manually
- Table Names & Column Names should NOT be changed

#### Step 2: SQL File Generation

```bash
go build . && ./CSV_App sql
```

- This generates initial schema in `data/db.sql` for the given CSV files
- CSV Data is validated against the constraints provided in `schema.json`

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

### Data Types

- integer
- float
- boolean
- text
- date
- time
- datetime (with timezone)
- array of these primitive types


### Column Constraints

#### Primary Key

- Composite Primary Keys are not allowed
- If no primary key is provided, a column named `__ID` is added by default

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

#### Constraint Codes:

| Code | Constraint  |
| ---- | ----------- |
| N    | Not Null    |
| U    | Unique      |
| P    | Primary Key |
| F    | Foreign Key |

#### Notes:

- The constraints should be mentioned before the column name
- Constraint code & column name should be separated by a colon (:)
- If a colon is present in a column label, the left part will be treated as constraint
- Presence of 'P' overrides all the other constraints, P also sets the Not Null & Unique
- Extra unwanted characters in constraints are ignored
- Eg: `P:columnName`, `U:columnName`, `NU:columnName`, `FN:columnName`, `FNU:columnName`

#### Foreign Key Mapping

If a column is marked as a foreign key, the referenced column is mapped with the following idea:

- if any other table has a column with the same name and datatype, it is marked as referenced column
- if multiple such columns exist, any one is choosen
- if no match is found, the `ForeignTable` & `ForeignColumn` fields are left with value `__`
