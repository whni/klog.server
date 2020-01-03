package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Azure/azure-storage-blob-go/azblob"
	"github.com/gin-gonic/gin"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var cloudMediaConfigHandlerTable = map[string]gin.HandlerFunc{
	"get":    cloudMediaGetHandler,
	"post":   cloudMediaPostHandler,
	"put":    cloudMediaPutHandler,
	"delete": cloudMediaDeleteHandler,
}

func cloudMediaGetHandler(ctx *gin.Context) {
	params := ginContextRequestParameter(ctx)
	response := GinResponse{
		Status: http.StatusOK,
	}
	defer func() {
		ginContextProcessResponse(ctx, &response)
	}()

	var cloudMediaSlice []*CloudMedia
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
	cloudMediaSlice, err = findCloudMedia(pid)
	if err != nil {
		response.Status = http.StatusConflict
		response.Message = err.Error()
		return
	}
	response.Payload = cloudMediaSlice
	return
}

func cloudMediaPostHandler(ctx *gin.Context) {
	params := ginContextRequestParameter(ctx)
	response := GinResponse{
		Status: http.StatusOK,
	}
	defer func() {
		ginContextProcessResponse(ctx, &response)
	}()

	var cloudMedia CloudMedia
	var cloudMediaPID primitive.ObjectID
	var err error

	if err = json.Unmarshal(params.Data, &cloudMedia); err != nil {
		response.Status = http.StatusBadRequest
		response.Message = fmt.Sprintf("[%s] - %s", serverErrorMessages[seInputJSONNotValid], err.Error())
		return
	}

	cloudMediaPID, err = createCloudMedia(&cloudMedia)
	if err != nil {
		response.Status = http.StatusConflict
		response.Message = err.Error()
	} else {
		response.Payload = cloudMediaPID
	}
	return
}

func cloudMediaPutHandler(ctx *gin.Context) {
	params := ginContextRequestParameter(ctx)
	response := GinResponse{
		Status: http.StatusOK,
	}
	defer func() {
		ginContextProcessResponse(ctx, &response)
	}()

	var cloudMedia CloudMedia
	var err error

	if err = json.Unmarshal(params.Data, &cloudMedia); err != nil {
		response.Status = http.StatusBadRequest
		response.Message = fmt.Sprintf("[%s] - %s", serverErrorMessages[seInputJSONNotValid], err.Error())
		return
	}

	err = updateCloudMedia(&cloudMedia)
	if err != nil {
		response.Status = http.StatusConflict
		response.Message = err.Error()
	} else {
		response.Payload = cloudMedia.PID
	}
	return
}

func cloudMediaDeleteHandler(ctx *gin.Context) {
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
	deletedRows, err = deleteCloudMedia(pid)
	if err != nil {
		response.Status = http.StatusConflict
		response.Message = err.Error()
		return
	}
	response.Payload = deletedRows
	return
}

