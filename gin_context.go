package main

import (
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
	"strings"
)

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

var ginAPITable = map[string]map[string]gin.HandlerFunc{
	"/api/0/config/institute": instituteHandlerTable,
}
