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

var storyConfigHandlerTable = map[string]gin.HandlerFunc{
	"get":    storyGetHandler,
	"post":   storyPostHandler,
	"put":    storyPutHandler,
	"delete": storyDeleteHandler,
}

func storyGetHandler(ctx *gin.Context) {
	params := ginContextRequestParameter(ctx)
	response := GinResponse{
		Status: http.StatusOK,
	}
	defer func() {
		ginContextProcessResponse(ctx, &response)
	}()

	var storys []*Story
	var err error
	var pid primitive.ObjectID
	var findFilter bson.M

	if params.PID == "all" {
		pid = primitive.NilObjectID
		if params.FKEY == "student_pid" {
			var fid primitive.ObjectID
			fid, err = primitive.ObjectIDFromHex(params.FID)
			if err != nil {
				response.Status = http.StatusBadRequest
				response.Message = fmt.Sprintf("[%s] - Please specifiy a valid FID (mongoDB ObjectID)", serverErrorMessages[seInputParamNotValid])
				return
			}
			findFilter = bson.M{"student_pid": fid}
		}
	} else {
		pid, err = primitive.ObjectIDFromHex(params.PID)
		if err != nil {
			response.Status = http.StatusBadRequest
			response.Message = fmt.Sprintf("[%s] - Please specifiy a valid PID (mongoDB ObjectID)", serverErrorMessages[seInputParamNotValid])
			return
		}
		findFilter = bson.M{"_id": pid}
	}

	// pid: nil objectid for all, others for specified one
	storys, err = findStory(pid, findFilter)
	if err != nil {
		response.Status = http.StatusConflict
		response.Message = err.Error()
		return
	}
	response.Payload = storys
	return
}

func storyPostHandler(ctx *gin.Context) {
	params := ginContextRequestParameter(ctx)
	response := GinResponse{
		Status: http.StatusOK,
	}
	defer func() {
		ginContextProcessResponse(ctx, &response)
	}()

	var story Story
	var storyPID primitive.ObjectID
	var err error

	if err = json.Unmarshal(params.Data, &story); err != nil {
		response.Status = http.StatusBadRequest
		response.Message = fmt.Sprintf("[%s] - %s", serverErrorMessages[seInputJSONNotValid], err.Error())
		return
	}

	storyPID, err = createStory(&story)
	if err != nil {
		response.Status = http.StatusConflict
		response.Message = err.Error()
	} else {
		response.Payload = storyPID
	}
	return
}

func storyPutHandler(ctx *gin.Context) {
	params := ginContextRequestParameter(ctx)
	response := GinResponse{
		Status: http.StatusOK,
	}
	defer func() {
		ginContextProcessResponse(ctx, &response)
	}()

	var story Story
	var err error

	if err = json.Unmarshal(params.Data, &story); err != nil {
		response.Status = http.StatusBadRequest
		response.Message = fmt.Sprintf("[%s] - %s", serverErrorMessages[seInputJSONNotValid], err.Error())
		return
	}

	err = updateStory(&story)
	if err != nil {
		response.Status = http.StatusConflict
		response.Message = err.Error()
	} else {
		response.Payload = story.PID
	}
	return
}

func storyDeleteHandler(ctx *gin.Context) {
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
	deletedRows, err = deleteStory(pid)
	if err != nil {
		response.Status = http.StatusConflict
		response.Message = err.Error()
		return
	}
	response.Payload = deletedRows
	return
}

// find story, return story slice, error
func findStory(pid primitive.ObjectID, findFilter bson.M) ([]*Story, error) {
	var err error
	defer func() {
		if err != nil {
			logging.Errormf(logModStoryMgmt, err.Error())
		}
	}()

	var findOptions = options.Find()
	if pid.IsZero() {
	} else {
		findOptions.SetLimit(1)
	}

	findStory, err := dbPool.Collection(DBCollectionStory).Find(context.TODO(), findFilter, findOptions)
	if err != nil {
		err = fmt.Errorf("[%s] - %s", serverErrorMessages[seDBResourceQuery], err.Error())
		return nil, err
	}

	storys := []*Story{}
	for findStory.Next(context.TODO()) {
		var story Story
		err = findStory.Decode(&story)
		if err != nil {
			err = fmt.Errorf("[%s] - %s", serverErrorMessages[seDBResourceQuery], err.Error())
			return nil, err
		}
		storys = append(storys, &story)
	}

	err = findStory.Err()
	if err != nil {
		err = fmt.Errorf("[%s] - %s", serverErrorMessages[seDBResourceQuery], err.Error())
		return nil, err
	}

	logging.Debugmf(logModStoryMgmt, "Found %d story results from DB (PID=%v)", len(storys), pid.Hex())
	return storys, nil
}

