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

	if studentInfo.PID.IsZero() || studentInfo.TeacherPID.IsZero() {
		response.Status = http.StatusConflict
		response.Message = fmt.Sprintf("[%s] - Please specified student PID and teacher PID", serverErrorMessages[seResourceNotFound])
		return
	}

	var students []*Student
	students, err = findStudent(studentInfo.PID)
	if err != nil || len(students) == 0 {
		response.Status = http.StatusConflict
		response.Message = fmt.Sprintf("[%s] - No student found with PID %s", serverErrorMessages[seResourceNotFound], studentInfo.PID.Hex())
		return
	}

	var studentFound = students[0]
	if studentInfo.TeacherPID.Hex() != studentFound.TeacherPID.Hex() {
		response.Status = http.StatusConflict
		response.Message = fmt.Sprintf("[%s] - teacher PID not matched (input:%s, found:%s)", serverErrorMessages[seResourceNotMatched],
			studentInfo.TeacherPID.Hex(), studentFound.TeacherPID.Hex())
		return
	}

	studentFound.BindingCode = xid.New().String()
	studentFound.BindingExpire = int64(time.Now().Unix()) + int64(3600*serverConfig.StudentBindingCodeLifeTime) // expired after one week
	updateStudent(studentFound)

	response.Payload = studentFound
	return
}
