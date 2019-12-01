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

var courseCommentConfigHandlerTable = map[string]gin.HandlerFunc{
	"get":    courseCommentGetHandler,
	"post":   courseCommentPostHandler,
	"put":    courseCommentPutHandler,
	"delete": courseCommentDeleteHandler,
}

func courseCommentGetHandler(ctx *gin.Context) {
	params := ginContextRequestParameter(ctx)
	response := GinResponse{
		Status: http.StatusOK,
	}
	defer func() {
		ginContextProcessResponse(ctx, &response)
	}()

	var courseComments []*CourseComment
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
	courseComments, err = findCourseComment(pid)
	if err != nil {
		response.Status = http.StatusConflict
		response.Message = err.Error()
		return
	}
	response.Payload = courseComments
	return
}

func courseCommentPostHandler(ctx *gin.Context) {
	params := ginContextRequestParameter(ctx)
	response := GinResponse{
		Status: http.StatusOK,
	}
	defer func() {
		ginContextProcessResponse(ctx, &response)
	}()

	var courseComment CourseComment
	var courseCommentPID primitive.ObjectID
	var err error

	if err = json.Unmarshal(params.Data, &courseComment); err != nil {
		response.Status = http.StatusBadRequest
		response.Message = fmt.Sprintf("[%s] - %s", serverErrorMessages[seInputJSONNotValid], err.Error())
		return
	}

	courseCommentPID, err = createCourseComment(&courseComment)
	if err != nil {
		response.Status = http.StatusConflict
		response.Message = err.Error()
	} else {
		response.Payload = courseCommentPID
	}
	return
}

func courseCommentPutHandler(ctx *gin.Context) {
	params := ginContextRequestParameter(ctx)
	response := GinResponse{
		Status: http.StatusOK,
	}
	defer func() {
		ginContextProcessResponse(ctx, &response)
	}()

	var courseComment CourseComment
	var err error

	if err = json.Unmarshal(params.Data, &courseComment); err != nil {
		response.Status = http.StatusBadRequest
		response.Message = fmt.Sprintf("[%s] - %s", serverErrorMessages[seInputJSONNotValid], err.Error())
		return
	}

	err = updateCourseComment(&courseComment)
	if err != nil {
		response.Status = http.StatusConflict
		response.Message = err.Error()
	} else {
		response.Payload = courseComment.PID
	}
	return
}

func courseCommentDeleteHandler(ctx *gin.Context) {
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
	deletedRows, err = deleteCourseComment(pid)
	if err != nil {
		response.Status = http.StatusConflict
		response.Message = err.Error()
		return
	}
	response.Payload = deletedRows
	return
}

