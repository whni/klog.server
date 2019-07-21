package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"net/http"
	_ "net/http/pprof"
	"runtime/debug"
	"time"
)

var serverStartTime time.Time

func main() {
	serverStartTime = time.Now()

	// global serverConfig variable
	var scErr error = nil
	serverConfig, scErr = readServerConfig("config.json")
	if scErr != nil {
		panic(scErr)
	}
	initDefaultServerConfig(serverConfig)

	// logging setup
	var loggingErr error = nil
	logging = loggingInitSetup(serverConfig)
	if loggingErr = loggingRegisterModules(logModEnabledTable); loggingErr != nil {
		panic(loggingErr)
	}
	if loggingErr = loggingErrRedirect(errLogFile); loggingErr != nil {
		panic(loggingErr)
	}
	logging.Infomln(logModMain, "Logging module loaded.")

	// gin web framework
	gin.SetMode(gin.DebugMode)
	gin.DefaultWriter = ljGinLogger
	//gin.DefaultErrorWriter = ljErrLogger

	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	// CORS middleware
	r.Use(cors.New(cors.Config{
		AllowOriginFunc: func(origin string) bool {
			return true
		},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	debug.SetTraceback("crash")

	r.POST("/ping", func(c *gin.Context) {
		var response = GinResponse{
			Status: http.StatusOK,
		}
		defer func() {
			ginContextProcessResponse(c, &response)
		}()

		var params = ginContextRequestParameter(c)
		response.Payload = params
		logging.Infoln(params)
	})
	/*
			for api, handler := range Api_maps {
				func_map := handler
				r.POST(api, func_map["Append"])
				r.GET(api, func_map["Get"])
				r.PUT(api, func_map["Change"])
				r.DELETE(api, func_map["Delete"])
			}

			for api, handler := range Api_controller_post {
				r.POST(api, handler)
			}

			for api, handler := range Api_history_post {
				r.POST(api, handler)
			}

			for api, handler := range Api_cwconf_post {
				r.POST(api, handler)
			}

			for api, handler := range Api_status_get {
				r.GET(api, handler)
			}

			for api, handler := range Api_sys_get {
				r.GET(api, handler)
			}

		DbInit()
	*/

	logging.Infomln(logModMain, "Server is listening and serving on 0.0.0.0:8080")
	r.Run(":8080") // listen and serve on 0.0.0.0:8080
	logging.Infomln(logModMain, "Server existed unexpectedly :(")
}
