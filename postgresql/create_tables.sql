CREATE DATABASE TimeTableDB;
ALTER ROLE denis SET client_encoding TO 'utf8';
GRANT ALL PRIVILEGES ON DATABASE TimeTableDB TO denis;

CREATE SEQUENCE lecturer_id_lecturer_seq START WITH 1;

CREATE TABLE IF NOT EXISTS lecturer
(
    id_lecturer integer NOT NULL DEFAULT nextval('lecturer_id_lecturer_seq'),
    full_name character(100) NOT NULL,
    CONSTRAINT lecturer_pkey PRIMARY KEY (id_lecturer)
);

CREATE TABLE IF NOT EXISTS call_schedule
(
    subject_number integer NOT NULL,
    time_start time without time zone NOT NULL,
    time_end time without time zone NOT NULL,
    CONSTRAINT call_schedule_pkey PRIMARY KEY (subject_number)
);

CREATE SEQUENCE id_of_course_seq START WITH 1;

CREATE TABLE IF NOT EXISTS course
(
    id_of_course integer NOT NULL DEFAULT nextval('id_of_course_seq'),
    CONSTRAINT course_pkey PRIMARY KEY (id_of_course)
);

CREATE SEQUENCE id_of_degree_seq START WITH 1;

CREATE TABLE IF NOT EXISTS degree
(
    id_of_degree integer NOT NULL DEFAULT nextval('id_of_degree_seq'),
    degree_of_study character(100) NOT NULL,
    CONSTRAINT degree_pkey PRIMARY KEY (id_of_degree)
);

CREATE SEQUENCE id_of_the_institute_seq START WITH 1;

CREATE TABLE IF NOT EXISTS institute
(
    id_of_the_institute integer NOT NULL DEFAULT nextval('id_of_the_institute_seq'),
    name_of_the_institute character(100) NOT NULL,
    CONSTRAINT institute_pkey PRIMARY KEY (id_of_the_institute)
);

CREATE SEQUENCE id_group_seq START WITH 1;

CREATE TABLE IF NOT EXISTS study_group
(
    id_group integer NOT NULL DEFAULT nextval('id_group_seq'),
    group_name character(20) NOT NULL,
    id_institute integer,
    id_of_course integer NOT NULL,
    id_degree integer,
    CONSTRAINT study_group_pkey PRIMARY KEY (id_group)
);

CREATE TABLE IF NOT EXISTS timetable
(
    id_to_group integer NOT NULL,
    subject_to_number integer NOT NULL,
    id_lectur integer,  
    subject_title character(500) NOT NULL,
    auditorium character(500),
    day_week integer NOT NULL,
    type_of_week integer NOT NULL,
    CONSTRAINT fkey PRIMARY KEY (id_to_group, subject_to_number, day_week, type_of_week)
);
