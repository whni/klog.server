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

var cloudMediaConfigHandlerTable = map[string]gin.HandlerFunc{
	"get":  cloudMediaGetHandler,
	"post": cloudMediaPostHandler,
	// "put":    cloudMediaPutHandler,
	// "delete": cloudMediaDeleteHandler,
}

func cloudMediaGetHandler(ctx *gin.Context) {
	params := ginContextRequestParameter(ctx)
	response := GinResponse{
		Status: http.StatusOK,
	}
	defer func() {
		ginContextProcessResponse(ctx, &response)
	}()

	var cloudMediaSlice []*CloudMedia
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
	cloudMediaSlice, err = findCloudMedia(pid)
	if err != nil {
		response.Status = http.StatusConflict
		response.Message = err.Error()
		return
	}
	response.Payload = cloudMediaSlice
	return
}

func cloudMediaPostHandler(ctx *gin.Context) {
	params := ginContextRequestParameter(ctx)
	response := GinResponse{
		Status: http.StatusOK,
	}
	defer func() {
		ginContextProcessResponse(ctx, &response)
	}()

	var cloudMedia CloudMedia
	var cloudMediaPID primitive.ObjectID
	var err error

	if err = json.Unmarshal(params.Data, &cloudMedia); err != nil {
		response.Status = http.StatusBadRequest
		response.Message = fmt.Sprintf("[%s] - %s", serverErrorMessages[seInputJSONNotValid], err.Error())
		return
	}

	cloudMediaPID, err = createCloudMedia(&cloudMedia)
	if err != nil {
		response.Status = http.StatusConflict
		response.Message = err.Error()
	} else {
		response.Payload = cloudMedia
	}
	return
}

// func studentPutHandler(ctx *gin.Context) {
// 	params := ginContextRequestParameter(ctx)
// 	response := GinResponse{
// 		Status: http.StatusOK,
// 	}
// 	defer func() {
// 		ginContextProcessResponse(ctx, &response)
// 	}()

// 	var student Student
// 	var err error

// 	if err = json.Unmarshal(params.Data, &student); err != nil {
// 		response.Status = http.StatusBadRequest
// 		response.Message = fmt.Sprintf("[%s] - %s", serverErrorMessages[seInputJSONNotValid], err.Error())
// 		return
// 	}

// 	err = updateStudent(&student)
// 	if err != nil {
// 		response.Status = http.StatusConflict
// 		response.Message = err.Error()
// 	} else {
// 		response.Payload = student.PID
// 	}
// 	return
// }

// func studentDeleteHandler(ctx *gin.Context) {
// 	params := ginContextRequestParameter(ctx)
// 	response := GinResponse{
// 		Status: http.StatusOK,
// 	}
// 	defer func() {
// 		ginContextProcessResponse(ctx, &response)
// 	}()

// 	var err error
// 	var deletedRows int
// 	var pid primitive.ObjectID
// 	if params.PID == "all" {
// 		pid = primitive.NilObjectID
// 	} else {
// 		pid, err = primitive.ObjectIDFromHex(params.PID)
// 		if err != nil {
// 			response.Status = http.StatusBadRequest
// 			response.Message = fmt.Sprintf("[%s] - Please specifiy a valid PID (mongoDB ObjectID)", serverErrorMessages[seInputParamNotValid])
// 			return
// 		}
// 	}

// 	// pid: nil objectid for all, others for specified one
// 	deletedRows, err = deleteStudent(pid)
// 	if err != nil {
// 		response.Status = http.StatusConflict
// 		response.Message = err.Error()
// 		return
// 	}
// 	response.Payload = deletedRows
// 	return
// }

