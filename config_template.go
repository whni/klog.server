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

var templateConfigHandlerTable = map[string]gin.HandlerFunc{
	"get":    templateGetHandler,
	"post":   templatePostHandler,
	"put":    templatePutHandler,
	"delete": templateDeleteHandler,
}

func templateGetHandler(ctx *gin.Context) {
	params := ginContextRequestParameter(ctx)
	response := GinResponse{
		Status: http.StatusOK,
	}
	defer func() {
		ginContextProcessResponse(ctx, &response)
	}()

	var templates []*Template
	var err error
	var pid primitive.ObjectID
	var findFilter bson.D
	findFilter = bson.D{{}}
	if params.PID == "all" {
		pid = primitive.NilObjectID
		if params.FKEY == "template_name" {
			findFilter = append(findFilter, bson.E{"template_name", params.FID})
		}
	} else {
		pid, err = primitive.ObjectIDFromHex(params.PID)
		if err != nil {
			response.Status = http.StatusBadRequest
			response.Message = fmt.Sprintf("[%s] - Please specifiy a valid PID (mongoDB ObjectID)", serverErrorMessages[seInputParamNotValid])
			return
		}
		findFilter = append(findFilter, bson.E{"_id", pid})
	}

	// pid: nil objectid for all, others for specified one
	templates, err = findTemplate(pid, findFilter)
	if err != nil {
		response.Status = http.StatusConflict
		response.Message = err.Error()
		return
	}
	response.Payload = templates
	return
}

func templatePostHandler(ctx *gin.Context) {
	params := ginContextRequestParameter(ctx)
	response := GinResponse{
		Status: http.StatusOK,
	}
	defer func() {
		ginContextProcessResponse(ctx, &response)
	}()

	var template Template
	var templatePID primitive.ObjectID
	var err error

	if err = json.Unmarshal(params.Data, &template); err != nil {
		response.Status = http.StatusBadRequest
		response.Message = fmt.Sprintf("[%s] - %s", serverErrorMessages[seInputJSONNotValid], err.Error())
		return
	}

	templatePID, err = createTemplate(&template)
	if err != nil {
		response.Status = http.StatusConflict
		response.Message = err.Error()
	} else {
		response.Payload = templatePID
	}
	return
}

func templatePutHandler(ctx *gin.Context) {
	params := ginContextRequestParameter(ctx)
	response := GinResponse{
		Status: http.StatusOK,
	}
	defer func() {
		ginContextProcessResponse(ctx, &response)
	}()

	var template Template
	var err error

	if err = json.Unmarshal(params.Data, &template); err != nil {
		response.Status = http.StatusBadRequest
		response.Message = fmt.Sprintf("[%s] - %s", serverErrorMessages[seInputJSONNotValid], err.Error())
		return
	}

	err = updateTemplate(&template)
	if err != nil {
		response.Status = http.StatusConflict
		response.Message = err.Error()
	} else {
		response.Payload = template.PID
	}
	return
}

func templateDeleteHandler(ctx *gin.Context) {
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
	deletedRows, err = deleteTemplate(pid)
	if err != nil {
		response.Status = http.StatusConflict
		response.Message = err.Error()
		return
	}
	response.Payload = deletedRows
	return
}

// find template, return template slice, error
func findTemplate(pid primitive.ObjectID, findFilter bson.D) ([]*Template, error) {
	var err error
	defer func() {
		if err != nil {
			logging.Errormf(logModTemplateMgmt, err.Error())
		}
	}()

	var findOptions = options.Find()
	if pid.IsZero() {
	} else {
		findOptions.SetLimit(1)
	}

	findTemplate, err := dbPool.Collection(DBCollectionTemplate).Find(context.TODO(), findFilter, findOptions)
	if err != nil {
		err = fmt.Errorf("[%s] - %s", serverErrorMessages[seDBResourceQuery], err.Error())
		return nil, err
	}

	templates := []*Template{}
	for findTemplate.Next(context.TODO()) {
		var template Template
		err = findTemplate.Decode(&template)
		if err != nil {
			err = fmt.Errorf("[%s] - %s", serverErrorMessages[seDBResourceQuery], err.Error())
			return nil, err
		}
		templates = append(templates, &template)
	}

	err = findTemplate.Err()
	if err != nil {
		err = fmt.Errorf("[%s] - %s", serverErrorMessages[seDBResourceQuery], err.Error())
		return nil, err
	}

	logging.Debugmf(logModTemplateMgmt, "Found %d template results from DB (PID=%v)", len(templates), pid.Hex())
	return templates, nil
}

