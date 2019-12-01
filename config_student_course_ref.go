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

var studentCourseRefConfigHandlerTable = map[string]gin.HandlerFunc{
	"get":    studentCourseRefGetHandler,
	"post":   studentCourseRefPostHandler,
	"put":    studentCourseRefPutHandler,
	"delete": studentCourseRefDeleteHandler,
}

func studentCourseRefGetHandler(ctx *gin.Context) {
	// params := ginContextRequestParameter(ctx)
	response := GinResponse{
		Status: http.StatusOK,
	}
	defer func() {
		ginContextProcessResponse(ctx, &response)
	}()

	var references []*StudentCourseRef
	var err error
	var studentPID primitive.ObjectID
	var coursePID primitive.ObjectID

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
	cPID := ctx.Request.URL.Query().Get("course_pid")
	if cPID == "all" {
		coursePID = primitive.NilObjectID
	} else {
		coursePID, err = primitive.ObjectIDFromHex(cPID)
		if err != nil {
			response.Status = http.StatusBadRequest
			response.Message = fmt.Sprintf("[%s] - Please specifiy a valid course PID (course_pid=?)", serverErrorMessages[seInputParamNotValid])
			return
		}
	}

	// pid: nil objectid for all, others for specified one
	references, err = findStudentCourseRef(studentPID, coursePID)
	if err != nil {
		response.Status = http.StatusConflict
		response.Message = err.Error()
		return
	}
	response.Payload = references
	return
}

func studentCourseRefPostHandler(ctx *gin.Context) {
	params := ginContextRequestParameter(ctx)
	response := GinResponse{
		Status: http.StatusOK,
	}
	defer func() {
		ginContextProcessResponse(ctx, &response)
	}()

	var reference StudentCourseRef
	var referencePID primitive.ObjectID
	var err error

	if err = json.Unmarshal(params.Data, &reference); err != nil {
		response.Status = http.StatusBadRequest
		response.Message = fmt.Sprintf("[%s] - %s", serverErrorMessages[seInputJSONNotValid], err.Error())
		return
	}

	referencePID, err = createStudentCourseRef(&reference)
	if err != nil {
		response.Status = http.StatusConflict
		response.Message = err.Error()
	} else {
		response.Payload = referencePID
	}
	return
}

func studentCourseRefPutHandler(ctx *gin.Context) {
	params := ginContextRequestParameter(ctx)
	response := GinResponse{
		Status: http.StatusOK,
	}
	defer func() {
		ginContextProcessResponse(ctx, &response)
	}()

	var reference StudentCourseRef
	var err error

	if err = json.Unmarshal(params.Data, &reference); err != nil {
		response.Status = http.StatusBadRequest
		response.Message = fmt.Sprintf("[%s] - %s", serverErrorMessages[seInputJSONNotValid], err.Error())
		return
	}

	err = updateStudentCourseRef(&reference)
	if err != nil {
		response.Status = http.StatusConflict
		response.Message = err.Error()
	} else {
		response.Payload = reference.PID
	}
	return
}

func studentCourseRefDeleteHandler(ctx *gin.Context) {
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
	var coursePID primitive.ObjectID

	sPID := ctx.Request.URL.Query().Get("student_pid")
	if sPID == "all" {
		studentPID = primitive.NilObjectID
	} else {
		studentPID, err = primitive.ObjectIDFromHex(sPID)
		if err != nil {
			response.Status = http.StatusBadRequest
			response.Message = fmt.Sprintf("[%s] - Please specifiy a valid student PID (pid=?)", serverErrorMessages[seInputParamNotValid])
			return
		}
	}
	cPID := ctx.Request.URL.Query().Get("course_pid")
	if cPID == "all" {
		coursePID = primitive.NilObjectID
	} else {
		coursePID, err = primitive.ObjectIDFromHex(cPID)
		if err != nil {
			response.Status = http.StatusBadRequest
			response.Message = fmt.Sprintf("[%s] - Please specifiy a valid course PID (sid=?)", serverErrorMessages[seInputParamNotValid])
			return
		}
	}

	// pid: nil objectid for all, others for specified one
	deletedRows, err = deleteStudentCourseRef(studentPID, coursePID)
	if err != nil {
		response.Status = http.StatusConflict
		response.Message = err.Error()
		return
	}
	response.Payload = deletedRows
	return
}

