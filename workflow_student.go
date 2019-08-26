package main

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/rs/xid"
	"net/http"
	"time"
)

func studentGenerateCodeHandler(ctx *gin.Context) {
	params := ginContextRequestParameter(ctx)
	response := GinResponse{
		Status: http.StatusOK,
	}
	defer func() {
		ginContextProcessResponse(ctx, &response)
	}()

	var studentInfo Student
	var err error
	if err = json.Unmarshal(params.Data, &studentInfo); err != nil {
		response.Status = http.StatusBadRequest
		response.Message = fmt.Sprintf("[%s] - %s", serverErrorMessages[seInputJSONNotValid], err.Error())
		return
	}

	// check input student PID and teacher PID
	if studentInfo.PID.IsZero() || studentInfo.TeacherPID.IsZero() {
		response.Status = http.StatusConflict
		response.Message = fmt.Sprintf("[%s] - Please specified student PID and teacher PID", serverErrorMessages[seInputJSONNotValid])
		return
	}

	// find student by PID
	var students []*Student
	students, err = findStudent(studentInfo.PID)
	if err != nil || len(students) == 0 {
		response.Status = http.StatusConflict
		response.Message = fmt.Sprintf("[%s] - No student found with PID %s", serverErrorMessages[seResourceNotFound], studentInfo.PID.Hex())
		return
	}

	// match student PID and teacher PID
	var studentFound = students[0]
	if studentInfo.TeacherPID.Hex() != studentFound.TeacherPID.Hex() {
		response.Status = http.StatusConflict
		response.Message = fmt.Sprintf("[%s] - teacher PID not matched (input:%s, found:%s)", serverErrorMessages[seResourceNotMatched],
			studentInfo.TeacherPID.Hex(), studentFound.TeacherPID.Hex())
		return
	}

	// generate binding code
	studentFound.BindingCode = xid.New().String()
	studentFound.BindingExpire = int64(time.Now().Unix()) + int64(3600*serverConfig.StudentBindingCodeLifeTime) // expired after one week
	updateStudent(studentFound)

	response.Payload = studentFound
	return
}

func studentBindingParentHandler(ctx *gin.Context) {
	params := ginContextRequestParameter(ctx)
	response := GinResponse{
		Status: http.StatusOK,
	}
	defer func() {
		ginContextProcessResponse(ctx, &response)
	}()

	var studentInfo Student
	var err error
	if err = json.Unmarshal(params.Data, &studentInfo); err != nil {
		response.Status = http.StatusBadRequest
		response.Message = fmt.Sprintf("[%s] - %s", serverErrorMessages[seInputJSONNotValid], err.Error())
		return
	}

	// check student parent wxID and binding code
	if studentInfo.ParentWXID == "" || studentInfo.BindingCode == "" {
		response.Status = http.StatusConflict
		response.Message = fmt.Sprintf("[%s] - Please specified valid parent wxID and binding code for student", serverErrorMessages[seInputJSONNotValid])
		return
	}

	// find student by binding code
	var studentFound *Student
	studentFound, err = findStudentByBindingCode(studentInfo.BindingCode)
	if err != nil || studentFound == nil {
		response.Status = http.StatusConflict
		response.Message = fmt.Sprintf("[%s] - No student found with binding code %s", serverErrorMessages[seResourceNotFound], studentInfo.BindingCode)
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

	// update binding information
	studentFound.ParentWXID = studentInfo.ParentWXID
	studentFound.ParentName = studentInfo.ParentName
	studentFound.PhoneNumber = studentInfo.PhoneNumber
	studentFound.Email = studentInfo.Email
	studentFound.BindingCode = ""
	studentFound.BindingExpire = 0
	updateStudent(studentFound)

	response.Payload = studentFound
	return
}
