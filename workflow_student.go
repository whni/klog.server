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

	// avoid repeated binding
	if studentFound.ParentWXID != "" {
		response.Status = http.StatusConflict
		response.Message = fmt.Sprintf("[%s] - Cannot binding parent wxID since parent (%s) has been already bound to student (%s)",
			serverErrorMessages[seResourceConflict], studentFound.ParentName, studentFound.StudentName)
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

func studentUnbindingParentHandler(ctx *gin.Context) {
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
	if studentInfo.ParentWXID == "" || studentInfo.PID.IsZero() {
		response.Status = http.StatusConflict
		response.Message = fmt.Sprintf("[%s] - Please specified valid parent wxID and student PID", serverErrorMessages[seInputJSONNotValid])
		return
	}

	// find student by binding code
	var students []*Student
	students, err = findStudent(studentInfo.PID)
	if err != nil || len(students) == 0 {
		response.Status = http.StatusConflict
		response.Message = fmt.Sprintf("[%s] - No student found with PID %s", serverErrorMessages[seResourceNotFound], studentInfo.PID.Hex())
		return
	}
	studentFound := students[0]

	// check parent wxid
	if studentFound.ParentWXID == "" {
		response.Status = http.StatusConflict
		response.Message = fmt.Sprintf("[%s] - No need to unbind parent wxID since nothing is in record", serverErrorMessages[seResourceNotFound])
		return
	}

	if studentFound.ParentWXID != studentInfo.ParentWXID {
		response.Status = http.StatusConflict
		response.Message = fmt.Sprintf("[%s] - Could not unbind parent wxID due to mismatched record (received %s recorded %s)", serverErrorMessages[seResourceNotMatched],
			studentFound.ParentWXID, studentInfo.ParentWXID)
		return
	}

	// update binding information
	studentFound.ParentWXID = ""
	studentFound.BindingCode = ""
	studentFound.BindingExpire = 0
	updateStudent(studentFound)

	response.Payload = studentFound
	return
}

// parent wechat login
func studentParentWXLogin(ctx *gin.Context) {
	params := ginContextRequestParameter(ctx)
	response := GinResponse{
		Status: http.StatusOK,
	}
	defer func() {
		ginContextProcessResponse(ctx, &response)
	}()

	var parentWXLoginInfo ParentWXLoginInfo
	var err error
	if err = json.Unmarshal(params.Data, &parentWXLoginInfo); err != nil {
		response.Status = http.StatusBadRequest
		response.Message = fmt.Sprintf("[%s] - %s", serverErrorMessages[seInputJSONNotValid], err.Error())
		return
	}

	// check student parent wxID and binding code
	wxLoginURL := fmt.Sprintf("%s?appid=%s&secret=%s&js_code=%s&grant_type=authorization_code", serverConfig.ParentWXLoginURL,
		parentWXLoginInfo.AppID, parentWXLoginInfo.Secret, parentWXLoginInfo.JSCode)
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
		"openid": fmt.Sprintf("%v", openID),
	}
	return
}
