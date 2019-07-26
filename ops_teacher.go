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

var teacherHandlerTable = map[string]gin.HandlerFunc{
	"get":    teacherGetHandler,
	"post":   teacherPostHandler,
	"put":    teacherPutHandler,
	"delete": teacherDeleteHandler,
}

func teacherGetHandler(ctx *gin.Context) {
	params := ginContextRequestParameter(ctx)
	response := GinResponse{
		Status: http.StatusOK,
	}
	defer func() {
		ginContextProcessResponse(ctx, &response)
	}()

	var teachers []*Teacher
	var err error
	var pid int
	pid, err = strconv.Atoi(params.PID)
	if err != nil || pid < 0 {
		response.Status = http.StatusBadRequest
		response.Message = fmt.Sprintf("[%s] - Please specifiy a valid PID (pid >= 0)", serverErrorMessages[seInputParamNotValid])
		return
	}

	// pid: 0 for all, > 0 for specified one
	teachers, err = findTeacher(pid)
	if err != nil {
		response.Status = http.StatusConflict
		response.Message = err.Error()
		return
	}
	response.Payload = teachers
	return
}

func teacherPostHandler(ctx *gin.Context) {
	params := ginContextRequestParameter(ctx)
	response := GinResponse{
		Status: http.StatusOK,
	}
	defer func() {
		ginContextProcessResponse(ctx, &response)
	}()

	var teacher Teacher
	var teacherPID int
	var err error

	if err = json.Unmarshal(params.Data, &teacher); err != nil {
		response.Status = http.StatusBadRequest
		response.Message = fmt.Sprintf("[%s] - %s", serverErrorMessages[seInputJSONNotValid], err.Error())
		return
	}

	teacherPID, err = createTeacher(&teacher)
	if err != nil {
		response.Status = http.StatusConflict
		response.Message = err.Error()
	} else {
		response.Payload = teacherPID
	}
	return
}

func teacherPutHandler(ctx *gin.Context) {
	params := ginContextRequestParameter(ctx)
	response := GinResponse{
		Status: http.StatusOK,
	}
	defer func() {
		ginContextProcessResponse(ctx, &response)
	}()

	var teacher Teacher
	var err error

	if err = json.Unmarshal(params.Data, &teacher); err != nil {
		response.Status = http.StatusBadRequest
		response.Message = fmt.Sprintf("[%s] - %s", serverErrorMessages[seInputJSONNotValid], err.Error())
		return
	}

	err = updateTeacher(&teacher)
	if err != nil {
		response.Status = http.StatusConflict
		response.Message = err.Error()
	} else {
		response.Payload = teacher.PID
	}
	return
}

func teacherDeleteHandler(ctx *gin.Context) {
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
	deletedRows, err = deleteTeacher(pid)
	if err != nil {
		response.Status = http.StatusConflict
		response.Message = err.Error()
		return
	}
	response.Payload = deletedRows
	return
}

// find teacher, return teacher ptr, error
func findTeacher(pid int) ([]*Teacher, error) {
	var rows *sql.Rows
	var err error
	defer func() {
		if err != nil {
			logging.Errormf(logModTeacherHandler, err.Error())
		}
	}()

	var dbQuery = "SELECT pid, teacher_uid, first_name, last_name, date_of_birth, address, phone_number, email, institute_pid, create_ts, modify_ts FROM teacher"
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

	teachers := []*Teacher{}
	for rows.Next() {
		var teacher Teacher
		err = rows.Scan(&teacher.PID, &teacher.TeacherUID, &teacher.FirstName, &teacher.LastName, &teacher.DateOfBirth,
			&teacher.Address, &teacher.PhoneNumber, &teacher.Email, &teacher.InstitutePID, &teacher.CreateTS, &teacher.ModifyTS)
		if err != nil {
			err = fmt.Errorf("[%s] - %s", serverErrorMessages[seDBResourceQuery], err.Error())
			return nil, err
		}
		teachers = append(teachers, &teacher)
	}

	err = rows.Err()
	if err != nil {
		err = fmt.Errorf("[%s] - %s", serverErrorMessages[seDBResourceQuery], err.Error())
		return nil, err
	}

	logging.Debugmf(logModTeacherHandler, "Found %d teacher results from DB (PID=%d)", len(teachers), pid)
	return teachers, nil
}

