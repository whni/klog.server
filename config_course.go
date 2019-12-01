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

var courseConfigHandlerTable = map[string]gin.HandlerFunc{
	"get":    courseGetHandler,
	"post":   coursePostHandler,
	"put":    coursePutHandler,
	"delete": courseDeleteHandler,
}

func courseGetHandler(ctx *gin.Context) {
	params := ginContextRequestParameter(ctx)
	response := GinResponse{
		Status: http.StatusOK,
	}
	defer func() {
		ginContextProcessResponse(ctx, &response)
	}()

	var courses []*Course
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
	courses, err = findCourse(pid)
	if err != nil {
		response.Status = http.StatusConflict
		response.Message = err.Error()
		return
	}
	response.Payload = courses
	return
}

func coursePostHandler(ctx *gin.Context) {
	params := ginContextRequestParameter(ctx)
	response := GinResponse{
		Status: http.StatusOK,
	}
	defer func() {
		ginContextProcessResponse(ctx, &response)
	}()

	var course Course
	var coursePID primitive.ObjectID
	var err error

	if err = json.Unmarshal(params.Data, &course); err != nil {
		response.Status = http.StatusBadRequest
		response.Message = fmt.Sprintf("[%s] - %s", serverErrorMessages[seInputJSONNotValid], err.Error())
		return
	}

	coursePID, err = createCourse(&course)
	if err != nil {
		response.Status = http.StatusConflict
		response.Message = err.Error()
	} else {
		response.Payload = coursePID
	}
	return
}

func coursePutHandler(ctx *gin.Context) {
	params := ginContextRequestParameter(ctx)
	response := GinResponse{
		Status: http.StatusOK,
	}
	defer func() {
		ginContextProcessResponse(ctx, &response)
	}()

	var course Course
	var err error

	if err = json.Unmarshal(params.Data, &course); err != nil {
		response.Status = http.StatusBadRequest
		response.Message = fmt.Sprintf("[%s] - %s", serverErrorMessages[seInputJSONNotValid], err.Error())
		return
	}

	err = updateCourse(&course)
	if err != nil {
		response.Status = http.StatusConflict
		response.Message = err.Error()
	} else {
		response.Payload = course.PID
	}
	return
}

func courseDeleteHandler(ctx *gin.Context) {
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
	deletedRows, err = deleteCourse(pid)
	if err != nil {
		response.Status = http.StatusConflict
		response.Message = err.Error()
		return
	}
	response.Payload = deletedRows
	return
}

