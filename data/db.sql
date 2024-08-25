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
	 "Branch_Id" text PRIMARY KEY,
	 "Branch_Name" text NOT NULL,
	 "Course_Id" text NOT NULL,
	 "HoD" text NOT NULL,
	 "Teachers" text[],
	 "added_by" text NOT NULL,
	 "college_id" text NOT NULL);

-- CREATE TABLE college
CREATE TABLE "college" (
	 "college_id" text PRIMARY KEY,
	 "college_name" text NOT NULL,
	 "principal_id" text NOT NULL);

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
	 "role" text NOT NULL,
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

-- DATA INSERTION "branches"
INSERT INTO "branches" ("Branch_Id", "Branch_Name", "college_id", "Course_Id", "Teachers", "added_by", "HoD")
VALUES
('branch_1', 'Computer Science', 'college_1', 'course_1', array['HA', 'PC']::text[], 'jethalal', 'cs_hod'),
('branch_2', 'Information Technology', 'college_1', 'course_1', array['LD', 'RK']::text[], 'jethalal', 'it_hod'),
('branch_3', 'Civil Engineering', 'college_1', 'course_1', NULL, 'jethalal', 'civil_hod');

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
('superuser', '$2a$10$dVDpdr50WzvBVyvP54fY/ukYboookNvft76FU7qPi.QFAD2qvOX.u', 'admin', NULL, NULL, NULL, NULL),
('jethalal', '$2a$10$evs4PTQj9PVltrbupLp0QuWzis6Vf6I8S72HW8TSbTyn63B1zwqlO', 'principal', 'college_1', NULL, NULL, 'superuser'),
('cs_hod', '$2a$10$JBDn86kFsKRu/hjtvh/s.u54lwn2SrG.enul97d0bCnv7e0k8Im7m', 'hod', 'college_1', 'course_1', 'branch_1', 'jethalal'),
('it_hod', '$2a$10$lBe/Bi4xQv5GtspJUmHJq.1VddXnhJdifZTlsP1wCsPKkof2phxfC', 'hod', 'college_1', 'course_1', 'branch_2', 'jethalal'),
('civil_hod', '$2a$10$QwTH0nky6Y6HnVrBinQyA.DkJWmpKKs4i0yDBCVYO6aZT3XRafFhG', 'hod', 'college_1', 'course_1', 'branch_3', 'jethalal');

-- branches Table Foreign Keys
ALTER TABLE "branches"
ADD CONSTRAINT "branches_Course_Id_fkey" FOREIGN KEY ("Course_Id")
REFERENCES "courses" ("Course_Id"),
ADD CONSTRAINT "branches_HoD_fkey" FOREIGN KEY ("HoD")
REFERENCES "login" ("username"),
ADD CONSTRAINT "branches_added_by_fkey" FOREIGN KEY ("added_by")
REFERENCES "login" ("username"),
ADD CONSTRAINT "branches_college_id_fkey" FOREIGN KEY ("college_id")
REFERENCES "college" ("college_id");

-- college Table Foreign Keys
ALTER TABLE "college"
ADD CONSTRAINT "college_principal_id_fkey" FOREIGN KEY ("principal_id")
REFERENCES "login" ("username");

-- courses Table Foreign Keys
ALTER TABLE "courses"
ADD CONSTRAINT "courses_added_by_fkey" FOREIGN KEY ("added_by")
REFERENCES "login" ("username"),
ADD CONSTRAINT "courses_college_id_fkey" FOREIGN KEY ("college_id")
REFERENCES "college" ("college_id");

-- login Table Foreign Keys
ALTER TABLE "login"
ADD CONSTRAINT "login_added_by_fkey" FOREIGN KEY ("added_by")
REFERENCES "login" ("username");

-- students Table Foreign Keys
ALTER TABLE "students"
ADD CONSTRAINT "students_Branch_Id_fkey" FOREIGN KEY ("Branch_Id")
REFERENCES "branches" ("Branch_Id"),
ADD CONSTRAINT "students_Course_Id_fkey" FOREIGN KEY ("Course_Id")
REFERENCES "courses" ("Course_Id"),
ADD CONSTRAINT "students_added_by_fkey" FOREIGN KEY ("added_by")
REFERENCES "login" ("username"),
ADD CONSTRAINT "students_college_id_fkey" FOREIGN KEY ("college_id")
REFERENCES "college" ("college_id");

-- subjects Table Foreign Keys
ALTER TABLE "subjects"
ADD CONSTRAINT "subjects_Branch_Id_fkey" FOREIGN KEY ("Branch_Id")
REFERENCES "branches" ("Branch_Id"),
ADD CONSTRAINT "subjects_added_by_fkey" FOREIGN KEY ("added_by")
REFERENCES "login" ("username"),
ADD CONSTRAINT "subjects_college_id_fkey" FOREIGN KEY ("college_id")
REFERENCES "college" ("college_id"),
ADD CONSTRAINT "subjects_course_id_fkey" FOREIGN KEY ("course_id")
REFERENCES "courses" ("Course_Id");

