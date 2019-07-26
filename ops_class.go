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

var classHandlerTable = map[string]gin.HandlerFunc{
	"get":    classGetHandler,
	"post":   classPostHandler,
	"put":    classPutHandler,
	"delete": classDeleteHandler,
}

func classGetHandler(ctx *gin.Context) {
	params := ginContextRequestParameter(ctx)
	response := GinResponse{
		Status: http.StatusOK,
	}
	defer func() {
		ginContextProcessResponse(ctx, &response)
	}()

	var classes []*Class
	var err error
	var pid int
	pid, err = strconv.Atoi(params.PID)
	if err != nil || pid < 0 {
		response.Status = http.StatusBadRequest
		response.Message = fmt.Sprintf("[%s] - Please specifiy a valid PID (pid >= 0)", serverErrorMessages[seInputParamNotValid])
		return
	}

	// pid: 0 for all, > 0 for specified one
	classes, err = findClass(pid)
	if err != nil {
		response.Status = http.StatusConflict
		response.Message = err.Error()
		return
	}
	response.Payload = classes
	return
}

func classPostHandler(ctx *gin.Context) {
	params := ginContextRequestParameter(ctx)
	response := GinResponse{
		Status: http.StatusOK,
	}
	defer func() {
		ginContextProcessResponse(ctx, &response)
	}()

	var class Class
	var classPID int
	var err error

	if err = json.Unmarshal(params.Data, &class); err != nil {
		response.Status = http.StatusBadRequest
		response.Message = fmt.Sprintf("[%s] - %s", serverErrorMessages[seInputJSONNotValid], err.Error())
		return
	}

	classPID, err = createClass(&class)
	if err != nil {
		response.Status = http.StatusConflict
		response.Message = err.Error()
	} else {
		response.Payload = classPID
	}
	return
}

func classPutHandler(ctx *gin.Context) {
	params := ginContextRequestParameter(ctx)
	response := GinResponse{
		Status: http.StatusOK,
	}
	defer func() {
		ginContextProcessResponse(ctx, &response)
	}()

	var class Class
	var err error

	if err = json.Unmarshal(params.Data, &class); err != nil {
		response.Status = http.StatusBadRequest
		response.Message = fmt.Sprintf("[%s] - %s", serverErrorMessages[seInputJSONNotValid], err.Error())
		return
	}

	err = updateClass(&class)
	if err != nil {
		response.Status = http.StatusConflict
		response.Message = err.Error()
	} else {
		response.Payload = class.PID
	}
	return
}

func classDeleteHandler(ctx *gin.Context) {
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
	deletedRows, err = deleteClass(pid)
	if err != nil {
		response.Status = http.StatusConflict
		response.Message = err.Error()
		return
	}
	response.Payload = deletedRows
	return
}

// find class, return class ptr, error
func findClass(pid int) ([]*Class, error) {
	var rows *sql.Rows
	var err error
	defer func() {
		if err != nil {
			logging.Errormf(logModClassHandler, err.Error())
		}
	}()

	var dbQuery = "SELECT pid, class_uid, class_name, location, institute_pid, create_ts, modify_ts FROM class"
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

	classes := []*Class{}
	for rows.Next() {
		var class Class
		err = rows.Scan(&class.PID, &class.ClassUID, &class.ClassName, &class.Location,
			&class.InstitutePID, &class.CreateTS, &class.ModifyTS)
		if err != nil {
			err = fmt.Errorf("[%s] - %s", serverErrorMessages[seDBResourceQuery], err.Error())
			return nil, err
		}
		classes = append(classes, &class)
	}

	err = rows.Err()
	if err != nil {
		err = fmt.Errorf("[%s] - %s", serverErrorMessages[seDBResourceQuery], err.Error())
		return nil, err
	}

	logging.Debugmf(logModClassHandler, "Found %d class results from DB (PID=%d)", len(classes), pid)
	return classes, nil
}

// create class, return PID, error
func createClass(class *Class) (int, error) {
	var err error
	var result sql.Result
	var createWithPID = false

	defer func() {
		if err != nil {
			logging.Errormf(logModClassHandler, err.Error())
		}
	}()

	if err = ginStructValidCheck(class); err != nil {
		return 0, err
	}

	if class.PID > 0 {
		if classes, errExist := findClass(class.PID); errExist == nil && len(classes) > 0 {
			err = fmt.Errorf("[%s] - Class (PID=%d) already exists", serverErrorMessages[seResourceDuplicated], class.PID)
			return 0, err
		}
		createWithPID = true
	} else if class.PID < 0 {
		err = fmt.Errorf("[%s] - PID (%d) not valid", serverErrorMessages[seInputParamNotValid], class.PID)
		return 0, err
	}

	var dbQuery string
	if createWithPID == true {
		dbQuery = "INSERT INTO class(pid, class_uid, class_name, location, institute_pid) VALUES (?, ?, ?, ?, ?)"
		result, err = dbPool.Exec(dbQuery, class.PID, class.ClassUID, class.ClassName, class.Location, class.InstitutePID)
	} else {
		dbQuery = "INSERT INTO class(class_uid, class_name, location, institute_pid) VALUES (?, ?, ?, ?)"
		result, err = dbPool.Exec(dbQuery, class.ClassUID, class.ClassName, class.Location, class.InstitutePID)
	}
	if err != nil {
		err = fmt.Errorf("[%s] - %s", serverErrorMessages[seDBResourceQuery], err.Error())
		return 0, err
	}

	lastInsertID, err := result.LastInsertId()
	if err != nil {
		lastInsertID = -1
		logging.Warnmf(logModClassHandler, "Cound not retrieve created class PID")
	}
	logging.Debugmf(logModClassHandler, "Created class in DB (LastInsertId,PID=%v)", lastInsertID)
	return int(lastInsertID), nil
}

// update class, return error
func updateClass(class *Class) error {
	var err error
	defer func() {
		if err != nil {
			logging.Errormf(logModClassHandler, err.Error())
		}
	}()

	if err = ginStructValidCheck(class); err != nil {
		return err
	}

	if class.PID <= 0 {
		err = fmt.Errorf("[%s] - PID (%d) not valid", serverErrorMessages[seInputParamNotValid], class.PID)
		return err
	}
	if classes, errExist := findClass(class.PID); errExist != nil || len(classes) == 0 {
		err = fmt.Errorf("[%s] - Class (PID=%d) does not exist ==> not updated", serverErrorMessages[seResourceNotFound], class.PID)
		return err
	}

	var dbQuery = "UPDATE class SET class_uid=?, class_name=?, location=?, institute_pid=? WHERE pid=?"
	_, err = dbPool.Exec(dbQuery, class.ClassUID, class.ClassName, class.Location, class.InstitutePID, class.PID)
	if err != nil {
		err = fmt.Errorf("[%s] - %s", serverErrorMessages[seDBResourceQuery], err.Error())
		return err
	}

	logging.Debugmf(logModClassHandler, "Updated class in DB (PID=%v)", class.PID)
	return nil
}

// delete class, return #delete rows, error
func deleteClass(pid int) (int, error) {
	var err error
	var result sql.Result

	defer func() {
		if err != nil {
			logging.Errormf(logModClassHandler, err.Error())
		}
	}()

	var dbQuery = "DELETE FROM class"
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
		logging.Warnmf(logModClassHandler, "Cound not count #deleted classes")
	}
	logging.Debugmf(logModClassHandler, "Deleted %d class results from DB", rowsAffected)
	return int(rowsAffected), nil
}
