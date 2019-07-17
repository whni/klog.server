/* create vlog db */
DROP DATABASE IF EXISTS vlog;
CREATE DATABASE vlog;
USE vlog;


/* table for institute info */
CREATE TABLE institute  (
	db_id INT NOT NULL AUTO_INCREMENT,
	institute_id VARCHAR(63) NOT NULL UNIQUE,
    institute_name VARCHAR(63) NOT NULL UNIQUE,
    address VARCHAR(255),
    country_code VARCHAR(3),
    PRIMARY KEY (db_id)
) ENGINE InnoDB;
ALTER TABLE institute AUTO_INCREMENT=1001;

/* table for class info */
CREATE TABLE class (
	db_id INT NOT NULL AUTO_INCREMENT,
	class_id VARCHAR(63) NOT NULL UNIQUE,
    class_name VARCHAR(63) NOT NULL,
    location VARCHAR(255),
    institute_id VARCHAR(63),
    PRIMARY KEY (db_id),
    FOREIGN KEY (institute_id) REFERENCES institute(institute_id)
    ON UPDATE CASCADE
    ON DELETE SET NULL
) ENGINE InnoDB;
ALTER TABLE class AUTO_INCREMENT=1001;

/* table for teacher info */
CREATE TABLE teacher (
	db_id INT NOT NULL AUTO_INCREMENT,
	teacher_id VARCHAR(63) NOT NULL UNIQUE,
    first_name VARCHAR(63) NOT NULL,
    last_name VARCHAR(63) NOT NULL,
    date_of_birth DATE,
    address VARCHAR(255),
    phone_number VARCHAR(31),
    email VARCHAR(127),
    PRIMARY KEY (db_id)
) ENGINE InnoDB;
ALTER TABLE teacher AUTO_INCREMENT=1001;

/* table for class-teacher relationship info */
CREATE TABLE class_teacher (
	db_id INT NOT NULL AUTO_INCREMENT,
    class_id VARCHAR(63) NOT NULL,
    teacher_id VARCHAR(63) NOT NULL,
    PRIMARY KEY (db_id),
    UNIQUE (class_id,teacher_id),
    FOREIGN KEY (class_id) REFERENCES class(class_id)
    ON UPDATE CASCADE
    ON DELETE CASCADE,
    FOREIGN KEY (teacher_id) REFERENCES teacher(teacher_id)
    ON UPDATE CASCADE
    ON DELETE CASCADE
) ENGINE InnoDB;
ALTER TABLE class_teacher AUTO_INCREMENT=1001;

/* table for student info */
CREATE TABLE student (
	db_id INT NOT NULL AUTO_INCREMENT,
	student_id VARCHAR(63) NOT NULL UNIQUE,
    first_name VARCHAR(63) NOT NULL,
    last_name VARCHAR(63) NOT NULL,
    date_of_birth DATE NOT NULL,
    media_location VARCHAR(255),
    class_id VARCHAR(63),
    PRIMARY KEY (db_id),
    FOREIGN KEY (class_id) REFERENCES class(class_id)
    ON UPDATE CASCADE
    ON DELETE SET NULL
) ENGINE InnoDB;
ALTER TABLE student AUTO_INCREMENT=1001;

/* table for parent info */
CREATE TABLE parent (
	db_id INT NOT NULL AUTO_INCREMENT,
	parent_id VARCHAR(63) NOT NULL UNIQUE,
    first_name VARCHAR(63) NOT NULL,
    last_name VARCHAR(63) NOT NULL,
    date_of_birth DATE,
	address VARCHAR(255),
	phone_number VARCHAR(31),
    email VARCHAR(127),
    occupation VARCHAR(31),
    PRIMARY KEY (db_id)
) ENGINE InnoDB;
ALTER TABLE parent AUTO_INCREMENT=1001;

/* table for student-parent relationship info */
CREATE TABLE student_parent (
	db_id INT NOT NULL AUTO_INCREMENT,
    student_id VARCHAR(63) NOT NULL,
    parent_id VARCHAR(63) NOT NULL,
    PRIMARY KEY (db_id),
    UNIQUE (student_id, parent_id),
    FOREIGN KEY (student_id) REFERENCES student(student_id)
    ON UPDATE CASCADE
    ON DELETE CASCADE,
    FOREIGN KEY (parent_id) REFERENCES parent(parent_id)
    ON UPDATE CASCADE
    ON DELETE CASCADE
) ENGINE InnoDB;
ALTER TABLE student_parent AUTO_INCREMENT=1001;


/* insert sample date */
LOAD DATA LOCAL INFILE './samples/institute.csv' INTO TABLE institute
FIELDS ENCLOSED BY '"' TERMINATED BY ',' LINES TERMINATED BY '\n'
IGNORE 1 ROWS;

LOAD DATA LOCAL INFILE './samples/class.csv' INTO TABLE class
FIELDS ENCLOSED BY '"' TERMINATED BY ',' LINES TERMINATED BY '\n'
IGNORE 1 ROWS;

LOAD DATA LOCAL INFILE './samples/teacher.csv' INTO TABLE teacher
FIELDS ENCLOSED BY '"' TERMINATED BY ',' LINES TERMINATED BY '\n'
IGNORE 1 ROWS
(db_id,teacher_id,first_name,last_name,@date_of_birth,address,phone_number,email)
SET date_of_birth = STR_TO_DATE(@date_of_birth, '%m/%d/%Y');

LOAD DATA LOCAL INFILE './samples/class_teacher.csv' INTO TABLE class_teacher
FIELDS ENCLOSED BY '"' TERMINATED BY ',' LINES TERMINATED BY '\n'
IGNORE 1 ROWS;

LOAD DATA LOCAL INFILE './samples/student.csv' INTO TABLE student
FIELDS ENCLOSED BY '"' TERMINATED BY ',' LINES TERMINATED BY '\n'
IGNORE 1 ROWS
(db_id,student_id,first_name,last_name,@date_of_birth,media_location,class_id)
SET date_of_birth = STR_TO_DATE(@date_of_birth, '%m/%d/%Y');

LOAD DATA LOCAL INFILE './samples/parent.csv' INTO TABLE parent
FIELDS ENCLOSED BY '"' TERMINATED BY ',' LINES TERMINATED BY '\n'
IGNORE 1 ROWS
(db_id,parent_id,first_name,last_name,@date_of_birth,address,phone_number,email,occupation)
SET date_of_birth = STR_TO_DATE(@date_of_birth, '%m/%d/%Y');

LOAD DATA LOCAL INFILE './samples/student_parent.csv' INTO TABLE student_parent
FIELDS ENCLOSED BY '"' TERMINATED BY ',' LINES TERMINATED BY '\n'
IGNORE 1 ROWS;