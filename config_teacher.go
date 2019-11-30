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

var teacherConfigHandlerTable = map[string]gin.HandlerFunc{
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
	var teacherPID primitive.ObjectID
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
	deletedRows, err = deleteTeacher(pid)
	if err != nil {
		response.Status = http.StatusConflict
		response.Message = err.Error()
		return
	}
	response.Payload = deletedRows
	return
}

// find teacher, return teacher slice, error
func findTeacher(pid primitive.ObjectID) ([]*Teacher, error) {
	var err error
	defer func() {
		if err != nil {
			logging.Errormf(logModTeacherMgmt, err.Error())
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

	findCursor, err := dbPool.Collection(DBCollectionTeacher).Find(context.TODO(), findFilter, findOptions)
	if err != nil {
		err = fmt.Errorf("[%s] - %s", serverErrorMessages[seDBResourceQuery], err.Error())
		return nil, err
	}

	teachers := []*Teacher{}
	for findCursor.Next(context.TODO()) {
		var teacher Teacher
		err = findCursor.Decode(&teacher)
		if err != nil {
			err = fmt.Errorf("[%s] - %s", serverErrorMessages[seDBResourceQuery], err.Error())
			return nil, err
		}
		teachers = append(teachers, &teacher)
	}

	err = findCursor.Err()
	if err != nil {
		err = fmt.Errorf("[%s] - %s", serverErrorMessages[seDBResourceQuery], err.Error())
		return nil, err
	}

	logging.Debugmf(logModTeacherMgmt, "Found %d teacher results from DB (PID=%v)", len(teachers), pid.Hex())
	return teachers, nil
}

// find teacher by teacherUID
func findTeacherByUID(teacherUID string) (*Teacher, error) {
	var err error
	defer func() {
		if err != nil {
			logging.Errormf(logModTeacherMgmt, err.Error())
		}
	}()

	var teacher Teacher
	findFilter := bson.D{{"teacher_uid", teacherUID}}
	err = dbPool.Collection(DBCollectionTeacher).FindOne(context.TODO(), findFilter).Decode(&teacher)
	if err != nil {
		err = fmt.Errorf("[%s] - %s", serverErrorMessages[seDBResourceQuery], err.Error())
		return nil, err
	}

	logging.Debugmf(logModTeacherMgmt, "Found teacher from DB (teacherUID=%s)", teacherUID)
	return &teacher, nil
}

// create teacher, return PID, error
func createTeacher(teacher *Teacher) (primitive.ObjectID, error) {
	var err error
	defer func() {
		if err != nil {
			logging.Errormf(logModTeacherMgmt, err.Error())
		}
	}()

	// institute PID check
	if teacher.InstitutePID.IsZero() {
		err = fmt.Errorf("[%s] - No institute PID specified", serverErrorMessages[seResourceNotFound])
		return primitive.NilObjectID, err
	}
	institutes, err := findInstitute(teacher.InstitutePID)
	if err != nil || len(institutes) == 0 {
		err = fmt.Errorf("[%s] - No institutes found with PID %s", serverErrorMessages[seResourceNotFound], teacher.InstitutePID.Hex())
		return primitive.NilObjectID, err
	}

	insertResult, err := dbPool.Collection(DBCollectionTeacher).InsertOne(context.TODO(), teacher)
	if err != nil {
		err = fmt.Errorf("[%s] - %s", serverErrorMessages[seDBResourceQuery], err.Error())
		return primitive.NilObjectID, err
	}

	lastInsertID := insertResult.InsertedID.(primitive.ObjectID)
	logging.Debugmf(logModTeacherMgmt, "Created teacher in DB (LastInsertID,PID=%s)", lastInsertID.Hex())
	return lastInsertID, nil
}

// update teacher, return error
func updateTeacher(teacher *Teacher) error {
	var err error
	defer func() {
		if err != nil {
			logging.Errormf(logModTeacherMgmt, err.Error())
		}
	}()

	// teacher PID check
	if teacher.PID.IsZero() {
		err = fmt.Errorf("[%s] - teacher PID is empty", serverErrorMessages[seInputJSONNotValid])
		return err
	}

	// institute PID check
	if teacher.InstitutePID.IsZero() {
		err = fmt.Errorf("[%s] - No institute PID specified", serverErrorMessages[seResourceNotFound])
		return err
	}
	institutes, err := findInstitute(teacher.InstitutePID)
	if err != nil || len(institutes) == 0 {
		err = fmt.Errorf("[%s] - No institutes found with PID %s", serverErrorMessages[seResourceNotFound], teacher.InstitutePID.Hex())
		return err
	}

	var updateFilter = bson.D{{"_id", teacher.PID}}
	var updateBSONDocument = bson.D{}
	teacherBSONData, err := bson.Marshal(teacher)
	if err != nil {
		err = fmt.Errorf("[%s] - could not convert teacher (PID %s) to bson data", serverErrorMessages[seInputBSONNotValid], teacher.PID.Hex())
		return err
	}
	err = bson.Unmarshal(teacherBSONData, &updateBSONDocument)
	if err != nil {
		err = fmt.Errorf("[%s] - could not convert teacher (PID %s) to bson document", serverErrorMessages[seInputBSONNotValid], teacher.PID.Hex())
		return err
	}
	var updateOptions = bson.D{{"$set", updateBSONDocument}}

	insertResult, err := dbPool.Collection(DBCollectionTeacher).UpdateOne(context.TODO(), updateFilter, updateOptions)
	if err != nil {
		err = fmt.Errorf("[%s] - %s", serverErrorMessages[seDBResourceQuery], err.Error())
		return err
	}

	logging.Debugmf(logModTeacherMgmt, "Update teacher (PID %s): matched %d modified %d",
		teacher.PID.Hex(), insertResult.MatchedCount, insertResult.ModifiedCount)
	if insertResult.MatchedCount == 0 {
		err = fmt.Errorf("[%s] - could not find teacher (PID %s)", serverErrorMessages[seResourceNotFound], teacher.PID.Hex())
		return err
	} else if insertResult.ModifiedCount == 0 {
		err = fmt.Errorf("[%s] - teacher (PID %s) not changed", serverErrorMessages[seResourceNotChange], teacher.PID.Hex())
		return err
	}
	return nil
}

// delete teacher, return #delete entries, error
func deleteTeacher(pid primitive.ObjectID) (int, error) {
	var err error
	defer func() {
		if err != nil {
			logging.Errormf(logModTeacherMgmt, err.Error())
		}
	}()

	var deleteFilter bson.D
	if pid.IsZero() {
		deleteFilter = bson.D{{}}
	} else {
		deleteFilter = bson.D{{"_id", pid}}
	}

	deleteResult, err := dbPool.Collection(DBCollectionTeacher).DeleteMany(context.TODO(), deleteFilter)
	if err != nil {
		err = fmt.Errorf("[%s] - %s", serverErrorMessages[seDBResourceQuery], err.Error())
		return 0, err
	}

	logging.Debugmf(logModTeacherMgmt, "Deleted %d teacher results from DB", deleteResult.DeletedCount)
	return int(deleteResult.DeletedCount), nil
}
