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

// main relative add or delete other relative
func relativeExtraEditHandler(ctx *gin.Context) {
	params := ginContextRequestParameter(ctx)
	response := GinResponse{
		Status: http.StatusOK,
	}
	defer func() {
		ginContextProcessResponse(ctx, &response)
	}()

	var studentRelativeEditInfo StudentRelativeEditInfo
	var err error
	if err = json.Unmarshal(params.Data, &studentRelativeEditInfo); err != nil {
		response.Status = http.StatusBadRequest
		response.Message = fmt.Sprintf("[%s] - %s", serverErrorMessages[seInputJSONNotValid], err.Error())
		return
	}

	if studentRelativeEditInfo.StudentPID.IsZero() ||
		studentRelativeEditInfo.RelativeWXID == "" || studentRelativeEditInfo.SecRelativeWXID == "" ||
		(studentRelativeEditInfo.Operation != "add" && studentRelativeEditInfo.Operation != "delete") {
		response.Status = http.StatusBadRequest
		response.Message = fmt.Sprintf("[%s] - Could not retrieve relative_wxid", serverErrorMessages[seInputJSONNotValid])
		return
	}

	// find relative by wechat id
	var relativeFound *Relative
	relativeFound, err = findRelativeByWXID(studentRelativeEditInfo.RelativeWXID)
	if err != nil || relativeFound == nil {
		response.Status = http.StatusConflict
		response.Message = fmt.Sprintf("[%s] - No relative found with wechat id \"%s\"", serverErrorMessages[seResourceNotFound], studentRelativeEditInfo.RelativeWXID)
		return
	}
	// check relative is main
	studentReferences, err := findStudentRelativeRef(studentRelativeEditInfo.StudentPID, relativeFound.PID)

	if err != nil {
		response.Status = http.StatusConflict
		response.Message = fmt.Sprintf("[%s] - %s -> could not search student-relative references with given relative_wxid %s",
			serverErrorMessages[seResourceNotFound], err.Error(), studentRelativeEditInfo.RelativeWXID)
		return
	}

	if len(studentReferences) == 0 || !studentReferences[0].IsMain {
		response.Status = http.StatusConflict
		response.Message = fmt.Sprintf("[%s] - %s -> could not search student-main-relative references with given relative_wxid %s",
			serverErrorMessages[seResourceNotFound], err.Error(), studentRelativeEditInfo.RelativeWXID)
		return
	}
	// find second relative by wechat id, if none, create a new one for add case or return for del case
	var secRelativeFound *Relative
	secRelativeFound, err = findRelativeByWXID(studentRelativeEditInfo.SecRelativeWXID)
	if err != nil || (secRelativeFound == nil && studentRelativeEditInfo.Operation == "delete") {
		response.Status = http.StatusConflict
		response.Message = fmt.Sprintf("[%s] - No relative found with wechat id \"%s\"", serverErrorMessages[seResourceNotFound], studentRelativeEditInfo.SecRelativeWXID)
		return
	}
	// create a new one
	if secRelativeFound == nil && studentRelativeEditInfo.Operation == "add" {
		var secRelative Relative
		secRelative.RelativeWXID = studentRelativeEditInfo.SecRelativeWXID
		_, err = createRelative(&secRelative)
		if err != nil {
			response.Status = http.StatusConflict
			response.Message = err.Error()
			return
		}
		secRelativeFound, err = findRelativeByWXID(studentRelativeEditInfo.SecRelativeWXID)
		if err != nil {
			response.Status = http.StatusConflict
			response.Message = fmt.Sprintf("[%s] - No relative found with wechat id after created \"%s\"", serverErrorMessages[seResourceNotFound], studentRelativeEditInfo.SecRelativeWXID)
			return
		}
	}
	var reference StudentRelativeRef
	var referencePID primitive.ObjectID
	reference.StudentPID = studentRelativeEditInfo.StudentPID
	reference.RelativePID = secRelativeFound.PID
	reference.IsMain = false
	reference.Relationship = studentRelativeEditInfo.Relationship

	if studentRelativeEditInfo.Operation == "add" {
		referencePID, err = createStudentRelativeRef(&reference)
		response.Payload = referencePID
	} else {
		_, err = deleteStudentRelativeRef(reference.StudentPID, reference.RelativePID)
		response.Payload = 1
	}
	if err != nil {
		response.Status = http.StatusConflict
		response.Message = fmt.Sprintf("[%s] - could not %s student-relative references with given student_id %s, relative_wxid %s and secrelative_wxid %s",
			serverErrorMessages[seResourceNotChange], studentRelativeEditInfo.Operation, studentRelativeEditInfo.StudentPID,
			studentRelativeEditInfo.RelativeWXID, studentRelativeEditInfo.SecRelativeWXID)
	}
	return
}
