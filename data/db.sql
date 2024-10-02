-- CREATE TABLE TypeTest
CREATE TABLE "TypeTest" (
	 "__ID" SERIAL PRIMARY KEY, 
	 "Bool" boolean,
	 "Bool_Arr" boolean[],
	 "Date" date CHECK ( "Date" IN ('2024-01-01','2024-02-01') ),
	 "DateTime" timestamptz CHECK ( "DateTime" IN ('2024-07-01T12:30:00+05:30') ),
	 "Date_arr" date[],
	 "Datetime_Arr" timestamptz[],
	 "Float" real CHECK ( "Float" IN (4.1,4.65) ),
	 "Float_arr" real[],
	 "Int" integer CHECK ( "Int" IN (1,2) ),
	 "Int_Arr" integer[],
	 "Str_Arr" text[],
	 "String" text CHECK ( "String" IN ('Ram') ),
	 "Time" time CHECK ( "Time" IN ('14:30:00') ),
	 "Time_Arr" time[]);

-- CREATE TABLE branches
CREATE TABLE "branches" (
	 "Branch_Id" text PRIMARY KEY,
	 "Branch_Name" text NOT NULL,
	 "Course_Id" text NOT NULL,
	 "HoD" text,
	 "Teachers" text[],
	 "added_by" text NOT NULL,
	 "college_id" text NOT NULL);

-- CREATE TABLE college
CREATE TABLE "college" (
	 "college_id" text PRIMARY KEY,
	 "college_name" text NOT NULL,
	 "principal_id" text);

-- CREATE TABLE courses
CREATE TABLE "courses" (
	 "Course_Id" text PRIMARY KEY,
	 "Course_Name" text NOT NULL,
	 "Lateral_Allowed" boolean,
	 "added_by" text NOT NULL,
	 "college_id" text NOT NULL);

-- CREATE TABLE login
CREATE TABLE "login" (
	 "added_by" text,
	 "branch_id" text,
	 "college_id" text,
	 "course_id" text,
	 "password" text NOT NULL,
	 "role" text CHECK ( "role" IN ('admin','principal','hod','teacher','student') ) NOT NULL,
	 "username" text PRIMARY KEY);

-- CREATE TABLE students
CREATE TABLE "students" (
	 "Branch_Id" text NOT NULL,
	 "Course_Id" text NOT NULL,
	 "Student_Father" text NOT NULL,
	 "Student_Id" integer PRIMARY KEY,
	 "Student_Name" text NOT NULL,
	 "added_by" text NOT NULL,
	 "college_id" text NOT NULL);

-- CREATE TABLE subjects
CREATE TABLE "subjects" (
	 "Branch_Id" text NOT NULL,
	 "Subject_Id" integer PRIMARY KEY,
	 "Subject_Name" text NOT NULL,
	 "added_by" text NOT NULL,
	 "college_id" text NOT NULL,
	 "course_id" text NOT NULL);

-- timestamptz Array Validator Function
CREATE FUNCTION validate_timestamptz_arr(
    arr timestamptz[] DEFAULT NULL,
    not_null boolean DEFAULT FALSE,
    min_arr_len integer DEFAULT NULL,
    max_arr_len integer DEFAULT NULL,
    
    min_ind timestamptz DEFAULT NULL, 
    max_ind timestamptz DEFAULT NULL,
    
    enum_arr timestamptz[] DEFAULT NULL)
RETURNS text AS $$
DECLARE
    val timestamptz;
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

-- text Array Validator Function
CREATE FUNCTION validate_text_arr(
    arr text[] DEFAULT NULL,
    not_null boolean DEFAULT FALSE,
    min_arr_len integer DEFAULT NULL,
    max_arr_len integer DEFAULT NULL,
    
    min_ind integer DEFAULT NULL, 
    max_ind integer DEFAULT NULL,
    
    enum_arr text[] DEFAULT NULL)
RETURNS text AS $$
DECLARE
    val text;
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
        IF min_ind IS NOT NULL AND LENGTH(val::text) < min_ind THEN
            RETURN FORMAT('Each element length should be at least %s', min_ind);
        END IF;

        IF max_ind IS NOT NULL AND LENGTH(val::text) > max_ind THEN
            RETURN FORMAT('Each element length should be at most %s', max_ind);
        END IF;

        IF enum_arr IS NOT NULL AND val::text NOT IN (SELECT * FROM unnest(enum_arr)) THEN
            RETURN FORMAT('%s element not present in enums', val);
        END IF;
    END LOOP;
    RETURN '';
