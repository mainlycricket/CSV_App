### Data Types

- integer
- float
- boolean
- string
- date
- time
- datetime
- array of these primitive types

### Schema Creation

- No validation
- Always successful unless a parsing error is encountered by `encoding/csv`
- The generated schema should be examined
- Special constraints like 'Min', 'Max', 'Enums' & 'Default' are always required to be added manually
- DB Name should be also set manually

### Column Constraints

- Composite Primary Keys are not allowed
- Array columns can't be primary keys

#### Min & Max

- Should always be mentioned in strings
- Should be empty for boolean values or if they aren't required

- Data is validated by value for integer, float, date, time, datetime
- E.g. `Min:"3"` and `Max:"10"` for `integer` mean `3 >= value <= 10`

- Data is validated by length for strings (should be a +ve integer)
- E.g. `Min:"3"` and `Max:"10"` for `string` mean `3 >= length(value) <= 10`

- For array: "array_length,individual_value_constraint"
- E.g. `Min: "2,3"` and `Max: "5,10"` for `[]integer` implies `2 >= len(arr) <= 5` and `3 >= each_element <= 10`
- E.g. `Min: "2,3"` and `Max: "4,10"` for `[]string` implies `2 >= len(arr) <= 4` and `3 >= length(each_element) <= 10`

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
- Presence of extra characters in constraints are ignored
- Eg: `P:columnName`, `U:columnName`, `NU:columnName`, `FN:columnName`, `FNU:columnName`

#### Foreign Key Mapping

- If a column is marked as a foreign key, the referenced column is mapped with the following idea:
