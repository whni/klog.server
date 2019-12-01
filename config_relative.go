package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var relativeConfigHandlerTable = map[string]gin.HandlerFunc{
	"get":    relativeGetHandler,
	"post":   relativePostHandler,
	"put":    relativePutHandler,
	"delete": relativeDeleteHandler,
}

func relativeGetHandler(ctx *gin.Context) {
	params := ginContextRequestParameter(ctx)
	response := GinResponse{
		Status: http.StatusOK,
	}
	defer func() {
		ginContextProcessResponse(ctx, &response)
	}()

	var relatives []*Relative
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
	relatives, err = findRelative(pid)
	if err != nil {
		response.Status = http.StatusConflict
		response.Message = err.Error()
		return
	}
	response.Payload = relatives
	return
}

func relativePostHandler(ctx *gin.Context) {
	params := ginContextRequestParameter(ctx)
	response := GinResponse{
		Status: http.StatusOK,
	}
	defer func() {
		ginContextProcessResponse(ctx, &response)
	}()

	var relative Relative
	var relativePID primitive.ObjectID
	var err error

	if err = json.Unmarshal(params.Data, &relative); err != nil {
		response.Status = http.StatusBadRequest
		response.Message = fmt.Sprintf("[%s] - %s", serverErrorMessages[seInputJSONNotValid], err.Error())
		return
	}

	relativePID, err = createRelative(&relative)
	if err != nil {
		response.Status = http.StatusConflict
		response.Message = err.Error()
	} else {
		response.Payload = relativePID
	}
	return
}

func relativePutHandler(ctx *gin.Context) {
	params := ginContextRequestParameter(ctx)
	response := GinResponse{
		Status: http.StatusOK,
	}
	defer func() {
		ginContextProcessResponse(ctx, &response)
	}()

	var relative Relative
	var err error

	if err = json.Unmarshal(params.Data, &relative); err != nil {
		response.Status = http.StatusBadRequest
		response.Message = fmt.Sprintf("[%s] - %s", serverErrorMessages[seInputJSONNotValid], err.Error())
		return
	}

	err = updateRelative(&relative)
	if err != nil {
		response.Status = http.StatusConflict
		response.Message = err.Error()
	} else {
		response.Payload = relative.PID
	}
	return
}

func relativeDeleteHandler(ctx *gin.Context) {
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
	deletedRows, err = deleteRelative(pid)
	if err != nil {
		response.Status = http.StatusConflict
		response.Message = err.Error()
		return
	}
	response.Payload = deletedRows
	return
}