// find student-course references, return references slice, error
func findStudentCourseRef(studentPID primitive.ObjectID, coursePID primitive.ObjectID) ([]*StudentCourseRef, error) {
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
	if !coursePID.IsZero() {
		findFilter = append(findFilter, bson.E{"course_pid", coursePID})
	}

	findCursor, err := dbPool.Collection(DBCollectionStudentCourseRef).Find(context.TODO(), findFilter, findOptions)
	if err != nil {
		err = fmt.Errorf("[%s] - %s", serverErrorMessages[seDBResourceQuery], err.Error())
		return nil, err
	}

	references := []*StudentCourseRef{}
	for findCursor.Next(context.TODO()) {
		var reference StudentCourseRef
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

	logging.Debugmf(logModReferenceMgmt, "Found %d student-course reference results from DB (studentPID=%v, coursePID=%v)",
		len(references), studentPID.Hex(), coursePID.Hex())
	return references, nil
}

// create student-course reference, return PID, error
func createStudentCourseRef(reference *StudentCourseRef) (primitive.ObjectID, error) {
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

	// course PID check
	if reference.CoursePID.IsZero() {
		err = fmt.Errorf("[%s] - No course PID specified", serverErrorMessages[seResourceNotFound])
		return primitive.NilObjectID, err
	}
	courses, err := findCourse(reference.CoursePID)
	if err != nil || len(courses) == 0 {
		err = fmt.Errorf("[%s] - No courses found with PID %s", serverErrorMessages[seResourceNotFound], reference.CoursePID.Hex())
		return primitive.NilObjectID, err
	}

	insertResult, err := dbPool.Collection(DBCollectionStudentCourseRef).InsertOne(context.TODO(), reference)
	if err != nil {
		err = fmt.Errorf("[%s] - %s", serverErrorMessages[seDBResourceQuery], err.Error())
		return primitive.NilObjectID, err
	}

	lastInsertID := insertResult.InsertedID.(primitive.ObjectID)
	logging.Debugmf(logModReferenceMgmt, "Created student-course reference in DB (LastInsertID,PID=%s)", lastInsertID.Hex())

	return lastInsertID, nil
}

// update student-course reference, return error
func updateStudentCourseRef(reference *StudentCourseRef) error {
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

	// course PID check
	if reference.CoursePID.IsZero() {
		err = fmt.Errorf("[%s] - No course PID specified", serverErrorMessages[seResourceNotFound])
		return err
	}
	courses, err := findCourse(reference.CoursePID)
	if err != nil || len(courses) == 0 {
		err = fmt.Errorf("[%s] - No courses found with PID %s", serverErrorMessages[seResourceNotFound], reference.CoursePID.Hex())
		return err
	}

	// update student
	var updateFilter = bson.D{{"_id", reference.PID}}
	var updateBSONDocument = bson.D{}
	referenceBSONData, err := bson.Marshal(reference)
	if err != nil {
		err = fmt.Errorf("[%s] - could not convert student-course reference (PID %s) to bson data", serverErrorMessages[seInputBSONNotValid], reference.PID.Hex())
		return err
	}
	err = bson.Unmarshal(referenceBSONData, &updateBSONDocument)
	if err != nil {
		err = fmt.Errorf("[%s] - could not convert student-course reference (PID %s) to bson document", serverErrorMessages[seInputBSONNotValid], reference.PID.Hex())
		return err
	}
	var updateOptions = bson.D{{"$set", updateBSONDocument}}

	insertResult, err := dbPool.Collection(DBCollectionStudentCourseRef).UpdateOne(context.TODO(), updateFilter, updateOptions)
	if err != nil {
		err = fmt.Errorf("[%s] - %s", serverErrorMessages[seDBResourceQuery], err.Error())
		return err
	}

	logging.Debugmf(logModReferenceMgmt, "Update student-course reference (PID %s): matched %d modified %d",
		reference.PID.Hex(), insertResult.MatchedCount, insertResult.ModifiedCount)
	if insertResult.MatchedCount == 0 {
		err = fmt.Errorf("[%s] - could not find student-course reference (PID %s)", serverErrorMessages[seResourceNotFound], reference.PID.Hex())
		return err
	} else if insertResult.ModifiedCount == 0 {
		err = fmt.Errorf("[%s] - student-course reference (PID %s) not changed", serverErrorMessages[seResourceNotChange], reference.PID.Hex())
		return err
	}
	return nil
}

// delete student-course reference, return #delete entries, error
func deleteStudentCourseRef(studentPID primitive.ObjectID, coursePID primitive.ObjectID) (int, error) {
	var err error
	defer func() {
		if err != nil {
			logging.Errormf(logModReferenceMgmt, err.Error())
		}
	}()

	studentCourseReferences, findErr := findStudentCourseRef(studentPID, coursePID)
	if findErr != nil {
		err = fmt.Errorf("[%s] - could not delete student-course reference DB entries due to query error", serverErrorMessages[seResourceNotFound])
		return 0, err
	}

	var deleteCnt int64
	for i := range studentCourseReferences {
		_, deleteRecordErr := deleteCourseRecordByStudentPIDAndCoursePID(studentCourseReferences[i].StudentPID, studentCourseReferences[i].CoursePID)
		if deleteRecordErr != nil {
			err = fmt.Errorf("[%s] - stop deleting course-record reference (student PID %s course PID %s) since course record could not be deleted: %s",
				serverErrorMessages[seCloudOpsError], studentCourseReferences[i].StudentPID, studentCourseReferences[i].CoursePID, deleteRecordErr.Error())
			return int(deleteCnt), err
		}

		deleteFilter := bson.D{{"_id", studentCourseReferences[i].PID}}
		deleteResult, err := dbPool.Collection(DBCollectionStudentCourseRef).DeleteMany(context.TODO(), deleteFilter)
		if err != nil {
			err = fmt.Errorf("[%s] - %s", serverErrorMessages[seDBResourceQuery], err.Error())
			return 0, err
		}

		deleteCnt += deleteResult.DeletedCount
	}

	logging.Debugmf(logModReferenceMgmt, "Deleted %d student-course references from DB", deleteCnt)
	return int(deleteCnt), nil
}
