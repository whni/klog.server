package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"net/http"
	"strconv"
)

var studentHandlerTable = map[string]gin.HandlerFunc{
	"get":    studentGetHandler,
	"post":   studentPostHandler,
	"put":    studentPutHandler,
	"delete": studentDeleteHandler,
}

func studentGetHandler(ctx *gin.Context) {
	params := ginContextRequestParameter(ctx)
	response := GinResponse{
		Status: http.StatusOK,
	}
	defer func() {
		ginContextProcessResponse(ctx, &response)
	}()

	var students []*Student
	var err error
	var pid int
	pid, err = strconv.Atoi(params.PID)
	if err != nil || pid < 0 {
		response.Status = http.StatusBadRequest
		response.Message = fmt.Sprintf("[%s] - Please specifiy a valid PID (pid >= 0)", serverErrorMessages[seInputParamNotValid])
		return
	}

	// pid: 0 for all, > 0 for specified one
	students, err = findStudent(pid)
	if err != nil {
		response.Status = http.StatusConflict
		response.Message = err.Error()
		return
	}
	response.Payload = students
	return
}

func studentPostHandler(ctx *gin.Context) {
	params := ginContextRequestParameter(ctx)
	response := GinResponse{
		Status: http.StatusOK,
	}
	defer func() {
		ginContextProcessResponse(ctx, &response)
	}()

	var student Student
	var studentPID int
	var err error

	if err = json.Unmarshal(params.Data, &student); err != nil {
		response.Status = http.StatusBadRequest
		response.Message = fmt.Sprintf("[%s] - %s", serverErrorMessages[seInputJSONNotValid], err.Error())
		return
	}

	studentPID, err = createStudent(&student)
	if err != nil {
		response.Status = http.StatusConflict
		response.Message = err.Error()
	} else {
		response.Payload = studentPID
	}
	return
}

func studentPutHandler(ctx *gin.Context) {
	params := ginContextRequestParameter(ctx)
	response := GinResponse{
		Status: http.StatusOK,
	}
	defer func() {
		ginContextProcessResponse(ctx, &response)
	}()

	var student Student
	var err error

	if err = json.Unmarshal(params.Data, &student); err != nil {
		response.Status = http.StatusBadRequest
		response.Message = fmt.Sprintf("[%s] - %s", serverErrorMessages[seInputJSONNotValid], err.Error())
		return
	}

	err = updateStudent(&student)
	if err != nil {
		response.Status = http.StatusConflict
		response.Message = err.Error()
	} else {
		response.Payload = student.PID
	}
	return
}

func studentDeleteHandler(ctx *gin.Context) {
	params := ginContextRequestParameter(ctx)
	response := GinResponse{
		Status: http.StatusOK,
	}
	defer func() {
		ginContextProcessResponse(ctx, &response)
	}()

	var err error
	var deletedRows int
	var pid int
	pid, err = strconv.Atoi(params.PID)
	if err != nil || pid < 0 {
		response.Status = http.StatusBadRequest
		response.Message = fmt.Sprintf("[%s] - Please specifiy a valid PID (pid >= 0)", serverErrorMessages[seInputParamNotValid])
		return
	}

	// pid: 0 for all, > 0 for specified one
	deletedRows, err = deleteStudent(pid)
	if err != nil {
		response.Status = http.StatusConflict
		response.Message = err.Error()
		return
	}
	response.Payload = deletedRows
	return
}

// find student, return student ptr, error
func findStudent(pid int) ([]*Student, error) {
	var rows *sql.Rows
	var err error
	defer func() {
		if err != nil {
			logging.Errormf(logModStudentHandler, err.Error())
		}
	}()

	var dbQuery = "SELECT pid, student_uid, first_name, last_name, date_of_birth, media_location, class_pid, create_ts, modify_ts FROM student"
	if pid == 0 {
		rows, err = dbPool.Query(dbQuery)
	} else if pid > 0 {
		rows, err = dbPool.Query(dbQuery+" WHERE pid = ?", pid)
	} else {
		err = fmt.Errorf("[%s] - PID (%d) not valid", serverErrorMessages[seInputParamNotValid], pid)
		return nil, err
	}
	if err != nil {
		err = fmt.Errorf("[%s] - %s", serverErrorMessages[seDBResourceQuery], err.Error())
		return nil, err
	}
	defer rows.Close()

	students := []*Student{}
	for rows.Next() {
		var student Student
		err = rows.Scan(&student.PID, &student.StudentUID, &student.FirstName, &student.LastName, &student.DateOfBirth,
			&student.MediaLocation, &student.ClassPID, &student.CreateTS, &student.ModifyTS)
		if err != nil {
			err = fmt.Errorf("[%s] - %s", serverErrorMessages[seDBResourceQuery], err.Error())
			return nil, err
		}
		students = append(students, &student)
	}

	err = rows.Err()
	if err != nil {
		err = fmt.Errorf("[%s] - %s", serverErrorMessages[seDBResourceQuery], err.Error())
		return nil, err
	}

	logging.Debugmf(logModStudentHandler, "Found %d student results from DB (PID=%d)", len(students), pid)
	return students, nil
}

