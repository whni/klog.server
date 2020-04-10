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

var registeruserConfigHandlerTable = map[string]gin.HandlerFunc{
	"get":    registeruserGetHandler,
	"post":   registeruserPostHandler,
	"put":    registeruserPutHandler,
	"delete": registeruserDeleteHandler,
}

func registeruserGetHandler(ctx *gin.Context) {
	params := ginContextRequestParameter(ctx)
	response := GinResponse{
		Status: http.StatusOK,
	}
	defer func() {
		ginContextProcessResponse(ctx, &response)
	}()

	var users []*RegisterUser
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
	users, err = findRegisterUser(pid, findFilter)
	if err != nil {
		response.Status = http.StatusConflict
		response.Message = err.Error()
		return
	}
	response.Payload = users
	return
}

func registeruserPostHandler(ctx *gin.Context) {
	params := ginContextRequestParameter(ctx)
	response := GinResponse{
		Status: http.StatusOK,
	}
	defer func() {
		ginContextProcessResponse(ctx, &response)
	}()

	var user RegisterUser
	var userPID primitive.ObjectID
	var err error

	if err = json.Unmarshal(params.Data, &user); err != nil {
		response.Status = http.StatusBadRequest
		response.Message = fmt.Sprintf("[%s] - %s", serverErrorMessages[seInputJSONNotValid], err.Error())
		return
	}

	userPID, err = createRegisterUser(&user)
	if err != nil {
		response.Status = http.StatusConflict
		response.Message = err.Error()
	} else {
		response.Payload = userPID
	}
	return
}

func registeruserPutHandler(ctx *gin.Context) {
	params := ginContextRequestParameter(ctx)
	response := GinResponse{
		Status: http.StatusOK,
	}
	defer func() {
		ginContextProcessResponse(ctx, &response)
	}()

	var user RegisterUser
	var err error

	if err = json.Unmarshal(params.Data, &user); err != nil {
		response.Status = http.StatusBadRequest
		response.Message = fmt.Sprintf("[%s] - %s", serverErrorMessages[seInputJSONNotValid], err.Error())
		return
	}

	err = updateRegisterUser(&user)
	if err != nil {
		response.Status = http.StatusConflict
		response.Message = err.Error()
	} else {
		response.Payload = user.PID
	}
	return
}

func registeruserDeleteHandler(ctx *gin.Context) {
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
	deletedRows, err = deleteRegisterUser(pid)
	if err != nil {
		response.Status = http.StatusConflict
		response.Message = err.Error()
		return
	}
	response.Payload = deletedRows
	return
}

