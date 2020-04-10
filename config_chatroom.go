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

var chatroomConfigHandlerTable = map[string]gin.HandlerFunc{
	"get":    chatroomGetHandler,
	"post":   chatroomPostHandler,
	"put":    chatroomPutHandler,
	"delete": chatroomDeleteHandler,
}

func chatroomGetHandler(ctx *gin.Context) {
	params := ginContextRequestParameter(ctx)
	response := GinResponse{
		Status: http.StatusOK,
	}
	defer func() {
		ginContextProcessResponse(ctx, &response)
	}()

	var users []*ChatRoom
	var err error
	var pid primitive.ObjectID
	var findFilter bson.M

	if params.PID == "all" {
		pid = primitive.NilObjectID
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
	users, err = findRoom(pid, findFilter)
	if err != nil {
		response.Status = http.StatusConflict
		response.Message = err.Error()
		return
	}
	response.Payload = users
	return
}

func chatroomPostHandler(ctx *gin.Context) {
	params := ginContextRequestParameter(ctx)
	response := GinResponse{
		Status: http.StatusOK,
	}
	defer func() {
		ginContextProcessResponse(ctx, &response)
	}()

	var user ChatRoom
	var userPID primitive.ObjectID
	var err error

	if err = json.Unmarshal(params.Data, &user); err != nil {
		response.Status = http.StatusBadRequest
		response.Message = fmt.Sprintf("[%s] - %s", serverErrorMessages[seInputJSONNotValid], err.Error())
		return
	}

	userPID, err = createRoom(&user)
	if err != nil {
		response.Status = http.StatusConflict
		response.Message = err.Error()
	} else {
		response.Payload = userPID
	}
	return
}

func chatroomPutHandler(ctx *gin.Context) {
	params := ginContextRequestParameter(ctx)
	response := GinResponse{
		Status: http.StatusOK,
	}
	defer func() {
		ginContextProcessResponse(ctx, &response)
	}()

	var user ChatRoom
	var err error

	if err = json.Unmarshal(params.Data, &user); err != nil {
		response.Status = http.StatusBadRequest
		response.Message = fmt.Sprintf("[%s] - %s", serverErrorMessages[seInputJSONNotValid], err.Error())
		return
	}

	err = updateRoom(&user)
	if err != nil {
		response.Status = http.StatusConflict
		response.Message = err.Error()
	} else {
		response.Payload = user.PID
	}
	return
}

func chatroomDeleteHandler(ctx *gin.Context) {
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
	deletedRows, err = deleteRoom(pid)
	if err != nil {
		response.Status = http.StatusConflict
		response.Message = err.Error()
		return
	}
	response.Payload = deletedRows
	return
}

// find user, return user slice, error
func findRoom(pid primitive.ObjectID, findFilter bson.M) ([]*ChatRoom, error) {
	var err error
	defer func() {
		if err != nil {
			logging.Errormf(logModChatRoomMgmt, err.Error())
		}
	}()

	var findOptions = options.Find()
	//var findFilter bson.D
	if pid.IsZero() {
		//findFilter = bson.D{{}}
	} else {
		findOptions.SetLimit(1)
		findFilter = bson.M{"_id": pid}
	}

	findUser, err := dbPool.Collection(DBCollectionChatRoom).Find(context.TODO(), findFilter, findOptions)
	if err != nil {
		err = fmt.Errorf("[%s] - %s", serverErrorMessages[seDBResourceQuery], err.Error())
		return nil, err
	}

	users := []*ChatRoom{}
	for findUser.Next(context.TODO()) {
		var user ChatRoom
		err = findUser.Decode(&user)
		if err != nil {
			err = fmt.Errorf("[%s] - %s", serverErrorMessages[seDBResourceQuery], err.Error())
			return nil, err
		}
		users = append(users, &user)
	}

	err = findUser.Err()
	if err != nil {
		err = fmt.Errorf("[%s] - %s", serverErrorMessages[seDBResourceQuery], err.Error())
		return nil, err
	}

	logging.Debugmf(logModChatRoomMgmt, "Found %d user results from DB (PID=%v)", len(users), pid.Hex())
	return users, nil
}

// create user, return PID, error
func createRoom(user *ChatRoom) (primitive.ObjectID, error) {
	var err error
	defer func() {
		if err != nil {
			logging.Errormf(logModChatRoomMgmt, err.Error())
		}
	}()

	insertResult, err := dbPool.Collection(DBCollectionChatRoom).InsertOne(context.TODO(), user)
	if err != nil {
		err = fmt.Errorf("[%s] - %s", serverErrorMessages[seDBResourceQuery], err.Error())
		return primitive.NilObjectID, err
	}

	lastInsertID := insertResult.InsertedID.(primitive.ObjectID)
	logging.Debugmf(logModChatRoomMgmt, "Created user in DB (LastInsertID,PID=%s)", lastInsertID.Hex())

	// update user image name/url

	return lastInsertID, nil
}

// update user, return error
func updateRoom(user *ChatRoom) error {
	var err error
	defer func() {
		if err != nil {
			logging.Errormf(logModChatRoomMgmt, err.Error())
		}
	}()

	// user PID check
	if user.PID.IsZero() {
		err = fmt.Errorf("[%s] - user PID is empty", serverErrorMessages[seInputJSONNotValid])
		return err
	}

	// user image name/url

	// update user
	var updateFilter = bson.D{{"_id", user.PID}}
	var updateBSONDocument = bson.D{}
	userBSONData, err := bson.Marshal(user)
	if err != nil {
		err = fmt.Errorf("[%s] - could not convert user (PID %s) to bson data", serverErrorMessages[seInputBSONNotValid], user.PID.Hex())
		return err
	}
	err = bson.Unmarshal(userBSONData, &updateBSONDocument)
	if err != nil {
		err = fmt.Errorf("[%s] - could not convert user (PID %s) to bson document", serverErrorMessages[seInputBSONNotValid], user.PID.Hex())
		return err
	}
	var updateOptions = bson.D{{"$set", updateBSONDocument}}

	insertResult, err := dbPool.Collection(DBCollectionChatRoom).UpdateOne(context.TODO(), updateFilter, updateOptions)
	if err != nil {
		err = fmt.Errorf("[%s] - %s", serverErrorMessages[seDBResourceQuery], err.Error())
		return err
	}

	logging.Debugmf(logModChatRoomMgmt, "Update user (PID %s): matched %d modified %d",
		user.PID.Hex(), insertResult.MatchedCount, insertResult.ModifiedCount)
	if insertResult.MatchedCount == 0 {
		err = fmt.Errorf("[%s] - could not find user (PID %s)", serverErrorMessages[seResourceNotFound], user.PID.Hex())
		return err
	} else if insertResult.ModifiedCount == 0 {
		err = fmt.Errorf("[%s] - user (PID %s) not changed", serverErrorMessages[seResourceNotChange], user.PID.Hex())
		return err
	}
	return nil
}

// delete user, return #delete entries, error
func deleteRoom(pid primitive.ObjectID) (int, error) {
	var err error
	defer func() {
		if err != nil {
			logging.Errormf(logModChatRoomMgmt, err.Error())
		}
	}()

	var findFilter bson.M

	users, findErr := findRoom(pid, findFilter)
	if findErr != nil {
		err = fmt.Errorf("[%s] - could not delete user (PID %s) due to DB query/find error occurs", serverErrorMessages[seDBResourceQuery], pid.Hex())
		return 0, err
	}

	var deleteCnt int64
	for i := range users {
		deleteFilter := bson.D{{"_id", users[i].PID}}
		deleteResult, err := dbPool.Collection(DBCollectionChatRoom).DeleteMany(context.TODO(), deleteFilter)
		if err != nil {
			err = fmt.Errorf("[%s] - %s", serverErrorMessages[seDBResourceQuery], err.Error())
			return 0, err
		}
		deleteCnt += deleteResult.DeletedCount
	}
	logging.Debugmf(logModChatRoomMgmt, "Deleted %d user results from DB", deleteCnt)
	return int(deleteCnt), nil
}
