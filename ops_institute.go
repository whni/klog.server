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

var instituteHandlerTable = map[string]gin.HandlerFunc{
	"get":    instituteGetHandler,
	"post":   institutePostHandler,
	"put":    institutePutHandler,
	"delete": instituteDeleteHandler,
}

func instituteGetHandler(ctx *gin.Context) {
	params := ginContextRequestParameter(ctx)
	response := GinResponse{
		Status: http.StatusOK,
	}
	defer func() {
		ginContextProcessResponse(ctx, &response)
	}()

	var institutes []*Institute
	var err error
	var pid int
	if params.PID == "all" {
		pid = 0
	} else {
		pid, err := strconv.Atoi(params.PID)
		if err != nil || pid <= 0 {
			response.Status = http.StatusBadRequest
			response.Message = fmt.Sprintf("[%s] - Please specifiy a valid institute PID (pid > 0)", serverErrorMessages[seInputParamNotValid])
			return
		}
	}

	institutes, err = findInstitute(pid)
	if err != nil {
		response.Status = http.StatusConflict
		response.Message = err.Error()
		return
	}
	response.Payload = institutes
	return
}

func institutePostHandler(ctx *gin.Context) {
	params := ginContextRequestParameter(ctx)
	response := GinResponse{
		Status: http.StatusOK,
	}
	defer func() {
		ginContextProcessResponse(ctx, &response)
	}()

	var institute Institute
	var institutePID int
	var err error

	if err = json.Unmarshal(params.Data, &institute); err != nil {
		response.Status = http.StatusBadRequest
		response.Message = fmt.Sprintf("[%s] - %s", serverErrorMessages[seInputJSONNotValid], err.Error())
		return
	}

	institutePID, err = createInstitute(&institute)
	if err != nil {
		response.Status = http.StatusConflict
		response.Message = err.Error()
	} else {
		response.Payload = institutePID
	}
	return
}

func institutePutHandler(ctx *gin.Context) {
	params := ginContextRequestParameter(ctx)
	response := GinResponse{
		Status: http.StatusOK,
	}
	defer func() {
		ginContextProcessResponse(ctx, &response)
	}()

	var institute Institute
	var err error

	if err = json.Unmarshal(params.Data, &institute); err != nil {
		response.Status = http.StatusBadRequest
		response.Message = fmt.Sprintf("[%s] - %s", serverErrorMessages[seInputJSONNotValid], err.Error())
		return
	}

	err = updateInstitute(&institute)
	if err != nil {
		response.Status = http.StatusConflict
		response.Message = err.Error()
	} else {
		response.Payload = institute.PID
	}
	return
}

func instituteDeleteHandler(ctx *gin.Context) {
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
	if params.PID == "all" {
		pid = 0
	} else {
		pid, err = strconv.Atoi(params.PID)
		if err != nil || pid <= 0 {
			response.Status = http.StatusBadRequest
			response.Message = fmt.Sprintf("[%s] - Please specifiy a valid institute PID (pid > 0)", serverErrorMessages[seInputParamNotValid])
			return
		}
	}

	deletedRows, err = deleteInstitute(pid)
	if err != nil {
		response.Status = http.StatusConflict
		response.Message = err.Error()
		return
	}
	response.Payload = deletedRows
	return
}

// find institute, return institute ptr, error
func findInstitute(pid int) ([]*Institute, error) {
	var institutes []*Institute
	var rows *sql.Rows
	var err error
	defer func() {
		if err != nil {
			logging.Errormf(logModInstituteHandler, err.Error())
		}
	}()

	var dbQuery = "SELECT pid, institute_uid, institute_name, address, country_code, create_ts, modify_ts FROM institute"
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

	for rows.Next() {
		var institute Institute
		err = rows.Scan(&institute.PID, &institute.InstituteUID, &institute.InstituteName, &institute.Address,
			&institute.CountryCode, &institute.CreateTS, &institute.ModifyTS)
		if err != nil {
			err = fmt.Errorf("[%s] - %s", serverErrorMessages[seDBResourceQuery], err.Error())
			return nil, err
		}
		institutes = append(institutes, &institute)
	}

	err = rows.Err()
	if err != nil {
		err = fmt.Errorf("[%s] - %s", serverErrorMessages[seDBResourceQuery], err.Error())
		return nil, err
	}

	logging.Debugmf(logModInstituteHandler, "Found %d institute results from DB (PID=%d)", len(institutes), pid)
	return institutes, nil
}

