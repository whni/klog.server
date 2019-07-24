package main

import (
	"database/sql"
	"encoding/json"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"net/http"
	"strconv"
)

var instituteHandlerTable = map[string]gin.HandlerFunc{
	"get":  instituteGetHandler,
	"post": institutePostHandler,
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
	var errCode int = seNoError
	if params.PID == "all" {
		institutes, errCode = findInstitute(0)
		if errCode < seNoError {
			response.Status = http.StatusConflict
			response.Message = serverErrorMessages[errCode]
			return
		}
		response.Payload = institutes
		return
	}

	pid, err := strconv.Atoi(params.PID)
	if err != nil || pid <= 0 {
		response.Status = http.StatusBadRequest
		response.Message = serverErrorMessages[seInputParamNotValid] + " - Please specifiy a valid institute PID (pkey > 0)."
		return
	}

	institutes, errCode = findInstitute(pid)
	if errCode < seNoError {
		response.Status = http.StatusConflict
		response.Message = serverErrorMessages[errCode]
		return
	}
	if len(institutes) == 0 {
		response.Status = http.StatusNotFound
		response.Message = serverErrorMessages[seResourceNotFound]
	} else {
		response.Payload = institutes[0]
	}
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
	var institutePID int = 0
	var errCode int = seNoError

	if err := json.Unmarshal(params.Data, &institute); err != nil {
		response.Status = http.StatusBadRequest
		response.Message = serverErrorMessages[seInputJSONNotValid]
		return
	}

	institutePID, errCode = createInstitute(&institute)
	if errCode < seNoError {
		response.Status = http.StatusConflict
		response.Message = serverErrorMessages[errCode]
	} else {
		response.Payload = institutePID
	}
	return
}

// find institute, return institute ptr, errCode
func findInstitute(pid int) ([]*Institute, int) {
	var institutes []*Institute
	var rows *sql.Rows
	var err error
	var dbQuery = "SELECT pid, institute_uid, institute_name, address, country_code, create_ts, modify_ts FROM institute"

	if pid == 0 {
		rows, err = dbPool.Query(dbQuery)
	} else if pid > 0 {
		rows, err = dbPool.Query(dbQuery+" WHERE pid = ?", pid)
	} else {
		logging.Errormf(logModInstituteHandler, serverErrorMessages[seInputParamNotValid]+" (PID not valid)")
		return nil, seInputParamNotValid
	}
	if err != nil {
		logging.Errormf(logModInstituteHandler, "Could not query DB rows - Error msg: %s", err.Error())
		return nil, seDBResourceQuery
	}
	defer rows.Close()

	for rows.Next() {
		var institute Institute
		err = rows.Scan(&institute.PID, &institute.InstituteUID, &institute.InstituteName, &institute.Address,
			&institute.CountryCode, &institute.CreateTS, &institute.ModifyTS)
		if err != nil {
			logging.Errormf(logModInstituteHandler, "Could not retrieve data from DB rows - Error msg: %s", err.Error())
			return nil, seDBResourceQuery
		}
		institutes = append(institutes, &institute)
	}

	err = rows.Err()
	if err != nil {
		logging.Errormf(logModInstituteHandler, "DB rows processing error occurs - Error msg: %s", err.Error())
		return nil, seDBResourceQuery
	}

	logging.Debugmf(logModInstituteHandler, "Found %d institute results from DB (PID=%d)", len(institutes), pid)
	return institutes, seNoError
}

// create institute, return PID, errCode
func createInstitute(institute *Institute) (int, int) {
	var err error
	var result sql.Result
	var dbQuery = "INSERT INTO institute(institute_uid, institute_name, address, country_code) VALUES (?, ?, ?, ?)"
	var createWithPID = false

	if !ginInputStructValid(institute) {
		return 0, seInputSchemaNotValid
	}

	if institute.PID > 0 {
		if institutes, errCode := findInstitute(institute.PID); errCode == seNoError && len(institutes) > 0 {
			return 0, seResourceDuplicated
		}
		createWithPID = true
		dbQuery = "INSERT INTO institute(pid, institute_uid, institute_name, address, country_code) VALUES (?, ?, ?, ?, ?)"
	} else if institute.PID < 0 {
		logging.Errormf(logModInstituteHandler, serverErrorMessages[seInputParamNotValid]+" (PID not valid)")
		return 0, seInputSchemaNotValid
	}

	if createWithPID == true {
		result, err = dbPool.Exec(dbQuery, institute.PID, institute.InstituteUID, institute.InstituteName, institute.Address, institute.CountryCode)
	} else {
		result, err = dbPool.Exec(dbQuery, institute.InstituteUID, institute.InstituteName, institute.Address, institute.CountryCode)
	}
	if err != nil {
		logging.Errormf(logModInstituteHandler, "Could not execute DB query correctly - Error msg: %s", err.Error())
		return 0, seDBResourceQuery
	}

	lastInsertID, err := result.LastInsertId()
	if err != nil {
		lastInsertID = -1
		logging.Warnmf(logModInstituteHandler, "Cound not retrieve created institute PID")
	}
	logging.Debugmf(logModInstituteHandler, "Created institute in DB (LastInsertId=%v)", lastInsertID)
	return int(lastInsertID), seNoError
}