// find course comment, return course comment slice, error
func findCourseComment(pid primitive.ObjectID) ([]*CourseComment, error) {
	var err error
	defer func() {
		if err != nil {
			logging.Errormf(logModCourseCommentMgmt, err.Error())
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

	findCursor, err := dbPool.Collection(DBCollectionCourseComment).Find(context.TODO(), findFilter, findOptions)
	if err != nil {
		err = fmt.Errorf("[%s] - %s", serverErrorMessages[seDBResourceQuery], err.Error())
		return nil, err
	}

	courseComments := []*CourseComment{}
	for findCursor.Next(context.TODO()) {
		var courseComment CourseComment
		err = findCursor.Decode(&courseComment)
		if err != nil {
			err = fmt.Errorf("[%s] - %s", serverErrorMessages[seDBResourceQuery], err.Error())
			return nil, err
		}
		courseComments = append(courseComments, &courseComment)
	}

	err = findCursor.Err()
	if err != nil {
		err = fmt.Errorf("[%s] - %s", serverErrorMessages[seDBResourceQuery], err.Error())
		return nil, err
	}

	logging.Debugmf(logModCourseCommentMgmt, "Found %d course comments from DB (PID=%v)", len(courseComments), pid.Hex())
	return courseComments, nil
}

// create course comment, return PID, error
func createCourseComment(courseComment *CourseComment) (primitive.ObjectID, error) {
	var err error
	defer func() {
		if err != nil {
			logging.Errormf(logModCourseCommentMgmt, err.Error())
		}
	}()

	// course record PID check
	if courseComment.CourseRecordPID.IsZero() {
		err = fmt.Errorf("[%s] - No course record PID specified", serverErrorMessages[seResourceNotFound])
		return primitive.NilObjectID, err
	}
	courseRecords, err := findCourseRecord(courseComment.CourseRecordPID)
	if err != nil || len(courseRecords) == 0 {
		err = fmt.Errorf("[%s] - No associate course records found with PID %s", serverErrorMessages[seResourceNotFound], courseComment.CourseRecordPID.Hex())
		return primitive.NilObjectID, err
	}

	// person type check
	if courseComment.CommentPersonType == CommentPersonTypeTeacher {
		var teachers []*Teacher
		teachers, err = findTeacher(courseComment.CommentPersonPID)
		if err != nil || len(teachers) == 0 {
			err = fmt.Errorf("[%s] - No comment person (%s) found with PID %s", serverErrorMessages[seResourceNotFound],
				courseComment.CommentPersonType, courseComment.CourseRecordPID.Hex())
			return primitive.NilObjectID, err
		}
	} else if courseComment.CommentPersonType == CommentPersonTypeRelative {
		var relatives []*Relative
		relatives, err = findRelative(courseComment.CommentPersonPID)
		if err != nil || len(relatives) == 0 {
			err = fmt.Errorf("[%s] - No comment person (%s) found with PID %s", serverErrorMessages[seResourceNotFound],
				courseComment.CommentPersonType, courseComment.CourseRecordPID.Hex())
			return primitive.NilObjectID, err
		}
	} else {
		err = fmt.Errorf("[%s] - Comment person type must be %s or %s", serverErrorMessages[seResourceNotFound],
			CommentPersonTypeTeacher, CommentPersonTypeRelative)
		return primitive.NilObjectID, err
	}

	insertResult, err := dbPool.Collection(DBCollectionCourseComment).InsertOne(context.TODO(), courseComment)
	if err != nil {
		err = fmt.Errorf("[%s] - %s", serverErrorMessages[seDBResourceQuery], err.Error())
		return primitive.NilObjectID, err
	}

	lastInsertID := insertResult.InsertedID.(primitive.ObjectID)
	logging.Debugmf(logModCourseCommentMgmt, "Created course comment in DB (LastInsertID,PID=%s)", lastInsertID.Hex())
	return lastInsertID, nil
}

// update course comment, return error
func updateCourseComment(courseComment *CourseComment) error {
	var err error
	defer func() {
		if err != nil {
			logging.Errormf(logModCourseCommentMgmt, err.Error())
		}
	}()

	// course Comment PID check
	if courseComment.PID.IsZero() {
		err = fmt.Errorf("[%s] - course comment PID is empty", serverErrorMessages[seInputJSONNotValid])
		return err
	}

	// course record PID check
	if courseComment.CourseRecordPID.IsZero() {
		err = fmt.Errorf("[%s] - No course record PID specified", serverErrorMessages[seResourceNotFound])
		return err
	}
	courseRecords, err := findCourseRecord(courseComment.CourseRecordPID)
	if err != nil || len(courseRecords) == 0 {
		err = fmt.Errorf("[%s] - No associate course records found with PID %s", serverErrorMessages[seResourceNotFound], courseComment.CourseRecordPID.Hex())
		return err
	}

	// person type check
	if courseComment.CommentPersonType == CommentPersonTypeTeacher {
		var teachers []*Teacher
		teachers, err = findTeacher(courseComment.CommentPersonPID)
		if err != nil || len(teachers) == 0 {
			err = fmt.Errorf("[%s] - No comment person (%s) found with PID %s", serverErrorMessages[seResourceNotFound],
				courseComment.CommentPersonType, courseComment.CourseRecordPID.Hex())
			return err
		}
	} else if courseComment.CommentPersonType == CommentPersonTypeRelative {
		var relatives []*Relative
		relatives, err = findRelative(courseComment.CommentPersonPID)
		if err != nil || len(relatives) == 0 {
			err = fmt.Errorf("[%s] - No comment person (%s) found with PID %s", serverErrorMessages[seResourceNotFound],
				courseComment.CommentPersonType, courseComment.CourseRecordPID.Hex())
			return err
		}
	} else {
		err = fmt.Errorf("[%s] - Comment person type must be %s or %s", serverErrorMessages[seResourceNotFound],
			CommentPersonTypeTeacher, CommentPersonTypeRelative)
		return err
	}

	var updateFilter = bson.D{{"_id", courseComment.PID}}
	var updateBSONDocument = bson.D{}
	courseCommentBSONData, err := bson.Marshal(courseComment)
	if err != nil {
		err = fmt.Errorf("[%s] - could not convert course comment (PID %s) to bson data", serverErrorMessages[seInputBSONNotValid], courseComment.PID.Hex())
		return err
	}
	err = bson.Unmarshal(courseCommentBSONData, &updateBSONDocument)
	if err != nil {
		err = fmt.Errorf("[%s] - could not convert course comment (PID %s) to bson document", serverErrorMessages[seInputBSONNotValid], courseComment.PID.Hex())
		return err
	}
	var updateOptions = bson.D{{"$set", updateBSONDocument}}

	insertResult, err := dbPool.Collection(DBCollectionCourseComment).UpdateOne(context.TODO(), updateFilter, updateOptions)
	if err != nil {
		err = fmt.Errorf("[%s] - %s", serverErrorMessages[seDBResourceQuery], err.Error())
		return err
	}

	logging.Debugmf(logModCourseCommentMgmt, "Update course comment (PID %s): matched %d modified %d",
		courseComment.PID.Hex(), insertResult.MatchedCount, insertResult.ModifiedCount)
	if insertResult.MatchedCount == 0 {
		err = fmt.Errorf("[%s] - could not find course comment (PID %s)", serverErrorMessages[seResourceNotFound], courseComment.PID.Hex())
		return err
	} else if insertResult.ModifiedCount == 0 {
		err = fmt.Errorf("[%s] - course comment (PID %s) not changed", serverErrorMessages[seResourceNotChange], courseComment.PID.Hex())
		return err
	}
	return nil
}

// delete course comment, return #delete entries, error
func deleteCourseComment(pid primitive.ObjectID) (int, error) {
	var err error
	defer func() {
		if err != nil {
			logging.Errormf(logModCourseCommentMgmt, err.Error())
		}
	}()

	var deleteFilter bson.D
	if pid.IsZero() {
		deleteFilter = bson.D{{}}
	} else {
		deleteFilter = bson.D{{"_id", pid}}
	}

	deleteResult, err := dbPool.Collection(DBCollectionCourseComment).DeleteMany(context.TODO(), deleteFilter)
	if err != nil {
		err = fmt.Errorf("[%s] - %s", serverErrorMessages[seDBResourceQuery], err.Error())
		return 0, err
	}

	logging.Debugmf(logModInstituteMgmt, "Deleted %d course comments from DB", deleteResult.DeletedCount)
	return int(deleteResult.DeletedCount), nil
}