END;
$$ LANGUAGE plpgsql;

-- date Array Validator Function
CREATE FUNCTION validate_date_arr(
    arr date[] DEFAULT NULL,
    not_null boolean DEFAULT FALSE,
    min_arr_len integer DEFAULT NULL,
    max_arr_len integer DEFAULT NULL,
    
    min_ind date DEFAULT NULL, 
    max_ind date DEFAULT NULL,
    
    enum_arr date[] DEFAULT NULL)
RETURNS text AS $$
DECLARE
    val date;
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

-- time Array Validator Function
CREATE FUNCTION validate_time_arr(
    arr time[] DEFAULT NULL,
    not_null boolean DEFAULT FALSE,
    min_arr_len integer DEFAULT NULL,
    max_arr_len integer DEFAULT NULL,
    
    min_ind time DEFAULT NULL, 
    max_ind time DEFAULT NULL,
    
    enum_arr time[] DEFAULT NULL)
RETURNS text AS $$
DECLARE
    val time;
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

res := validate_date_arr(NEW."Date_arr", false, NULL, NULL, NULL, NULL, array['2024-07-01', '2024-02-01']::date[]);
IF res != '' THEN
RAISE EXCEPTION 'Error in "Date_arr" column in "TypeTest" table: %', res;
END IF;

res := validate_timestamptz_arr(NEW."Datetime_Arr", false, NULL, NULL, NULL, NULL, array['2024-07-01T12:30:00+05:30', '2024-02-01T12:30:00+05:30']::timestamptz[]);
IF res != '' THEN
RAISE EXCEPTION 'Error in "Datetime_Arr" column in "TypeTest" table: %', res;
END IF;

res := validate_real_arr(NEW."Float_arr", false, NULL, NULL, NULL, NULL, array[4.5, 4.61]::real[]);
IF res != '' THEN
RAISE EXCEPTION 'Error in "Float_arr" column in "TypeTest" table: %', res;
END IF;

res := validate_integer_arr(NEW."Int_Arr", false, NULL, NULL, NULL, NULL, array[1, 2, 3]::integer[]);
IF res != '' THEN
RAISE EXCEPTION 'Error in "Int_Arr" column in "TypeTest" table: %', res;
END IF;

res := validate_text_arr(NEW."Str_Arr", false, NULL, NULL, NULL, NULL, array['val', 'val2']::text[]);
IF res != '' THEN
RAISE EXCEPTION 'Error in "Str_Arr" column in "TypeTest" table: %', res;
END IF;

res := validate_time_arr(NEW."Time_Arr", false, NULL, NULL, NULL, NULL, array['14:30:00']::time[]);
IF res != '' THEN
RAISE EXCEPTION 'Error in "Time_Arr" column in "TypeTest" table: %', res;
END IF;

RETURN NEW;

END;
$$ LANGUAGE plpgsql;

-- TypeTest Table Trigger
CREATE TRIGGER validate_table_TypeTest_trigger
BEFORE INSERT OR UPDATE ON "TypeTest"
FOR EACH ROW
EXECUTE FUNCTION validate_TypeTest_trigger();

-- DATA INSERTION "college"
INSERT INTO "college" ("college_id", "college_name", "principal_id")
VALUES
('college_1', 'IIT Delhi', 'jethalal');

-- DATA INSERTION "subjects"
INSERT INTO "subjects" ("Subject_Id", "Subject_Name", "college_id", "course_id", "Branch_Id", "added_by")
VALUES
(1, 'DS', 'college_1', 'course_1', 'branch_1', 'cs_hod'),
(2, 'COA', 'college_1', 'course_1', 'branch_1', 'jethalal'),
(3, 'WT', 'college_1', 'course_1', 'branch_2', 'it_hod'),
(4, 'Java', 'college_1', 'course_1', 'branch_2', 'jethalal'),
(5, 'RCC', 'college_1', 'course_1', 'branch_3', 'civil_hod');

