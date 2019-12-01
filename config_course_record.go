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

var courseRecordConfigHandlerTable = map[string]gin.HandlerFunc{
	"get":    courseRecordGetHandler,
	"post":   courseRecordPostHandler,
	"put":    courseRecordPutHandler,
	"delete": courseRecordDeleteHandler,
}

func courseRecordGetHandler(ctx *gin.Context) {
	params := ginContextRequestParameter(ctx)
	response := GinResponse{
		Status: http.StatusOK,
	}
	defer func() {
		ginContextProcessResponse(ctx, &response)
	}()

	var courseRecords []*CourseRecord
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
	courseRecords, err = findCourseRecord(pid)
	if err != nil {
		response.Status = http.StatusConflict
		response.Message = err.Error()
		return
	}
	response.Payload = courseRecords
	return
}

func courseRecordPostHandler(ctx *gin.Context) {
	params := ginContextRequestParameter(ctx)
	response := GinResponse{
		Status: http.StatusOK,
	}
	defer func() {
		ginContextProcessResponse(ctx, &response)
	}()

	var courseRecord CourseRecord
	var courseRecordPID primitive.ObjectID
	var err error

	if err = json.Unmarshal(params.Data, &courseRecord); err != nil {
		response.Status = http.StatusBadRequest
		response.Message = fmt.Sprintf("[%s] - %s", serverErrorMessages[seInputJSONNotValid], err.Error())
		return
	}

	courseRecordPID, err = createCourseRecord(&courseRecord)
	if err != nil {
		response.Status = http.StatusConflict
		response.Message = err.Error()
	} else {
		response.Payload = courseRecordPID
	}
	return
}

func courseRecordPutHandler(ctx *gin.Context) {
	params := ginContextRequestParameter(ctx)
	response := GinResponse{
		Status: http.StatusOK,
	}
	defer func() {
		ginContextProcessResponse(ctx, &response)
	}()

	var courseRecord CourseRecord
	var err error

	if err = json.Unmarshal(params.Data, &courseRecord); err != nil {
		response.Status = http.StatusBadRequest
		response.Message = fmt.Sprintf("[%s] - %s", serverErrorMessages[seInputJSONNotValid], err.Error())
		return
	}

	err = updateCourseRecord(&courseRecord)
	if err != nil {
		response.Status = http.StatusConflict
		response.Message = err.Error()
	} else {
		response.Payload = courseRecord.PID
	}
	return
}

func courseRecordDeleteHandler(ctx *gin.Context) {
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
	deletedRows, err = deleteCourseRecord(pid)
	if err != nil {
		response.Status = http.StatusConflict
		response.Message = err.Error()
		return
	}
	response.Payload = deletedRows
	return
}