// find user, return user slice, error
func findRegisterUser(pid primitive.ObjectID, findFilter bson.M) ([]*RegisterUser, error) {
	var err error
	defer func() {
		if err != nil {
			logging.Errormf(logModRegisterUserMgmt, err.Error())
		}
	}()

	var findOptions = options.Find()
	//var findFilter bson.D
	if pid.IsZero() {

	} else {
		findOptions.SetLimit(1)
		findFilter = bson.M{"_id": pid}
	}

	findUser, err := dbPool.Collection(DBCollectionRegisterUser).Find(context.TODO(), findFilter, findOptions)
	if err != nil {
		err = fmt.Errorf("[%s] - %s", serverErrorMessages[seDBResourceQuery], err.Error())
		return nil, err
	}

	users := []*RegisterUser{}
	for findUser.Next(context.TODO()) {
		var user RegisterUser
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

	logging.Debugmf(logModRegisterUserMgmt, "Found %d register user results from DB (PID=%v)", len(users), pid.Hex())
	return users, nil
}

// create user, return PID, error
func createRegisterUser(user *RegisterUser) (primitive.ObjectID, error) {
	var err error
	defer func() {
		if err != nil {
			logging.Errormf(logModRegisterUserMgmt, err.Error())
		}
	}()
	logging.Debugmf(logModRegisterUserMgmt, "%+v\n", user)
	insertResult, err := dbPool.Collection(DBCollectionRegisterUser).InsertOne(context.TODO(), user)
	if err != nil {
		err = fmt.Errorf("[%s] - %s", serverErrorMessages[seDBResourceQuery], err.Error())
		return primitive.NilObjectID, err
	}

	lastInsertID := insertResult.InsertedID.(primitive.ObjectID)
	logging.Debugmf(logModRegisterUserMgmt, "Created register user in DB (LastInsertID,PID=%s)", lastInsertID.Hex())

	return lastInsertID, nil
}

// update user, return error
func updateRegisterUser(user *RegisterUser) error {
	var err error
	defer func() {
		if err != nil {
			logging.Errormf(logModRegisterUserMgmt, err.Error())
		}
	}()

	// user PID check
	if user.PID.IsZero() {
		err = fmt.Errorf("[%s] - user PID is empty", serverErrorMessages[seInputJSONNotValid])
		return err
	}

	// update user
	var updateFilter = bson.D{{"_id", user.PID}}
	var updateBSONDocument = bson.D{}
	userBSONData, err := bson.Marshal(user)
	if err != nil {
		err = fmt.Errorf("[%s] - could not convert register user (PID %s) to bson data", serverErrorMessages[seInputBSONNotValid], user.PID.Hex())
		return err
	}
	err = bson.Unmarshal(userBSONData, &updateBSONDocument)
	if err != nil {
		err = fmt.Errorf("[%s] - could not convert register user (PID %s) to bson document", serverErrorMessages[seInputBSONNotValid], user.PID.Hex())
		return err
	}
	var updateOptions = bson.D{{"$set", updateBSONDocument}}

	insertResult, err := dbPool.Collection(DBCollectionRegisterUser).UpdateOne(context.TODO(), updateFilter, updateOptions)
	if err != nil {
		err = fmt.Errorf("[%s] - %s", serverErrorMessages[seDBResourceQuery], err.Error())
		return err
	}

	logging.Debugmf(logModRegisterUserMgmt, "Update register user (PID %s): matched %d modified %d",
		user.PID.Hex(), insertResult.MatchedCount, insertResult.ModifiedCount)
	if insertResult.MatchedCount == 0 {
		err = fmt.Errorf("[%s] - could not find register user (PID %s)", serverErrorMessages[seResourceNotFound], user.PID.Hex())
		return err
	} else if insertResult.ModifiedCount == 0 {
		err = fmt.Errorf("[%s] - register user (PID %s) not changed", serverErrorMessages[seResourceNotChange], user.PID.Hex())
		return err
	}
	return nil
}

// delete user, return #delete entries, error
func deleteRegisterUser(pid primitive.ObjectID) (int, error) {
	var err error
	defer func() {
		if err != nil {
			logging.Errormf(logModRegisterUserMgmt, err.Error())
		}
	}()

	var findFilter bson.M

	users, findErr := findRegisterUser(pid, findFilter)
	if findErr != nil {
		err = fmt.Errorf("[%s] - could not delete register user (PID %s) due to DB query/find error occurs", serverErrorMessages[seDBResourceQuery], pid.Hex())
		return 0, err
	}

	var deleteCnt int64
	for i := range users {
		deleteFilter := bson.D{{"_id", users[i].PID}}
		deleteResult, err := dbPool.Collection(DBCollectionRegisterUser).DeleteMany(context.TODO(), deleteFilter)
		if err != nil {
			err = fmt.Errorf("[%s] - %s", serverErrorMessages[seDBResourceQuery], err.Error())
			return 0, err
		}
		deleteCnt += deleteResult.DeletedCount
	}
	logging.Debugmf(logModRegisterUserMgmt, "Deleted %d register user results from DB", deleteCnt)
	return int(deleteCnt), nil
}

func findRegisteruserByID(userEmail string) (*RegisterUser, error) {
	var err error
	defer func() {
		if err != nil {
			logging.Errormf(logModRegisterUserMgmt, err.Error())
		}
	}()

	var teacher RegisterUser
	findFilter := bson.D{{"user_email", userEmail}}
	err = dbPool.Collection(DBCollectionRegisterUser).FindOne(context.TODO(), findFilter).Decode(&teacher)
	if err != nil {
		err = fmt.Errorf("[%s] - %s", serverErrorMessages[seDBResourceQuery], err.Error())
		return nil, err
	}

	logging.Debugmf(logModRegisterUserMgmt, "Found registeruser from DB (userEmail=%s)", userEmail)
	return &teacher, nil
}
