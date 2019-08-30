package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/Azure/azure-storage-blob-go/azblob"
	"github.com/gin-gonic/gin"
	"net/http"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var studentConfigHandlerTable = map[string]gin.HandlerFunc{
	"get":    studentGetHandler,
	"post":   studentPostHandler,
	"put":    studentPutHandler,
	"delete": studentDeleteHandler,
}

func studentGetImageName(student *Student) string {
	if student == nil {
		return ""
	}
	return "image-student-" + student.PID.Hex() + ".jpg"
}

func studentGetHandler(ctx *gin.Context) {
	params := ginContextRequestParameter(ctx)
	response := GinResponse{
		Status: http.StatusOK,
	}
	defer func() {
		ginContextProcessResponse(ctx, &response)
	}()

	var students []*Student
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
	students, err = findStudent(pid)
	if err != nil {
		response.Status = http.StatusConflict
		response.Message = err.Error()
		return
	}
	response.Payload = students
	return
}

func studentPostHandler(ctx *gin.Context) {
	params := ginContextRequestParameter(ctx)
	response := GinResponse{
		Status: http.StatusOK,
	}
	defer func() {
		ginContextProcessResponse(ctx, &response)
	}()

	var student Student
	var studentPID primitive.ObjectID
	var err error

	if err = json.Unmarshal(params.Data, &student); err != nil {
		response.Status = http.StatusBadRequest
		response.Message = fmt.Sprintf("[%s] - %s", serverErrorMessages[seInputJSONNotValid], err.Error())
		return
	}

	studentPID, err = createStudent(&student)
	if err != nil {
		response.Status = http.StatusConflict
		response.Message = err.Error()
	} else {
		response.Payload = studentPID
	}
	return
}

func studentPutHandler(ctx *gin.Context) {
	params := ginContextRequestParameter(ctx)
	response := GinResponse{
		Status: http.StatusOK,
	}
	defer func() {
		ginContextProcessResponse(ctx, &response)
	}()

	var student Student
	var err error

	if err = json.Unmarshal(params.Data, &student); err != nil {
		response.Status = http.StatusBadRequest
		response.Message = fmt.Sprintf("[%s] - %s", serverErrorMessages[seInputJSONNotValid], err.Error())
		return
	}

	err = updateStudent(&student)
	if err != nil {
		response.Status = http.StatusConflict
		response.Message = err.Error()
	} else {
		response.Payload = student.PID
	}
	return
}

func studentDeleteHandler(ctx *gin.Context) {
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
	deletedRows, err = deleteStudent(pid)
	if err != nil {
		response.Status = http.StatusConflict
		response.Message = err.Error()
		return
	}
	response.Payload = deletedRows
	return
}