// create teacher, return PID, error
func createTeacher(teacher *Teacher) (int, error) {
	var err error
	var result sql.Result

	defer func() {
		if err != nil {
			logging.Errormf(logModTeacherHandler, err.Error())
		}
	}()

	if err = ginStructValidCheck(teacher); err != nil {
		return 0, err
	}

	if teacher.PID > 0 {
		if teachers, errExist := findTeacher(teacher.PID); errExist == nil && len(teachers) > 0 {
			err = fmt.Errorf("[%s] - Teacher (PID=%d) already exists", serverErrorMessages[seResourceDuplicated], teacher.PID)
			return 0, err
		}
	} else if teacher.PID < 0 {
		err = fmt.Errorf("[%s] - PID (%d) not valid", serverErrorMessages[seInputParamNotValid], teacher.PID)
		return 0, err
	}

	var dbQuery string
	if teacher.PID > 0 {
		dbQuery = "INSERT INTO teacher(pid, teacher_uid, first_name, last_name, date_of_birth, address, phone_number, email, institute_pid) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)"
		result, err = dbPool.Exec(dbQuery, teacher.PID, teacher.TeacherUID, teacher.FirstName, teacher.LastName, teacher.DateOfBirth, teacher.Address,
			teacher.PhoneNumber, teacher.Email, teacher.InstitutePID)
	} else {
		dbQuery = "INSERT INTO teacher(teacher_uid, first_name, last_name, date_of_birth, address, phone_number, email, institute_pid) VALUES (?, ?, ?, ?, ?, ?, ?, ?)"
		result, err = dbPool.Exec(dbQuery, teacher.TeacherUID, teacher.FirstName, teacher.LastName, teacher.DateOfBirth, teacher.Address,
			teacher.PhoneNumber, teacher.Email, teacher.InstitutePID)
	}
	if err != nil {
		err = fmt.Errorf("[%s] - %s", serverErrorMessages[seDBResourceQuery], err.Error())
		return 0, err
	}

	lastInsertID, err := result.LastInsertId()
	if err != nil {
		lastInsertID = -1
		logging.Warnmf(logModTeacherHandler, "Cound not retrieve created teacher PID")
	}
	logging.Debugmf(logModTeacherHandler, "Created teacher in DB (LastInsertId,PID=%v)", lastInsertID)
	return int(lastInsertID), nil
}

// update teacher, return error
func updateTeacher(teacher *Teacher) error {
	var err error
	defer func() {
		if err != nil {
			logging.Errormf(logModTeacherHandler, err.Error())
		}
	}()

	if err = ginStructValidCheck(teacher); err != nil {
		return err
	}

	if teacher.PID <= 0 {
		err = fmt.Errorf("[%s] - PID (%d) not valid", serverErrorMessages[seInputParamNotValid], teacher.PID)
		return err
	}
	if teachers, errExist := findTeacher(teacher.PID); errExist != nil || len(teachers) == 0 {
		err = fmt.Errorf("[%s] - Teacher (PID=%d) does not exist ==> not updated", serverErrorMessages[seResourceNotFound], teacher.PID)
		return err
	}

	var dbQuery = "UPDATE teacher SET teacher_uid=?, first_name=?, last_name=?, date_of_birth=?, address=?, phone_number=?, email=?, institute_pid=? WHERE pid=?"
	_, err = dbPool.Exec(dbQuery, teacher.TeacherUID, teacher.FirstName, teacher.LastName, teacher.DateOfBirth, teacher.Address,
		teacher.PhoneNumber, teacher.Email, teacher.InstitutePID, teacher.PID)
	if err != nil {
		err = fmt.Errorf("[%s] - %s", serverErrorMessages[seDBResourceQuery], err.Error())
		return err
	}

	logging.Debugmf(logModTeacherHandler, "Updated teacher in DB (PID=%v)", teacher.PID)
	return nil
}

// delete teacher, return #delete rows, error
func deleteTeacher(pid int) (int, error) {
	var err error
	var result sql.Result

	defer func() {
		if err != nil {
			logging.Errormf(logModTeacherHandler, err.Error())
		}
	}()

	var dbQuery = "DELETE FROM teacher"
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
		logging.Warnmf(logModTeacherHandler, "Cound not count #deleted teachers")
	}
	logging.Debugmf(logModTeacherHandler, "Deleted %d teacher results from DB", rowsAffected)
	return int(rowsAffected), nil
}
