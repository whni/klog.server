package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func registerUserLoginHandler(ctx *gin.Context) {
	params := ginContextRequestParameter(ctx)
	response := GinResponse{
		Status: http.StatusOK,
	}
	defer func() {
		ginContextProcessResponse(ctx, &response)
	}()

	var user RegisterUser
	var err error
	if err = json.Unmarshal(params.Data, &user); err != nil {
		response.Status = http.StatusBadRequest
		response.Message = fmt.Sprintf("[%s] - %s", serverErrorMessages[seInputJSONNotValid], err.Error())
		return
	}

	// find teacher by login information (teacher uid and password)
	var userFound *RegisterUser
	userFound, err = findRegisteruserByID(user.UserEmail)
	if err != nil {
		response.Status = http.StatusConflict
		response.Message = err.Error()
	} else if userFound.UserPassWord != user.UserPassWord {
		response.Status = http.StatusConflict
		response.Message = fmt.Sprintf("[%s] - user key not matched (email=%s)", serverErrorMessages[seResourceNotMatched], user.UserEmail)
	} else {
		response.Payload = user
	}
	return
}

func registerUserLogoutHandler(ctx *gin.Context) {

	return
}