// find student, return student slice, error
func findStudent(pid primitive.ObjectID) ([]*Student, error) {
	var err error
	defer func() {
		if err != nil {
			logging.Errormf(logModStudentHandler, err.Error())
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

	findCursor, err := dbPool.Collection(DBCollectionStudent).Find(context.TODO(), findFilter, findOptions)
	if err != nil {
		err = fmt.Errorf("[%s] - %s", serverErrorMessages[seDBResourceQuery], err.Error())
		return nil, err
	}

	students := []*Student{}
	for findCursor.Next(context.TODO()) {
		var student Student
		err = findCursor.Decode(&student)
		if err != nil {
			err = fmt.Errorf("[%s] - %s", serverErrorMessages[seDBResourceQuery], err.Error())
			return nil, err
		}
		students = append(students, &student)
	}

	err = findCursor.Err()
	if err != nil {
		err = fmt.Errorf("[%s] - %s", serverErrorMessages[seDBResourceQuery], err.Error())
		return nil, err
	}

	logging.Debugmf(logModStudentHandler, "Found %d student results from DB (PID=%v)", len(students), pid.Hex())
	return students, nil
}

// find student by binding code
func findStudentByBindingCode(bindingCode string) (*Student, error) {
	var err error
	defer func() {
		if err != nil {
			logging.Errormf(logModStudentHandler, err.Error())
		}
	}()

	var student Student
	findFilter := bson.D{{"binding_code", bindingCode}}
	err = dbPool.Collection(DBCollectionStudent).FindOne(context.TODO(), findFilter).Decode(&student)
	if err != nil {
		err = fmt.Errorf("[%s] - %s", serverErrorMessages[seDBResourceQuery], err.Error())
		return nil, err
	}

	logging.Debugmf(logModStudentHandler, "Found student from DB (studentPID=%s, bindingCode=%s)", student.PID.Hex(), bindingCode)
	return &student, nil
}

// find student by parent wechat id
func findStudentByParentWXID(parentWXID string) (*Student, error) {
	var err error
	defer func() {
		if err != nil {
			logging.Errormf(logModStudentHandler, err.Error())
		}
	}()

	var student Student
	findFilter := bson.D{{"parent_wxid", parentWXID}}
	err = dbPool.Collection(DBCollectionStudent).FindOne(context.TODO(), findFilter).Decode(&student)
	if err != nil {
		err = fmt.Errorf("[%s] - %s", serverErrorMessages[seDBResourceQuery], err.Error())
		return nil, err
	}

	logging.Debugmf(logModStudentHandler, "Found student from DB (studentPID=%s, parentWXID=%s)", student.PID.Hex(), parentWXID)
	return &student, nil
}

// create student, return PID, error
func createStudent(student *Student) (primitive.ObjectID, error) {
	var err error
	defer func() {
		if err != nil {
			logging.Errormf(logModStudentHandler, err.Error())
		}
	}()

	// teacher PID check
	if student.TeacherPID.IsZero() {
		err = fmt.Errorf("[%s] - No teacher PID specified", serverErrorMessages[seResourceNotFound])
		return primitive.NilObjectID, err
	}
	teachers, err := findTeacher(student.TeacherPID)
	if err != nil || len(teachers) == 0 {
		err = fmt.Errorf("[%s] - No teachers found with PID %s", serverErrorMessages[seResourceNotFound], student.TeacherPID.Hex())
		return primitive.NilObjectID, err
	}

	// student image name/url
	studentImageName := studentGetImageName(student)
	if _, imagePropErr := azureStorageGetBlobProperties(azMediaContainerURL, studentImageName); imagePropErr == nil {
		student.StudentImageName = studentImageName
		student.StudentImageURL = azMediaContainerURL.String() + "/" + student.StudentImageName
	} else {
		student.StudentImageName = ""
		student.StudentImageURL = ""
	}

	insertResult, err := dbPool.Collection(DBCollectionStudent).InsertOne(context.TODO(), student)
	if err != nil {
		err = fmt.Errorf("[%s] - %s", serverErrorMessages[seDBResourceQuery], err.Error())
		return primitive.NilObjectID, err
	}

	lastInsertID := insertResult.InsertedID.(primitive.ObjectID)
	logging.Debugmf(logModStudentHandler, "Created student in DB (LastInsertID,PID=%s)", lastInsertID.Hex())
	return lastInsertID, nil
}

// update student, return error
func updateStudent(student *Student) error {
	var err error
	defer func() {
		if err != nil {
			logging.Errormf(logModStudentHandler, err.Error())
		}
	}()

	// student PID check
	if student.PID.IsZero() {
		err = fmt.Errorf("[%s] - student PID is empty", serverErrorMessages[seInputJSONNotValid])
		return err
	}

	// teacher PID check
	if student.TeacherPID.IsZero() {
		err = fmt.Errorf("[%s] - No teacher PID specified", serverErrorMessages[seResourceNotFound])
		return err
	}
	teachers, err := findTeacher(student.TeacherPID)
	if err != nil || len(teachers) == 0 {
		err = fmt.Errorf("[%s] - No teacher found with PID %s", serverErrorMessages[seResourceNotFound], student.TeacherPID.Hex())
		return err
	}

	// student image name/url
	studentImageName := studentGetImageName(student)
	if _, imagePropErr := azureStorageGetBlobProperties(azMediaContainerURL, studentImageName); imagePropErr == nil {
		student.StudentImageName = studentImageName
		student.StudentImageURL = azMediaContainerURL.String() + "/" + student.StudentImageName
	} else {
		student.StudentImageName = ""
		student.StudentImageURL = ""
	}

	// update student
	var updateFilter = bson.D{{"_id", student.PID}}
	var updateBSONDocument = bson.D{}
	studentBSONData, err := bson.Marshal(student)
	if err != nil {
		err = fmt.Errorf("[%s] - could not convert student (PID %s) to bson data", serverErrorMessages[seInputBSONNotValid], student.PID.Hex())
		return err
	}
	err = bson.Unmarshal(studentBSONData, &updateBSONDocument)
	if err != nil {
		err = fmt.Errorf("[%s] - could not convert student (PID %s) to bson document", serverErrorMessages[seInputBSONNotValid], student.PID.Hex())
		return err
	}
	var updateOptions = bson.D{{"$set", updateBSONDocument}}

	insertResult, err := dbPool.Collection(DBCollectionStudent).UpdateOne(context.TODO(), updateFilter, updateOptions)
	if err != nil {
		err = fmt.Errorf("[%s] - %s", serverErrorMessages[seDBResourceQuery], err.Error())
		return err
	}

	logging.Debugmf(logModStudentHandler, "Update student (PID %s): matched %d modified %d",
		student.PID.Hex(), insertResult.MatchedCount, insertResult.ModifiedCount)
	if insertResult.MatchedCount == 0 {
		err = fmt.Errorf("[%s] - could not find student (PID %s)", serverErrorMessages[seResourceNotFound], student.PID.Hex())
		return err
	} else if insertResult.ModifiedCount == 0 {
		err = fmt.Errorf("[%s] - student (PID %s) not changed", serverErrorMessages[seResourceNotChange], student.PID.Hex())
		return err
	}
	return nil
}

// delete student, return #delete entries, error
func deleteStudent(pid primitive.ObjectID) (int, error) {
	var err error
	defer func() {
		if err != nil {
			logging.Errormf(logModStudentHandler, err.Error())
		}
	}()

	students, findErr := findStudent(pid)
	if findErr != nil {
		err = fmt.Errorf("[%s] - could not delete student (PID %s) due to DB query/find error occurs", serverErrorMessages[seDBResourceQuery], pid.Hex())
		return 0, err
	}

	var deleteCnt int64
	for i := range students {
		_, deleteCloudMediaErr := deleteCloudMediaByStudentPID(students[i].PID)
		if deleteCloudMediaErr != nil {
			err = fmt.Errorf("[%s] - stop deleting student (PID %s) since cloud media could not be deleted: %s",
				serverErrorMessages[seCloudOpsError], pid.Hex(), deleteCloudMediaErr.Error())
			return int(deleteCnt), err
		}

		if students[i].StudentImageName != "" {
			if deleteStudentImageErr := azureStorageDeleteBlob(azMediaContainerURL, students[i].StudentImageName); deleteStudentImageErr != nil {
				if serr, ok := deleteStudentImageErr.(azblob.StorageError); !ok || serr.ServiceCode() != azblob.ServiceCodeBlobNotFound {
					err = fmt.Errorf("[%s] - could not delete student image at cloud (PID: %s image name:%s image url:%s) due to error: [%s]",
						serverErrorMessages[seCloudOpsError], students[i].PID.Hex(), students[i].StudentImageName, students[i].StudentImageURL, deleteStudentImageErr.Error())
					return int(deleteCnt), err
				}
			}
		}

		deleteFilter := bson.D{{"_id", pid}}
		deleteResult, err := dbPool.Collection(DBCollectionStudent).DeleteMany(context.TODO(), deleteFilter)
		if err != nil {
			err = fmt.Errorf("[%s] - %s", serverErrorMessages[seDBResourceQuery], err.Error())
			return 0, err
		}

		deleteCnt += deleteResult.DeletedCount
	}

	logging.Debugmf(logModStudentHandler, "Deleted %d student results from DB", deleteCnt)
	return int(deleteCnt), nil
}