// create story, return PID, error
func createStory(story *Story) (primitive.ObjectID, error) {
	var err error
	defer func() {
		if err != nil {
			logging.Errormf(logModStoryMgmt, err.Error())
		}
	}()

	// student PID check
	if story.StudentPID.IsZero() {
		err = fmt.Errorf("[%s] - No student PID specified", serverErrorMessages[seResourceNotFound])
		return primitive.NilObjectID, err
	}

	var findFilter bson.M
	students, err := findStudent(story.StudentPID, findFilter)
	if err != nil || len(students) == 0 {
		err = fmt.Errorf("[%s] - No associate student found with PID %s", serverErrorMessages[seResourceNotFound], story.StudentPID.Hex())
		return primitive.NilObjectID, err
	}

	insertResult, err := dbPool.Collection(DBCollectionStory).InsertOne(context.TODO(), story)
	if err != nil {
		err = fmt.Errorf("[%s] - %s", serverErrorMessages[seDBResourceQuery], err.Error())
		return primitive.NilObjectID, err
	}

	lastInsertID := insertResult.InsertedID.(primitive.ObjectID)
	logging.Debugmf(logModStoryMgmt, "Created story in DB (LastInsertID,PID=%s)", lastInsertID.Hex())
	return lastInsertID, nil
}

// update story, return error
func updateStory(story *Story) error {
	var err error
	defer func() {
		if err != nil {
			logging.Errormf(logModStoryMgmt, err.Error())
		}
	}()

	// story PID check
	if story.PID.IsZero() {
		err = fmt.Errorf("[%s] - story PID is empty", serverErrorMessages[seInputJSONNotValid])
		return err
	}

	// student PID check
	if story.StudentPID.IsZero() {
		err = fmt.Errorf("[%s] - No student PID specified", serverErrorMessages[seResourceNotFound])
		return err
	}

	var findFilter bson.M
	students, err := findStudent(story.StudentPID, findFilter)
	if err != nil || len(students) == 0 {
		err = fmt.Errorf("[%s] - No associate student found with PID %s", serverErrorMessages[seResourceNotFound], story.StudentPID.Hex())
		return err
	}

	var updateFilter = bson.D{{"_id", story.PID}}
	var updateBSONDocument = bson.D{}
	storyBSONData, err := bson.Marshal(story)
	if err != nil {
		err = fmt.Errorf("[%s] - could not convert story (PID %s) to bson data", serverErrorMessages[seInputBSONNotValid], story.PID.Hex())
		return err
	}
	err = bson.Unmarshal(storyBSONData, &updateBSONDocument)
	if err != nil {
		err = fmt.Errorf("[%s] - could not convert story (PID %s) to bson document", serverErrorMessages[seInputBSONNotValid], story.PID.Hex())
		return err
	}
	var updateOptions = bson.D{{"$set", updateBSONDocument}}

	insertResult, err := dbPool.Collection(DBCollectionStory).UpdateOne(context.TODO(), updateFilter, updateOptions)
	if err != nil {
		err = fmt.Errorf("[%s] - %s", serverErrorMessages[seDBResourceQuery], err.Error())
		return err
	}

	logging.Debugmf(logModStoryMgmt, "Update story (PID %s): matched %d modified %d",
		story.PID.Hex(), insertResult.MatchedCount, insertResult.ModifiedCount)
	if insertResult.MatchedCount == 0 {
		err = fmt.Errorf("[%s] - could not find story (PID %s)", serverErrorMessages[seResourceNotFound], story.PID.Hex())
		return err
	} else if insertResult.ModifiedCount == 0 {
		err = fmt.Errorf("[%s] - story (PID %s) not changed", serverErrorMessages[seResourceNotChange], story.PID.Hex())
		return err
	}
	return nil
}

// delete story, return #delete entries, error
func deleteStory(pid primitive.ObjectID) (int, error) {
	var err error
	defer func() {
		if err != nil {
			logging.Errormf(logModStoryMgmt, err.Error())
		}
	}()

	var deleteFilter bson.D = bson.D{}

	if !pid.IsZero() {
		deleteFilter = append(deleteFilter, bson.E{"_id", pid})
	}

	deleteResult, err := dbPool.Collection(DBCollectionStory).DeleteMany(context.TODO(), deleteFilter)
	if err != nil {
		err = fmt.Errorf("[%s] - %s", serverErrorMessages[seDBResourceQuery], err.Error())
		return 0, err
	}

	logging.Debugmf(logModStoryMgmt, "Deleted %d story results from DB", deleteResult.DeletedCount)
	return int(deleteResult.DeletedCount), nil
}
