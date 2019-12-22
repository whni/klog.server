package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/xid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func studentGenerateCodeHandler(ctx *gin.Context) {
	params := ginContextRequestParameter(ctx)
	response := GinResponse{
		Status: http.StatusOK,
	}
	defer func() {
		ginContextProcessResponse(ctx, &response)
	}()

	var studentRelativeBindInfo StudentRelativeBindInfo
	var err error
	if err = json.Unmarshal(params.Data, &studentRelativeBindInfo); err != nil {
		response.Status = http.StatusBadRequest
		response.Message = fmt.Sprintf("[%s] - %s", serverErrorMessages[seInputJSONNotValid], err.Error())
		return
	}

	// check input student PID
	if studentRelativeBindInfo.StudentPID.IsZero() {
		response.Status = http.StatusConflict
		response.Message = fmt.Sprintf("[%s] - Please specified student PID", serverErrorMessages[seInputJSONNotValid])
		return
	}

	// find student by PID
	var findFilter bson.M
	var students []*Student
	students, err = findStudent(studentRelativeBindInfo.StudentPID, findFilter)
	if err != nil || len(students) == 0 {
		response.Status = http.StatusConflict
		response.Message = fmt.Sprintf("[%s] - No student found with PID %s", serverErrorMessages[seResourceNotFound], studentRelativeBindInfo.StudentPID.Hex())
		return
	}
	var studentFound = students[0]

	// do not generate code if a main relative relationship exists
	studentReferences, err := findStudentRelativeRef(studentFound.PID, primitive.NilObjectID)
	var mainReference *StudentRelativeRef = nil
	if err == nil && len(studentReferences) > 0 {
		for i := 0; i < len(studentReferences); i++ {
			if studentReferences[i].IsMain {
				mainReference = studentReferences[i]
				break
			}
		}
	}
	if mainReference != nil {
		response.Status = http.StatusConflict
		response.Message = fmt.Sprintf("[%s] - Student (PID %s) already has a main relative binding (PID %s) -> no binding code generated",
			serverErrorMessages[seResourceConflict], studentFound.PID.Hex(), mainReference.RelativePID.Hex())
		return
	}

	// generate binding code
	studentFound.BindingCode = xid.New().String()
	studentFound.BindingExpire = int64(time.Now().Unix()) + int64(3600*serverConfig.StudentBindingCodeLifeTime) // expired after one week
	updateStudent(studentFound)

	response.Payload = studentFound
	return
}

// bind student with a main relative
func studentBindingMainRelativeHandler(ctx *gin.Context) {
	params := ginContextRequestParameter(ctx)
	response := GinResponse{
		Status: http.StatusOK,
	}
	defer func() {
		ginContextProcessResponse(ctx, &response)
	}()

	var studentRelativeBindInfo StudentRelativeBindInfo
	var err error
	if err = json.Unmarshal(params.Data, &studentRelativeBindInfo); err != nil {
		response.Status = http.StatusBadRequest
		response.Message = fmt.Sprintf("[%s] - %s", serverErrorMessages[seInputJSONNotValid], err.Error())
		return
	}

	// check student binding code and relative wechat id
	if studentRelativeBindInfo.RelativeWXID == "" || studentRelativeBindInfo.BindingCode == "" {
		response.Status = http.StatusConflict
		response.Message = fmt.Sprintf("[%s] - Please specified valid relative_wxid and binding_code", serverErrorMessages[seInputJSONNotValid])
		return
	}

	// find relative by wechat id
	var relativeFound *Relative
	relativeFound, err = findRelativeByWXID(studentRelativeBindInfo.RelativeWXID)
	if err != nil || relativeFound == nil {
		response.Status = http.StatusConflict
		response.Message = fmt.Sprintf("[%s] - No relative found with wechat id \"%s\"", serverErrorMessages[seResourceNotFound], studentRelativeBindInfo.RelativeWXID)
		return
	}

	// find student by binding code
	var studentFound *Student
	studentFound, err = findStudentByBindingCode(studentRelativeBindInfo.BindingCode)
	if err != nil || studentFound == nil {
		response.Status = http.StatusConflict
		response.Message = fmt.Sprintf("[%s] - No student found with binding code \"%s\"", serverErrorMessages[seResourceNotFound], studentRelativeBindInfo.BindingCode)
		return
	}

	// check if binding code is expired
	var curTS = int64(time.Now().Unix())
	if curTS > studentFound.BindingExpire {
		response.Status = http.StatusConflict
		response.Message = fmt.Sprintf("[%s] - Binding code (%s) for student (PID %s) is expired at %v", serverErrorMessages[seResourceExpired],
			studentFound.BindingCode, studentFound.PID.Hex(), time.Unix(studentFound.BindingExpire, 0).Format(time.RFC3339))
		return
	}

	var reference StudentRelativeRef
	reference.StudentPID = studentFound.PID
	reference.RelativePID = relativeFound.PID
	reference.IsMain = true
	if studentRelativeBindInfo.Relationship == "" {
		reference.Relationship = "Unknown"
	} else {
		reference.Relationship = studentRelativeBindInfo.Relationship
	}

	var referencePID = primitive.NilObjectID
	referencePID, err = createStudentRelativeRef(&reference)
	if err != nil {
		response.Status = http.StatusConflict
		response.Message = fmt.Sprintf("[%s] - %s --> could not bind student (PID %s) and relative (PID %s)", serverErrorMessages[seResourceConflict],
			err.Error(), studentFound.PID.Hex(), relativeFound.PID.Hex())
		return
	}

	// remove binding code from student information
	studentFound.BindingCode = ""
	studentFound.BindingExpire = 0
	updateStudent(studentFound)

	response.Payload = referencePID
	return
}

