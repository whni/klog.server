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

var parentHandlerTable = map[string]gin.HandlerFunc{
	"get":    parentGetHandler,
	"post":   parentPostHandler,
	"put":    parentPutHandler,
	"delete": parentDeleteHandler,
}

func parentGetHandler(ctx *gin.Context) {
	params := ginContextRequestParameter(ctx)
	response := GinResponse{
		Status: http.StatusOK,
	}
	defer func() {
		ginContextProcessResponse(ctx, &response)
	}()

	var parents []*Parent
	var err error
	var pid int
	pid, err = strconv.Atoi(params.PID)
	if err != nil || pid < 0 {
		response.Status = http.StatusBadRequest
		response.Message = fmt.Sprintf("[%s] - Please specifiy a valid PID (pid >= 0)", serverErrorMessages[seInputParamNotValid])
		return
	}

	// pid: 0 for all, > 0 for specified one
	parents, err = findParent(pid)
	if err != nil {
		response.Status = http.StatusConflict
		response.Message = err.Error()
		return
	}
	response.Payload = parents
	return
}

func parentPostHandler(ctx *gin.Context) {
	params := ginContextRequestParameter(ctx)
	response := GinResponse{
		Status: http.StatusOK,
	}
	defer func() {
		ginContextProcessResponse(ctx, &response)
	}()

	var parent Parent
	var parentPID int
	var err error

	if err = json.Unmarshal(params.Data, &parent); err != nil {
		response.Status = http.StatusBadRequest
		response.Message = fmt.Sprintf("[%s] - %s", serverErrorMessages[seInputJSONNotValid], err.Error())
		return
	}

	parentPID, err = createParent(&parent)
	if err != nil {
		response.Status = http.StatusConflict
		response.Message = err.Error()
	} else {
		response.Payload = parentPID
	}
	return
}

func parentPutHandler(ctx *gin.Context) {
	params := ginContextRequestParameter(ctx)
	response := GinResponse{
		Status: http.StatusOK,
	}
	defer func() {
		ginContextProcessResponse(ctx, &response)
	}()

	var parent Parent
	var err error

	if err = json.Unmarshal(params.Data, &parent); err != nil {
		response.Status = http.StatusBadRequest
		response.Message = fmt.Sprintf("[%s] - %s", serverErrorMessages[seInputJSONNotValid], err.Error())
		return
	}

	err = updateParent(&parent)
	if err != nil {
		response.Status = http.StatusConflict
		response.Message = err.Error()
	} else {
		response.Payload = parent.PID
	}
	return
}

func parentDeleteHandler(ctx *gin.Context) {
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
	deletedRows, err = deleteParent(pid)
	if err != nil {
		response.Status = http.StatusConflict
		response.Message = err.Error()
		return
	}
	response.Payload = deletedRows
	return
}

// find parent, return parent slice, error
func findParent(pid int) ([]*Parent, error) {
	var rows *sql.Rows
	var err error
	defer func() {
		if err != nil {
			logging.Errormf(logModParentHandler, err.Error())
		}
	}()

	var dbQuery = "SELECT pid, parent_uid, first_name, last_name, date_of_birth, address, phone_number, email, occupation, create_ts, modify_ts FROM parent"
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

	parents := []*Parent{}
	for rows.Next() {
		var parent Parent
		err = rows.Scan(&parent.PID, &parent.ParentUID, &parent.FirstName, &parent.LastName, &parent.DateOfBirth,
			&parent.Address, &parent.PhoneNumber, &parent.Email, &parent.Occupation, &parent.CreateTS, &parent.ModifyTS)
		if err != nil {
			err = fmt.Errorf("[%s] - %s", serverErrorMessages[seDBResourceQuery], err.Error())
			return nil, err
		}
		parents = append(parents, &parent)
	}

	err = rows.Err()
	if err != nil {
		err = fmt.Errorf("[%s] - %s", serverErrorMessages[seDBResourceQuery], err.Error())
		return nil, err
	}

	logging.Debugmf(logModParentHandler, "Found %d parent results from DB (PID=%d)", len(parents), pid)
	return parents, nil
}