// find cloud media, return cloud media slice, error
func findCloudMedia(pid primitive.ObjectID) ([]*CloudMedia, error) {
	var err error
	defer func() {
		if err != nil {
			logging.Errormf(logModCloudMediaMgmt, err.Error())
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

	findCursor, err := dbPool.Collection(DBCollectionCloudMedia).Find(context.TODO(), findFilter, findOptions)
	if err != nil {
		err = fmt.Errorf("[%s] - %s", serverErrorMessages[seDBResourceQuery], err.Error())
		return nil, err
	}

	cloudMediaSlice := []*CloudMedia{}
	for findCursor.Next(context.TODO()) {
		var cloudMedia CloudMedia
		err = findCursor.Decode(&cloudMedia)
		if err != nil {
			err = fmt.Errorf("[%s] - %s", serverErrorMessages[seDBResourceQuery], err.Error())
			return nil, err
		}
		cloudMediaSlice = append(cloudMediaSlice, &cloudMedia)
	}

	err = findCursor.Err()
	if err != nil {
		err = fmt.Errorf("[%s] - %s", serverErrorMessages[seDBResourceQuery], err.Error())
		return nil, err
	}

	logging.Debugmf(logModCloudMediaMgmt, "Found %d cloud media results from DB (PID=%v)", len(cloudMediaSlice), pid.Hex())
	return cloudMediaSlice, nil
}

// find cloud media by student pid, return cloud media slice, error
func findCloudMediaByStudentPID(studentPID primitive.ObjectID, onlyNilCourseRecord bool) ([]*CloudMedia, error) {
	var err error
	defer func() {
		if err != nil {
			logging.Errormf(logModCloudMediaMgmt, err.Error())
		}
	}()

	var findOptions = options.Find()
	var findFilter = bson.D{{"student_pid", studentPID}}
	if onlyNilCourseRecord {
		findFilter = append(findFilter, bson.E{"course_record_pid", primitive.NilObjectID})
	}

	findCursor, err := dbPool.Collection(DBCollectionCloudMedia).Find(context.TODO(), findFilter, findOptions)
	if err != nil {
		err = fmt.Errorf("[%s] - %s", serverErrorMessages[seDBResourceQuery], err.Error())
		return nil, err
	}

	cloudMediaSlice := []*CloudMedia{}
	for findCursor.Next(context.TODO()) {
		var cloudMedia CloudMedia
		err = findCursor.Decode(&cloudMedia)
		if err != nil {
			err = fmt.Errorf("[%s] - %s", serverErrorMessages[seDBResourceQuery], err.Error())
			return nil, err
		}
		cloudMediaSlice = append(cloudMediaSlice, &cloudMedia)
	}

	err = findCursor.Err()
	if err != nil {
		err = fmt.Errorf("[%s] - %s", serverErrorMessages[seDBResourceQuery], err.Error())
		return nil, err
	}

	logging.Debugmf(logModCloudMediaMgmt, "Found %d cloud media results from DB (studentPID=%v)", len(cloudMediaSlice), studentPID.Hex())
	return cloudMediaSlice, nil
}

// find cloud media by course record pid, return cloud media slice, error
func findCloudMediaByRecordPID(courseRecordPID primitive.ObjectID) ([]*CloudMedia, error) {
	var err error
	defer func() {
		if err != nil {
			logging.Errormf(logModCloudMediaMgmt, err.Error())
		}
	}()

	var findOptions = options.Find()
	var findFilter = bson.D{{"course_record_pid", courseRecordPID}}

	findCursor, err := dbPool.Collection(DBCollectionCloudMedia).Find(context.TODO(), findFilter, findOptions)
	if err != nil {
		err = fmt.Errorf("[%s] - %s", serverErrorMessages[seDBResourceQuery], err.Error())
		return nil, err
	}

	cloudMediaSlice := []*CloudMedia{}
	for findCursor.Next(context.TODO()) {
		var cloudMedia CloudMedia
		err = findCursor.Decode(&cloudMedia)
		if err != nil {
			err = fmt.Errorf("[%s] - %s", serverErrorMessages[seDBResourceQuery], err.Error())
			return nil, err
		}
		cloudMediaSlice = append(cloudMediaSlice, &cloudMedia)
	}

	err = findCursor.Err()
	if err != nil {
		err = fmt.Errorf("[%s] - %s", serverErrorMessages[seDBResourceQuery], err.Error())
		return nil, err
	}

	logging.Debugmf(logModCloudMediaMgmt, "Found %d cloud media results from DB (courseRecordPID=%v)", len(cloudMediaSlice), courseRecordPID.Hex())
	return cloudMediaSlice, nil
}

// create cloud media, return PID, error
func createCloudMedia(cloudMedia *CloudMedia) (primitive.ObjectID, error) {
	var err error
	defer func() {
		if err != nil {
			logging.Errormf(logModCloudMediaMgmt, err.Error())
		}
	}()

	// student PID check
	if cloudMedia.StudentPID.IsZero() {
		err = fmt.Errorf("[%s] - No student PID specified", serverErrorMessages[seResourceNotFound])
		return primitive.NilObjectID, err
	}

	var findFilter bson.M
	students, err := findStudent(cloudMedia.StudentPID, findFilter)
	if err != nil || len(students) == 0 {
		err = fmt.Errorf("[%s] - No associate student found with PID %s", serverErrorMessages[seResourceNotFound], cloudMedia.StudentPID.Hex())
		return primitive.NilObjectID, err
	}

	// course record PID check (pid = nil -> not related to any course record)
	if !cloudMedia.CourseRecordPID.IsZero() {
		var courseRecords []*CourseRecord
		courseRecords, err = findCourseRecord(cloudMedia.CourseRecordPID)
		if err != nil || len(courseRecords) == 0 {
			err = fmt.Errorf("[%s] - No course record found with PID %s", serverErrorMessages[seResourceNotFound], cloudMedia.CourseRecordPID.Hex())
			return primitive.NilObjectID, err
		}
		if courseRecords[0].StudentPID != cloudMedia.StudentPID {
			err = fmt.Errorf("[%s] - Student PID %s does not match course record (PID %s) student PID %s", serverErrorMessages[seResourceNotMatched],
				cloudMedia.StudentPID, cloudMedia.CourseRecordPID.Hex(), courseRecords[0].StudentPID)
			return primitive.NilObjectID, err
		}
	}

	// check if media exists at cloud and fill media information
	azProp, azPropErr := azureStorageGetBlobProperties(azMediaContainerURL, cloudMedia.MediaName)
	if azPropErr != nil {
		err = fmt.Errorf("[%s] - No media properties (URL: %s/%s) found at cloud (please check cloud connection and blob contents)",
			serverErrorMessages[seResourceNotFound], azMediaContainerURL.String(), cloudMedia.MediaName)
		return primitive.NilObjectID, err
	}
	cloudMedia.MediaURL = azMediaContainerURL.String() + "/" + cloudMedia.MediaName
	cloudMedia.CreateTS = azProp.CreateTS // NOTE: or we can use current timestamp
	cloudMedia.ContentLength = azProp.ContentLength

	// media type check
	if _, ok := cloudMediaTypeMap[cloudMedia.MediaType]; !ok {
		err = fmt.Errorf("[%s] - Cloud media type can only be %s/%s/%s", serverErrorMessages[seInputSchemaNotValid], CloudMediaTypeVideo, CloudMediaTypeImage, CloudMediaTypeOthers)
		return primitive.NilObjectID, err
	}

	insertResult, err := dbPool.Collection(DBCollectionCloudMedia).InsertOne(context.TODO(), cloudMedia)
	if err != nil {
		err = fmt.Errorf("[%s] - %s", serverErrorMessages[seDBResourceQuery], err.Error())
		return primitive.NilObjectID, err
	}

	lastInsertID := insertResult.InsertedID.(primitive.ObjectID)
	logging.Debugmf(logModCloudMediaMgmt, "Created cloud media in DB (LastInsertID,PID=%s)", lastInsertID.Hex())
	return lastInsertID, nil
}

// update cloud media (only rank score and student pid), return error
func updateCloudMedia(cloudMedia *CloudMedia) error {
	var err error
	defer func() {
		if err != nil {
			logging.Errormf(logModCloudMediaMgmt, err.Error())
		}
	}()

	// cloud media PID check
	if cloudMedia.PID.IsZero() {
		err = fmt.Errorf("[%s] - Cloud media PID is empty", serverErrorMessages[seInputJSONNotValid])
		return err
	}
	cloudMediaSlice, err := findCloudMedia(cloudMedia.PID)
	if err != nil || len(cloudMediaSlice) == 0 {
		err = fmt.Errorf("[%s] - No cloud media found with PID %s", serverErrorMessages[seResourceNotFound], cloudMedia.PID.Hex())
		return err
	}
	cloudMediaFound := cloudMediaSlice[0]

	// student PID check
	if cloudMedia.StudentPID.IsZero() {
		err = fmt.Errorf("[%s] - No student PID associated", serverErrorMessages[seResourceNotFound])
		return err
	}
	var findFilter bson.M
	students, err := findStudent(cloudMedia.StudentPID, findFilter)
	if err != nil || len(students) == 0 {
		err = fmt.Errorf("[%s] - No associated student found with PID %s", serverErrorMessages[seResourceNotFound], cloudMedia.StudentPID.Hex())
		return err
	}

	// course record PID check (pid = nil -> not related to any course record)
	if !cloudMedia.CourseRecordPID.IsZero() {
		var courseRecords []*CourseRecord
		courseRecords, err = findCourseRecord(cloudMedia.CourseRecordPID)
		if err != nil || len(courseRecords) == 0 {
			err = fmt.Errorf("[%s] - No course record found with PID %s", serverErrorMessages[seResourceNotFound], cloudMedia.CourseRecordPID.Hex())
			return err
		}
		if courseRecords[0].StudentPID != cloudMedia.StudentPID {
			err = fmt.Errorf("[%s] - Student PID %s does not match course record (PID %s) student PID %s", serverErrorMessages[seResourceNotMatched],
				cloudMedia.StudentPID, cloudMedia.CourseRecordPID.Hex(), courseRecords[0].StudentPID)
			return err
		}
	}

	// NOTE: only update partial attributes
	cloudMediaFound.StudentPID = cloudMedia.StudentPID
	cloudMediaFound.CourseRecordPID = cloudMedia.CourseRecordPID
	cloudMediaFound.MediaTags = []MediaTag{}
	cloudMediaFound.MediaTags = append(cloudMediaFound.MediaTags, cloudMedia.MediaTags...)

	var updateFilter = bson.D{{"_id", cloudMediaFound.PID}}
	var updateBSONDocument = bson.D{}
	cloudMediaBSONData, err := bson.Marshal(cloudMediaFound)
	if err != nil {
		err = fmt.Errorf("[%s] - could not convert cloud media (PID %s) to bson data", serverErrorMessages[seInputBSONNotValid], cloudMediaFound.PID.Hex())
		return err
	}
	err = bson.Unmarshal(cloudMediaBSONData, &updateBSONDocument)
	if err != nil {
		err = fmt.Errorf("[%s] - could not convert cloud media (PID %s) to bson document", serverErrorMessages[seInputBSONNotValid], cloudMediaFound.PID.Hex())
		return err
	}
	var updateOptions = bson.D{{"$set", updateBSONDocument}}

	insertResult, err := dbPool.Collection(DBCollectionCloudMedia).UpdateOne(context.TODO(), updateFilter, updateOptions)
	if err != nil {
		err = fmt.Errorf("[%s] - %s", serverErrorMessages[seDBResourceQuery], err.Error())
		return err
	}

	logging.Debugmf(logModCloudMediaMgmt, "Update cloud media (PID %s): matched %d modified %d",
		cloudMediaFound.PID.Hex(), insertResult.MatchedCount, insertResult.ModifiedCount)
	if insertResult.MatchedCount == 0 {
		err = fmt.Errorf("[%s] - could not find cloud media (PID %s)", serverErrorMessages[seResourceNotFound], cloudMediaFound.PID.Hex())
		return err
	} else if insertResult.ModifiedCount == 0 {
		err = fmt.Errorf("[%s] - cloud media (PID %s) not changed", serverErrorMessages[seResourceNotChange], cloudMediaFound.PID.Hex())
		return err
	}
	return nil
}

// delete cloud media, return #delete entries, error
func deleteCloudMedia(pid primitive.ObjectID) (int, error) {
	var err error
	defer func() {
		if err != nil {
			logging.Errormf(logModCloudMediaMgmt, err.Error())
		}
	}()

	// try to delete media files at cloud side
	cloudMediaSlice, findErr := findCloudMedia(pid)
	if findErr != nil {
		err = fmt.Errorf("[%s] - could not delete cloud media DB entries due to query error", serverErrorMessages[seResourceNotFound])
		return 0, err
	}

	var deleteCount int64
	for i := range cloudMediaSlice {
		cloudMedia := cloudMediaSlice[i]
		var azBlobDeleteErr error = azureStorageDeleteBlob(azMediaContainerURL, cloudMedia.MediaName)
		// return if error occurs (except blob not found)
		if azBlobDeleteErr != nil {
			if serr, ok := azBlobDeleteErr.(azblob.StorageError); !ok || serr.ServiceCode() != azblob.ServiceCodeBlobNotFound {
				err = fmt.Errorf("[%s] - could not delete cloud media blob at cloud (PID: %s name:%s type:%s) due to error: \"%s\"",
					serverErrorMessages[seCloudOpsError], cloudMedia.PID.Hex(), cloudMedia.MediaName, cloudMedia.MediaType, azBlobDeleteErr.Error())
				return int(deleteCount), err
			}
		}

		deleteFilter := bson.D{{"_id", cloudMedia.PID}}
		deleteResult, err := dbPool.Collection(DBCollectionCloudMedia).DeleteMany(context.TODO(), deleteFilter)
		if err != nil {
			err = fmt.Errorf("[%s] - %s", serverErrorMessages[seDBResourceQuery], err.Error())
			return int(deleteCount), err
		}

		deleteCount += deleteResult.DeletedCount
	}

	logging.Debugmf(logModCloudMediaMgmt, "Deleted %d cloud media results from DB (cloud media PID %s)", deleteCount, pid.Hex())
	return int(deleteCount), nil
}

// delete cloud media by student pid, return #delete entries, error
func deleteCloudMediaByStudentPID(studentPID primitive.ObjectID, onlyNilCourseRecord bool) (int, error) {
	var err error
	defer func() {
		if err != nil {
			logging.Errormf(logModCloudMediaMgmt, err.Error())
		}
	}()

	// try to delete media files at cloud side
	cloudMediaSlice, findErr := findCloudMediaByStudentPID(studentPID, onlyNilCourseRecord)
	if findErr != nil {
		err = fmt.Errorf("[%s] - could not delete cloud media DB entries due to query error", serverErrorMessages[seResourceNotFound])
		return 0, err
	}

	var deleteCount int64
	for i := range cloudMediaSlice {
		cloudMedia := cloudMediaSlice[i]
		var azBlobDeleteErr error = azureStorageDeleteBlob(azMediaContainerURL, cloudMedia.MediaName)
		// return if error occurs (except blob not found)
		if azBlobDeleteErr != nil {
			if serr, ok := azBlobDeleteErr.(azblob.StorageError); !ok || serr.ServiceCode() != azblob.ServiceCodeBlobNotFound {
				err = fmt.Errorf("[%s] - could not delete cloud media blob at cloud (PID: %s name:%s type:%s) due to error: [%s]",
					serverErrorMessages[seCloudOpsError], cloudMedia.PID.Hex(), cloudMedia.MediaName, cloudMedia.MediaType, azBlobDeleteErr.Error())
				return int(deleteCount), err
			}
		}

		deleteFilter := bson.D{{"_id", cloudMedia.PID}}
		deleteResult, err := dbPool.Collection(DBCollectionCloudMedia).DeleteMany(context.TODO(), deleteFilter)
		if err != nil {
			err = fmt.Errorf("[%s] - %s", serverErrorMessages[seDBResourceQuery], err.Error())
			return int(deleteCount), err
		}

		deleteCount += deleteResult.DeletedCount
	}

	logging.Debugmf(logModCloudMediaMgmt, "Deleted %d cloud media results from DB (studentPID %s)", deleteCount, studentPID.Hex())
	return int(deleteCount), nil
}

// delete cloud media by course record pid, return #delete entries, error
func deleteCloudMediaByRecordPID(courseRecordPID primitive.ObjectID) (int, error) {
	var err error
	defer func() {
		if err != nil {
			logging.Errormf(logModCloudMediaMgmt, err.Error())
		}
	}()

	// try to delete media files at cloud side
	cloudMediaSlice, findErr := findCloudMediaByRecordPID(courseRecordPID)
	if findErr != nil {
		err = fmt.Errorf("[%s] - could not delete cloud media DB entries due to query error", serverErrorMessages[seResourceNotFound])
		return 0, err
	}

	var deleteCount int64
	for i := range cloudMediaSlice {
		cloudMedia := cloudMediaSlice[i]
		var azBlobDeleteErr error = azureStorageDeleteBlob(azMediaContainerURL, cloudMedia.MediaName)
		// return if error occurs (except blob not found)
		if azBlobDeleteErr != nil {
			if serr, ok := azBlobDeleteErr.(azblob.StorageError); !ok || serr.ServiceCode() != azblob.ServiceCodeBlobNotFound {
				err = fmt.Errorf("[%s] - could not delete cloud media blob at cloud (PID: %s name:%s type:%s) due to error: [%s]",
					serverErrorMessages[seCloudOpsError], cloudMedia.PID.Hex(), cloudMedia.MediaName, cloudMedia.MediaType, azBlobDeleteErr.Error())
				return int(deleteCount), err
			}
		}

		deleteFilter := bson.D{{"_id", cloudMedia.PID}}
		deleteResult, err := dbPool.Collection(DBCollectionCloudMedia).DeleteMany(context.TODO(), deleteFilter)
		if err != nil {
			err = fmt.Errorf("[%s] - %s", serverErrorMessages[seDBResourceQuery], err.Error())
			return int(deleteCount), err
		}

		deleteCount += deleteResult.DeletedCount
	}

	logging.Debugmf(logModCloudMediaMgmt, "Deleted %d cloud media results from DB (courseRecordPID %s)", deleteCount, courseRecordPID.Hex())
	return int(deleteCount), nil
}
