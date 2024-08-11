-- CREATE TABLE TypeTest
CREATE TABLE "TypeTest" (
	 "__ID" SERIAL PRIMARY KEY, 
	 "Bool" boolean,
	 "Bool_Arr" boolean[],
	 "Date" date,
	 "DateTime" timestamptz,
	 "Date_arr" date[],
	 "Datetime_Arr" timestamptz[],
	 "Float" real,
	 "Float_arr" real[],
	 "Int" integer,
	 "Int_Arr" integer[],
	 "Str_Arr" text[],
	 "String" text,
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

-- DATA INSERTION "branches"
INSERT INTO "branches" ("Branch_Id", "Branch_Name", "Course_Id", "Teachers")
VALUES
(1, 'Computer Science', 1, array['HA', 'PC']::text[]),
(2, 'Information Technology', 1, array['LD', 'RK']::text[]),
(3, 'Civil Engineering', 1, NULL);

-- DATA INSERTION "students"
INSERT INTO "students" ("Student_Id", "Student_Name", "Student_Father", "Course_Id", "Branch_Id")
VALUES
(1, 'Tushar', 'Ajay', 1, 1),
(2, 'Akshay', 'Nand', 1, 1),
(3, 'Saurabh', 'Jagganath', 1, 2),
(4, 'Harsh', 'Ramesh', 1, 2);

-- DATA INSERTION "subjects"
INSERT INTO "subjects" ("Subject_Id", "Subject_Name", "Branch_Id")
VALUES
(1, 'DS', 1),
(2, 'COA', 1),
(3, 'WT', 2),
(4, 'Java', 2);

-- DATA INSERTION "courses"
INSERT INTO "courses" ("Course_Id", "Course_Name", "Lateral_Allowed")
VALUES
(1, 'B. Tech.', true),
(2, 'M. Tech.', false);

-- DATA INSERTION "TypeTest"
INSERT INTO "TypeTest" ("Int", "String", "Float", "Date", "Time", "DateTime", "Bool", "Int_Arr", "Str_Arr", "Float_arr", "Date_arr", "Time_Arr", "Datetime_Arr", "Bool_Arr")
VALUES
(1, 'Ram', 4.1, '2024-01-01', '14:30:00', NULL, false, array[1, 2]::integer[], array['val', 'val2']::text[], array[4.5, 4.61]::real[], array['2024-07-01', '2024-07-01']::date[], array['14:30:00', '14:30:00']::time[], array['2024-07-01T12:30:00+05:30', '2024-07-01T12:30:00+05:30']::timestamptz[], array[true, false]::boolean[]),
(1, 'Ram', 4.1, '2024-01-01', '14:30:00', '2024-07-01T12:30:00+05:30', true, array[1, 2]::integer[], array['val', 'val2']::text[], array[4.5, 4.61]::real[], array['2024-07-01', '2024-07-01']::date[], array['14:30:00', '14:30:00']::time[], array['2024-07-01T12:30:00+05:30', '2024-07-01T12:30:00+05:30']::timestamptz[], array[true, false]::boolean[]),
(1, 'Ram', 4.1, '2024-01-01', '14:30:00', NULL, true, array[1, 2]::integer[], array['val', 'val2']::text[], array[4.5, 4.61]::real[], array['2024-07-01', '2024-07-01']::date[], array['14:30:00', '14:30:00']::time[], array['2024-07-01T12:30:00+05:30', '2024-07-01T12:30:00+05:30']::timestamptz[], array[true, false]::boolean[]);

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

