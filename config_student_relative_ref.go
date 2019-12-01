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

var studentRelativeRefConfigHandlerTable = map[string]gin.HandlerFunc{
	"get":    studentRelativeRefGetHandler,
	"post":   studentRelativeRefPostHandler,
	"put":    studentRelativeRefPutHandler,
	"delete": studentRelativeRefDeleteHandler,
}

func studentRelativeRefGetHandler(ctx *gin.Context) {
	// params := ginContextRequestParameter(ctx)
	response := GinResponse{
		Status: http.StatusOK,
	}
	defer func() {
		ginContextProcessResponse(ctx, &response)
	}()

	var references []*StudentRelativeRef
	var err error
	var studentPID primitive.ObjectID
	var relativePID primitive.ObjectID

	sPID := ctx.Request.URL.Query().Get("student_pid")
	if sPID == "all" {
		studentPID = primitive.NilObjectID
	} else {
		studentPID, err = primitive.ObjectIDFromHex(sPID)
		if err != nil {
			response.Status = http.StatusBadRequest
			response.Message = fmt.Sprintf("[%s] - Please specifiy a valid student PID (student_pid=?)", serverErrorMessages[seInputParamNotValid])
			return
		}
	}
	rPID := ctx.Request.URL.Query().Get("relative_pid")
	if rPID == "all" {
		relativePID = primitive.NilObjectID
	} else {
		relativePID, err = primitive.ObjectIDFromHex(rPID)
		if err != nil {
			response.Status = http.StatusBadRequest
			response.Message = fmt.Sprintf("[%s] - Please specifiy a valid relative PID (relative_pid=?)", serverErrorMessages[seInputParamNotValid])
			return
		}
	}

	// pid: nil objectid for all, others for specified one
	references, err = findStudentRelativeRef(studentPID, relativePID)
	if err != nil {
		response.Status = http.StatusConflict
		response.Message = err.Error()
		return
	}
	response.Payload = references
	return
}

func studentRelativeRefPostHandler(ctx *gin.Context) {
	params := ginContextRequestParameter(ctx)
	response := GinResponse{
		Status: http.StatusOK,
	}
	defer func() {
		ginContextProcessResponse(ctx, &response)
	}()

	var reference StudentRelativeRef
	var referencePID primitive.ObjectID
	var err error

	if err = json.Unmarshal(params.Data, &reference); err != nil {
		response.Status = http.StatusBadRequest
		response.Message = fmt.Sprintf("[%s] - %s", serverErrorMessages[seInputJSONNotValid], err.Error())
		return
	}

	referencePID, err = createStudentRelativeRef(&reference)
	if err != nil {
		response.Status = http.StatusConflict
		response.Message = err.Error()
	} else {
		response.Payload = referencePID
	}
	return
}

func studentRelativeRefPutHandler(ctx *gin.Context) {
	params := ginContextRequestParameter(ctx)
	response := GinResponse{
		Status: http.StatusOK,
	}
	defer func() {
		ginContextProcessResponse(ctx, &response)
	}()

	var reference StudentRelativeRef
	var err error

	if err = json.Unmarshal(params.Data, &reference); err != nil {
		response.Status = http.StatusBadRequest
		response.Message = fmt.Sprintf("[%s] - %s", serverErrorMessages[seInputJSONNotValid], err.Error())
		return
	}

	err = updateStudentRelativeRef(&reference)
	if err != nil {
		response.Status = http.StatusConflict
		response.Message = err.Error()
	} else {
		response.Payload = reference.PID
	}
	return
}

func studentRelativeRefDeleteHandler(ctx *gin.Context) {
	// params := ginContextRequestParameter(ctx)
	response := GinResponse{
		Status: http.StatusOK,
	}
	defer func() {
		ginContextProcessResponse(ctx, &response)
	}()

	var err error
	var deletedRows int
	var studentPID primitive.ObjectID
	var relativePID primitive.ObjectID

	sPID := ctx.Request.URL.Query().Get("student_pid")
	if sPID == "all" {
		studentPID = primitive.NilObjectID
	} else {
		studentPID, err = primitive.ObjectIDFromHex(sPID)
		if err != nil {
			response.Status = http.StatusBadRequest
			response.Message = fmt.Sprintf("[%s] - Please specifiy a valid student PID (student_pid=?)", serverErrorMessages[seInputParamNotValid])
			return
		}
	}
	rPID := ctx.Request.URL.Query().Get("relative_pid")
	if rPID == "all" {
		relativePID = primitive.NilObjectID
	} else {
		relativePID, err = primitive.ObjectIDFromHex(rPID)
		if err != nil {
			response.Status = http.StatusBadRequest
			response.Message = fmt.Sprintf("[%s] - Please specifiy a valid relative PID (relative_pid=?)", serverErrorMessages[seInputParamNotValid])
			return
		}
	}

	// pid: nil objectid for all, others for specified one
	deletedRows, err = deleteStudentRelativeRef(studentPID, relativePID)
	if err != nil {
		response.Status = http.StatusConflict
		response.Message = err.Error()
		return
	}
	response.Payload = deletedRows
	return
}

