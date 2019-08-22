package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"

	"go.mongodb.org/mongo-driver/bson"
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

// find teacher by teacherUID
func findTeacherByUID(teacherUID string) (*Teacher, error) {
	var err error
	defer func() {
		if err != nil {
			logging.Errormf(logModTeacherHandler, err.Error())
		}
	}()

	var teacher Teacher
	findFilter := bson.D{{"teacher_uid", teacherUID}}
	err = dbPool.Collection(DBCollectionTeacher).FindOne(context.TODO(), findFilter).Decode(&teacher)
	if err != nil {
		err = fmt.Errorf("[%s] - %s", serverErrorMessages[seDBResourceQuery], err.Error())
		return nil, err
	}

	logging.Debugmf(logModTeacherHandler, "Found teacher from DB (teacherUID=%s)", teacherUID)
	return &teacher, nil
}