// find cloud media, return cloud media slice, error
func findCloudMedia(pid primitive.ObjectID) ([]*CloudMedia, error) {
	var err error
	defer func() {
		if err != nil {
			logging.Errormf(logModCloudMediaHandler, err.Error())
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

	findCursor, err := dbPool.Collection(DBCollectionCloudMedia).Find(context.TODO(), findFilter, findOptions)
	if err != nil {
		err = fmt.Errorf("[%s] - %s", serverErrorMessages[seDBResourceQuery], err.Error())
		return nil, err
	}

	cloudMediaSlice := []*CloudMedia{}
	for findCursor.Next(context.TODO()) {
		var cloudMedia CloudMedia
		err = findCursor.Decode(&cloudMedia)
		if err != nil {
			err = fmt.Errorf("[%s] - %s", serverErrorMessages[seDBResourceQuery], err.Error())
			return nil, err
		}
		cloudMediaSlice = append(cloudMediaSlice, &cloudMedia)
	}

	err = findCursor.Err()
	if err != nil {
		err = fmt.Errorf("[%s] - %s", serverErrorMessages[seDBResourceQuery], err.Error())
		return nil, err
	}

	logging.Debugmf(logModCloudMediaHandler, "Found %d cloud media results from DB (PID=%v)", len(cloudMediaSlice), pid)
	return cloudMediaSlice, nil
}

// create cloud media, return PID, error
func createCloudMedia(cloudMedia *CloudMedia) (primitive.ObjectID, error) {
	var err error
	defer func() {
		if err != nil {
			logging.Errormf(logModCloudMediaHandler, err.Error())
		}
	}()

	// teacher PID check
	if cloudMedia.StudentPID.IsZero() {
		err = fmt.Errorf("[%s] - No student PID associated", serverErrorMessages[seResourceNotFound])
		return primitive.NilObjectID, err
	}
	students, err := findStudent(cloudMedia.StudentPID)
	if err != nil || len(students) == 0 {
		err = fmt.Errorf("[%s] - No associate student found with PID %s", serverErrorMessages[seResourceNotFound], cloudMedia.StudentPID.Hex())
		return primitive.NilObjectID, err
	}

	// check if media uploaded to cloud already

	// insertResult, err := dbPool.Collection(DBCollectionCloudMedia).InsertOne(context.TODO(), student)
	// if err != nil {
	// 	err = fmt.Errorf("[%s] - %s", serverErrorMessages[seDBResourceQuery], err.Error())
	// 	return primitive.NilObjectID, err
	// }

	// lastInsertID := insertResult.InsertedID.(primitive.ObjectID)
	// logging.Debugmf(logModStudentHandler, "Created student in DB (LastInsertID,PID=%s)", lastInsertID.Hex())
	// return lastInsertID, nil
}

// // update student, return error
// func updateStudent(student *Student) error {
// 	var err error
// 	defer func() {
// 		if err != nil {
// 			logging.Errormf(logModStudentHandler, err.Error())
// 		}
// 	}()

// 	// student PID check
// 	if student.PID.IsZero() {
// 		err = fmt.Errorf("[%s] - student PID is empty", serverErrorMessages[seInputJSONNotValid])
// 		return err
// 	}

// 	// teacher PID check
// 	if student.TeacherPID.IsZero() {
// 		err = fmt.Errorf("[%s] - No teacher PID specified", serverErrorMessages[seResourceNotFound])
// 		return err
// 	}
// 	teachers, err := findTeacher(student.TeacherPID)
// 	if err != nil || len(teachers) == 0 {
// 		err = fmt.Errorf("[%s] - No teacher found with PID %s", serverErrorMessages[seResourceNotFound], student.TeacherPID.Hex())
// 		return err
// 	}

// 	var updateFilter = bson.D{{"_id", student.PID}}
// 	var updateBSONDocument = bson.D{}
// 	studentBSONData, err := bson.Marshal(student)
// 	if err != nil {
// 		err = fmt.Errorf("[%s] - could not convert student (PID %s) to bson data", serverErrorMessages[seInputBSONNotValid], student.PID.Hex())
// 		return err
// 	}
// 	err = bson.Unmarshal(studentBSONData, &updateBSONDocument)
// 	if err != nil {
// 		err = fmt.Errorf("[%s] - could not convert student (PID %s) to bson document", serverErrorMessages[seInputBSONNotValid], student.PID.Hex())
// 		return err
// 	}
// 	var updateOptions = bson.D{{"$set", updateBSONDocument}}

// 	insertResult, err := dbPool.Collection(DBCollectionStudent).UpdateOne(context.TODO(), updateFilter, updateOptions)
// 	if err != nil {
// 		err = fmt.Errorf("[%s] - %s", serverErrorMessages[seDBResourceQuery], err.Error())
// 		return err
// 	}

// 	logging.Debugmf(logModStudentHandler, "Update student (PID %s): matched %d modified %d",
// 		student.PID.Hex(), insertResult.MatchedCount, insertResult.ModifiedCount)
// 	if insertResult.MatchedCount == 0 {
// 		err = fmt.Errorf("[%s] - could not find student (PID %s)", serverErrorMessages[seResourceNotFound], student.PID.Hex())
// 		return err
// 	} else if insertResult.ModifiedCount == 0 {
// 		err = fmt.Errorf("[%s] - student (PID %s) not changed", serverErrorMessages[seResourceNotChange], student.PID.Hex())
// 		return err
// 	}
// 	return nil
// }

// // delete student, return #delete entries, error
// func deleteStudent(pid primitive.ObjectID) (int, error) {
// 	var err error
// 	defer func() {
// 		if err != nil {
// 			logging.Errormf(logModStudentHandler, err.Error())
// 		}
// 	}()

// 	var deleteFilter bson.D
// 	if pid.IsZero() {
// 		deleteFilter = bson.D{{}}
// 	} else {
// 		deleteFilter = bson.D{{"_id", pid}}
// 	}

// 	deleteResult, err := dbPool.Collection(DBCollectionStudent).DeleteMany(context.TODO(), deleteFilter)
// 	if err != nil {
// 		err = fmt.Errorf("[%s] - %s", serverErrorMessages[seDBResourceQuery], err.Error())
// 		return 0, err
// 	}

// 	logging.Debugmf(logModStudentHandler, "Deleted %d student results from DB", deleteResult.DeletedCount)
// 	return int(deleteResult.DeletedCount), nil
// }