-- DATA INSERTION "courses"
INSERT INTO "courses" ("Course_Id", "college_id", "Course_Name", "Lateral_Allowed", "added_by")
VALUES
('course_1', 'college_1', 'B. Tech.', true, 'jethalal'),
('course_2', 'college_1', 'M. Tech.', false, 'jethalal');

-- DATA INSERTION "students"
INSERT INTO "students" ("Student_Id", "Student_Name", "Student_Father", "college_id", "Course_Id", "Branch_Id", "added_by")
VALUES
(1, 'Tushar', 'Ajay', 'college_1', 'course_1', 'branch_1', 'jethalal'),
(2, 'Akshay', 'Nand', 'college_1', 'course_1', 'branch_1', 'cs_hod'),
(3, 'Saurabh', 'Jagganath', 'college_1', 'course_1', 'branch_2', 'it_hod'),
(4, 'Harsh', 'Ramesh', 'college_1', 'course_1', 'branch_3', 'civil_hod');

-- DATA INSERTION "branches"
INSERT INTO "branches" ("Branch_Id", "Branch_Name", "college_id", "Course_Id", "Teachers", "added_by", "HoD")
VALUES
('branch_1', 'Computer Science', 'college_1', 'course_1', array['HA', 'PC']::text[], 'jethalal', 'cs_hod'),
('branch_2', 'Information Technology', 'college_1', 'course_1', array['LD', 'RK']::text[], 'jethalal', 'it_hod'),
('branch_3', 'Civil Engineering', 'college_1', 'course_1', NULL, 'jethalal', 'civil_hod');

-- DATA INSERTION "TypeTest"
INSERT INTO "TypeTest" ("Int", "String", "Float", "Date", "Time", "DateTime", "Bool", "Int_Arr", "Str_Arr", "Float_arr", "Date_arr", "Time_Arr", "Datetime_Arr", "Bool_Arr")
VALUES
(1, 'Ram', 4.1, '2024-01-01', '14:30:00', NULL, false, array[1, 2]::integer[], array['val', 'val2']::text[], array[4.5, 4.61]::real[], array['2024-07-01', '2024-07-01']::date[], array['14:30:00', '14:30:00']::time[], array['2024-07-01T12:30:00+05:30', '2024-07-01T12:30:00+05:30']::timestamptz[], array[true, false]::boolean[]),
(1, 'Ram', 4.1, '2024-01-01', '14:30:00', '2024-07-01T12:30:00+05:30', true, array[1, 2]::integer[], array['val', 'val2']::text[], array[4.5, 4.61]::real[], array['2024-07-01', '2024-07-01']::date[], array['14:30:00', '14:30:00']::time[], array['2024-07-01T12:30:00+05:30', '2024-07-01T12:30:00+05:30']::timestamptz[], array[true, false]::boolean[]),
(1, 'Ram', 4.1, '2024-01-01', '14:30:00', NULL, true, array[1, 2]::integer[], NULL, array[4.5, 4.61]::real[], array['2024-07-01', '2024-07-01']::date[], array['14:30:00', '14:30:00']::time[], array['2024-07-01T12:30:00+05:30', '2024-07-01T12:30:00+05:30']::timestamptz[], array[true, false]::boolean[]),
(1, NULL, 4.1, '2024-01-01', '14:30:00', NULL, true, array[1, 2]::integer[], array[]::text[], array[4.5, 4.61]::real[], array['2024-07-01', '2024-07-01']::date[], array['14:30:00', '14:30:00']::time[], array['2024-07-01T12:30:00+05:30', '2024-07-01T12:30:00+05:30']::timestamptz[], array[true, false]::boolean[]);

-- DATA INSERTION "login"
INSERT INTO "login" ("username", "password", "role", "college_id", "course_id", "branch_id", "added_by")
VALUES
('superuser', '$2a$10$AP/OqoTkIcW2Ku8PGyQRF.X797PY1P5rbxgfs7CviNcfAXrb3Qryy', 'admin', NULL, NULL, NULL, NULL),
('jethalal', '$2a$10$mkRVU8LrXlwE49/XrikKAuX6UsQkJsBnBc3jcws.btiCT5x9z/eoO', 'principal', 'college_1', NULL, NULL, 'superuser'),
('cs_hod', '$2a$10$3KYwMX5c36pER8RL.L317.pc3D.eorcpbprhNB2NNuCSAiXuMCp7y', 'hod', 'college_1', 'course_1', 'branch_1', 'jethalal'),
('it_hod', '$2a$10$cOW7cnvwnlGDcQ3MDS2OXe0Gj5xHCJlql12yvlYmX/z6taLSfDP.S', 'hod', 'college_1', 'course_1', 'branch_2', 'jethalal'),
('civil_hod', '$2a$10$xiF3YR7zKZMwZmFHQFXX4.3q3hu9QGFkeC4L3xp937TkX/mQOg3Dq', 'hod', 'college_1', 'course_1', 'branch_3', 'jethalal');

