package main

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

// parent wechat login
func parentWeChatLoginHandler(ctx *gin.Context) {
	params := ginContextRequestParameter(ctx)
	response := GinResponse{
		Status: http.StatusOK,
	}
	defer func() {
		ginContextProcessResponse(ctx, &response)
	}()

	var parentWeChatLoginInfo ParentWeChatLoginInfo
	var err error
	if err = json.Unmarshal(params.Data, &parentWeChatLoginInfo); err != nil {
		response.Status = http.StatusBadRequest
		response.Message = fmt.Sprintf("[%s] - %s", serverErrorMessages[seInputJSONNotValid], err.Error())
		return
	}

	// check student parent wxID and binding code
	wxLoginURL := fmt.Sprintf("%s?appid=%s&secret=%s&js_code=%s&grant_type=authorization_code", serverConfig.ParentWeChatLoginURL,
		parentWeChatLoginInfo.AppID, parentWeChatLoginInfo.Secret, parentWeChatLoginInfo.JSCode)
	wxLoginResp, err := http.Get(wxLoginURL)
	if err != nil {
		response.Status = http.StatusBadRequest
		response.Message = fmt.Sprintf("[%s] - Could not make wechat login request: %s", serverErrorMessages[seInputParamNotValid], err.Error())
		return
	}

	var loginRespMap = map[string]interface{}{}
	err = json.NewDecoder(wxLoginResp.Body).Decode(&loginRespMap)
	if err != nil {
		response.Status = http.StatusConflict
		response.Message = fmt.Sprintf("[%s] - Could not retrieve wechat login response: %s", serverErrorMessages[seDataParseError], err.Error())
		return
	}

	if _, errExist := loginRespMap["errcode"]; errExist {
		response.Status = http.StatusConflict
		response.Message = fmt.Sprintf("[%s] - WeChat login errcode: %v, errmsg: %v", serverErrorMessages[seWeChatLoginError], loginRespMap["errcode"], loginRespMap["errmsg"])
		return
	}

	openID, openIDExist := loginRespMap["openid"]
	if !openIDExist {
		response.Status = http.StatusConflict
		response.Message = fmt.Sprintf("[%s] - Could not retrieve WeChat OpenID", serverErrorMessages[seWeChatLoginError])
		return
	}

	response.Payload = map[string]string{
		"parent_wxid": fmt.Sprintf("%v", openID),
	}
	return
}

// parent search student
func parentFindBoundStudentHandler(ctx *gin.Context) {
	params := ginContextRequestParameter(ctx)
	response := GinResponse{
		Status: http.StatusOK,
	}
	defer func() {
		ginContextProcessResponse(ctx, &response)
	}()

	var parentWXIDMap = map[string]string{}
	var err error
	if err = json.Unmarshal(params.Data, &parentWXIDMap); err != nil {
		response.Status = http.StatusBadRequest
		response.Message = fmt.Sprintf("[%s] - %s", serverErrorMessages[seInputJSONNotValid], err.Error())
		return
	}

	parentWXID, parentWXIDExist := parentWXIDMap["parent_wxid"]
	if !parentWXIDExist {
		response.Status = http.StatusBadRequest
		response.Message = fmt.Sprintf("[%s] - Could not retrieve parent_wxid", serverErrorMessages[seInputJSONNotValid])
		return
	}

	student, err := findStudentByParentWXID(parentWXID)
	if err != nil || student == nil {
		response.Status = http.StatusConflict
		response.Message = fmt.Sprintf("[%s] - Could not find student with given parent_wxid %s", serverErrorMessages[seResourceNotFound], parentWXID)
		return
	}

	response.Payload = student
	return
}