// unbind all relatives for a student (delete main relative-> other should be deleted)
func studentUnbindingAllRelativeHandler(ctx *gin.Context) {
	params := ginContextRequestParameter(ctx)
	response := GinResponse{
		Status: http.StatusOK,
	}
	defer func() {
		ginContextProcessResponse(ctx, &response)
	}()

	var studentRelativeBindInfo StudentRelativeBindInfo
	var err error
	if err = json.Unmarshal(params.Data, &studentRelativeBindInfo); err != nil {
		response.Status = http.StatusBadRequest
		response.Message = fmt.Sprintf("[%s] - %s", serverErrorMessages[seInputJSONNotValid], err.Error())
		return
	}

	// check student PID
	if studentRelativeBindInfo.StudentPID.IsZero() {
		response.Status = http.StatusConflict
		response.Message = fmt.Sprintf("[%s] - Please specified valid student PID", serverErrorMessages[seInputJSONNotValid])
		return
	}

	var deleteCnt int = 0
	deleteCnt, err = deleteStudentRelativeRef(studentRelativeBindInfo.StudentPID, primitive.NilObjectID)
	if err != nil {
		response.Status = http.StatusConflict
		response.Message = fmt.Sprintf("[%s] - Error occurs during unbind student (PID %s) relative references - %s",
			serverErrorMessages[seInputJSONNotValid], studentRelativeBindInfo.StudentPID.Hex(), err.Error())
		return
	}

	response.Payload = fmt.Sprintf("Unbind %d relatives for student (PID %s)", deleteCnt, studentRelativeBindInfo.StudentPID.Hex())
	return
}

func studentMediaQueryHandler(ctx *gin.Context) {
	params := ginContextRequestParameter(ctx)
	response := GinResponse{
		Status: http.StatusOK,
	}
	defer func() {
		ginContextProcessResponse(ctx, &response)
	}()

	var mediaReq StudentMediaQueryReq
	var err error
	if err = json.Unmarshal(params.Data, &mediaReq); err != nil {
		response.Status = http.StatusBadRequest
		response.Message = fmt.Sprintf("[%s] - %s", serverErrorMessages[seInputJSONNotValid], err.Error())
		return
	}

	// find student by PID
	var findFilter bson.M
	var students []*Student
	students, err = findStudent(mediaReq.StudentPID, findFilter)
	if err != nil || len(students) == 0 {
		response.Status = http.StatusConflict
		response.Message = fmt.Sprintf("[%s] - No student found with PID %s", serverErrorMessages[seResourceNotFound], mediaReq.StudentPID.Hex())
		return
	}

	cloudMediaSlice, err := findCloudMediaByStudentPID(mediaReq.StudentPID, false)
	if err != nil {
		response.Status = http.StatusConflict
		response.Message = fmt.Sprintf("[%s] - Error occurs when searching cloud media for student (PID %s): %s",
			serverErrorMessages[seResourceNotFound], mediaReq.StudentPID.Hex(), err.Error())
		return
	}

	cloudMediaRes := []*CloudMedia{}
	for i := range cloudMediaSlice {
		if cloudMediaSlice[i].CreateTS <= mediaReq.EndTS && cloudMediaSlice[i].CreateTS >= mediaReq.StartTS {
			cloudMediaRes = append(cloudMediaRes, cloudMediaSlice[i])
		}
	}

	sort.Slice(cloudMediaRes, func(i, j int) bool {
		return cloudMediaRes[i].CreateTS < cloudMediaRes[j].CreateTS
	})

	response.Payload = cloudMediaRes
	return
}

func studentStoryQueryHandler(ctx *gin.Context) {
	params := ginContextRequestParameter(ctx)
	response := GinResponse{
		Status: http.StatusOK,
	}
	defer func() {
		ginContextProcessResponse(ctx, &response)
	}()

	var storyReq StudentMediaQueryReq
	var err error
	if err = json.Unmarshal(params.Data, &storyReq); err != nil {
		response.Status = http.StatusBadRequest
		response.Message = fmt.Sprintf("[%s] - %s", serverErrorMessages[seInputJSONNotValid], err.Error())
		return
	}

	// find student by PID
	var findFilter bson.M
	var students []*Student
	students, err = findStudent(storyReq.StudentPID, findFilter)
	if err != nil || len(students) == 0 {
		response.Status = http.StatusConflict
		response.Message = fmt.Sprintf("[%s] - No student found with PID %s", serverErrorMessages[seResourceNotFound], storyReq.StudentPID.Hex())
		return
	}
	var storys []*Story
	var pid primitive.ObjectID
	pid = primitive.NilObjectID

	findFilter = bson.M{"student_pid": storyReq.StudentPID, "store_ts": bson.M{"$gte": storyReq.StartTS, "$lte": storyReq.EndTS}}

	storys, err = findStory(pid, findFilter)
	if err != nil {
		response.Status = http.StatusConflict
		response.Message = err.Error()
		return
	}
	response.Payload = storys

	return
}