// create institute, return PID, error
func createInstitute(institute *Institute) (int, error) {
	var err error
	var result sql.Result
	var createWithPID = false

	defer func() {
		if err != nil {
			logging.Errormf(logModInstituteHandler, err.Error())
		}
	}()

	if err = ginStructValidCheck(institute); err != nil {
		return 0, err
	}

	if institute.PID > 0 {
		if institutes, errExist := findInstitute(institute.PID); errExist == nil && len(institutes) > 0 {
			err = fmt.Errorf("[%s] - Institute (PID=%d) already exists", serverErrorMessages[seResourceDuplicated], institute.PID)
			return 0, err
		}
		createWithPID = true
	} else if institute.PID < 0 {
		err = fmt.Errorf("[%s] - PID (%d) not valid", serverErrorMessages[seInputParamNotValid], institute.PID)
		return 0, err
	}

	var dbQuery string
	if createWithPID == true {
		dbQuery = "INSERT INTO institute(pid, institute_uid, institute_name, address, country_code) VALUES (?, ?, ?, ?, ?)"
		result, err = dbPool.Exec(dbQuery, institute.PID, institute.InstituteUID, institute.InstituteName, institute.Address, institute.CountryCode)
	} else {
		dbQuery = "INSERT INTO institute(institute_uid, institute_name, address, country_code) VALUES (?, ?, ?, ?)"
		result, err = dbPool.Exec(dbQuery, institute.InstituteUID, institute.InstituteName, institute.Address, institute.CountryCode)
	}
	if err != nil {
		err = fmt.Errorf("[%s] - %s", serverErrorMessages[seDBResourceQuery], err.Error())
		return 0, err
	}

	lastInsertID, err := result.LastInsertId()
	if err != nil {
		lastInsertID = -1
		logging.Warnmf(logModInstituteHandler, "Cound not retrieve created institute PID")
	}
	logging.Debugmf(logModInstituteHandler, "Created institute in DB (LastInsertId,PID=%v)", lastInsertID)
	return int(lastInsertID), nil
}

// update institute, return error
func updateInstitute(institute *Institute) error {
	var err error
	var institutes []*Institute

	defer func() {
		if err != nil {
			logging.Errormf(logModInstituteHandler, err.Error())
		}
	}()

	if err = ginStructValidCheck(institute); err != nil {
		return err
	}

	if institute.PID <= 0 {
		err = fmt.Errorf("[%s] - PID (%d) not valid", serverErrorMessages[seInputParamNotValid], institute.PID)
		return err
	}
	if institutes, err = findInstitute(institute.PID); err != nil || len(institutes) == 0 {
		err = fmt.Errorf("[%s] - Institute (PID=%d) does not exist ==> not updated", serverErrorMessages[seResourceNotFound], institute.PID)
		return err
	}

	var dbQuery = "UPDATE institute SET institute_uid=?, institute_name=?, address=?, country_code=? WHERE pid=?"
	_, err = dbPool.Exec(dbQuery, institute.InstituteUID, institute.InstituteName, institute.Address, institute.CountryCode, institute.PID)
	if err != nil {
		err = fmt.Errorf("[%s] - %s", serverErrorMessages[seDBResourceQuery], err.Error())
		return err
	}

	logging.Debugmf(logModInstituteHandler, "Updated institute in DB (PID=%v)", institute.PID)
	return nil
}

// delete institute, return #d(elete rows, error
func deleteInstitute(pid int) (int, error) {
	var err error
	var result sql.Result

	defer func() {
		if err != nil {
			logging.Errormf(logModInstituteHandler, err.Error())
		}
	}()

	var dbQuery = "DELETE FROM institute"
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
		logging.Warnmf(logModInstituteHandler, "Cound not count #deleted institutes")
	}
	logging.Debugmf(logModInstituteHandler, "Deleted %d institute results from DB", rowsAffected)
	return int(rowsAffected), nil
}
