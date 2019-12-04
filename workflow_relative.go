package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// relative wechat login
func relativeWeChatLoginHandler(ctx *gin.Context) {
	params := ginContextRequestParameter(ctx)
	response := GinResponse{
		Status: http.StatusOK,
	}
	defer func() {
		ginContextProcessResponse(ctx, &response)
	}()

	var relativeWeChatLoginInfo RelativeWeChatLoginInfo
	var err error
	if err = json.Unmarshal(params.Data, &relativeWeChatLoginInfo); err != nil {
		response.Status = http.StatusBadRequest
		response.Message = fmt.Sprintf("[%s] - %s", serverErrorMessages[seInputJSONNotValid], err.Error())
		return
	}

	// check student relative wxID and binding code
	wxLoginURL := fmt.Sprintf("%s?appid=%s&secret=%s&js_code=%s&grant_type=authorization_code", serverConfig.RelativeWeChatLoginURL,
		relativeWeChatLoginInfo.AppID, relativeWeChatLoginInfo.Secret, relativeWeChatLoginInfo.JSCode)
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
		"relative_wxid": fmt.Sprintf("%v", openID),
	}
	return
}

// relative search student
func relativeFindBoundStudentHandler(ctx *gin.Context) {
	params := ginContextRequestParameter(ctx)
	response := GinResponse{
		Status: http.StatusOK,
	}
	defer func() {
		ginContextProcessResponse(ctx, &response)
	}()

	var relativeWXIDMap = map[string]string{}
	var err error
	if err = json.Unmarshal(params.Data, &relativeWXIDMap); err != nil {
		response.Status = http.StatusBadRequest
		response.Message = fmt.Sprintf("[%s] - %s", serverErrorMessages[seInputJSONNotValid], err.Error())
		return
	}

	relativeWXID, relativeWXIDExist := relativeWXIDMap["relative_wxid"]
	if !relativeWXIDExist {
		response.Status = http.StatusBadRequest
		response.Message = fmt.Sprintf("[%s] - Could not retrieve relative_wxid", serverErrorMessages[seInputJSONNotValid])
		return
	}

	// find relative by wechat id
	var relativeFound *Relative
	relativeFound, err = findRelativeByWXID(relativeWXID)
	if err != nil || relativeFound == nil {
		response.Status = http.StatusConflict
		response.Message = fmt.Sprintf("[%s] - No relative found with wechat id \"%s\"", serverErrorMessages[seResourceNotFound], relativeWXID)
		return
	}

	var references []*StudentRelativeRef
	references, err = findStudentRelativeRef(primitive.NilObjectID, relativeFound.PID)
	if err != nil {
		response.Status = http.StatusConflict
		response.Message = fmt.Sprintf("[%s] - %s -> could not search student-relative references with given relative_wxid %s",
			serverErrorMessages[seResourceNotFound], err.Error(), relativeWXID)
		return
	}

	response.Payload = references
	return
}
