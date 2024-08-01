

CREATE TABLE "TypeTest" ("__ID" SERIAL PRIMARY KEY, "Bool" boolean,"Bool_Arr" boolean[] DEFAULT array['true', 'false', 'true']::boolean[],"Date" date,"DateTime" timestamptz,"Date_arr" date[],"Datetime_Arr" timestamptz[],"Float" real CHECK ( "Float" > -100 AND "Float" < 100 ) DEFAULT -5.34,"Float_arr" real[] DEFAULT array[4.5, 4.64]::real[],"Int" integer DEFAULT 4,"Int_Arr" integer[] DEFAULT array[4, 5]::integer[],"Str_Arr" text[] DEFAULT array['Hi', 'Bro']::text[],"String" text DEFAULT 'Text',"Time" time,"Time_Arr" time[]);

CREATE TABLE "branches" ("Branch_Id" integer PRIMARY KEY,"Branch_Name" text NOT NULL,"Course_Id" integer NOT NULL,"Teachers" text[]);

CREATE TABLE "courses" ("Course_Id" integer PRIMARY KEY,"Course_Name" text NOT NULL UNIQUE,"Lateral_Allowed" boolean);

CREATE TABLE "students" ("Branch_Id" integer NOT NULL,"Course_Id" integer NOT NULL,"Student_Father" text NOT NULL,"Student_Id" integer PRIMARY KEY,"Student_Name" text NOT NULL);

CREATE TABLE "subjects" ("Branch_Id" integer NOT NULL,"Subject_Id" integer PRIMARY KEY,"Subject_Name" text NOT NULL);

ALTER TABLE "branches"
ADD CONSTRAINT "branches_Course_Id_fkey" FOREIGN KEY ("Course_Id")
REFERENCES "courses" ("Course_Id");

ALTER TABLE "students"
ADD CONSTRAINT "students_Branch_Id_fkey" FOREIGN KEY ("Branch_Id")
REFERENCES "branches" ("Branch_Id"),
ADD CONSTRAINT "students_Course_Id_fkey" FOREIGN KEY ("Course_Id")
REFERENCES "courses" ("Course_Id");

ALTER TABLE "subjects"
ADD CONSTRAINT "subjects_Branch_Id_fkey" FOREIGN KEY ("Branch_Id")
REFERENCES "branches" ("Branch_Id");

CREATE FUNCTION validate_real_arr(
    arr real[] DEFAULT NULL,
    not_null boolean DEFAULT FALSE,
    min_arr_len integer DEFAULT NULL,
    max_arr_len integer DEFAULT NULL,
    min_ind real DEFAULT NULL, 
    max_ind real DEFAULT NULL,
    enum_arr real[] DEFAULT NULL)
RETURNS text AS $$
DECLARE
    val real;
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
            RETURN FORMAT('Each element value should be at least %s', min_ind);
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
            RETURN FORMAT('Each element value should be at least %s', min_ind);
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

CREATE OR REPLACE FUNCTION validate_TypeTest_trigger()
RETURNS TRIGGER AS $$
DECLARE
res TEXT;

BEGIN

res := validate_real_arr(NEW."Float_arr", false, NULL, NULL, NULL, NULL, NULL, 2, 2, 4, 5);
IF res != '' THEN
RAISE EXCEPTION 'Error in "Float_arr" column in "TypeTest" table: %', res;
END IF;

res := validate_integer_arr(NEW."Int_Arr", false, NULL, NULL, NULL, NULL, NULL, array[1, 2, 3, 4, 5]::integer[]);
IF res != '' THEN
RAISE EXCEPTION 'Error in "Int_Arr" column in "TypeTest" table: %', res;
END IF;

RETURN NEW;

END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER validate_table_TypeTest_trigger
BEFORE INSERT OR UPDATE ON "TypeTest"
FOR EACH ROW
EXECUTE FUNCTION validate_TypeTest_trigger();