-- branches Table Foreign Keys
ALTER TABLE "branches"
ADD CONSTRAINT "branches_Course_Id_fkey" FOREIGN KEY ("Course_Id")
REFERENCES "courses" ("Course_Id")
ON UPDATE CASCADE
ON DELETE CASCADE,
ADD CONSTRAINT "branches_HoD_fkey" FOREIGN KEY ("HoD")
REFERENCES "login" ("username")
ON UPDATE CASCADE
ON DELETE CASCADE,
ADD CONSTRAINT "branches_added_by_fkey" FOREIGN KEY ("added_by")
REFERENCES "login" ("username")
ON UPDATE CASCADE
ON DELETE CASCADE,
ADD CONSTRAINT "branches_college_id_fkey" FOREIGN KEY ("college_id")
REFERENCES "college" ("college_id")
ON UPDATE CASCADE
ON DELETE CASCADE;

-- college Table Foreign Keys
ALTER TABLE "college"
ADD CONSTRAINT "college_principal_id_fkey" FOREIGN KEY ("principal_id")
REFERENCES "login" ("username")
ON UPDATE CASCADE
ON DELETE CASCADE;

-- courses Table Foreign Keys
ALTER TABLE "courses"
ADD CONSTRAINT "courses_added_by_fkey" FOREIGN KEY ("added_by")
REFERENCES "login" ("username")
ON UPDATE CASCADE
ON DELETE CASCADE,
ADD CONSTRAINT "courses_college_id_fkey" FOREIGN KEY ("college_id")
REFERENCES "college" ("college_id")
ON UPDATE CASCADE
ON DELETE CASCADE;

-- login Table Foreign Keys
ALTER TABLE "login"
ADD CONSTRAINT "login_added_by_fkey" FOREIGN KEY ("added_by")
REFERENCES "login" ("username")
ON UPDATE CASCADE
ON DELETE CASCADE;

-- students Table Foreign Keys
ALTER TABLE "students"
ADD CONSTRAINT "students_Branch_Id_fkey" FOREIGN KEY ("Branch_Id")
REFERENCES "branches" ("Branch_Id")
ON UPDATE CASCADE
ON DELETE CASCADE,
ADD CONSTRAINT "students_Course_Id_fkey" FOREIGN KEY ("Course_Id")
REFERENCES "courses" ("Course_Id")
ON UPDATE CASCADE
ON DELETE CASCADE,
ADD CONSTRAINT "students_added_by_fkey" FOREIGN KEY ("added_by")
REFERENCES "login" ("username")
ON UPDATE CASCADE
ON DELETE CASCADE,
ADD CONSTRAINT "students_college_id_fkey" FOREIGN KEY ("college_id")
REFERENCES "college" ("college_id")
ON UPDATE CASCADE
ON DELETE CASCADE;

-- subjects Table Foreign Keys
ALTER TABLE "subjects"
ADD CONSTRAINT "subjects_Branch_Id_fkey" FOREIGN KEY ("Branch_Id")
REFERENCES "branches" ("Branch_Id")
ON UPDATE CASCADE
ON DELETE CASCADE,
ADD CONSTRAINT "subjects_added_by_fkey" FOREIGN KEY ("added_by")
REFERENCES "login" ("username")
ON UPDATE CASCADE
ON DELETE CASCADE,
ADD CONSTRAINT "subjects_college_id_fkey" FOREIGN KEY ("college_id")
REFERENCES "college" ("college_id")
ON UPDATE CASCADE
ON DELETE CASCADE,
ADD CONSTRAINT "subjects_course_id_fkey" FOREIGN KEY ("course_id")
REFERENCES "courses" ("Course_Id")
ON UPDATE CASCADE
ON DELETE CASCADE;