// find relative, return relative slice, error
func findRelative(pid primitive.ObjectID) ([]*Relative, error) {
	var err error
	defer func() {
		if err != nil {
			logging.Errormf(logModRelativeMgmt, err.Error())
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

	findCursor, err := dbPool.Collection(DBCollectionRelative).Find(context.TODO(), findFilter, findOptions)
	if err != nil {
		err = fmt.Errorf("[%s] - %s", serverErrorMessages[seDBResourceQuery], err.Error())
		return nil, err
	}

	relatives := []*Relative{}
	for findCursor.Next(context.TODO()) {
		var relative Relative
		err = findCursor.Decode(&relative)
		if err != nil {
			err = fmt.Errorf("[%s] - %s", serverErrorMessages[seDBResourceQuery], err.Error())
			return nil, err
		}
		relatives = append(relatives, &relative)
	}

	err = findCursor.Err()
	if err != nil {
		err = fmt.Errorf("[%s] - %s", serverErrorMessages[seDBResourceQuery], err.Error())
		return nil, err
	}

	logging.Debugmf(logModRelativeMgmt, "Found %d relative results from DB (PID=%v)", len(relatives), pid.Hex())
	return relatives, nil
}

// create relative, return PID, error
func createRelative(relative *Relative) (primitive.ObjectID, error) {
	var err error
	defer func() {
		if err != nil {
			logging.Errormf(logModRelativeMgmt, err.Error())
		}
	}()

	insertResult, err := dbPool.Collection(DBCollectionRelative).InsertOne(context.TODO(), relative)
	if err != nil {
		err = fmt.Errorf("[%s] - %s", serverErrorMessages[seDBResourceQuery], err.Error())
		return primitive.NilObjectID, err
	}

	lastInsertID := insertResult.InsertedID.(primitive.ObjectID)
	logging.Debugmf(logModRelativeMgmt, "Created relative in DB (LastInsertID,PID=%s)", lastInsertID.Hex())
	return lastInsertID, nil
}

// update relative, return error
func updateRelative(relative *Relative) error {
	var err error
	defer func() {
		if err != nil {
			logging.Errormf(logModRelativeMgmt, err.Error())
		}
	}()

	// check relative PID
	if relative.PID.IsZero() {
		err = fmt.Errorf("[%s] - relative PID is empty", serverErrorMessages[seInputJSONNotValid])
		return err
	}

	var updateFilter = bson.D{{"_id", relative.PID}}
	var updateBSONDocument = bson.D{}
	relativeBSONData, err := bson.Marshal(relative)
	if err != nil {
		err = fmt.Errorf("[%s] - could not convert relative (PID %s) to bson data", serverErrorMessages[seInputBSONNotValid], relative.PID.Hex())
		return err
	}
	err = bson.Unmarshal(relativeBSONData, &updateBSONDocument)
	if err != nil {
		err = fmt.Errorf("[%s] - could not convert relative (PID %s) to bson document", serverErrorMessages[seInputBSONNotValid], relative.PID.Hex())
		return err
	}
	var updateOptions = bson.D{{"$set", updateBSONDocument}}

	insertResult, err := dbPool.Collection(DBCollectionRelative).UpdateOne(context.TODO(), updateFilter, updateOptions)
	if err != nil {
		err = fmt.Errorf("[%s] - %s", serverErrorMessages[seDBResourceQuery], err.Error())
		return err
	}

	logging.Debugmf(logModRelativeMgmt, "Update relative (PID %s): matched %d modified %d",
		relative.PID.Hex(), insertResult.MatchedCount, insertResult.ModifiedCount)
	if insertResult.MatchedCount == 0 {
		err = fmt.Errorf("[%s] - could not find relative (PID %s)", serverErrorMessages[seResourceNotFound], relative.PID.Hex())
		return err
	} else if insertResult.ModifiedCount == 0 {
		err = fmt.Errorf("[%s] - relative (PID %s) not changed", serverErrorMessages[seResourceNotChange], relative.PID.Hex())
		return err
	}
	return nil
}

// delete relative, return #delete entries, error
func deleteRelative(pid primitive.ObjectID) (int, error) {
	var err error
	defer func() {
		if err != nil {
			logging.Errormf(logModRelativeMgmt, err.Error())
		}
	}()

	// check relative-student dependency
	studentRelativeReferences, err := findStudentRelativeRef(primitive.NilObjectID, pid)
	if err == nil && len(studentRelativeReferences) > 0 {
		err = fmt.Errorf("[%s] - student-relative dependency unresolved (e.g. student PID %s relative PID %s)",
			serverErrorMessages[seDependencyIssue], studentRelativeReferences[0].StudentPID.Hex(), studentRelativeReferences[0].RelativePID.Hex())
		return 0, err
	}

	var deleteFilter bson.D = bson.D{}
	if !pid.IsZero() {
		deleteFilter = append(deleteFilter, bson.E{"_id", pid})
	}

	deleteResult, err := dbPool.Collection(DBCollectionRelative).DeleteMany(context.TODO(), deleteFilter)
	if err != nil {
		err = fmt.Errorf("[%s] - %s", serverErrorMessages[seDBResourceQuery], err.Error())
		return 0, err
	}

	logging.Debugmf(logModRelativeMgmt, "Deleted %d relative results from DB", deleteResult.DeletedCount)
	return int(deleteResult.DeletedCount), nil
}