// find course, return course slice, error
func findCourse(pid primitive.ObjectID) ([]*Course, error) {
	var err error
	defer func() {
		if err != nil {
			logging.Errormf(logModCourseMgmt, err.Error())
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

	findCursor, err := dbPool.Collection(DBCollectionCourse).Find(context.TODO(), findFilter, findOptions)
	if err != nil {
		err = fmt.Errorf("[%s] - %s", serverErrorMessages[seDBResourceQuery], err.Error())
		return nil, err
	}

	courses := []*Course{}
	for findCursor.Next(context.TODO()) {
		var course Course
		err = findCursor.Decode(&course)
		if err != nil {
			err = fmt.Errorf("[%s] - %s", serverErrorMessages[seDBResourceQuery], err.Error())
			return nil, err
		}
		courses = append(courses, &course)
	}

	err = findCursor.Err()
	if err != nil {
		err = fmt.Errorf("[%s] - %s", serverErrorMessages[seDBResourceQuery], err.Error())
		return nil, err
	}

	logging.Debugmf(logModCourseMgmt, "Found %d course results from DB (PID=%v)", len(courses), pid.Hex())
	return courses, nil
}

// find course by courseUID
func findCourseByUID(courseUID string) (*Course, error) {
	var err error
	defer func() {
		if err != nil {
			logging.Errormf(logModCourseMgmt, err.Error())
		}
	}()

	var course Course
	findFilter := bson.D{{"course_uid", courseUID}}
	err = dbPool.Collection(DBCollectionCourse).FindOne(context.TODO(), findFilter).Decode(&course)
	if err != nil {
		err = fmt.Errorf("[%s] - %s", serverErrorMessages[seDBResourceQuery], err.Error())
		return nil, err
	}

	logging.Debugmf(logModCourseMgmt, "Found course from DB (courseUID=%s)", courseUID)
	return &course, nil
}

// create course, return PID, error
func createCourse(course *Course) (primitive.ObjectID, error) {
	var err error
	defer func() {
		if err != nil {
			logging.Errormf(logModCourseMgmt, err.Error())
		}
	}()

	// institute PID check
	if course.InstitutePID.IsZero() {
		err = fmt.Errorf("[%s] - No institute PID specified", serverErrorMessages[seResourceNotFound])
		return primitive.NilObjectID, err
	}
	institutes, err := findInstitute(course.InstitutePID)
	if err != nil || len(institutes) == 0 {
		err = fmt.Errorf("[%s] - No institutes found with PID %s", serverErrorMessages[seResourceNotFound], course.InstitutePID.Hex())
		return primitive.NilObjectID, err
	}

	// teacher PID check
	if course.TeacherPID.IsZero() {
		err = fmt.Errorf("[%s] - No teacher PID specified", serverErrorMessages[seResourceNotFound])
		return primitive.NilObjectID, err
	}
	teachers, err := findTeacher(course.TeacherPID)
	if err != nil || len(teachers) == 0 {
		err = fmt.Errorf("[%s] - No teachers found with PID %s", serverErrorMessages[seResourceNotFound], course.TeacherPID.Hex())
		return primitive.NilObjectID, err
	}

	// assistant PID check
	if !course.AssistantPID.IsZero() {
		if course.AssistantPID.Hex() == course.TeacherPID.Hex() {
			err = fmt.Errorf("[%s] - Identical teacher/assistant PID %s", serverErrorMessages[seInputParamNotValid], course.AssistantPID.Hex())
			return primitive.NilObjectID, err
		}

		assistants, err := findTeacher(course.AssistantPID)
		if err != nil || len(assistants) == 0 {
			err = fmt.Errorf("[%s] - No assistants found with PID %s", serverErrorMessages[seResourceNotFound], course.AssistantPID.Hex())
			return primitive.NilObjectID, err
		}
	}

	insertResult, err := dbPool.Collection(DBCollectionCourse).InsertOne(context.TODO(), course)
	if err != nil {
		err = fmt.Errorf("[%s] - %s", serverErrorMessages[seDBResourceQuery], err.Error())
		return primitive.NilObjectID, err
	}

	lastInsertID := insertResult.InsertedID.(primitive.ObjectID)
	logging.Debugmf(logModCourseMgmt, "Created course in DB (LastInsertID,PID=%s)", lastInsertID.Hex())
	return lastInsertID, nil
}

// update course, return error
func updateCourse(course *Course) error {
	var err error
	defer func() {
		if err != nil {
			logging.Errormf(logModCourseMgmt, err.Error())
		}
	}()

	// course PID check
	if course.PID.IsZero() {
		err = fmt.Errorf("[%s] - course PID is empty", serverErrorMessages[seInputJSONNotValid])
		return err
	}

	// institute PID check
	if course.InstitutePID.IsZero() {
		err = fmt.Errorf("[%s] - No institute PID specified", serverErrorMessages[seResourceNotFound])
		return err
	}
	institutes, err := findInstitute(course.InstitutePID)
	if err != nil || len(institutes) == 0 {
		err = fmt.Errorf("[%s] - No institutes found with PID %s", serverErrorMessages[seResourceNotFound], course.InstitutePID.Hex())
		return err
	}

	// teacher PID check
	if course.TeacherPID.IsZero() {
		err = fmt.Errorf("[%s] - No teacher PID specified", serverErrorMessages[seResourceNotFound])
		return err
	}
	teachers, err := findTeacher(course.TeacherPID)
	if err != nil || len(teachers) == 0 {
		err = fmt.Errorf("[%s] - No teachers found with PID %s", serverErrorMessages[seResourceNotFound], course.TeacherPID.Hex())
		return err
	}

	// assistant PID check
	if !course.AssistantPID.IsZero() {
		if course.AssistantPID.Hex() == course.TeacherPID.Hex() {
			err = fmt.Errorf("[%s] - Identical teacher/assistant PID %s", serverErrorMessages[seInputParamNotValid], course.AssistantPID.Hex())
			return err
		}

		assistants, err := findTeacher(course.AssistantPID)
		if err != nil || len(assistants) == 0 {
			err = fmt.Errorf("[%s] - No assistants found with PID %s", serverErrorMessages[seResourceNotFound], course.AssistantPID.Hex())
			return err
		}
	}

	var updateFilter = bson.D{{"_id", course.PID}}
	var updateBSONDocument = bson.D{}
	courseBSONData, err := bson.Marshal(course)
	if err != nil {
		err = fmt.Errorf("[%s] - could not convert course (PID %s) to bson data", serverErrorMessages[seInputBSONNotValid], course.PID.Hex())
		return err
	}
	err = bson.Unmarshal(courseBSONData, &updateBSONDocument)
	if err != nil {
		err = fmt.Errorf("[%s] - could not convert course (PID %s) to bson document", serverErrorMessages[seInputBSONNotValid], course.PID.Hex())
		return err
	}
	var updateOptions = bson.D{{"$set", updateBSONDocument}}

	insertResult, err := dbPool.Collection(DBCollectionCourse).UpdateOne(context.TODO(), updateFilter, updateOptions)
	if err != nil {
		err = fmt.Errorf("[%s] - %s", serverErrorMessages[seDBResourceQuery], err.Error())
		return err
	}

	logging.Debugmf(logModCourseMgmt, "Update course (PID %s): matched %d modified %d",
		course.PID.Hex(), insertResult.MatchedCount, insertResult.ModifiedCount)
	if insertResult.MatchedCount == 0 {
		err = fmt.Errorf("[%s] - could not find course (PID %s)", serverErrorMessages[seResourceNotFound], course.PID.Hex())
		return err
	} else if insertResult.ModifiedCount == 0 {
		err = fmt.Errorf("[%s] - course (PID %s) not changed", serverErrorMessages[seResourceNotChange], course.PID.Hex())
		return err
	}
	return nil
}

// delete course, return #delete entries, error
func deleteCourse(pid primitive.ObjectID) (int, error) {
	var err error
	defer func() {
		if err != nil {
			logging.Errormf(logModCourseMgmt, err.Error())
		}
	}()

	// check course-student dependency
	studentCourseReferences, err := findStudentCourseRef(primitive.NilObjectID, pid)
	if err == nil && len(studentCourseReferences) > 0 {
		err = fmt.Errorf("[%s] - student-course dependency unresolved (e.g. student PID %s course PID %s)",
			serverErrorMessages[seDependencyIssue], studentCourseReferences[0].StudentPID.Hex(), studentCourseReferences[0].CoursePID.Hex())
		return 0, err
	}

	var deleteFilter bson.D = bson.D{}
	if !pid.IsZero() {
		deleteFilter = append(deleteFilter, bson.E{"_id", pid})
	}

	deleteResult, err := dbPool.Collection(DBCollectionCourse).DeleteMany(context.TODO(), deleteFilter)
	if err != nil {
		err = fmt.Errorf("[%s] - %s", serverErrorMessages[seDBResourceQuery], err.Error())
		return 0, err
	}

	logging.Debugmf(logModCourseMgmt, "Deleted %d course results from DB", deleteResult.DeletedCount)
	return int(deleteResult.DeletedCount), nil
}
