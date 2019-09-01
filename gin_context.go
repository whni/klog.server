package main

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"

	//"reflect"
	"strings"
)

var ginConfigAPITable = map[string]map[string]gin.HandlerFunc{
	"/api/0/config/institute":  instituteConfigHandlerTable,
	"/api/0/config/teacher":    teacherConfigHandlerTable,
	"/api/0/config/student":    studentConfigHandlerTable,
	"/api/0/config/cloudmedia": cloudMediaConfigHandlerTable,
}

var ginWorkflowAPITable = map[string]gin.HandlerFunc{
	"/api/0/workflow/teacher/login":        teacherLoginHandler,
	"/api/0/workflow/student/generatecode": studentGenerateCodeHandler,
	"/api/0/workflow/student/bind":         studentBindingParentHandler,
	"/api/0/workflow/student/unbind":       studentUnbindingParentHandler,
	"/api/0/workflow/student/mediaquery":   studentMediaQueryHandler,
	"/api/0/workflow/parent/wxlogin":       parentWeChatLoginHandler,
	"/api/0/workflow/parent/findstudent":   parentFindBoundStudentHandler,
}

// GinParameter a generic paramter wrapper for gin web framework handler
type GinParameter struct {
	Token string // auth token
	PID   string // primary id
	SID   string // secondary id
	Data  []byte // data content
}

// GinResponse a generic response struct for gin web framework handler
type GinResponse struct {
	Status  int         `json:"-"`
	Message string      `json:"message"`
	Payload interface{} `json:"payload"`
}

func ginContextRequestParameter(ctx *gin.Context) *GinParameter {
	var tokenString string = ""
	if token, exist := ctx.Get("JWT_TOKEN"); exist {
		tokenString = token.(string)
	}

	var primaryKey = ctx.Request.URL.Query().Get("pid")
	var secondarykey = ctx.Request.URL.Query().Get("sid")

	var bodyData []byte = nil
	switch strings.ToLower(ctx.Request.Method) {
	case "post", "put":
		if reqBuffer, err := ioutil.ReadAll(ctx.Request.Body); err != nil {
			logging.Warnmf(logModGinContext, "Http request body read err: %v\n", err.Error())
		} else {
			bodyData = reqBuffer
		}
	default:
	}

	return &GinParameter{tokenString, primaryKey, secondarykey, bodyData}
}

func ginContextProcessResponse(ctx *gin.Context, response *GinResponse) {
	responseContent := gin.H{}
	if response.Status == http.StatusOK {
		if response.Payload != nil {
			responseContent["payload"] = response.Payload
		} else {
			// payload should exist with http status ok ==> internal error occurs
			response.Status = http.StatusInternalServerError
			responseContent["message"] = fmt.Sprintf("[%s] - Please check error log", serverErrorMessages[seUnresolvedError])
		}
		if len(response.Message) > 0 {
			responseContent["message"] = response.Message
		}
	} else {
		responseContent["message"] = response.Message
	}
	ctx.JSON(response.Status, responseContent)
	logging.Debugmf(logModGinContext, "GIN: [%s] FROM %v | URL %v | RESPONSE CODE %v", ctx.Request.Method, ctx.Request.RemoteAddr, ctx.Request.URL, response.Status)
}

/* gin input struct check */
// func ginStructValidCheck(input interface{}) error {
// 	val := reflect.ValueOf(input)
// 	if val.Kind() == reflect.Ptr {
// 		val = reflect.Indirect(val)
// 	}

// 	if val.Kind() != reflect.Struct {
// 		return fmt.Errorf("[%s] - unexpected type (%s) - struct required", serverErrorMessages[seInputSchemaNotValid], val.Type().Name())
// 	}
// 	structType := val.Type()

// 	for i := 0; i < structType.NumField(); i++ {
// 		field := structType.Field(i)
// 		fieldName := field.Name
// 		logging.Tracemf(logModGinContext, "check struct (%s) field: %s", val.Type().Name(), fieldName)

