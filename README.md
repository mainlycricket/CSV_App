### Data Types

- integer
- float
- boolean
- string
- date
- time
- datetime
- array of these primitive types

```sql
CREATE OR REPLACE FUNCTION validate_time_arr(time[], min_arr, max_arr, min_ind, max_ind)
RETURNS boolean AS $$
DECLARE
    t time;
BEGIN
	IF array_length($1, 1) < min_arr OR array_length($1, 1) > max_arr THEN
		RETURN FALSE;
	END IF;
    FOREACH t IN ARRAY $1 LOOP
        IF t < min_ind OR t > max_ind THEN
            RETURN FALSE;
        END IF;
    END LOOP;
    RETURN TRUE;
END;
$$ LANGUAGE plpgsql;

CREATE TABLE "AllTypes" (
	"ID" serial PRIMARY KEY,
	"Int_" integer DEFAULT 1 CHECK("Int_" >= 1 AND "Int_" <= 50),
	"int_arr" integer[] DEFAULT array[1, 2, 3]::integer[],
	"Float_" real DEFAULT 4.3 CHECK ("Float_" IN (4.3, 5.1)),
	"Float_arr" real[] DEFAULT array[3.4]::real[],
	"Bool_" bool DEFAULT false,
	"Bool_arr" bool[] DEFAULT array[true, false]::bool[],
	"String" text DEFAULT 'text' REFERENCES "DataTypeTest"("String") CHECK(LENGTH("String") > 1 AND LENGTH("String") < 8),
	"String_arr" text[] DEFAULT array['text1', 'text2']::text[],
	"Date_" date DEFAULT '2024-01-01',
	"Date_Arr" date[] DEFAULT array['2024-07-01']::date[],
	"Time_" time DEFAULT '13:45:59',
	"Time_arr" time[] DEFAULT array['10:00:00']::time[] CHECK(validate_time_arr("Time_arr")),
	"Datetime" timestamptz DEFAULT '2024-07-01T12:30:00+05:30',
	"Datetime_arr" timestamptz[] DEFAULT array['2024-07-01T12:30:00+05:30']::timestamptz[]
);
```

### Schema Creation

- No validation
- Always successful unless a parsing error is encountered by `encoding/csv`
- The generated schema should be examined
- Special constraints like 'Min', 'Max', 'Enums' & 'Default' are always required to be added manually
- DB Name should be also set manually

### Column Constraints

- Composite Primary Keys are not allowed

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