// create parent, return PID, error
func createParent(parent *Parent) (int, error) {
	var err error
	var result sql.Result

	defer func() {
		if err != nil {
			logging.Errormf(logModParentHandler, err.Error())
		}
	}()

	if err = ginStructValidCheck(parent); err != nil {
		return 0, err
	}

	if parent.PID > 0 {
		if parents, errExist := findParent(parent.PID); errExist == nil && len(parents) > 0 {
			err = fmt.Errorf("[%s] - Parent (PID=%d) already exists", serverErrorMessages[seResourceDuplicated], parent.PID)
			return 0, err
		}
	} else if parent.PID < 0 {
		err = fmt.Errorf("[%s] - PID (%d) not valid", serverErrorMessages[seInputParamNotValid], parent.PID)
		return 0, err
	}

	var dbQuery string
	if parent.PID > 0 {
		dbQuery = "INSERT INTO parent(pid, parent_uid, first_name, last_name, date_of_birth, address, phone_number, email, occupation) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)"
		result, err = dbPool.Exec(dbQuery, parent.PID, parent.ParentUID, parent.FirstName, parent.LastName, parent.DateOfBirth, parent.Address,
			parent.PhoneNumber, parent.Email, parent.Occupation)
	} else {
		dbQuery = "INSERT INTO parent(parent_uid, first_name, last_name, date_of_birth, address, phone_number, email, occupation) VALUES (?, ?, ?, ?, ?, ?, ?, ?)"
		result, err = dbPool.Exec(dbQuery, parent.ParentUID, parent.FirstName, parent.LastName, parent.DateOfBirth, parent.Address,
			parent.PhoneNumber, parent.Email, parent.Occupation)
	}
	if err != nil {
		err = fmt.Errorf("[%s] - %s", serverErrorMessages[seDBResourceQuery], err.Error())
		return 0, err
	}

	lastInsertID, err := result.LastInsertId()
	if err != nil {
		lastInsertID = -1
		logging.Warnmf(logModParentHandler, "Cound not retrieve created parent PID")
	}
	logging.Debugmf(logModParentHandler, "Created parent in DB (LastInsertId,PID=%v)", lastInsertID)
	return int(lastInsertID), nil
}

// update parent, return error
func updateParent(parent *Parent) error {
	var err error
	defer func() {
		if err != nil {
			logging.Errormf(logModParentHandler, err.Error())
		}
	}()

	if err = ginStructValidCheck(parent); err != nil {
		return err
	}

	if parent.PID <= 0 {
		err = fmt.Errorf("[%s] - PID (%d) not valid", serverErrorMessages[seInputParamNotValid], parent.PID)
		return err
	}
	if parents, errExist := findParent(parent.PID); errExist != nil || len(parents) == 0 {
		err = fmt.Errorf("[%s] - Parent (PID=%d) does not exist ==> not updated", serverErrorMessages[seResourceNotFound], parent.PID)
		return err
	}

	var dbQuery = "UPDATE parent SET parent_uid=?, first_name=?, last_name=?, date_of_birth=?, address=?, phone_number=?, email=?, occupation=? WHERE pid=?"
	_, err = dbPool.Exec(dbQuery, parent.ParentUID, parent.FirstName, parent.LastName, parent.DateOfBirth, parent.Address,
		parent.PhoneNumber, parent.Email, parent.Occupation, parent.PID)
	if err != nil {
		err = fmt.Errorf("[%s] - %s", serverErrorMessages[seDBResourceQuery], err.Error())
		return err
	}

	logging.Debugmf(logModParentHandler, "Updated parent in DB (PID=%v)", parent.PID)
	return nil
}

// delete parent, return #delete rows, error
func deleteParent(pid int) (int, error) {
	var err error
	var result sql.Result

	defer func() {
		if err != nil {
			logging.Errormf(logModParentHandler, err.Error())
		}
	}()

	var dbQuery = "DELETE FROM parent"
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
		logging.Warnmf(logModParentHandler, "Cound not count #deleted parents")
	}
	logging.Debugmf(logModParentHandler, "Deleted %d parents results from DB", rowsAffected)
	return int(rowsAffected), nil
}