// 		switch fieldName {
// 		case "InstituteUID", "ClassUID", "TeacherUID", "StudentUID", "ParentUID",
// 			"Name", "InstituteName", "ClassName", "StudentName", "ParentName",
// 			"Address", "CountryCode", "Location", "MediaLocation",
// 			"PhoneNumber", "Email", "Occupation":
// 			if len(val.FieldByName(fieldName).Interface().(string)) < 1 {
// 				return fmt.Errorf("[%s] - invalid struct (%s) field: %s", serverErrorMessages[seInputSchemaNotValid], structType.Name(), fieldName)
// 			}
// 			break
// 		case "Password":
// 			if len(val.FieldByName(fieldName).Interface().(string)) < 5 {
// 				return fmt.Errorf("[%s] - invalid struct (%s) field: %s (at least 5 length)", serverErrorMessages[seInputSchemaNotValid], structType.Name(), fieldName)
// 			}
// 			break
// 		case "PID":
// 			if val.FieldByName(fieldName).Interface().(int) < 0 {
// 				return fmt.Errorf("[%s] - invalid struct (%s) field: %s (require postive value)", serverErrorMessages[seInputSchemaNotValid], structType.Name(), fieldName)
// 			}
// 			break
// 		case "InstitutePID", "ClassPID", "StudentPID", "TeacherPID", "ParentPID":
// 			pid := val.FieldByName(fieldName).Interface().(null.Int)
// 			if pid.Valid && pid.ValueOrZero() <= 0 {
// 				return fmt.Errorf("[%s] - invalid struct (%s) field: %s (require postive value)", serverErrorMessages[seInputSchemaNotValid], structType.Name(), fieldName)
// 			}
// 			break
// 		}
// 	}
// 	return nil
// }

// func ginStructEqualCheck(x, y interface{}) error {
// 	valx := reflect.ValueOf(x)
// 	if valx.Kind() == reflect.Ptr {
// 		valx = reflect.Indirect(valx)
// 	}

// 	valy := reflect.ValueOf(y)
// 	if valy.Kind() == reflect.Ptr {
// 		valy = reflect.Indirect(valy)
// 	}

// 	if valx.Kind() != reflect.Struct || valy.Kind() != reflect.Struct {
// 		return fmt.Errorf("[%s] - unexpected type (%s, %s) - struct required", serverErrorMessages[seInputSchemaNotValid], valx.Type().Name(), valy.Type().Name())
// 	}
// 	structTypex := valx.Type()
// 	structTypey := valy.Type()

// 	if structTypex.Name() != structTypey.Name() {
// 		return fmt.Errorf("[%s] - inconsistent type (%s != %s) ", serverErrorMessages[seInputSchemaNotValid], valx.Type().Name(), valy.Type().Name())
// 	}

// 	for i := 0; i < structTypex.NumField(); i++ {
// 		field := structTypex.Field(i)
// 		fieldName := field.Name
// 		logging.Tracemf(logModGinContext, "compare struct [%s]: field [%s]", valx.Type().Name(), fieldName)

// 		switch fieldName {
// 		case "InstituteUID", "ClassUID", "TeacherUID", "StudentUID", "ParentUID",
// 			"Name", "InstituteName", "ClassName", "FirstName", "LastName",
// 			"Address", "CountryCode", "Location", "MediaLocation",
// 			"DateOfBirth", "PhoneNumber", "Email", "Occupation":
// 			if valx.FieldByName(fieldName).Interface().(string) != valy.FieldByName(fieldName).Interface().(string) {
// 				return fmt.Errorf("[%s] - field (%s) not equal in struct (%s)", serverErrorMessages[seInputSchemaNotValid], fieldName, valx.Type().Name())
// 			}
// 			break
// 		case "PID":
// 			if valx.FieldByName(fieldName).Interface().(int) != valy.FieldByName(fieldName).Interface().(int) {
// 				return fmt.Errorf("[%s] - field (%s) not equal in struct (%s)", serverErrorMessages[seInputSchemaNotValid], fieldName, valx.Type().Name())
// 			}
// 			break
// 		case "InstitutePID", "ClassPID", "TeacherPID", "StudentPID", "ParentPID":
// 			if valx.FieldByName(fieldName).Interface().(int) != valy.FieldByName(fieldName).Interface().(int) {
// 				return fmt.Errorf("[%s] - field (%s) not equal in struct (%s)", serverErrorMessages[seInputSchemaNotValid], fieldName, valx.Type().Name())
// 			}
// 			break
// 		case "Enabled":
// 			if valx.FieldByName(fieldName).Interface().(bool) != valy.FieldByName(fieldName).Interface().(bool) {
// 				return fmt.Errorf("[%s] - field (%s) not equal in struct (%s)", serverErrorMessages[seInputSchemaNotValid], fieldName, valx.Type().Name())
// 			}
// 			break
// 		case "PIDs":
// 			if !isIntListEqual(valx.FieldByName(fieldName).Interface().([]int), valy.FieldByName(fieldName).Interface().([]int)) {
// 				return fmt.Errorf("[%s] - field (%s) not equal in struct (%s)", serverErrorMessages[seInputSchemaNotValid], fieldName, valx.Type().Name())
// 			}
// 			break
// 		}
// 	}
// 	return nil
// }