// find student-relative references, return references slice, error
func findStudentRelativeRef(studentPID primitive.ObjectID, relativePID primitive.ObjectID) ([]*StudentRelativeRef, error) {
	var err error
	defer func() {
		if err != nil {
			logging.Errormf(logModReferenceMgmt, err.Error())
		}
	}()

	var findOptions = options.Find()
	var findFilter bson.D = bson.D{}
	if !studentPID.IsZero() {
		findFilter = append(findFilter, bson.E{"student_pid", studentPID})
	}
	if !relativePID.IsZero() {
		findFilter = append(findFilter, bson.E{"relative_pid", relativePID})
	}

	findCursor, err := dbPool.Collection(DBCollectionStudentRelativeRef).Find(context.TODO(), findFilter, findOptions)
	if err != nil {
		err = fmt.Errorf("[%s] - %s", serverErrorMessages[seDBResourceQuery], err.Error())
		return nil, err
	}

	references := []*StudentRelativeRef{}
	for findCursor.Next(context.TODO()) {
		var reference StudentRelativeRef
		err = findCursor.Decode(&reference)
		if err != nil {
			err = fmt.Errorf("[%s] - %s", serverErrorMessages[seDBResourceQuery], err.Error())
			return nil, err
		}
		references = append(references, &reference)
	}

	err = findCursor.Err()
	if err != nil {
		err = fmt.Errorf("[%s] - %s", serverErrorMessages[seDBResourceQuery], err.Error())
		return nil, err
	}

	logging.Debugmf(logModReferenceMgmt, "Found %d student-relative reference results from DB (studentPID=%v, relativePID=%v)",
		len(references), studentPID.Hex(), relativePID.Hex())
	return references, nil
}

// create student-relative reference, return PID, error
func createStudentRelativeRef(reference *StudentRelativeRef) (primitive.ObjectID, error) {
	var err error
	defer func() {
		if err != nil {
			logging.Errormf(logModReferenceMgmt, err.Error())
		}
	}()

	// student PID check
	if reference.StudentPID.IsZero() {
		err = fmt.Errorf("[%s] - No student PID specified", serverErrorMessages[seResourceNotFound])
		return primitive.NilObjectID, err
	}
	students, err := findStudent(reference.StudentPID)
	if err != nil || len(students) == 0 {
		err = fmt.Errorf("[%s] - No students found with PID %s", serverErrorMessages[seResourceNotFound], reference.StudentPID.Hex())
		return primitive.NilObjectID, err
	}

	// relative PID check
	if reference.RelativePID.IsZero() {
		err = fmt.Errorf("[%s] - No relative PID specified", serverErrorMessages[seResourceNotFound])
		return primitive.NilObjectID, err
	}
	relatives, err := findRelative(reference.RelativePID)
	if err != nil || len(relatives) == 0 {
		err = fmt.Errorf("[%s] - No relatives found with PID %s", serverErrorMessages[seResourceNotFound], reference.RelativePID.Hex())
		return primitive.NilObjectID, err
	}

	// only one relative must be main relationship
	studentReferences, err := findStudentRelativeRef(reference.StudentPID, primitive.NilObjectID)
	var mainReference *StudentRelativeRef = nil
	if err == nil && len(studentReferences) > 0 {
		for i := 0; i < len(studentReferences); i++ {
			if studentReferences[i].IsMain {
				mainReference = studentReferences[i]
				break
			}
		}
	}
	if mainReference != nil && reference.IsMain {
		err = fmt.Errorf("[%s] - Student (PID %v) already has a main relative (PID %v) -> could not add more", serverErrorMessages[seResourceConflict],
			mainReference.StudentPID.Hex(), mainReference.RelativePID.Hex())
		return primitive.NilObjectID, err
	}

	if mainReference == nil && !reference.IsMain {
		err = fmt.Errorf("[%s] - Student (PID %v) already has no main relative -> required at least one", serverErrorMessages[seResourceConflict],
			reference.StudentPID.Hex())
		return primitive.NilObjectID, err
	}

	insertResult, err := dbPool.Collection(DBCollectionStudentRelativeRef).InsertOne(context.TODO(), reference)
	if err != nil {
		err = fmt.Errorf("[%s] - %s", serverErrorMessages[seDBResourceQuery], err.Error())
		return primitive.NilObjectID, err
	}

	lastInsertID := insertResult.InsertedID.(primitive.ObjectID)
	logging.Debugmf(logModReferenceMgmt, "Created student-relative reference in DB (LastInsertID,PID=%s)", lastInsertID.Hex())

	return lastInsertID, nil
}

