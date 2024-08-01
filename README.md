### Data Types

- integer
- float
- boolean
- string
- date
- time
- datetime
- array of these primitive types

```bash
psql -h localhost -U postgres -c 'CREATE DATABASE "CSV_App"'
psql -h localhost -U postgres -d "CSV_App" -f db.sql
```

```sql
CREATE FUNCTION validate_integer_arr(
    arr integer[] DEFAULT NULL,
    not_null boolean DEFAULT FALSE,
    min_arr_len integer DEFAULT NULL,
    max_arr_len integer DEFAULT NULL,
    min_ind integer DEFAULT NULL, 
    max_ind integer DEFAULT NULL,
    enum_arr integer[] DEFAULT NULL)
RETURNS text AS $$
DECLARE
    val integer;
BEGIN
    IF arr IS NULL AND not_null THEN
        RETURN 'Empty Array';
    END IF;

    IF arr IS NULL THEN
        RETURN '';
    END IF;

    IF min_arr_len IS NOT NULL AND array_length(arr, 1) < min_arr_len THEN
        RETURN FORMAT('Array length should be at least %s', min_arr_len);
    END IF;

    IF max_arr_len IS NOT NULL AND array_length(arr, 1) > max_arr_len THEN
        RETURN FORMAT('Array length should be at most %s', max_arr_len);
    END IF;

    FOREACH val IN ARRAY arr LOOP
        IF min_ind IS NOT NULL AND val < min_ind THEN
            RETURN FORMAT('Each element value should be at least %s', max_ind);
        END IF;

        IF max_ind IS NOT NULL AND val > max_ind THEN
            RETURN FORMAT('Each element value should be at most %s', max_ind);
        END IF;

        IF enum_arr IS NOT NULL AND val NOT IN (SELECT * FROM unnest(enum_arr)) THEN
            RETURN FORMAT('%s element not present in enums', val);
        END IF;
    END LOOP;
    RETURN '';
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

-- Trigger Function Validator
CREATE OR REPLACE FUNCTION validate_TypeTest_trigger()
RETURNS TRIGGER AS $$
DECLARE
res TEXT;
BEGIN
	res := validate_integer_arr(NEW."Int_Arr", false, 2, 2, 4, 5, NULL);
    IF res != '' THEN
        RAISE EXCEPTION 'Error in "Int_Arr": %', res;
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Final Trigger
CREATE TRIGGER validate_table_TypeTest_trigger
BEFORE INSERT OR UPDATE ON "TypeTest"
FOR EACH ROW
EXECUTE FUNCTION validate_TypeTest_trigger();
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
