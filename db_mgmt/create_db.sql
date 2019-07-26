/* create klog business model db */
DROP DATABASE IF EXISTS klog_business;
CREATE DATABASE klog_business;
USE klog_business;


/* table for institute info */
CREATE TABLE institute (
	pid INT NOT NULL AUTO_INCREMENT,
    institute_uid VARCHAR(63) NOT NULL,
    institute_name VARCHAR(63) NOT NULL,
    address VARCHAR(255) NOT NULL,
    country_code VARCHAR(3) NOT NULL,
    create_ts TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    modify_ts TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (pid),
    UNIQUE (institute_uid)
) ENGINE InnoDB;
ALTER TABLE institute AUTO_INCREMENT=1001;

/* table for class info */
CREATE TABLE class (
	pid INT NOT NULL AUTO_INCREMENT,
	class_uid VARCHAR(63) NOT NULL,
    class_name VARCHAR(63) NOT NULL,
    location VARCHAR(255) NOT NULL,
    institute_pid INT,
    create_ts TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    modify_ts TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (pid),
    UNIQUE (class_uid),
    FOREIGN KEY (institute_pid) REFERENCES institute(pid)
    ON UPDATE CASCADE
    ON DELETE SET NULL
) ENGINE InnoDB;
ALTER TABLE class AUTO_INCREMENT=1001;

/* table for teacher info */
CREATE TABLE teacher (
	pid INT NOT NULL AUTO_INCREMENT,
	teacher_uid VARCHAR(63) NOT NULL,
    first_name VARCHAR(63) NOT NULL,
    last_name VARCHAR(63) NOT NULL,
    date_of_birth DATE NOT NULL,
    address VARCHAR(255) NOT NULL,
    phone_number VARCHAR(31) NOT NULL,
    email VARCHAR(127) NOT NULL,
    institute_pid INT,
    create_ts TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    modify_ts TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (pid),
    UNIQUE (teacher_uid),
    FOREIGN KEY (institute_pid) REFERENCES institute(pid)
    ON UPDATE CASCADE
    ON DELETE SET NULL
) ENGINE InnoDB;
ALTER TABLE teacher AUTO_INCREMENT=1001;

/* table for class-teacher relationship info */
CREATE TABLE class_teacher (
	pid INT NOT NULL AUTO_INCREMENT,
    class_pid INT NOT NULL,
    teacher_pid INT NOT NULL,
    PRIMARY KEY (pid),
    UNIQUE (class_pid, teacher_pid),
    FOREIGN KEY (class_pid) REFERENCES class(pid)
    ON UPDATE CASCADE
    ON DELETE CASCADE,
    FOREIGN KEY (teacher_pid) REFERENCES teacher(pid)
    ON UPDATE CASCADE
    ON DELETE CASCADE
) ENGINE InnoDB;
ALTER TABLE class_teacher AUTO_INCREMENT=1001;

/* table for student info */
CREATE TABLE student (
	pid INT NOT NULL AUTO_INCREMENT,
	student_uid VARCHAR(63) NOT NULL,
    first_name VARCHAR(63) NOT NULL,
    last_name VARCHAR(63) NOT NULL,
    date_of_birth DATE NOT NULL,
    media_location VARCHAR(255) NOT NULL,
    class_pid INT,
    create_ts TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    modify_ts TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (pid),
    UNIQUE (student_uid),
    FOREIGN KEY (class_pid) REFERENCES class(pid)
    ON UPDATE CASCADE
    ON DELETE SET NULL
) ENGINE InnoDB;
ALTER TABLE student AUTO_INCREMENT=1001;

/* table for parent info */
CREATE TABLE parent (
	pid INT NOT NULL AUTO_INCREMENT,
	parent_uid VARCHAR(63) NOT NULL,
    first_name VARCHAR(63) NOT NULL,
    last_name VARCHAR(63) NOT NULL,
    date_of_birth DATE NOT NULL,
	address VARCHAR(255) NOT NULL,
	phone_number VARCHAR(31) NOT NULL,
    email VARCHAR(127) NOT NULL,
    occupation VARCHAR(31) NOT NULL,
    create_ts TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    modify_ts TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (pid),
    UNIQUE (parent_uid)
) ENGINE InnoDB;
ALTER TABLE parent AUTO_INCREMENT=1001;

/* table for student-parent relationship info */
CREATE TABLE student_parent (
	pid INT NOT NULL AUTO_INCREMENT,
    student_pid INT NOT NULL,
    parent_pid INT NOT NULL,
    PRIMARY KEY (pid),
    UNIQUE (student_pid, parent_pid),
    FOREIGN KEY (student_pid) REFERENCES student(pid)
    ON UPDATE CASCADE
    ON DELETE CASCADE,
    FOREIGN KEY (parent_pid) REFERENCES parent(pid)
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
(pid,teacher_uid,first_name,last_name,@date_of_birth,address,phone_number,email,institute_pid,create_ts,modify_ts)
SET date_of_birth = STR_TO_DATE(@date_of_birth, '%m/%d/%Y');

LOAD DATA LOCAL INFILE './samples/class_teacher.csv' INTO TABLE class_teacher
FIELDS ENCLOSED BY '"' TERMINATED BY ',' LINES TERMINATED BY '\n'
IGNORE 1 ROWS;

LOAD DATA LOCAL INFILE './samples/student.csv' INTO TABLE student
FIELDS ENCLOSED BY '"' TERMINATED BY ',' LINES TERMINATED BY '\n'
IGNORE 1 ROWS
(pid,student_uid,first_name,last_name,@date_of_birth,media_location,class_pid,create_ts,modify_ts)
SET date_of_birth = STR_TO_DATE(@date_of_birth, '%m/%d/%Y');

LOAD DATA LOCAL INFILE './samples/parent.csv' INTO TABLE parent
FIELDS ENCLOSED BY '"' TERMINATED BY ',' LINES TERMINATED BY '\n'
IGNORE 1 ROWS
(pid,parent_uid,first_name,last_name,@date_of_birth,address,phone_number,email,occupation,create_ts,modify_ts)
SET date_of_birth = STR_TO_DATE(@date_of_birth, '%m/%d/%Y');

LOAD DATA LOCAL INFILE './samples/student_parent.csv' INTO TABLE student_parent
FIELDS ENCLOSED BY '"' TERMINATED BY ',' LINES TERMINATED BY '\n'
IGNORE 1 ROWS;