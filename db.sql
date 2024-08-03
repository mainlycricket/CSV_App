-- CREATE TABLE TypeTest
CREATE TABLE "TypeTest" (
	 "__ID" SERIAL PRIMARY KEY, 
	 "Bool" boolean,
	 "Bool_Arr" boolean[] DEFAULT array[true, false, true]::boolean[],
	 "Date" date,
	 "DateTime" timestamptz,
	 "Date_arr" date[],
	 "Datetime_Arr" timestamptz[],
	 "Float" real CHECK ( "Float" > -100 AND "Float" < 100 ) DEFAULT -5.34,
	 "Float_arr" real[] DEFAULT array[4.5, 4.64]::real[],
	 "Int" integer DEFAULT 4,
	 "Int_Arr" integer[] DEFAULT array[4, 5]::integer[],
	 "Str_Arr" text[] DEFAULT array['Hi', 'Bro']::text[],
	 "String" text DEFAULT 'Text',
	 "Time" time,
	 "Time_Arr" time[]);

-- CREATE TABLE branches
CREATE TABLE "branches" (
	 "Branch_Id" integer PRIMARY KEY,
	 "Branch_Name" text NOT NULL,
	 "Course_Id" integer NOT NULL,
	 "Teachers" text[]);

-- CREATE TABLE courses
CREATE TABLE "courses" (
	 "Course_Id" integer PRIMARY KEY,
	 "Course_Name" text NOT NULL UNIQUE,
	 "Lateral_Allowed" boolean);

-- CREATE TABLE students
CREATE TABLE "students" (
	 "Branch_Id" integer NOT NULL,
	 "Course_Id" integer NOT NULL,
	 "Student_Father" text NOT NULL,
	 "Student_Id" integer PRIMARY KEY,
	 "Student_Name" text NOT NULL);

-- CREATE TABLE subjects
CREATE TABLE "subjects" (
	 "Branch_Id" integer NOT NULL,
	 "Subject_Id" integer PRIMARY KEY,
	 "Subject_Name" text NOT NULL);

-- real Array Validator Function
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

-- integer Array Validator Function
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

-- TypeTest Table Validator Trigger Function
CREATE OR REPLACE FUNCTION validate_TypeTest_trigger()
RETURNS TRIGGER AS $$
DECLARE
res TEXT;

BEGIN

res := validate_real_arr(NEW."Float_arr", false, 2, 2, 4, 5, NULL);
IF res != '' THEN
RAISE EXCEPTION 'Error in "Float_arr" column in "TypeTest" table: %', res;
END IF;

res := validate_integer_arr(NEW."Int_Arr", false, NULL, NULL, NULL, NULL, array[1, 2, 3, 4, 5]::integer[]);
IF res != '' THEN
RAISE EXCEPTION 'Error in "Int_Arr" column in "TypeTest" table: %', res;
END IF;

RETURN NEW;

END;
$$ LANGUAGE plpgsql;

-- TypeTest Table Trigger
CREATE TRIGGER validate_table_TypeTest_trigger
BEFORE INSERT OR UPDATE ON "TypeTest"
FOR EACH ROW
EXECUTE FUNCTION validate_TypeTest_trigger();

-- DATA INSERTION "TypeTest"
INSERT INTO "TypeTest"("Int","String","Float","Date","Time","DateTime","Bool","Int_Arr","Str_Arr","Float_arr","Date_arr","Time_Arr","Datetime_Arr","Bool_Arr")
VALUES
(1, 'Ram', 4.1, '2024-01-01', '14:30:00', NULL, false, array[1, 2]::integer[], array['val', 'val2']::text[], array[4, 4]::real[], array['2024-07-01', '2024-07-01']::date[], array['14:30:00', '14:30:00']::time[], array['2024-07-01T12:30:00+05:30', '2024-07-01T12:30:00+05:30']::timestamptz[], array[true, false]::boolean[]),
(1, 'Ram', 4.1, '2024-01-01', '14:30:00', '2024-07-01T12:30:00+05:30', true, array[1, 2]::integer[], array['val', 'val2']::text[], array[4, 4]::real[], array['2024-07-01', '2024-07-01']::date[], array['14:30:00', '14:30:00']::time[], array['2024-07-01T12:30:00+05:30', '2024-07-01T12:30:00+05:30']::timestamptz[], array[true, false]::boolean[]),
(1, 'Ram', 4.1, '2024-01-01', '14:30:00', NULL, true, array[1, 2]::integer[], array['val', 'val2']::text[], array[4, 4]::real[], array['2024-07-01', '2024-07-01']::date[], array['14:30:00', '14:30:00']::time[], array['2024-07-01T12:30:00+05:30', '2024-07-01T12:30:00+05:30']::timestamptz[], array[true, false]::boolean[]);

-- DATA INSERTION "branches"
INSERT INTO "branches"("Branch_Id","Branch_Name","Course_Id","Teachers")
VALUES
(1, 'Computer Science', 1, array['HA', 'PC']::text[]),
(2, 'Information Technology', 1, array['LD', 'RK']::text[]),
(3, 'Civil Engineering', 1, NULL);

-- DATA INSERTION "courses"
INSERT INTO "courses"("Course_Id","Course_Name","Lateral_Allowed")
VALUES
(1, 'B. Tech.', true),
(2, 'M. Tech.', false);

-- DATA INSERTION "students"
INSERT INTO "students"("Student_Id","Student_Name","Student_Father","Course_Id","Branch_Id")
VALUES
(1, 'Tushar', 'Ajay', 1, 1),
(2, 'Akshay', 'Nand', 1, 1),
(3, 'Saurabh', 'Jagganath', 1, 2),
(4, 'Harsh', 'Ramesh', 1, 2);

-- DATA INSERTION "subjects"
INSERT INTO "subjects"("Subject_Id","Subject_Name","Branch_Id")
VALUES
(1, 'DS', 1),
(2, 'COA', 1),
(3, 'WT', 2),
(4, 'Java', 2);

-- branches Table Foreign Keys
ALTER TABLE "branches"
ADD CONSTRAINT "branches_Course_Id_fkey" FOREIGN KEY ("Course_Id")
REFERENCES "courses" ("Course_Id");

-- students Table Foreign Keys
ALTER TABLE "students"
ADD CONSTRAINT "students_Branch_Id_fkey" FOREIGN KEY ("Branch_Id")
REFERENCES "branches" ("Branch_Id"),
ADD CONSTRAINT "students_Course_Id_fkey" FOREIGN KEY ("Course_Id")
REFERENCES "courses" ("Course_Id");

-- subjects Table Foreign Keys
ALTER TABLE "subjects"
ADD CONSTRAINT "subjects_Branch_Id_fkey" FOREIGN KEY ("Branch_Id")
REFERENCES "branches" ("Branch_Id");

