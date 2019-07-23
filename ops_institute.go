package main

import (
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"net/http"
	"strconv"
)

var instituteHandlerTable = map[string]gin.HandlerFunc{
	"get": instituteGetHandler,
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
	if params.Pkey == "all" {
		institutes, err = findInstitute(0)
		if err != nil {
			response.Status = http.StatusInternalServerError
			response.Message = err.Error()
			return
		}
		response.Payload = institutes
		return
	}

	pid, err := strconv.Atoi(params.Pkey)
	if err != nil || pid <= 0 {
		response.Status = http.StatusBadRequest
		response.Message = serverErrorMessages[serrInputNotValid] + " Please specifiy a valid institute PID (pkey > 0)."
		return
	}

	institutes, err = findInstitute(pid)
	if err != nil {
		response.Status = http.StatusInternalServerError
		response.Message = err.Error()
		return
	}
	if len(institutes) == 0 {
		response.Status = http.StatusNotFound
		response.Message = fmt.Sprintf(serverErrorMessages[serrResourceNotFound], "institute")
	} else {
		response.Payload = institutes[0]
	}
}

func findInstitute(pid int) ([]*Institute, error) {
	var institutes []*Institute
	var rows *sql.Rows
	var err error
	var dbQuery = "SELECT pid, institute_uid, institute_name, address, country_code, create_ts, modify_ts FROM institute"

	if pid == 0 {
		rows, err = dbPool.Query(dbQuery)
	} else if pid > 0 {
		rows, err = dbPool.Query(dbQuery+" WHERE pid = ?", pid)
	} else {
		err = fmt.Errorf("invalid institute PID (%d) given", pid)
		logging.Errormf(logModInstituteHandler, err.Error())
		return nil, err
	}
	if err != nil {
		logging.Errormf(logModInstituteHandler, "Could not query DB rows - Error msg: %s", err.Error())
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var institute Institute
		err = rows.Scan(&institute.PID, &institute.InstituteUID, &institute.InstituteName, &institute.Address,
			&institute.CountryCode, &institute.CreateTS, &institute.ModifyTS)
		if err != nil {
			logging.Errormf(logModInstituteHandler, "Could not retrieve data from DB rows - Error msg: %s", err.Error())
			return nil, err
		}
		institutes = append(institutes, &institute)
	}

	err = rows.Err()
	if err != nil {
		logging.Errormf(logModInstituteHandler, "DB rows processing error occurs - Error msg: %s", err.Error())
		return nil, err
	}

	return institutes, nil
}

func createInstitute(institute *Institute) error {

	return nil
}