// create template, return PID, error
func createTemplate(template *Template) (primitive.ObjectID, error) {
	var err error
	defer func() {
		if err != nil {
			logging.Errormf(logModTemplateMgmt, err.Error())
		}
	}()

	insertResult, err := dbPool.Collection(DBCollectionTemplate).InsertOne(context.TODO(), template)
	if err != nil {
		err = fmt.Errorf("[%s] - %s", serverErrorMessages[seDBResourceQuery], err.Error())
		return primitive.NilObjectID, err
	}

	lastInsertID := insertResult.InsertedID.(primitive.ObjectID)
	logging.Debugmf(logModTemplateMgmt, "Created template in DB (LastInsertID,PID=%s)", lastInsertID.Hex())
	return lastInsertID, nil
}

// update template, return error
func updateTemplate(template *Template) error {
	var err error
	defer func() {
		if err != nil {
			logging.Errormf(logModTemplateMgmt, err.Error())
		}
	}()

	// template PID check
	if template.PID.IsZero() {
		err = fmt.Errorf("[%s] - template PID is empty", serverErrorMessages[seInputJSONNotValid])
		return err
	}

	var updateFilter = bson.D{{"_id", template.PID}}
	var updateBSONDocument = bson.D{}
	templateBSONData, err := bson.Marshal(template)
	if err != nil {
		err = fmt.Errorf("[%s] - could not convert template (PID %s) to bson data", serverErrorMessages[seInputBSONNotValid], template.PID.Hex())
		return err
	}
	err = bson.Unmarshal(templateBSONData, &updateBSONDocument)
	if err != nil {
		err = fmt.Errorf("[%s] - could not convert template (PID %s) to bson document", serverErrorMessages[seInputBSONNotValid], template.PID.Hex())
		return err
	}
	var updateOptions = bson.D{{"$set", updateBSONDocument}}

	insertResult, err := dbPool.Collection(DBCollectionTemplate).UpdateOne(context.TODO(), updateFilter, updateOptions)
	if err != nil {
		err = fmt.Errorf("[%s] - %s", serverErrorMessages[seDBResourceQuery], err.Error())
		return err
	}

	logging.Debugmf(logModTemplateMgmt, "Update template (PID %s): matched %d modified %d",
		template.PID.Hex(), insertResult.MatchedCount, insertResult.ModifiedCount)
	if insertResult.MatchedCount == 0 {
		err = fmt.Errorf("[%s] - could not find template (PID %s)", serverErrorMessages[seResourceNotFound], template.PID.Hex())
		return err
	} else if insertResult.ModifiedCount == 0 {
		err = fmt.Errorf("[%s] - template (PID %s) not changed", serverErrorMessages[seResourceNotChange], template.PID.Hex())
		return err
	}
	return nil
}

// delete template, return #delete entries, error
func deleteTemplate(pid primitive.ObjectID) (int, error) {
	var err error
	defer func() {
		if err != nil {
			logging.Errormf(logModTemplateMgmt, err.Error())
		}
	}()

	var deleteFilter bson.D = bson.D{}

	if !pid.IsZero() {
		deleteFilter = append(deleteFilter, bson.E{"_id", pid})
	}

	deleteResult, err := dbPool.Collection(DBCollectionTemplate).DeleteMany(context.TODO(), deleteFilter)
	if err != nil {
		err = fmt.Errorf("[%s] - %s", serverErrorMessages[seDBResourceQuery], err.Error())
		return 0, err
	}

	logging.Debugmf(logModTemplateMgmt, "Deleted %d template results from DB", deleteResult.DeletedCount)
	return int(deleteResult.DeletedCount), nil
}
