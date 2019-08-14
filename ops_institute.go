package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	_ "go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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
	var pid primitive.ObjectID
	if params.PID == "all" {
		pid = primitive.NilObjectID
	} else {
		pid, err = primitive.ObjectIDFromHex(params.PID)
		if err != nil {
			response.Status = http.StatusBadRequest
			response.Message = fmt.Sprintf("[%s] - Please specifiy a valid PID (mongoDB ObjectID)", serverErrorMessages[seInputParamNotValid])
			return
		}
	}

	// pid: nil objectid for all, others for specified one
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
	var institutePID primitive.ObjectID
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
	pid, err = strconv.Atoi(params.PID)
	if err != nil || pid < 0 {
		response.Status = http.StatusBadRequest
		response.Message = fmt.Sprintf("[%s] - Please specifiy a valid PID (pid >= 0)", serverErrorMessages[seInputParamNotValid])
		return
	}

	// pid: 0 for all, > 0 for specified one
	deletedRows, err = deleteInstitute(pid)
	if err != nil {
		response.Status = http.StatusConflict
		response.Message = err.Error()
		return
	}
	response.Payload = deletedRows
	return
}

func deleteInstitute(pid int) (int, error) {
	return 0, nil
}

// find institute, return institute slice, error
func findInstitute(pid primitive.ObjectID) ([]*Institute, error) {
	var err error
	defer func() {
		if err != nil {
			logging.Errormf(logModInstituteHandler, err.Error())
		}
	}()

	var findOptions = options.Find()
	var findFilter bson.D
	if pid.IsZero() {
		findFilter = bson.D{{}}
	} else {
		findOptions.SetLimit(1)
		findFilter = bson.D{{"_id", pid}}
	}

	findCursor, err := dbPool.Collection(DBCollectionInstitute).Find(context.TODO(), findFilter, findOptions)
	if err != nil {
		err = fmt.Errorf("[%s] - %s", serverErrorMessages[seDBResourceQuery], err.Error())
		return nil, err
	}

	institutes := []*Institute{}
	for findCursor.Next(context.TODO()) {
		var institute Institute
		err = findCursor.Decode(&institute)
		if err != nil {
			err = fmt.Errorf("[%s] - %s", serverErrorMessages[seDBResourceQuery], err.Error())
			return nil, err
		}
		institutes = append(institutes, &institute)
	}

	err = findCursor.Err()
	if err != nil {
		err = fmt.Errorf("[%s] - %s", serverErrorMessages[seDBResourceQuery], err.Error())
		return nil, err
	}

	logging.Debugmf(logModInstituteHandler, "Found %d institute results from DB (PID=%v)", len(institutes), pid)
	return institutes, nil
}

// create institute, return PID, error
func createInstitute(institute *Institute) (primitive.ObjectID, error) {
	var err error
	defer func() {
		if err != nil {
			logging.Errormf(logModInstituteHandler, err.Error())
		}
	}()

	insertResult, err := dbPool.Collection(DBCollectionInstitute).InsertOne(context.TODO(), institute)
	if err != nil {
		err = fmt.Errorf("[%s] - %s", serverErrorMessages[seDBResourceQuery], err.Error())
		return primitive.NilObjectID, err
	}

	lastInsertID := insertResult.InsertedID.(primitive.ObjectID)
	logging.Debugmf(logModInstituteHandler, "Created institute in DB (LastInsertID,PID=%s)", lastInsertID.Hex())
	return lastInsertID, nil
}

// update institute, return error
func updateInstitute(institute *Institute) error {
	var err error
	defer func() {
		if err != nil {
			logging.Errormf(logModInstituteHandler, err.Error())
		}
	}()

	if institute.PID.IsZero() {
		err = fmt.Errorf("[%s] - institute PID is empty", serverErrorMessages[seInputJSONNotValid])
		return err
	}

	var updateFilter = bson.D{{"_id", institute.PID}}
	var updateBSONDocument = bson.D{}
	instituteBSONData, err := bson.Marshal(institute)
	if err != nil {
		err = fmt.Errorf("[%s] - could not convert institute (PID %s) to bson data", serverErrorMessages[seInputBSONNotValid], institute.PID.Hex())
		return err
	}
	err = bson.Unmarshal(instituteBSONData, &updateBSONDocument)
	if err != nil {
		err = fmt.Errorf("[%s] - could not convert institute (PID %s) to bson document", serverErrorMessages[seInputBSONNotValid], institute.PID.Hex())
		return err
	}
	var updateOptions = bson.D{{"$set", updateBSONDocument}}

	insertResult, err := dbPool.Collection(DBCollectionInstitute).UpdateOne(context.TODO(), updateFilter, updateOptions)
	if err != nil {
		err = fmt.Errorf("[%s] - %s", serverErrorMessages[seDBResourceQuery], err.Error())
		return err
	}

	logging.Debugmf(logModInstituteHandler, "Update institute (PID %s): matched %d modified %d",
		institute.PID.Hex(), insertResult.MatchedCount, insertResult.ModifiedCount)
	if insertResult.MatchedCount == 0 {
		err = fmt.Errorf("[%s] - could not find institute (PID %s)", serverErrorMessages[seResourceNotFound], institute.PID.Hex())
		return err
	} else if insertResult.ModifiedCount == 0 {
		err = fmt.Errorf("[%s] - institute (PID %s) not changed", serverErrorMessages[seResourceNotChange], institute.PID.Hex())
		return err
	}
	return nil
}

/*
// delete institute, return #delete rows, error
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

*/
