package main

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func teacherLoginHandler(ctx *gin.Context) {
	params := ginContextRequestParameter(ctx)
	response := GinResponse{
		Status: http.StatusOK,
	}
	defer func() {
		ginContextProcessResponse(ctx, &response)
	}()

	var teacher Teacher
	var err error
	if err = json.Unmarshal(params.Data, &teacher); err != nil {
		response.Status = http.StatusBadRequest
		response.Message = fmt.Sprintf("[%s] - %s", serverErrorMessages[seInputJSONNotValid], err.Error())
		return
	}

	// find teacher by login information (teacher uid and password)
	var teacherFound *Teacher
	teacherFound, err = findTeacherByUID(teacher.TeacherUID)
	if err != nil {
		response.Status = http.StatusConflict
		response.Message = err.Error()
	} else if teacherFound.TeacherKey != teacher.TeacherKey {
		response.Status = http.StatusConflict
		response.Message = fmt.Sprintf("[%s] - teacher key not matched (UID=%s)", serverErrorMessages[seResourceNotMatched], teacher.TeacherUID)
	} else {
		response.Payload = teacherFound
	}
	return
}
