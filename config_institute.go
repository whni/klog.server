package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var instituteConfigHandlerTable = map[string]gin.HandlerFunc{
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
	deletedRows, err = deleteInstitute(pid)
	if err != nil {
		response.Status = http.StatusConflict
		response.Message = err.Error()
		return
	}
	response.Payload = deletedRows
	return
}

// find institute, return institute slice, error
func findInstitute(pid primitive.ObjectID) ([]*Institute, error) {
	var err error
	defer func() {
		if err != nil {
			logging.Errormf(logModInstituteMgmt, err.Error())
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

	logging.Debugmf(logModInstituteMgmt, "Found %d institute results from DB (PID=%v)", len(institutes), pid.Hex())
	return institutes, nil
}

// create institute, return PID, error
func createInstitute(institute *Institute) (primitive.ObjectID, error) {
	var err error
	defer func() {
		if err != nil {
			logging.Errormf(logModInstituteMgmt, err.Error())
		}
	}()

	insertResult, err := dbPool.Collection(DBCollectionInstitute).InsertOne(context.TODO(), institute)
	if err != nil {
		err = fmt.Errorf("[%s] - %s", serverErrorMessages[seDBResourceQuery], err.Error())
		return primitive.NilObjectID, err
	}

	lastInsertID := insertResult.InsertedID.(primitive.ObjectID)
	logging.Debugmf(logModInstituteMgmt, "Created institute in DB (LastInsertID,PID=%s)", lastInsertID.Hex())
	return lastInsertID, nil
}

// update institute, return error
func updateInstitute(institute *Institute) error {
	var err error
	defer func() {
		if err != nil {
			logging.Errormf(logModInstituteMgmt, err.Error())
		}
	}()

	// check institute PID
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

	logging.Debugmf(logModInstituteMgmt, "Update institute (PID %s): matched %d modified %d",
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

// delete institute, return #delete entries, error
func deleteInstitute(pid primitive.ObjectID) (int, error) {
	var err error
	defer func() {
		if err != nil {
			logging.Errormf(logModInstituteMgmt, err.Error())
		}
	}()

	var deleteFilter bson.D = bson.D{}
	var dependencyFindFilter bson.D = bson.D{}
	if !pid.IsZero() {
		deleteFilter = append(deleteFilter, bson.E{"_id", pid})
		dependencyFindFilter = append(dependencyFindFilter, bson.E{"institute_pid", pid})
	}

	// check teacher dependency
	var teacher Teacher
	if dbPool.Collection(DBCollectionTeacher).FindOne(context.TODO(), dependencyFindFilter).Decode(&teacher) == nil {
		err = fmt.Errorf("[%s] - teacher-institute dependency unresolved (e.g. teacher PID %s institute PID %s)",
			serverErrorMessages[seDependencyIssue], teacher.PID.Hex(), teacher.InstitutePID.Hex())
		return 0, err
	}

	// check course dependency
	var course Course
	if dbPool.Collection(DBCollectionCourse).FindOne(context.TODO(), dependencyFindFilter).Decode(&course) == nil {
		err = fmt.Errorf("[%s] - course-institute dependency unresolved (e.g. course PID %s institute PID %s)",
			serverErrorMessages[seDependencyIssue], course.PID.Hex(), course.InstitutePID.Hex())
		return 0, err
	}

	deleteResult, err := dbPool.Collection(DBCollectionInstitute).DeleteMany(context.TODO(), deleteFilter)
	if err != nil {
		err = fmt.Errorf("[%s] - %s", serverErrorMessages[seDBResourceQuery], err.Error())
		return 0, err
	}

	logging.Debugmf(logModInstituteMgmt, "Deleted %d institute results from DB", deleteResult.DeletedCount)
	return int(deleteResult.DeletedCount), nil
}