// create student, return PID, error
func createStudent(student *Student) (int, error) {
	var err error
	var result sql.Result

	defer func() {
		if err != nil {
			logging.Errormf(logModStudentHandler, err.Error())
		}
	}()

	if err = ginStructValidCheck(student); err != nil {
		return 0, err
	}

	if student.PID > 0 {
		if students, errExist := findStudent(student.PID); errExist == nil && len(students) > 0 {
			err = fmt.Errorf("[%s] - Student (PID=%d) already exists", serverErrorMessages[seResourceDuplicated], student.PID)
			return 0, err
		}
	} else if student.PID < 0 {
		err = fmt.Errorf("[%s] - PID (%d) not valid", serverErrorMessages[seInputParamNotValid], student.PID)
		return 0, err
	}

	var dbQuery string
	if student.PID > 0 {
		dbQuery = "INSERT INTO student(pid, student_uid, first_name, last_name, date_of_birth, media_location, class_pid) VALUES (?, ?, ?, ?, ?, ?, ?)"
		result, err = dbPool.Exec(dbQuery, student.PID, student.StudentUID, student.FirstName, student.LastName, student.DateOfBirth, student.MediaLocation, student.ClassPID)
	} else {
		dbQuery = "INSERT INTO student(student_uid, first_name, last_name, date_of_birth, media_location, class_pid) VALUES (?, ?, ?, ?, ?, ?)"
		result, err = dbPool.Exec(dbQuery, student.StudentUID, student.FirstName, student.LastName, student.DateOfBirth, student.MediaLocation, student.ClassPID)
	}
	if err != nil {
		err = fmt.Errorf("[%s] - %s", serverErrorMessages[seDBResourceQuery], err.Error())
		return 0, err
	}

	lastInsertID, err := result.LastInsertId()
	if err != nil {
		lastInsertID = -1
		logging.Warnmf(logModStudentHandler, "Cound not retrieve created student PID")
	}
	logging.Debugmf(logModStudentHandler, "Created student in DB (LastInsertId,PID=%v)", lastInsertID)
	return int(lastInsertID), nil
}

// update student, return error
func updateStudent(student *Student) error {
	var err error
	defer func() {
		if err != nil {
			logging.Errormf(logModStudentHandler, err.Error())
		}
	}()

	if err = ginStructValidCheck(student); err != nil {
		return err
	}

	if student.PID <= 0 {
		err = fmt.Errorf("[%s] - PID (%d) not valid", serverErrorMessages[seInputParamNotValid], student.PID)
		return err
	}
	if students, errExist := findStudent(student.PID); errExist != nil || len(students) == 0 {
		err = fmt.Errorf("[%s] - Student (PID=%d) does not exist ==> not updated", serverErrorMessages[seResourceNotFound], student.PID)
		return err
	}

	var dbQuery = "UPDATE student SET student_uid=?, first_name=?, last_name=?, date_of_birth=?, media_location=?, class_pid=? WHERE pid=?"
	_, err = dbPool.Exec(dbQuery, student.StudentUID, student.FirstName, student.LastName, student.DateOfBirth, student.MediaLocation, student.ClassPID, student.PID)
	if err != nil {
		err = fmt.Errorf("[%s] - %s", serverErrorMessages[seDBResourceQuery], err.Error())
		return err
	}

	logging.Debugmf(logModStudentHandler, "Updated student in DB (PID=%v)", student.PID)
	return nil
}

// delete student, return #delete rows, error
func deleteStudent(pid int) (int, error) {
	var err error
	var result sql.Result

	defer func() {
		if err != nil {
			logging.Errormf(logModStudentHandler, err.Error())
		}
	}()

	var dbQuery = "DELETE FROM student"
	if pid == 0 {
		result, err = dbPool.Exec(dbQuery)
	} else if pid > 0 {
		result, err = dbPool.Exec(dbQuery+" WHERE pid = ?", pid)
	} else {
		err = fmt.Errorf("[%s] - PID (%d) not valid", serverErrorMessages[seInputParamNotValid], pid)
		return 0, err
	}
	if err != nil {
		err = fmt.Errorf("[%s] - %s", serverErrorMessages[seDBResourceQuery], err.Error())
		return 0, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		rowsAffected = -1
		logging.Warnmf(logModStudentHandler, "Cound not count #deleted students")
	}
	logging.Debugmf(logModStudentHandler, "Deleted %d student results from DB", rowsAffected)
	return int(rowsAffected), nil
}