// update student-relative reference, return error
func updateStudentRelativeRef(reference *StudentRelativeRef) error {
	var err error
	defer func() {
		if err != nil {
			logging.Errormf(logModReferenceMgmt, err.Error())
		}
	}()

	// reference PID check
	if reference.PID.IsZero() {
		err = fmt.Errorf("[%s] - student-course reference PID is empty", serverErrorMessages[seInputJSONNotValid])
		return err
	}
	// student PID check
	if reference.StudentPID.IsZero() {
		err = fmt.Errorf("[%s] - No student PID specified", serverErrorMessages[seResourceNotFound])
		return err
	}
	students, err := findStudent(reference.StudentPID)
	if err != nil || len(students) == 0 {
		err = fmt.Errorf("[%s] - No students found with PID %s", serverErrorMessages[seResourceNotFound], reference.StudentPID.Hex())
		return err
	}

	// relative PID check
	if reference.RelativePID.IsZero() {
		err = fmt.Errorf("[%s] - No relative PID specified", serverErrorMessages[seResourceNotFound])
		return err
	}
	relatives, err := findRelative(reference.RelativePID)
	if err != nil || len(relatives) == 0 {
		err = fmt.Errorf("[%s] - No relatives found with PID %s", serverErrorMessages[seResourceNotFound], reference.RelativePID.Hex())
		return err
	}

	// only one relative can be main relationship
	studentReferences, err := findStudentRelativeRef(reference.StudentPID, primitive.NilObjectID)
	var mainReference *StudentRelativeRef = nil
	if err == nil && len(studentReferences) > 0 {
		for i := 0; i < len(studentReferences); i++ {
			if studentReferences[i].IsMain {
				mainReference = studentReferences[i]
				break
			}
		}
	}
	if mainReference != nil && reference.IsMain {
		err = fmt.Errorf("[%s] - Student (PID %v) already has a main relative (PID %v) -> could not add more", serverErrorMessages[seResourceConflict],
			mainReference.StudentPID.Hex(), mainReference.RelativePID.Hex())
		return err
	}

	if mainReference == nil && !reference.IsMain {
		err = fmt.Errorf("[%s] - Student (PID %v) already has no main relative -> required at least one", serverErrorMessages[seResourceConflict],
			reference.StudentPID.Hex())
		return err
	}

	// update student
	var updateFilter = bson.D{{"_id", reference.PID}}
	var updateBSONDocument = bson.D{}
	referenceBSONData, err := bson.Marshal(reference)
	if err != nil {
		err = fmt.Errorf("[%s] - could not convert student-relative reference (PID %s) to bson data", serverErrorMessages[seInputBSONNotValid], reference.PID.Hex())
		return err
	}
	err = bson.Unmarshal(referenceBSONData, &updateBSONDocument)
	if err != nil {
		err = fmt.Errorf("[%s] - could not convert student-relative reference (PID %s) to bson document", serverErrorMessages[seInputBSONNotValid], reference.PID.Hex())
		return err
	}
	var updateOptions = bson.D{{"$set", updateBSONDocument}}

	insertResult, err := dbPool.Collection(DBCollectionStudentRelativeRef).UpdateOne(context.TODO(), updateFilter, updateOptions)
	if err != nil {
		err = fmt.Errorf("[%s] - %s", serverErrorMessages[seDBResourceQuery], err.Error())
		return err
	}

	logging.Debugmf(logModReferenceMgmt, "Update student-relative reference (PID %s): matched %d modified %d",
		reference.PID.Hex(), insertResult.MatchedCount, insertResult.ModifiedCount)
	if insertResult.MatchedCount == 0 {
		err = fmt.Errorf("[%s] - could not find student-relative reference (PID %s)", serverErrorMessages[seResourceNotFound], reference.PID.Hex())
		return err
	} else if insertResult.ModifiedCount == 0 {
		err = fmt.Errorf("[%s] - student-relative reference (PID %s) not changed", serverErrorMessages[seResourceNotChange], reference.PID.Hex())
		return err
	}
	return nil
}

// delete student-relative reference, return #delete entries, error
func deleteStudentRelativeRef(studentPID primitive.ObjectID, relativePID primitive.ObjectID) (int, error) {
	var err error
	defer func() {
		if err != nil {
			logging.Errormf(logModReferenceMgmt, err.Error())
		}
	}()

	var deleteFilter bson.D = bson.D{}
	if !studentPID.IsZero() {
		deleteFilter = append(deleteFilter, bson.E{"student_pid", studentPID})
	}
	if !relativePID.IsZero() {
		deleteFilter = append(deleteFilter, bson.E{"relative_pid", relativePID})
	}

	deleteResult, err := dbPool.Collection(DBCollectionStudentRelativeRef).DeleteMany(context.TODO(), deleteFilter)
	if err != nil {
		err = fmt.Errorf("[%s] - %s", serverErrorMessages[seDBResourceQuery], err.Error())
		return 0, err
	}

	logging.Debugmf(logModReferenceMgmt, "Deleted %d student-relative references from DB", deleteResult.DeletedCount)
	return int(deleteResult.DeletedCount), nil
}
