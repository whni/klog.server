package main

import (
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
	"reflect"
	"strings"
)

var ginAPITable = map[string]map[string]gin.HandlerFunc{
	"/api/0/config/institute": instituteHandlerTable,
}

// GinParameter a generic paramter wrapper for gin web framework handler
type GinParameter struct {
	Token string
	Pkey  string
	Skey  string
	Data  []byte
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

	var primaryKey = ctx.Request.URL.Query().Get("pkey")
	var secondarykey = ctx.Request.URL.Query().Get("skey")

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
	reponseContent := gin.H{}
	if response.Status == http.StatusOK {
		if response.Payload != nil {
			reponseContent["payload"] = response.Payload
		} else {
			reponseContent["payload"] = nil
		}
		if len(response.Message) > 0 {
			reponseContent["message"] = response.Message
		}
	} else {
		reponseContent["message"] = response.Message
	}
	ctx.JSON(response.Status, reponseContent)
}

/* gin input struct check */
func ginInputStructValid(input interface{}) bool {
	val := reflect.ValueOf(input)
	if val.Kind() == reflect.Ptr {
		val = reflect.Indirect(val)
	}

	if val.Kind() != reflect.Struct {
		logging.Errormf(logModGinContext, "unexpected type (%s) - struct required", val.Type().Name())
		return false
	}
	structType := val.Type()

	for i := 0; i < structType.NumField(); i++ {
		field := structType.Field(i)
		fieldName := field.Name

		switch fieldName {
		case "InstituteUID", "ClassUID", "TeacherUID", "StudentUID", "ParentUID",
			"Name", "InstituteName", "ClassName", "FirstName", "LastName",
			"Address", "CountryCode", "Location", "MediaLocation",
			"DateOfBirth", "PhoneNumber", "Email", "Occupation":
			if len(val.FieldByName(fieldName).Interface().(string)) < 1 {
				return false
			}
			break
		case "Password":
			if len(val.FieldByName(fieldName).Interface().(string)) < 5 {
				return false
			}
			break
		case "InstitutePID":
			pid := val.FieldByName(fieldName).Interface().(int)
			if institutes, errCode := findInstitute(pid); errCode > seNoError || len(institutes) == 0 {
				return false
			}
			break
		}
	}
	return true
}

func ginInputStructEqual(x, y interface{}) bool {
	valx := reflect.ValueOf(x)
	if valx.Kind() == reflect.Ptr {
		valx = reflect.Indirect(valx)
	}

	valy := reflect.ValueOf(y)
	if valy.Kind() == reflect.Ptr {
		valy = reflect.Indirect(valy)
	}

	if valx.Kind() != reflect.Struct || valy.Kind() != reflect.Struct {
		logging.Errormf(logModGinContext, "unexpected type (x-%s, y-%s) - struct required", valx.Type().Name(), valy.Type().Name())
		return false
	}
	structTypex := valx.Type()
	structTypey := valy.Type()

	if structTypex.Name() != structTypey.Name() {
		return false
	}

	for i := 0; i < structTypex.NumField(); i++ {
		field := structTypex.Field(i)
		fieldName := field.Name
		switch fieldName {
		case "InstituteUID", "ClassUID", "TeacherUID", "StudentUID", "ParentUID",
			"Name", "InstituteName", "ClassName", "FirstName", "LastName",
			"Address", "CountryCode", "Location", "MediaLocation",
			"DateOfBirth", "PhoneNumber", "Email", "Occupation":
			if valx.FieldByName(fieldName).Interface().(string) != valy.FieldByName(fieldName).Interface().(string) {
				return false
			}
			break
		case "PID", "InstitutePID", "ClassPID", "TeacherPID", "StudentPID", "ParentPID":
			if valx.FieldByName(fieldName).Interface().(int) != valy.FieldByName(fieldName).Interface().(int) {
				return false
			}
			break
		case "Enabled":
			if valx.FieldByName(fieldName).Interface().(bool) != valy.FieldByName(fieldName).Interface().(bool) {
				return false
			}
			break
		case "PIDs":
			if !isIntListEqual(valx.FieldByName(fieldName).Interface().([]int), valy.FieldByName(fieldName).Interface().([]int)) {
				return false
			}
			break
		}
	}
	return true
}