// find course record, return course record slice, error
func findCourseRecord(pid primitive.ObjectID) ([]*CourseRecord, error) {
	var err error
	defer func() {
		if err != nil {
			logging.Errormf(logModCourseRecordMgmt, err.Error())
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

	findCursor, err := dbPool.Collection(DBCollectionCourseRecord).Find(context.TODO(), findFilter, findOptions)
	if err != nil {
		err = fmt.Errorf("[%s] - %s", serverErrorMessages[seDBResourceQuery], err.Error())
		return nil, err
	}

	courseRecords := []*CourseRecord{}
	for findCursor.Next(context.TODO()) {
		var courseRecord CourseRecord
		err = findCursor.Decode(&courseRecord)
		if err != nil {
			err = fmt.Errorf("[%s] - %s", serverErrorMessages[seDBResourceQuery], err.Error())
			return nil, err
		}
		courseRecords = append(courseRecords, &courseRecord)
	}

	err = findCursor.Err()
	if err != nil {
		err = fmt.Errorf("[%s] - %s", serverErrorMessages[seDBResourceQuery], err.Error())
		return nil, err
	}

	logging.Debugmf(logModCourseRecordMgmt, "Found %d course records from DB (PID=%v)", len(courseRecords), pid.Hex())
	return courseRecords, nil
}

// find course record by student pid and course pid, return course record slice, error
func findCourseRecordByStudentPIDAndCoursePID(studentPID primitive.ObjectID, coursePID primitive.ObjectID) ([]*CourseRecord, error) {
	var err error
	defer func() {
		if err != nil {
			logging.Errormf(logModCourseRecordMgmt, err.Error())
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

	findCursor, err := dbPool.Collection(DBCollectionCourseRecord).Find(context.TODO(), findFilter, findOptions)
	if err != nil {
		err = fmt.Errorf("[%s] - %s", serverErrorMessages[seDBResourceQuery], err.Error())
		return nil, err
	}

	courseRecords := []*CourseRecord{}
	for findCursor.Next(context.TODO()) {
		var courseRecord CourseRecord
		err = findCursor.Decode(&courseRecord)
		if err != nil {
			err = fmt.Errorf("[%s] - %s", serverErrorMessages[seDBResourceQuery], err.Error())
			return nil, err
		}
		courseRecords = append(courseRecords, &courseRecord)
	}

	err = findCursor.Err()
	if err != nil {
		err = fmt.Errorf("[%s] - %s", serverErrorMessages[seDBResourceQuery], err.Error())
		return nil, err
	}

	logging.Debugmf(logModCourseRecordMgmt, "Found %d course records from DB (studentPID=%v, coursePID=%v)",
		len(courseRecords), studentPID.Hex(), coursePID.Hex())
	return courseRecords, nil
}

// create course record, return PID, error
func createCourseRecord(courseRecord *CourseRecord) (primitive.ObjectID, error) {
	var err error
	defer func() {
		if err != nil {
			logging.Errormf(logModCourseRecordMgmt, err.Error())
		}
	}()

	// student PID check
	if courseRecord.StudentPID.IsZero() {
		err = fmt.Errorf("[%s] - No student PID specified", serverErrorMessages[seResourceNotFound])
		return primitive.NilObjectID, err
	}

	// course PID check
	if courseRecord.CoursePID.IsZero() {
		err = fmt.Errorf("[%s] - No course PID specified", serverErrorMessages[seResourceNotFound])
		return primitive.NilObjectID, err
	}

	// student-course reference check
	studentCourseReferences, err := findStudentCourseRef(courseRecord.StudentPID, courseRecord.CoursePID)
	if err != nil || len(studentCourseReferences) == 0 {
		err = fmt.Errorf("[%s] - No student-course reference found with student PID %s and course PID %s -> cannot generate record",
			serverErrorMessages[seResourceNotFound], courseRecord.StudentPID.Hex(), courseRecord.CoursePID.Hex())
		return primitive.NilObjectID, err
	}

	// course target tags check
	courses, err := findCourse(courseRecord.CoursePID)
	if err != nil || len(courses) == 0 {
		err = fmt.Errorf("[%s] - No courses found with PID %s", serverErrorMessages[seResourceNotFound], courseRecord.CoursePID.Hex())
		return primitive.NilObjectID, err
	}
	courseTargets := courses[0].CourseTargets
	courseTargetTagValid := false
	for i := 0; i < len(courseTargets); i++ {
		if courseRecord.TargetTag == courseTargets[i].Tag {
			courseTargetTagValid = true
			break
		}
	}
	if courseTargetTagValid == false || courseRecord.TargetTag == "" {
		err = fmt.Errorf("[%s] - No valid course target tag (\"%v\") specified", serverErrorMessages[seResourceNotFound], courseRecord.TargetTag)
		return primitive.NilObjectID, err
	}

	insertResult, err := dbPool.Collection(DBCollectionCourseRecord).InsertOne(context.TODO(), courseRecord)
	if err != nil {
		err = fmt.Errorf("[%s] - %s", serverErrorMessages[seDBResourceQuery], err.Error())
		return primitive.NilObjectID, err
	}

	lastInsertID := insertResult.InsertedID.(primitive.ObjectID)
	logging.Debugmf(logModCourseRecordMgmt, "Created course record in DB (LastInsertID,PID=%s)", lastInsertID.Hex())

	return lastInsertID, nil
}

// update course record, return error
func updateCourseRecord(courseRecord *CourseRecord) error {
	var err error
	defer func() {
		if err != nil {
			logging.Errormf(logModCourseRecordMgmt, err.Error())
		}
	}()

	// course record PID check
	if courseRecord.PID.IsZero() {
		err = fmt.Errorf("[%s] - course record PID is empty", serverErrorMessages[seInputJSONNotValid])
		return err
	}

	// student PID check
	if courseRecord.StudentPID.IsZero() {
		err = fmt.Errorf("[%s] - No student PID specified", serverErrorMessages[seResourceNotFound])
		return err
	}

	// course PID check
	if courseRecord.CoursePID.IsZero() {
		err = fmt.Errorf("[%s] - No course PID specified", serverErrorMessages[seResourceNotFound])
		return err
	}

	// student-course reference check
	studentCourseReferences, err := findStudentCourseRef(courseRecord.StudentPID, courseRecord.CoursePID)
	if err != nil || len(studentCourseReferences) == 0 {
		err = fmt.Errorf("[%s] - No student-course reference found with student PID %s and course PID %s -> cannot generate record",
			serverErrorMessages[seResourceNotFound], courseRecord.StudentPID.Hex(), courseRecord.CoursePID.Hex())
		return err
	}

	// course target tags check
	courses, err := findCourse(courseRecord.CoursePID)
	if err != nil || len(courses) == 0 {
		err = fmt.Errorf("[%s] - No courses found with PID %s", serverErrorMessages[seResourceNotFound], courseRecord.CoursePID.Hex())
		return err
	}
	courseTargets := courses[0].CourseTargets
	courseTargetTagValid := false
	for i := 0; i < len(courseTargets); i++ {
		if courseRecord.TargetTag == courseTargets[i].Tag {
			courseTargetTagValid = true
			break
		}
	}
	if courseTargetTagValid == false || courseRecord.TargetTag == "" {
		err = fmt.Errorf("[%s] - No valid course target tag (\"%v\") specified", serverErrorMessages[seResourceNotFound], courseRecord.TargetTag)
		return err
	}

	// update course record
	var updateFilter = bson.D{{"_id", courseRecord.PID}}
	var updateBSONDocument = bson.D{}
	courseRecordBSONData, err := bson.Marshal(courseRecord)
	if err != nil {
		err = fmt.Errorf("[%s] - could not convert course record (PID %s) to bson data", serverErrorMessages[seInputBSONNotValid], courseRecord.PID.Hex())
		return err
	}
	err = bson.Unmarshal(courseRecordBSONData, &updateBSONDocument)
	if err != nil {
		err = fmt.Errorf("[%s] - could not convert course record (PID %s) to bson document", serverErrorMessages[seInputBSONNotValid], courseRecord.PID.Hex())
		return err
	}
	var updateOptions = bson.D{{"$set", updateBSONDocument}}

	insertResult, err := dbPool.Collection(DBCollectionCourseRecord).UpdateOne(context.TODO(), updateFilter, updateOptions)
	if err != nil {
		err = fmt.Errorf("[%s] - %s", serverErrorMessages[seDBResourceQuery], err.Error())
		return err
	}

	logging.Debugmf(logModCourseRecordMgmt, "Update course record (PID %s): matched %d modified %d",
		courseRecord.PID.Hex(), insertResult.MatchedCount, insertResult.ModifiedCount)
	if insertResult.MatchedCount == 0 {
		err = fmt.Errorf("[%s] - could not find course record (PID %s)", serverErrorMessages[seResourceNotFound], courseRecord.PID.Hex())
		return err
	} else if insertResult.ModifiedCount == 0 {
		err = fmt.Errorf("[%s] - course record (PID %s) not changed", serverErrorMessages[seResourceNotChange], courseRecord.PID.Hex())
		return err
	}
	return nil
}

// delete course record, return #delete entries, error
func deleteCourseRecord(pid primitive.ObjectID) (int, error) {
	var err error
	defer func() {
		if err != nil {
			logging.Errormf(logModCourseRecordMgmt, err.Error())
		}
	}()

	courseRecords, findErr := findCourseRecord(pid)
	if findErr != nil {
		err = fmt.Errorf("[%s] - could not delete course records (PID %s) due to DB query/find error occurs", serverErrorMessages[seDBResourceQuery], pid.Hex())
		return 0, err
	}

	var deleteCnt int64
	for i := range courseRecords {
		_, deleteCommentErr := deleteCourseCommentByRecordPID(courseRecords[i].PID)
		if deleteCommentErr != nil {
			err = fmt.Errorf("[%s] - stop deleting course record (PID %s) since course comment could not be deleted: %s",
				serverErrorMessages[seCloudOpsError], courseRecords[i].PID, deleteCommentErr.Error())
			return int(deleteCnt), err
		}
		_, deleteMediaErr := deleteCloudMediaByRecordPID(courseRecords[i].PID)
		if deleteMediaErr != nil {
			err = fmt.Errorf("[%s] - stop deleting course record (PID %s) since cloud media could not be deleted: %s",
				serverErrorMessages[seCloudOpsError], courseRecords[i].PID, deleteMediaErr.Error())
			return int(deleteCnt), err
		}

		deleteFilter := bson.D{{"_id", courseRecords[i].PID}}
		deleteResult, err := dbPool.Collection(DBCollectionCourseRecord).DeleteMany(context.TODO(), deleteFilter)
		if err != nil {
			err = fmt.Errorf("[%s] - %s", serverErrorMessages[seDBResourceQuery], err.Error())
			return 0, err
		}

		deleteCnt += deleteResult.DeletedCount
	}

	logging.Debugmf(logModCourseRecordMgmt, "Deleted %d course records from DB", deleteCnt)
	return int(deleteCnt), nil
}

// delete course record by student/course pid, return #delete entries, error
func deleteCourseRecordByStudentPIDAndCoursePID(studentPID primitive.ObjectID, coursePID primitive.ObjectID) (int, error) {
	var err error
	defer func() {
		if err != nil {
			logging.Errormf(logModCourseRecordMgmt, err.Error())
		}
	}()

	courseRecords, findErr := findCourseRecordByStudentPIDAndCoursePID(studentPID, coursePID)
	if findErr != nil {
		err = fmt.Errorf("[%s] - could not delete cloud media DB entries due to query error", serverErrorMessages[seResourceNotFound])
		return 0, err
	}

	var deleteCnt int64
	for i := range courseRecords {
		_, deleteCommentErr := deleteCourseCommentByRecordPID(courseRecords[i].PID)
		if deleteCommentErr != nil {
			err = fmt.Errorf("[%s] - stop deleting course record (student PID %s course PID %s) since course comment could not be deleted: %s",
				serverErrorMessages[seCloudOpsError], courseRecords[i].StudentPID, courseRecords[i].CoursePID, deleteCommentErr.Error())
			return int(deleteCnt), err
		}
		_, deleteMediaErr := deleteCloudMediaByRecordPID(courseRecords[i].PID)
		if deleteMediaErr != nil {
			err = fmt.Errorf("[%s] - stop deleting course record (student PID %s course PID %s) since cloud media could not be deleted: %s",
				serverErrorMessages[seCloudOpsError], courseRecords[i].StudentPID, courseRecords[i].CoursePID, deleteMediaErr.Error())
			return int(deleteCnt), err
		}

		deleteFilter := bson.D{{"_id", courseRecords[i].PID}}
		deleteResult, err := dbPool.Collection(DBCollectionCourseRecord).DeleteMany(context.TODO(), deleteFilter)
		if err != nil {
			err = fmt.Errorf("[%s] - %s", serverErrorMessages[seDBResourceQuery], err.Error())
			return 0, err
		}

		deleteCnt += deleteResult.DeletedCount
	}

	logging.Debugmf(logModCourseRecordMgmt, "Deleted %d course records from DB", deleteCnt)
	return int(deleteCnt), nil
}
